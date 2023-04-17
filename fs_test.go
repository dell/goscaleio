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
	"fmt"
	"math"
	"net/http"
	"net/http/httptest"
	"testing"

	types "github.com/dell/goscaleio/types/v1"
	"github.com/stretchr/testify/assert"
)



func TestGetFileSystemByID(t *testing.T) {
	type checkFn func(*testing.T, *types.FileSystem, error)
	check := func(fns ...checkFn) []checkFn { return fns }

	hasNoError := func(t *testing.T, resp *types.FileSystem, err error) {
		if err != nil {
			t.Fatalf("expected no error")
		}
	}

	hasError := func(t *testing.T, resp *types.FileSystem, err error) {
		if err == nil {
			t.Fatalf("expected error")
		}
	}

	checkResp := func(fsID string) func(t *testing.T, resp *types.FileSystem, err error) {
		return func(t *testing.T, resp *types.FileSystem, err error) {
			assert.Equal(t, fsID, resp.ID)
		}
	}

	tests := map[string]func(t *testing.T) (*httptest.Server, []checkFn){
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
						ID:		"64366a19-54e8-1544-f3d7-2a50fb1ccff3",
						Name:	"fs-test-1",
				}

				respData, err := json.Marshal(resp)
				if err != nil {
					t.Fatal(err)
				}
				fmt.Fprintln(w, string(respData))
			}))
			return ts, check(hasNoError, checkResp("64366a19-54e8-1544-f3d7-2a50fb1ccff3"))
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

	var testCaseFSIds = map[string]string{
		"success":   "64366a19-54e8-1544-f3d7-2a50fb1ccff3",
		"not found": "6436aa58-e6a1-a4e2-de7b-2a50fb1ccff3",
	}

	for id, tc := range tests {
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

			resp, err := s.GetFileSystemByID(testCaseFSIds[id])
			for _, checkFn := range checkFns {
				checkFn(t, resp, err)
			}
		})
	}
}

func TestGetFileSystemByName(t *testing.T) {
	type checkFn func(*testing.T, *types.FileSystem, error)
	check := func(fns ...checkFn) []checkFn { return fns }

	hasNoError := func(t *testing.T, resp *types.FileSystem, err error) {
		if err != nil {
			t.Fatalf("expected no error")
		}
	}

	hasError := func(t *testing.T, resp *types.FileSystem, err error) {
		if err == nil {
			t.Fatalf("expected error")
		}
	}

	checkResp := func(fsName string) func(t *testing.T, resp *types.FileSystem, err error) {
		return func(t *testing.T, resp *types.FileSystem, err error) {
			assert.Equal(t, fsName, resp.Name)
		}
	}

	tests := map[string]func(t *testing.T) (*httptest.Server, []checkFn){
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
						ID:		"64366a19-54e8-1544-f3d7-2a50fb1ccff3",
						Name:	"fs-test-1",
					},
					{
						ID:		"6436aa58-e6a1-a4e2-de7b-2a50fb1ccff3",
						Name:	"fs-test-2",
					},
				}

				respData, err := json.Marshal(resp)
				if err != nil {
					t.Fatal(err)
				}
				fmt.Fprintln(w, string(respData))
			}))
			return ts, check(hasNoError, checkResp("fs-test-2"))
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
						ID:		"64366a19-54e8-1544-f3d7-2a50fb1ccff3",
						Name:	"fs-test-1",
					},
					{
						ID:		"6436aa58-e6a1-a4e2-de7b-2a50fb1ccff3",
						Name:	"fs-test-2",
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

	var testCaseFSNames = map[string]string{
		"success":   "fs-test-2",
		"not found": "fs-test-3",
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

			resp, err := s.GetFileSystemByName(testCaseFSNames[name])
			for _, checkFn := range checkFns {
				checkFn(t, resp, err)
			}
		})
	}
}

func TestCreateFileSystem(t *testing.T) {
	type checkFn func(*testing.T, *types.FileSystemResp, error)
	check := func(fns ...checkFn) []checkFn { return fns }

	hasNoError := func(t *testing.T, resp *types.FileSystemResp, err error) {
		if err != nil {
			t.Fatalf("expected no error")
		}
	}

	hasError := func(t *testing.T, resp *types.FileSystemResp, err error) {
		if err == nil {
			t.Fatalf("expected error")
		}
	}

	checkResp := func(fsId string) func(t *testing.T, resp *types.FileSystemResp, err error) {
		return func(t *testing.T, resp *types.FileSystemResp, err error) {
			assert.Equal(t, fsId, resp.ID)
		}
	}

	tests := map[string]func(t *testing.T) (*httptest.Server, []checkFn){
		"success": func(t *testing.T) (*httptest.Server, []checkFn) {

			href := fmt.Sprintf("/rest/v1/file-systems")

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
			href := fmt.Sprintf("/rest/v1/file-systems")

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
