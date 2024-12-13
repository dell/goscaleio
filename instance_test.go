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
	"context"
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"net/http/httptest"
	"testing"

	types "github.com/dell/goscaleio/types/v1"
)

func mockInstanceServerHandler(resp http.ResponseWriter, req *http.Request) {
	switch req.RequestURI {
	case "/api/types/System/instances":
		if req.Method == http.MethodGet {
			resp.WriteHeader(http.StatusOK)
			response := []types.System{
				{
					ID:   "mock-system-id",
					Name: "mock-system-name",
				},
			}
			content, err := json.Marshal(response)
			if err != nil {
				http.Error(resp, err.Error(), http.StatusNotFound)
			}
			resp.Write(content)
		}
	case "/api/types/Volume/instances/action/queryIdByKey":
		resp.WriteHeader(http.StatusOK)
		response := "mock-volume-id"

		content, err := json.Marshal(response)
		if err != nil {
			http.Error(resp, err.Error(), http.StatusNotFound)
		}

		resp.Write(content)
	case fmt.Sprintf("/api/instances/Volume::%s", "mock-volume-id"):
		resp.WriteHeader(http.StatusOK)
		response := types.Volume{
			ID:               "mock-volume-id",
			Name:             "mock-volume-name",
			AncestorVolumeID: "mock-ancestore-volume-id",
		}
		content, err := json.Marshal(response)
		if err != nil {
			http.Error(resp, err.Error(), http.StatusNotFound)
		}

		resp.Write(content)
	case "/api/types/Volume/instances":
		if req.Method == http.MethodGet {
			resp.WriteHeader(http.StatusOK)
			response := []types.Volume{
				{ID: "mock-volume-id"},
			}
			content, err := json.Marshal(response)
			if err != nil {
				http.Error(resp, err.Error(), http.StatusNotFound)
			}
			resp.Write(content)
		} else if req.Method == http.MethodPost {
			resp.WriteHeader(http.StatusOK)
			response := types.VolumeResp{
				ID: "mock-volume-id",
			}
			content, err := json.Marshal(response)
			if err != nil {
				http.Error(resp, err.Error(), http.StatusNotFound)
			}
			resp.Write(content)
		}
	case "/api/types/StoragePool/instances":
		resp.WriteHeader(http.StatusOK)
		response := []types.StoragePool{
			{
				ID:                 "mock-storage-pool-id",
				Name:               "mock-storage-pool-name",
				ProtectionDomainID: "mock-protection-domain-id",
			},
		}
		content, err := json.Marshal(response)
		if err != nil {
			http.Error(resp, err.Error(), http.StatusNotFound)
		}
		resp.Write(content)
	case fmt.Sprintf("/api/instances/StoragePool::%s", "mock-storage-pool-id"):
		resp.WriteHeader(http.StatusOK)
		response := types.StoragePool{
			ID:                 "mock-storage-pool-id",
			Name:               "mock-storage-pool-name",
			ProtectionDomainID: "mock-protection-domain-id",
		}
		content, err := json.Marshal(response)
		if err != nil {
			http.Error(resp, err.Error(), http.StatusNotFound)
		}
		resp.Write(content)
	case fmt.Sprintf("/api/instances/StoragePool::%s/relationships/Volume", "mock-storage-pool-id"):
		resp.WriteHeader(http.StatusOK)
		response := []types.Volume{
			{
				Name:          "mock-volume-name-1",
				ID:            "mock-volume-id-1",
				StoragePoolID: "mock-storage-pool-id",
				VTreeID:       "mock-vtree-id",
			},
			{
				Name:          "mock-volume-name-2",
				ID:            "mock-volume-id-2",
				StoragePoolID: "mock-storage-pool-id-2",
				VTreeID:       "mock-vtree-id",
			},
		}
		content, err := json.Marshal(response)
		if err != nil {
			http.Error(resp, err.Error(), http.StatusNotFound)
		}
		resp.Write(content)
	case "/api/types/SnapshotPolicy/instances/action/queryIdByKey":
		resp.WriteHeader(http.StatusOK)
		response := "mock-snapshot-policy-id"
		content, err := json.Marshal(response)
		if err != nil {
			http.Error(resp, err.Error(), http.StatusNotFound)
		}
		resp.Write(content)
	case fmt.Sprintf("/api/instances/SnapshotPolicy::%s", "mock-snapshot-policy-id"):
		resp.WriteHeader(http.StatusOK)
		response := types.SnapshotPolicy{
			ID:   "mock-snapshot-policy-id",
			Name: "mock-snapshot-policy-name",
		}
		content, err := json.Marshal(response)
		if err != nil {
			http.Error(resp, err.Error(), http.StatusNotFound)
		}
		resp.Write(content)
	case "/api/types/SnapshotPolicy/instances":
		resp.WriteHeader(http.StatusOK)
		response := []types.SnapshotPolicy{
			{
				ID:   "mock-snapshot-policy-id",
				Name: "mock-snapshot-policy-name",
			},
		}
		content, err := json.Marshal(response)
		if err != nil {
			http.Error(resp, err.Error(), http.StatusNotFound)
		}
		resp.Write(content)

	default:
		resp.WriteHeader(http.StatusNoContent)
	}
}

func TestGetInstance(t *testing.T) {
	mockServer := httptest.NewServer(http.HandlerFunc(mockInstanceServerHandler))
	defer mockServer.Close()
	tests := []struct {
		name       string
		systemhref string
		error      string
	}{
		{
			name:       "system href valid",
			systemhref: "/api/instances/System::mock-system-id",
			error:      "",
		},
		{
			name:       "system href null",
			systemhref: "",
			error:      "",
		},
	}

	for _, tc := range tests {
		client, err := NewClientWithArgs(mockServer.URL, "3.6", math.MaxInt64, true, false)
		if err != nil {
			t.Fatal(err)
		}
		_, err = client.GetInstance(context.Background(), tc.systemhref)
		if err != nil {
			if tc.error != err.Error() {
				t.Fatal(err)
			}
		}
	}
}
func TestGetVolume(t *testing.T) {
	mockServer := httptest.NewServer(http.HandlerFunc(mockInstanceServerHandler))
	defer mockServer.Close()

	tests := []struct {
		name              string
		volumeid          string
		volumename        string
		volumehref        string
		ancestorevolumeid string
		snapshots         bool
		error             string
	}{
		{
			name:              "get volume name not null",
			volumeid:          "",
			volumehref:        "",
			ancestorevolumeid: "mock-ancestor-volume-id",
			volumename:        "mock-volume-name",
			snapshots:         false,
			error:             "volume not found",
		},
		{
			name:              "get volume id not null volume name and href null",
			volumeid:          "mock-volume-id",
			volumehref:        "",
			ancestorevolumeid: "mock-ancestor-volume-id",
			volumename:        "",
			snapshots:         false,
			error:             "",
		},
		{
			name:              "get volume volume id href and volume name null",
			volumeid:          "",
			volumehref:        "",
			ancestorevolumeid: "mock-ancestor-volume-id",
			volumename:        "",
			snapshots:         false,
			error:             "",
		},
		{
			name:              "get volume same id and ancestor volume id",
			volumeid:          "",
			volumehref:        "/api/types/Volume/instances",
			ancestorevolumeid: "mock-volume-id",
			volumename:        "mock-volume-name",
			snapshots:         false,
			error:             "volume not found",
		},
		{
			name:              "get volume same id and ancestor volume id",
			volumeid:          "",
			volumehref:        "",
			ancestorevolumeid: "mock-volume-id",
			volumename:        "mock-volume-name",
			snapshots:         true,
			error:             "volume not found",
		},
	}

	for _, tc := range tests {
		client, err := NewClientWithArgs(mockServer.URL, "3.6", math.MaxInt64, true, false)
		if err != nil {
			t.Fatal(err)
		}

		_, err = client.GetVolume(context.Background(), tc.volumehref, tc.volumeid, tc.ancestorevolumeid, tc.volumename, tc.snapshots)
		if err != nil {
			if tc.error != err.Error() {
				t.Fatal(err)
			}
		}
	}
}

func TestGetStoragePool(t *testing.T) {
	mockServer := httptest.NewServer(http.HandlerFunc(mockInstanceServerHandler))
	defer mockServer.Close()

	tests := []struct {
		name            string
		storagepoolhref string
		error           string
	}{
		{
			name:            "storage pool valid",
			storagepoolhref: "/api/instances/StoragePool::mock-storage-pool-id",
			error:           "storage pool not found",
		},
		{
			name:            "get volume valid and href null",
			storagepoolhref: "",
			error:           "",
		},
	}

	for _, tc := range tests {
		client, err := NewClientWithArgs(mockServer.URL, "3.6", math.MaxInt64, true, false)
		if err != nil {
			t.Fatal(err)
		}

		_, err = client.GetStoragePool(context.Background(), tc.storagepoolhref)
		if err != nil {
			if tc.error != err.Error() {
				t.Fatal(err)
			}
		}
	}
}

func TestFindStoragePool(t *testing.T) {
	mockServer := httptest.NewServer(http.HandlerFunc(mockInstanceServerHandler))
	defer mockServer.Close()

	tests := []struct {
		name               string
		poolid             string
		poolname           string
		storagepoolhref    string
		protectiondomainid string
		error              string
	}{
		{
			name:               "storage pool valid",
			poolid:             "mock-storage-pool-id",
			poolname:           "mock-storage-pool-name",
			protectiondomainid: "mock-protection-domain-id",
			storagepoolhref:    "",
			error:              "",
		},
		{

			name:               "storage pool valid href valid",
			poolid:             "mock-storage-pool-id",
			poolname:           "mock-storage-pool-name",
			protectiondomainid: "mock-protection-domain-id",
			storagepoolhref:    "/api/instances/StoragePool::mock-storage-pool-id",
			error:              "",
		},
		{

			name:               "storage pool invalid",
			poolid:             "mock-storage-pool-invalid-id",
			poolname:           "mock-storage-pool-invalid-name",
			protectiondomainid: "mock-protection-domain",
			storagepoolhref:    "",
			error:              "Couldn't find storage pool",
		},
	}

	for _, tc := range tests {
		client, err := NewClientWithArgs(mockServer.URL, "3.6", math.MaxInt64, true, false)
		if err != nil {
			t.Fatal(err)
		}

		_, err = client.FindStoragePool(context.Background(), tc.poolid, tc.poolname, tc.storagepoolhref, tc.protectiondomainid)
		if err != nil {
			if tc.error != err.Error() {
				t.Fatal(err)
			}
		}
	}
}

func TestGetStoragePoolVolumes(t *testing.T) {
	mockServer := httptest.NewServer(http.HandlerFunc(mockInstanceServerHandler))
	defer mockServer.Close()

	tests := []struct {
		name          string
		storagepoolid string
		error         string
	}{
		{
			name:          "storage pool valid",
			storagepoolid: "mock-storage-pool-id",
			error:         "",
		},
		{
			name:          "storage pool invalid",
			storagepoolid: "mock-storage-pool-id-invalid",
			error:         "storage pool not found",
		},
	}

	for _, tc := range tests {
		client, err := NewClientWithArgs(mockServer.URL, "3.6", math.MaxInt64, true, false)
		if err != nil {
			t.Fatal(err)
		}

		_, err = client.GetStoragePoolVolumes(context.Background(), tc.storagepoolid)
		if err != nil {
			if tc.error != err.Error() {
				t.Fatal(err)
			}
		}
	}
}

func TestCreateVolume(t *testing.T) {
	mockServer := httptest.NewServer(http.HandlerFunc(mockInstanceServerHandler))
	defer mockServer.Close()

	tests := []struct {
		name             string
		volumeparam      *types.VolumeParam
		protectiondomain string
		storagepoolname  string
		error            string
	}{
		{
			name: "valid volume create",
			volumeparam: &types.VolumeParam{
				Name:               "mock-volume-name",
				VolumeSizeInKb:     "1024",
				StoragePoolID:      "mock-storage-pool-id",
				ProtectionDomainID: "mock-protection-domain-id",
			},
			storagepoolname:  "mock-storage-pool-name",
			protectiondomain: "mock-protection-domain-id",
			error:            "",
		},
	}

	for _, tc := range tests {
		client, err := NewClientWithArgs(mockServer.URL, "3.6", math.MaxInt64, true, false)
		if err != nil {
			t.Fatal(err)
		}

		_, err = client.CreateVolume(tc.volumeparam, tc.storagepoolname, tc.protectiondomain)
		if err != nil {
			if tc.error != err.Error() {
				t.Fatal(err)
			}
		}
	}
}

func TestGetSnapshotPolicyI(t *testing.T) {
	mockServer := httptest.NewServer(http.HandlerFunc(mockInstanceServerHandler))
	defer mockServer.Close()

	tests := []struct {
		name               string
		snapshotpolicyid   string
		snapshotpolicyname string
		error              string
	}{
		{
			name:               "snapshot policy name empty",
			snapshotpolicyid:   "mock-snapshot-policy-id",
			snapshotpolicyname: "",
			error:              "",
		},
		{
			name:               "snapshot policy d empty",
			snapshotpolicyid:   "",
			snapshotpolicyname: "mock-snapshot-policy-name",
			error:              "",
		},
	}

	for _, tc := range tests {
		client, err := NewClientWithArgs(mockServer.URL, "3.6", math.MaxInt64, true, false)
		if err != nil {
			t.Fatal(err)
		}

		_, err = client.GetSnapshotPolicy(context.Background(), tc.snapshotpolicyname, tc.snapshotpolicyid)
		if err != nil {
			if tc.error != err.Error() {
				t.Fatal(err)
			}
		}
	}
}

func TestNewSnapshotPolicyI(t *testing.T) {
	mockServer := httptest.NewServer(http.HandlerFunc(mockInstanceServerHandler))
	defer mockServer.Close()

	client, err := NewClientWithArgs(mockServer.URL, "3.6", math.MaxInt64, true, false)
	if err != nil {
		t.Fatal(err)
	}
	_ = NewSnapshotPolicy(client)
}
