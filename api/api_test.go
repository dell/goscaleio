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
)

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

func TestDoWithHeaders(t *testing.T) {
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
					w.Write([]byte(tt.expectedErr.Error()))
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
				"",
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
	reader io.Reader
}

func (r *TestReadCloser) Read(p []byte) (n int, err error) {
	return r.reader.Read(p)
}

func (r *TestReadCloser) Close() error {
	return nil
}

func TestDoAndGetResponseBody(t *testing.T) {
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
		"successful GET request with custom body": {
			method:       http.MethodGet,
			path:         "/api/test",
			headers:      map[string]string{"Content-Type": "application/json"},
			body:         SamplePayload{},
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

			// Call the DoAndGetResponseBody function with the test parameters
			res, err := c.DoAndGetResponseBody(
				context.Background(),
				tt.method,
				tt.path,
				tt.headers,
				tt.body, "4.0")

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
}