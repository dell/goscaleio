/*
 *
 * Copyright © 2021-2024 Dell Inc. or its subsidiaries. All Rights Reserved.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */

// Copyright © 2024 Dell Inc. or its subsidiaries. All Rights Reserved.
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

package api

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
	"time"

	types "github.com/dell/goscaleio/types/v1"
	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	tests := []struct {
		name        string
		host        string
		opts        ClientOptions
		debug       bool
		expectedErr error
	}{
		{
			name:        "missing host",
			host:        "",
			opts:        ClientOptions{},
			expectedErr: errors.New("missing endpoint"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := New(context.Background(), tt.host, tt.opts, true)
			if !reflect.DeepEqual(err, tt.expectedErr) {
				t.Errorf("Got error: %v, expected: %v", err, tt.expectedErr)
				return
			}
		})
	}
}

func TestGet(t *testing.T) {
	tests := map[string]struct {
		method       string
		path         string
		headers      map[string]string
		body         interface{}
		expectedErr  error
		expectedBody string
	}{
		"successful GET request": {
			method:       http.MethodGet,
			path:         "/api/test",
			headers:      map[string]string{"Content-Type": "application/json"},
			body:         nil,
			expectedErr:  nil,
			expectedBody: `{"message":"success"}`,
		},
		"failed GET request": {
			method:       http.MethodGet,
			path:         "/api/test",
			headers:      map[string]string{"Content-Type": "application/json"},
			body:         nil,
			expectedErr:  fmt.Errorf("400 Bad Request"),
			expectedBody: "",
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			// Create a test server to handle the request
			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				for k, v := range tt.headers {
					w.Header().Set(k, v)
				}

				if tt.expectedErr != nil {
					w.WriteHeader(http.StatusBadRequest)
					data, err := json.Marshal(tt.expectedErr)
					if err != nil {
						t.Fatal(err)
					}
					w.Write(data)
				} else {
					w.WriteHeader(http.StatusOK)
					w.Write([]byte(tt.expectedBody))
				}
			}))
			defer ts.Close()

			// Create a new client and set the host to the test server
			c, err := New(context.Background(), ts.URL, ClientOptions{
				Timeout:  10 * time.Second,
				Insecure: true,
				ShowHTTP: true,
			}, true)
			if err != nil {
				t.Fatal(err)
			}
			c.SetToken("token-1")

			if c.GetToken() != "token-1" {
				t.Errorf("expected token %s, got %s", "token-1", c.GetToken())
			}

			// Call the Get function with the test parameters
			var resp interface{}
			err = c.Get(context.Background(), tt.path, tt.headers, &resp)

			// Check if the error matches the expected error
			if tt.expectedErr != nil {
				if err == nil || err.Error() != tt.expectedErr.Error() {
					t.Errorf("expected error %v, got %v", tt.expectedErr, err)
				}
			} else if err != nil {
				t.Errorf("unexpected error: %v", err)
			}

			// Check if the response body matches the expected body
			if tt.expectedBody != "" {
				b, err := json.Marshal(resp)
				if err != nil {
					t.Fatal(err)
				}
				if string(b) != tt.expectedBody {
					t.Errorf("expected response body %s, got %s", tt.expectedBody, string(b))
				}
			}
		})
	}
}

func TestPost(t *testing.T) {
	tests := map[string]struct {
		method       string
		path         string
		headers      map[string]string
		body         interface{}
		expectedErr  error
		expectedBody string
	}{
		"successful POST request": {
			method:       http.MethodPost,
			path:         "/api/test",
			headers:      map[string]string{"Content-Type": "application/json"},
			body:         nil,
			expectedErr:  nil,
			expectedBody: `{"message":"success"}`,
		},
		"failed POST request": {
			method:       http.MethodPost,
			path:         "/api/test",
			headers:      map[string]string{"Content-Type": "application/json"},
			body:         nil,
			expectedErr:  fmt.Errorf("400 Bad Request"),
			expectedBody: "",
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			// Create a test server to handle the request
			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method != tt.method {
					t.Fatal(fmt.Errorf("wrong method. Expected %s; but got %s", tt.method, r.Method))
				}

				if r.URL.Path != tt.path {
					t.Fatal(fmt.Errorf("wrong path. Expected %s; but got %s", tt.path, r.URL.Path))
				}

				for k, v := range tt.headers {
					if r.Header.Get(k) != v {
						t.Fatal(fmt.Errorf("wrong header value. Expected %s; but got %s", v, r.Header.Get(k)))
					}
				}

				if tt.expectedErr != nil {
					w.WriteHeader(http.StatusBadRequest)
					data, err := json.Marshal(tt.expectedErr)
					if err != nil {
						t.Fatal(err)
					}
					w.Write(data)
				} else {
					w.WriteHeader(http.StatusOK)
					w.Write([]byte(tt.expectedBody))
				}
			}))
			defer ts.Close()

			// Create a new client and set the host to the test server
			c, err := New(context.Background(), ts.URL, ClientOptions{}, false)
			if err != nil {
				t.Fatal(err)
			}

			// Call the Post function with the test parameters
			var resp interface{}
			err = c.Post(context.Background(), tt.path, tt.headers, tt.body, &resp)

			// Check if the error matches the expected error
			if tt.expectedErr != nil {
				if err == nil || err.Error() != tt.expectedErr.Error() {
					t.Errorf("expected error %v, got %v", tt.expectedErr, err)
				}
			} else if err != nil {
				t.Errorf("unexpected error: %v", err)
			}

			// Check if the response body matches the expected body
			if tt.expectedBody != "" {
				b, err := json.Marshal(resp)
				if err != nil {
					t.Fatal(err)
				}
				if string(b) != tt.expectedBody {
					t.Errorf("expected response body %s, got %s", tt.expectedBody, string(b))
				}
			}
		})
	}
}

func TestPut(t *testing.T) {
	tests := map[string]struct {
		method       string
		path         string
		headers      map[string]string
		body         interface{}
		expectedErr  error
		expectedBody string
	}{
		"successful PUT request": {
			method:       http.MethodPut,
			path:         "/api/test",
			headers:      map[string]string{"Content-Type": "application/json"},
			body:         nil,
			expectedErr:  nil,
			expectedBody: `{"message":"success"}`,
		},
		"failed PUT request": {
			method:       http.MethodPut,
			path:         "/api/test",
			headers:      map[string]string{"Content-Type": "application/json"},
			body:         nil,
			expectedErr:  fmt.Errorf("400 Bad Request"),
			expectedBody: "",
		},
		"invalid headers and body": {
			method:       http.MethodPut,
			path:         "/api/test",
			headers:      nil,
			body:         "nil",
			expectedErr:  fmt.Errorf("400 Bad Request"),
			expectedBody: "",
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			// Create a test server to handle the request
			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method != tt.method {
					t.Fatal(fmt.Errorf("wrong method. Expected %s; but got %s", tt.method, r.Method))
				}

				if r.URL.Path != tt.path {
					t.Fatal(fmt.Errorf("wrong path. Expected %s; but got %s", tt.path, r.URL.Path))
				}

				for k, v := range tt.headers {
					if r.Header.Get(k) != v {
						t.Fatal(fmt.Errorf("wrong header value. Expected %s; but got %s", v, r.Header.Get(k)))
					}
				}

				if tt.expectedErr != nil {
					w.WriteHeader(http.StatusBadRequest)
					data, err := json.Marshal(tt.expectedErr)
					if err != nil {
						t.Fatal(err)
					}
					w.Write(data)
				} else {
					w.WriteHeader(http.StatusOK)
					w.Write([]byte(tt.expectedBody))
				}
			}))
			defer ts.Close()

			// Create a new client and set the host to the test server
			c, err := New(context.Background(), ts.URL, ClientOptions{}, false)
			if err != nil {
				t.Fatal(err)
			}

			// Call the Put function with the test parameters
			var resp interface{}
			err = c.Put(context.Background(), tt.path, tt.headers, tt.body, &resp)

			// Check if the error matches the expected error
			if tt.expectedErr != nil {
				if err == nil || err.Error() != tt.expectedErr.Error() {
					t.Errorf("expected error %v, got %v", tt.expectedErr, err)
				}
			} else if err != nil {
				t.Errorf("unexpected error: %v", err)
			}

			// Check if the response body matches the expected body
			if tt.expectedBody != "" {
				b, err := json.Marshal(resp)
				if err != nil {
					t.Fatal(err)
				}
				if string(b) != tt.expectedBody {
					t.Errorf("expected response body %s, got %s", tt.expectedBody, string(b))
				}
			}
		})
	}
}

func TestDelete(t *testing.T) {
	tests := map[string]struct {
		method      string
		path        string
		headers     map[string]string
		body        interface{}
		expectedErr error
	}{
		"successful DELETE request": {
			method:      http.MethodDelete,
			path:        "/api/test",
			headers:     map[string]string{"Content-Type": "application/json"},
			body:        nil,
			expectedErr: nil,
		},
		"failed DELETE request": {
			method:      http.MethodDelete,
			path:        "/api/test",
			headers:     map[string]string{"Content-Type": "application/json"},
			body:        nil,
			expectedErr: fmt.Errorf("400 Bad Request"),
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			// Create a test server to handle the request
			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method != tt.method {
					t.Fatal(fmt.Errorf("wrong method. Expected %s; but got %s", tt.method, r.Method))
				}

				if r.URL.Path != tt.path {
					t.Fatal(fmt.Errorf("wrong path. Expected %s; but got %s", tt.path, r.URL.Path))
				}

				for k, v := range tt.headers {
					if r.Header.Get(k) != v {
						t.Fatal(fmt.Errorf("wrong header value. Expected %s; but got %s", v, r.Header.Get(k)))
					}
				}

				if tt.expectedErr != nil {
					w.WriteHeader(http.StatusBadRequest)
					data, err := json.Marshal(tt.expectedErr)
					if err != nil {
						t.Fatal(err)
					}
					w.Write(data)
				} else {
					w.WriteHeader(http.StatusOK)
				}
			}))
			defer ts.Close()

			// Create a new client and set the host to the test server
			c, err := New(context.Background(), ts.URL, ClientOptions{}, false)
			if err != nil {
				t.Fatal(err)
			}

			// Call the Delete function with the test parameters
			err = c.Delete(context.Background(), tt.path, tt.headers, tt.body)

			// Check if the error matches the expected error
			if tt.expectedErr != nil {
				if err == nil || err.Error() != tt.expectedErr.Error() {
					t.Errorf("expected error %v, got %v", tt.expectedErr, err)
				}
			} else if err != nil {
				t.Errorf("unexpected error: %v", err)
			}
		})
	}
}

func TestDo(t *testing.T) {
	tests := map[string]struct {
		method       string
		path         string
		body         interface{}
		expectedErr  error
		expectedBody string
	}{
		"successful GET request": {
			method:       http.MethodGet,
			path:         "/api/test",
			body:         nil,
			expectedErr:  nil,
			expectedBody: `{"message":"success"}`,
		},
		"failed GET request": {
			method:       http.MethodGet,
			path:         "/api/test",
			body:         nil,
			expectedErr:  fmt.Errorf("400 Bad Request"),
			expectedBody: "",
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			// Create a test server to handle the request
			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method != tt.method {
					t.Fatal(fmt.Errorf("wrong method. Expected %s; but got %s", tt.method, r.Method))
				}

				if r.URL.Path != tt.path {
					t.Fatal(fmt.Errorf("wrong path. Expected %s; but got %s", tt.path, r.URL.Path))
				}

				if tt.expectedErr != nil {
					w.WriteHeader(http.StatusBadRequest)
					data, err := json.Marshal(tt.expectedErr)
					if err != nil {
						t.Fatal(err)
					}
					w.Write(data)
				} else {
					w.WriteHeader(http.StatusOK)
					w.Write([]byte(tt.expectedBody))
				}
			}))
			defer ts.Close()

			// Create a new client and set the host to the test server
			c, err := New(context.Background(), ts.URL, ClientOptions{}, false)
			if err != nil {
				t.Fatal(err)
			}

			// Call the Do function with the test parameters
			var resp interface{}
			err = c.Do(context.Background(), tt.method, tt.path, tt.body, &resp)

			// Check if the error matches the expected error
			if tt.expectedErr != nil {
				if err == nil || err.Error() != tt.expectedErr.Error() {
					t.Errorf("expected error %v, got %v", tt.expectedErr, err)
				}
			} else if err != nil {
				t.Errorf("unexpected error: %v", err)
			}

			// Check if the response body matches the expected body
			if tt.expectedBody != "" {
				b, err := json.Marshal(resp)
				if err != nil {
					t.Fatal(err)
				}
				if string(b) != tt.expectedBody {
					t.Errorf("expected response body %s, got %s", tt.expectedBody, string(b))
				}
			}
		})
	}
}

func TestDoXMLRequest_FailDecode(t *testing.T) {
	tests := map[string]struct {
		method       string
		path         string
		body         interface{}
		resp         interface{}
		req          string
		version      string
		expectedErr  error
		expectedBody string
	}{
		"Fail decode response": {
			method:       http.MethodPost,
			path:         "/api/test",
			body:         nil,
			resp:         "",
			version:      "3.6",
			expectedBody: `{"message":"failed to decode response"}`,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			// Create a test server to handle the request
			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method != tt.method {
					t.Fatal(fmt.Errorf("wrong method. Expected %s; but got %s", tt.method, r.Method))
				}

				if tt.req != "" {
					w.WriteHeader(http.StatusBadRequest)
					w.Write([]byte(tt.req))
				} else {
					w.WriteHeader(http.StatusOK)
					w.Write([]byte(tt.expectedBody))
				}
			}))
			defer ts.Close()

			// Create a new client and set the host to the test server
			c, err := New(context.Background(), ts.URL, ClientOptions{}, false)
			if err != nil {
				t.Fatal(err)
			}

			// Call the DoWithHeaders function with the test parameters
			_, err = c.DoXMLRequest(
				context.Background(),
				tt.method,
				tt.path,
				tt.version,
				tt.body,
				tt.resp,
			)
			assert.NotNil(t, err)
		})
	}
}

func TestDoXMLRequest(t *testing.T) {
	tests := map[string]struct {
		method       string
		path         string
		body         interface{}
		resp         interface{}
		req          string
		version      string
		token        string
		expectedErr  error
		expectedBody string
	}{
		"successful GET request": {
			method:       http.MethodGet,
			path:         "/api/test",
			body:         nil,
			resp:         "",
			expectedErr:  nil,
			expectedBody: `{"message":"success"}`,
		},
		"successful Post request Body not Nil": {
			method: http.MethodPost,
			path:   "/api/test",
			body: types.SwitchCredentialWrapper{
				IomCredential: types.IomCredential{
					Username: "test",
					Password: "test",
				},
			},
			resp:         "",
			expectedErr:  nil,
			expectedBody: `{"message":"success"}`,
		},
		"successful Post request Body without begining slash": {
			method: http.MethodPost,
			path:   "api/test",
			body: types.SwitchCredentialWrapper{
				IomCredential: types.IomCredential{
					Username: "test",
					Password: "test",
				},
			},
			resp:         "",
			expectedErr:  nil,
			expectedBody: `{"message":"success"}`,
		},
		"Fail Invalid Body": {
			method:       http.MethodPost,
			path:         "api/test",
			body:         map[string]string{"test": "test"},
			resp:         "",
			expectedErr:  fmt.Errorf("xml: unsupported type: map[string]string"),
			expectedBody: "",
		},
		"Invalid response": {
			method: http.MethodPost,
			path:   "/api/test/notexist",
			body: types.SwitchCredentialWrapper{
				IomCredential: types.IomCredential{
					Username: "test",
					Password: "test",
				},
			},
			resp:         "",
			req:          "{}",
			expectedErr:  fmt.Errorf(`400 Bad Request`),
			expectedBody: "",
		},
		"Invalid status code": {
			method: http.MethodGet,
			path:   "api/test",
			body: types.SwitchCredentialWrapper{
				IomCredential: types.IomCredential{
					Username: "test",
					Password: "test",
				},
			},
			resp:         nil,
			expectedBody: "",
		},
		"Fail parse version": {
			method:      http.MethodPost,
			path:        "/api/test",
			body:        nil,
			resp:        "",
			version:     "invalid_version",
			expectedErr: fmt.Errorf("strconv.ParseFloat: parsing \"invalid_version\": invalid syntax"),
		},
		"Token with version 4": {
			method:       http.MethodGet,
			path:         "/api/test",
			body:         nil,
			resp:         "",
			version:      "4.0",
			token:        "test-token",
			expectedErr:  nil,
			expectedBody: `{"message":"success"}`,
		},
		"Token with version less than 4": {
			method:       http.MethodGet,
			path:         "/api/test",
			body:         nil,
			resp:         "",
			version:      "2.0",
			token:        "test-token",
			expectedErr:  nil,
			expectedBody: `{"message":"success"}`,
		},
		"Token with no version": {
			method:       http.MethodGet,
			path:         "/api/test",
			body:         nil,
			resp:         "",
			version:      "",
			token:        "test-token",
			expectedErr:  nil,
			expectedBody: `{"message":"success"}`,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			// Create a test server to handle the request
			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method != tt.method {
					t.Fatal(fmt.Errorf("wrong method. Expected %s; but got %s", tt.method, r.Method))
				}

				if tt.req != "" {
					w.WriteHeader(http.StatusBadRequest)
					w.Write([]byte(tt.req))
				} else {
					w.WriteHeader(http.StatusOK)
					w.Write([]byte(tt.expectedBody))
				}
			}))
			defer ts.Close()

			// Create a new client and set the host to the test server
			c, err := New(context.Background(), ts.URL, ClientOptions{}, false)
			if err != nil {
				t.Fatal(err)
			}

			c.SetToken(tt.token)

			// Call the DoWithHeaders function with the test parameters
			var resp interface{}
			res := tt.resp
			if tt.resp != nil {
				res = &resp
			}
			_, err = c.DoXMLRequest(
				context.Background(),
				tt.method,
				tt.path,
				tt.version,
				tt.body,
				res,
			)

			// Check if the error matches the expected error
			if tt.expectedErr != nil {
				if err == nil || err.Error() != tt.expectedErr.Error() {
					t.Errorf("expected error %v, got %v", tt.expectedErr, err)
				}
			} else if err != nil {
				t.Errorf("unexpected error: %v", err)
			}

			// Check if the response body matches the expected body
			if tt.expectedBody != "" {
				b, err := json.Marshal(resp)
				if err != nil {
					t.Fatal(err)
				}
				t.Logf("Response: %s", string(b))
				if string(b) != tt.expectedBody {
					t.Errorf("expected response body %s, got %s", tt.expectedBody, string(b))
				}
			}
		})
	}
}

func TestDoWithHeaders(t *testing.T) {
	tests := map[string]struct {
		method       string
		path         string
		headers      map[string]string
		body         interface{}
		version      string
		expectedErr  error
		expectedBody string
	}{
		"successful GET request": {
			method:       http.MethodGet,
			path:         "/api/test",
			headers:      map[string]string{"Content-Type": "application/json"},
			body:         nil,
			expectedErr:  nil,
			expectedBody: `{"message":"success"}`,
		},
		"successful POST request": {
			method:       http.MethodPost,
			path:         "/api/test",
			headers:      map[string]string{"Content-Type": "application/json"},
			body:         nil,
			expectedErr:  nil,
			expectedBody: `{"message":"success"}`,
		},
		"successful PUT request": {
			method:       http.MethodPut,
			path:         "/api/test",
			headers:      map[string]string{"Content-Type": "application/json"},
			body:         nil,
			expectedErr:  nil,
			expectedBody: `{"message":"success"}`,
		},
		"successful DELETE request": {
			method:      http.MethodDelete,
			path:        "/api/test",
			headers:     map[string]string{"Content-Type": "application/json"},
			body:        nil,
			expectedErr: nil,
		},
		"Invalid version": {
			method:      http.MethodGet,
			path:        "api/test",
			version:     "invalid_version",
			expectedErr: fmt.Errorf("strconv.ParseFloat: parsing \"invalid_version\": invalid syntax"),
		},
		"read close error": {
			method: http.MethodGet,
			path:   "/api/test",
			body: &TestReadCloser{
				reader:   strings.NewReader("Test request body"),
				closeErr: errors.New("transport connection broken"),
			},
			version:     "4.0",
			expectedErr: errors.New("transport connection broken"),
		},
		"Invalid Json": {
			method:      http.MethodGet,
			path:        "/api/test",
			version:     "4.0",
			expectedErr: errors.New("looking for beginning of value"),
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			// Create a test server to handle the request
			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method != tt.method {
					t.Fatal(fmt.Errorf("wrong method. Expected %s; but got %s", tt.method, r.Method))
				}

				if r.URL.Path != tt.path {
					t.Fatal(fmt.Errorf("wrong path. Expected %s; but got %s", tt.path, r.URL.Path))
				}

				for k, v := range tt.headers {
					if r.Header.Get(k) != v {
						t.Fatal(fmt.Errorf("wrong header value. Expected %s; but got %s", v, r.Header.Get(k)))
					}
				}

				if tt.expectedErr != nil && tt.expectedErr.Error() != "looking for beginning of value" {
					w.WriteHeader(http.StatusBadRequest)
					w.Write([]byte(tt.expectedErr.Error()))
				} else if tt.expectedErr != nil && tt.expectedErr.Error() == "looking for beginning of value" {
					w.WriteHeader(http.StatusOK)
					// return an invalid JSON that will cause an error when decoding
					w.Write([]byte(`invalid json`))
				} else {
					w.WriteHeader(http.StatusOK)
					w.Write([]byte(tt.expectedBody))
				}
			}))
			defer ts.Close()

			// Create a new client and set the host to the test server
			c, err := New(context.Background(), ts.URL, ClientOptions{}, false)
			if err != nil {
				t.Fatal(err)
			}

			// Call the DoWithHeaders function with the test parameters
			var resp interface{}
			err = c.DoWithHeaders(
				context.Background(),
				tt.method,
				tt.path,
				tt.headers,
				tt.body,
				&resp,
				tt.version,
			)

			// Check if the error matches the expected error
			if tt.expectedErr != nil {
				if err == nil || !strings.Contains(err.Error(), tt.expectedErr.Error()) {
					t.Errorf("expected error %v, got %v", tt.expectedErr, err)
				}
			} else if err != nil {
				t.Errorf("unexpected error: %v", err)
			}

			// Check if the response body matches the expected body
			if tt.expectedBody != "" {
				b, err := json.Marshal(resp)
				if err != nil {
					t.Fatal(err)
				}
				if string(b) != tt.expectedBody {
					t.Errorf("expected response body %s, got %s", tt.expectedBody, string(b))
				}
			}
		})
	}
}

type SamplePayload struct {
	Message string `json:"message"`
}

type TestReadCloser struct {
	reader   io.Reader
	closeErr error
}

func (r *TestReadCloser) Read(p []byte) (n int, err error) {
	return r.reader.Read(p)
}

func (r *TestReadCloser) Close() error {
	return r.closeErr
}

func TestDoAndGetResponseBody(t *testing.T) {
	tests := map[string]struct {
		method       string
		path         string
		headers      map[string]string
		body         interface{}
		version      string
		token        string
		expectedErr  error
		expectedBody string
	}{
		"successful GET request": {
			method:       http.MethodGet,
			path:         "/api/test",
			headers:      map[string]string{"Content-Type": "application/json"},
			body:         nil,
			version:      "4.0",
			expectedErr:  nil,
			expectedBody: `{"message":"success"}`,
		},
		"successful GET request with custom body": {
			method:       http.MethodGet,
			path:         "/api/test",
			headers:      map[string]string{"Content-Type": "application/json"},
			body:         SamplePayload{},
			version:      "4.0",
			expectedErr:  nil,
			expectedBody: `{"message":"success"}`,
		},
		"successful GET request with readcloser body": {
			method:  http.MethodGet,
			path:    "/api/test",
			headers: map[string]string{"Content-Type": "application/json"},
			body: &TestReadCloser{
				reader: strings.NewReader("Test request body"),
			},
			version:      "4.0",
			expectedErr:  nil,
			expectedBody: `{"message":"success"}`,
		},
		"No headers": {
			method:  http.MethodGet,
			path:    "/api/test",
			headers: nil,
			body: &TestReadCloser{
				reader: strings.NewReader("Test request body"),
			},
			version:      "4.0",
			expectedErr:  nil,
			expectedBody: `{"message":"success"}`,
		},
		"Version 4 with token": {
			method:       http.MethodGet,
			path:         "/api/test",
			headers:      nil,
			body:         nil,
			token:        "test-token",
			version:      "4.0",
			expectedErr:  nil,
			expectedBody: `{"message":"success"}`,
		},
		"Version less than 4 with token": {
			method:       http.MethodGet,
			path:         "/api/test",
			headers:      nil,
			body:         nil,
			token:        "test-token",
			version:      "2.0",
			expectedErr:  nil,
			expectedBody: `{"message":"success"}`,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			// Create a test server to handle the request
			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method != tt.method {
					t.Fatal(fmt.Errorf("wrong method. Expected %s; but got %s", tt.method, r.Method))
				}

				if r.URL.Path != tt.path {
					t.Fatal(fmt.Errorf("wrong path. Expected %s; but got %s", tt.path, r.URL.Path))
				}

				for k, v := range tt.headers {
					if r.Header.Get(k) != v {
						t.Fatal(fmt.Errorf("wrong header value. Expected %s; but got %s", v, r.Header.Get(k)))
					}
				}

				if tt.expectedErr != nil {
					w.WriteHeader(http.StatusBadRequest)
					data, err := json.Marshal(tt.expectedErr)
					if err != nil {
						t.Fatal(err)
					}
					w.Write(data)
				} else {
					w.WriteHeader(http.StatusOK)
					w.Write([]byte(tt.expectedBody))
				}
			}))
			defer ts.Close()

			// Create a new client and set the host to the test server
			c, err := New(context.Background(), ts.URL, ClientOptions{}, false)
			if err != nil {
				t.Fatal(err)
			}

			c.SetToken(tt.token)

			// Call the DoAndGetResponseBody function with the test parameters
			res, err := c.DoAndGetResponseBody(
				context.Background(),
				tt.method,
				tt.path,
				tt.headers,
				tt.body, tt.version)

			// Check if the error matches the expected error
			if tt.expectedErr != nil {
				if err == nil || err.Error() != tt.expectedErr.Error() {
					t.Errorf("expected error %v, got %v", tt.expectedErr, err)
				}
			} else if err != nil {
				t.Errorf("unexpected error: %v", err)
			}

			// Check if the response body matches the expected body
			if tt.expectedBody != "" {
				b, err := ioutil.ReadAll(res.Body)
				if err != nil {
					t.Fatal(err)
				}
				if string(b) != tt.expectedBody {
					t.Errorf("expected response body %s, got %s", tt.expectedBody, string(b))
				}
			}
		})
	}
}

func TestParseJSONError(t *testing.T) {
	tests := map[string]struct {
		name        string
		response    *http.Response
		expectedErr error
	}{
		"JSON response": {
			response: &http.Response{
				StatusCode: http.StatusBadRequest,
				Header:     http.Header{"Content-Type": []string{"application/json"}},
				Body:       ioutil.NopCloser(strings.NewReader(`{"message":"Bad Request"}`)),
			},
			expectedErr: &types.Error{
				HTTPStatusCode: http.StatusBadRequest,
				Message:        "Bad Request",
			},
		},
		"HTML response": {
			response: &http.Response{
				StatusCode: http.StatusBadRequest,
				Status:     "Bad Request",
				Header:     http.Header{"Content-Type": []string{"text/html"}},
				Body:       ioutil.NopCloser(strings.NewReader("<html><body>Bad Request</body></html>")),
			},
			expectedErr: &types.Error{
				HTTPStatusCode: http.StatusBadRequest,
				Message:        "Bad Request",
			},
		},
		"No content type": {
			response: &http.Response{
				StatusCode: http.StatusBadRequest,
				Body:       ioutil.NopCloser(strings.NewReader(`{"message":"Bad Request"}`)),
			},
			expectedErr: &types.Error{
				HTTPStatusCode: http.StatusBadRequest,
				Message:        "Bad Request",
			},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			c := &client{}
			err := c.ParseJSONError(tt.response)
			if !reflect.DeepEqual(err, tt.expectedErr) {
				t.Errorf("Expected error %v, got %v", tt.expectedErr, err)
			}
		})
	}

	t.Run("Bad response", func(t *testing.T) {
		c := &client{}
		response := &http.Response{
			Body: ioutil.NopCloser(strings.NewReader("{Bad response}")),
		}
		err := c.ParseJSONError(response)
		assert.NotNil(t, err)
	})
}
