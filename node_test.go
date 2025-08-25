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

package goscaleio

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	types "github.com/dell/goscaleio/types/v1"
)

func TestGetNodeByID(t *testing.T) {
	responseJSON := `{ "refId": "softwareOnlyServer-1.1.1.1", "refType": null, "ipAddress": "1.1.1.1", "currentIpAddress": "1.1.1.1", "serviceTag": "VMware-42 05 a8 96 26 f7 98 2c-a6 72 b9 1a 26 94 a9 9c-SW", "model": "VMware Virtual Platform", "deviceType": "SoftwareOnlyServer", "discoverDeviceType": "SOFTWAREONLYSERVER_SLES", "displayName": "pfmc-k8s-20230809-1", "managedState": "MANAGED", "state": "READY", "inUse": false, "serviceReferences": [], "statusMessage": null, "firmwareName": "Default Catalog - PowerFlex 4.5.2.0", "customFirmware": false, "needsAttention": false, "manufacturer": "VMware, Inc.", "systemId": null, "health": "NA", "healthMessage": null, "operatingSystem": "N/A", "numberOfCPUs": 0, "cpuType": null, "nics": 0, "memoryInGB": 0, "infraTemplateDate": null, "infraTemplateId": null, "serverTemplateDate": null, "serverTemplateId": null, "inventoryDate": null, "complianceCheckDate": "2024-05-08T11:16:52.951+00:00", "discoveredDate": "2024-05-08T11:16:51.805+00:00", "deviceGroupList": { "paging": null, "deviceGroup": [ { "link": null, "groupSeqId": -1, "groupName": "Global", "groupDescription": null, "createdDate": null, "createdBy": "admin", "updatedDate": null, "updatedBy": null, "managedDeviceList": null, "groupUserList": null } ] }, "detailLink": { "title": "softwareOnlyServer-1.1.1.1", "href": "/AsmManager/ManagedDevice/softwareOnlyServer-1.1.1.1", "rel": "describedby", "type": null }, "credId": "3f5869e6-6525-4dee-bb0c-fab3fe60771d", "compliance": "NONCOMPLIANT", "failuresCount": 0, "chassisId": null, "parsedFacts": null, "config": null, "hostname": "pfmc-k8s-20230809-1", "osIpAddress": null, "osAdminCredential": null, "osImageType": null, "lastJobs": null, "puppetCertName": "sles-1.1.1.1", "svmAdminCredential": null, "svmName": null, "svmIpAddress": null, "svmImageType": null, "flexosMaintMode": 0, "esxiMaintMode": 0, "vmList": [] }`

	type testCase struct {
		version     string
		server      *httptest.Server
		expectedErr error
	}

	cases := map[string]testCase{
		"succeed: 3.7 version": {
			version: "3.7",
			server: httptest.NewServer(http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
				if req.Method == http.MethodGet && strings.Contains(req.URL.Path, "/Api/V1/ManagedDevice/") {
					resp.WriteHeader(http.StatusOK)
					_, err := resp.Write([]byte(responseJSON))
					if err != nil {
						t.Fatalf("Error writing response: %v", err)
					}
					return
				}
				http.NotFound(resp, req)
			})),
			expectedErr: nil,
		},
		"succeed: 4.0 version": {
			version: "4.0",
			server: httptest.NewServer(http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
				if req.Method == http.MethodGet && strings.Contains(req.URL.Path, "/Api/V1/ManagedDevice/") {
					resp.WriteHeader(http.StatusOK)
					_, err := resp.Write([]byte(responseJSON))
					if err != nil {
						t.Fatalf("Error writing response: %v", err)
					}
					return
				}
				http.NotFound(resp, req)
			})),
			expectedErr: nil,
		},
		"error: bad request": {
			server: httptest.NewServer(http.HandlerFunc(func(resp http.ResponseWriter, _ *http.Request) {
				resp.WriteHeader(http.StatusBadRequest)
				resp.Write([]byte(`{"message":"bad request","httpStatusCode":400,"errorCode":0}`))
			})),
			expectedErr: fmt.Errorf("Couldn't find nodes with the given filter"),
		},
		"error: unable to unmarshal": {
			version: "4.0",
			server: httptest.NewServer(http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
				if req.Method == http.MethodGet && strings.Contains(req.URL.Path, "/Api/V1/ManagedDevice/") {
					resp.WriteHeader(http.StatusOK)
					return
				}
				http.NotFound(resp, req)
			})),
			expectedErr: fmt.Errorf("Error While Parsing Response Data For Node: unexpected end of JSON input"),
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			gc := &GatewayClient{
				http:     &http.Client{},
				host:     tc.server.URL,
				username: "test_username",
				password: "test_password",
				version:  tc.version,
			}
			defer tc.server.Close()

			id := "softwareOnlyServer-1.1.1.1"
			_, err := gc.GetNodeByID(id)
			if err != nil {
				if tc.expectedErr.Error() != err.Error() {
					t.Fatal(err)
				}
			}
		})
	}
}

func TestGetAllNodes(t *testing.T) {
	responseJSON := `[{ "refId": "softwareOnlyServer-1.1.1.1", "refType": null, "ipAddress": "1.1.1.1", "currentIpAddress": "1.1.1.1", "serviceTag": "VMware-42 05 a8 96 26 f7 98 2c-a6 72 b9 1a 26 94 a9 9c-SW", "model": "VMware Virtual Platform", "deviceType": "SoftwareOnlyServer", "discoverDeviceType": "SOFTWAREONLYSERVER_SLES", "displayName": "pfmc-k8s-20230809-1", "managedState": "MANAGED", "state": "READY", "inUse": false, "serviceReferences": [], "statusMessage": null, "firmwareName": "Default Catalog - PowerFlex 4.5.2.0", "customFirmware": false, "needsAttention": false, "manufacturer": "VMware, Inc.", "systemId": null, "health": "NA", "healthMessage": null, "operatingSystem": "N/A", "numberOfCPUs": 0, "cpuType": null, "nics": 0, "memoryInGB": 0, "infraTemplateDate": null, "infraTemplateId": null, "serverTemplateDate": null, "serverTemplateId": null, "inventoryDate": null, "complianceCheckDate": "2024-05-08T11:16:52.951+00:00", "discoveredDate": "2024-05-08T11:16:51.805+00:00", "deviceGroupList": { "paging": null, "deviceGroup": [ { "link": null, "groupSeqId": -1, "groupName": "Global", "groupDescription": null, "createdDate": null, "createdBy": "admin", "updatedDate": null, "updatedBy": null, "managedDeviceList": null, "groupUserList": null } ] }, "detailLink": { "title": "softwareOnlyServer-1.1.1.1", "href": "/AsmManager/ManagedDevice/softwareOnlyServer-1.1.1.1", "rel": "describedby", "type": null }, "credId": "3f5869e6-6525-4dee-bb0c-fab3fe60771d", "compliance": "NONCOMPLIANT", "failuresCount": 0, "chassisId": null, "parsedFacts": null, "config": null, "hostname": "pfmc-k8s-20230809-1", "osIpAddress": null, "osAdminCredential": null, "osImageType": null, "lastJobs": null, "puppetCertName": "sles-1.1.1.1", "svmAdminCredential": null, "svmName": null, "svmIpAddress": null, "svmImageType": null, "flexosMaintMode": 0, "esxiMaintMode": 0, "vmList": [] }]`

	type testCase struct {
		version     string
		server      *httptest.Server
		expectedErr error
	}

	cases := map[string]testCase{
		"succeed: 3.7 version": {
			version: "3.7",
			server: httptest.NewServer(http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
				if req.Method == http.MethodGet && strings.Contains(req.URL.Path, "/Api/V1/ManagedDevice") {
					resp.WriteHeader(http.StatusOK)
					_, err := resp.Write([]byte(responseJSON))
					if err != nil {
						t.Fatalf("Error writing response: %v", err)
					}
					return
				}
				http.NotFound(resp, req)
			})),
			expectedErr: nil,
		},
		"succeed: 4.0 version": {
			version: "4.0",
			server: httptest.NewServer(http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
				if req.Method == http.MethodGet && strings.Contains(req.URL.Path, "/Api/V1/ManagedDevice") {
					resp.WriteHeader(http.StatusOK)
					_, err := resp.Write([]byte(responseJSON))
					if err != nil {
						t.Fatalf("Error writing response: %v", err)
					}
					return
				}
				http.NotFound(resp, req)
			})),
			expectedErr: nil,
		},
		"error: bad request": {
			server: httptest.NewServer(http.HandlerFunc(func(resp http.ResponseWriter, _ *http.Request) {
				resp.WriteHeader(http.StatusBadRequest)
				resp.Write([]byte(`{"message":"bad request","httpStatusCode":400,"errorCode":0}`))
			})),
			expectedErr: fmt.Errorf("Couldn't find nodes with the given filter"),
		},
		"error: unable to unmarshal": {
			version: "4.0",
			server: httptest.NewServer(http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
				if req.Method == http.MethodGet && strings.Contains(req.URL.Path, "/Api/V1/ManagedDevice") {
					resp.WriteHeader(http.StatusOK)
					return
				}
				http.NotFound(resp, req)
			})),
			expectedErr: fmt.Errorf("Error While Parsing Response Data For Node: unexpected end of JSON input"),
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			gc := &GatewayClient{
				http:     &http.Client{},
				host:     tc.server.URL,
				username: "test_username",
				password: "test_password",
				version:  tc.version,
			}
			defer tc.server.Close()

			_, err := gc.GetAllNodes()
			if err != nil {
				if tc.expectedErr.Error() != err.Error() {
					t.Fatal(err)
				}
			}
		})
	}
}

func TestGetNodeByFilters(t *testing.T) {
	responseJSON := `[{ "refId": "softwareOnlyServer-1.1.1.1", "refType": null, "ipAddress": "1.1.1.1", "currentIpAddress": "1.1.1.1", "serviceTag": "VMware-42 05 a8 96 26 f7 98 2c-a6 72 b9 1a 26 94 a9 9c-SW", "model": "VMware Virtual Platform", "deviceType": "SoftwareOnlyServer", "discoverDeviceType": "SOFTWAREONLYSERVER_SLES", "displayName": "pfmc-k8s-20230809-1", "managedState": "MANAGED", "state": "READY", "inUse": false, "serviceReferences": [], "statusMessage": null, "firmwareName": "Default Catalog - PowerFlex 4.5.2.0", "customFirmware": false, "needsAttention": false, "manufacturer": "VMware, Inc.", "systemId": null, "health": "NA", "healthMessage": null, "operatingSystem": "N/A", "numberOfCPUs": 0, "cpuType": null, "nics": 0, "memoryInGB": 0, "infraTemplateDate": null, "infraTemplateId": null, "serverTemplateDate": null, "serverTemplateId": null, "inventoryDate": null, "complianceCheckDate": "2024-05-08T11:16:52.951+00:00", "discoveredDate": "2024-05-08T11:16:51.805+00:00", "deviceGroupList": { "paging": null, "deviceGroup": [ { "link": null, "groupSeqId": -1, "groupName": "Global", "groupDescription": null, "createdDate": null, "createdBy": "admin", "updatedDate": null, "updatedBy": null, "managedDeviceList": null, "groupUserList": null } ] }, "detailLink": { "title": "softwareOnlyServer-1.1.1.1", "href": "/AsmManager/ManagedDevice/softwareOnlyServer-1.1.1.1", "rel": "describedby", "type": null }, "credId": "3f5869e6-6525-4dee-bb0c-fab3fe60771d", "compliance": "NONCOMPLIANT", "failuresCount": 0, "chassisId": null, "parsedFacts": null, "config": null, "hostname": "pfmc-k8s-20230809-1", "osIpAddress": null, "osAdminCredential": null, "osImageType": null, "lastJobs": null, "puppetCertName": "sles-1.1.1.1", "svmAdminCredential": null, "svmName": null, "svmIpAddress": null, "svmImageType": null, "flexosMaintMode": 0, "esxiMaintMode": 0, "vmList": [] }]`

	type testCase struct {
		version     string
		server      *httptest.Server
		expectedErr error
	}

	cases := map[string]testCase{
		"succeed: 3.7 version": {
			version: "3.7",
			server: httptest.NewServer(http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
				if req.Method == http.MethodGet && strings.Contains(req.URL.Path, "/Api/V1/ManagedDevice") {
					resp.WriteHeader(http.StatusOK)
					_, err := resp.Write([]byte(responseJSON))
					if err != nil {
						t.Fatalf("Error writing response: %v", err)
					}
					return
				}
				http.NotFound(resp, req)
			})),
			expectedErr: nil,
		},
		"succeed: 4.0 version": {
			version: "4.0",
			server: httptest.NewServer(http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
				if req.Method == http.MethodGet && strings.Contains(req.URL.Path, "/Api/V1/ManagedDevice") {
					resp.WriteHeader(http.StatusOK)
					_, err := resp.Write([]byte(responseJSON))
					if err != nil {
						t.Fatalf("Error writing response: %v", err)
					}
					return
				}
				http.NotFound(resp, req)
			})),
			expectedErr: nil,
		},
		"error: bad request": {
			server: httptest.NewServer(http.HandlerFunc(func(resp http.ResponseWriter, _ *http.Request) {
				resp.WriteHeader(http.StatusBadRequest)
				resp.Write([]byte(`{"message":"bad request","httpStatusCode":400,"errorCode":0}`))
			})),
			expectedErr: fmt.Errorf("Couldn't find nodes with the given filter"),
		},
		"error: unable to unmarshal": {
			version: "4.0",
			server: httptest.NewServer(http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
				if req.Method == http.MethodGet && strings.Contains(req.URL.Path, "/Api/V1/ManagedDevice") {
					resp.WriteHeader(http.StatusOK)
					return
				}
				http.NotFound(resp, req)
			})),
			expectedErr: fmt.Errorf("Error While Parsing Response Data For Node: unexpected end of JSON input"),
		},
		"error: no nodes": {
			version: "4.0",
			server: httptest.NewServer(http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
				if req.Method == http.MethodGet && strings.Contains(req.URL.Path, "/Api/V1/ManagedDevice") {
					resp.WriteHeader(http.StatusOK)
					content, err := json.Marshal([]types.NodeDetails{})
					if err != nil {
						t.Fatalf("Error writing response: %v", err)
					}
					resp.Write(content)
					return
				}
				http.NotFound(resp, req)
			})),
			expectedErr: fmt.Errorf("Couldn't find nodes with the given filter"),
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			gc := &GatewayClient{
				http:     &http.Client{},
				host:     tc.server.URL,
				username: "test_username",
				password: "test_password",
				version:  tc.version,
			}
			defer tc.server.Close()

			key := "ipAddress"
			value := "1.1.1.1"
			_, err := gc.GetNodeByFilters(key, value)
			if err != nil {
				if tc.expectedErr.Error() != err.Error() {
					t.Fatal(err)
				}
			}
		})
	}
}

func TestGetNodePoolByID(t *testing.T) {
	responseJSON := `{ "link": null, "groupSeqId": 123, "groupName": "Test", "groupDescription": "", "createdDate": "2024-05-08T11:27:46.144+00:00", "createdBy": "admin", "updatedDate": "2024-05-08T11:27:46.144+00:00", "updatedBy": "admin", "managedDeviceList": { "paging": null, "totalCount": 1, "managedDevices": [ { "refId": "softwareOnlyServer-1.1.1.1", "refType": null, "ipAddress": "1.1.1.1", "currentIpAddress": "1.1.1.1", "serviceTag": "VMware-42 05 a8 96 26 f7 98 2c-a6 72 b9 1a 26 94 a9 9c-SW", "model": "VMware Virtual Platform", "deviceType": "SoftwareOnlyServer", "discoverDeviceType": "SOFTWAREONLYSERVER_SLES", "displayName": "pfmc-k8s-20230809-1", "managedState": "MANAGED", "state": "READY", "inUse": false, "serviceReferences": [], "statusMessage": null, "firmwareName": "Default Catalog - PowerFlex 4.5.2.0", "customFirmware": false, "needsAttention": false, "manufacturer": "VMware, Inc.", "systemId": null, "health": "NA", "healthMessage": null, "operatingSystem": "N/A", "numberOfCPUs": 0, "cpuType": null, "nics": 0, "memoryInGB": 0, "infraTemplateDate": null, "infraTemplateId": null, "serverTemplateDate": null, "serverTemplateId": null, "inventoryDate": null, "complianceCheckDate": "2024-05-08T11:16:52.951+00:00", "discoveredDate": "2024-05-08T11:16:51.805+00:00", "deviceGroupList": null, "detailLink": null, "credId": "3f5869e6-6525-4dee-bb0c-fab3fe60771d", "compliance": "NONCOMPLIANT", "failuresCount": 0, "chassisId": null, "parsedFacts": null, "config": null, "hostname": "pfmc-k8s-20230809-1", "osIpAddress": null, "osAdminCredential": null, "osImageType": null, "lastJobs": null, "puppetCertName": "sles-1.1.1.1", "svmAdminCredential": null, "svmName": null, "svmIpAddress": null, "svmImageType": null, "flexosMaintMode": 0, "esxiMaintMode": 0, "vmList": [] } ] }, "groupUserList": { "totalRecords": 1, "groupUsers": [ { "userSeqId": "03569bce-5d9b-47a1-addf-2ec44f91f1b9", "userName": "admin", "firstName": "admin", "lastName": "admin", "role": "SuperUser", "enabled": true } ] } }`

	type testCase struct {
		version     string
		server      *httptest.Server
		expectedErr error
	}

	cases := map[string]testCase{
		"succeed: 3.7 version": {
			version: "3.7",
			server: httptest.NewServer(http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
				if req.Method == http.MethodGet && strings.Contains(req.URL.Path, "/Api/V1/nodepool/") {
					resp.WriteHeader(http.StatusOK)
					_, err := resp.Write([]byte(responseJSON))
					if err != nil {
						t.Fatalf("Error writing response: %v", err)
					}
					return
				}
				http.NotFound(resp, req)
			})),
			expectedErr: nil,
		},
		"succeed: 4.0 version": {
			version: "4.0",
			server: httptest.NewServer(http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
				if req.Method == http.MethodGet && strings.Contains(req.URL.Path, "/Api/V1/nodepool/") {
					resp.WriteHeader(http.StatusOK)
					_, err := resp.Write([]byte(responseJSON))
					if err != nil {
						t.Fatalf("Error writing response: %v", err)
					}
					return
				}
				http.NotFound(resp, req)
			})),
			expectedErr: nil,
		},
		"error: bad request": {
			server: httptest.NewServer(http.HandlerFunc(func(resp http.ResponseWriter, _ *http.Request) {
				resp.WriteHeader(http.StatusBadRequest)
				resp.Write([]byte(`{"message":"bad request","httpStatusCode":400,"errorCode":0}`))
			})),
			expectedErr: fmt.Errorf("Couldn't find nodes with the given filter"),
		},
		"error: unable to unmarshal": {
			version: "4.0",
			server: httptest.NewServer(http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
				if req.Method == http.MethodGet && strings.Contains(req.URL.Path, "/Api/V1/nodepool/") {
					resp.WriteHeader(http.StatusOK)
					resp.Write([]byte(`{abc}`))
					return
				}
				http.NotFound(resp, req)
			})),
			expectedErr: fmt.Errorf("Error While Parsing Response Data For Nodepool: invalid character 'a' looking for beginning of object key string"),
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			gc := &GatewayClient{
				http:     &http.Client{},
				host:     tc.server.URL,
				username: "test_username",
				password: "test_password",
				version:  tc.version,
			}
			defer tc.server.Close()

			id := 123
			_, err := gc.GetNodePoolByID(id)
			if err != nil {
				if tc.expectedErr.Error() != err.Error() {
					t.Fatal(err)
				}
			}
		})
	}
}

func TestGetNodePoolByName(t *testing.T) {
	responseJSON := `{"deviceGroup": [ { "link": null, "groupSeqId": 43, "groupName": "Test", "groupDescription": "", "createdDate": "2024-05-08T11:27:46.144+00:00", "createdBy": "admin", "updatedDate": "2024-05-08T11:27:46.144+00:00", "updatedBy": "admin", "managedDeviceList": null, "groupUserList": null } ] }`

	responseJSONID := `{ "link": null, "groupSeqId": 123, "groupName": "Test", "groupDescription": "", "createdDate": "2024-05-08T11:27:46.144+00:00", "createdBy": "admin", "updatedDate": "2024-05-08T11:27:46.144+00:00", "updatedBy": "admin", "managedDeviceList": { "paging": null, "totalCount": 1, "managedDevices": [ { "refId": "softwareOnlyServer-1.1.1.1", "refType": null, "ipAddress": "1.1.1.1", "currentIpAddress": "1.1.1.1", "serviceTag": "VMware-42 05 a8 96 26 f7 98 2c-a6 72 b9 1a 26 94 a9 9c-SW", "model": "VMware Virtual Platform", "deviceType": "SoftwareOnlyServer", "discoverDeviceType": "SOFTWAREONLYSERVER_SLES", "displayName": "pfmc-k8s-20230809-1", "managedState": "MANAGED", "state": "READY", "inUse": false, "serviceReferences": [], "statusMessage": null, "firmwareName": "Default Catalog - PowerFlex 4.5.2.0", "customFirmware": false, "needsAttention": false, "manufacturer": "VMware, Inc.", "systemId": null, "health": "NA", "healthMessage": null, "operatingSystem": "N/A", "numberOfCPUs": 0, "cpuType": null, "nics": 0, "memoryInGB": 0, "infraTemplateDate": null, "infraTemplateId": null, "serverTemplateDate": null, "serverTemplateId": null, "inventoryDate": null, "complianceCheckDate": "2024-05-08T11:16:52.951+00:00", "discoveredDate": "2024-05-08T11:16:51.805+00:00", "deviceGroupList": null, "detailLink": null, "credId": "3f5869e6-6525-4dee-bb0c-fab3fe60771d", "compliance": "NONCOMPLIANT", "failuresCount": 0, "chassisId": null, "parsedFacts": null, "config": null, "hostname": "pfmc-k8s-20230809-1", "osIpAddress": null, "osAdminCredential": null, "osImageType": null, "lastJobs": null, "puppetCertName": "sles-1.1.1.1", "svmAdminCredential": null, "svmName": null, "svmIpAddress": null, "svmImageType": null, "flexosMaintMode": 0, "esxiMaintMode": 0, "vmList": [] } ] }, "groupUserList": { "totalRecords": 1, "groupUsers": [ { "userSeqId": "03569bce-5d9b-47a1-addf-2ec44f91f1b9", "userName": "admin", "firstName": "admin", "lastName": "admin", "role": "SuperUser", "enabled": true } ] } }`

	type testCase struct {
		poolName    string
		server      *httptest.Server
		expectedErr error
	}

	cases := map[string]testCase{
		"succeed": {
			poolName: "Test",
			server: httptest.NewServer(http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
				if req.Method == http.MethodGet && strings.Contains(req.URL.Path, "/Api/V1/nodepool/") {
					resp.WriteHeader(http.StatusOK)
					_, err := resp.Write([]byte(responseJSONID))
					if err != nil {
						t.Fatalf("Unexpected error: %v", err)
					}
					return
				} else if req.Method == http.MethodGet && strings.Contains(req.URL.Path, "/Api/V1/nodepool") {
					resp.WriteHeader(http.StatusOK)
					_, err := resp.Write([]byte(responseJSON))
					if err != nil {
						t.Fatalf("Unexpected error: %v", err)
					}
					return
				}

				http.NotFound(resp, req)
			})),
			expectedErr: nil,
		},
		"error: getting all node pools": {
			poolName: "Test",
			server: httptest.NewServer(http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
				if req.Method == http.MethodGet && strings.Contains(req.URL.Path, "/Api/V1/nodepool/") {
					resp.WriteHeader(http.StatusOK)
					_, err := resp.Write([]byte(responseJSONID))
					if err != nil {
						t.Fatalf("Unexpected error: %v", err)
					}
					return
				} else if req.Method == http.MethodGet && strings.Contains(req.URL.Path, "/Api/V1/nodepool") {
					resp.WriteHeader(http.StatusBadRequest)
					return
				}

				http.NotFound(resp, req)
			})),
			expectedErr: fmt.Errorf("Couldn't find nodes with the given filter"),
		},
		"error: cannot find pool with name": {
			poolName: "InvalidName",
			server: httptest.NewServer(http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
				if req.Method == http.MethodGet && strings.Contains(req.URL.Path, "/Api/V1/nodepool/") {
					resp.WriteHeader(http.StatusOK)
					_, err := resp.Write([]byte(responseJSONID))
					if err != nil {
						t.Fatalf("Unexpected error: %v", err)
					}
					return
				} else if req.Method == http.MethodGet && strings.Contains(req.URL.Path, "/Api/V1/nodepool") {
					resp.WriteHeader(http.StatusOK)
					_, err := resp.Write([]byte(responseJSON))
					if err != nil {
						t.Fatalf("Unexpected error: %v", err)
					}
					return
				}

				http.NotFound(resp, req)
			})),
			expectedErr: fmt.Errorf("no node pool found with name %s", "InvalidName"),
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			gc := &GatewayClient{
				http:     &http.Client{},
				host:     tc.server.URL,
				username: "test_username",
				password: "test_password",
			}
			defer tc.server.Close()

			_, err := gc.GetNodePoolByName(tc.poolName)
			if err != nil {
				if tc.expectedErr.Error() != err.Error() {
					t.Fatal(err)
				}
			}
		})
	}
}

func TestGetAllNodePools(t *testing.T) {
	responseJSON := `{"deviceGroup": [ { "link": null, "groupSeqId": 43, "groupName": "Test", "groupDescription": "", "createdDate": "2024-05-08T11:27:46.144+00:00", "createdBy": "admin", "updatedDate": "2024-05-08T11:27:46.144+00:00", "updatedBy": "admin", "managedDeviceList": null, "groupUserList": null } ] }`

	type testCase struct {
		version     string
		server      *httptest.Server
		expectedErr error
	}

	cases := map[string]testCase{
		"succeed: 3.7 version": {
			version: "3.7",
			server: httptest.NewServer(http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
				if req.Method == http.MethodGet && strings.Contains(req.URL.Path, "/Api/V1/nodepool") {
					resp.WriteHeader(http.StatusOK)
					_, err := resp.Write([]byte(responseJSON))
					if err != nil {
						t.Fatalf("Error writing response: %v", err)
					}
					return
				}
				http.NotFound(resp, req)
			})),
			expectedErr: nil,
		},
		"succeed: 4.0 version": {
			version: "4.0",
			server: httptest.NewServer(http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
				if req.Method == http.MethodGet && strings.Contains(req.URL.Path, "/Api/V1/nodepool") {
					resp.WriteHeader(http.StatusOK)
					_, err := resp.Write([]byte(responseJSON))
					if err != nil {
						t.Fatalf("Error writing response: %v", err)
					}
					return
				}
				http.NotFound(resp, req)
			})),
			expectedErr: nil,
		},
		"error: bad request": {
			server: httptest.NewServer(http.HandlerFunc(func(resp http.ResponseWriter, _ *http.Request) {
				resp.WriteHeader(http.StatusBadRequest)
				resp.Write([]byte(`{"message":"bad request","httpStatusCode":400,"errorCode":0}`))
			})),
			expectedErr: fmt.Errorf("Couldn't find nodes with the given filter"),
		},
		"error: unable to unmarshal": {
			version: "4.0",
			server: httptest.NewServer(http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
				if req.Method == http.MethodGet && strings.Contains(req.URL.Path, "/Api/V1/nodepool") {
					resp.WriteHeader(http.StatusOK)
					resp.Write([]byte(`{abc}`))
					return
				}
				http.NotFound(resp, req)
			})),
			expectedErr: fmt.Errorf("Error While Parsing Response Data For Nodepool: invalid character 'a' looking for beginning of object key string"),
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			gc := &GatewayClient{
				http:     &http.Client{},
				host:     tc.server.URL,
				username: "test_username",
				password: "test_password",
				version:  tc.version,
			}
			defer tc.server.Close()

			_, err := gc.GetAllNodePools()
			if err != nil {
				if tc.expectedErr.Error() != err.Error() {
					t.Fatal(err)
				}
			}
		})
	}
}
