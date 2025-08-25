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

// Copyright © 2020 - 2022 Dell Inc. or its subsidiaries. All Rights Reserved.
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
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestGetVolumeStatistics(t *testing.T) {
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

	tests := map[string]func(t *testing.T) (*httptest.Server, *types.Volume, []checkFn){
		"success": func(t *testing.T) (*httptest.Server, *types.Volume, []checkFn) {
			VolumeID := "000001111a2222b"
			href := fmt.Sprintf("/api/instances/Volume::%s/relationships/Statistics", VolumeID)
			vol := types.Volume{
				ID: VolumeID,
				Links: []*types.Link{
					{
						Rel:  "/api/Volume/relationship/Statistics",
						HREF: fmt.Sprintf("/api/instances/Volume::%s/relationships/Statistics", VolumeID),
					},
				},
			}
			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method != http.MethodGet {
					t.Fatal(fmt.Errorf("wrong method. Expected %s; but got %s", http.MethodGet, r.Method))
				}
				if r.URL.Path != href {
					t.Fatal(fmt.Errorf("wrong path. Expected %s; but got %s", href, r.URL.Path))
				}
				volumeStats := types.VolumeStatistics{}
				respData, err := json.Marshal(volumeStats)
				if err != nil {
					t.Fatal(err)
				}
				fmt.Fprintln(w, string(respData))
			}))
			return ts, &vol, check(hasNoError)
		},
		"error from getJSONWithRetry": func(t *testing.T) (*httptest.Server, *types.Volume, []checkFn) {
			VolumeID := "000001111a2222b"
			href := fmt.Sprintf("/api/instances/Volume::%s/relationships/Statistics", VolumeID)
			vol := types.Volume{
				ID: VolumeID,
				Links: []*types.Link{
					{
						Rel:  "/api/Volume/relationship/Statistics",
						HREF: fmt.Sprintf("/api/instances/Volume::%s/relationships/Statistics", VolumeID),
					},
				},
			}
			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method != http.MethodGet {
					t.Fatal(fmt.Errorf("wrong method. Expected %s; but got %s", http.MethodGet, r.Method))
				}
				if r.URL.Path != href {
					t.Fatal(fmt.Errorf("wrong path. Expected %s; but got %s", href, r.URL.Path))
				}
				http.NotFound(w, r)
			}))
			return ts, &vol, check(hasError)
		},
		"error from GetLink": func(t *testing.T) (*httptest.Server, *types.Volume, []checkFn) {
			VolumeID := "000001111a2222b"
			vol := types.Volume{
				ID: VolumeID,
				Links: []*types.Link{
					{
						Rel:  "noLink error",
						HREF: fmt.Sprintf("/api/instances/Volume::%s/relationships/Statistics", VolumeID),
					},
				},
			}
			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				volumeStats := types.VolumeStatistics{}
				respData, err := json.Marshal(volumeStats)
				if err != nil {
					t.Fatal(err)
				}
				fmt.Fprintln(w, string(respData))
			}))
			return ts, &vol, check(hasError)
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			ts, vol, checkFns := tc(t)
			defer ts.Close()

			client, err := NewClientWithArgs(ts.URL, "", math.MaxInt64, true, false)
			if err != nil {
				t.Fatal(err)
			}

			volClient := NewVolume(client)
			volClient.Volume = vol
			_, err = volClient.GetVolumeStatistics()
			for _, checkFn := range checkFns {
				checkFn(t, err)
			}
		})
	}
}

func TestGetVolumeSP(t *testing.T) {
	searchVolumeID := uuid.NewString()
	storagePoolID := uuid.NewString()
	type testCase struct {
		storagePool  types.StoragePool
		volumeHref   string
		volumeName   string
		getSnapshots bool
		server       *httptest.Server
		expectedErr  error
	}

	cases := map[string]testCase{
		"success: via volume name": {
			volumeHref:   "",
			volumeName:   "myVolumeName",
			getSnapshots: false,
			storagePool: types.StoragePool{
				Links: []*types.Link{
					{
						Rel:  "/api/StoragePool/relationship/Volume",
						HREF: fmt.Sprintf("/api/instances/StoragePool::%s/relationships/Volume", storagePoolID),
					},
				},
			},
			server: httptest.NewServer(http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
				switch req.RequestURI {
				case "/api/types/Volume/instances/action/queryIdByKey":
					resp.WriteHeader(http.StatusOK)
					resp.Write([]byte(searchVolumeID))
				case fmt.Sprintf("/api/instances/Volume::%s", searchVolumeID):
					resp.WriteHeader(http.StatusOK)
					content, err := json.Marshal(types.Volume{
						Name: "myVolume",
						ID:   searchVolumeID,
					})
					if err != nil {
						t.Fatalf("failed to marshal volume: %v", err)
					}
					resp.Write(content)
				default:
					resp.WriteHeader(http.StatusBadRequest)
					resp.Write([]byte(`{"message":"no route handled","httpStatusCode":400,"errorCode":0}`))
				}
			})),
			expectedErr: nil,
		},
		"success: via empty volumehref": {
			volumeHref:   "",
			volumeName:   "",
			getSnapshots: false,
			storagePool: types.StoragePool{
				Links: []*types.Link{
					{
						Rel:  "/api/StoragePool/relationship/Volume",
						HREF: fmt.Sprintf("/api/instances/StoragePool::%s/relationships/Volume", storagePoolID),
					},
				},
			},
			server: httptest.NewServer(http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
				switch req.RequestURI {
				case "/api/types/Volume/instances/action/queryIdByKey":
					resp.WriteHeader(http.StatusOK)
					resp.Write([]byte(searchVolumeID))
				case fmt.Sprintf("/api/instances/StoragePool::%s/relationships/Volume", storagePoolID):
					resp.WriteHeader(http.StatusOK)
					content, err := json.Marshal([]types.Volume{{
						Name: "myVolume",
						ID:   searchVolumeID,
					}})
					if err != nil {
						t.Fatalf("failed to marshal volume: %v", err)
					}
					resp.Write(content)
				default:
					resp.WriteHeader(http.StatusBadRequest)
					resp.Write([]byte(`{"message":"no route handled","httpStatusCode":400,"errorCode":0}`))
				}
			})),
			expectedErr: nil,
		},
		"success: via volumehref": {
			volumeHref:   fmt.Sprintf("/api/instances/StoragePool::%s/relationships/Volume", storagePoolID),
			volumeName:   "",
			getSnapshots: false,
			storagePool:  types.StoragePool{},
			server: httptest.NewServer(http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
				switch req.RequestURI {
				case "/api/types/Volume/instances/action/queryIdByKey":
					resp.WriteHeader(http.StatusOK)
					resp.Write([]byte(searchVolumeID))
				case fmt.Sprintf("/api/instances/StoragePool::%s/relationships/Volume", storagePoolID):
					resp.WriteHeader(http.StatusOK)
					content, err := json.Marshal(types.Volume{
						Name: "myVolume",
						ID:   searchVolumeID,
					})
					if err != nil {
						t.Fatalf("failed to marshal volume: %v", err)
					}
					resp.Write(content)
				default:
					resp.WriteHeader(http.StatusBadRequest)
					resp.Write([]byte(`{"message":"no route handled","httpStatusCode":400,"errorCode":0}`))
				}
			})),
			expectedErr: nil,
		},
		"success: unable to find ID": {
			volumeHref:   "",
			volumeName:   "myVolumeName",
			getSnapshots: false,
			storagePool: types.StoragePool{
				Links: []*types.Link{
					{
						Rel:  "/api/StoragePool/relationship/Volume",
						HREF: fmt.Sprintf("/api/instances/StoragePool::%s/relationships/Volume", storagePoolID),
					},
				},
			},
			server: httptest.NewServer(http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
				switch req.RequestURI {
				case "/api/types/Volume/instances/action/queryIdByKey":
					resp.WriteHeader(http.StatusNotFound)
					resp.Write([]byte(`{"message":"Not found","httpStatusCode":404,"errorCode":0}`))
				default:
					resp.WriteHeader(http.StatusBadRequest)
					resp.Write([]byte(`{"message":"no route handled","httpStatusCode":400,"errorCode":0}`))
				}
			})),
			expectedErr: nil,
		},
		"error: bad request": {
			volumeHref:   "",
			volumeName:   "myVolumeName",
			getSnapshots: false,
			storagePool: types.StoragePool{
				Links: []*types.Link{
					{
						Rel:  "/api/StoragePool/relationship/Volume",
						HREF: fmt.Sprintf("/api/instances/StoragePool::%s/relationships/Volume", storagePoolID),
					},
				},
			},
			server: httptest.NewServer(http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
				switch req.RequestURI {
				case "/api/types/Volume/instances/action/queryIdByKey":
					resp.WriteHeader(http.StatusBadRequest)
					resp.Write([]byte(`{"message":"bad request","httpStatusCode":400,"errorCode":0}`))
				default:
					resp.WriteHeader(http.StatusBadRequest)
					resp.Write([]byte(`{"message":"no route handled","httpStatusCode":400,"errorCode":0}`))
				}
			})),
			expectedErr: fmt.Errorf("Error: problem finding volume: bad request"),
		},
		"error: via volumehref": {
			volumeHref:   "",
			volumeName:   "",
			getSnapshots: false,
			storagePool:  types.StoragePool{},
			server: httptest.NewServer(http.HandlerFunc(func(resp http.ResponseWriter, _ *http.Request) {
				resp.WriteHeader(http.StatusBadRequest)
				resp.Write([]byte(`{"message":"no route handled","httpStatusCode":400,"errorCode":0}`))
			})),
			expectedErr: fmt.Errorf("Error: problem finding link"),
		},
		"error: getting volume instance": {
			volumeHref:   "",
			volumeName:   "myVolumeName",
			getSnapshots: false,
			storagePool:  types.StoragePool{},
			server: httptest.NewServer(http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
				switch req.RequestURI {
				case "/api/types/Volume/instances/action/queryIdByKey":
					resp.WriteHeader(http.StatusOK)
					resp.Write([]byte(searchVolumeID))
				case fmt.Sprintf("/api/instances/Volume::%s", searchVolumeID):
					resp.WriteHeader(http.StatusBadRequest)
					resp.Write([]byte(`{"message":"bad request","httpStatusCode":400,"errorCode":0}`))
				default:
					resp.WriteHeader(http.StatusBadRequest)
					resp.Write([]byte(`{"message":"no route handled","httpStatusCode":400,"errorCode":0}`))
				}
			})),
			expectedErr: fmt.Errorf("bad request"),
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			defer tc.server.Close()

			client, err := NewClientWithArgs(tc.server.URL, "3.6", math.MaxInt64, true, false)
			if err != nil {
				t.Fatal(err)
			}

			sp := NewStoragePool(client)
			sp.StoragePool = &tc.storagePool

			_, err = sp.GetVolume(tc.volumeHref, "", "", tc.volumeName, tc.getSnapshots)
			if err != nil {
				if tc.expectedErr.Error() != err.Error() {
					t.Fatal(err)
				}
			}
		})
	}
}

func TestGetLocalVolumeMapByRegex(t *testing.T) {
	type testCase struct {
		directoryPrefix string
		systenIDRegex   string
		volumeIDRegex   string
		expectedLength  int
	}

	cases := map[string]testCase{
		"success: mocked location": {
			directoryPrefix: "mocks",
			systenIDRegex:   "",
			volumeIDRegex:   "",
			expectedLength:  1,
		},
		"success: empty response": {
			directoryPrefix: "mocks",
			systenIDRegex:   "mySystemID*",
			volumeIDRegex:   "myVolumeID*",
			expectedLength:  0,
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			FSDevDirectoryPrefix = tc.directoryPrefix

			mappedVolumes, err := GetLocalVolumeMapByRegex(tc.systenIDRegex, tc.volumeIDRegex)
			assert.Nil(t, err)
			assert.Equal(t, tc.expectedLength, len(mappedVolumes))
		})
	}
}

func TestGetLocalVolumeMap(t *testing.T) {
	FSDevDirectoryPrefix = "mocks"

	mappedVolumes, err := GetLocalVolumeMap()
	assert.Nil(t, err)
	assert.Equal(t, 1, len(mappedVolumes))
}

func TestCreateVolumeSP(t *testing.T) {
	type testCase struct {
		storagePool types.StoragePool
		server      *httptest.Server
		expectedErr error
	}

	cases := map[string]testCase{
		"success": {
			storagePool: types.StoragePool{
				ID:                 uuid.NewString(),
				ProtectionDomainID: uuid.NewString(),
			},
			server: httptest.NewServer(http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
				switch req.RequestURI {
				case "/api/types/Volume/instances":
					resp.WriteHeader(http.StatusOK)
					content, err := json.Marshal(types.VolumeResp{
						ID: uuid.NewString(),
					})
					if err != nil {
						t.Fatalf("failed to marshal volume: %v", err)
					}
					resp.Write(content)
				default:
					resp.WriteHeader(http.StatusBadRequest)
					resp.Write([]byte(`{"message":"no route handled","httpStatusCode":400,"errorCode":0}`))
				}
			})),
			expectedErr: nil,
		},
		"error: bad request": {
			storagePool: types.StoragePool{
				ID:                 uuid.NewString(),
				ProtectionDomainID: uuid.NewString(),
			},
			server: httptest.NewServer(http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
				switch req.RequestURI {
				case "/api/types/Volume/instances":
					resp.WriteHeader(http.StatusBadRequest)
					resp.Write([]byte(`{"message":"bad request","httpStatusCode":400,"errorCode":0}`))
				default:
					resp.WriteHeader(http.StatusBadRequest)
					resp.Write([]byte(`{"message":"no route handled","httpStatusCode":400,"errorCode":0}`))
				}
			})),
			expectedErr: fmt.Errorf("bad request"),
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			defer tc.server.Close()

			client, err := NewClientWithArgs(tc.server.URL, "3.6", math.MaxInt64, true, false)
			if err != nil {
				t.Fatal(err)
			}

			sp := NewStoragePool(client)
			sp.StoragePool = &tc.storagePool

			volume := types.VolumeParam{
				Name: "myVolume",
			}

			_, err = sp.CreateVolume(&volume)
			if err != nil {
				if tc.expectedErr.Error() != err.Error() {
					t.Fatal(err)
				}
			}
		})
	}
}

func TestGetVTree(t *testing.T) {
	vtreeID := uuid.NewString()
	type testCase struct {
		volume      types.Volume
		server      *httptest.Server
		expectedErr error
	}

	cases := map[string]testCase{
		"success": {
			volume: types.Volume{
				ID: uuid.NewString(),
				Links: []*types.Link{
					{
						Rel:  "/api/parent/relationship/vtreeId",
						HREF: fmt.Sprintf("/api/instances/VTree::%s", vtreeID),
					},
				},
			},
			server: httptest.NewServer(http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
				switch req.RequestURI {
				case fmt.Sprintf("/api/instances/VTree::%s", vtreeID):
					resp.WriteHeader(http.StatusOK)
					content, err := json.Marshal(types.VTree{
						ID:   uuid.NewString(),
						Name: "myVTree",
					})
					if err != nil {
						t.Fatalf("failed to marshal volume: %v", err)
					}
					resp.Write(content)
				default:
					resp.WriteHeader(http.StatusBadRequest)
					resp.Write([]byte(`{"message":"no route handled","httpStatusCode":400,"errorCode":0}`))
				}
			})),
			expectedErr: nil,
		},
		"error: problem finding link": {
			volume: types.Volume{
				ID: uuid.NewString(),
			},
			server: httptest.NewServer(http.HandlerFunc(func(resp http.ResponseWriter, _ *http.Request) {
				resp.WriteHeader(http.StatusBadRequest)
				resp.Write([]byte(`{"message":"no route handled","httpStatusCode":400,"errorCode":0}`))
			})),
			expectedErr: fmt.Errorf("Error: problem finding link"),
		},
		"error: bad request": {
			volume: types.Volume{
				ID: uuid.NewString(),
				Links: []*types.Link{
					{
						Rel:  "/api/parent/relationship/vtreeId",
						HREF: fmt.Sprintf("/api/instances/VTree::%s", vtreeID),
					},
				},
			},
			server: httptest.NewServer(http.HandlerFunc(func(resp http.ResponseWriter, _ *http.Request) {
				resp.WriteHeader(http.StatusBadRequest)
				resp.Write([]byte(`{"message":"bad request","httpStatusCode":400,"errorCode":0}`))
			})),
			expectedErr: fmt.Errorf("bad request"),
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			defer tc.server.Close()

			client, err := NewClientWithArgs(tc.server.URL, "3.6", math.MaxInt64, true, false)
			if err != nil {
				t.Fatal(err)
			}

			vl := NewVolume(client)
			vl.Volume = &tc.volume

			_, err = vl.GetVTree()
			if err != nil {
				if tc.expectedErr.Error() != err.Error() {
					t.Fatal(err)
				}
			}
		})
	}
}

func TestRemoveVolume(t *testing.T) {
	selfID := uuid.NewString()
	type testCase struct {
		volume      types.Volume
		server      *httptest.Server
		expectedErr error
	}

	cases := map[string]testCase{
		"success": {
			volume: types.Volume{
				ID: uuid.NewString(),
				Links: []*types.Link{
					{
						Rel:  "self",
						HREF: fmt.Sprintf("/api/instances/Volume::%s", selfID),
					},
				},
			},
			server: httptest.NewServer(http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
				switch req.RequestURI {
				case fmt.Sprintf("/api/instances/Volume::%s/action/removeVolume", selfID):
					resp.WriteHeader(http.StatusNoContent)
				default:
					resp.WriteHeader(http.StatusBadRequest)
					resp.Write([]byte(`{"message":"no route handled","httpStatusCode":400,"errorCode":0}`))
				}
			})),
			expectedErr: nil,
		},
		"error: problem finding link": {
			volume: types.Volume{
				ID: uuid.NewString(),
			},
			server: httptest.NewServer(http.HandlerFunc(func(resp http.ResponseWriter, _ *http.Request) {
				resp.WriteHeader(http.StatusBadRequest)
				resp.Write([]byte(`{"message":"no route handled","httpStatusCode":400,"errorCode":0}`))
			})),
			expectedErr: fmt.Errorf("Error: problem finding link"),
		},
		"error: bad request": {
			volume: types.Volume{
				ID: uuid.NewString(),
				Links: []*types.Link{
					{
						Rel:  "self",
						HREF: fmt.Sprintf("/api/instances/Volume::%s", selfID),
					},
				},
			},
			server: httptest.NewServer(http.HandlerFunc(func(resp http.ResponseWriter, _ *http.Request) {
				resp.WriteHeader(http.StatusBadRequest)
				resp.Write([]byte(`{"message":"bad request","httpStatusCode":400,"errorCode":0}`))
			})),
			expectedErr: fmt.Errorf("bad request"),
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			defer tc.server.Close()

			client, err := NewClientWithArgs(tc.server.URL, "3.6", math.MaxInt64, true, false)
			if err != nil {
				t.Fatal(err)
			}

			vl := NewVolume(client)
			vl.Volume = &tc.volume

			err = vl.RemoveVolume("")
			if err != nil {
				if tc.expectedErr.Error() != err.Error() {
					t.Fatal(err)
				}
			}
		})
	}
}

func TestSetVolumeName(t *testing.T) {
	volumeID := uuid.NewString()
	type testCase struct {
		volume      types.Volume
		server      *httptest.Server
		expectedErr error
	}

	cases := map[string]testCase{
		"success": {
			volume: types.Volume{
				ID: volumeID,
			},
			server: httptest.NewServer(http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
				switch req.RequestURI {
				case fmt.Sprintf("/api/instances/Volume::%s/action/setVolumeName", volumeID):
					resp.WriteHeader(http.StatusOK)
				default:
					resp.WriteHeader(http.StatusBadRequest)
					resp.Write([]byte(`{"message":"no route handled","httpStatusCode":400,"errorCode":0}`))
				}
			})),
			expectedErr: nil,
		},
		"error: bad request": {
			volume: types.Volume{
				ID: volumeID,
			},
			server: httptest.NewServer(http.HandlerFunc(func(resp http.ResponseWriter, _ *http.Request) {
				resp.WriteHeader(http.StatusBadRequest)
				resp.Write([]byte(`{"message":"bad request","httpStatusCode":400,"errorCode":0}`))
			})),
			expectedErr: fmt.Errorf("bad request"),
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			defer tc.server.Close()

			client, err := NewClientWithArgs(tc.server.URL, "3.6", math.MaxInt64, true, false)
			if err != nil {
				t.Fatal(err)
			}

			vl := NewVolume(client)
			vl.Volume = &tc.volume

			err = vl.SetVolumeName("newVolumeName")
			if err != nil {
				if tc.expectedErr.Error() != err.Error() {
					t.Fatal(err)
				}
			}
		})
	}
}

func TestSetVolumeSize(t *testing.T) {
	volumeID := uuid.NewString()
	type testCase struct {
		volume      types.Volume
		server      *httptest.Server
		expectedErr error
	}

	cases := map[string]testCase{
		"success": {
			volume: types.Volume{
				ID: volumeID,
				Links: []*types.Link{
					{
						Rel:  "self",
						HREF: fmt.Sprintf("/api/instances/Volume::%s", volumeID),
					},
				},
			},
			server: httptest.NewServer(http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
				switch req.RequestURI {
				case fmt.Sprintf("/api/instances/Volume::%s/action/setVolumeSize", volumeID):
					resp.WriteHeader(http.StatusOK)
				default:
					resp.WriteHeader(http.StatusBadRequest)
					resp.Write([]byte(`{"message":"no route handled","httpStatusCode":400,"errorCode":0}`))
				}
			})),
			expectedErr: nil,
		},
		"error: bad request": {
			volume: types.Volume{
				ID: volumeID,
				Links: []*types.Link{
					{
						Rel:  "self",
						HREF: fmt.Sprintf("/api/instances/Volume::%s", volumeID),
					},
				},
			},
			server: httptest.NewServer(http.HandlerFunc(func(resp http.ResponseWriter, _ *http.Request) {
				resp.WriteHeader(http.StatusBadRequest)
				resp.Write([]byte(`{"message":"bad request","httpStatusCode":400,"errorCode":0}`))
			})),
			expectedErr: fmt.Errorf("bad request"),
		},
		"error: problem finding link": {
			volume: types.Volume{
				ID: uuid.NewString(),
			},
			server: httptest.NewServer(http.HandlerFunc(func(resp http.ResponseWriter, _ *http.Request) {
				resp.WriteHeader(http.StatusBadRequest)
				resp.Write([]byte(`{"message":"no route handled","httpStatusCode":400,"errorCode":0}`))
			})),
			expectedErr: fmt.Errorf("Error: problem finding link"),
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			defer tc.server.Close()

			client, err := NewClientWithArgs(tc.server.URL, "3.6", math.MaxInt64, true, false)
			if err != nil {
				t.Fatal(err)
			}

			vl := NewVolume(client)
			vl.Volume = &tc.volume

			err = vl.SetVolumeSize("32Gi")
			if err != nil {
				if tc.expectedErr.Error() != err.Error() {
					t.Fatal(err)
				}
			}
		})
	}
}

func TestLockAutoSnapshot(t *testing.T) {
	volumeID := uuid.NewString()
	type testCase struct {
		volume      types.Volume
		server      *httptest.Server
		expectedErr error
	}

	cases := map[string]testCase{
		"success": {
			volume: types.Volume{
				ID:         volumeID,
				VolumeType: "Snapshot",
				Links: []*types.Link{
					{
						Rel:  "self",
						HREF: fmt.Sprintf("/api/instances/Volume::%s", volumeID),
					},
				},
			},
			server: httptest.NewServer(http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
				switch req.RequestURI {
				case fmt.Sprintf("/api/instances/Volume::%s/action/lockAutoSnapshot", volumeID):
					resp.WriteHeader(http.StatusOK)
				default:
					resp.WriteHeader(http.StatusBadRequest)
					resp.Write([]byte(`{"message":"no route handled","httpStatusCode":400,"errorCode":0}`))
				}
			})),
			expectedErr: nil,
		},
		"error: bad request": {
			volume: types.Volume{
				ID:         volumeID,
				VolumeType: "Snapshot",
				Links: []*types.Link{
					{
						Rel:  "self",
						HREF: fmt.Sprintf("/api/instances/Volume::%s", volumeID),
					},
				},
			},
			server: httptest.NewServer(http.HandlerFunc(func(resp http.ResponseWriter, _ *http.Request) {
				resp.WriteHeader(http.StatusBadRequest)
				resp.Write([]byte(`{"message":"bad request","httpStatusCode":400,"errorCode":0}`))
			})),
			expectedErr: fmt.Errorf("bad request"),
		},
		"error: problem finding link": {
			volume: types.Volume{
				ID:         uuid.NewString(),
				VolumeType: "Snapshot",
			},
			server: httptest.NewServer(http.HandlerFunc(func(resp http.ResponseWriter, _ *http.Request) {
				resp.WriteHeader(http.StatusBadRequest)
				resp.Write([]byte(`{"message":"no route handled","httpStatusCode":400,"errorCode":0}`))
			})),
			expectedErr: fmt.Errorf("Error: problem finding link"),
		},
		"error: incorrect volume type": {
			volume: types.Volume{
				ID:         uuid.NewString(),
				VolumeType: "Volume",
			},
			server:      httptest.NewServer(http.HandlerFunc(func(_ http.ResponseWriter, _ *http.Request) {})),
			expectedErr: fmt.Errorf("Volume type should be snapshot"),
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			defer tc.server.Close()

			client, err := NewClientWithArgs(tc.server.URL, "3.6", math.MaxInt64, true, false)
			if err != nil {
				t.Fatal(err)
			}

			vl := NewVolume(client)
			vl.Volume = &tc.volume

			err = vl.LockAutoSnapshot()
			if err != nil {
				if tc.expectedErr.Error() != err.Error() {
					t.Fatal(err)
				}
			}
		})
	}
}

func TestUnlockAutoSnapshot(t *testing.T) {
	volumeID := uuid.NewString()
	type testCase struct {
		volume      types.Volume
		server      *httptest.Server
		expectedErr error
	}

	cases := map[string]testCase{
		"success": {
			volume: types.Volume{
				ID:         volumeID,
				VolumeType: "Snapshot",
				Links: []*types.Link{
					{
						Rel:  "self",
						HREF: fmt.Sprintf("/api/instances/Volume::%s", volumeID),
					},
				},
			},
			server: httptest.NewServer(http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
				switch req.RequestURI {
				case fmt.Sprintf("/api/instances/Volume::%s/action/unlockAutoSnapshot", volumeID):
					resp.WriteHeader(http.StatusOK)
				default:
					resp.WriteHeader(http.StatusBadRequest)
					resp.Write([]byte(`{"message":"no route handled","httpStatusCode":400,"errorCode":0}`))
				}
			})),
			expectedErr: nil,
		},
		"error: bad request": {
			volume: types.Volume{
				ID:         volumeID,
				VolumeType: "Snapshot",
				Links: []*types.Link{
					{
						Rel:  "self",
						HREF: fmt.Sprintf("/api/instances/Volume::%s", volumeID),
					},
				},
			},
			server: httptest.NewServer(http.HandlerFunc(func(resp http.ResponseWriter, _ *http.Request) {
				resp.WriteHeader(http.StatusBadRequest)
				resp.Write([]byte(`{"message":"bad request","httpStatusCode":400,"errorCode":0}`))
			})),
			expectedErr: fmt.Errorf("bad request"),
		},
		"error: problem finding link": {
			volume: types.Volume{
				ID:         uuid.NewString(),
				VolumeType: "Snapshot",
			},
			server: httptest.NewServer(http.HandlerFunc(func(resp http.ResponseWriter, _ *http.Request) {
				resp.WriteHeader(http.StatusBadRequest)
				resp.Write([]byte(`{"message":"no route handled","httpStatusCode":400,"errorCode":0}`))
			})),
			expectedErr: fmt.Errorf("Error: problem finding link"),
		},
		"error: incorrect volume type": {
			volume: types.Volume{
				ID:         uuid.NewString(),
				VolumeType: "Volume",
			},
			server:      httptest.NewServer(http.HandlerFunc(func(_ http.ResponseWriter, _ *http.Request) {})),
			expectedErr: fmt.Errorf("Volume type should be snapshot"),
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			defer tc.server.Close()

			client, err := NewClientWithArgs(tc.server.URL, "3.6", math.MaxInt64, true, false)
			if err != nil {
				t.Fatal(err)
			}

			vl := NewVolume(client)
			vl.Volume = &tc.volume

			err = vl.UnlockAutoSnapshot()
			if err != nil {
				if tc.expectedErr.Error() != err.Error() {
					t.Fatal(err)
				}
			}
		})
	}
}

func TestSetVolumeAccessModeLimit(t *testing.T) {
	volumeID := uuid.NewString()
	type testCase struct {
		volume      types.Volume
		server      *httptest.Server
		expectedErr error
	}

	cases := map[string]testCase{
		"success": {
			volume: types.Volume{
				ID: volumeID,
				Links: []*types.Link{
					{
						Rel:  "self",
						HREF: fmt.Sprintf("/api/instances/Volume::%s", volumeID),
					},
				},
			},
			server: httptest.NewServer(http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
				switch req.RequestURI {
				case fmt.Sprintf("/api/instances/Volume::%s/action/setVolumeAccessModeLimit", volumeID):
					resp.WriteHeader(http.StatusOK)
				default:
					resp.WriteHeader(http.StatusBadRequest)
					resp.Write([]byte(`{"message":"no route handled","httpStatusCode":400,"errorCode":0}`))
				}
			})),
			expectedErr: nil,
		},
		"error: bad request": {
			volume: types.Volume{
				ID: volumeID,
				Links: []*types.Link{
					{
						Rel:  "self",
						HREF: fmt.Sprintf("/api/instances/Volume::%s", volumeID),
					},
				},
			},
			server: httptest.NewServer(http.HandlerFunc(func(resp http.ResponseWriter, _ *http.Request) {
				resp.WriteHeader(http.StatusBadRequest)
				resp.Write([]byte(`{"message":"bad request","httpStatusCode":400,"errorCode":0}`))
			})),
			expectedErr: fmt.Errorf("bad request"),
		},
		"error: problem finding link": {
			volume: types.Volume{
				ID: uuid.NewString(),
			},
			server: httptest.NewServer(http.HandlerFunc(func(resp http.ResponseWriter, _ *http.Request) {
				resp.WriteHeader(http.StatusBadRequest)
				resp.Write([]byte(`{"message":"no route handled","httpStatusCode":400,"errorCode":0}`))
			})),
			expectedErr: fmt.Errorf("Error: problem finding link"),
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			defer tc.server.Close()

			client, err := NewClientWithArgs(tc.server.URL, "3.6", math.MaxInt64, true, false)
			if err != nil {
				t.Fatal(err)
			}

			vl := NewVolume(client)
			vl.Volume = &tc.volume

			err = vl.SetVolumeAccessModeLimit("ReadWrite")
			if err != nil {
				if tc.expectedErr.Error() != err.Error() {
					t.Fatal(err)
				}
			}
		})
	}
}

func TestSetSnapshotSecurity(t *testing.T) {
	volumeID := uuid.NewString()
	type testCase struct {
		volume      types.Volume
		server      *httptest.Server
		expectedErr error
	}

	cases := map[string]testCase{
		"success": {
			volume: types.Volume{
				ID: volumeID,
				Links: []*types.Link{
					{
						Rel:  "self",
						HREF: fmt.Sprintf("/api/instances/Volume::%s", volumeID),
					},
				},
			},
			server: httptest.NewServer(http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
				switch req.RequestURI {
				case fmt.Sprintf("/api/instances/Volume::%s/action/setSnapshotSecurity", volumeID):
					resp.WriteHeader(http.StatusOK)
				default:
					resp.WriteHeader(http.StatusBadRequest)
					resp.Write([]byte(`{"message":"no route handled","httpStatusCode":400,"errorCode":0}`))
				}
			})),
			expectedErr: nil,
		},
		"error: bad request": {
			volume: types.Volume{
				ID: volumeID,
				Links: []*types.Link{
					{
						Rel:  "self",
						HREF: fmt.Sprintf("/api/instances/Volume::%s", volumeID),
					},
				},
			},
			server: httptest.NewServer(http.HandlerFunc(func(resp http.ResponseWriter, _ *http.Request) {
				resp.WriteHeader(http.StatusBadRequest)
				resp.Write([]byte(`{"message":"bad request","httpStatusCode":400,"errorCode":0}`))
			})),
			expectedErr: fmt.Errorf("bad request"),
		},
		"error: problem finding link": {
			volume: types.Volume{
				ID: uuid.NewString(),
			},
			server: httptest.NewServer(http.HandlerFunc(func(resp http.ResponseWriter, _ *http.Request) {
				resp.WriteHeader(http.StatusBadRequest)
				resp.Write([]byte(`{"message":"no route handled","httpStatusCode":400,"errorCode":0}`))
			})),
			expectedErr: fmt.Errorf("Error: problem finding link"),
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			defer tc.server.Close()

			client, err := NewClientWithArgs(tc.server.URL, "3.6", math.MaxInt64, true, false)
			if err != nil {
				t.Fatal(err)
			}

			vl := NewVolume(client)
			vl.Volume = &tc.volume

			err = vl.SetSnapshotSecurity("Insecure")
			if err != nil {
				if tc.expectedErr.Error() != err.Error() {
					t.Fatal(err)
				}
			}
		})
	}
}

func TestSetVolumeMappingAccessMode(t *testing.T) {
	volumeID := uuid.NewString()
	type testCase struct {
		volume      types.Volume
		server      *httptest.Server
		expectedErr error
	}

	cases := map[string]testCase{
		"success": {
			volume: types.Volume{
				ID: volumeID,
				Links: []*types.Link{
					{
						Rel:  "self",
						HREF: fmt.Sprintf("/api/instances/Volume::%s", volumeID),
					},
				},
			},
			server: httptest.NewServer(http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
				switch req.RequestURI {
				case fmt.Sprintf("/api/instances/Volume::%s/action/setVolumeMappingAccessMode", volumeID):
					resp.WriteHeader(http.StatusOK)
				default:
					resp.WriteHeader(http.StatusBadRequest)
					resp.Write([]byte(`{"message":"no route handled","httpStatusCode":400,"errorCode":0}`))
				}
			})),
			expectedErr: nil,
		},
		"error: bad request": {
			volume: types.Volume{
				ID: volumeID,
				Links: []*types.Link{
					{
						Rel:  "self",
						HREF: fmt.Sprintf("/api/instances/Volume::%s", volumeID),
					},
				},
			},
			server: httptest.NewServer(http.HandlerFunc(func(resp http.ResponseWriter, _ *http.Request) {
				resp.WriteHeader(http.StatusBadRequest)
				resp.Write([]byte(`{"message":"bad request","httpStatusCode":400,"errorCode":0}`))
			})),
			expectedErr: fmt.Errorf("bad request"),
		},
		"error: problem finding link": {
			volume: types.Volume{
				ID: uuid.NewString(),
			},
			server:      httptest.NewServer(http.HandlerFunc(func(_ http.ResponseWriter, _ *http.Request) {})),
			expectedErr: fmt.Errorf("Error: problem finding link"),
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			defer tc.server.Close()

			client, err := NewClientWithArgs(tc.server.URL, "3.6", math.MaxInt64, true, false)
			if err != nil {
				t.Fatal(err)
			}

			vl := NewVolume(client)
			vl.Volume = &tc.volume

			err = vl.SetVolumeMappingAccessMode("ReadOnly", uuid.NewString())
			if err != nil {
				if tc.expectedErr.Error() != err.Error() {
					t.Fatal(err)
				}
			}
		})
	}
}

func TestSetVolumeUseRmCache(t *testing.T) {
	volumeID := uuid.NewString()
	type testCase struct {
		volume      types.Volume
		server      *httptest.Server
		expectedErr error
	}

	cases := map[string]testCase{
		"success": {
			volume: types.Volume{
				ID: volumeID,
				Links: []*types.Link{
					{
						Rel:  "self",
						HREF: fmt.Sprintf("/api/instances/Volume::%s", volumeID),
					},
				},
			},
			server: httptest.NewServer(http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
				switch req.RequestURI {
				case fmt.Sprintf("/api/instances/Volume::%s/action/setVolumeUseRmcache", volumeID):
					resp.WriteHeader(http.StatusOK)
				default:
					resp.WriteHeader(http.StatusBadRequest)
					resp.Write([]byte(`{"message":"no route handled","httpStatusCode":400,"errorCode":0}`))
				}
			})),
			expectedErr: nil,
		},
		"error: bad request": {
			volume: types.Volume{
				ID: volumeID,
				Links: []*types.Link{
					{
						Rel:  "self",
						HREF: fmt.Sprintf("/api/instances/Volume::%s", volumeID),
					},
				},
			},
			server: httptest.NewServer(http.HandlerFunc(func(resp http.ResponseWriter, _ *http.Request) {
				resp.WriteHeader(http.StatusBadRequest)
				resp.Write([]byte(`{"message":"bad request","httpStatusCode":400,"errorCode":0}`))
			})),
			expectedErr: fmt.Errorf("bad request"),
		},
		"error: problem finding link": {
			volume: types.Volume{
				ID: uuid.NewString(),
			},
			server:      httptest.NewServer(http.HandlerFunc(func(_ http.ResponseWriter, _ *http.Request) {})),
			expectedErr: fmt.Errorf("Error: problem finding link"),
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			defer tc.server.Close()

			client, err := NewClientWithArgs(tc.server.URL, "3.6", math.MaxInt64, true, false)
			if err != nil {
				t.Fatal(err)
			}

			vl := NewVolume(client)
			vl.Volume = &tc.volume

			err = vl.SetVolumeUseRmCache(true)
			if err != nil {
				if tc.expectedErr.Error() != err.Error() {
					t.Fatal(err)
				}
			}
		})
	}
}

func TestSetCompressionMethod(t *testing.T) {
	volumeID := uuid.NewString()
	type testCase struct {
		volume      types.Volume
		server      *httptest.Server
		expectedErr error
	}

	cases := map[string]testCase{
		"success": {
			volume: types.Volume{
				ID: volumeID,
				Links: []*types.Link{
					{
						Rel:  "self",
						HREF: fmt.Sprintf("/api/instances/Volume::%s", volumeID),
					},
				},
			},
			server: httptest.NewServer(http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
				switch req.RequestURI {
				case fmt.Sprintf("/api/instances/Volume::%s/action/modifyCompressionMethod", volumeID):
					resp.WriteHeader(http.StatusOK)
				default:
					resp.WriteHeader(http.StatusBadRequest)
					resp.Write([]byte(`{"message":"no route handled","httpStatusCode":400,"errorCode":0}`))
				}
			})),
			expectedErr: nil,
		},
		"error: bad request": {
			volume: types.Volume{
				ID: volumeID,
				Links: []*types.Link{
					{
						Rel:  "self",
						HREF: fmt.Sprintf("/api/instances/Volume::%s", volumeID),
					},
				},
			},
			server: httptest.NewServer(http.HandlerFunc(func(resp http.ResponseWriter, _ *http.Request) {
				resp.WriteHeader(http.StatusBadRequest)
				resp.Write([]byte(`{"message":"bad request","httpStatusCode":400,"errorCode":0}`))
			})),
			expectedErr: fmt.Errorf("bad request"),
		},
		"error: problem finding link": {
			volume: types.Volume{
				ID: uuid.NewString(),
			},
			server:      httptest.NewServer(http.HandlerFunc(func(_ http.ResponseWriter, _ *http.Request) {})),
			expectedErr: fmt.Errorf("Error: problem finding link"),
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			defer tc.server.Close()

			client, err := NewClientWithArgs(tc.server.URL, "3.6", math.MaxInt64, true, false)
			if err != nil {
				t.Fatal(err)
			}

			vl := NewVolume(client)
			vl.Volume = &tc.volume

			err = vl.SetCompressionMethod("gzip")
			if err != nil {
				if tc.expectedErr.Error() != err.Error() {
					t.Fatal(err)
				}
			}
		})
	}
}

func TestUnmarkForReplication(t *testing.T) {
	volumeID := uuid.NewString()
	type testCase struct {
		volume      types.Volume
		server      *httptest.Server
		expectedErr error
	}

	cases := map[string]testCase{
		"success": {
			volume: types.Volume{
				ID: volumeID,
			},
			server: httptest.NewServer(http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
				switch req.RequestURI {
				case fmt.Sprintf("/api/instances/Volume::%s/action/unmarkForReplication", volumeID):
					resp.WriteHeader(http.StatusOK)
				default:
					resp.WriteHeader(http.StatusBadRequest)
					resp.Write([]byte(`{"message":"no route handled","httpStatusCode":400,"errorCode":0}`))
				}
			})),
			expectedErr: nil,
		},
		"error: bad request": {
			volume: types.Volume{
				ID: volumeID,
			},
			server: httptest.NewServer(http.HandlerFunc(func(resp http.ResponseWriter, _ *http.Request) {
				resp.WriteHeader(http.StatusBadRequest)
				resp.Write([]byte(`{"message":"bad request","httpStatusCode":400,"errorCode":0}`))
			})),
			expectedErr: fmt.Errorf("bad request"),
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			defer tc.server.Close()

			client, err := NewClientWithArgs(tc.server.URL, "3.6", math.MaxInt64, true, false)
			if err != nil {
				t.Fatal(err)
			}

			vl := NewVolume(client)
			vl.Volume = &tc.volume

			err = vl.UnmarkForReplication()
			if err != nil {
				if tc.expectedErr.Error() != err.Error() {
					t.Fatal(err)
				}
			}
		})
	}
}
