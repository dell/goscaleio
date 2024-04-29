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
	"errors"
	"fmt"
	"math"
	"net/http"
	"net/http/httptest"
	"testing"

	types "github.com/dell/goscaleio/types/v1"
	"github.com/stretchr/testify/assert"
)

func TestTreeQuotaByID(t *testing.T) {
	type checkFn func(*testing.T, *types.TreeQuota, error)
	check := func(fns ...checkFn) []checkFn { return fns }

	hasNoError := func(t *testing.T, _ *types.TreeQuota, err error) {
		if err != nil {
			t.Fatalf("expected no error")
		}
	}

	hasError := func(t *testing.T, _ *types.TreeQuota, err error) {
		if err == nil {
			t.Fatalf("expected error")
		}
	}

	checkRespID := func(quotaID string) func(t *testing.T, resp *types.TreeQuota, err error) {
		return func(t *testing.T, resp *types.TreeQuota, _ error) {
			assert.Equal(t, quotaID, resp.ID)
		}
	}

	testsID := map[string]func(t *testing.T) (*httptest.Server, []checkFn){
		"success": func(t *testing.T) (*httptest.Server, []checkFn) {
			quotaID := "00000003-006a-0000-0600-000000000000"
			href := fmt.Sprintf("/rest/v1/file-tree-quotas/%s", quotaID)

			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method != http.MethodGet {
					t.Fatal(fmt.Errorf("wrong method. Expected %s; but got %s", http.MethodGet, r.Method))
				}

				if r.URL.Path != href {
					t.Fatal(fmt.Errorf("wrong path. Expected %s; but got %s", href, r.URL.Path))
				}

				resp := types.TreeQuota{
					ID: "00000003-006a-0000-0600-000000000000",
				}

				respData, err := json.Marshal(resp)
				if err != nil {
					t.Fatal(err)
				}
				fmt.Fprintln(w, string(respData))
			}))
			return ts, check(hasNoError, checkRespID("00000003-006a-0000-0600-000000000000"))
		},
		"not found": func(t *testing.T) (*httptest.Server, []checkFn) {
			quotaID := "00000003-006a-0000-0700-000000000000"
			href := fmt.Sprintf("/rest/v1/file-tree-quotas/%s", quotaID)

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

	testCaseIDs := map[string]string{
		"success":   "00000003-006a-0000-0600-000000000000",
		"not found": "00000003-006a-0000-0700-000000000000",
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

			resp, err := s.GetTreeQuotaByID(testCaseIDs[id])
			for _, checkFn := range checkFns {
				checkFn(t, resp, err)
			}
		})
	}
}

func TestCreateTreeQuota(t *testing.T) {
	type checkFn func(*testing.T, *types.TreeQuotaCreateResponse, error)
	check := func(fns ...checkFn) []checkFn { return fns }

	hasNoError := func(t *testing.T, _ *types.TreeQuotaCreateResponse, err error) {
		if err != nil {
			t.Fatalf("expected no error")
		}
	}

	hasError := func(t *testing.T, _ *types.TreeQuotaCreateResponse, err error) {
		if err == nil {
			t.Fatalf("expected error")
		}
	}

	checkResp := func(quotaId string) func(t *testing.T, resp *types.TreeQuotaCreateResponse, err error) {
		return func(t *testing.T, resp *types.TreeQuotaCreateResponse, _ error) {
			assert.Equal(t, quotaId, resp.ID)
		}
	}

	tests := map[string]func(t *testing.T) (*httptest.Server, []checkFn){
		"success": func(t *testing.T) (*httptest.Server, []checkFn) {
			href := fmt.Sprintf("/rest/v1/file-tree-quotas")

			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method != http.MethodPost {
					t.Fatal(fmt.Errorf("wrong method. Expected %s; but got %s", http.MethodGet, r.Method))
				}

				if r.URL.Path != href {
					t.Fatal(fmt.Errorf("wrong path. Expected %s; but got %s", href, r.URL.Path))
				}

				resp := types.TreeQuotaCreateResponse{
					ID: "00000003-006a-0000-0600-000000000000",
				}

				respData, err := json.Marshal(resp)
				if err != nil {
					t.Fatal(err)
				}
				fmt.Fprintln(w, string(respData))
			}))
			return ts, check(hasNoError, checkResp("00000003-006a-0000-0600-000000000000"))
		},
		"bad request": func(t *testing.T) (*httptest.Server, []checkFn) {
			href := fmt.Sprintf("/rest/v1/file-tree-quotas")

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

			resp, err := s.CreateTreeQuota(&types.TreeQuotaCreate{
				FileSystemID: "64b3ceca-046f-eb3a-da83-3a7645b0a943",
				Path:         "/fs111",
			})
			for _, checkFn := range checkFns {
				checkFn(t, resp, err)
			}
		})
	}
}

func TestDeleteTreeQuota(t *testing.T) {
	id := "00000003-006a-0000-0600-000000000000"

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
	}

	err = s.DeleteTreeQuota(id)
	assert.Nil(t, err)
}

func TestModifyTreeQuota(t *testing.T) {
	type testCase struct {
		QuotaID     string
		SoftLimit   int
		description string
		expected    error
	}
	cases := []testCase{
		{
			"00000003-006a-0000-0600-000000000000",
			1000,
			"",
			nil,
		},
		{
			"",
			1100,
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

			quotaParam := &types.TreeQuotaModify{
				SoftLimit:   tc.SoftLimit,
				Description: tc.description,
			}

			// calling ModifyTreeQuota with mock value
			err = s.ModifyTreeQuota(quotaParam, tc.QuotaID)
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
