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
	"net/http"
	"net/http/httptest"
	"strings"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetNodeByID(t *testing.T) {
	responseJSON := `{ "refId": "softwareOnlyServer-1.1.1.1", "refType": null, "ipAddress": "1.1.1.1", "currentIpAddress": "1.1.1.1", "serviceTag": "VMware-42 05 a8 96 26 f7 98 2c-a6 72 b9 1a 26 94 a9 9c-SW", "model": "VMware Virtual Platform", "deviceType": "SoftwareOnlyServer", "discoverDeviceType": "SOFTWAREONLYSERVER_SLES", "displayName": "pfmc-k8s-20230809-1", "managedState": "MANAGED", "state": "READY", "inUse": false, "serviceReferences": [], "statusMessage": null, "firmwareName": "Default Catalog - PowerFlex 4.5.2.0", "customFirmware": false, "needsAttention": false, "manufacturer": "VMware, Inc.", "systemId": null, "health": "NA", "healthMessage": null, "operatingSystem": "N/A", "numberOfCPUs": 0, "cpuType": null, "nics": 0, "memoryInGB": 0, "infraTemplateDate": null, "infraTemplateId": null, "serverTemplateDate": null, "serverTemplateId": null, "inventoryDate": null, "complianceCheckDate": "2024-05-08T11:16:52.951+00:00", "discoveredDate": "2024-05-08T11:16:51.805+00:00", "deviceGroupList": { "paging": null, "deviceGroup": [ { "link": null, "groupSeqId": -1, "groupName": "Global", "groupDescription": null, "createdDate": null, "createdBy": "admin", "updatedDate": null, "updatedBy": null, "managedDeviceList": null, "groupUserList": null } ] }, "detailLink": { "title": "softwareOnlyServer-1.1.1.1", "href": "/AsmManager/ManagedDevice/softwareOnlyServer-1.1.1.1", "rel": "describedby", "type": null }, "credId": "3f5869e6-6525-4dee-bb0c-fab3fe60771d", "compliance": "NONCOMPLIANT", "failuresCount": 0, "chassisId": null, "parsedFacts": null, "config": null, "hostname": "pfmc-k8s-20230809-1", "osIpAddress": null, "osAdminCredential": null, "osImageType": null, "lastJobs": null, "puppetCertName": "sles-1.1.1.1", "svmAdminCredential": null, "svmName": null, "svmIpAddress": null, "svmImageType": null, "flexosMaintMode": 0, "esxiMaintMode": 0, "vmList": [] }`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" && strings.Contains(r.URL.Path, "/Api/V1/ManagedDevice/") {
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

	id := "softwareOnlyServer-1.1.1.1"
	nodeDetails, err := gc.GetNodeByID(id)
	defer server.Close()

	gc := &GatewayClient{
		http:     &http.Client{},
		host:     server.URL,
		username: "test_username",
		password: "test_password",
	}

	id := "softwareOnlyServer-1.1.1.1"
	nodeDetails, err := gc.GetNodeByID(id)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
		t.Fatalf("Unexpected error: %v", err)
	}

	assert.NotNil(t, nodeDetails, "Expected non-nil response")
	assert.EqualValues(t, nodeDetails.RefID, "softwareOnlyServer-1.1.1.1")
}

func TestGetAllNodes(t *testing.T) {
	responseJSON := `[{ "refId": "softwareOnlyServer-1.1.1.1", "refType": null, "ipAddress": "1.1.1.1", "currentIpAddress": "1.1.1.1", "serviceTag": "VMware-42 05 a8 96 26 f7 98 2c-a6 72 b9 1a 26 94 a9 9c-SW", "model": "VMware Virtual Platform", "deviceType": "SoftwareOnlyServer", "discoverDeviceType": "SOFTWAREONLYSERVER_SLES", "displayName": "pfmc-k8s-20230809-1", "managedState": "MANAGED", "state": "READY", "inUse": false, "serviceReferences": [], "statusMessage": null, "firmwareName": "Default Catalog - PowerFlex 4.5.2.0", "customFirmware": false, "needsAttention": false, "manufacturer": "VMware, Inc.", "systemId": null, "health": "NA", "healthMessage": null, "operatingSystem": "N/A", "numberOfCPUs": 0, "cpuType": null, "nics": 0, "memoryInGB": 0, "infraTemplateDate": null, "infraTemplateId": null, "serverTemplateDate": null, "serverTemplateId": null, "inventoryDate": null, "complianceCheckDate": "2024-05-08T11:16:52.951+00:00", "discoveredDate": "2024-05-08T11:16:51.805+00:00", "deviceGroupList": { "paging": null, "deviceGroup": [ { "link": null, "groupSeqId": -1, "groupName": "Global", "groupDescription": null, "createdDate": null, "createdBy": "admin", "updatedDate": null, "updatedBy": null, "managedDeviceList": null, "groupUserList": null } ] }, "detailLink": { "title": "softwareOnlyServer-1.1.1.1", "href": "/AsmManager/ManagedDevice/softwareOnlyServer-1.1.1.1", "rel": "describedby", "type": null }, "credId": "3f5869e6-6525-4dee-bb0c-fab3fe60771d", "compliance": "NONCOMPLIANT", "failuresCount": 0, "chassisId": null, "parsedFacts": null, "config": null, "hostname": "pfmc-k8s-20230809-1", "osIpAddress": null, "osAdminCredential": null, "osImageType": null, "lastJobs": null, "puppetCertName": "sles-1.1.1.1", "svmAdminCredential": null, "svmName": null, "svmIpAddress": null, "svmImageType": null, "flexosMaintMode": 0, "esxiMaintMode": 0, "vmList": [] }]`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" && strings.Contains(r.URL.Path, "/Api/V1/ManagedDevice") {
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

	nodes, err := gc.GetAllNodes()
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	assert.NotNil(t, nodes, "Expected non-nil response")
	assert.EqualValues(t, nodes[0].RefID, "softwareOnlyServer-1.1.1.1")
}

func TestGetNodeByFilters(t *testing.T) {
	responseJSON := `[{ "refId": "softwareOnlyServer-1.1.1.1", "refType": null, "ipAddress": "1.1.1.1", "currentIpAddress": "1.1.1.1", "serviceTag": "VMware-42 05 a8 96 26 f7 98 2c-a6 72 b9 1a 26 94 a9 9c-SW", "model": "VMware Virtual Platform", "deviceType": "SoftwareOnlyServer", "discoverDeviceType": "SOFTWAREONLYSERVER_SLES", "displayName": "pfmc-k8s-20230809-1", "managedState": "MANAGED", "state": "READY", "inUse": false, "serviceReferences": [], "statusMessage": null, "firmwareName": "Default Catalog - PowerFlex 4.5.2.0", "customFirmware": false, "needsAttention": false, "manufacturer": "VMware, Inc.", "systemId": null, "health": "NA", "healthMessage": null, "operatingSystem": "N/A", "numberOfCPUs": 0, "cpuType": null, "nics": 0, "memoryInGB": 0, "infraTemplateDate": null, "infraTemplateId": null, "serverTemplateDate": null, "serverTemplateId": null, "inventoryDate": null, "complianceCheckDate": "2024-05-08T11:16:52.951+00:00", "discoveredDate": "2024-05-08T11:16:51.805+00:00", "deviceGroupList": { "paging": null, "deviceGroup": [ { "link": null, "groupSeqId": -1, "groupName": "Global", "groupDescription": null, "createdDate": null, "createdBy": "admin", "updatedDate": null, "updatedBy": null, "managedDeviceList": null, "groupUserList": null } ] }, "detailLink": { "title": "softwareOnlyServer-1.1.1.1", "href": "/AsmManager/ManagedDevice/softwareOnlyServer-1.1.1.1", "rel": "describedby", "type": null }, "credId": "3f5869e6-6525-4dee-bb0c-fab3fe60771d", "compliance": "NONCOMPLIANT", "failuresCount": 0, "chassisId": null, "parsedFacts": null, "config": null, "hostname": "pfmc-k8s-20230809-1", "osIpAddress": null, "osAdminCredential": null, "osImageType": null, "lastJobs": null, "puppetCertName": "sles-1.1.1.1", "svmAdminCredential": null, "svmName": null, "svmIpAddress": null, "svmImageType": null, "flexosMaintMode": 0, "esxiMaintMode": 0, "vmList": [] }]`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" && strings.Contains(r.URL.Path, "/Api/V1/ManagedDevice") {
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

	key := "ipAddress"
	value := "1.1.1.1"
	nodes, err := gc.GetNodeByFilters(key, value)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	assert.NotNil(t, nodes, "Expected non-nil response")
	assert.EqualValues(t, nodes[0].RefID, "softwareOnlyServer-1.1.1.1")
}

func TestGetNodePoolByID(t *testing.T) {
	responseJSON := `{ "link": null, "groupSeqId": 123, "groupName": "Test", "groupDescription": "", "createdDate": "2024-05-08T11:27:46.144+00:00", "createdBy": "admin", "updatedDate": "2024-05-08T11:27:46.144+00:00", "updatedBy": "admin", "managedDeviceList": { "paging": null, "totalCount": 1, "managedDevices": [ { "refId": "softwareOnlyServer-1.1.1.1", "refType": null, "ipAddress": "1.1.1.1", "currentIpAddress": "1.1.1.1", "serviceTag": "VMware-42 05 a8 96 26 f7 98 2c-a6 72 b9 1a 26 94 a9 9c-SW", "model": "VMware Virtual Platform", "deviceType": "SoftwareOnlyServer", "discoverDeviceType": "SOFTWAREONLYSERVER_SLES", "displayName": "pfmc-k8s-20230809-1", "managedState": "MANAGED", "state": "READY", "inUse": false, "serviceReferences": [], "statusMessage": null, "firmwareName": "Default Catalog - PowerFlex 4.5.2.0", "customFirmware": false, "needsAttention": false, "manufacturer": "VMware, Inc.", "systemId": null, "health": "NA", "healthMessage": null, "operatingSystem": "N/A", "numberOfCPUs": 0, "cpuType": null, "nics": 0, "memoryInGB": 0, "infraTemplateDate": null, "infraTemplateId": null, "serverTemplateDate": null, "serverTemplateId": null, "inventoryDate": null, "complianceCheckDate": "2024-05-08T11:16:52.951+00:00", "discoveredDate": "2024-05-08T11:16:51.805+00:00", "deviceGroupList": null, "detailLink": null, "credId": "3f5869e6-6525-4dee-bb0c-fab3fe60771d", "compliance": "NONCOMPLIANT", "failuresCount": 0, "chassisId": null, "parsedFacts": null, "config": null, "hostname": "pfmc-k8s-20230809-1", "osIpAddress": null, "osAdminCredential": null, "osImageType": null, "lastJobs": null, "puppetCertName": "sles-1.1.1.1", "svmAdminCredential": null, "svmName": null, "svmIpAddress": null, "svmImageType": null, "flexosMaintMode": 0, "esxiMaintMode": 0, "vmList": [] } ] }, "groupUserList": { "totalRecords": 1, "groupUsers": [ { "userSeqId": "03569bce-5d9b-47a1-addf-2ec44f91f1b9", "userName": "admin", "firstName": "admin", "lastName": "admin", "role": "SuperUser", "enabled": true } ] } }`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" && strings.Contains(r.URL.Path, "/Api/V1/nodepool/") {
			w.WriteHeader(http.StatusOK)
			_, err := w.Write([]byte(responseJSON))
	responseJSON := `{ "link": null, "groupSeqId": 123, "groupName": "Test", "groupDescription": "", "createdDate": "2024-05-08T11:27:46.144+00:00", "createdBy": "admin", "updatedDate": "2024-05-08T11:27:46.144+00:00", "updatedBy": "admin", "managedDeviceList": { "paging": null, "totalCount": 1, "managedDevices": [ { "refId": "softwareOnlyServer-1.1.1.1", "refType": null, "ipAddress": "1.1.1.1", "currentIpAddress": "1.1.1.1", "serviceTag": "VMware-42 05 a8 96 26 f7 98 2c-a6 72 b9 1a 26 94 a9 9c-SW", "model": "VMware Virtual Platform", "deviceType": "SoftwareOnlyServer", "discoverDeviceType": "SOFTWAREONLYSERVER_SLES", "displayName": "pfmc-k8s-20230809-1", "managedState": "MANAGED", "state": "READY", "inUse": false, "serviceReferences": [], "statusMessage": null, "firmwareName": "Default Catalog - PowerFlex 4.5.2.0", "customFirmware": false, "needsAttention": false, "manufacturer": "VMware, Inc.", "systemId": null, "health": "NA", "healthMessage": null, "operatingSystem": "N/A", "numberOfCPUs": 0, "cpuType": null, "nics": 0, "memoryInGB": 0, "infraTemplateDate": null, "infraTemplateId": null, "serverTemplateDate": null, "serverTemplateId": null, "inventoryDate": null, "complianceCheckDate": "2024-05-08T11:16:52.951+00:00", "discoveredDate": "2024-05-08T11:16:51.805+00:00", "deviceGroupList": null, "detailLink": null, "credId": "3f5869e6-6525-4dee-bb0c-fab3fe60771d", "compliance": "NONCOMPLIANT", "failuresCount": 0, "chassisId": null, "parsedFacts": null, "config": null, "hostname": "pfmc-k8s-20230809-1", "osIpAddress": null, "osAdminCredential": null, "osImageType": null, "lastJobs": null, "puppetCertName": "sles-1.1.1.1", "svmAdminCredential": null, "svmName": null, "svmIpAddress": null, "svmImageType": null, "flexosMaintMode": 0, "esxiMaintMode": 0, "vmList": [] } ] }, "groupUserList": { "totalRecords": 1, "groupUsers": [ { "userSeqId": "03569bce-5d9b-47a1-addf-2ec44f91f1b9", "userName": "admin", "firstName": "admin", "lastName": "admin", "role": "SuperUser", "enabled": true } ] } }`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" && strings.Contains(r.URL.Path, "/Api/V1/nodepool/") {
			w.WriteHeader(http.StatusOK)
			_, err := w.Write([]byte(responseJSON))
			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
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

	id := 123
	nodePoolDetails, err := gc.GetNodePoolByID(id)
				t.Fatalf("Unexpected error: %v", err)
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

	id := 123
	nodePoolDetails, err := gc.GetNodePoolByID(id)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
		t.Fatalf("Unexpected error: %v", err)
	}

	assert.NotNil(t, nodePoolDetails, "Expected non-nil response")
	assert.EqualValues(t, nodePoolDetails.ManagedDeviceList.ManagedDevices[0].RefID, "softwareOnlyServer-1.1.1.1")
}

func TestGetNodePoolByName(t *testing.T) {
	responseJSON := `{"deviceGroup": [ { "link": null, "groupSeqId": 43, "groupName": "Test", "groupDescription": "", "createdDate": "2024-05-08T11:27:46.144+00:00", "createdBy": "admin", "updatedDate": "2024-05-08T11:27:46.144+00:00", "updatedBy": "admin", "managedDeviceList": null, "groupUserList": null } ] }`

	responseJSONID := `{ "link": null, "groupSeqId": 123, "groupName": "Test", "groupDescription": "", "createdDate": "2024-05-08T11:27:46.144+00:00", "createdBy": "admin", "updatedDate": "2024-05-08T11:27:46.144+00:00", "updatedBy": "admin", "managedDeviceList": { "paging": null, "totalCount": 1, "managedDevices": [ { "refId": "softwareOnlyServer-1.1.1.1", "refType": null, "ipAddress": "1.1.1.1", "currentIpAddress": "1.1.1.1", "serviceTag": "VMware-42 05 a8 96 26 f7 98 2c-a6 72 b9 1a 26 94 a9 9c-SW", "model": "VMware Virtual Platform", "deviceType": "SoftwareOnlyServer", "discoverDeviceType": "SOFTWAREONLYSERVER_SLES", "displayName": "pfmc-k8s-20230809-1", "managedState": "MANAGED", "state": "READY", "inUse": false, "serviceReferences": [], "statusMessage": null, "firmwareName": "Default Catalog - PowerFlex 4.5.2.0", "customFirmware": false, "needsAttention": false, "manufacturer": "VMware, Inc.", "systemId": null, "health": "NA", "healthMessage": null, "operatingSystem": "N/A", "numberOfCPUs": 0, "cpuType": null, "nics": 0, "memoryInGB": 0, "infraTemplateDate": null, "infraTemplateId": null, "serverTemplateDate": null, "serverTemplateId": null, "inventoryDate": null, "complianceCheckDate": "2024-05-08T11:16:52.951+00:00", "discoveredDate": "2024-05-08T11:16:51.805+00:00", "deviceGroupList": null, "detailLink": null, "credId": "3f5869e6-6525-4dee-bb0c-fab3fe60771d", "compliance": "NONCOMPLIANT", "failuresCount": 0, "chassisId": null, "parsedFacts": null, "config": null, "hostname": "pfmc-k8s-20230809-1", "osIpAddress": null, "osAdminCredential": null, "osImageType": null, "lastJobs": null, "puppetCertName": "sles-1.1.1.1", "svmAdminCredential": null, "svmName": null, "svmIpAddress": null, "svmImageType": null, "flexosMaintMode": 0, "esxiMaintMode": 0, "vmList": [] } ] }, "groupUserList": { "totalRecords": 1, "groupUsers": [ { "userSeqId": "03569bce-5d9b-47a1-addf-2ec44f91f1b9", "userName": "admin", "firstName": "admin", "lastName": "admin", "role": "SuperUser", "enabled": true } ] } }`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" && strings.Contains(r.URL.Path, "/Api/V1/nodepool/") {
			w.WriteHeader(http.StatusOK)
			_, err := w.Write([]byte(responseJSONID))
			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}
			return
		} else if r.Method == "GET" && strings.Contains(r.URL.Path, "/Api/V1/nodepool") {
			w.WriteHeader(http.StatusOK)
			_, err := w.Write([]byte(responseJSON))
			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}
			return
		}

		http.NotFound(w, r)
	responseJSON := `{"deviceGroup": [ { "link": null, "groupSeqId": 43, "groupName": "Test", "groupDescription": "", "createdDate": "2024-05-08T11:27:46.144+00:00", "createdBy": "admin", "updatedDate": "2024-05-08T11:27:46.144+00:00", "updatedBy": "admin", "managedDeviceList": null, "groupUserList": null } ] }`

	responseJSONID := `{ "link": null, "groupSeqId": 123, "groupName": "Test", "groupDescription": "", "createdDate": "2024-05-08T11:27:46.144+00:00", "createdBy": "admin", "updatedDate": "2024-05-08T11:27:46.144+00:00", "updatedBy": "admin", "managedDeviceList": { "paging": null, "totalCount": 1, "managedDevices": [ { "refId": "softwareOnlyServer-1.1.1.1", "refType": null, "ipAddress": "1.1.1.1", "currentIpAddress": "1.1.1.1", "serviceTag": "VMware-42 05 a8 96 26 f7 98 2c-a6 72 b9 1a 26 94 a9 9c-SW", "model": "VMware Virtual Platform", "deviceType": "SoftwareOnlyServer", "discoverDeviceType": "SOFTWAREONLYSERVER_SLES", "displayName": "pfmc-k8s-20230809-1", "managedState": "MANAGED", "state": "READY", "inUse": false, "serviceReferences": [], "statusMessage": null, "firmwareName": "Default Catalog - PowerFlex 4.5.2.0", "customFirmware": false, "needsAttention": false, "manufacturer": "VMware, Inc.", "systemId": null, "health": "NA", "healthMessage": null, "operatingSystem": "N/A", "numberOfCPUs": 0, "cpuType": null, "nics": 0, "memoryInGB": 0, "infraTemplateDate": null, "infraTemplateId": null, "serverTemplateDate": null, "serverTemplateId": null, "inventoryDate": null, "complianceCheckDate": "2024-05-08T11:16:52.951+00:00", "discoveredDate": "2024-05-08T11:16:51.805+00:00", "deviceGroupList": null, "detailLink": null, "credId": "3f5869e6-6525-4dee-bb0c-fab3fe60771d", "compliance": "NONCOMPLIANT", "failuresCount": 0, "chassisId": null, "parsedFacts": null, "config": null, "hostname": "pfmc-k8s-20230809-1", "osIpAddress": null, "osAdminCredential": null, "osImageType": null, "lastJobs": null, "puppetCertName": "sles-1.1.1.1", "svmAdminCredential": null, "svmName": null, "svmIpAddress": null, "svmImageType": null, "flexosMaintMode": 0, "esxiMaintMode": 0, "vmList": [] } ] }, "groupUserList": { "totalRecords": 1, "groupUsers": [ { "userSeqId": "03569bce-5d9b-47a1-addf-2ec44f91f1b9", "userName": "admin", "firstName": "admin", "lastName": "admin", "role": "SuperUser", "enabled": true } ] } }`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" && strings.Contains(r.URL.Path, "/Api/V1/nodepool/") {
			w.WriteHeader(http.StatusOK)
			_, err := w.Write([]byte(responseJSONID))
			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}
			return
		} else if r.Method == "GET" && strings.Contains(r.URL.Path, "/Api/V1/nodepool") {
			w.WriteHeader(http.StatusOK)
			_, err := w.Write([]byte(responseJSON))
			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
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

	name := "Test"
	nodePoolDetails, err := gc.GetNodePoolByName(name)
	defer server.Close()

	gc := &GatewayClient{
		http:     &http.Client{},
		host:     server.URL,
		username: "test_username",
		password: "test_password",
	}

	name := "Test"
	nodePoolDetails, err := gc.GetNodePoolByName(name)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
		t.Fatalf("Unexpected error: %v", err)
	}

	assert.NotNil(t, nodePoolDetails, "Expected non-nil response")
	assert.GreaterOrEqual(t, len(nodePoolDetails.ManagedDeviceList.ManagedDevices), 1)
}

func TestGetAllNodePools(t *testing.T) {
	responseJSON := `{"deviceGroup": [ { "link": null, "groupSeqId": 43, "groupName": "Test", "groupDescription": "", "createdDate": "2024-05-08T11:27:46.144+00:00", "createdBy": "admin", "updatedDate": "2024-05-08T11:27:46.144+00:00", "updatedBy": "admin", "managedDeviceList": null, "groupUserList": null } ] }`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" && strings.Contains(r.URL.Path, "/Api/V1/nodepool") {
			w.WriteHeader(http.StatusOK)
			_, err := w.Write([]byte(responseJSON))
			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}
			return
		}
		http.NotFound(w, r)
func TestGetAllNodePools(t *testing.T) {
	responseJSON := `{"deviceGroup": [ { "link": null, "groupSeqId": 43, "groupName": "Test", "groupDescription": "", "createdDate": "2024-05-08T11:27:46.144+00:00", "createdBy": "admin", "updatedDate": "2024-05-08T11:27:46.144+00:00", "updatedBy": "admin", "managedDeviceList": null, "groupUserList": null } ] }`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" && strings.Contains(r.URL.Path, "/Api/V1/nodepool") {
			w.WriteHeader(http.StatusOK)
			_, err := w.Write([]byte(responseJSON))
			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
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

	nodePoolDetails, err := gc.GetAllNodePools()
	defer server.Close()

	gc := &GatewayClient{
		http:     &http.Client{},
		host:     server.URL,
		username: "test_username",
		password: "test_password",
	}

	nodePoolDetails, err := gc.GetAllNodePools()
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
		t.Fatalf("Unexpected error: %v", err)
	}

	assert.NotNil(t, nodePoolDetails, "Expected non-nil response")
	assert.GreaterOrEqual(t, len(nodePoolDetails.NodePoolDetails), 1)
}
