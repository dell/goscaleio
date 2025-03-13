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

package goscaleio

import (
	"fmt"
	"io/ioutil"
	"math"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetTemplateByID(t *testing.T) {
	tests := map[string]struct {
		id       string
		server   *httptest.Server
		version  string
		expected error
	}{
		"error due to parsing response": {
			id: "12345",
			server: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				w.WriteHeader(http.StatusOK)
				w.Write([]byte("invalid json"))
			})),
			expected: fmt.Errorf("Error parsing response data for template: invalid character 'i' looking for beginning of value"),
		},
		"success version 4.0": {
			id: "12345",
			server: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				responseJSONFile := "response/template_response.json"
				responseData, err := ioutil.ReadFile(responseJSONFile)
				if err != nil {
					t.Fatalf("Failed to read response JSON file: %v", err)
				}
				if strings.Contains(r.URL.Path, "/Api/V1/template/") {
					if r.Method == http.MethodGet {
						w.WriteHeader(http.StatusOK)
						w.Write(responseData)
						return
					}
				}
				http.NotFound(w, r)
			})),
			version:  "4.0",
			expected: nil,
		},
		"success version < 4.0": {
			id: "12345",
			server: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				responseJSONFile := "response/template_response.json"
				responseData, err := ioutil.ReadFile(responseJSONFile)
				if err != nil {
					t.Fatalf("Failed to read response JSON file: %v", err)
				}
				if strings.Contains(r.URL.Path, "/Api/V1/template/") {
					if r.Method == http.MethodGet {
						w.WriteHeader(http.StatusOK)
						w.Write(responseData)
						return
					}
				}
				http.NotFound(w, r)
			})),
			version:  "3.0",
			expected: nil,
		},
		"error due to template not found": {
			id: "12345",
			server: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				http.NotFound(w, r)
			})),
			expected: fmt.Errorf("Template not found"),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			defer tc.server.Close()
			gc := &GatewayClient{
				http:     &http.Client{},
				host:     tc.server.URL,
				username: "test_username",
				password: "test_password",
				version:  tc.version,
			}

			template, err := gc.GetTemplateByID(tc.id)

			if tc.expected == nil {
				assert.Nil(t, err)
				assert.NotNil(t, template)
			} else {
				assert.NotNil(t, err)
				assert.Equal(t, tc.expected.Error(), err.Error())
			}
		})
	}
}

func TestGetTemplateByFilters(t *testing.T) {
	tests := map[string]struct {
		server   *httptest.Server
		version  string
		expected error
	}{
		"success with version 4.0": {
			server: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				responseJSONFile := "response/templates_response.json"
				responseData, err := ioutil.ReadFile(responseJSONFile)
				if err != nil {
					t.Fatalf("Failed to read response JSON file: %v", err)
				}
				if strings.Contains(r.URL.Path, "/Api/V1/template") {
					if r.Method == http.MethodGet {
						w.WriteHeader(http.StatusOK)
						w.Write(responseData)
						return
					}
				}
				http.NotFound(w, r)
			})),
			version:  "4.0",
			expected: nil,
		},
		"success with version 3.0": {
			server: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				responseJSONFile := "response/templates_response.json"
				responseData, err := ioutil.ReadFile(responseJSONFile)
				if err != nil {
					t.Fatalf("Failed to read response JSON file: %v", err)
				}
				if strings.Contains(r.URL.Path, "/Api/V1/template") {
					if r.Method == http.MethodGet {
						w.WriteHeader(http.StatusOK)
						w.Write(responseData)
						return
					}
				}
				http.NotFound(w, r)
			})),
			version:  "3.0",
			expected: nil,
		},
		"error due to parsing response": {
			server: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				w.WriteHeader(http.StatusOK)
				w.Write([]byte("invalid json"))
			})),
			expected: fmt.Errorf("Error While Parsing Response Data For Template: invalid character 'i' looking for beginning of value"),
		},
		"error due to template not found": {
			server: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				http.NotFound(w, r)
			})),
			expected: fmt.Errorf("Template not found"),
		},
		"error due to template details is empty": {
			server: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				w.WriteHeader(http.StatusOK)
				w.Write([]byte("{}"))
			})),
			expected: fmt.Errorf("Template not found"),
		},
	}

	for _, tc := range tests {
		t.Run("", func(t *testing.T) {
			defer tc.server.Close()
			gc := &GatewayClient{
				http:     &http.Client{},
				host:     tc.server.URL,
				username: "test_username",
				password: "test_password",
				version:  tc.version,
			}

			template, err := gc.GetTemplateByFilters("name", "template1")

			if tc.expected == nil {
				assert.Nil(t, err)
				assert.NotNil(t, template)
			} else {
				assert.NotNil(t, err)
				assert.Equal(t, tc.expected.Error(), err.Error())
			}
		})
	}
}

func TestGetAllTemplates(t *testing.T) {
	tests := map[string]struct {
		server   *httptest.Server
		version  string
		expected error
	}{
		"success with version 4.0": {
			server: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				responseJSONFile := "response/templates_response.json"
				responseData, err := ioutil.ReadFile(responseJSONFile)
				if err != nil {
					t.Fatalf("Failed to read response JSON file: %v", err)
				}
				if strings.Contains(r.URL.Path, "/Api/V1/template") {
					if r.Method == http.MethodGet {
						w.WriteHeader(http.StatusOK)
						w.Write(responseData)
						return
					}
				}
				http.NotFound(w, r)
			})),
			version:  "4.0",
			expected: nil,
		},
		"success with version 3.0": {
			server: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				responseJSONFile := "response/templates_response.json"
				responseData, err := ioutil.ReadFile(responseJSONFile)
				if err != nil {
					t.Fatalf("Failed to read response JSON file: %v", err)
				}
				if strings.Contains(r.URL.Path, "/Api/V1/template") {
					if r.Method == http.MethodGet {
						w.WriteHeader(http.StatusOK)
						w.Write(responseData)
						return
					}
				}
				http.NotFound(w, r)
			})),
			version:  "3.0",
			expected: nil,
		},
		"error due to parsing response": {
			server: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				w.WriteHeader(http.StatusOK)
				w.Write([]byte("invalid json"))
			})),
			expected: fmt.Errorf("Error While Parsing Response Data For Template: invalid character 'i' looking for beginning of value"),
		},
	}

	for _, tc := range tests {
		t.Run("", func(t *testing.T) {
			defer tc.server.Close()
			gc := &GatewayClient{
				http:     &http.Client{},
				host:     tc.server.URL,
				username: "test_username",
				password: "test_password",
				version:  tc.version,
			}

			template, err := gc.GetAllTemplates()

			if tc.expected == nil {
				assert.Nil(t, err)
				assert.NotNil(t, template)
			} else {
				assert.NotNil(t, err)
				assert.Equal(t, tc.expected.Error(), err.Error())
			}
		})
	}
}

func TestCloneTemplate(t *testing.T) {
	tests := []struct {
		name         string
		originID     string
		templateName string
		expectedErr  error
	}{
		{
			name:         "success",
			originID:     "1234567",
			templateName: "Test-Copy",
			expectedErr:  nil,
		},
		{
			name:         "error due to cloning a non-existant template",
			originID:     "template_id_does_not_exist",
			templateName: "Test-Copy",
			expectedErr:  fmt.Errorf("Error While Cloning Template: Template not found"),
		},
		{
			name:         "error due to cloning an existing template",
			originID:     "1234567",
			templateName: "Test",
			expectedErr:  fmt.Errorf("Error While Cloning Template: Template already exists please use a different name"),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			encodedValue := url.QueryEscape("1234567")
			// Create a mock server
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method == http.MethodPost && r.URL.Path == "/Api/V1/ServiceTemplate/cloneTemplate" && tc.templateName == "Test-Copy" {
					w.WriteHeader(http.StatusOK)
					return
				}
				responseJSONFile := "response/templates_response.json"
				responseData, err := ioutil.ReadFile(responseJSONFile)
				if err != nil {
					t.Fatalf("Failed to read response JSON file: %v", err)
				}
				if strings.Contains(r.URL.Path, "/Api/V1/template") && strings.Contains(r.URL.RawQuery, encodedValue) {
					if r.Method == http.MethodGet {
						w.WriteHeader(http.StatusOK)
						w.Write(responseData)
						return
					}
				}
				http.NotFound(w, r)
			}))
			defer server.Close()

			client, err := NewClientWithArgs(server.URL, "4.0", math.MaxInt64, true, false)
			if err != nil {
				t.Fatal(err)
			}
			s := System{
				client: client,
			}
			gc := &GatewayClient{
				http:     &http.Client{},
				host:     server.URL,
				username: "test_username",
				password: "test_password",
				version:  "4.0",
			}

			err = gc.CloneTemplate(&s, tc.originID, tc.templateName)

			if tc.expectedErr == nil {
				assert.Nil(t, err)
			} else {
				assert.NotNil(t, err)
				assert.Equal(t, tc.expectedErr.Error(), err.Error())
			}
		})
	}
}
