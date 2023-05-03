/*
 *
 * Copyright © 2020-2023 Dell Inc. or its subsidiaries. All Rights Reserved.
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

func TestGetNFSExportByIDName(t *testing.T) {
	type checkFn func(*testing.T, *types.NFSExport, error)
	check := func(fns ...checkFn) []checkFn { return fns }

	hasNoError := func(t *testing.T, resp *types.NFSExport, err error) {
		if err != nil {
			t.Fatalf("expected no error")
		}
	}

	hasError := func(t *testing.T, resp *types.NFSExport, err error) {
		if err == nil {
			t.Fatalf("expected error")
		}
	}

	checkRespName := func(nfsName string) func(t *testing.T, resp *types.NFSExport, err error) {
		return func(t *testing.T, resp *types.NFSExport, err error) {
			assert.Equal(t, nfsName, resp.Name)
		}
	}

	checkRespID := func(nfsID string) func(t *testing.T, resp *types.NFSExport, err error) {
		return func(t *testing.T, resp *types.NFSExport, err error) {
			assert.Equal(t, nfsID, resp.ID)
		}
	}

	testsName := map[string]func(t *testing.T) (*httptest.Server, []checkFn){
		"success": func(t *testing.T) (*httptest.Server, []checkFn) {
			href := "/rest/v1/nfs-exports"

			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method != http.MethodGet {
					t.Fatal(fmt.Errorf("wrong method. Expected %s; but got %s", http.MethodGet, r.Method))
				}

				if r.URL.Path != href {
					t.Fatal(fmt.Errorf("wrong path. Expected %s; but got %s", href, r.URL.Path))
				}

				resp := []types.NFSExport{
					{
						ID:   "64366a19-54e8-1544-f3d7-2a50fb1ccff3",
						Name: "nfs-test-1",
					},
					{
						ID:   "6436aa58-e6a1-a4e2-de7b-2a50fb1ccff3",
						Name: "nfs-test-2",
					},
				}

				respData, err := json.Marshal(resp)
				if err != nil {
					t.Fatal(err)
				}
				fmt.Fprintln(w, string(respData))
			}))
			return ts, check(hasNoError, checkRespName("nfs-test-2"))
		},
		"not found": func(t *testing.T) (*httptest.Server, []checkFn) {
			href := "/rest/v1/nfs-exports"

			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method != http.MethodGet {
					t.Fatal(fmt.Errorf("wrong method. Expected %s; but got %s", http.MethodGet, r.Method))
				}

				if r.URL.Path != href {
					t.Fatal(fmt.Errorf("wrong path. Expected %s; but got %s", href, r.URL.Path))
				}

				resp := []types.NFSExport{
					{
						ID:   "64366a19-54e8-1544-f3d7-2a50fb1ccff3",
						Name: "nfs-test-1",
					},
					{
						ID:   "6436aa58-e6a1-a4e2-de7b-2a50fb1ccff3",
						Name: "nfs-test-2",
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
			nfsID := "64242c06-7a78-1773-50f4-2a50fb1ccff3"
			href := fmt.Sprintf("/rest/v1/nfs-exports/%s", nfsID)

			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method != http.MethodGet {
					t.Fatal(fmt.Errorf("wrong method. Expected %s; but got %s", http.MethodGet, r.Method))
				}

				if r.URL.Path != href {
					t.Fatal(fmt.Errorf("wrong path. Expected %s; but got %s", href, r.URL.Path))
				}

				resp := types.NFSExport{
					ID:           "64242c06-7a78-1773-50f4-2a50fb1ccff3",
					Name:         "nfs-test-1",
					FileSystemID: "64242bfb-d188-e87b-c144-2a50fb1ccff3",
				}

				respData, err := json.Marshal(resp)
				if err != nil {
					t.Fatal(err)
				}
				fmt.Fprintln(w, string(respData))
			}))
			return ts, check(hasNoError, checkRespID("64242c06-7a78-1773-50f4-2a50fb1ccff3"))
		},
		"not found": func(t *testing.T) (*httptest.Server, []checkFn) {
			nfsID := "6433a2b2-6d60-f737-9f3b-2a50fb1ccff3"
			href := fmt.Sprintf("/rest/v1/nfs-exports/%s", nfsID)

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

	var testCaseFSNames = map[string]string{
		"success":   "nfs-test-2",
		"not found": "nfs-test-3",
	}

	var testCaseFSIds = map[string]string{
		"success":   "64242c06-7a78-1773-50f4-2a50fb1ccff3",
		"not found": "6433a2b2-6d60-f737-9f3b-2a50fb1ccff3",
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

			resp, err := s.client.GetNFSExportByIDName("", testCaseFSNames[name])
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

			resp, err := s.client.GetNFSExportByIDName(testCaseFSIds[id], "")
			for _, checkFn := range checkFns {
				checkFn(t, resp, err)
			}
		})
	}
}

func TestDeleteNFSExport(t *testing.T) {

	id := "new-nas"

	//mock a powerflex endpoint
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	}))
	defer svr.Close()

	client, err := NewClientWithArgs(svr.URL, "", math.MaxInt64, true, false)
	client.configConnect.Version = "4.0"
	if err != nil {
		t.Fatal(err)
	}

	err = client.DeleteNFSExport(id)
	assert.Nil(t, err)
}

func TestCreateNFSExport(t *testing.T) {
	type checkFn func(*testing.T, *types.NFSExportCreateResponse, error)
	check := func(fns ...checkFn) []checkFn { return fns }

	hasNoError := func(t *testing.T, resp *types.NFSExportCreateResponse, err error) {
		if err != nil {
			t.Fatalf("expected no error")
		}
	}

	hasError := func(t *testing.T, resp *types.NFSExportCreateResponse, err error) {
		if err == nil {
			t.Fatalf("expected error")
		}
	}

	checkResp := func(nfsId string) func(t *testing.T, resp *types.NFSExportCreateResponse, err error) {
		return func(t *testing.T, resp *types.NFSExportCreateResponse, err error) {
			assert.Equal(t, nfsId, resp.ID)
		}
	}

	tests := map[string]func(t *testing.T) (*httptest.Server, []checkFn){
		"success": func(t *testing.T) (*httptest.Server, []checkFn) {

			href := fmt.Sprintf("/rest/v1/nfs-exports")

			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method != http.MethodPost {
					t.Fatal(fmt.Errorf("wrong method. Expected %s; but got %s", http.MethodGet, r.Method))
				}

				if r.URL.Path != href {
					t.Fatal(fmt.Errorf("wrong path. Expected %s; but got %s", href, r.URL.Path))
				}

				resp := types.NFSExportCreateResponse{
					ID: "64385158-97a1-bb86-4fd9-2a50fb1ccff3",
				}

				respData, err := json.Marshal(resp)
				if err != nil {
					t.Fatal(err)
				}
				fmt.Fprintln(w, string(respData))
			}))
			return ts, check(hasNoError, checkResp("64385158-97a1-bb86-4fd9-2a50fb1ccff3"))
		},
		"bad request": func(t *testing.T) (*httptest.Server, []checkFn) {
			href := fmt.Sprintf("/rest/v1/nfs-exports")

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

			resp, err := client.CreateNFSExport(&types.NFSExportCreate{
				Name:         "twee-kk",
				FileSystemID: "6436aa58-e6a1-a4e2-de7b-2a50fb1ccff3",
				Path:         "/twee-fs11",
			})
			for _, checkFn := range checkFns {
				checkFn(t, resp, err)
			}

		})
	}
}
