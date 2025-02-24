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
	"errors"
	"fmt"
	"math"
	"net/http"
	"net/http/httptest"
	"testing"

	types "github.com/dell/goscaleio/types/v1"
	"github.com/google/uuid"
)

var (
	ID2     string
	errSnap error
)

func TestCreateSnapshotPolicy(t *testing.T) {
	type testCase struct {
		fs       *types.SnapshotPolicyCreateParam
		expected error
	}
	cases := []testCase{
		{
			fs: &types.SnapshotPolicyCreateParam{
				Name:                             "testSnapshotPolicy",
				AutoSnapshotCreationCadenceInMin: "5",
				NumOfRetainedSnapshotsPerLevel:   []string{"1"},
			},
			expected: nil,
		},
		{
			fs: &types.SnapshotPolicyCreateParam{
				Name:                             "testSnapshotPolicy2",
				AutoSnapshotCreationCadenceInMin: "5",
				NumOfRetainedSnapshotsPerLevel:   []string{"1"},
				SnapshotAccessMode:               "Invalid",
			},
			expected: errors.New("500 Internal Server Error"),
		},
	}
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/types/SnapshotPolicy/instances" {
			if r.Method == http.MethodPost {
				var param types.SnapshotPolicyCreateParam
				_ = json.NewDecoder(r.Body).Decode(&param)
				switch param.Name {
				case "testSnapshotPolicy":
					w.WriteHeader(http.StatusOK)
				case "testSnapshotPolicy2":
					w.WriteHeader(http.StatusInternalServerError)
					w.Write([]byte(`{"error": "accessMode should get one of the following values: ReadWrite, ReadOnly, but its value is Invalid"}`))
				}
			}
		}
	}))
	defer svr.Close()

	for _, tc := range cases {
		tc := tc
		t.Run("", func(_ *testing.T) {
			client, err := NewClientWithArgs(svr.URL, "", math.MaxInt64, true, false)
			if err != nil {
				t.Fatal(err)
			}

			s := System{
				client: client,
			}

			ID2, errSnap = s.CreateSnapshotPolicy(tc.fs)
			if errSnap != nil {
				if tc.expected == nil {
					t.Errorf("Creating Snapshot Policy did not work as expected, \n\tgot: %s \n\twant: %v", errSnap, tc.expected)
				} else {
					if errSnap.Error() != tc.expected.Error() {
						t.Errorf("Creating Snapshot Policy did not work as expected, \n\tgot: %s \n\twant: %s", errSnap, tc.expected)
					}
				}
			}
		})
	}
}

func TestRenameSnapshotPolicy(t *testing.T) {
	type testCase struct {
		name     string
		id       string
		expected error
	}
	cases := []testCase{
		{
			id:       ID2,
			name:     "renameSnapshotPolicy",
			expected: nil,
		},
		{
			id:       "1234",
			name:     "renameSnapshotPolicy",
			expected: errors.New("500 Internal Server Error"),
		},
	}
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/instances/SnapshotPolicy::"+ID2+"/action/renameSnapshotPolicy" {
			w.WriteHeader(http.StatusOK)
		} else if r.URL.Path == "/api/instances/SnapshotPolicy::"+"1234"+"/action/renameSnapshotPolicy" {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(`{"error": "id (Invalid) must be a hexadecimal number (unsigned long)"}`))
		}

	}))
	defer svr.Close()

	for _, tc := range cases {
		tc := tc
		t.Run("", func(_ *testing.T) {
			client, err := NewClientWithArgs(svr.URL, "", math.MaxInt64, true, false)
			if err != nil {
				t.Fatal(err)
			}

			s := System{
				client: client,
			}
			err2 := s.RenameSnapshotPolicy(tc.id, tc.name)
			if err2 != nil {
				if tc.expected == nil {
					t.Errorf("Modifying snapshot policy did not work as expected, \n\tgot: %s \n\twant: %v", err2, tc.expected)
				} else {
					if err2.Error() != tc.expected.Error() {
						t.Errorf("Modifying snapshot policy did not work as expected, \n\tgot: %s \n\twant: %s", err2, tc.expected)
					}
				}
			}
		})
	}
}

func TestModifySnapshotPolicy(t *testing.T) {
	type testCase struct {
		id       string
		snap     *types.SnapshotPolicyModifyParam
		expected error
	}
	cases := []testCase{
		{
			id: ID2,
			snap: &types.SnapshotPolicyModifyParam{
				AutoSnapshotCreationCadenceInMin: "6",
				NumOfRetainedSnapshotsPerLevel:   []string{"2", "3"},
			},
			expected: nil,
		},
		{
			id: "Invalid",
			snap: &types.SnapshotPolicyModifyParam{
				AutoSnapshotCreationCadenceInMin: "6",
				NumOfRetainedSnapshotsPerLevel:   []string{"2", "3"},
			},
			expected: errors.New("500 Internal Server Error"),
		},
	}
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/instances/SnapshotPolicy::"+ID2+"/action/modifySnapshotPolicy" {
			w.WriteHeader(http.StatusOK)
		} else if r.URL.Path == "/api/instances/SnapshotPolicy::"+"Invalid"+"/action/modifySnapshotPolicy" {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(`{"error": "id (Invalid) must be a hexadecimal number (unsigned long)"}`))
		}
	}))
	defer svr.Close()

	for _, tc := range cases {
		tc := tc
		t.Run("", func(_ *testing.T) {
			client, err := NewClientWithArgs(svr.URL, "", math.MaxInt64, true, false)
			if err != nil {
				t.Fatal(err)
			}

			s := System{
				client: client,
			}

			err2 := s.ModifySnapshotPolicy(tc.snap, tc.id)
			if err2 != nil {
				if tc.expected == nil {
					t.Errorf("Modifying snapshot policy did not work as expected, \n\tgot: %s \n\twant: %v", err2, tc.expected)
				} else {
					if err2.Error() != tc.expected.Error() {
						t.Errorf("Modifying snapshot policy did not work as expected, \n\tgot: %s \n\twant: %s", err2, tc.expected)
					}
				}
			}
		})
	}
}

func TestAssignVolumeToSnapshotPolicy(t *testing.T) {
	type testCase struct {
		id       string
		snap     *types.AssignVolumeToSnapshotPolicyParam
		expected error
	}
	cases := []testCase{
		{
			id: ID2,
			snap: &types.AssignVolumeToSnapshotPolicyParam{
				SourceVolumeID: "edba1bff00000001",
			},
			expected: nil,
		},
		{
			id: "Invalid",
			snap: &types.AssignVolumeToSnapshotPolicyParam{
				SourceVolumeID: "edba1bff00000001",
			},
			expected: errors.New("500 Internal Server Error"),
		},
	}
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/instances/SnapshotPolicy::"+ID2+"/action/addSourceVolumeToSnapshotPolicy" {
			w.WriteHeader(http.StatusOK)
		} else if r.URL.Path == "/api/instances/SnapshotPolicy::"+"Invalid"+"/action/addSourceVolumeToSnapshotPolicy" {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(`{"error": "id (Invalid) must be a hexadecimal number (unsigned long)"}`))
		}
	}))
	defer svr.Close()

	for _, tc := range cases {
		tc := tc
		t.Run("", func(_ *testing.T) {
			client, err := NewClientWithArgs(svr.URL, "", math.MaxInt64, true, false)
			if err != nil {
				t.Fatal(err)
			}

			s := System{
				client: client,
			}

			err2 := s.AssignVolumeToSnapshotPolicy(tc.snap, tc.id)
			if err2 != nil {
				if tc.expected == nil {
					t.Errorf("Assigning volume to snapshot policy did not work as expected, \n\tgot: %s \n\twant: %v", err2, tc.expected)
				} else {
					if err2.Error() != tc.expected.Error() {
						t.Errorf("Assigning volume to snapshot policy did not work as expected, \n\tgot: %s \n\twant: %s", err2, tc.expected)
					}
				}
			}
		})
	}
}

func TestUnassignVolumeFromSnapshotPolicy(t *testing.T) {
	policyID := uuid.NewString()
	systemID := uuid.NewString()
	type testCase struct {
		server      *httptest.Server
		expectedErr error
	}

	cases := map[string]testCase{
		"success": {
			server: httptest.NewServer(http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
				switch req.RequestURI {
				case fmt.Sprintf("/api/instances/SnapshotPolicy::%s/action/removeSourceVolumeFromSnapshotPolicy", policyID):
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
	}

	for _, tc := range cases {
		client, err := NewClientWithArgs(tc.server.URL, "3.6", math.MaxInt64, true, false)
		if err != nil {
			t.Fatal(err)
		}
		defer tc.server.Close()

		system := System{
			System: &types.System{
				ID: systemID,
			},
			client: client,
		}

		params := &types.AssignVolumeToSnapshotPolicyParam{
			SourceVolumeID: uuid.NewString(),
		}

		err = system.UnassignVolumeFromSnapshotPolicy(params, policyID)
		if err != nil {
			if tc.expectedErr.Error() != err.Error() {
				t.Fatal(err)
			}
		}
	}
}

func TestPauseSnapshotPolicy(t *testing.T) {
	type testCase struct {
		id       string
		expected error
	}
	cases := []testCase{
		{
			id:       ID2,
			expected: nil,
		},
		{
			id:       "Invalid",
			expected: errors.New("500 Internal Server Error"),
		},
	}
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/instances/SnapshotPolicy::"+ID2+"/action/pauseSnapshotPolicy" {
			w.WriteHeader(http.StatusOK)
		} else if r.URL.Path == "/api/instances/SnapshotPolicy::"+"Invalid"+"/action/pauseSnapshotPolicy" {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(`{"error": "id (Invalid) must be a hexadecimal number (unsigned long)"}`))
		}
	}))
	defer svr.Close()

	for _, tc := range cases {
		tc := tc
		t.Run("", func(_ *testing.T) {
			client, err := NewClientWithArgs(svr.URL, "", math.MaxInt64, true, false)
			if err != nil {
				t.Fatal(err)
			}

			s := System{
				client: client,
			}

			err2 := s.PauseSnapshotPolicy(tc.id)
			if err2 != nil {
				if tc.expected == nil {
					t.Errorf("Pausing snapshot policy did not work as expected, \n\tgot: %s \n\twant: %v", err2, tc.expected)
				} else {
					if err2.Error() != tc.expected.Error() {
						t.Errorf("Pausing snapshot policy did not work as expected, \n\tgot: %s \n\twant: %s", err2, tc.expected)
					}
				}
			}
		})
	}
}

func TestResumeSnapshotPolicy(t *testing.T) {
	type testCase struct {
		id       string
		expected error
	}
	cases := []testCase{
		{
			id:       ID2,
			expected: nil,
		},
		{
			id:       "Invalid",
			expected: errors.New("500 Internal Server Error"),
		},
	}
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/instances/SnapshotPolicy::"+ID2+"/action/resumeSnapshotPolicy" {
			w.WriteHeader(http.StatusOK)
		} else if r.URL.Path == "/api/instances/SnapshotPolicy::"+"Invalid"+"/action/resumeSnapshotPolicy" {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(`{"error": "id (Invalid) must be a hexadecimal number (unsigned long)"}`))
		}
	}))
	defer svr.Close()

	for _, tc := range cases {
		tc := tc
		t.Run("", func(_ *testing.T) {
			client, err := NewClientWithArgs(svr.URL, "", math.MaxInt64, true, false)
			if err != nil {
				t.Fatal(err)
			}

			s := System{
				client: client,
			}

			err2 := s.ResumeSnapshotPolicy(tc.id)
			if err2 != nil {
				if tc.expected == nil {
					t.Errorf("Resuming snapshot policy did not work as expected, \n\tgot: %s \n\twant: %v", err2, tc.expected)
				} else {
					if err2.Error() != tc.expected.Error() {
						t.Errorf("Resuming snapshot policy did not work as expected, \n\tgot: %s \n\twant: %s", err2, tc.expected)
					}
				}
			}
		})
	}
}

func TestRemoveSnapshotPolicy(t *testing.T) {
	type testCase struct {
		id       string
		expected error
	}
	cases := []testCase{
		{
			id:       ID2,
			expected: nil,
		},
		{
			id:       "Invalid",
			expected: errors.New("500 Internal Server Error"),
		},
	}
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/instances/SnapshotPolicy::"+ID2+"/action/removeSnapshotPolicy" {
			w.WriteHeader(http.StatusOK)
		} else if r.URL.Path == "/api/instances/SnapshotPolicy::"+"Invalid"+"/action/removeSnapshotPolicy" {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(`{"error": "id (Invalid) must be a hexadecimal number (unsigned long)"}`))
		}
	}))
	defer svr.Close()

	for _, tc := range cases {
		tc := tc
		t.Run("", func(_ *testing.T) {
			client, err := NewClientWithArgs(svr.URL, "", math.MaxInt64, true, false)
			if err != nil {
				t.Fatal(err)
			}

			s := System{
				client: client,
			}

			err2 := s.RemoveSnapshotPolicy(tc.id)
			if err2 != nil {
				if tc.expected == nil {
					t.Errorf("Removing snapshot policy did not work as expected, \n\tgot: %s \n\twant: %v", err2, tc.expected)
				} else {
					if err2.Error() != tc.expected.Error() {
						t.Errorf("Removing snapshot policy did not work as expected, \n\tgot: %s \n\twant: %s", err2, tc.expected)
					}
				}
			}
		})
	}
}

func TestGetSourceVolume(t *testing.T) {
	type testCase struct {
		id       string
		expected error
	}
	cases := []testCase{
		{
			id:       ID2,
			expected: nil,
		},
		{
			id:       "Invalid",
			expected: errors.New("500 Internal Server Error"),
		},
	}
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/instances/SnapshotPolicy::"+ID2+"/relationships/SourceVolume" {
			w.WriteHeader(http.StatusOK)
		} else if r.URL.Path == "/api/instances/SnapshotPolicy::"+"Invalid"+"/relationships/SourceVolume" {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(`{"error": "id (Invalid) must be a hexadecimal number (unsigned long)"}`))
		}
	}))
	defer svr.Close()

	for _, tc := range cases {
		tc := tc
		t.Run("", func(_ *testing.T) {
			client, err := NewClientWithArgs(svr.URL, "", math.MaxInt64, true, false)
			if err != nil {
				t.Fatal(err)
			}

			s := System{
				client: client,
			}

			_, err2 := s.GetSourceVolume(tc.id)
			if err2 != nil {
				if tc.expected == nil {
					t.Errorf("Assigning volume to snapshot policy did not work as expected, \n\tgot: %s \n\twant: %v", err2, tc.expected)
				} else {
					if err2.Error() != tc.expected.Error() {
						t.Errorf("Assigning volume to snapshot policy did not work as expected, \n\tgot: %s \n\twant: %s", err2, tc.expected)
					}
				}
			}
		})
	}
}
