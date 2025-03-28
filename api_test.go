// Copyright © 2019 - 2022 Dell Inc. or its subsidiaries. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//      http://www.apache.org/licenses/LICENSE-2.0
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package goscaleio

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"strings"
	"sync"
	"testing"
	"time"

	v1 "github.com/dell/goscaleio/types/v1"
	"github.com/stretchr/testify/assert"
)

func setupClient(t *testing.T, hostAddr string) *Client {
	os.Setenv("GOSCALEIO_ENDPOINT", hostAddr+"/api")
	client, err := NewClient()
	if err != nil {
		t.Fatal(err)
	}
	// test ok
	_, err = client.Authenticate(&ConfigConnect{
		Username: "ScaleIOUser",
		Password: "password",
		Version:  "2.0",
	})
	if err != nil {
		t.Fatal(err)
	}
	return client
}

func requestAuthOK(resp http.ResponseWriter, req *http.Request) bool {
	_, pwd, _ := req.BasicAuth()
	if pwd == "" {
		resp.WriteHeader(http.StatusUnauthorized)
		resp.Write([]byte(`{"message":"Unauthorized","httpStatusCode":401,"errorCode":0}`))
		return false
	}
	return true
}

func handleAuthToken(resp http.ResponseWriter, req *http.Request) {
	if !requestAuthOK(resp, req) {
		return
	}
	resp.WriteHeader(http.StatusOK)
	resp.Write([]byte(`"012345678901234567890123456789"`))
}

type goScaleioTestServer struct {
	*httptest.Server
	// Server will return this as http code if injectErrorNum > 0
	injectError int
	// The number of times to inject the error
	injectErrorNum int
	// If not "-", server will use it as resp body once only
	injectResp string
}

func newGoScaleioTestServer() *goScaleioTestServer {
	ts := &goScaleioTestServer{
		injectError: 0,
		injectResp:  "-",
	}
	ts.Server = httptest.NewServer(http.HandlerFunc(
		func(resp http.ResponseWriter, req *http.Request) {
			switch req.RequestURI {
			case "/api/version":
				if ts.injectErrorNum > 0 {
					resp.WriteHeader(ts.injectError)
					ts.injectErrorNum--
					return
				}
				// Check for valid authentication token
				authHeader := req.Header.Get("Authorization")
				if authHeader != "Bearer valid_token" {
					// Respond with 401 Unauthorized only if the token is missing or invalid
					if authHeader == "" {
						resp.WriteHeader(http.StatusUnauthorized)
						resp.Write([]byte(`Unauthorized 401`))
						return
					}
					// For any other case, respond with version
					resp.WriteHeader(http.StatusOK)
					if ts.injectResp != "-" {
						resp.Write([]byte(ts.injectResp))
						ts.injectResp = "-"
					} else {
						resp.Write([]byte(`"4.0"`))
					}
					return
				}
				// Respond with version 4.0
				resp.WriteHeader(http.StatusOK)
				if ts.injectResp != "-" {
					resp.Write([]byte(ts.injectResp))
					ts.injectResp = "-"
				} else {
					resp.Write([]byte(`"4.0"`))
				}
			case "/api/login":
				// Check basic authentication
				uname, pwd, basic := req.BasicAuth()
				if !basic {
					// Respond with 401 Unauthorized if basic auth is not provided
					resp.WriteHeader(http.StatusUnauthorized)
					resp.Write([]byte(`{"message":"Unauthorized","httpStatusCode":401,"errorCode":0}`))
					return
				}

				if uname != "ScaleIOUser" || pwd != "password" {
					// Respond with 401 Unauthorized if credentials are invalid
					resp.WriteHeader(http.StatusUnauthorized)
					resp.Write([]byte(`{"message":"Unauthorized","httpStatusCode":401,"errorCode":0}`))
					return
				}
				// Respond with a valid token
				resp.WriteHeader(http.StatusOK)
				if ts.injectResp != "-" {
					resp.Write([]byte(ts.injectResp))
					ts.injectResp = "-"
				} else {
					resp.Write([]byte(`"012345678901234567890123456789"`))
				}
			default:
				// Respond with 404 Not Found for any other endpoint
				http.Error(resp, "Expecting endpoint /api/login got "+req.RequestURI, http.StatusNotFound)
			}
		},
	))

	return ts
}

func TestClientVersion(t *testing.T) {
	server := newGoScaleioTestServer()
	defer server.Close()

	// Set the environment variable for the endpoint
	os.Setenv("GOSCALEIO_ENDPOINT", server.URL+"/api")

	t.Run("successful authentication", func(t *testing.T) {
		client, err := NewClient()
		assert.NoError(t, err)

		// Test successful authentication
		_, err = client.Authenticate(&ConfigConnect{
			Username: "ScaleIOUser",
			Password: "password",
			Version:  "4.0",
		})
		assert.NoError(t, err)

		ver, err := client.GetVersion()
		assert.NoError(t, err)
		assert.Equal(t, "4.0", ver)
	})

	t.Run("auto re-authentication", func(t *testing.T) {
		client, err := NewClient()
		assert.NoError(t, err)

		_, err = client.Authenticate(&ConfigConnect{
			Username: "ScaleIOUser",
			Password: "password",
			Version:  "4.0",
		})
		assert.NoError(t, err)

		server.injectError = http.StatusUnauthorized
		server.injectErrorNum = 1

		// Expect to get StatusUnauthorized, then auto-reauthenticate and retry get version
		ver, err := client.GetVersion()
		assert.NoError(t, err)
		assert.Equal(t, "4.0", ver)
	})

	t.Run("server returns error status", func(t *testing.T) {
		client, err := NewClient()
		assert.NoError(t, err)

		_, err = client.Authenticate(&ConfigConnect{
			Username: "ScaleIOUser",
			Password: "password",
			Version:  "4.0",
		})
		assert.NoError(t, err)

		server.injectError = http.StatusBadRequest
		server.injectErrorNum = 2

		// Expect to get StatusUnauthorized, then auto-reauthenticate and retry get version
		ver, err := client.GetVersion()
		assert.Error(t, err)
		assert.Equal(t, "", ver)

		err = client.updateVersion()
		assert.Error(t, err)
	})

	t.Run("request after failed authentication", func(t *testing.T) {
		// Initialize the client
		client, err := NewClient()
		assert.NoError(t, err)

		// Test unauthorized authentication
		_, err = client.Authenticate(&ConfigConnect{
			Username: "ScaleIOUser",
			Password: "badpassword",
		})
		assert.Error(t, err)

		// Try to get version after failed authentication
		_, err = client.GetVersion()
		assert.ErrorContains(t, err, "Unauthorized")

		// Re-authenticate properly this time
		_, err = client.Authenticate(&ConfigConnect{
			Username: "ScaleIOUser",
			Password: "password",
		})
		assert.NoError(t, err)

		ver, err := client.GetVersion()
		assert.NoError(t, err)
		assert.Equal(t, "4.0", ver)
	})

	t.Run("authenticate with malformed version", func(t *testing.T) {
		os.Setenv("GOSCALEIO_VERSION", "malformed")
		defer os.Unsetenv("GOSCALEIO_VERSION")

		// Initialize the client
		client, err := NewClient()
		assert.NoError(t, err)

		// Test client configured with bad version
		_, err = client.Authenticate(&ConfigConnect{
			Username: "ScaleIOUser",
			Password: "password",
		})
		assert.Error(t, err)
	})

	t.Run("get version with malformed version", func(t *testing.T) {
		os.Setenv("GOSCALEIO_VERSION", "malformed")
		defer os.Unsetenv("GOSCALEIO_VERSION")

		// Initialize the client
		client, err := NewClient()
		assert.NoError(t, err)

		// Test client configured with bad version
		_, err = client.GetVersion()
		assert.Error(t, err)
	})

	t.Run("server returns malformed version", func(t *testing.T) {
		// Initialize the client
		client, err := NewClient()
		assert.NoError(t, err)

		_, err = client.Authenticate(&ConfigConnect{
			Username: "ScaleIOUser",
			Password: "password",
		})
		assert.NoError(t, err)

		server.injectResp = "malformed"

		ver, err := client.GetVersion()
		assert.NoError(t, err)
		assert.Equal(t, "malformed", ver)
	})
}

func TestClientLogin(t *testing.T) {
	server := newGoScaleioTestServer()
	defer server.Close()

	os.Setenv("GOSCALEIO_ENDPOINT", server.URL+"/api")

	t.Run("successful login", func(t *testing.T) {
		client, err := NewClient()
		assert.NoError(t, err)

		_, err = client.Authenticate(&ConfigConnect{
			Username: "ScaleIOUser",
			Password: "password",
		})
		assert.NoError(t, err)
		assert.Equal(t, "012345678901234567890123456789", client.GetToken())

		client.SetToken("012345678901234567890123456789")

		cc := client.GetConfigConnect()
		assert.NotNil(t, cc)
		assert.Equal(t, "ScaleIOUser", cc.Username)
		assert.Equal(t, "4.0", cc.Version)
	})

	t.Run("failed login", func(t *testing.T) {
		client, err := NewClient()
		assert.NoError(t, err)

		_, err = client.Authenticate(&ConfigConnect{
			Username: "ScaleIOUser",
			Password: "badPassWord",
		})
		assert.Error(t, err)
	})
}

type stubTypeWithMetaData struct{}

func (s stubTypeWithMetaData) MetaData() http.Header {
	h := make(http.Header)
	h.Set("foo", "bar")
	return h
}

func Test_addMetaData(t *testing.T) {
	tests := []struct {
		name           string
		givenHeader    map[string]string
		expectedHeader map[string]string
		body           interface{}
	}{
		{"nil header is a noop", nil, nil, nil},
		{"nil body is a noop", nil, nil, nil},
		{"header is updated", make(map[string]string), map[string]string{"Foo": "bar"}, stubTypeWithMetaData{}},
		{"header is not updated", make(map[string]string), map[string]string{}, struct{}{}},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			addMetaData(tt.givenHeader, tt.body)

			switch {
			case tt.givenHeader == nil:
				if tt.givenHeader != nil {
					t.Errorf("(%s): expected %s, actual %s", tt.body, tt.expectedHeader, tt.givenHeader)
				}
			case tt.body == nil:
				if len(tt.givenHeader) != 0 {
					t.Errorf("(%s): expected %s, actual %s", tt.body, tt.expectedHeader, tt.givenHeader)
				}
			default:
				if !reflect.DeepEqual(tt.expectedHeader, tt.givenHeader) {
					t.Errorf("(%s): expected %s, actual %s", tt.body, tt.expectedHeader, tt.givenHeader)
				}
			}
		})
	}
}

func Test_updateHeaders(_ *testing.T) {
	var wg sync.WaitGroup
	for i := 0; i < 3; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			updateHeaders("3.5")
		}()
	}
	wg.Wait()
}

func TestGetJSONWithRetry(t *testing.T) {
	t.Run("retried request is similar to the original", func(t *testing.T) {
		var (
			paths     []string      // record the requested paths in order.
			bodies    []string      // record the request bodies in order.
			headers   []http.Header // record the headers in order.
			callCount int           // how many times our endpoint was requested.
		)
		checkHeaders := []string{"Accept"} // only check these headers.

		// mock a PowerFlex endpoint.
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Record the requested paths in order.
			paths = append(paths, fmt.Sprintf("%s %s", r.Method, r.URL.Path))

			switch r.URL.Path {
			case "/testing":
				callCount++
				b := testReadAll(t, r.Body)
				bodies = append(bodies, string(b))
				headers = append(headers, testFilterHeaders(r.Header, checkHeaders))
				// First request to error with HTTP 401  and trigger a login request.
				if callCount == 1 {
					w.WriteHeader(http.StatusUnauthorized)
					testjsonEncode(t, w, testBuildError(http.StatusUnauthorized))
				}
			case "/api/login":
				fmt.Fprintf(w, `"fakesessiontoken"`)
			default:
				t.Fatalf("unexpected path: %q", r.URL.Path)
			}
		}))
		defer ts.Close()
		c, err := NewClientWithArgs(ts.URL, "3.5", math.MaxInt64, true, false)
		if err != nil {
			t.Fatal(err)
		}

		// Call getJSONWithRetry with a dummy request and some
		// map as the request body. We don't care about the
		// response so pass in nil.
		m := map[string]string{"foo": "bar"}
		wantBody, err := json.Marshal(&m)
		if err != nil {
			t.Fatal(err)
		}
		c.getJSONWithRetry(http.MethodPost, "/testing", wantBody, nil)

		// Assert the call order was as expected.
		wantPaths := []string{"POST /testing", "GET /api/login", "POST /testing"}
		if !reflect.DeepEqual(paths, wantPaths) {
			t.Errorf("paths: got %+v, want %+v", paths, wantPaths)
		}
		// Assert the second body was the same as the first.
		gotBodies, wantBodies := bodies[1], bodies[0]
		if !reflect.DeepEqual(gotBodies, wantBodies) {
			t.Errorf("retried body: got %q, want %q", gotBodies, wantBodies)
		}
		// Assert the headers for both requests were the same.
		gotHeaders, wantHeaders := headers[1], headers[0]
		if !reflect.DeepEqual(gotHeaders, wantHeaders) {
			t.Errorf("retried headers: got %q, want %q", gotHeaders, wantHeaders)
		}
	})
}

// testFilterHeaders accepts a header and a list of header names
// to filter on (inclusive).  The returned http.Header will include only
// header fields with these names.
func testFilterHeaders(h http.Header, filter []string) http.Header {
	result := make(http.Header)
	for _, v := range filter {
		if _, ok := h[v]; !ok {
			continue
		}
		result.Set(v, h.Get(v))
	}
	return result
}

func testBuildError(code int) error {
	return &v1.Error{
		Message:        "test message",
		HTTPStatusCode: code,
		ErrorCode:      0,
		ErrorDetails:   nil,
	}
}

func testReadAll(t *testing.T, rc io.ReadCloser) []byte {
	t.Helper()
	b, err := io.ReadAll(rc)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		rc.Close()
	})
	return b
}

func testjsonEncode(t *testing.T, w io.Writer, v interface{}) {
	t.Helper()
	err := json.NewEncoder(w).Encode(&v)
	if err != nil {
		t.Fatal(err)
	}
}

func TestWithFields(t *testing.T) {
	tests := []struct {
		name   string
		fields map[string]interface{}
	}{
		{
			name:   "No fields",
			fields: nil,
		},
		{
			name: "With fields",
			fields: map[string]interface{}{
				"key1": "test1",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := withFields(tt.fields, "")
			if err == nil {
				t.Errorf("withFieldsE() expected error got nil")
			}
		})
	}
}

func TestNewClientWithArgs(t *testing.T) {
	tests := []struct {
		name     string
		endpoint string
		setEnv   func()
		wantErr  bool
	}{
		{
			name:     "success",
			endpoint: "/testing",
			wantErr:  false,
		},
		{
			name:     "failure",
			endpoint: "",
			setEnv: func() {
				os.Setenv("GOSCALEIO_SHOWHTTP", "true")
			},
			wantErr: true,
		},
		{
			name:     "Set GOSCALEIO_SHOWHTTP to true",
			endpoint: "/testing",
			setEnv: func() {
				showHTTP = true
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setEnv != nil {
				tt.setEnv()
			}
			_, err := NewClientWithArgs(tt.endpoint, "3.5", math.MaxInt64, true, false)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewClientWithArgs() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestWithFieldsE(t *testing.T) {
	tests := []struct {
		name   string
		fields map[string]interface{}
		inner  error
	}{
		{
			name:   "No fields",
			fields: nil,
			inner:  fmt.Errorf("test"),
		},
		{
			name: "With fields",
			fields: map[string]interface{}{
				"key1": "test1",
			},
			inner: fmt.Errorf("test"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := withFieldsE(tt.fields, "", tt.inner)
			if err == nil {
				t.Errorf("withFieldsE() expected error got nil")
			}
		})
	}
}

func TestGetStringWithRetry(t *testing.T) {
	tests := []struct {
		name string
		URL  string
	}{
		{
			name: "Unauthorized, Need re-auth",
			URL:  "/testing",
		},
		{
			name: "Re-authentication not required",
			URL:  "/success",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				switch r.URL.Path {
				case "/success":
					w.WriteHeader(http.StatusOK)
				case "/testing":
					w.WriteHeader(http.StatusUnauthorized)
					testjsonEncode(t, w, testBuildError(http.StatusUnauthorized))
				case "/api/login":
					_, err := fmt.Fprintf(w, `"fakesessiontoken"`)
					if err != nil {
						return
					}
				default:
					t.Fatalf("unexpected path: %q", r.URL.Path)
				}
			}))
			defer ts.Close()
			client, err := NewClientWithArgs(ts.URL, "3.5", math.MaxInt64, true, false)
			if err != nil {
				t.Fatal(err)
			}
			_, err = client.getStringWithRetry(http.MethodPost, tt.URL, nil)
		})
	}
}

func TestWithContext(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
		switch req.RequestURI {
		case "/api/login":
			resp.WriteHeader(http.StatusOK)
		default:
			resp.WriteHeader(http.StatusBadRequest)
			resp.Write([]byte(`{"message":"no route handled","httpStatusCode":400,"errorCode":0}`))
		}
	}))

	defer server.Close()

	client, err := NewClientWithArgs(server.URL, "3.6", math.MaxInt64, true, false)
	if err != nil {
		t.Fatal(err)
	}

	parentCtx := context.Background()
	ctx, cancel := context.WithTimeout(parentCtx, 10*time.Millisecond)
	defer cancel()

	// Test with Context
	_, err = client.WithContext(ctx).Authenticate(&ConfigConnect{
		Username: "ScaleIOUser",
		Password: "password",
		Endpoint: "",
		Version:  "2.0",
	})
	if err != nil {
		t.Fatal(err)
	}

	// Delay to test calls without context paramter set.
	time.Sleep(1 * time.Second)

	// Test with without context
	_, err = client.Authenticate(&ConfigConnect{
		Username: "ScaleIOUser",
		Password: "password",
		Endpoint: "",
		Version:  "2.0",
	})
	if err != nil {
		t.Fatal(err)
	}
}

type failingReadCloser struct{}

func (r *failingReadCloser) Read(_ []byte) (n int, err error) {
	return 0, io.ErrUnexpectedEOF
}

func (r *failingReadCloser) Close() error {
	return io.ErrUnexpectedEOF
}

func TestExtractString(t *testing.T) {
	res := &http.Response{
		Body: &failingReadCloser{},
	}

	// Body read error
	_, err := extractString(res)
	assert.Error(t, err)

	// Successful extraction
	res = &http.Response{
		Body: io.NopCloser(strings.NewReader(`{"message":"success"}`)),
	}
	s, err := extractString(res)
	assert.NoError(t, err)
	assert.Equal(t, `{"message":"success"}`, s)
}
