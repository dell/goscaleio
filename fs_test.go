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
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"net/http"
	"net/http/httptest"
	"testing"

	types "github.com/dell/goscaleio/types/v1"
	"github.com/stretchr/testify/assert"
)

func TestNewFileSystem(t *testing.T) {
	client := &Client{}
	existingFileSystem := &types.FileSystem{}
	fs := NewFileSystem(client, existingFileSystem)

	assert.NotNil(t, fs)
	assert.Equal(t, existingFileSystem, fs.FileSystem)
	assert.Equal(t, client, fs.client)
}

func TestGetFileSystemByIDName(t *testing.T) {
	type checkFn func(*testing.T, *types.FileSystem, error)
	check := func(fns ...checkFn) []checkFn { return fns }

	hasNoError := func(t *testing.T, _ *types.FileSystem, err error) {
		if err != nil {
			t.Fatalf("expected no error")
		}
	}

	hasError := func(t *testing.T, _ *types.FileSystem, err error) {
		if err == nil {
			t.Fatalf("expected error")
		}
	}

	checkRespName := func(fsName string) func(t *testing.T, resp *types.FileSystem, err error) {
		return func(t *testing.T, resp *types.FileSystem, _ error) {
			assert.Equal(t, fsName, resp.Name)
		}
	}

	checkRespID := func(fsID string) func(t *testing.T, resp *types.FileSystem, err error) {
		return func(t *testing.T, resp *types.FileSystem, _ error) {
			assert.Equal(t, fsID, resp.ID)
		}
	}

	testsName := map[string]func(t *testing.T) (*httptest.Server, []checkFn){
		"success": func(t *testing.T) (*httptest.Server, []checkFn) {
			href := "/rest/v1/file-systems"

			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method != http.MethodGet {
					t.Fatal(fmt.Errorf("wrong method. Expected %s; but got %s", http.MethodGet, r.Method))
				}

				if r.URL.Path != href {
					t.Fatal(fmt.Errorf("wrong path. Expected %s; but got %s", href, r.URL.Path))
				}

				resp := []types.FileSystem{
					{
						ID:   "64366a19-54e8-1544-f3d7-2a50fb1ccff3",
						Name: "fs-test-1",
					},
					{
						ID:   "6436aa58-e6a1-a4e2-de7b-2a50fb1ccff3",
						Name: "fs-test-2",
					},
				}

				respData, err := json.Marshal(resp)
				if err != nil {
					t.Fatal(err)
				}
				fmt.Fprintln(w, string(respData))
			}))
			return ts, check(hasNoError, checkRespName("fs-test-2"))
		},
		"not found": func(t *testing.T) (*httptest.Server, []checkFn) {
			href := "/rest/v1/file-systems"

			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method != http.MethodGet {
					t.Fatal(fmt.Errorf("wrong method. Expected %s; but got %s", http.MethodGet, r.Method))
				}

				if r.URL.Path != href {
					t.Fatal(fmt.Errorf("wrong path. Expected %s; but got %s", href, r.URL.Path))
				}

				resp := []types.FileSystem{
					{
						ID:   "64366a19-54e8-1544-f3d7-2a50fb1ccff3",
						Name: "fs-test-1",
					},
					{
						ID:   "6436aa58-e6a1-a4e2-de7b-2a50fb1ccff3",
						Name: "fs-test-2",
					},
				}

				respData, err := json.Marshal(resp)
				if err != nil {
					t.Fatal(err)
				}
				fmt.Fprintln(w, string(respData))
			}))
			return ts, check(hasError)
		},
	}

	testsID := map[string]func(t *testing.T) (*httptest.Server, []checkFn){
		"success": func(t *testing.T) (*httptest.Server, []checkFn) {
			fsID := "64366a19-54e8-1544-f3d7-2a50fb1ccff3"
			href := fmt.Sprintf("/rest/v1/file-systems/%s", fsID)

			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method != http.MethodGet {
					t.Fatal(fmt.Errorf("wrong method. Expected %s; but got %s", http.MethodGet, r.Method))
				}

				if r.URL.Path != href {
					t.Fatal(fmt.Errorf("wrong path. Expected %s; but got %s", href, r.URL.Path))
				}

				resp := types.FileSystem{
					ID:   "64366a19-54e8-1544-f3d7-2a50fb1ccff3",
					Name: "fs-test-1",
				}

				respData, err := json.Marshal(resp)
				if err != nil {
					t.Fatal(err)
				}
				fmt.Fprintln(w, string(respData))
			}))
			return ts, check(hasNoError, checkRespID("64366a19-54e8-1544-f3d7-2a50fb1ccff3"))
		},
		"not found": func(t *testing.T) (*httptest.Server, []checkFn) {
			fsID := "6436aa58-e6a1-a4e2-de7b-2a50fb1ccff3"
			href := fmt.Sprintf("/rest/v1/file-systems/%s", fsID)

			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method != http.MethodGet {
					t.Fatal(fmt.Errorf("wrong method. Expected %s; but got %s", http.MethodGet, r.Method))
				}

				if r.URL.Path != href {
					t.Fatal(fmt.Errorf("wrong path. Expected %s; but got %s", href, r.URL.Path))
				}

				http.Error(w, "nas not found", http.StatusNotFound)
			}))
			return ts, check(hasError)
		},
	}

	testCaseFSNames := map[string]string{
		"success":   "fs-test-2",
		"not found": "fs-test-3",
	}

	testCaseFSIds := map[string]string{
		"success":   "64366a19-54e8-1544-f3d7-2a50fb1ccff3",
		"not found": "6436aa58-e6a1-a4e2-de7b-2a50fb1ccff3",
	}

	for name, tc := range testsName {
		t.Run(name, func(t *testing.T) {
			ts, checkFns := tc(t)
			defer ts.Close()

			client, err := NewClientWithArgs(ts.URL, "", math.MaxInt64, true, false)
			client.configConnect.Version = "4.0"
			if err != nil {
				t.Fatal(err)
			}

			s := System{
				client: client,
			}

			resp, err := s.GetFileSystemByIDName("", testCaseFSNames[name])
			for _, checkFn := range checkFns {
				checkFn(t, resp, err)
			}
		})
	}

	for id, tc := range testsID {
		t.Run(id, func(t *testing.T) {
			ts, checkFns := tc(t)
			defer ts.Close()

			client, err := NewClientWithArgs(ts.URL, "", math.MaxInt64, true, false)
			client.configConnect.Version = "4.0"
			if err != nil {
				t.Fatal(err)
			}

			s := System{
				client: client,
			}

			resp, err := s.GetFileSystemByIDName(testCaseFSIds[id], "")
			for _, checkFn := range checkFns {
				checkFn(t, resp, err)
			}
		})
	}
}

func TestCreateFileSystem(t *testing.T) {
	type checkFn func(*testing.T, *types.FileSystemResp, error)
	check := func(fns ...checkFn) []checkFn { return fns }

	hasNoError := func(t *testing.T, _ *types.FileSystemResp, err error) {
		if err != nil {
			t.Fatalf("expected no error")
		}
	}

	hasError := func(t *testing.T, _ *types.FileSystemResp, err error) {
		if err == nil {
			t.Fatalf("expected error")
		}
	}

	checkResp := func(fsId string) func(t *testing.T, resp *types.FileSystemResp, err error) {
		return func(t *testing.T, resp *types.FileSystemResp, _ error) {
			assert.Equal(t, fsId, resp.ID)
		}
	}

	tests := map[string]func(t *testing.T) (*httptest.Server, []checkFn){
		"success": func(t *testing.T) (*httptest.Server, []checkFn) {
			href := "/rest/v1/file-systems"

			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method != http.MethodPost {
					t.Fatal(fmt.Errorf("wrong method. Expected %s; but got %s", http.MethodGet, r.Method))
				}

				if r.URL.Path != href {
					t.Fatal(fmt.Errorf("wrong path. Expected %s; but got %s", href, r.URL.Path))
				}

				resp := types.FileSystemResp{
					ID: "64366a19-54e8-1544-f3d7-2a50fb1ccff3",
				}

				respData, err := json.Marshal(resp)
				if err != nil {
					t.Fatal(err)
				}
				fmt.Fprintln(w, string(respData))
			}))
			return ts, check(hasNoError, checkResp("64366a19-54e8-1544-f3d7-2a50fb1ccff3"))
		},
		"bad request": func(t *testing.T) (*httptest.Server, []checkFn) {
			href := "/rest/v1/file-systems"

			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method != http.MethodPost {
					t.Fatal(fmt.Errorf("wrong method. Expected %s; but got %s", http.MethodGet, r.Method))
				}

				if r.URL.Path != href {
					t.Fatal(fmt.Errorf("wrong path. Expected %s; but got %s", href, r.URL.Path))
				}

				http.Error(w, "bad Request", http.StatusBadRequest)
			}))
			return ts, check(hasError)
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			ts, checkFns := tc(t)
			defer ts.Close()

			client, err := NewClientWithArgs(ts.URL, "", math.MaxInt64, true, false)
			client.configConnect.Version = "4.0"
			if err != nil {
				t.Fatal(err)
			}

			s := System{
				client: client,
			}

			fs := &types.FsCreate{
				Name:          "test-FS",
				SizeTotal:     16106127360,
				StoragePoolID: "28515fee00000000",
				NasServerID:   "64132f37-d33e-9d4a-89ba-d625520a4779",
			}

			resp, err := s.CreateFileSystem(fs)
			for _, checkFn := range checkFns {
				checkFn(t, resp, err)
			}
		})
	}
}

func TestCreateFileSystemSnapshot(t *testing.T) {
	type checkFn func(*testing.T, *types.CreateFileSystemSnapshotResponse, error)
	check := func(fns ...checkFn) []checkFn { return fns }

	hasNoError := func(t *testing.T, _ *types.CreateFileSystemSnapshotResponse, err error) {
		if err != nil {
			t.Fatalf("expected no error")
		}
	}

	hasError := func(t *testing.T, _ *types.CreateFileSystemSnapshotResponse, err error) {
		if err == nil {
			t.Fatalf("expected error")
		}
	}

	checkResp := func(snapId string) func(t *testing.T, resp *types.CreateFileSystemSnapshotResponse, err error) {
		return func(t *testing.T, resp *types.CreateFileSystemSnapshotResponse, _ error) {
			assert.Equal(t, snapId, resp.ID)
		}
	}

	tests := map[string]func(t *testing.T) (*httptest.Server, []checkFn){
		"success": func(t *testing.T) (*httptest.Server, []checkFn) {
			href := fmt.Sprintf("/rest/v1/file-systems/%v/snapshot", "64366a19-54e8-1544-f3d7-2a50fb1ccff3")

			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method != http.MethodPost {
					t.Fatal(fmt.Errorf("wrong method. Expected %s; but got %s", http.MethodGet, r.Method))
				}

				if r.URL.Path != href {
					t.Fatal(fmt.Errorf("wrong path. Expected %s; but got %s", href, r.URL.Path))
				}

				resp := types.CreateFileSystemSnapshotResponse{
					ID: "64366a19-54e8-1544-f3d7-2a50fb1ccff3",
				}

				respData, err := json.Marshal(resp)
				if err != nil {
					t.Fatal(err)
				}
				fmt.Fprintln(w, string(respData))
			}))
			return ts, check(hasNoError, checkResp("64366a19-54e8-1544-f3d7-2a50fb1ccff3"))
		},
		"bad request": func(t *testing.T) (*httptest.Server, []checkFn) {
			href := fmt.Sprintf("/rest/v1/file-systems/%v/snapshot", "64366a19-54e8-1544-f3d7-2a50fb1ccff3")

			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method != http.MethodPost {
					t.Fatal(fmt.Errorf("wrong method. Expected %s; but got %s", http.MethodGet, r.Method))
				}

				if r.URL.Path != href {
					t.Fatal(fmt.Errorf("wrong path. Expected %s; but got %s", href, r.URL.Path))
				}

				http.Error(w, "bad Request", http.StatusBadRequest)
			}))
			return ts, check(hasError)
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			ts, checkFns := tc(t)
			defer ts.Close()

			client, err := NewClientWithArgs(ts.URL, "", math.MaxInt64, true, false)
			client.configConnect.Version = "4.0"
			if err != nil {
				t.Fatal(err)
			}

			s := System{
				client: client,
			}

			fsID := "64366a19-54e8-1544-f3d7-2a50fb1ccff3"

			fsSnapRequest := &types.CreateFileSystemSnapshotParam{
				Name: "test-snapshot",
			}

			resp, err := s.CreateFileSystemSnapshot(fsSnapRequest, fsID)
			for _, checkFn := range checkFns {
				checkFn(t, resp, err)
			}
		})
	}
}

func TestGetFsSnapshotsByVolumeID(t *testing.T) {
	type checkFn func(*testing.T, []types.FileSystem, error)
	check := func(fns ...checkFn) []checkFn { return fns }

	hasNoError := func(t *testing.T, _ []types.FileSystem, err error) {
		if err != nil {
			t.Fatalf("expected no error")
		}
	}

	hasError := func(t *testing.T, _ []types.FileSystem, err error) {
		if err == nil {
			t.Fatalf("expected error")
		}
	}

	checkResp := func(snapLength int) func(t *testing.T, resp []types.FileSystem, err error) {
		return func(t *testing.T, resp []types.FileSystem, _ error) {
			assert.Equal(t, snapLength, len(resp))
		}
	}

	tests := map[string]func(t *testing.T) (*httptest.Server, []checkFn){
		"success": func(t *testing.T) (*httptest.Server, []checkFn) {
			href := "/rest/v1/file-systems"
			var resp []types.FileSystem

			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method != http.MethodGet {
					t.Fatal(fmt.Errorf("wrong method. Expected %s; but got %s", http.MethodGet, r.Method))
				}

				if r.URL.Path != href {
					t.Fatal(fmt.Errorf("wrong path. Expected %s; but got %s", href, r.URL.Path))
				}

				resp = []types.FileSystem{
					{
						ID:   "64366a19-54e8-1544-f3d7-2a50fb1ccff3",
						Name: "fs-test-1",
					},
					{
						ID:   "6436aa58-e6a1-a4e2-de7b-2a50fb1ccff3",
						Name: "fs-test-2",
					},
				}

				respData, err := json.Marshal(resp)
				if err != nil {
					t.Fatal(err)
				}
				fmt.Fprintln(w, string(respData))
			}))
			return ts, check(hasNoError, checkResp(len(resp)))
		},

		"operation-failed": func(t *testing.T) (*httptest.Server, []checkFn) {
			href := "/rest/v1/file-systems"

			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method != http.MethodGet {
					t.Fatal(fmt.Errorf("wrong method. Expected %s; but got %s", http.MethodGet, r.Method))
				}

				if r.URL.Path != href {
					t.Fatal(fmt.Errorf("wrong path. Expected %s; but got %s", href, r.URL.Path))
				}

				http.Error(w, "operation failed", http.StatusUnprocessableEntity)
			}))
			return ts, check(hasError)
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			ts, checkFns := tc(t)
			defer ts.Close()

			client, err := NewClientWithArgs(ts.URL, "", math.MaxInt64, true, false)
			client.configConnect.Version = "4.0"
			if err != nil {
				t.Fatal(err)
			}

			s := System{
				client: client,
			}

			fsID := "64366a19-54e8-1544-f3d7-2a50fb1ccff3"

			resp, err := s.GetFsSnapshotsByVolumeID(fsID)
			for _, checkFn := range checkFns {
				checkFn(t, resp, err)
			}
		})
	}
}

func TestRestoreFileSystemFromSnapshot(t *testing.T) {
	type checkFn func(*testing.T, *types.RestoreFsSnapResponse, error)
	check := func(fns ...checkFn) []checkFn { return fns }

	hasNoError := func(t *testing.T, _ *types.RestoreFsSnapResponse, err error) {
		if err != nil {
			t.Fatalf("expected no error")
		}
	}

	hasError := func(t *testing.T, _ *types.RestoreFsSnapResponse, err error) {
		if err == nil {
			t.Fatalf("expected error")
		}
	}

	checkResp := func(snapId string) func(t *testing.T, resp *types.RestoreFsSnapResponse, err error) {
		return func(t *testing.T, resp *types.RestoreFsSnapResponse, _ error) {
			assert.Equal(t, snapId, resp.ID)
		}
	}

	tests := map[string]func(t *testing.T) (*httptest.Server, []checkFn){
		"successNoContent": func(t *testing.T) (*httptest.Server, []checkFn) {
			href := fmt.Sprintf("/rest/v1/file-systems/%v/restore", "64366a19-54e8-1544-f3d7-2a50fb1ccff3")

			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method != http.MethodPost {
					t.Fatal(fmt.Errorf("wrong method. Expected %s; but got %s", http.MethodGet, r.Method))
				}

				if r.URL.Path != href {
					t.Fatal(fmt.Errorf("wrong path. Expected %s; but got %s", href, r.URL.Path))
				}

				w.WriteHeader(http.StatusNoContent)
			}))
			return ts, check(hasNoError)
		},

		"successWithContent": func(t *testing.T) (*httptest.Server, []checkFn) {
			href := fmt.Sprintf("/rest/v1/file-systems/%v/restore", "64366a19-54e8-1544-f3d7-2a50fb1ccff3")

			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method != http.MethodPost {
					t.Fatal(fmt.Errorf("wrong method. Expected %s; but got %s", http.MethodGet, r.Method))
				}

				if r.URL.Path != href {
					t.Fatal(fmt.Errorf("wrong path. Expected %s; but got %s", href, r.URL.Path))
				}

				resp := types.RestoreFsSnapResponse{
					ID: "64366a19-54e8-1544-f3d7-2a50fb1cckk3",
				}

				respData, err := json.Marshal(resp)
				if err != nil {
					t.Fatal(err)
				}
				fmt.Fprintln(w, string(respData))
			}))
			return ts, check(hasNoError, checkResp("64366a19-54e8-1544-f3d7-2a50fb1cckk3"))
		},
		"not-found": func(t *testing.T) (*httptest.Server, []checkFn) {
			href := fmt.Sprintf("/rest/v1/file-systems/%v/restore", "64366a19-54e8-1544-f3d7-2a50fb1ccff3")

			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method != http.MethodPost {
					t.Fatal(fmt.Errorf("wrong method. Expected %s; but got %s", http.MethodGet, r.Method))
				}

				if r.URL.Path != href {
					t.Fatal(fmt.Errorf("wrong path. Expected %s; but got %s", href, r.URL.Path))
				}

				http.Error(w, "not found", http.StatusNotFound)
			}))
			return ts, check(hasError)
		},

		"operation-failed": func(t *testing.T) (*httptest.Server, []checkFn) {
			href := fmt.Sprintf("/rest/v1/file-systems/%v/restore", "64366a19-54e8-1544-f3d7-2a50fb1ccff3")

			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method != http.MethodPost {
					t.Fatal(fmt.Errorf("wrong method. Expected %s; but got %s", http.MethodGet, r.Method))
				}

				if r.URL.Path != href {
					t.Fatal(fmt.Errorf("wrong path. Expected %s; but got %s", href, r.URL.Path))
				}

				http.Error(w, "operation failed", http.StatusUnprocessableEntity)
			}))
			return ts, check(hasError)
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			ts, checkFns := tc(t)
			defer ts.Close()

			client, err := NewClientWithArgs(ts.URL, "", math.MaxInt64, true, false)
			client.configConnect.Version = "4.0"
			if err != nil {
				t.Fatal(err)
			}

			s := System{
				client: client,
			}

			fsID := "64366a19-54e8-1544-f3d7-2a50fb1ccff3"
			restoreSnapshotRequest := new(types.RestoreFsSnapParam)
			if name == "successWithContent" {
				restoreSnapshotRequest = &types.RestoreFsSnapParam{
					SnapshotID: "64366a19-54e8-1544-f3d7-2a50fb1ccdd3",
					CopyName:   "test-name",
				}
			} else {
				restoreSnapshotRequest = &types.RestoreFsSnapParam{
					SnapshotID: "64366a19-54e8-1544-f3d7-2a50fb1ccdd3",
				}
			}

			resp, err := s.RestoreFileSystemFromSnapshot(restoreSnapshotRequest, fsID)
			for _, checkFn := range checkFns {
				checkFn(t, resp, err)
			}
		})
	}
}

func TestDeleteFileSystem(t *testing.T) {
	name := "new-fs"

	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	}))
	defer svr.Close()

	client, err := NewClientWithArgs(svr.URL, "", math.MaxInt64, true, false)
	client.configConnect.Version = "4.0"
	if err != nil {
		t.Fatal(err)
	}

	s := System{
		client: client,
	}

	err = s.DeleteFileSystem(name)
	assert.NotNil(t, err)
}

func TestModifyFileSystem(t *testing.T) {
	type testCase struct {
		fsID        string
		newSize     int
		description string
		expected    error
	}
	cases := []testCase{
		{
			"64a6b2a4-1ff6-acaf-39a0-5643ff849351",
			10737418240,
			"",
			nil,
		},
		{
			"",
			10737418240,
			"",
			errors.New("file system name or ID is mandatory, please enter a valid value"),
		},
	}
	// mock a powerflex endpoint
	svr := httptest.NewServer(http.HandlerFunc(func(_ http.ResponseWriter, _ *http.Request) {
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

			fsParam := &types.FSModify{
				Size:        tc.newSize,
				Description: tc.description,
			}

			// calling ModifyFileSystem with mock value
			err = s.ModifyFileSystem(fsParam, tc.fsID)
			if err != nil {
				if tc.expected == nil {
					t.Errorf("Modifying FS did not work as expected, \n\tgot: %s \n\twant: %v", err, tc.expected)
				} else {
					if err.Error() != tc.expected.Error() {
						t.Errorf("Modifying FS did not work as expected, \n\tgot: %s \n\twant: %s", err, tc.expected)
					}
				}
			}
		})
	}
}
