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
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"math"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"sync"
	"testing"

	v1 "github.com/dell/goscaleio/types/v1"
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

func TestClientVersion(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(
		func(resp http.ResponseWriter, req *http.Request) {
			if req.RequestURI != "/api/version" {
				t.Fatal("Expecting endpoint /api/version got", req.RequestURI)
			}
			resp.WriteHeader(http.StatusOK)
			resp.Write([]byte(`"2.0"`))
		},
	))
	defer server.Close()
	hostAddr := server.URL
	os.Setenv("GOSCALEIO_ENDPOINT", hostAddr+"/api")
	client, err := NewClient()
	if err != nil {
		t.Fatal(err)
	}
	ver, err := client.GetVersion()
	if err != nil {
		t.Fatal(err)
	}
	if ver != "2.0" {
		t.Fatal("Expecting version string \"2.0\", got ", ver)
	}
}

func TestClientLogin(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(
		func(resp http.ResponseWriter, req *http.Request) {
			switch req.RequestURI {
			case "/api/version":
				resp.WriteHeader(http.StatusOK)
				resp.Write([]byte(`"2.0"`))
			case "/api/login":
				//accept := req.Header.Get("Accept")
				// check Accept header
				//if ver := strings.Split(accept, ";"); len(ver) != 2 {
				//	t.Fatal("Expecting Accept header to include version")
				//} else {
				//	if !strings.HasPrefix(ver[1], "version=") {
				//		t.Fatal("Header Accept must include version")
				//	}
				//}

				uname, pwd, basic := req.BasicAuth()
				if !basic {
					t.Fatal("Client only support basic auth")
				}

				if uname != "ScaleIOUser" || pwd != "password" {
					resp.WriteHeader(http.StatusUnauthorized)
					resp.Write([]byte(`{"message":"Unauthorized","httpStatusCode":401,"errorCode":0}`))
					return
				}
				resp.WriteHeader(http.StatusOK)
				resp.Write([]byte(`"012345678901234567890123456789"`))
			default:
				t.Fatal("Expecting endpoint /api/login got", req.RequestURI)
			}

		},
	))
	defer server.Close()
	hostAddr := server.URL
	os.Setenv("GOSCALEIO_ENDPOINT", hostAddr+"/api")
	client, err := NewClient()
	if err != nil {
		t.Fatal(err)
	}
	// test ok
	_, err = client.Authenticate(&ConfigConnect{
		Username: "ScaleIOUser",
		Password: "password",
		Endpoint: "",
		Version:  "2.0",
	})
	if err != nil {
		t.Fatal(err)
	}
	if client.GetToken() != "012345678901234567890123456789" {
		t.Fatal("Expecting token 012345678901234567890123456789, got", client.GetToken())
	}

	// test bad login
	_, err = client.Authenticate(&ConfigConnect{
		Username: "ScaleIOUser",
		Password: "badPassWord",
		Endpoint: "",
		Version:  "2.0",
	})
	if err == nil {
		t.Fatal("Expecting an error for bad Login, but did not")
	}
}

type stubTypeWithMetaData struct{}

func (s stubTypeWithMetaData) MetaData() http.Header {
	h := make(http.Header)
	h.Set("foo", "bar")
	return h
}

func Test_addMetaData(t *testing.T) {
	var tests = []struct {
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

func Test_updateHeaders(t *testing.T) {
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

func Test_getJSONWithRetry(t *testing.T) {
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
	b, err := ioutil.ReadAll(rc)
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
