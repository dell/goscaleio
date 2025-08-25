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
	"fmt"
	"math"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	types "github.com/dell/goscaleio/types/v1"
	"github.com/stretchr/testify/assert"
)

func Test_GetAllSdts(t *testing.T) {
	type checkFn func(*testing.T, []types.Sdt, error)
	check := func(fns ...checkFn) []checkFn { return fns }

	responseJSON := `[
	{
		"ipList": [
			{
				"role": "StorageAndHost",
				"ip": "10.0.0.23"
			}
		],
		"sdtState": "Normal",
		"systemId": "mock-system-id",
		"name": "yulan1sdt3",
		"protectionDomainId": "5e4640aa00000000",
		"storagePort": 12200,
		"nvmePort": 4420,
		"discoveryPort": 8009,
		"certificateInfo": {},
		"mdmConnectionState": "Connected",
		"membershipState": "Joined",
		"faultSetId": null,
		"softwareVersionInfo": "R4_5.2100.0",
		"maintenanceState": "NoMaintenance",
		"authenticationError": "None",
		"persistentDiscoveryControllersNum": 0,
		"id": "mock-id",
		"links": [
			{
				"rel": "self",
				"href": "/api/instances/Sdt::mock-id"
			},
			{
				"rel": "/api/Sdt/relationship/Statistics",
				"href": "/api/instances/Sdt::mock-id/relationships/Statistics"
			},
			{
				"rel": "/api/parent/relationship/protectionDomainId",
				"href": "/api/instances/ProtectionDomain::5e4640aa00000000"
			}
		]
	}
]`

	hasNoError := func(t *testing.T, _ []types.Sdt, err error) {
		if err != nil {
			t.Fatalf("expected no error")
		}
	}

	checkLength := func(length int) func(t *testing.T, hosts []types.Sdt, err error) {
		return func(t *testing.T, hosts []types.Sdt, _ error) {
			assert.Equal(t, length, len(hosts))
		}
	}

	hasError := func(t *testing.T, _ []types.Sdt, err error) {
		if err == nil {
			t.Fatalf("expected error")
		}
	}

	tests := map[string]func(t *testing.T) (*httptest.Server, types.System, []checkFn){
		"success": func(t *testing.T) (*httptest.Server, types.System, []checkFn) {
			systemID := "mock-system-id"
			url := "/api/types/Sdt/instances"
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

			client, _ := NewClientWithArgs(server.URL, "", math.MaxInt64, true, false)
			system := System{
				System: &sys,
				client: client,
			}
			sdts, err := system.GetAllSdts()

			for _, checkFn := range checkFns {
				checkFn(t, sdts, err)
			}
		})
	}
}

func Test_GetSdtByID(t *testing.T) {
	type checkFn func(*testing.T, *types.Sdt, error)
	check := func(fns ...checkFn) []checkFn { return fns }

	responseJSON := `{
		"ipList": [
			{
				"role": "StorageAndHost",
				"ip": "10.0.0.23"
			}
		],
		"sdtState": "Normal",
		"systemId": "mock-system-id",
		"name": "mock-name",
		"protectionDomainId": "mock-pd-id",
		"storagePort": 12200,
		"nvmePort": 4420,
		"discoveryPort": 8009,
		"certificateInfo": {},
		"mdmConnectionState": "Connected",
		"membershipState": "Joined",
		"faultSetId": null,
		"softwareVersionInfo": "R4_5.2100.0",
		"maintenanceState": "NoMaintenance",
		"authenticationError": "None",
		"persistentDiscoveryControllersNum": 0,
		"id": "mock-id",
		"links": [
			{
				"rel": "self",
				"href": "/api/instances/Sdt::mock-id"
			},
			{
				"rel": "/api/Sdt/relationship/Statistics",
				"href": "/api/instances/Sdt::mock-id/relationships/Statistics"
			},
			{
				"rel": "/api/parent/relationship/protectionDomainId",
				"href": "/api/instances/ProtectionDomain::mock-pd-id"
			}
		]
	}`

	hasNoError := func(t *testing.T, _ *types.Sdt, err error) {
		if err != nil {
			t.Fatalf("expected no error")
		}
	}

	checkName := func(name string) func(t *testing.T, host *types.Sdt, err error) {
		return func(t *testing.T, sdt *types.Sdt, _ error) {
			assert.Equal(t, sdt.Name, name)
		}
	}

	hasError := func(t *testing.T, _ *types.Sdt, err error) {
		if err == nil {
			t.Fatalf("expected error")
		}
	}

	tests := map[string]func(t *testing.T) (*httptest.Server, types.System, []checkFn){
		"success": func(t *testing.T) (*httptest.Server, types.System, []checkFn) {
			systemID := "mock-system-id"
			id := "mock-id"
			url := fmt.Sprintf("/api/instances/Sdt::%v", id)
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

			client, _ := NewClientWithArgs(server.URL, "", math.MaxInt64, true, false)
			system := System{
				System: &sys,
				client: client,
			}
			sdt, err := system.GetSdtByID("mock-id")

			for _, checkFn := range checkFns {
				checkFn(t, sdt, err)
			}
		})
	}
}

func Test_CreateSdt(t *testing.T) {
	type checkFn func(*testing.T, *types.SdtResp, error)
	check := func(fns ...checkFn) []checkFn { return fns }

	hasNoError := func(t *testing.T, _ *types.SdtResp, err error) {
		if err != nil {
			t.Fatalf("expected no error")
		}
	}

	checkID := func(id string) func(t *testing.T, resp *types.SdtResp, err error) {
		return func(t *testing.T, resp *types.SdtResp, _ error) {
			assert.Equal(t, resp.ID, id)
		}
	}

	checkEmptyIPList := func(t *testing.T, _ *types.SdtResp, err error) {
		assert.Equal(t, err.Error(), "Must provide at least 1 SDT IP")
	}

	hasError := func(t *testing.T, _ *types.SdtResp, err error) {
		if err == nil {
			t.Fatalf("expected error")
		}
	}

	tests := map[string]func(t *testing.T) (*httptest.Server, *types.SdtParam, []checkFn){
		"success": func(t *testing.T) (*httptest.Server, *types.SdtParam, []checkFn) {
			sdtParam := &types.SdtParam{
				Name:               "example-sdt",
				IPList:             []*types.SdtIP{{IP: "192.168.0.1", Role: "StorageAndHost"}},
				StoragePort:        12200,
				NvmePort:           4420,
				DiscoveryPort:      8009,
				ProtectionDomainID: "mock-pd-id",
			}

			url := "/api/types/Sdt/instances"
			id := "mock-id"
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
			return server, sdtParam, check(hasNoError, checkID(id))
		},

		"empty IP list": func(_ *testing.T) (*httptest.Server, *types.SdtParam, []checkFn) {
			sdtParam := &types.SdtParam{
				Name:               "example-sdt",
				IPList:             []*types.SdtIP{},
				StoragePort:        12200,
				NvmePort:           4420,
				DiscoveryPort:      8009,
				ProtectionDomainID: "mock-pd-id",
			}

			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				http.Error(w, "Server Error", http.StatusInternalServerError)
			}))
			return server, sdtParam, check(checkEmptyIPList)
		},

		"error response": func(_ *testing.T) (*httptest.Server, *types.SdtParam, []checkFn) {
			sdtParam := &types.SdtParam{
				Name:               "example-sdt",
				IPList:             []*types.SdtIP{{IP: "192.168.0.1", Role: "StorageAndHost"}},
				StoragePort:        12200,
				NvmePort:           4420,
				DiscoveryPort:      8009,
				ProtectionDomainID: "mock-pd-id",
			}

			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				http.Error(w, "Server Error", http.StatusInternalServerError)
			}))
			return server, sdtParam, check(hasError)
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			server, sdtParam, checkFns := tc(t)
			defer server.Close()

			client, _ := NewClientWithArgs(server.URL, "", math.MaxInt64, true, false)
			pd := NewProtectionDomain(client)
			resp, err := pd.CreateSdt(sdtParam)

			for _, checkFn := range checkFns {
				checkFn(t, resp, err)
			}
		})
	}
}

func Test_RenameSdt(t *testing.T) {
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
			url := fmt.Sprintf("/api/instances/Sdt::%v/action/renameSdt", id)
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

			client, _ := NewClientWithArgs(server.URL, "", math.MaxInt64, true, false)
			system := System{
				System: &sys,
				client: client,
			}
			err := system.RenameSdt("mock-id", "new-mock-name")

			for _, checkFn := range checkFns {
				checkFn(t, err)
			}
		})
	}
}

func Test_SetSdtNvmePort(t *testing.T) {
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
			url := fmt.Sprintf("/api/instances/Sdt::%v/action/modifyNvmePort", id)
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

			client, _ := NewClientWithArgs(server.URL, "", math.MaxInt64, true, false)
			system := System{
				System: &sys,
				client: client,
			}
			err := system.SetSdtNvmePort("mock-id", 1234)

			for _, checkFn := range checkFns {
				checkFn(t, err)
			}
		})
	}
}

func Test_SetSdtStoragePort(t *testing.T) {
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
			url := fmt.Sprintf("/api/instances/Sdt::%v/action/modifyStoragePort", id)
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

			client, _ := NewClientWithArgs(server.URL, "", math.MaxInt64, true, false)
			system := System{
				System: &sys,
				client: client,
			}
			err := system.SetSdtStoragePort("mock-id", 1234)

			for _, checkFn := range checkFns {
				checkFn(t, err)
			}
		})
	}
}

func Test_SetSdtDiscoveryPort(t *testing.T) {
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
			url := fmt.Sprintf("/api/instances/Sdt::%v/action/modifyDiscoveryPort", id)
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

			client, _ := NewClientWithArgs(server.URL, "", math.MaxInt64, true, false)
			system := System{
				System: &sys,
				client: client,
			}
			err := system.SetSdtDiscoveryPort("mock-id", 1234)

			for _, checkFn := range checkFns {
				checkFn(t, err)
			}
		})
	}
}

func Test_AddSdtTargetIP(t *testing.T) {
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
			url := fmt.Sprintf("/api/instances/Sdt::%v/action/addIp", id)
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

			client, _ := NewClientWithArgs(server.URL, "", math.MaxInt64, true, false)
			system := System{
				System: &sys,
				client: client,
			}
			err := system.AddSdtTargetIP("mock-id", "192.168.0.2", "StorageAndHost")

			for _, checkFn := range checkFns {
				checkFn(t, err)
			}
		})
	}
}

func Test_RemoveSdtTargetIP(t *testing.T) {
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
			url := fmt.Sprintf("/api/instances/Sdt::%v/action/removeIp", id)
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

			client, _ := NewClientWithArgs(server.URL, "", math.MaxInt64, true, false)
			system := System{
				System: &sys,
				client: client,
			}
			err := system.RemoveSdtTargetIP("mock-id", "192.168.0.2")

			for _, checkFn := range checkFns {
				checkFn(t, err)
			}
		})
	}
}

func Test_ModifySdtIPRole(t *testing.T) {
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
			url := fmt.Sprintf("/api/instances/Sdt::%v/action/modifyIpRole", id)
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

			client, _ := NewClientWithArgs(server.URL, "", math.MaxInt64, true, false)
			system := System{
				System: &sys,
				client: client,
			}
			err := system.ModifySdtIPRole("mock-id", "192.168.0.2", "StorageOnly")

			for _, checkFn := range checkFns {
				checkFn(t, err)
			}
		})
	}
}

func Test_DeleteSdt(t *testing.T) {
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
			url := fmt.Sprintf("/api/instances/Sdt::%v/action/removeSdt", id)
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

			client, _ := NewClientWithArgs(server.URL, "", math.MaxInt64, true, false)
			system := System{
				System: &sys,
				client: client,
			}
			err := system.DeleteSdt("mock-id")

			for _, checkFn := range checkFns {
				checkFn(t, err)
			}
		})
	}
}

func Test_EnterSdtMaintenanceMode(t *testing.T) {
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
			url := fmt.Sprintf("/api/instances/Sdt::%v/action/enterMaintenanceMode", id)
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

			client, _ := NewClientWithArgs(server.URL, "", math.MaxInt64, true, false)
			system := System{
				System: &sys,
				client: client,
			}
			err := system.EnterSdtMaintenanceMode("mock-id")

			for _, checkFn := range checkFns {
				checkFn(t, err)
			}
		})
	}
}

func Test_ExitSdtMaintenanceMode(t *testing.T) {
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
			url := fmt.Sprintf("/api/instances/Sdt::%v/action/exitMaintenanceMode", id)
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

			client, _ := NewClientWithArgs(server.URL, "", math.MaxInt64, true, false)
			system := System{
				System: &sys,
				client: client,
			}
			err := system.ExitSdtMaintenanceMode("mock-id")

			for _, checkFn := range checkFns {
				checkFn(t, err)
			}
		})
	}
}

func TestNewSdt(t *testing.T) {
	assert.NotNil(t, NewSdt(nil, nil))
}
