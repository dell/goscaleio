// Copyright © 2024 - 2026 Dell Inc. or its subsidiaries. All Rights Reserved.
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
	"math"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	types "github.com/dell/goscaleio/types/v1"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func Test_GetAllNvmeHosts(t *testing.T) {
	type checkFn func(*testing.T, []types.NvmeHost, error)
	check := func(fns ...checkFn) []checkFn { return fns }

	responseJSON := `[{
    "name": "mock-name",
    "hostOsFullType": "mock-host-os-full-type",
    "systemID": "mock-system-id",
    "sdcApproved": null,
    "sdcAgentActive": null,
    "mdmIpAddressesCurrent": null,
    "sdcIp": null,
    "sdcIps": null,
    "osType": null,
    "perfProfile": null,
    "peerMdmId": null,
    "sdtId": null,
    "mdmConnectionState": null,
    "softwareVersionInfo": null,
    "socketAllocationFailure": null,
    "memoryAllocationFailure": null,
    "sdcGuid": null,
    "installedSoftwareVersionInfo": null,
    "kernelVersion": null,
    "kernelBuildNumber": null,
    "sdcApprovedIps": null,
    "hostType": "NVMeHost",
    "sdrId": null,
    "versionInfo": null,
    "sdcType": null,
    "nqn": "mock-nqn",
    "maxNumPaths": null,
    "maxNumSysPorts": null,
    "id": "mock-id",
    "links": [
        {
            "rel": "self",
            "href": "/api/instances/Host::mock-id"
        },
        {
            "rel": "/api/Host/relationship/Volume",
            "href": "/api/instances/Host::mock-id/relationships/Volume"
        },
        {
            "rel": "/api/Host/relationship/NvmeController",
            "href": "/api/instances/Host::mock-id/relationships/NvmeController"
        },
        {
            "rel": "/api/parent/relationship/systemID",
            "href": "/api/instances/System::mock-system-id"
        }
    ]
}]`

	hasNoError := func(t *testing.T, _ []types.NvmeHost, err error) {
		if err != nil {
			t.Fatalf("expected no error")
		}
	}

	checkLength := func(length int) func(t *testing.T, hosts []types.NvmeHost, err error) {
		return func(t *testing.T, hosts []types.NvmeHost, _ error) {
			assert.Equal(t, length, len(hosts))
		}
	}

	hasError := func(t *testing.T, _ []types.NvmeHost, err error) {
		if err == nil {
			t.Fatalf("expected error")
		}
	}

	tests := map[string]func(t *testing.T) (*httptest.Server, types.System, []checkFn){
		"success": func(t *testing.T) (*httptest.Server, types.System, []checkFn) {
			systemID := "mock-system-id"
			url := fmt.Sprintf("/api/instances/System::%v/relationships/Sdc", systemID)
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method == http.MethodGet && strings.EqualFold(r.URL.Path, url) {
					w.WriteHeader(http.StatusOK)
					_, err := w.Write([]byte(responseJSON))
					if err != nil {
						t.Fatalf("Unexpected error: %v", err)
					}
					return
				}
				http.NotFound(w, r)
			}))
			system := types.System{
				ID: systemID,
			}
			return server, system, check(hasNoError, checkLength(1))
		},
		"error response": func(_ *testing.T) (*httptest.Server, types.System, []checkFn) {
			systemID := "mock-system-id"
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				http.NotFound(w, r)
			}))
			system := types.System{
				ID: systemID,
			}
			return server, system, check(hasError)
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			server, sys, checkFns := tc(t)
			defer server.Close()

			client, _ := NewClientWithArgs(server.URL, "", math.MaxInt64, true, false, "")
			system := System{
				System: &sys,
				client: client,
			}
			nvmeHosts, err := system.GetAllNvmeHosts()

			for _, checkFn := range checkFns {
				checkFn(t, nvmeHosts, err)
			}
		})
	}
}

func Test_GetNvmeHostByID(t *testing.T) {
	type checkFn func(*testing.T, *types.NvmeHost, error)
	check := func(fns ...checkFn) []checkFn { return fns }

	responseJSON := `{
    "name": "mock-name",
    "hostOsFullType": "mock-host-os-full-type",
    "systemID": "mock-system-id",
    "sdcApproved": null,
    "sdcAgentActive": null,
    "mdmIpAddressesCurrent": null,
    "sdcIp": null,
    "sdcIps": null,
    "osType": null,
    "perfProfile": null,
    "peerMdmId": null,
    "sdtId": null,
    "mdmConnectionState": null,
    "softwareVersionInfo": null,
    "socketAllocationFailure": null,
    "memoryAllocationFailure": null,
    "sdcGuid": null,
    "installedSoftwareVersionInfo": null,
    "kernelVersion": null,
    "kernelBuildNumber": null,
    "sdcApprovedIps": null,
    "hostType": "NVMeHost",
    "sdrId": null,
    "versionInfo": null,
    "sdcType": null,
    "nqn": "mock-nqn",
    "maxNumPaths": null,
    "maxNumSysPorts": null,
    "id": "mock-id",
    "links": [
        {
            "rel": "self",
            "href": "/api/instances/Host::mock-id"
        },
        {
            "rel": "/api/Host/relationship/Volume",
            "href": "/api/instances/Host::mock-id/relationships/Volume"
        },
        {
            "rel": "/api/Host/relationship/NvmeController",
            "href": "/api/instances/Host::mock-id/relationships/NvmeController"
        },
        {
            "rel": "/api/parent/relationship/systemID",
            "href": "/api/instances/System::mock-system-id"
        }
    ]
}`

	hasNoError := func(t *testing.T, _ *types.NvmeHost, err error) {
		if err != nil {
			t.Fatalf("expected no error")
		}
	}

	checkName := func(name string) func(t *testing.T, host *types.NvmeHost, err error) {
		return func(t *testing.T, host *types.NvmeHost, _ error) {
			assert.Equal(t, host.Name, name)
		}
	}

	hasError := func(t *testing.T, _ *types.NvmeHost, err error) {
		if err == nil {
			t.Fatalf("expected error")
		}
	}

	tests := map[string]func(t *testing.T) (*httptest.Server, types.System, []checkFn){
		"success": func(t *testing.T) (*httptest.Server, types.System, []checkFn) {
			systemID := "mock-system-id"
			id := "mock-id"
			url := fmt.Sprintf("/api/instances/Sdc::%v", id)
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method == http.MethodGet && strings.EqualFold(r.URL.Path, url) {
					w.WriteHeader(http.StatusOK)
					_, err := w.Write([]byte(responseJSON))
					if err != nil {
						t.Fatalf("Unexpected error: %v", err)
					}
					return
				}
				http.NotFound(w, r)
			}))
			system := types.System{
				ID: systemID,
			}
			return server, system, check(hasNoError, checkName("mock-name"))
		},
		"error response": func(_ *testing.T) (*httptest.Server, types.System, []checkFn) {
			systemID := "mock-system-id"
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				http.Error(w, "Server Error", http.StatusInternalServerError)
			}))
			system := types.System{
				ID: systemID,
			}
			return server, system, check(hasError)
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			server, sys, checkFns := tc(t)
			defer server.Close()

			client, _ := NewClientWithArgs(server.URL, "", math.MaxInt64, true, false, "")
			system := System{
				System: &sys,
				client: client,
			}
			nvmeHost, err := system.GetNvmeHostByID("mock-id")

			for _, checkFn := range checkFns {
				checkFn(t, nvmeHost, err)
			}
		})
	}
}

func Test_CreateNvmeHost(t *testing.T) {
	type checkFn func(*testing.T, *types.NvmeHostResp, error)
	check := func(fns ...checkFn) []checkFn { return fns }

	hasNoError := func(t *testing.T, _ *types.NvmeHostResp, err error) {
		if err != nil {
			t.Fatalf("expected no error")
		}
	}

	checkID := func(id string) func(t *testing.T, resp *types.NvmeHostResp, err error) {
		return func(t *testing.T, resp *types.NvmeHostResp, _ error) {
			assert.Equal(t, resp.ID, id)
		}
	}

	hasError := func(t *testing.T, _ *types.NvmeHostResp, err error) {
		if err == nil {
			t.Fatalf("expected error")
		}
	}

	tests := map[string]func(t *testing.T) (*httptest.Server, types.System, []checkFn){
		"success": func(t *testing.T) (*httptest.Server, types.System, []checkFn) {
			systemID := "mock-system-id"
			url := "/api/types/Host/instances"
			id := "4d2a628100010000"
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method == http.MethodPost && strings.EqualFold(r.URL.Path, url) {
					w.WriteHeader(http.StatusOK)
					_, err := w.Write([]byte(fmt.Sprintf("{\"id\": \"%v\"}", id)))
					if err != nil {
						t.Fatalf("Unexpected error: %v", err)
					}
					return
				}
				http.NotFound(w, r)
			}))
			system := types.System{
				ID: systemID,
			}
			return server, system, check(hasNoError, checkID(id))
		},
		"error response": func(_ *testing.T) (*httptest.Server, types.System, []checkFn) {
			systemID := "mock-system-id"
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				http.Error(w, "Server Error", http.StatusInternalServerError)
			}))
			system := types.System{
				ID: systemID,
			}
			return server, system, check(hasError)
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			server, sys, checkFns := tc(t)
			defer server.Close()

			client, _ := NewClientWithArgs(server.URL, "", math.MaxInt64, true, false, "")
			system := System{
				System: &sys,
				client: client,
			}
			nvmeHostParam := types.NvmeHostParam{
				Name:           "mock-name",
				Nqn:            "mock-nqn",
				MaxNumPaths:    4,
				MaxNumSysPorts: 10,
			}
			resp, err := system.CreateNvmeHost(nvmeHostParam)

			for _, checkFn := range checkFns {
				checkFn(t, resp, err)
			}
		})
	}
}

func Test_ChangeNvmeHostName(t *testing.T) {
	type checkFn func(*testing.T, error)
	check := func(fns ...checkFn) []checkFn { return fns }

	hasNoError := func(t *testing.T, err error) {
		if err != nil {
			t.Fatalf("expected no error")
		}
	}

	hasError := func(t *testing.T, err error) {
		if err == nil {
			t.Fatalf("expected error")
		}
	}

	tests := map[string]func(t *testing.T) (*httptest.Server, types.System, []checkFn){
		"success": func(_ *testing.T) (*httptest.Server, types.System, []checkFn) {
			systemID := "mock-system-id"
			id := "mock-id"
			url := fmt.Sprintf("/api/instances/Sdc::%v/action/setSdcName", id)
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method == http.MethodPost && strings.EqualFold(r.URL.Path, url) {
					w.WriteHeader(http.StatusOK)
					return
				}
				http.NotFound(w, r)
			}))
			system := types.System{
				ID: systemID,
			}
			return server, system, check(hasNoError)
		},
		"error response": func(_ *testing.T) (*httptest.Server, types.System, []checkFn) {
			systemID := "mock-system-id"
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				http.Error(w, "Server Error", http.StatusInternalServerError)
			}))
			system := types.System{
				ID: systemID,
			}
			return server, system, check(hasError)
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			server, sys, checkFns := tc(t)
			defer server.Close()

			client, _ := NewClientWithArgs(server.URL, "", math.MaxInt64, true, false, "")
			system := System{
				System: &sys,
				client: client,
			}
			err := system.ChangeNvmeHostName("mock-id", "new-mock-name")

			for _, checkFn := range checkFns {
				checkFn(t, err)
			}
		})
	}
}

func Test_ChangeNvmeHostMaxNumPaths(t *testing.T) {
	type checkFn func(*testing.T, error)
	check := func(fns ...checkFn) []checkFn { return fns }

	hasNoError := func(t *testing.T, err error) {
		if err != nil {
			t.Fatalf("expected no error")
		}
	}

	hasError := func(t *testing.T, err error) {
		if err == nil {
			t.Fatalf("expected error")
		}
	}

	tests := map[string]func(t *testing.T) (*httptest.Server, types.System, []checkFn){
		"success": func(_ *testing.T) (*httptest.Server, types.System, []checkFn) {
			systemID := "mock-system-id"
			id := "mock-id"
			url := fmt.Sprintf("/api/instances/Host::%v/action/modifyMaxNumPaths", id)
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method == http.MethodPost && strings.EqualFold(r.URL.Path, url) {
					w.WriteHeader(http.StatusOK)
					return
				}
				http.NotFound(w, r)
			}))
			system := types.System{
				ID: systemID,
			}
			return server, system, check(hasNoError)
		},
		"error response": func(_ *testing.T) (*httptest.Server, types.System, []checkFn) {
			systemID := "mock-system-id"
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				http.Error(w, "Server Error", http.StatusInternalServerError)
			}))
			system := types.System{
				ID: systemID,
			}
			return server, system, check(hasError)
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			server, sys, checkFns := tc(t)
			defer server.Close()

			client, _ := NewClientWithArgs(server.URL, "", math.MaxInt64, true, false, "")
			system := System{
				System: &sys,
				client: client,
			}
			err := system.ChangeNvmeHostMaxNumPaths("mock-id", 8)

			for _, checkFn := range checkFns {
				checkFn(t, err)
			}
		})
	}
}

func Test_ChangeNvmeHostMaxNumSysPorts(t *testing.T) {
	type checkFn func(*testing.T, error)
	check := func(fns ...checkFn) []checkFn { return fns }

	hasNoError := func(t *testing.T, err error) {
		if err != nil {
			t.Fatalf("expected no error")
		}
	}

	hasError := func(t *testing.T, err error) {
		if err == nil {
			t.Fatalf("expected error")
		}
	}

	tests := map[string]func(t *testing.T) (*httptest.Server, types.System, []checkFn){
		"success": func(_ *testing.T) (*httptest.Server, types.System, []checkFn) {
			systemID := "mock-system-id"
			id := "mock-id"
			url := fmt.Sprintf("/api/instances/Host::%v/action/modifyMaxNumSysPorts", id)
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method == http.MethodPost && strings.EqualFold(r.URL.Path, url) {
					w.WriteHeader(http.StatusOK)
					return
				}
				http.NotFound(w, r)
			}))
			system := types.System{
				ID: systemID,
			}
			return server, system, check(hasNoError)
		},
		"error response": func(_ *testing.T) (*httptest.Server, types.System, []checkFn) {
			systemID := "mock-system-id"
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				http.Error(w, "Server Error", http.StatusInternalServerError)
			}))
			system := types.System{
				ID: systemID,
			}
			return server, system, check(hasError)
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			server, sys, checkFns := tc(t)
			defer server.Close()

			client, _ := NewClientWithArgs(server.URL, "", math.MaxInt64, true, false, "")
			system := System{
				System: &sys,
				client: client,
			}
			err := system.ChangeNvmeHostMaxNumSysPorts("mock-id", 8)

			for _, checkFn := range checkFns {
				checkFn(t, err)
			}
		})
	}
}

func Test_DeleteNvmeHost(t *testing.T) {
	type checkFn func(*testing.T, error)
	check := func(fns ...checkFn) []checkFn { return fns }

	hasNoError := func(t *testing.T, err error) {
		if err != nil {
			t.Fatalf("expected no error")
		}
	}

	hasError := func(t *testing.T, err error) {
		if err == nil {
			t.Fatalf("expected error")
		}
	}

	tests := map[string]func(t *testing.T) (*httptest.Server, types.System, []checkFn){
		"success": func(_ *testing.T) (*httptest.Server, types.System, []checkFn) {
			systemID := "mock-system-id"
			id := "mock-id"
			url := fmt.Sprintf("/api/instances/Sdc::%v/action/removeSdc", id)
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method == http.MethodPost && strings.EqualFold(r.URL.Path, url) {
					w.WriteHeader(http.StatusOK)
					return
				}
				http.NotFound(w, r)
			}))
			system := types.System{
				ID: systemID,
			}
			return server, system, check(hasNoError)
		},
		"error response": func(_ *testing.T) (*httptest.Server, types.System, []checkFn) {
			systemID := "mock-system-id"
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				http.Error(w, "Server Error", http.StatusInternalServerError)
			}))
			system := types.System{
				ID: systemID,
			}
			return server, system, check(hasError)
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			server, sys, checkFns := tc(t)
			defer server.Close()

			client, _ := NewClientWithArgs(server.URL, "", math.MaxInt64, true, false, "")
			system := System{
				System: &sys,
				client: client,
			}
			err := system.DeleteNvmeHost("mock-id")

			for _, checkFn := range checkFns {
				checkFn(t, err)
			}
		})
	}
}

func Test_GetHostNvmeControllers(t *testing.T) {
	type checkFn func(*testing.T, []types.NvmeController, error)
	check := func(fns ...checkFn) []checkFn { return fns }

	responseJSON := `[
	{
		"name": null,
		"isConnected": true,
		"sdtId": "cd16ff4300000000",
		"hostIp": null,
		"hostId": "4d2a627100010000",
		"controllerId": 1,
		"sysPortId": 2,
		"sysPortIp": "10.0.0.22",
		"subsystem": "Io",
		"isAssigned": true,
		"id": "cc00000001000011",
		"links": [
			{
				"rel": "self",
				"href": "/api/instances/NvmeController::cc00000001000011"
			}
		]
	}]`

	hasNoError := func(t *testing.T, _ []types.NvmeController, err error) {
		if err != nil {
			t.Fatalf("expected no error")
		}
	}

	checkLength := func(length int) func(t *testing.T, slice []types.NvmeController, err error) {
		return func(t *testing.T, slice []types.NvmeController, _ error) {
			assert.Equal(t, length, len(slice))
		}
	}

	checkSdtID := func(index int, id string) func(t *testing.T, slice []types.NvmeController, err error) {
		return func(t *testing.T, slice []types.NvmeController, _ error) {
			assert.Equal(t, slice[index].SdtID, id)
		}
	}

	checkHostID := func(index int, id string) func(t *testing.T, slice []types.NvmeController, err error) {
		return func(t *testing.T, slice []types.NvmeController, _ error) {
			assert.Equal(t, slice[index].HostID, id)
		}
	}

	checkConnected := func(index int, connected bool) func(t *testing.T, slice []types.NvmeController, err error) {
		return func(t *testing.T, slice []types.NvmeController, _ error) {
			assert.Equal(t, slice[index].IsConnected, connected)
		}
	}

	hasError := func(t *testing.T, _ []types.NvmeController, err error) {
		if err == nil {
			t.Fatalf("expected error")
		}
	}

	tests := map[string]func(t *testing.T) (*httptest.Server, types.System, []checkFn){
		"success": func(_ *testing.T) (*httptest.Server, types.System, []checkFn) {
			systemID := "mock-system-id"
			id := "mock-id"
			url := fmt.Sprintf("/api/instances/Host::%v/relationships/NvmeController", id)
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method == http.MethodGet && strings.EqualFold(r.URL.Path, url) {
					w.WriteHeader(http.StatusOK)
					_, err := w.Write([]byte(responseJSON))
					if err != nil {
						t.Fatalf("Unexpected error: %v", err)
					}
					return
				}
				http.NotFound(w, r)
			}))
			system := types.System{
				ID: systemID,
			}
			return server, system, check(hasNoError, checkLength(1), checkSdtID(0, "cd16ff4300000000"), checkHostID(0, "4d2a627100010000"), checkConnected(0, true))
		},
		"error response": func(_ *testing.T) (*httptest.Server, types.System, []checkFn) {
			systemID := "mock-system-id"
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				http.Error(w, "Server Error", http.StatusInternalServerError)
			}))
			system := types.System{
				ID: systemID,
			}
			return server, system, check(hasError)
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			server, sys, checkFns := tc(t)
			defer server.Close()

			client, _ := NewClientWithArgs(server.URL, "", math.MaxInt64, true, false, "")
			system := System{
				System: &sys,
				client: client,
			}
			host := types.NvmeHost{
				ID: "mock-id",
			}
			nvmeHostControllers, err := system.GetHostNvmeControllers(host)

			for _, checkFn := range checkFns {
				checkFn(t, nvmeHostControllers, err)
			}
		})
	}
}

func TestGetNewNvmeHost(t *testing.T) {
	assert.NotNil(t, NewNvmeHost(nil, nil))
}

func TestMapVolumeNVMe(t *testing.T) {
	volumeID := uuid.NewString()
	mapVolumeNVMeParam := &types.MapVolumeNVMeParam{}
	type testCase struct {
		server      *httptest.Server
		expectedErr error
	}

	cases := map[string]testCase{
		"succeed": {
			server: httptest.NewServer(http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
				switch req.RequestURI {
				case fmt.Sprintf("/api/instances/Volume::%s/action/addMappedHost", volumeID):
					content, err := json.Marshal(types.NvmeHost{
						ID: volumeID,
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
		"error: bad request": {
			server: httptest.NewServer(http.HandlerFunc(func(resp http.ResponseWriter, _ *http.Request) {
				resp.WriteHeader(http.StatusBadRequest)
				resp.Write([]byte(`{"message":"bad request","httpStatusCode":400,"errorCode":0}`))
			})),
			expectedErr: fmt.Errorf("bad request"),
		},
		"error: could not find nvme host": {
			server: httptest.NewServer(http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
				switch req.RequestURI {
				case fmt.Sprintf("/api/instances/Volume::%s/action/addMappedHost", volumeID):
					content, err := json.Marshal([]types.NvmeHost{
						{
							ID: uuid.NewString(),
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
			expectedErr: fmt.Errorf("Couldn't find NVMe host"),
		},
	}

	for _, tc := range cases {
		client, err := NewClientWithArgs(tc.server.URL, "3.6", math.MaxInt64, true, false, "")
		if err != nil {
			t.Fatal(err)
		}
		defer tc.server.Close()

		volume := Volume{
			Volume: &types.Volume{
				ID: volumeID,
			},
			client: client,
		}

		err = volume.MapVolumeNVMe(mapVolumeNVMeParam)
		if err != nil {
			if tc.expectedErr.Error() != err.Error() {
				t.Fatal(err)
			}
		}
	}
}

func TestRemoveMappedHost(t *testing.T) {
	volumeID := uuid.NewString()
	unmapVolumeNVMeParam := &types.UnmapVolumeNVMeParam{}

	type testCase struct {
		server      *httptest.Server
		expectedErr error
	}

	cases := map[string]testCase{
		"succeed": {
			server: httptest.NewServer(http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
				if req.RequestURI == fmt.Sprintf("/api/instances/Volume::%s/action/removeMappedHost", volumeID) {
					resp.WriteHeader(http.StatusOK)
					resp.Write([]byte(`{}`)) // success response
				} else {
					resp.WriteHeader(http.StatusBadRequest)
					resp.Write([]byte(`{"message":"no route handled"}`))
				}
			})),
			expectedErr: nil,
		},
		"error: bad request": {
			server: httptest.NewServer(http.HandlerFunc(func(resp http.ResponseWriter, _ *http.Request) {
				resp.WriteHeader(http.StatusBadRequest)
				resp.Write([]byte(`{"message":"bad request"}`))
			})),
			expectedErr: fmt.Errorf("bad request"),
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			client, err := NewClientWithArgs(tc.server.URL, "3.6", math.MaxInt64, true, false, "")
			if err != nil {
				t.Fatal(err)
			}
			defer tc.server.Close()

			volume := Volume{
				Volume: &types.Volume{ID: volumeID},
				client: client,
			}

			err = volume.RemoveMappedHost(unmapVolumeNVMeParam)
			if err != nil {
				if tc.expectedErr == nil || tc.expectedErr.Error() != err.Error() {
					t.Fatalf("expected %v, got %v", tc.expectedErr, err)
				}
			}
		})
	}
}
