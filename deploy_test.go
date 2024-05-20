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
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

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

	defer server.Close()

	gc, err := NewGateway(server.URL, "test_username", "test_password", false, false)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	assert.Nil(t, err, "Unexpected error")
	assert.NotNil(t, gc, "GatewayClient is nil")
	assert.Equal(t, "mock_access_token", gc.token, "Unexpected access token")
	assert.Equal(t, "4.0", gc.version, "Unexpected version")
}

// TestGetVersion tests the GetVersion function.
func TestGetVersion(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet && r.URL.Path == "/api/version" {
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
	assert.NoError(t, err)
	assert.Equal(t, "4.0", version)
}

// TestUploadPackages tests the UploadPackages function.
func TestUploadPackages(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost && r.URL.Path == "/im/types/installationPackages/instances/actions/uploadPackages" {
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
	assert.Error(t, err)

	expectedErrorMsg := "stat mock_file.tar: no such file or directory"
	assert.EqualError(t, err, expectedErrorMsg)
}

func TestParseCSV(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost && r.URL.Path == "/im/types/Configuration/instances/actions/parseFromCSV" {
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
	assert.Error(t, err)

	expectedErrorMsg := "open test_file.csv: no such file or directory"
	assert.EqualError(t, err, expectedErrorMsg)
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
	}

	packageDetails, err := gc.GetPackageDetails()
	assert.NoError(t, err)
	assert.NotNil(t, packageDetails)
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
	assert.NoError(t, err)
	assert.Equal(t, 200, packageResponse.StatusCode)
}

func TestBeginInstallation(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost && strings.HasPrefix(r.URL.Path, "/im/types/Configuration/actions/install") {
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
	assert.Error(t, err)

	expectedErrorMsg := "unexpected end of JSON input"
	assert.EqualError(t, err, expectedErrorMsg)
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
		http:     &http.Client{},
		host:     server.URL,
		username: "test_username",
		password: "test_password",
	}

	gatewayResponse, err := gc.MoveToNextPhase()
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, gatewayResponse.StatusCode)
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
		http:     &http.Client{},
		host:     server.URL,
		username: "test_username",
		password: "test_password",
	}

	gatewayResponse, err := gc.RetryPhase()
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, gatewayResponse.StatusCode)
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
		http:     &http.Client{},
		host:     server.URL,
		username: "test_username",
		password: "test_password",
	}

	gatewayResponse, err := gc.AbortOperation()
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, gatewayResponse.StatusCode)
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
		http:     &http.Client{},
		host:     server.URL,
		username: "test_username",
		password: "test_password",
	}

	gatewayResponse, err := gc.ClearQueueCommand()
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, gatewayResponse.StatusCode)
}

func TestMoveToIdlePhase(t *testing.T) {
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
}

func TestCheckForCompletionQueueCommands(t *testing.T) {
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
	}

	gatewayResponse, err := gc.CheckForCompletionQueueCommands("Query")
	assert.NoError(t, err)
	assert.NotNil(t, gatewayResponse)
	assert.Equal(t, http.StatusOK, gatewayResponse.StatusCode)
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
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, gatewayResponse.StatusCode)
}
