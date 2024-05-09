// Copyright © 2021 - 2023 Dell Inc. or its subsidiaries. All Rights Reserved.
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
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// TestNewGateway tests the NewGateway function.
func TestNewGateway(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" && r.URL.Path == "/rest/auth/login" {

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			fmt.Fprintln(w, `{"access_token":"mock_access_token"}`)
			return
		}
		if r.Method == "GET" && r.URL.Path == "/api/version" {
			w.Header().Set("Content-Type", "text/plain")
			w.WriteHeader(http.StatusOK)
			fmt.Fprintln(w, "4.0")
			return
		}
		http.NotFound(w, r)
	}))

	defer server.Close()

	gc, err := NewGateway(server.URL, "test_username", "test_password", false, false)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if gc == nil {
		t.Fatal("GatewayClient is nil")
	}
	if gc.token != "mock_access_token" {
		t.Errorf("Unexpected access token: %s", gc.token)
	}
	if gc.version != "4.0" {
		t.Errorf("Unexpected version: %s", gc.version)
	}
}

// TestGetVersion tests the GetVersion function.
func TestGetVersion(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" && r.URL.Path == "/api/version" {
			w.Header().Set("Content-Type", "text/plain")
			w.WriteHeader(http.StatusOK)
			fmt.Fprintln(w, "4.0")
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

	version, err := gc.GetVersion()
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if version != "4.0" {
		t.Errorf("Unexpected version: %s", version)
	}
}

// TestUploadPackages tests the UploadPackages function.
func TestUploadPackages(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" && r.URL.Path == "/im/types/installationPackages/instances/actions/uploadPackages" {
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

	_, err := gc.UploadPackages([]string{"mock_file.tar"})
	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	expectedErrorMsg := "stat mock_file.tar: no such file or directory"
	if err.Error() != expectedErrorMsg {
		t.Errorf("Unexpected error message: %s", err.Error())
	}
}

func TestParseCSV(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" && r.URL.Path == "/im/types/Configuration/instances/actions/parseFromCSV" {
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

	_, err := gc.ParseCSV("test_file.csv")
	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	expectedErrorMsg := "open test_file.csv: no such file or directory"
	if err.Error() != expectedErrorMsg {
		t.Errorf("Unexpected error message: %s", err.Error())
	}
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
		if r.Method == "GET" && r.URL.Path == "/im/types/installationPackages/instances" {
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

	packageDetails, err := gc.GetPackageDetails()
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if packageDetails == nil {
		t.Error("Package details are nil")
	}
}

func TestDeletePackage(t *testing.T) {
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
		username: "test_username",
		password: "test_password",
	}

	packageResponse, err := gc.DeletePackage("test_package")
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if packageResponse.StatusCode != 200 {
		t.Errorf("Unexpected status code: %d", packageResponse.StatusCode)
	}
}

func TestBeginInstallation(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" && strings.HasPrefix(r.URL.Path, "/im/types/Configuration/actions/install") {
			w.WriteHeader(http.StatusAccepted)
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

	_, err := gc.BeginInstallation("", "mdm_user", "mdm_password", "lia_password", true, true, true, false)
	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	expectedErrorMsg := "unexpected end of JSON input"
	if err.Error() != expectedErrorMsg {
		t.Errorf("Unexpected error message: %s", err.Error())
	}
}

func TestMoveToNextPhase(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" && r.URL.Path == "/im/types/ProcessPhase/actions/moveToNextPhase" {
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

	gatewayResponse, err := gc.MoveToNextPhase()
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if gatewayResponse.StatusCode != http.StatusOK {
		t.Errorf("Unexpected status code: %d", gatewayResponse.StatusCode)
	}
}

func TestRetryPhase(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" && strings.HasPrefix(r.URL.Path, "/im/types/Command/instances/actions/retry") {
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

	gatewayResponse, err := gc.RetryPhase()
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if gatewayResponse.StatusCode != http.StatusOK {
		t.Errorf("Unexpected status code: %d", gatewayResponse.StatusCode)
	}
}

func TestAbortOperation(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" && r.URL.Path == "/im/types/Command/instances/actions/abort" {
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

	gatewayResponse, err := gc.AbortOperation()
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if gatewayResponse.StatusCode != http.StatusOK {
		t.Errorf("Unexpected status code: %d", gatewayResponse.StatusCode)
	}
}

func TestClearQueueCommand(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" && r.URL.Path == "/im/types/Command/instances/actions/clear" {
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

	gatewayResponse, err := gc.ClearQueueCommand()
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if gatewayResponse.StatusCode != http.StatusOK {
		t.Errorf("Unexpected status code: %d", gatewayResponse.StatusCode)
	}
}

func TestMoveToIdlePhase(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" && r.URL.Path == "/im/types/ProcessPhase/actions/moveToIdlePhase" {
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
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if gatewayResponse.StatusCode != http.StatusOK {
		t.Errorf("Unexpected status code: %d", gatewayResponse.StatusCode)
	}
}

func TestCheckForCompletionQueueCommands(t *testing.T) {
	responseJSON := `{
		"MDM Commands": []
	}`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" && r.URL.Path == "/im/types/Command/instances" {
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

	gatewayResponse, err := gc.CheckForCompletionQueueCommands("Query")
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if gatewayResponse == nil {
		t.Error("Gateway response is nil")
	}

	if gatewayResponse.StatusCode != http.StatusOK {
		t.Errorf("Unexpected status code: %d", gatewayResponse.StatusCode)
	}
}

func TestUninstallCluster(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" && strings.Contains(r.URL.Path, "/im/types/Configuration/actions/uninstall") {
			w.WriteHeader(http.StatusAccepted)
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

	gatewayResponse, err := gc.UninstallCluster(jsonStr, mdmUsername, mdmPassword, liaPassword, allowNonSecureCommunicationWithMdm, allowNonSecureCommunicationWithLia, disableNonMgmtComponentsAuth, false)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if gatewayResponse.StatusCode != http.StatusOK {
		t.Errorf("Unexpected status code: %d", gatewayResponse.StatusCode)
	}
}
