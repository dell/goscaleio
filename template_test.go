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
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestGetTemplateByID(t *testing.T) {
	responseJSONFile := "response/template_response.json"
	responseData, err := ioutil.ReadFile(responseJSONFile)
	if err != nil {
		t.Fatalf("Failed to read response JSON file: %v", err)
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.Contains(r.URL.Path, "/Api/V1/template/") {
			if r.Method == http.MethodGet {
				w.WriteHeader(http.StatusOK)
				w.Write(responseData)
				return
			}
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

	templateID := "453c41eb-d72a-4ed1-ad16-bacdffbdd766"

	templateResponse, err := gc.GetTemplateByID(templateID)

	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if templateResponse == nil {
		t.Error("Template response is nil")
	}

}

func TestGetTemplateByFilters(t *testing.T) {
	responseJSONFile := "response/templates_response.json"
	responseData, err := ioutil.ReadFile(responseJSONFile)
	if err != nil {
		t.Fatalf("Failed to read response JSON file: %v", err)
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.Contains(r.URL.Path, "/Api/V1/template") {
			if r.Method == http.MethodGet {
				w.WriteHeader(http.StatusOK)
				w.Write(responseData)
				return
			}
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

	filter := "name"
	value := "Test"

	templateResponse, err := gc.GetTemplateByFilters(filter, value)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if templateResponse == nil {
		t.Error("Template response is nil")
	}
}

func TestGetAllTemplates(t *testing.T) {
	responseJSONFile := "response/templates_response.json"
	responseData, err := ioutil.ReadFile(responseJSONFile)
	if err != nil {
		t.Fatalf("Failed to read response JSON file: %v", err)
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.Contains(r.URL.Path, "/Api/V1/template") {
			if r.Method == http.MethodGet {
				w.WriteHeader(http.StatusOK)
				w.Write(responseData)
				return
			}
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

	templateResponse, err := gc.GetAllTemplates()
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if templateResponse == nil {
		t.Error("Template response is nil")
	}
}
