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

func TestDeployService(t *testing.T) {

	firmwareResponse := `{ "id": "67890", "name": "PowerFlex 4.5.0.0", "sourceLocation": "PowerFlex_Software_4.5.0.0_287_r1.zip", "sourceType": null, "diskLocation": "/opt/Dell/ASM/temp/RCM_8aaaee188f38ea00018f3d4dc8ea0075/catalog", "filename": "catalog.xml", "md5Hash": null, "username": null, "password": null, "downloadStatus": "error", "createdDate": "2024-05-03T07:14:18.986+00:00", "createdBy": "admin", "updatedDate": "2024-05-06T05:59:33.696+00:00", "updatedBy": "system", "defaultCatalog": false, "embedded": false, "state": "errors", "softwareComponents": [], "softwareBundles": [], "deployments": [], "bundleCount": 0, "componentCount": 0, "userBundleCount": 0, "minimal": true, "downloadProgress": 100, "extractProgress": 0, "fileSizeInGigabytes": 4.6, "signedKeySourceLocation": null, "signature": "Unsigned", "custom": false, "needsAttention": false, "jobId": "Job-2cf0b7b7-c794-4fa4-9256-784c261ebbc9", "rcmapproved": false }`

	serviceTemplateJSONFile := "response/service_template_response.json"
	serviceTemplateResponse, err := ioutil.ReadFile(serviceTemplateJSONFile)
	if err != nil {
		t.Fatalf("Failed to read response JSON file: %v", err)
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.Contains(r.URL.Path, "/Api/V1/FirmwareRepository/") {
			w.WriteHeader(http.StatusOK)
			_, err := w.Write([]byte(firmwareResponse))
			if err != nil {
				t.Fatalf("Error writing response: %v", err)
			}
			return
		} else if strings.Contains(r.URL.Path, "/Api/V1/ServiceTemplate/") {
			w.WriteHeader(http.StatusOK)
			_, err := w.Write([]byte(serviceTemplateResponse))
			if err != nil {
				t.Fatalf("Error writing response: %v", err)
			}
			return
		} else if strings.Contains(r.URL.Path, "/Api/V1/Deployment") {
			w.WriteHeader(http.StatusOK)
			responseJSON := `{"StatusCode":200,"Messages":[{"DisplayMessage":"Service deployed successfully"}]}`
			w.Write([]byte(responseJSON))
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

	deploymentName := "Test Deployment"
	deploymentDesc := "Test Deployment Description"
	serviceTemplateID := "12345"
	firmwareRepositoryID := "67890"
	nodes := "3"

	serviceResponse, err := gc.DeployService(deploymentName, deploymentDesc, serviceTemplateID, firmwareRepositoryID, nodes)

	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if serviceResponse == nil {
		t.Error("Service response is nil")
	}

	if serviceResponse.StatusCode != 200 {
		t.Errorf("Expected status code 200, got %d", serviceResponse.StatusCode)
	}

	expectedMessage := "Service deployed successfully"
	if serviceResponse.Messages[0].DisplayMessage != expectedMessage {
		t.Errorf("Expected message '%s', got '%s'", expectedMessage, serviceResponse.Messages[0].DisplayMessage)
	}
}

func TestUpdateService(t *testing.T) {
	responseJSONFile := "response/update_service_response.json"
	responseData, err := ioutil.ReadFile(responseJSONFile)
	if err != nil {
		t.Fatalf("Failed to read response JSON file: %v", err)
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.Contains(r.URL.Path, "/Api/V1/Deployment/") {
			if r.Method == http.MethodGet {
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(responseData))
				return
			} else if r.Method == http.MethodPut {
				w.WriteHeader(http.StatusOK)
				responseJSON := `{"StatusCode":200,"Messages":[{"DisplayMessage":"Service updated successfully"}]}`
				w.Write([]byte(responseJSON))
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

	deploymentID := "12345"
	deploymentName := "Updated Deployment"
	deploymentDesc := "Updated Deployment Description"
	nodes := "4"
	nodename := "pfmc-k8s-20230809-160-1"

	serviceResponse, err := gc.UpdateService(deploymentID, deploymentName, deploymentDesc, nodes, nodename)

	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if serviceResponse == nil {
		t.Error("Service response is nil")
	}

	if serviceResponse.StatusCode != 200 {
		t.Errorf("Expected status code 200, got %d", serviceResponse.StatusCode)
	}

	expectedMessage := "Service updated successfully"
	if serviceResponse.Messages[0].DisplayMessage != expectedMessage {
		t.Errorf("Expected message '%s', got '%s'", expectedMessage, serviceResponse.Messages[0].DisplayMessage)
	}
}

func TestGetServiceDetailsByID(t *testing.T) {
	responseJSONFile := "response/update_service_response.json"
	responseData, err := ioutil.ReadFile(responseJSONFile)
	if err != nil {
		t.Fatalf("Failed to read response JSON file: %v", err)
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/rest/auth/login" && r.Method == http.MethodPost {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"access_token": "mock_access_token"}`))
			return
		} else if strings.Contains(r.URL.Path, "/Api/V1/Deployment/") && r.Method == http.MethodGet {
			w.WriteHeader(http.StatusOK)
			w.Write(responseData)
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

	deploymentID := "12345"
	newToken := true

	serviceResponse, err := gc.GetServiceDetailsByID(deploymentID, newToken)

	if err != nil {
		t.Fatalf("Error while getting service details: %v", err)
	}

	if serviceResponse == nil {
		t.Fatalf("Expected non-nil response, got nil")
	}
}

func TestGetServiceDetailsByFilter(t *testing.T) {
	responseJSONFile := "response/services_response.json"
	responseData, err := ioutil.ReadFile(responseJSONFile)
	if err != nil {
		t.Fatalf("Failed to read response JSON file: %v", err)
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.Contains(r.URL.Path, "/Api/V1/Deployment") {
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
	value := "TestCreate"

	serviceResponse, err := gc.GetServiceDetailsByFilter(filter, value)

	if err != nil {
		t.Fatalf("Error while getting service details: %v", err)
	}

	if serviceResponse == nil {
		t.Fatalf("Expected non-nil response, got nil")
	}
}

func TestGetAllServiceDetails(t *testing.T) {
	responseJSONFile := "response/services_response.json"
	responseData, err := ioutil.ReadFile(responseJSONFile)
	if err != nil {
		t.Fatalf("Failed to read response JSON file: %v", err)
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.Contains(r.URL.Path, "/Api/V1/Deployment") {
			if r.Method == http.MethodGet {
				w.WriteHeader(http.StatusOK)
				w.Write(responseData)
				return
			}
		}
		http.NotFound(w, r)
	}))
	defer server.Close()

	// Creating a GatewayClient with the mocked server's URL
	gc := &GatewayClient{
		http:     &http.Client{},
		host:     server.URL,
		username: "test_username",
		password: "test_password",
	}

	serviceResponse, err := gc.GetAllServiceDetails()

	if err != nil {
		t.Fatalf("Error while getting service details: %v", err)
	}

	if serviceResponse == nil {
		t.Fatalf("Expected non-nil response, got nil")
	}
}
