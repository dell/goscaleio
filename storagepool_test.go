// Copyright © 2019 - 2022 Dell Inc. or its subsidiaries. All Rights Reserved.
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
	"testing"

	types "github.com/dell/goscaleio/types/v1"
	"github.com/stretchr/testify/assert"
)

func SetUpProtectionDomain(url string) (ProtectionDomain, error) {
	pdID := "12345678-1234-1234-1234-123456789012"
	protectionDomain := &types.ProtectionDomain{
		Name: "domain1",
		ID:   pdID,
	}

	client, err := NewClientWithArgs(url, "", math.MaxInt64, true, false)
	if err != nil {
		return ProtectionDomain{}, err
	}
	pd := ProtectionDomain{
		ProtectionDomain: protectionDomain,
		client:           client,
	}
	return pd, nil
}

func TestCreateStoragePool(t *testing.T) {
	type checkFn func(*testing.T, string, error)
	check := func(fns ...checkFn) []checkFn { return fns }

	hasNoError := func(t *testing.T, _ string, err error) {
		if err != nil {
			t.Fatalf("expected no error")
		}
	}

	hasError := func(t *testing.T, _ string, err error) {
		if err == nil {
			t.Fatalf("expected error")
		}
	}

	checkResp := func(storagePoolID string) func(t *testing.T, resp string, err error) {
		return func(t *testing.T, resp string, _ error) {
			assert.Equal(t, storagePoolID, resp)
		}
	}

	tests := map[string]func(t *testing.T) (*httptest.Server, *types.ProtectionDomain, []checkFn){
		"success": func(t *testing.T) (*httptest.Server, *types.ProtectionDomain, []checkFn) {
			href := "/api/types/StoragePool/instances"
			protectionDomain := types.ProtectionDomain{
				Name: "domain1",
			}

			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method != http.MethodPost {
					t.Fatal(fmt.Errorf("wrong method. Expected %s; but got %s", http.MethodGet, r.Method))
				}

				if r.URL.Path != href {
					t.Fatal(fmt.Errorf("wrong path. Expected %s; but got %s", href, r.URL.Path))
				}

				resp := types.StoragePoolResp{
					ID: "1a2b345c-123a-123a-ab1c-abc1fd234567e",
				}
				respData, err := json.Marshal(resp)
				if err != nil {
					t.Fatal(err)
				}
				_, err = fmt.Fprintln(w, string(respData))
				if err != nil {
					return
				}
			}))
			return ts, &protectionDomain, check(hasNoError, checkResp("1a2b345c-123a-123a-ab1c-abc1fd234567e"))
		},
		"bad request": func(*testing.T) (*httptest.Server, *types.ProtectionDomain, []checkFn) {
			protectionDomain := types.ProtectionDomain{
				Name: "domain1",
			}
			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				w.WriteHeader(http.StatusBadRequest)
			}))
			return ts, &protectionDomain, check(hasError)
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			ts, protectionDomain, checkFns := tc(t)
			defer ts.Close()

			client, err := NewClientWithArgs(ts.URL, "", math.MaxInt64, true, false)
			if err != nil {
				t.Fatal(err)
			}

			pd := ProtectionDomain{
				ProtectionDomain: protectionDomain,
				client:           client,
			}

			storagePoolParams := &types.StoragePoolParam{
				Name: "pool1",
			}

			resp, err := pd.CreateStoragePool(storagePoolParams)
			for _, checkFn := range checkFns {
				checkFn(t, resp, err)
			}
		})
	}
}

func TestModifyStoragePoolName(t *testing.T) {
	poolID := "1a2b345c-123a-123a-ab1c-abc1fd234567e"
	type checkFn func(*testing.T, string, error)
	check := func(fns ...checkFn) []checkFn { return fns }

	hasNoError := func(t *testing.T, _ string, err error) {
		if err != nil {
			t.Fatalf("expected no error")
		}
	}

	hasError := func(t *testing.T, _ string, err error) {
		if err == nil {
			t.Fatalf("expected error")
		}
	}

	checkResp := func(storagePoolID string) func(t *testing.T, resp string, err error) {
		return func(t *testing.T, resp string, _ error) {
			assert.Equal(t, storagePoolID, resp)
		}
	}

	tests := map[string]func(t *testing.T) (*httptest.Server, []checkFn){
		"success": func(t *testing.T) (*httptest.Server, []checkFn) {
			href := fmt.Sprintf("/api/instances/StoragePool::%v/action/setStoragePoolName", poolID)
			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method != http.MethodPost {
					t.Fatal(fmt.Errorf("wrong method. Expected %s; but got %s", http.MethodGet, r.Method))
				}
				if r.URL.Path != href {
					t.Fatal(fmt.Errorf("wrong path. Expected %s; but got %s", href, r.URL.Path))
				}

				resp := types.StoragePoolResp{
					ID: poolID,
				}
				respData, err := json.Marshal(resp)
				if err != nil {
					t.Fatal(err)
				}
				_, err = fmt.Fprintln(w, string(respData))
				if err != nil {
					return
				}
			}))
			return ts, check(hasNoError, checkResp(poolID))
		},
		"bad request": func(*testing.T) (*httptest.Server, []checkFn) {
			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				w.WriteHeader(http.StatusBadRequest)
			}))
			return ts, check(hasError)
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			ts, checkFns := tc(t)
			defer ts.Close()

			client, err := NewClientWithArgs(ts.URL, "", math.MaxInt64, true, false)
			if err != nil {
				t.Fatal(err)
			}

			pd := ProtectionDomain{
				client: client,
			}

			resp, err := pd.ModifyStoragePoolName(poolID, "pool1")
			for _, checkFn := range checkFns {
				checkFn(t, resp, err)
			}
		})
	}
}

func TestDeleteStoragePool(t *testing.T) {
	poolID := "1a2b345c-123a-123a-ab1c-abc1fd234567e"
	protectionDomain := &types.ProtectionDomain{
		Name: "domain1",
		Links: []*types.Link{
			{
				Rel:  "/api/ProtectionDomain/relationship/StoragePool",
				HREF: fmt.Sprintf("/api/instances/StoragePool::%s", poolID),
			},
		},
	}

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case fmt.Sprintf("/api/instances/StoragePool::%s", poolID):
			resp := []types.StoragePool{
				{
					ID:   poolID,
					Name: "pool1",
					Links: []*types.Link{
						{
							Rel:  "self",
							HREF: fmt.Sprintf("/api/instances/StoragePool::%s", poolID),
						},
					},
				},
			}
			respData, err := json.Marshal(resp)
			if err != nil {
				t.Fatal(err)
			}
			_, err = fmt.Fprintln(w, string(respData))
			if err != nil {
				return
			}
			return
		}
	}))
	defer ts.Close()

	client, err := NewClientWithArgs(ts.URL, "", math.MaxInt64, true, false)
	if err != nil {
		t.Fatal(err)
	}

	pd := ProtectionDomain{
		ProtectionDomain: protectionDomain,
		client:           client,
	}
	err = pd.DeleteStoragePool("pool1")
	assert.Nil(t, err)
}

func TestModifyStoragePoolMedia(t *testing.T) {
	poolID := "1a2b345c-123a-123a-ab1c-abc1fd234567e"
	type checkFn func(*testing.T, string, error)
	check := func(fns ...checkFn) []checkFn { return fns }

	hasNoError := func(t *testing.T, _ string, err error) {
		if err != nil {
			t.Fatalf("expected no error")
		}
	}

	hasError := func(t *testing.T, _ string, err error) {
		if err == nil {
			t.Fatalf("expected error")
		}
	}

	checkResp := func(storagePoolID string) func(t *testing.T, resp string, err error) {
		return func(t *testing.T, resp string, _ error) {
			assert.Equal(t, storagePoolID, resp)
		}
	}

	tests := map[string]func(t *testing.T) (*httptest.Server, []checkFn){
		"success": func(t *testing.T) (*httptest.Server, []checkFn) {
			href := fmt.Sprintf("/api/instances/StoragePool::%v/action/setMediaType", poolID)
			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method != http.MethodPost {
					t.Fatal(fmt.Errorf("wrong method. Expected %s; but got %s", http.MethodGet, r.Method))
				}
				if r.URL.Path != href {
					t.Fatal(fmt.Errorf("wrong path. Expected %s; but got %s", href, r.URL.Path))
				}

				resp := types.StoragePoolResp{
					ID: poolID,
				}
				respData, err := json.Marshal(resp)
				if err != nil {
					t.Fatal(err)
				}
				_, err = fmt.Fprintln(w, string(respData))
				if err != nil {
					return
				}
			}))
			return ts, check(hasNoError, checkResp(poolID))
		},
		"bad request": func(*testing.T) (*httptest.Server, []checkFn) {
			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				w.WriteHeader(http.StatusBadRequest)
			}))
			return ts, check(hasError)
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			ts, checkFns := tc(t)
			defer ts.Close()

			client, err := NewClientWithArgs(ts.URL, "", math.MaxInt64, true, false)
			if err != nil {
				t.Fatal(err)
			}

			pd := ProtectionDomain{
				client: client,
			}

			resp, err := pd.ModifyStoragePoolMedia(poolID, "media1")
			for _, checkFn := range checkFns {
				checkFn(t, resp, err)
			}
		})
	}
}

func TestSetReplicationJournalCapacity(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	}))
	defer ts.Close()

	pd, err := SetUpProtectionDomain(ts.URL)
	if err != nil {
		t.Fatal(err)
	}

	err = pd.SetReplicationJournalCapacity("", "")
	assert.Nil(t, err)
}

func TestEnableOrDisableZeroPadding(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	}))
	defer ts.Close()

	pd, err := SetUpProtectionDomain(ts.URL)
	if err != nil {
		t.Fatal(err)
	}
	err = pd.EnableOrDisableZeroPadding("", "")
	assert.Nil(t, err)
}

func TestSetCapacityAlertThreshold(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	}))
	defer ts.Close()

	pd, err := SetUpProtectionDomain(ts.URL)
	if err != nil {
		t.Fatal(err)
	}

	capacityThresholdParameters := &types.CapacityAlertThresholdParam{}
	err = pd.SetCapacityAlertThreshold("", capacityThresholdParameters)
	assert.Nil(t, err)
}

func TestSetProtectedMaintenanceModeIoPriorityPolicy(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	}))
	defer ts.Close()

	pd, err := SetUpProtectionDomain(ts.URL)
	if err != nil {
		t.Fatal(err)
	}

	protectedMaintenanceModeParams := &types.ProtectedMaintenanceModeParam{}
	err = pd.SetProtectedMaintenanceModeIoPriorityPolicy("", protectedMaintenanceModeParams)
	assert.Nil(t, err)
}

func TestGetPoolStorage(t *testing.T) {
	poolID := "1a2b345c-123a-123a-ab1c-abc1fd234567e"
	tests := []struct {
		name             string
		storagePoolHref  string
		protectionDomain *types.ProtectionDomain
		ts               *httptest.Server
	}{
		{
			name:            "Success: Empty StoragePoolHref",
			storagePoolHref: "",
			protectionDomain: &types.ProtectionDomain{
				Name: "domain1",
				Links: []*types.Link{
					{
						Rel:  "/api/ProtectionDomain/relationship/StoragePool",
						HREF: fmt.Sprintf("/api/instances/StoragePool::%s", poolID),
					},
				},
			},
			ts: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				resp := []types.StoragePool{
					{
						ID:   poolID,
						Name: "pool1",
					},
				}

				respData, err := json.Marshal(resp)
				if err != nil {
					t.Fatal(err)
				}
				_, err = fmt.Fprintln(w, string(respData))
				if err != nil {
					return
				}
			})),
		},
		{
			name:            "Success: StoragePoolHref",
			storagePoolHref: fmt.Sprintf("/api/instances/StoragePool::%s", poolID),
			protectionDomain: &types.ProtectionDomain{
				Name: "domain1",
			},
			ts: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				resp := types.StoragePool{
					ID:   poolID,
					Name: "pool1",
				}

				respData, err := json.Marshal(resp)
				if err != nil {
					t.Fatal(err)
				}
				_, err = fmt.Fprintln(w, string(respData))
				if err != nil {
					return
				}
			})),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ts := tt.ts
			defer ts.Close()

			client, err := NewClientWithArgs(ts.URL, "", math.MaxInt64, true, false)
			if err != nil {
				t.Fatal(err)
			}

			pd := ProtectionDomain{
				ProtectionDomain: tt.protectionDomain,
				client:           client,
			}
			_, err = pd.GetStoragePool(tt.storagePoolHref)
			assert.Nil(t, err)
		})
	}
}

func TestGetStoragePoolByID(t *testing.T) {
	poolID := "1a2b345c-123a-123a-ab1c-abc1fd234567e"
	type checkFn func(*testing.T, *types.StoragePool, error)
	check := func(fns ...checkFn) []checkFn { return fns }

	hasNoError := func(t *testing.T, _ *types.StoragePool, err error) {
		if err != nil {
			t.Fatalf("expected no error")
		}
	}

	hasError := func(t *testing.T, _ *types.StoragePool, err error) {
		if err == nil {
			t.Fatalf("expected error")
		}
	}

	checkResp := func(storagePoolID string) func(t *testing.T, resp *types.StoragePool, err error) {
		return func(t *testing.T, resp *types.StoragePool, _ error) {
			assert.Equal(t, storagePoolID, resp.ID)
		}
	}

	tests := map[string]func(t *testing.T) (*httptest.Server, *types.System, []checkFn){
		"success": func(t *testing.T) (*httptest.Server, *types.System, []checkFn) {
			systemID := "0000aaacccddd1111"
			href := fmt.Sprintf("/api/instances/StoragePool::%s", poolID)
			system := types.System{
				ID: systemID,
			}

			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method != http.MethodGet {
					t.Fatal(fmt.Errorf("wrong method. Expected %s; but got %s", http.MethodGet, r.Method))
				}

				if r.URL.Path != href {
					t.Fatal(fmt.Errorf("wrong path. Expected %s; but got %s", href, r.URL.Path))
				}

				resp := types.StoragePool{
					ID: "1a2b345c-123a-123a-ab1c-abc1fd234567e",
				}
				respData, err := json.Marshal(resp)
				if err != nil {
					t.Fatal(err)
				}
				_, err = fmt.Fprintln(w, string(respData))
				if err != nil {
					return
				}
			}))
			return ts, &system, check(hasNoError, checkResp("1a2b345c-123a-123a-ab1c-abc1fd234567e"))
		},
		"bad request": func(*testing.T) (*httptest.Server, *types.System, []checkFn) {
			systemID := "0000aaacccddd1111"
			system := types.System{
				ID: systemID,
			}
			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				w.WriteHeader(http.StatusBadRequest)
			}))
			return ts, &system, check(hasError)
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			ts, system, checkFns := tc(t)
			defer ts.Close()

			client, err := NewClientWithArgs(ts.URL, "", math.MaxInt64, true, false)
			if err != nil {
				t.Fatal(err)
			}

			s := System{
				System: system,
				client: client,
			}

			resp, err := s.GetStoragePoolByID(poolID)
			for _, checkFn := range checkFns {
				checkFn(t, resp, err)
			}
		})
	}
}

func TestGetAllStoragePools(t *testing.T) {
	type checkFn func(*testing.T, []types.StoragePool, error)
	check := func(fns ...checkFn) []checkFn { return fns }

	hasNoError := func(t *testing.T, _ []types.StoragePool, err error) {
		if err != nil {
			t.Fatalf("expected no error")
		}
	}

	hasError := func(t *testing.T, _ []types.StoragePool, err error) {
		if err == nil {
			t.Fatalf("expected error")
		}
	}

	checkResp := func(storagePoolID string) func(t *testing.T, resp []types.StoragePool, err error) {
		return func(t *testing.T, resp []types.StoragePool, _ error) {
			assert.Equal(t, storagePoolID, resp[0].ID)
		}
	}

	tests := map[string]func(t *testing.T) (*httptest.Server, *types.System, []checkFn){
		"success": func(t *testing.T) (*httptest.Server, *types.System, []checkFn) {
			systemID := "0000aaacccddd1111"
			href := "/api/types/StoragePool/instances"
			system := types.System{
				ID: systemID,
			}

			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method != http.MethodGet {
					t.Fatal(fmt.Errorf("wrong method. Expected %s; but got %s", http.MethodGet, r.Method))
				}

				if r.URL.Path != href {
					t.Fatal(fmt.Errorf("wrong path. Expected %s; but got %s", href, r.URL.Path))
				}

				resp := []types.StoragePool{
					{
						ID: "1a2b345c-123a-123a-ab1c-abc1fd234567e",
					},
				}
				respData, err := json.Marshal(resp)
				if err != nil {
					t.Fatal(err)
				}
				_, err = fmt.Fprintln(w, string(respData))
				if err != nil {
					return
				}
			}))
			return ts, &system, check(hasNoError, checkResp("1a2b345c-123a-123a-ab1c-abc1fd234567e"))
		},
		"bad request": func(*testing.T) (*httptest.Server, *types.System, []checkFn) {
			systemID := "0000aaacccddd1111"
			system := types.System{
				ID: systemID,
			}
			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				w.WriteHeader(http.StatusBadRequest)
			}))
			return ts, &system, check(hasError)
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			ts, system, checkFns := tc(t)
			defer ts.Close()

			client, err := NewClientWithArgs(ts.URL, "", math.MaxInt64, true, false)
			if err != nil {
				t.Fatal(err)
			}

			s := System{
				System: system,
				client: client,
			}

			resp, err := s.GetAllStoragePools()
			for _, checkFn := range checkFns {
				checkFn(t, resp, err)
			}
		})
	}
}

func TestSetRebalanceEnabled(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	}))
	defer ts.Close()

	pd, err := SetUpProtectionDomain(ts.URL)
	if err != nil {
		t.Fatal(err)
	}

	err = pd.SetRebalanceEnabled("", "")
	assert.Nil(t, err)
}

func TestSetRebalanceIoPriorityPolicy(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	}))
	defer ts.Close()

	pd, err := SetUpProtectionDomain(ts.URL)
	if err != nil {
		t.Fatal(err)
	}

	protectedMaintenanceModeParam := &types.ProtectedMaintenanceModeParam{}
	err = pd.SetRebalanceIoPriorityPolicy("", protectedMaintenanceModeParam)
	assert.Nil(t, err)
}

func TestSetVTreeMigrationIOPriorityPolicy(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	}))
	defer ts.Close()

	pd, err := SetUpProtectionDomain(ts.URL)
	if err != nil {
		t.Fatal(err)
	}

	protectedMaintenanceModeParam := &types.ProtectedMaintenanceModeParam{}
	err = pd.SetVTreeMigrationIOPriorityPolicy("", protectedMaintenanceModeParam)
	assert.Nil(t, err)
}

func TestSetSparePercentage(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	}))
	defer ts.Close()

	pd, err := SetUpProtectionDomain(ts.URL)
	if err != nil {
		t.Fatal(err)
	}

	err = pd.SetSparePercentage("", "")
	assert.Nil(t, err)
}

func TestSetRMcacheWriteHandlingMode(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	}))
	defer ts.Close()

	pd, err := SetUpProtectionDomain(ts.URL)
	if err != nil {
		t.Fatal(err)
	}

	err = pd.SetRMcacheWriteHandlingMode("", "")
	assert.Nil(t, err)
}

func TestModifyRMCache(t *testing.T) {
	poolID := "1a2b345c-123a-123a-ab1c-abc1fd234567e"
	storagePool := &types.StoragePool{
		Name: "pool1",
		ID:   poolID,
		Links: []*types.Link{
			{
				Rel:  "self",
				HREF: fmt.Sprintf("/api/instances/StoragePool::%s", poolID),
			},
		},
	}

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	}))
	defer ts.Close()
	client, err := NewClientWithArgs(ts.URL, "", math.MaxInt64, true, false)
	if err != nil {
		t.Fatal(err)
	}

	sp := StoragePool{
		StoragePool: storagePool,
		client:      client,
	}
	err = sp.ModifyRMCache(poolID)
	assert.Nil(t, err)
}

func TestSetRebuildEnabled(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	}))
	defer ts.Close()

	pd, err := SetUpProtectionDomain(ts.URL)
	if err != nil {
		t.Fatal(err)
	}

	err = pd.SetRebuildEnabled("", "")
	assert.Nil(t, err)
}

func TestSetRebuildRebalanceParallelismParam(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	}))
	defer ts.Close()

	pd, err := SetUpProtectionDomain(ts.URL)
	if err != nil {
		t.Fatal(err)
	}

	err = pd.SetRebuildRebalanceParallelismParam("", "")
	assert.Nil(t, err)
}

func TestEnableRFCache(t *testing.T) {
	poolID := "1a2b345c-123a-123a-ab1c-abc1fd234567e"
	type checkFn func(*testing.T, string, error)
	check := func(fns ...checkFn) []checkFn { return fns }

	hasNoError := func(t *testing.T, _ string, err error) {
		if err != nil {
			t.Fatalf("expected no error")
		}
	}
	hasError := func(t *testing.T, _ string, err error) {
		if err == nil {
			t.Fatalf("expected error")
		}
	}

	checkResp := func(storagePoolID string) func(t *testing.T, resp string, err error) {
		return func(t *testing.T, resp string, _ error) {
			assert.Equal(t, storagePoolID, resp)
		}
	}

	tests := map[string]func(*testing.T) (*httptest.Server, []checkFn){
		"success": func(t *testing.T) (*httptest.Server, []checkFn) {
			href := fmt.Sprintf("/api/instances/StoragePool::%v/action/enableRfcache", poolID)

			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method != http.MethodPost {
					t.Fatal(fmt.Errorf("wrong method. Expected %s; but got %s", http.MethodGet, r.Method))
				}

				if r.URL.Path != href {
					t.Fatal(fmt.Errorf("wrong path. Expected %s; but got %s", href, r.URL.Path))
				}

				resp := types.StoragePoolResp{
					ID: "1a2b345c-123a-123a-ab1c-abc1fd234567e",
				}
				respData, err := json.Marshal(resp)
				if err != nil {
					t.Fatal(err)
				}
				_, err = fmt.Fprintln(w, string(respData))
				if err != nil {
					return
				}
			}))
			return ts, check(hasNoError, checkResp("1a2b345c-123a-123a-ab1c-abc1fd234567e"))
		},
		"bad request": func(*testing.T) (*httptest.Server, []checkFn) {
			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				w.WriteHeader(http.StatusBadRequest)
			}))
			return ts, check(hasError)
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			ts, checkFns := tc(t)
			defer ts.Close()

			client, err := NewClientWithArgs(ts.URL, "", math.MaxInt64, true, false)
			if err != nil {
				t.Fatal(err)
			}

			pd := ProtectionDomain{
				client: client,
			}

			resp, err := pd.EnableRFCache(poolID)
			for _, checkFn := range checkFns {
				checkFn(t, resp, err)
			}
		})
	}
}

func TestDisableRFCache(t *testing.T) {
	poolID := "1a2b345c-123a-123a-ab1c-abc1fd234567e"
	type checkFn func(*testing.T, string, error)
	check := func(fns ...checkFn) []checkFn { return fns }

	hasNoError := func(t *testing.T, _ string, err error) {
		if err != nil {
			t.Fatalf("expected no error")
		}
	}
	hasError := func(t *testing.T, _ string, err error) {
		if err == nil {
			t.Fatalf("expected error")
		}
	}

	checkResp := func(storagePoolID string) func(t *testing.T, resp string, err error) {
		return func(t *testing.T, resp string, _ error) {
			assert.Equal(t, storagePoolID, resp)
		}
	}

	tests := map[string]func(t *testing.T) (*httptest.Server, []checkFn){
		"success": func(t *testing.T) (*httptest.Server, []checkFn) {
			href := fmt.Sprintf("/api/instances/StoragePool::%v/action/disableRfcache", poolID)

			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method != http.MethodPost {
					t.Fatal(fmt.Errorf("wrong method. Expected %s; but got %s", http.MethodGet, r.Method))
				}

				if r.URL.Path != href {
					t.Fatal(fmt.Errorf("wrong path. Expected %s; but got %s", href, r.URL.Path))
				}

				resp := types.StoragePoolResp{
					ID: "1a2b345c-123a-123a-ab1c-abc1fd234567e",
				}
				respData, err := json.Marshal(resp)
				if err != nil {
					t.Fatal(err)
				}
				_, err = fmt.Fprintln(w, string(respData))
				if err != nil {
					return
				}
			}))
			return ts, check(hasNoError, checkResp("1a2b345c-123a-123a-ab1c-abc1fd234567e"))
		},
		"bad request": func(*testing.T) (*httptest.Server, []checkFn) {
			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				w.WriteHeader(http.StatusBadRequest)
			}))
			return ts, check(hasError)
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			ts, checkFns := tc(t)
			defer ts.Close()

			client, err := NewClientWithArgs(ts.URL, "", math.MaxInt64, true, false)
			if err != nil {
				t.Fatal(err)
			}

			pd := ProtectionDomain{
				client: client,
			}

			resp, err := pd.DisableRFCache(poolID)
			for _, checkFn := range checkFns {
				checkFn(t, resp, err)
			}
		})
	}
}

func TestFragmentation(t *testing.T) {
	poolID := ""
	tests := []struct {
		name  string
		path  string
		value bool
	}{
		{
			name:  "Enable Fragmentation",
			path:  fmt.Sprintf("/api/instances/StoragePool::%v/action/enableFragmentation", poolID),
			value: true,
		},
		{
			name:  "Disable Fragmentation",
			path:  fmt.Sprintf("/api/instances/StoragePool::%v/action/disableFragmentation", poolID),
			value: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method != http.MethodPost {
					t.Fatal(fmt.Errorf("wrong method. Expected %s; but got %s", http.MethodGet, r.Method))
				}
				if r.URL.Path != tt.path {
					t.Fatal(fmt.Errorf("wrong path. Expected %s; but got %s", tt.path, r.URL.Path))
				}
				w.WriteHeader(http.StatusNoContent)
			}))
			defer ts.Close()

			client, err := NewClientWithArgs(ts.URL, "", math.MaxInt64, true, false)
			if err != nil {
				t.Fatal(err)
			}
			pd := ProtectionDomain{
				client: client,
			}
			err = pd.Fragmentation(poolID, tt.value)
			assert.Nil(t, err)
		})
	}
}

func TestGetStatistics(t *testing.T) {
	poolID := "1a2b345c-123a-123a-ab1c-abc1fd234567e"
	type checkFn func(*testing.T, *types.Statistics, error)
	check := func(fns ...checkFn) []checkFn { return fns }

	hasNoError := func(t *testing.T, _ *types.Statistics, err error) {
		if err != nil {
			t.Fatalf("expected no error")
		}
	}
	hasError := func(t *testing.T, _ *types.Statistics, err error) {
		if err == nil {
			t.Fatalf("expected error")
		}
	}

	checkResp := func(numOfSDC int) func(t *testing.T, resp *types.Statistics, err error) {
		return func(t *testing.T, resp *types.Statistics, _ error) {
			assert.Equal(t, numOfSDC, resp.NumOfSdc)
		}
	}

	tests := map[string]func(t *testing.T) (*httptest.Server, *types.StoragePool, []checkFn){
		"success": func(t *testing.T) (*httptest.Server, *types.StoragePool, []checkFn) {
			storagePool := types.StoragePool{
				Links: []*types.Link{
					{
						Rel:  "/api/StoragePool/relationship/Statistics",
						HREF: fmt.Sprintf("/api/instances/StoragePool::%s", poolID),
					},
				},
			}
			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method != http.MethodGet {
					t.Fatal(fmt.Errorf("wrong method. Expected %s; but got %s", http.MethodGet, r.Method))
				}

				resp := types.Statistics{
					NumOfSdc: 1,
				}
				respData, err := json.Marshal(resp)
				if err != nil {
					t.Fatal(err)
				}
				_, err = fmt.Fprintln(w, string(respData))
				if err != nil {
					return
				}
			}))
			return ts, &storagePool, check(hasNoError, checkResp(1))
		},
		"bad request": func(*testing.T) (*httptest.Server, *types.StoragePool, []checkFn) {
			storagePool := types.StoragePool{}
			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				w.WriteHeader(http.StatusBadRequest)
			}))
			return ts, &storagePool, check(hasError)
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			ts, storagePool, checkFns := tc(t)
			defer ts.Close()

			client, err := NewClientWithArgs(ts.URL, "", math.MaxInt64, true, false)
			if err != nil {
				t.Fatal(err)
			}

			sp := StoragePool{
				client:      client,
				StoragePool: storagePool,
			}
			resp, err := sp.GetStatistics()
			for _, checkFn := range checkFns {
				checkFn(t, resp, err)
			}
		})
	}
}

func TestGetSDSStoragePool(t *testing.T) {
	sdsID := "1a2b345c-123a-123a-ab1c-abc1fd234567e"
	type checkFn func(*testing.T, []types.Sds, error)
	check := func(fns ...checkFn) []checkFn { return fns }

	hasNoError := func(t *testing.T, _ []types.Sds, err error) {
		if err != nil {
			t.Fatalf("expected no error")
		}
	}
	hasError := func(t *testing.T, _ []types.Sds, err error) {
		if err == nil {
			t.Fatalf("expected error")
		}
	}

	checkResp := func(sdsID string) func(t *testing.T, resp []types.Sds, err error) {
		return func(t *testing.T, resp []types.Sds, _ error) {
			assert.Equal(t, sdsID, resp[0].ID)
		}
	}

	tests := map[string]func(t *testing.T) (*httptest.Server, *types.StoragePool, []checkFn){
		"success": func(t *testing.T) (*httptest.Server, *types.StoragePool, []checkFn) {
			storagePool := types.StoragePool{
				Links: []*types.Link{
					{
						Rel:  "/api/StoragePool/relationship/SpSds",
						HREF: fmt.Sprintf("/api/instances/StoragePool::%s", sdsID),
					},
				},
			}
			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method != http.MethodGet {
					t.Fatal(fmt.Errorf("wrong method. Expected %s; but got %s", http.MethodGet, r.Method))
				}

				resp := []types.Sds{
					{
						ID: sdsID,
					},
				}
				respData, err := json.Marshal(resp)
				if err != nil {
					t.Fatal(err)
				}
				_, err = fmt.Fprintln(w, string(respData))
				if err != nil {
					return
				}
			}))
			return ts, &storagePool, check(hasNoError, checkResp(sdsID))
		},
		"bad request": func(*testing.T) (*httptest.Server, *types.StoragePool, []checkFn) {
			storagePool := types.StoragePool{}
			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				w.WriteHeader(http.StatusBadRequest)
			}))
			return ts, &storagePool, check(hasError)
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			ts, storagePool, checkFns := tc(t)
			defer ts.Close()

			client, err := NewClientWithArgs(ts.URL, "", math.MaxInt64, true, false)
			if err != nil {
				t.Fatal(err)
			}

			sp := StoragePool{
				client:      client,
				StoragePool: storagePool,
			}
			resp, err := sp.GetSDSStoragePool()
			for _, checkFn := range checkFns {
				checkFn(t, resp, err)
			}
		})
	}
}
