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
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	types "github.com/dell/goscaleio/types/v1"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
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
	assert.NotNil(t, serviceResponse, "Expected non-nil response")
	assert.Equal(t, 200, serviceResponse.StatusCode, "Expected status code 200")
	assert.Equal(t, "Service deployed successfully", serviceResponse.Messages[0].DisplayMessage, "Expected message 'Service deployed successfully'")
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

	assert.NotNil(t, serviceResponse, "Expected non-nil response")
	assert.Equal(t, 200, serviceResponse.StatusCode, "Expected status code 200")
	assert.Equal(t, "Service updated successfully", serviceResponse.Messages[0].DisplayMessage, "Expected message 'Service updated successfully'")
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

	assert.NotNil(t, serviceResponse, "Expected non-nil response")
	assert.EqualValues(t, serviceResponse.ID, "12345")
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

	assert.NotNil(t, serviceResponse, "Expected non-nil response")
	assert.EqualValues(t, serviceResponse[0].DeploymentName, "TestCreate")
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

	assert.NotNil(t, serviceResponse, "Expected non-nil response")
	assert.EqualValues(t, serviceResponse[0].DeploymentName, "TestCreate")
}

func TestDeleteService(t *testing.T) {
	serviceID := uuid.NewString()
	serversInInventory := "myService"
	serversManagedState := "myServiceState"

	type testCase struct {
		version     string
		server      *httptest.Server
		expectedErr error
	}

	cases := map[string]testCase{
		"succeed: 3.7 version": {
			version: "3.7",
			server: httptest.NewServer(http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
				switch req.RequestURI {
				case fmt.Sprintf("/Api/V1/Deployment/%s?serversInInventory=%s&serversManagedState=%s", serviceID, serversInInventory, serversManagedState):
					resp.WriteHeader(http.StatusNoContent)
				default:
					resp.WriteHeader(http.StatusBadRequest)
					resp.Write([]byte(`{"message":"no route handled","httpStatusCode":400,"errorCode":0}`))
				}
			})),
			expectedErr: nil,
		},
		"succeed: 4.0 version": {
			version: "4.0",
			server: httptest.NewServer(http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
				switch req.RequestURI {
				case fmt.Sprintf("/Api/V1/Deployment/%s?serversInInventory=%s&serversManagedState=%s", serviceID, serversInInventory, serversManagedState):
					resp.WriteHeader(http.StatusNoContent)
				default:
					resp.WriteHeader(http.StatusBadRequest)
					resp.Write([]byte(`{"message":"no route handled","httpStatusCode":400,"errorCode":0}`))
				}
			})),
			expectedErr: nil,
		},
		"error: couldn't delete service": {
			version: "3.7",
			server: httptest.NewServer(http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
				resp.WriteHeader(http.StatusBadRequest)
				resp.Write([]byte(`{"message":"no route handled","httpStatusCode":400,"errorCode":0}`))
			})),
			expectedErr: fmt.Errorf("Couldn't delete service"),
		},
	}

	for _, tc := range cases {
		gc := &GatewayClient{
			http:     &http.Client{},
			host:     tc.server.URL,
			username: "test_username",
			password: "test_password",
			version:  tc.version,
		}

		_, err := gc.DeleteService(serviceID, serversInInventory, serversManagedState)
		if err != nil {
			if tc.expectedErr.Error() != err.Error() {
				t.Fatal(err)
			}
		}
	}
}

func TestGetServiceComplianceDetails(t *testing.T) {
	deploymentID := uuid.NewString()

	type testCase struct {
		version     string
		server      *httptest.Server
		expectedErr error
	}

	cases := map[string]testCase{
		"succeed: 3.7 version": {
			version: "3.7",
			server: httptest.NewServer(http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
				switch req.RequestURI {
				case fmt.Sprintf("/Api/V1/Deployment/%s/firmware/compliancereport", deploymentID):
					content, err := json.Marshal([]types.ComplianceReport{
						{
							ID:        uuid.NewString(),
							Compliant: true,
						},
						{
							ID:        uuid.NewString(),
							Compliant: false,
						},
					})
					if err != nil {
						t.Fatal(err)
					}

					resp.Write(content)
					resp.WriteHeader(http.StatusOK)
				default:
					resp.WriteHeader(http.StatusBadRequest)
					resp.Write([]byte(`{"message":"no route handled","httpStatusCode":400,"errorCode":0}`))
				}
			})),
			expectedErr: nil,
		},
		"succeed: 4.0 version": {
			version: "4.0",
			server: httptest.NewServer(http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
				switch req.RequestURI {
				case fmt.Sprintf("/Api/V1/Deployment/%s/firmware/compliancereport", deploymentID):
					content, err := json.Marshal([]types.ComplianceReport{
						{
							ID:        uuid.NewString(),
							Compliant: true,
						},
						{
							ID:        uuid.NewString(),
							Compliant: false,
						},
					})
					if err != nil {
						t.Fatal(err)
					}

					resp.Write(content)
					resp.WriteHeader(http.StatusOK)
				default:
					resp.WriteHeader(http.StatusBadRequest)
					resp.Write([]byte(`{"message":"no route handled","httpStatusCode":400,"errorCode":0}`))
				}
			})),
			expectedErr: nil,
		},
		"error: error parsing response": {
			version: "3.7",
			server: httptest.NewServer(http.HandlerFunc(func(resp http.ResponseWriter, _ *http.Request) {
				resp.WriteHeader(http.StatusBadRequest)
				resp.Write([]byte(`{"message":"bad request","httpStatusCode":400,"errorCode":0}`))
			})),
			expectedErr: fmt.Errorf("Couldn't find compliance report for given deployment"),
		},
		"error: couldn't find compliance report": {
			version: "3.7",
			server: httptest.NewServer(http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
				switch req.RequestURI {
				case fmt.Sprintf("/Api/V1/Deployment/%s/firmware/compliancereport", deploymentID):
					resp.WriteHeader(http.StatusOK)
				default:
					resp.WriteHeader(http.StatusBadRequest)
					resp.Write([]byte(`{"message":"no route handled","httpStatusCode":400,"errorCode":0}`))
				}
			})),
			expectedErr: fmt.Errorf("Error while parsing response data for compliance report: unexpected end of JSON input"),
		},
		"error: empty compliance report": {
			version: "3.7",
			server: httptest.NewServer(http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
				switch req.RequestURI {
				case fmt.Sprintf("/Api/V1/Deployment/%s/firmware/compliancereport", deploymentID):
					content, err := json.Marshal([]types.ComplianceReport{})
					if err != nil {
						t.Fatal(err)
					}

					resp.Write(content)
					resp.WriteHeader(http.StatusOK)
				default:
					resp.WriteHeader(http.StatusBadRequest)
					resp.Write([]byte(`{"message":"no route handled","httpStatusCode":400,"errorCode":0}`))
				}
			})),
			expectedErr: fmt.Errorf("Couldn't find compliance report for given deployment"),
		},
	}

	for _, tc := range cases {
		gc := &GatewayClient{
			http:     &http.Client{},
			host:     tc.server.URL,
			username: "test_username",
			password: "test_password",
			version:  tc.version,
		}

		_, err := gc.GetServiceComplianceDetails(deploymentID)
		if err != nil {
			if tc.expectedErr.Error() != err.Error() {
				t.Fatal(err)
			}
		}
	}
}

func TestGetServiceComplianceDetailsByFilter(t *testing.T) {
	deploymentID := uuid.NewString()
	complianceID := uuid.NewString()
	complianceReports := []types.ComplianceReport{
		{
			ID:         complianceID,
			IPAddress:  "127.0.0.1",
			Compliant:  true,
			ServiceTag: "myServiceTag",
			HostName:   "myHostName",
		},
	}

	type testCase struct {
		version     string
		filter      string
		value       string
		server      *httptest.Server
		expectedErr error
	}

	cases := map[string]testCase{
		"success: IP Address": {
			version: "3.7",
			filter:  "IpAddress",
			value:   "127.0.0.1",
			server: httptest.NewServer(http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
				content, err := json.Marshal(complianceReports)
				if err != nil {
					t.Fatal(err)
				}

				resp.Write(content)
				resp.WriteHeader(http.StatusOK)
			})),
			expectedErr: nil,
		},
		"success: Service Tag": {
			version: "3.7",
			filter:  "ServiceTag",
			value:   "myServiceTag",
			server: httptest.NewServer(http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
				content, err := json.Marshal(complianceReports)
				if err != nil {
					t.Fatal(err)
				}

				resp.Write(content)
				resp.WriteHeader(http.StatusOK)
			})),
			expectedErr: nil,
		},
		"success: Compliant": {
			version: "3.7",
			filter:  "Compliant",
			value:   "true",
			server: httptest.NewServer(http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
				content, err := json.Marshal(complianceReports)
				if err != nil {
					t.Fatal(err)
				}

				resp.Write(content)
				resp.WriteHeader(http.StatusOK)
			})),
			expectedErr: nil,
		},
		"success: Host Name": {
			version: "3.7",
			filter:  "HostName",
			value:   "myHostName",
			server: httptest.NewServer(http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
				content, err := json.Marshal(complianceReports)
				if err != nil {
					t.Fatal(err)
				}

				resp.Write(content)
				resp.WriteHeader(http.StatusOK)
			})),
			expectedErr: nil,
		},
		"success: ID": {
			version: "3.7",
			filter:  "ID",
			value:   complianceID,
			server: httptest.NewServer(http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
				content, err := json.Marshal(complianceReports)
				if err != nil {
					t.Fatal(err)
				}

				resp.Write(content)
				resp.WriteHeader(http.StatusOK)
			})),
			expectedErr: nil,
		},
		"error: empty compliance report": {
			version: "3.7",
			filter:  "ID",
			value:   complianceID,
			server: httptest.NewServer(http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
				content, err := json.Marshal([]types.ComplianceReport{})
				if err != nil {
					t.Fatal(err)
				}

				resp.Write(content)
				resp.WriteHeader(http.StatusOK)
			})),
			expectedErr: fmt.Errorf("Couldn't find compliance report for the given deployment"),
		},
		"error: invalid filter": {
			version: "3.7",
			filter:  "InvalidFilter",
			value:   "InvalidValue",
			server: httptest.NewServer(http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
				content, err := json.Marshal(complianceReports)
				if err != nil {
					t.Fatal(err)
				}

				resp.Write(content)
				resp.WriteHeader(http.StatusOK)
			})),
			expectedErr: fmt.Errorf("Invalid filter provided"),
		},
	}

	for _, tc := range cases {
		gc := &GatewayClient{
			http:     &http.Client{},
			host:     tc.server.URL,
			username: "test_username",
			password: "test_password",
			version:  tc.version,
		}

		_, err := gc.GetServiceComplianceDetailsByFilter(deploymentID, tc.filter, tc.value)
		if err != nil {
			if tc.expectedErr.Error() != err.Error() {
				t.Fatal(err)
			}
		}
	}
}
