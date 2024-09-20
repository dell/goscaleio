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

func TestGetNasByIDName(t *testing.T) {
	type checkFn func(*testing.T, *types.NAS, error)
	check := func(fns ...checkFn) []checkFn { return fns }

	hasNoError := func(t *testing.T, _ *types.NAS, err error) {
		if err != nil {
			t.Fatalf("expected no error")
		}
	}

	hasError := func(t *testing.T, _ *types.NAS, err error) {
		if err == nil {
			t.Fatalf("expected error")
		}
	}

	checkRespName := func(nasName string) func(t *testing.T, resp *types.NAS, err error) {
		return func(t *testing.T, resp *types.NAS, _ error) {
			assert.Equal(t, nasName, resp.Name)
		}
	}

	checkRespID := func(nasId string) func(t *testing.T, resp *types.NAS, err error) {
		return func(t *testing.T, resp *types.NAS, _ error) {
			assert.Equal(t, nasId, resp.ID)
		}
	}

	testsName := map[string]func(t *testing.T) (*httptest.Server, *types.System, []checkFn){
		"success": func(t *testing.T) (*httptest.Server, *types.System, []checkFn) {
			systemID := "0000aaacccddd1111"
			href := "/rest/v1/nas-servers"
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

				resp := []types.NAS{
					{
						ID:                 "5e8d8e8e-671b-336f-db4e-cee0fbdc981e",
						Name:               "test-nas1",
						ProtectionDomainID: "test-pd",
					},
					{
						ID:                 "6e8d8e8e-671b-336f-eb4e-dee0fbdc981e",
						Name:               "test-nas2",
						ProtectionDomainID: "test-pd",
					},
				}

				respData, err := json.Marshal(resp)
				if err != nil {
					t.Fatal(err)
				}
				fmt.Fprintln(w, string(respData))
			}))
			return ts, &system, check(hasNoError, checkRespName("test-nas1"))
		},
		"not found": func(t *testing.T) (*httptest.Server, *types.System, []checkFn) {
			systemID := "0000aaacccddd1111"
			href := "/rest/v1/nas-servers"
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

				resp := []types.NAS{
					{
						ID:                 "5e8d8e8e-671b-336f-db4e-cee0fbdc981e",
						Name:               "test-nas1",
						ProtectionDomainID: "test-pd",
					},
					{
						ID:                 "6e8d8e8e-671b-336f-eb4e-dee0fbdc981e",
						Name:               "test-nas2",
						ProtectionDomainID: "test-pd",
					},
				}

				respData, err := json.Marshal(resp)
				if err != nil {
					t.Fatal(err)
				}
				fmt.Fprintln(w, string(respData))
			}))
			return ts, &system, check(hasError)
		},
	}

	testsID := map[string]func(t *testing.T) (*httptest.Server, *types.System, []checkFn){
		"success": func(t *testing.T) (*httptest.Server, *types.System, []checkFn) {
			nasID := "5e8d8e8e-671b-336f-db4e-cee0fbdc981e"
			systemID := "0000aaacccddd1111"
			href := fmt.Sprintf("/rest/v1/nas-servers/%s", nasID)
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

				resp := types.NAS{
					ID:                 "5e8d8e8e-671b-336f-db4e-cee0fbdc981e",
					Name:               "test-nas",
					ProtectionDomainID: "test-pd",
				}

				respData, err := json.Marshal(resp)
				if err != nil {
					t.Fatal(err)
				}
				fmt.Fprintln(w, string(respData))
			}))
			return ts, &system, check(hasNoError, checkRespID("5e8d8e8e-671b-336f-db4e-cee0fbdc981e"))
		},
		"not found": func(t *testing.T) (*httptest.Server, *types.System, []checkFn) {
			nasID := "6e8d8e8e-671b-336f-eb4e-dee0fbdc981f"
			systemID := "0000aaacccddd1111"
			href := fmt.Sprintf("/rest/v1/nas-servers/%s", nasID)
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

				http.Error(w, "nas not found", http.StatusNotFound)
			}))
			return ts, &system, check(hasError)
		},
	}

	testCaseNasNames := map[string]string{
		"success":   "test-nas1",
		"not found": "test-nas3",
	}

	testCaseNasIDs := map[string]string{
		"success":   "5e8d8e8e-671b-336f-db4e-cee0fbdc981e",
		"not found": "6e8d8e8e-671b-336f-eb4e-dee0fbdc981f",
	}

	for name, tc := range testsName {
		t.Run(name, func(t *testing.T) {
			ts, system, checkFns := tc(t)
			defer ts.Close()

			client, err := NewClientWithArgs(ts.URL, "", math.MaxInt64, true, false)
			client.configConnect.Version = "4.0"
			if err != nil {
				t.Fatal(err)
			}

			s := System{
				client: client,
				System: system,
			}

			resp, err := s.GetNASByIDName("", testCaseNasNames[name])
			for _, checkFn := range checkFns {
				checkFn(t, resp, err)
			}
		})
	}

	for name, tc := range testsID {
		t.Run(name, func(t *testing.T) {
			ts, system, checkFns := tc(t)
			defer ts.Close()

			client, err := NewClientWithArgs(ts.URL, "", math.MaxInt64, true, false)
			client.configConnect.Version = "4.0"
			if err != nil {
				t.Fatal(err)
			}

			s := System{
				client: client,
				System: system,
			}

			resp, err := s.GetNASByIDName(testCaseNasIDs[name], "")
			for _, checkFn := range checkFns {
				checkFn(t, resp, err)
			}
		})
	}
}

func TestCreateNAS(t *testing.T) {
	type checkFn func(*testing.T, *types.CreateNASResponse, error)
	check := func(fns ...checkFn) []checkFn { return fns }

	hasNoError := func(t *testing.T, _ *types.CreateNASResponse, err error) {
		if err != nil {
			t.Fatalf("expected no error")
		}
	}

	hasError := func(t *testing.T, _ *types.CreateNASResponse, err error) {
		if err == nil {
			t.Fatalf("expected error")
		}
	}

	checkResp := func(nasId string) func(t *testing.T, resp *types.CreateNASResponse, err error) {
		return func(t *testing.T, resp *types.CreateNASResponse, _ error) {
			assert.Equal(t, nasId, resp.ID)
		}
	}

	tests := map[string]func(t *testing.T) (*httptest.Server, *types.System, []checkFn){
		"success": func(t *testing.T) (*httptest.Server, *types.System, []checkFn) {
			systemID := "0000aaacccddd1111"
			href := "/rest/v1/nas-servers"
			system := types.System{
				ID: systemID,
			}

			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method != http.MethodPost {
					t.Fatal(fmt.Errorf("wrong method. Expected %s; but got %s", http.MethodGet, r.Method))
				}

				if r.URL.Path != href {
					t.Fatal(fmt.Errorf("wrong path. Expected %s; but got %s", href, r.URL.Path))
				}

				resp := types.CreateNASResponse{
					ID: "5e8d8e8e-671b-336f-db4e-cee0fbdc981e",
				}

				respData, err := json.Marshal(resp)
				if err != nil {
					t.Fatal(err)
				}
				fmt.Fprintln(w, string(respData))
			}))
			return ts, &system, check(hasNoError, checkResp("5e8d8e8e-671b-336f-db4e-cee0fbdc981e"))
		},
		"bad request": func(t *testing.T) (*httptest.Server, *types.System, []checkFn) {
			href := "/rest/v1/nas-servers"
			systemID := "0000aaacccddd1111"
			system := types.System{
				ID: systemID,
			}

			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method != http.MethodPost {
					t.Fatal(fmt.Errorf("wrong method. Expected %s; but got %s", http.MethodGet, r.Method))
				}

				if r.URL.Path != href {
					t.Fatal(fmt.Errorf("wrong path. Expected %s; but got %s", href, r.URL.Path))
				}

				http.Error(w, "bad Request", http.StatusBadRequest)
			}))
			return ts, &system, check(hasError)
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			ts, system, checkFns := tc(t)
			defer ts.Close()

			client, err := NewClientWithArgs(ts.URL, "", math.MaxInt64, true, false)
			client.configConnect.Version = "4.0"
			if err != nil {
				t.Fatal(err)
			}

			s := System{
				client: client,
				System: system,
			}

			resp, err := s.CreateNAS("test-nas1", "pd1")
			for _, checkFn := range checkFns {
				checkFn(t, resp, err)
			}
		})
	}
}

func TestDeleteNAS(t *testing.T) {
	id := "new-nas"
	systemID := "0000aaacccddd1111"
	system := types.System{
		ID: systemID,
	}
	system1 := &system

	// mock a powerflex endpoint
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
		System: system1,
	}

	err = s.DeleteNAS(id)
	assert.Nil(t, err)
}

func TestPingNAS(t *testing.T) {
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

	testsName := map[string]func(t *testing.T) (*httptest.Server, *types.System, []checkFn){
		"success": func(t *testing.T) (*httptest.Server, *types.System, []checkFn) {
			systemID := "0000aaacccddd1111"
			href := "/rest/v1/nas-servers/655374ea-13d7-c2d5-458c-4ec4ea9bb086/ping"
			system := types.System{
				ID: systemID,
			}

			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method != http.MethodPost {
					t.Fatal(fmt.Errorf("wrong method. Expected %s; but got %s", http.MethodGet, r.Method))
				}

				if r.URL.Path != href {
					t.Fatal(fmt.Errorf("wrong path. Expected %s; but got %s", href, r.URL.Path))
				}
			}))
			return ts, &system, check(hasNoError)
		},
		"failure": func(t *testing.T) (*httptest.Server, *types.System, []checkFn) {
			systemID := "0000aaacccddd1111"
			href := "/rest/v1/nas-servers/6e8d8e8e-671b-336f-eb4e-dee0fbdc981f/ping"
			system := types.System{
				ID: systemID,
			}

			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method != http.MethodPost {
					t.Fatal(fmt.Errorf("wrong method. Expected %s; but got %s", http.MethodPost, r.Method))
				}

				if r.URL.Path != href {
					t.Fatal(fmt.Errorf("wrong path. Expected %s; but got %s", href, r.URL.Path))
				}

				http.Error(w, "could not ping NAS server", http.StatusNotFound)
			}))
			return ts, &system, check(hasError)
		},
	}

	testCaseNasServers := map[string][]string{
		"success": {"655374ea-13d7-c2d5-458c-4ec4ea9bb086", "10.20.30.40"},
		"failure": {"6e8d8e8e-671b-336f-eb4e-dee0fbdc981f", "11.22.33.44"},
	}

	for name, tc := range testsName {
		t.Run(name, func(t *testing.T) {
			ts, system, checkFns := tc(t)
			defer ts.Close()

			client, err := NewClientWithArgs(ts.URL, "", math.MaxInt64, true, false)
			client.configConnect.Version = "4.0"
			if err != nil {
				t.Fatal(err)
			}

			s := System{
				client: client,
				System: system,
			}

			err = s.PingNAS(testCaseNasServers[name][0], testCaseNasServers[name][1])
			for _, checkFn := range checkFns {
				checkFn(t, err)
			}
		})
	}
}

func TestGeFileInterfaace(t *testing.T) {
	type checkFn func(*testing.T, *types.FileInterface, error)
	check := func(fns ...checkFn) []checkFn { return fns }

	hasNoError := func(t *testing.T, _ *types.FileInterface, err error) {
		if err != nil {
			t.Fatalf("expected no error")
		}
	}

	hasError := func(t *testing.T, _ *types.FileInterface, err error) {
		if err == nil {
			t.Fatalf("expected error")
		}
	}

	checkRespID := func(fileInterface string) func(t *testing.T, resp *types.FileInterface, err error) {
		return func(t *testing.T, resp *types.FileInterface, _ error) {
			assert.Equal(t, fileInterface, resp.ID)
		}
	}

	testsID := map[string]func(t *testing.T) (*httptest.Server, *types.System, []checkFn){
		"success": func(t *testing.T) (*httptest.Server, *types.System, []checkFn) {
			systemID := "0000aaacccddd1111"
			fileInterfaceID := "5e8d8e8e-671b-336f-db4e-cee0fbdc981e"
			href := fmt.Sprintf("/rest/v1/file-interfaces/%s", fileInterfaceID)
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

				resp := types.FileInterface{
					ID:        "5e8d8e8e-671b-336f-db4e-cee0fbdc981e",
					IPAddress: "10.20.30.40",
				}

				respData, err := json.Marshal(resp)
				if err != nil {
					t.Fatal(err)
				}
				fmt.Fprintln(w, string(respData))
			}))
			return ts, &system, check(hasNoError, checkRespID("5e8d8e8e-671b-336f-db4e-cee0fbdc981e"))
		},
		"not found": func(t *testing.T) (*httptest.Server, *types.System, []checkFn) {
			systemID := "0000aaacccddd1111"
			fileInterfaceId := "6e8d8e8e-671b-336f-eb4e-dee0fbdc981f"
			href := fmt.Sprintf("/rest/v1/file-interfaces/%s", fileInterfaceId)
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

				http.Error(w, "could not find the File interface using id", http.StatusNotFound)
			}))
			return ts, &system, check(hasError)
		},
	}

	testCaseFileInterfaceIds := map[string]string{
		"success":   "5e8d8e8e-671b-336f-db4e-cee0fbdc981e",
		"not found": "6e8d8e8e-671b-336f-eb4e-dee0fbdc981f",
	}

	for name, tc := range testsID {
		t.Run(name, func(t *testing.T) {
			ts, system, checkFns := tc(t)
			defer ts.Close()

			client, err := NewClientWithArgs(ts.URL, "", math.MaxInt64, true, false)
			client.configConnect.Version = "4.0"
			if err != nil {
				t.Fatal(err)
			}

			s := System{
				client: client,
				System: system,
			}

			resp, err := s.GetFileInterface(testCaseFileInterfaceIds[name])
			for _, checkFn := range checkFns {
				checkFn(t, resp, err)
			}
		})
	}
}
