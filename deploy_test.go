// Copyright © 2023 - 2024 Dell Inc. or its subsidiaries. All Rights Reserved.
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
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"strings"
	"testing"

	types "github.com/dell/goscaleio/types/v1"
	"github.com/stretchr/testify/assert"
)

// TestNewGateway tests the NewGateway function.
func TestNewGateway(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost && r.URL.Path == "/rest/auth/login" {

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			fmt.Fprintln(w, `{"access_token":"mock_access_token"}`)
			return
		}
		if r.Method == http.MethodGet && r.URL.Path == "/api/version" {
			w.Header().Set("Content-Type", "text/plain")
			w.WriteHeader(http.StatusOK)
			fmt.Fprintln(w, "4.0")
			return
		}
		http.NotFound(w, r)
	}))

	gc, err := NewGateway(server.URL, "test_username", "test_password", false, false)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	assert.Nil(t, err, "Unexpected error")
	assert.NotNil(t, gc, "GatewayClient is nil")
	assert.Equal(t, "mock_access_token", gc.token, "Unexpected access token")
	assert.Equal(t, "4.0", gc.version, "Unexpected version")

	// error test - empty host
	gc, err = NewGateway("", "test_username", "test_password", false, false)
	assert.Nil(t, gc, "GatewayClient is not nil")
	assert.Error(t, err, "missing endpoint")

	server.Close()

	////////////////
	// version response tests
	////////////////
	// error response
	server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost && r.URL.Path == "/rest/auth/login" {

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			fmt.Fprintln(w, `{"access_token":"mock_access_token"}`)
			return
		}
		if r.Method == http.MethodGet && r.URL.Path == "/api/version" {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
	}))
	gc, err = NewGateway(server.URL, "test_username", "test_password", false, false)
	assert.NotNil(t, err, "Expected error")
	assert.Nil(t, gc, "GatewayClient is nil")
	server.Close()

	// version 3.5 path
	server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost && r.URL.Path == "/rest/auth/login" {

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			fmt.Fprintln(w, `{"access_token":"mock_access_token"}`)
			return
		}
		if r.Method == http.MethodGet && r.URL.Path == "/api/version" {
			w.Header().Set("Content-Type", "text/plain")
			w.WriteHeader(http.StatusOK)
			fmt.Fprintln(w, "3.5")
			return
		}
		http.NotFound(w, r)
	}))
	gc, err = NewGateway(server.URL, "test_username", "test_password", false, false)
	assert.Nil(t, err, "Unexpected error")
	assert.NotNil(t, gc, "GatewayClient is nil")
	assert.Equal(t, "", gc.token, "") // no token for 3.5
	assert.Equal(t, "3.5", gc.version, "Unexpected version")
}

func TestNewGatewayInsecure(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost && r.URL.Path == "/rest/auth/login" {

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			fmt.Fprintln(w, `{"access_token":"mock_access_token"}`)
			return
		}
		if r.Method == http.MethodGet && r.URL.Path == "/api/version" {
			w.Header().Set("Content-Type", "text/plain")
			w.WriteHeader(http.StatusOK)
			fmt.Fprintln(w, "4.0")
			return
		}
		http.NotFound(w, r)
	}))

	defer server.Close()

	gc, err := NewGateway(server.URL, "test_username", "test_password", true, false)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	assert.Nil(t, err, "Unexpected error")
	assert.NotNil(t, gc, "GatewayClient is nil")
	assert.Equal(t, "mock_access_token", gc.token, "Unexpected access token")
	assert.Equal(t, "4.0", gc.version, "Unexpected version")
}

// errorTransport simulates an error during response body reading
type errorTransport struct{}

func (t *errorTransport) RoundTrip(_ *http.Request) (*http.Response, error) {
	return &http.Response{
		StatusCode: http.StatusOK,
		Body:       io.NopCloser(&errorReader{}),
	}, nil
}

type errorReader struct{}

func (r *errorReader) Read(_ []byte) (n int, err error) {
	return 0, errBodyRead
}

func (r *errorReader) Close() error {
	return nil
}

// TestGetVersion tests the GetVersion function.
func TestGetVersion(t *testing.T) {
	tests := []struct {
		name        string
		setup       func() *GatewayClient
		expected    string
		expectedErr string
	}{
		{
			name: "successful retrieval with basic auth",
			setup: func() *GatewayClient {
				server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					if r.Method == http.MethodGet && r.URL.Path == "/api/version" {
						w.Header().Set("Content-Type", "text/plain")
						w.WriteHeader(http.StatusOK)
						fmt.Fprintln(w, "4.0")
						return
					}
					http.NotFound(w, r)
				}))
				return &GatewayClient{
					host:     server.URL,
					http:     server.Client(),
					username: "test_username",
					password: "test_password",
				}
			},
			expected:    "4.0",
			expectedErr: "",
		},
		{
			name: "successful retrieval with bearer token",
			setup: func() *GatewayClient {
				server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					if r.Method == http.MethodGet && r.URL.Path == "/api/version" {
						w.Header().Set("Content-Type", "text/plain")
						w.WriteHeader(http.StatusOK)
						fmt.Fprintln(w, "4.0")
						return
					}
					http.NotFound(w, r)
				}))
				return &GatewayClient{
					host:  server.URL,
					http:  server.Client(),
					token: "dummy_token",
				}
			},
			expected:    "4.0",
			expectedErr: "",
		},
		{
			name: "http request creation error",
			setup: func() *GatewayClient {
				return &GatewayClient{
					host: "http://[::1]:namedport",
					http: &http.Client{},
				}
			},
			expected:    "",
			expectedErr: "invalid port \":namedport\" after host",
		},
		{
			name: "non-2xx status code",
			setup: func() *GatewayClient {
				server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
					http.Error(w, "error", http.StatusBadRequest)
				}))
				return &GatewayClient{
					host: server.URL,
					http: server.Client(),
				}
			},
			expected:    "",
			expectedErr: "Error: 400 Bad Request",
		},
		{
			name: "error extracting version string",
			setup: func() *GatewayClient {
				server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
					w.WriteHeader(http.StatusOK)
					w.Write([]byte(`{"version": "invalid"}`))
				}))
				return &GatewayClient{
					host: server.URL,
					http: &http.Client{
						Transport: &errorTransport{},
					},
				}
			},
			expected:    "",
			expectedErr: "error reading body",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gc := tt.setup()
			version, err := gc.GetVersion()
			if tt.expectedErr != "" {
				assert.NotNil(t, err)
				assert.Contains(t, err.Error(), tt.expectedErr)
			} else {
				assert.Nil(t, err)
				assert.Equal(t, tt.expected, version)
			}
		})
	}
}

// TestUploadPackages tests the UploadPackages function.
func TestUploadPackages(t *testing.T) {
	respStatus := http.StatusOK
	respMessage := ""
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost && r.URL.Path == "/im/types/installationPackages/instances/actions/uploadPackages" {
			w.WriteHeader(respStatus)
			if respMessage != "" {
				w.Write([]byte(respMessage))
			}
			return
		}
		http.NotFound(w, r)
	}))

	defer server.Close()

	gc := &GatewayClient{
		http:     &http.Client{},
		host:     server.URL,
		username: "test_username",
		password: "test_password",
	}

	t.Run("file does not exist", func(t *testing.T) {
		_, err := gc.UploadPackages([]string{"mock_file.tar"})
		assert.Error(t, err)

		expectedErrorMsg := "file mock_file.tar does not exist"
		assert.EqualError(t, err, expectedErrorMsg)
	})

	t.Run("wrong file type", func(t *testing.T) {
		name := "test_file.log"
		err := os.WriteFile(name, []byte("package data"), 0600)
		assert.NoError(t, err)
		defer os.Remove(name)

		_, err = gc.UploadPackages([]string{name})
		assert.ErrorContains(t, err, "invalid file type")
	})

	t.Run("successful upload", func(t *testing.T) {
		name := "test_file.tar"
		err := os.WriteFile(name, []byte("package data"), 0600)
		assert.NoError(t, err)

		defer os.Remove(name)

		_, err = gc.UploadPackages([]string{name})
		assert.NoError(t, err)

		gc.version = "4.0"
		defer func() {
			gc.version = ""
		}()

		_, err = gc.UploadPackages([]string{name})
		assert.NoError(t, err)
	})

	t.Run("bad response code", func(t *testing.T) {
		name := "test_file.tar"
		err := os.WriteFile(name, []byte("package data"), 0600)
		assert.NoError(t, err)
		defer os.Remove(name)

		// Induce a bad response code
		respStatus = http.StatusConflict
		defer func() {
			respStatus = http.StatusOK
		}()

		_, err = gc.UploadPackages([]string{name})
		assert.ErrorContains(t, err, "failed to parse response body")

		// Induce a bad response code with message
		respMessage = `{"message":"induced failure"}`
		defer func() {
			respMessage = ""
		}()

		_, err = gc.UploadPackages([]string{name})
		assert.ErrorContains(t, err, "received bad response")
	})

	t.Run("error - set cookie", func(t *testing.T) {
		defaultCookieFunc := setCookieFunc
		setCookieFunc = func(_ http.Header, _ string) error {
			return errors.New("cookie error")
		}
		name := "test_file.tar"
		err := os.WriteFile(name, []byte("package data"), 0600)
		assert.NoError(t, err)

		defer os.Remove(name)

		gc.version = "4.0"

		_, err = gc.UploadPackages([]string{name})
		assert.Error(t, err)
		setCookieFunc = defaultCookieFunc
	})
}

func TestParseCSV(t *testing.T) {
	respStatus := http.StatusOK
	respBody := "-"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(respStatus)
		if respBody != "-" {
			w.Write([]byte(respBody))
		} else {
			w.Write([]byte(`{"masterMdm": "data"}`))
		}
	}))
	defer server.Close()

	t.Run("successful parse with bearer token", func(t *testing.T) {
		file, err := os.CreateTemp("", "test_file.csv")
		if err != nil {
			t.Fatalf("Failed to create temp file: %v", err)
		}
		defer file.Close()

		_, err = file.WriteString("header1,header2\nvalue1,value2")
		if err != nil {
			t.Fatalf("Failed to write to temp file: %v", err)
		}

		defer os.Remove(file.Name())

		gc := &GatewayClient{
			http:    server.Client(),
			host:    server.URL,
			version: "4.0",
			token:   "test_token",
		}

		response, err := gc.ParseCSV(file.Name())
		assert.NoError(t, err)
		assert.Equal(t, 200, response.StatusCode)
	})

	t.Run("successful parse with basic auth", func(t *testing.T) {
		file, err := os.CreateTemp("", "test_file.csv")
		if err != nil {
			t.Fatalf("Failed to create temp file: %v", err)
		}
		defer file.Close()

		_, err = file.WriteString("header1,header2\nvalue1,value2")
		if err != nil {
			t.Fatalf("Failed to write to temp file: %v", err)
		}

		defer os.Remove(file.Name())

		gc := &GatewayClient{
			http:     server.Client(),
			host:     server.URL,
			username: "test_username",
			password: "test_password",
		}

		response, err := gc.ParseCSV(file.Name())
		assert.NoError(t, err)
		assert.Equal(t, 200, response.StatusCode)
	})

	t.Run("bad response code", func(t *testing.T) {
		name := "test_file.csv"

		err := os.WriteFile(name, []byte("header1,header2\nvalue1,value2"), 0600)
		assert.NoError(t, err)
		defer os.Remove(name)

		gc := &GatewayClient{
			http:     server.Client(),
			host:     server.URL,
			username: "test_username",
			password: "test_password",
		}

		// Induce a bad response code with json message
		respStatus = http.StatusInternalServerError
		respBody = `{"message": "error"}`
		defer func() {
			respStatus = http.StatusOK
			respBody = "-"
		}()

		_, err = gc.ParseCSV(name)
		assert.Error(t, err)

		// Induce a bad response code with malformed json message
		respBody = `malformed json`
		_, err = gc.ParseCSV(name)
		assert.Error(t, err)
	})

	t.Run("good response code, but no mdm", func(t *testing.T) {
		name := "test_file.csv"
		err := os.WriteFile(name, []byte("header1,header2\nvalue1,value2"), 0600)
		assert.NoError(t, err)
		defer os.Remove(name)

		gc := &GatewayClient{
			http:     server.Client(),
			host:     server.URL,
			username: "test_username",
			password: "test_password",
		}

		// Induce a dummy response body
		respBody = `{}`
		defer func() {
			respStatus = http.StatusOK
			respBody = "-"
		}()

		_, err = gc.ParseCSV(name)
		assert.Error(t, err)
	})

	t.Run("file not found", func(t *testing.T) {
		gc := &GatewayClient{}
		_, err := gc.ParseCSV("nonexistent.csv")
		assert.Error(t, err)
	})

	t.Run("cookie error", func(t *testing.T) {
		defaultCookieFunc := setCookieFunc
		setCookieFunc = func(_ http.Header, _ string) error {
			return errors.New("cookie error")
		}
		file, err := os.CreateTemp("", "test_file.csv")
		if err != nil {
			t.Fatalf("Failed to create temp file: %v", err)
		}
		defer file.Close()

		_, err = file.WriteString("header1,header2\nvalue1,value2")
		if err != nil {
			t.Fatalf("Failed to write to temp file: %v", err)
		}

		defer os.Remove(file.Name())

		gc := &GatewayClient{
			http:     server.Client(),
			host:     server.URL,
			username: "test_username",
			password: "test_password",
			version:  "4.0",
		}

		_, err = gc.ParseCSV(file.Name())
		assert.Error(t, err)
		setCookieFunc = defaultCookieFunc
	})
}

func TestGetPackageDetails(t *testing.T) {
	// Define the desired response JSON
	responseJSON := `[
    {
        "version": "4.5-0.287",
        "sioPatchNumber": 0,
        "type": "mdm",
        "size": 72378708,
        "label": "0.287.sles15.3.x86_64",
        "operatingSystem": "linux",
        "linuxFlavour": "sles15_3",
        "activemqPackage": false,
        "filename": "EMC-ScaleIO-mdm-4.5-0.287.sles15.3.x86_64.rpm",
        "activemqRpmPackage": false,
        "activemqUbuntuPackage": false,
        "latest": true
    },
    {
        "version": "5.16-4.62",
        "sioPatchNumber": 0,
        "type": "activemq",
        "size": 65279904,
        "label": "62",
        "operatingSystem": "linux",
        "linuxFlavour": "all_linux_rpm_flavors",
        "activemqPackage": true,
        "filename": "EMC-ScaleIO-activemq-5.16.4-62.noarch.rpm",
        "activemqRpmPackage": true,
        "activemqUbuntuPackage": false,
        "latest": true
    }]`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet && r.URL.Path == "/im/types/installationPackages/instances" {
			w.WriteHeader(http.StatusOK)
			_, err := w.Write([]byte(responseJSON))
			if err != nil {
				t.Fatalf("Error writing response: %v", err)
			}
			return
		}
		http.NotFound(w, r)
	}))

	defer server.Close()

	gc := &GatewayClient{
		http:     &http.Client{},
		host:     server.URL,
		username: "test_username",
		password: "test_password",
		version:  "4.0",
	}

	t.Run("successful response with bearer token", func(t *testing.T) {
		packageDetails, err := gc.GetPackageDetails()
		assert.NoError(t, err)
		assert.NotNil(t, packageDetails)
	})
	t.Run("set cookie error", func(t *testing.T) {
		defaultCookieFunc := setCookieFunc
		setCookieFunc = func(_ http.Header, _ string) error {
			return errors.New("cookie error")
		}
		_, err := gc.GetPackageDetails()
		assert.Error(t, err)
		setCookieFunc = defaultCookieFunc
	})
	t.Run("successful response with basic auth", func(t *testing.T) {
		gc.version = "3.0"
		packageDetails, err := gc.GetPackageDetails()
		assert.NoError(t, err)
		assert.NotNil(t, packageDetails)
	})
	t.Run("HTTP request creation error", func(t *testing.T) {
		gc.host = "http://[::1]:namedport"
		_, err := gc.GetPackageDetails()
		assert.NotNil(t, err)
	})
	t.Run("Error unmarshalling response", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`invalid json`))
		}))
		defer server.Close()

		gc.host = server.URL
		_, err := gc.GetPackageDetails()
		assert.NotNil(t, err)
		assert.Contains(t, err.Error(), "Error For Get Package Details")
	})
}

func TestDeletePackage(t *testing.T) {
	t.Run("successful response with basic auth", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusOK)
		}))
		defer server.Close()

		gc := &GatewayClient{
			http:    &http.Client{},
			host:    server.URL,
			version: "4.0",
			token:   "dummy_token",
		}

		packageResponse, err := gc.DeletePackage("test_package")
		assert.NoError(t, err)
		assert.Equal(t, 200, packageResponse.StatusCode)
	})
	t.Run("error - set cookies", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusOK)
		}))
		defer server.Close()
		defaultSetCookieFunc := setCookieFunc
		setCookieFunc = func(_ http.Header, _ string) error {
			return errors.New("cookie error")
		}

		gc := &GatewayClient{
			http:    &http.Client{},
			host:    server.URL,
			version: "4.0",
			token:   "dummy_token",
		}

		_, err := gc.DeletePackage("test_package")
		assert.Error(t, err)
		setCookieFunc = defaultSetCookieFunc
	})
	t.Run("successful response with bearer token", func(t *testing.T) {
		responseJSON := `{
		"StatusCode": 200
		}`

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method == "DELETE" && strings.HasPrefix(r.URL.Path, "/im/types/installationPackages/instances/actions/delete") {
				w.WriteHeader(http.StatusOK)
				_, err := w.Write([]byte(responseJSON))
				if err != nil {
					t.Fatalf("Error writing response: %v", err)
				}
				return
			}
			http.NotFound(w, r)
		}))

		defer server.Close()

		gc := &GatewayClient{
			http:     &http.Client{},
			host:     server.URL,
			version:  "3.0",
			username: "test_username",
			password: "test_password",
		}
		packageResponse, err := gc.DeletePackage("test_package")
		assert.NoError(t, err)
		assert.Equal(t, 200, packageResponse.StatusCode)
	})
	t.Run("non 200 status code", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusBadRequest) // Simulate a non-200 status code
			response := types.GatewayResponse{
				Message: "Bad Request",
			}
			_ = json.NewEncoder(w).Encode(response)
		}))
		defer server.Close()

		gc := &GatewayClient{
			http:    &http.Client{},
			host:    server.URL,
			version: "4.0",
			token:   "dummy_token",
		}

		packageResponse, err := gc.DeletePackage("test_package")
		assert.NoError(t, err)
		assert.Equal(t, "Bad Request", packageResponse.Message)
	})
}

func TestBeginInstallation(t *testing.T) {
	defaultCookiesFunc := setCookieFunc
	after := func() {
		setCookieFunc = defaultCookiesFunc
	}
	tests := map[string]struct {
		server           *httptest.Server
		version          string
		expectedResponse *types.GatewayResponse
		expectedErr      error
		expectedStatus   int
		setup            func()
	}{
		"success with version 4.0": {
			server: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method == http.MethodPost && strings.HasPrefix(r.URL.Path, "/im/types/Configuration/actions/install") {
					w.WriteHeader(http.StatusAccepted)
					return
				}
				http.NotFound(w, r)
			})),
			version:        "4.0",
			expectedStatus: http.StatusOK,
		},
		"fail - setCookie": {
			server: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method == http.MethodPost && strings.HasPrefix(r.URL.Path, "/im/types/Configuration/actions/install") {
					w.WriteHeader(http.StatusAccepted)
					return
				}
				http.NotFound(w, r)
			})),
			version:        "4.0",
			expectedStatus: -1,
			expectedErr:    errors.New("Error While Handling Cookie: cookie error"),
			setup: func() {
				setCookieFunc = func(_ http.Header, _ string) error {
					return errors.New("cookie error")
				}
			},
		},
		"success with version < 4.0": {
			server: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method == http.MethodPost && strings.HasPrefix(r.URL.Path, "/im/types/Configuration/actions/install") {
					w.WriteHeader(http.StatusAccepted)
					return
				}
				http.NotFound(w, r)
			})),
			version:        "3.6",
			expectedStatus: http.StatusOK,
		},
		"non 200 status code": {
			server: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				w.WriteHeader(http.StatusBadRequest) // Simulate a non-200 status code
				response := types.GatewayResponse{
					Message: "Bad Request",
				}
				_ = json.NewEncoder(w).Encode(response)
			})),
			expectedResponse: &types.GatewayResponse{
				Message: "Bad Request",
			},
			expectedErr: nil,
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			defer tt.server.Close()
			defer after()

			gc := &GatewayClient{
				http:     &http.Client{},
				host:     tt.server.URL,
				username: "test_username",
				password: "test_password",
				version:  tt.version,
			}
			if tt.setup != nil {
				tt.setup()
			}

			resp, err := gc.BeginInstallation("{}", "mdm_user", "mdm_password", "lia_password", true, true, true, false)
			if tt.expectedErr != nil {
				assert.EqualError(t, err, tt.expectedErr.Error())
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, resp)
				assert.Equal(t, tt.expectedStatus, resp.StatusCode)
			}
		})
	}
}

func TestMoveToNextPhase(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost && r.URL.Path == "/im/types/ProcessPhase/actions/moveToNextPhase" {
			w.WriteHeader(http.StatusOK)
			return
		}
		http.NotFound(w, r)
	}))
	defer server.Close()

	gc := &GatewayClient{
		http: &http.Client{},
		host: server.URL,
	}

	t.Run("successful response with basic auth", func(t *testing.T) {
		gc.username = "test_username"
		gc.password = "test_password"

		gatewayResponse, err := gc.MoveToNextPhase()
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, gatewayResponse.StatusCode)
	})
	t.Run("fail - setCookie", func(t *testing.T) {
		temp := setCookieFunc
		setCookieFunc = func(_ http.Header, _ string) error {
			return errors.New("cookie error")
		}
		gc.username = "test_username"
		gc.password = "test_password"
		gc.version = "4.0"

		_, err := gc.MoveToNextPhase()
		assert.Error(t, err)
		setCookieFunc = temp
	})
	t.Run("successful response with bearer token", func(t *testing.T) {
		gc.version = "4.0"
		gc.token = "dummy_token"

		gatewayResponse, err := gc.MoveToNextPhase()
		assert.NoError(t, err)
		assert.NotNil(t, gatewayResponse)
		assert.Equal(t, http.StatusOK, gatewayResponse.StatusCode)
	})
	t.Run("non 200 status code", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusBadRequest) // Simulate a non-200 status code
			response := types.GatewayResponse{
				Message: "Bad Request",
			}
			_ = json.NewEncoder(w).Encode(response)
		}))
		defer server.Close()

		gc.host = server.URL
		gc.version = "4.0"
		gc.token = "dummy_token"

		gatewayResponse, err := gc.MoveToNextPhase()
		assert.NoError(t, err)
		assert.Equal(t, "Bad Request", gatewayResponse.Message)
	})
}

func TestRetryPhase(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost && strings.HasPrefix(r.URL.Path, "/im/types/Command/instances/actions/retry") {
			w.WriteHeader(http.StatusOK)
			return
		}
		http.NotFound(w, r)
	}))
	defer server.Close()

	gc := &GatewayClient{
		http: &http.Client{},
		host: server.URL,
	}

	t.Run("successful response with basic auth", func(t *testing.T) {
		gc.username = "test_username"
		gc.password = "test_password"

		gatewayResponse, err := gc.RetryPhase()
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, gatewayResponse.StatusCode)
	})
	t.Run("fail - setCookie", func(t *testing.T) {
		temp := setCookieFunc
		setCookieFunc = func(_ http.Header, _ string) error {
			return errors.New("cookie error")
		}
		gc.username = "test_username"
		gc.password = "test_password"
		gc.version = "4.0"

		_, err := gc.RetryPhase()
		assert.Error(t, err)
		setCookieFunc = temp
	})
	t.Run("successful response with bearer token", func(t *testing.T) {
		gc.version = "4.0"
		gc.token = "dummy_token"

		gatewayResponse, err := gc.RetryPhase()
		assert.NoError(t, err)
		assert.NotNil(t, gatewayResponse)
		assert.Equal(t, http.StatusOK, gatewayResponse.StatusCode)
	})
	t.Run("non 200 status code", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusBadRequest) // Simulate a non-200 status code
			response := types.GatewayResponse{
				Message: "Bad Request",
			}
			_ = json.NewEncoder(w).Encode(response)
		}))
		defer server.Close()

		gc.host = server.URL
		gc.version = "4.0"
		gc.token = "dummy_token"

		gatewayResponse, err := gc.RetryPhase()
		assert.NoError(t, err)
		assert.Equal(t, "Bad Request", gatewayResponse.Message)
	})
}

func TestAbortOperation(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost && r.URL.Path == "/im/types/Command/instances/actions/abort" {
			w.WriteHeader(http.StatusOK)
			return
		}
		http.NotFound(w, r)
	}))
	defer server.Close()

	gc := &GatewayClient{
		http: &http.Client{},
		host: server.URL,
	}

	t.Run("successful response with basic auth", func(t *testing.T) {
		gc.username = "test_username"
		gc.password = "test_password"

		gatewayResponse, err := gc.AbortOperation()
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, gatewayResponse.StatusCode)
	})
	t.Run("fail - setCookie", func(t *testing.T) {
		temp := setCookieFunc
		setCookieFunc = func(_ http.Header, _ string) error {
			return errors.New("cookie error")
		}
		gc.username = "test_username"
		gc.password = "test_password"
		gc.version = "4.0"

		_, err := gc.AbortOperation()
		assert.Error(t, err)
		setCookieFunc = temp
	})
	t.Run("successful response with bearer token", func(t *testing.T) {
		gc.version = "4.0"
		gc.token = "dummy_token"

		gatewayResponse, err := gc.AbortOperation()
		assert.NoError(t, err)
		assert.NotNil(t, gatewayResponse)
		assert.Equal(t, http.StatusOK, gatewayResponse.StatusCode)
	})
	t.Run("non 200 status code", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusBadRequest) // Simulate a non-200 status code
			response := types.GatewayResponse{
				Message: "Bad Request",
			}
			_ = json.NewEncoder(w).Encode(response)
		}))
		defer server.Close()

		gc.host = server.URL
		gc.version = "4.0"
		gc.token = "dummy_token"

		gatewayResponse, err := gc.AbortOperation()
		assert.NoError(t, err)
		assert.Equal(t, "Bad Request", gatewayResponse.Message)
	})
}

func TestClearQueueCommand(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost && r.URL.Path == "/im/types/Command/instances/actions/clear" {
			w.WriteHeader(http.StatusOK)
			return
		}
		http.NotFound(w, r)
	}))
	defer server.Close()

	gc := &GatewayClient{
		http: &http.Client{},
		host: server.URL,
	}

	t.Run("successful response with basic auth", func(t *testing.T) {
		gc.username = "test_username"
		gc.password = "test_password"

		gatewayResponse, err := gc.ClearQueueCommand()
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, gatewayResponse.StatusCode)
	})
	t.Run("fail - setCookie", func(t *testing.T) {
		temp := setCookieFunc
		setCookieFunc = func(_ http.Header, _ string) error {
			return errors.New("cookie error")
		}
		gc.username = "test_username"
		gc.password = "test_password"
		gc.version = "4.0"

		_, err := gc.ClearQueueCommand()
		assert.Error(t, err)
		setCookieFunc = temp
	})
	t.Run("successful response with bearer token", func(t *testing.T) {
		gc.version = "4.0"
		gc.token = "dummy_token"

		gatewayResponse, err := gc.ClearQueueCommand()
		assert.NoError(t, err)
		assert.NotNil(t, gatewayResponse)
		assert.Equal(t, http.StatusOK, gatewayResponse.StatusCode)
	})
	t.Run("non 200 status code", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusBadRequest) // Simulate a non-200 status code
			response := types.GatewayResponse{
				Message: "Bad Request",
			}
			_ = json.NewEncoder(w).Encode(response)
		}))
		defer server.Close()

		gc.host = server.URL
		gc.version = "4.0"
		gc.token = "dummy_token"

		gatewayResponse, err := gc.ClearQueueCommand()
		assert.NoError(t, err)
		assert.Equal(t, "Bad Request", gatewayResponse.Message)
	})
}

func TestMoveToIdlePhase(t *testing.T) {
	t.Run("happy path with basic auth", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method == http.MethodPost && r.URL.Path == "/im/types/ProcessPhase/actions/moveToIdlePhase" {
				w.WriteHeader(http.StatusOK)
				return
			}
			http.NotFound(w, r)
		}))
		defer server.Close()

		gc := &GatewayClient{
			http:     &http.Client{},
			host:     server.URL,
			username: "test_username",
			password: "test_password",
		}

		gatewayResponse, err := gc.MoveToIdlePhase()
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, gatewayResponse.StatusCode)
	})
	t.Run("fail - setCookie", func(t *testing.T) {
		temp := setCookieFunc
		setCookieFunc = func(_ http.Header, _ string) error {
			return errors.New("cookie error")
		}
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method == http.MethodPost && r.URL.Path == "/im/types/ProcessPhase/actions/moveToIdlePhase" {
				w.WriteHeader(http.StatusOK)
				return
			}
			http.NotFound(w, r)
		}))
		defer server.Close()
		gc := &GatewayClient{
			http:     &http.Client{},
			host:     server.URL,
			username: "test_username",
			password: "test_password",
			version:  "4.0",
		}

		_, err := gc.MoveToIdlePhase()
		assert.Error(t, err)
		setCookieFunc = temp
	})
	t.Run("happy path with bearer token", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method == http.MethodPost && r.URL.Path == "/im/types/ProcessPhase/actions/moveToIdlePhase" {
				w.WriteHeader(http.StatusOK)
				return
			}
			http.NotFound(w, r)
		}))
		defer server.Close()

		gc := &GatewayClient{
			http:     &http.Client{},
			host:     server.URL,
			username: "test_username",
			password: "test_password",
			version:  "4.0",
		}

		gatewayResponse, err := gc.MoveToIdlePhase()
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, gatewayResponse.StatusCode)
	})
	t.Run("fail to move to idle phase", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			response := types.GatewayResponse{
				Message: "Bad Request",
			}
			_ = json.NewEncoder(w).Encode(response)
		}))
		defer server.Close()

		gc := &GatewayClient{
			http:     &http.Client{},
			host:     server.URL,
			username: "test_username",
			password: "test_password",
			version:  "4.0",
		}

		gatewayResponse, err := gc.MoveToIdlePhase()
		assert.Nil(t, err)
		assert.Equal(t, "Bad Request", gatewayResponse.Message)
		assert.NotEqual(t, http.StatusOK, gatewayResponse.StatusCode)
	})
}

func TestCheckForCompletionQueueCommands(t *testing.T) {
	t.Run("happy path", func(t *testing.T) {
		responseJSON := `{
		"MDM Commands": []
		}`
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method == http.MethodGet && r.URL.Path == "/im/types/Command/instances" {
				w.WriteHeader(http.StatusOK)
				_, err := w.Write([]byte(responseJSON))
				if err != nil {
					t.Fatalf("Error writing response: %v", err)
				}
				return
			}
			http.NotFound(w, r)
		}))
		defer server.Close()

		gc := &GatewayClient{
			http:     &http.Client{},
			host:     server.URL,
			username: "test_username",
			password: "test_password",
			version:  "4.0",
		}

		gatewayResponse, err := gc.CheckForCompletionQueueCommands("Query")
		assert.NoError(t, err)
		assert.NotNil(t, gatewayResponse)
		assert.Equal(t, http.StatusOK, gatewayResponse.StatusCode)
	})
	t.Run("pending command", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			response := map[string][]interface{}{
				"commands": {
					map[string]interface{}{
						"AllowedPhase": "test-pending",
						"CommandState": "pending",
						"Message":      "error message",
					},
				},
			}
			_ = json.NewEncoder(w).Encode(response)
		}))
		defer server.Close()

		gc := &GatewayClient{
			http:     &http.Client{},
			host:     server.URL,
			username: "test_username",
			password: "test_password",
			version:  "4.0",
		}

		gatewayResponse, err := gc.CheckForCompletionQueueCommands("test-pending")
		assert.Nil(t, err)
		assert.Equal(t, "Running", gatewayResponse.Data)
		assert.Equal(t, http.StatusOK, gatewayResponse.StatusCode)
	})

	t.Run("failed command", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			response := map[string][]interface{}{
				"commands": {
					map[string]interface{}{
						"AllowedPhase": "test-failed",
						"CommandState": "failed",
						"Message":      "error message",
					},
				},
			}
			_ = json.NewEncoder(w).Encode(response)
		}))
		defer server.Close()

		gc := &GatewayClient{
			http:     &http.Client{},
			host:     server.URL,
			username: "test_username",
			password: "test_password",
			version:  "4.0",
		}

		gatewayResponse, err := gc.CheckForCompletionQueueCommands("test-failed")
		assert.Nil(t, err)
		assert.Equal(t, "Failed", gatewayResponse.Data)
		assert.Equal(t, http.StatusOK, gatewayResponse.StatusCode)
	})
}

func TestUninstallCluster(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost && strings.Contains(r.URL.Path, "/im/types/Configuration/actions/uninstall") {
			w.WriteHeader(http.StatusAccepted)
			return
		}
		http.NotFound(w, r)
	}))
	defer server.Close()

	gc := &GatewayClient{
		http: &http.Client{},
		host: server.URL,
	}

	jsonStr := `{
		"snmpIp": null,
		"installationId": null,
		"systemId": null,
		"ingressIp": null,
		"mdmIPs": []
	}`
	mdmUsername := "mdm_username"
	mdmPassword := "mdm_password"
	liaPassword := "lia_password"
	allowNonSecureCommunicationWithMdm := true
	allowNonSecureCommunicationWithLia := true
	disableNonMgmtComponentsAuth := true

	t.Run("successful repsonse with bearer token", func(t *testing.T) {
		gc.version = "4.0"
		gc.token = "dummy_token"

		gatewayResponse, err := gc.UninstallCluster(jsonStr, mdmUsername, mdmPassword, liaPassword, allowNonSecureCommunicationWithMdm, allowNonSecureCommunicationWithLia, disableNonMgmtComponentsAuth, false)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, gatewayResponse.StatusCode)
	})
	t.Run("error - set cookies", func(t *testing.T) {
		temp := setCookieFunc
		setCookieFunc = func(_ http.Header, _ string) error {
			return errors.New("cookie error")
		}
		gc.version = "4.0"
		gc.token = "dummy_token"

		_, err := gc.UninstallCluster(jsonStr, mdmUsername, mdmPassword, liaPassword, allowNonSecureCommunicationWithMdm, allowNonSecureCommunicationWithLia, disableNonMgmtComponentsAuth, false)
		assert.Error(t, err)
		setCookieFunc = temp
	})
	t.Run("successful repsonse with basic auth", func(t *testing.T) {
		gc.username = "test_username"
		gc.password = "test_password"
		gc.version = "3.0"

		gatewayResponse, err := gc.UninstallCluster(jsonStr, mdmUsername, mdmPassword, liaPassword, allowNonSecureCommunicationWithMdm, allowNonSecureCommunicationWithLia, disableNonMgmtComponentsAuth, false)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, gatewayResponse.StatusCode)
	})

	t.Run("non 200 status code", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusBadRequest) // Simulate a non-200 status code
			response := types.GatewayResponse{
				Message: "Bad Request",
			}
			_ = json.NewEncoder(w).Encode(response)
		}))
		defer server.Close()

		gc.host = server.URL
		gc.version = "4.0"
		gc.token = "dummy_token"

		gatewayResponse, err := gc.UninstallCluster(jsonStr, mdmUsername, mdmPassword, liaPassword, allowNonSecureCommunicationWithMdm, allowNonSecureCommunicationWithLia, disableNonMgmtComponentsAuth, false)
		assert.NoError(t, err)
		assert.Equal(t, "Bad Request", gatewayResponse.Message)
	})
}

func TestRenewInstallationCookie(t *testing.T) {
	responseJSON := `[
	{
        "version": "4.5-0.287",
        "sioPatchNumber": 0,
        "type": "mdm",
        "size": 72378708,
        "label": "0.287.sles15.3.x86_64",
        "operatingSystem": "linux",
        "linuxFlavour": "sles15_3",
        "activemqPackage": false,
        "filename": "EMC-ScaleIO-mdm-4.5-0.287.sles15.3.x86_64.rpm",
        "activemqRpmPackage": false,
        "activemqUbuntuPackage": false,
        "latest": true
    }]`
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet && strings.Contains(r.URL.Path, "/im/types/installationPackages/instances") {
			cookie := &http.Cookie{
				Name:  "LEGACYGWCOOKIE",
				Value: "123456789",
				Path:  "/",
			}
			http.SetCookie(w, cookie)
			w.WriteHeader(http.StatusOK)
			_, err := w.Write([]byte(responseJSON))
			if err != nil {
				t.Fatalf("Error writing response: %v", err)
			}
			return
		}
		http.NotFound(w, r)
	}))
	defer server.Close()

	gc := &GatewayClient{
		http:     &http.Client{},
		host:     server.URL,
		username: "test_username",
		password: "test_password",
	}

	err := gc.RenewInstallationCookie(5)
	assert.NoError(t, err)
}

func TestValidateMDMDetails(t *testing.T) {
	defaultCookiesFunc := setCookieFunc
	after := func() {
		setCookieFunc = defaultCookiesFunc
	}
	tests := map[string]struct {
		mdmTopologyParam []byte
		server           *httptest.Server
		expectedResponse *types.GatewayResponse
		version          string
		expectedErr      error
		setup            func()
	}{
		"success with version 4.0": {
			mdmTopologyParam: []byte(`{"mdmUser": "admin", "mdmPassword": "password", "mdmIps": ["192.168.0.1"]}`),
			server: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				resp := types.MDMTopologyDetails{
					SdcIps: []string{"10.0.0.1", "10.0.0.2"},
				}

				data, err := json.Marshal(resp)
				if err != nil {
					t.Fatal(err)
				}
				w.WriteHeader(http.StatusOK)
				w.Write(data)
			})),
			expectedResponse: &types.GatewayResponse{
				StatusCode: 200,
				Data:       "10.0.0.1,10.0.0.2",
			},
			version:     "4.0",
			expectedErr: nil,
		},
		"failure - cookie error": {
			mdmTopologyParam: []byte(`{"mdmUser": "admin", "mdmPassword": "password", "mdmIps": ["192.168.0.1"]}`),
			server: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				resp := types.MDMTopologyDetails{
					SdcIps: []string{"10.0.0.1", "10.0.0.2"},
				}

				data, err := json.Marshal(resp)
				if err != nil {
					t.Fatal(err)
				}
				w.WriteHeader(http.StatusOK)
				w.Write(data)
			})),
			expectedResponse: nil,
			version:          "4.0",
			expectedErr:      errors.New("Error While Handling Cookie: Cookie error"),
			setup: func() {
				setCookieFunc = func(_ http.Header, host string) error {
					return errors.New("Cookie error")
				}
			},
		},
		"success with version < 4.0": {
			mdmTopologyParam: []byte(`{"mdmUser": "admin", "mdmPassword": "password", "mdmIps": ["192.168.0.1"]}`),
			server: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				resp := types.MDMTopologyDetails{
					SdcIps: []string{"10.0.0.1", "10.0.0.2"},
				}

				data, err := json.Marshal(resp)
				if err != nil {
					t.Fatal(err)
				}
				w.WriteHeader(http.StatusOK)
				w.Write(data)
			})),
			expectedResponse: &types.GatewayResponse{
				StatusCode: 200,
				Data:       "10.0.0.1,10.0.0.2",
			},
			version:     "3.6",
			expectedErr: nil,
		},
		"error primary mdm ip": {
			mdmTopologyParam: []byte(`{"mdmUser": "admin", "mdmPassword": "password", "mdmIps": ["10.10.0.2"]}`),
			server: httptest.NewServer(http.HandlerFunc(func(_ http.ResponseWriter, _ *http.Request) {
			})),

			expectedErr: errors.New("Wrong Primary MDM IP, Please provide valid Primary MDM IP"),
		},
		"non 200 status code": {
			mdmTopologyParam: []byte(`{"mdmUser": "admin", "mdmPassword": "password", "mdmIps": ["10.10.0.1"]}`),
			server: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				w.WriteHeader(http.StatusBadRequest) // Simulate a non-200 status code
				response := types.GatewayResponse{
					Message: "Bad Request",
				}
				_ = json.NewEncoder(w).Encode(response)
			})),
			expectedResponse: &types.GatewayResponse{
				Message: "Bad Request",
			},
			expectedErr: nil,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			defer tt.server.Close()
			defer after()

			gc := &GatewayClient{
				http:     &http.Client{},
				host:     tt.server.URL,
				username: "test_username",
				password: "test_password",
				version:  tt.version,
			}
			if tt.setup != nil {
				tt.setup()
			}

			res, err := gc.ValidateMDMDetails(tt.mdmTopologyParam)

			if tt.expectedErr != nil {
				assert.EqualError(t, err, tt.expectedErr.Error())
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, res)
				assert.Equal(t, tt.expectedResponse, res)
			}
		})
	}
}

func TestGetClusterDetails(t *testing.T) {
	defaultCookiesFunc := setCookieFunc
	after := func() {
		setCookieFunc = defaultCookiesFunc
	}
	tests := map[string]struct {
		mdmTopologyParam   []byte
		requireJSONOutput  bool
		server             *httptest.Server
		version            string
		expectedErr        error
		expectedStatusCode int
		expectedResponse   *types.GatewayResponse
		setup              func()
	}{
		"success with version 4.0": {
			mdmTopologyParam:  []byte(`{"mdmIps": ["192.168.0.1"]}`),
			requireJSONOutput: false,
			server: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				resp := types.MDMTopologyDetails{
					SdcIps: []string{"10.0.0.1", "10.0.0.2"},
				}
				data, err := json.Marshal(resp)
				if err != nil {
					t.Fatal(err)
				}
				w.WriteHeader(http.StatusOK)
				w.Write(data)
			})),
			version:            "4.0",
			expectedErr:        nil,
			expectedStatusCode: http.StatusOK,
			expectedResponse: &types.GatewayResponse{
				StatusCode: 200,
				ClusterDetails: types.MDMTopologyDetails{
					SdcIps: []string{"10.0.0.1", "10.0.0.2"},
				},
			},
		},
		"error - set cookies": {
			mdmTopologyParam:  []byte(`{"mdmIps": ["192.168.0.1"]}`),
			requireJSONOutput: false,
			server: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				resp := types.MDMTopologyDetails{
					SdcIps: []string{"10.0.0.1", "10.0.0.2"},
				}
				data, err := json.Marshal(resp)
				if err != nil {
					t.Fatal(err)
				}
				w.WriteHeader(http.StatusOK)
				w.Write(data)
			})),
			version:            "4.0",
			expectedStatusCode: http.StatusOK,
			expectedErr:        errors.New("Error While Handling Cookie: Cookie error"),
			setup: func() {
				setCookieFunc = func(_ http.Header, host string) error {
					return errors.New("Cookie error")
				}
			},
		},
		"error getting cluster details": {
			mdmTopologyParam:   []byte(`{"invalid": "data"}`),
			requireJSONOutput:  false,
			expectedErr:        errors.New("Error Getting Cluster Details"),
			expectedStatusCode: http.StatusBadRequest,
			server: httptest.NewServer(http.HandlerFunc(func(_ http.ResponseWriter, _ *http.Request) {
				// no response
			})),
			expectedResponse: &types.GatewayResponse{
				StatusCode:     200,
				ClusterDetails: types.MDMTopologyDetails{},
			},
		},
		"non 200 status code": {
			mdmTopologyParam: []byte(`{"mdmUser": "admin", "mdmPassword": "password", "mdmIps": ["192.168.0.1"]}`),
			server: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				w.WriteHeader(http.StatusBadRequest) // Simulate a non-200 status code
				response := types.GatewayResponse{
					Message: "Bad Request",
				}
				_ = json.NewEncoder(w).Encode(response)
			})),
			expectedResponse: &types.GatewayResponse{
				Message: "Bad Request",
			},
			expectedErr: nil,
		},
		"success with JSON output": {
			mdmTopologyParam:  []byte(`{"mdmUser": "admin", "mdmPassword": "password", "mdmIps": ["10.10.0.1"]}`),
			requireJSONOutput: true,
			server: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				resp := `{"sdcIps":["10.0.0.1","10.0.0.2"]}`
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(resp))
			})),
			expectedResponse: &types.GatewayResponse{
				StatusCode: 200,
				Data:       `{"sdcIps":["10.0.0.1","10.0.0.2"]}`,
			},
			version:            "4.0",
			expectedErr:        nil,
			expectedStatusCode: http.StatusOK,
		},
		"success with structured output": {
			mdmTopologyParam:  []byte(`{"mdmUser": "admin", "mdmPassword": "password", "mdmIps": ["10.10.0.1"]}`),
			requireJSONOutput: false,
			server: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				resp := types.MDMTopologyDetails{
					SdcIps: []string{"10.0.0.1", "10.0.0.2"},
				}
				data, err := json.Marshal(resp)
				if err != nil {
					t.Fatal(err)
				}
				w.WriteHeader(http.StatusOK)
				w.Write(data)
			})),
			expectedResponse: &types.GatewayResponse{
				StatusCode: 200,
				ClusterDetails: types.MDMTopologyDetails{
					SdcIps: []string{"10.0.0.1", "10.0.0.2"},
				},
			},
			version:            "4.0",
			expectedErr:        nil,
			expectedStatusCode: http.StatusOK,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			defer tt.server.Close()
			defer after()

			gc := &GatewayClient{
				http:     &http.Client{},
				host:     tt.server.URL,
				username: "test_username",
				password: "test_password",
				version:  tt.version,
			}
			if tt.setup != nil {
				tt.setup()
			}

			res, err := gc.GetClusterDetails(tt.mdmTopologyParam, tt.requireJSONOutput)

			if tt.expectedErr != nil {
				assert.EqualError(t, err, tt.expectedErr.Error())
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, res)
				assert.Equal(t, tt.expectedStatusCode, res.StatusCode)
				assert.Equal(t, tt.expectedResponse, res)
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
			err := ParseJSONError(tt.response)
			if !reflect.DeepEqual(err, tt.expectedErr) {
				t.Errorf("Expected error %v, got %v", tt.expectedErr, err)
			}
		})
	}
}

func TestGatewayClient_NewTokenGeneration(t *testing.T) {
	// create mock HTTP servers for various needs
	successServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost && r.URL.Path == "/rest/auth/login" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			fmt.Fprintln(w, `{"access_token":"mock_access_token"}`)
			return
		}
		http.NotFound(w, r)
	}))

	failServerLogin := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost && r.URL.Path == "/rest/auth/login" {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		http.NotFound(w, r)
	}))

	closedServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost && r.URL.Path == "/rest/auth/login" {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		http.NotFound(w, r)
	}))
	closedServer.Close()

	defer func() {
		successServer.Close()
		failServerLogin.Close()
	}()
	type fields struct {
		http     *http.Client
		host     string
		username string
		password string
		token    string
		version  string
		insecure bool
	}
	tests := []struct {
		name    string
		fields  fields
		want    string
		wantErr bool
	}{
		{
			name: "success",
			fields: fields{
				http:     &http.Client{},
				host:     successServer.URL,
				username: "admin",
				password: "password",
				token:    "",
				version:  "4.0",
				insecure: true,
			},
			want:    "mock_access_token",
			wantErr: false,
		},
		{
			name: "login error",
			fields: fields{
				http:     &http.Client{},
				host:     failServerLogin.URL,
				username: "admin",
				password: "password",
				token:    "",
				version:  "4.0",
				insecure: true,
			},
			want:    "",
			wantErr: true,
		},
		{
			name: "bad request",
			fields: fields{
				http:     &http.Client{},
				host:     failServerLogin.URL + "*?",
				username: "admin",
				password: "password",
				token:    "",
				version:  "4.0",
				insecure: true,
			},
			want:    "",
			wantErr: true,
		},
		{
			name: "closed server",
			fields: fields{
				http:     &http.Client{},
				host:     closedServer.URL,
				username: "admin",
				password: "password",
				token:    "",
				version:  "4.0",
				insecure: true,
			},
			want:    "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gc := &GatewayClient{
				http:     tt.fields.http,
				host:     tt.fields.host,
				username: tt.fields.username,
				password: tt.fields.password,
				token:    tt.fields.token,
				version:  tt.fields.version,
				insecure: tt.fields.insecure,
			}
			got, err := gc.NewTokenGeneration()
			if (err != nil) != tt.wantErr {
				t.Errorf("GatewayClient.NewTokenGeneration() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GatewayClient.NewTokenGeneration() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGatewayClient_GetInQueueCommand(t *testing.T) {
	defaultSetCookieFunc := setCookieFunc
	after := func() {
		setCookieFunc = defaultSetCookieFunc
	}
	successServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost && r.URL.Path == "/rest/auth/login" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			fmt.Fprintln(w, `{"access_token":"mock_access_token"}`)
			return
		}
		http.NotFound(w, r)
	}))
	type fields struct {
		http     *http.Client
		host     string
		username string
		password string
		token    string
		version  string
		insecure bool
	}
	tests := []struct {
		name    string
		fields  fields
		want    []types.MDMQueueCommandDetails
		wantErr bool
		setup   func()
	}{
		{name: "success case",
			fields: fields{
				http:     &http.Client{},
				host:     successServer.URL,
				username: "admin",
				password: "password",
				token:    "",
				version:  "4.0",
				insecure: true,
			},
			want:    []types.MDMQueueCommandDetails{},
			wantErr: false,
		},
		{name: "fail - set cookies",
			fields: fields{
				http:     &http.Client{},
				host:     successServer.URL,
				username: "admin",
				password: "password",
				token:    "",
				version:  "4.0",
				insecure: true,
			},
			want:    nil,
			wantErr: true,
			setup: func() {
				setCookieFunc = func(_ http.Header, _ string) error {
					return errors.New("cookie error")
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer after()
			gc := &GatewayClient{
				http:     tt.fields.http,
				host:     tt.fields.host,
				username: tt.fields.username,
				password: tt.fields.password,
				token:    tt.fields.token,
				version:  tt.fields.version,
				insecure: tt.fields.insecure,
			}
			if tt.setup != nil {
				tt.setup()
			}
			_, err := gc.GetInQueueCommand()
			if (err != nil) != tt.wantErr {
				t.Errorf("GatewayClient.GetInQueueCommand() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
