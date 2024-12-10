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
	"context"
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"net/http/httptest"
	"testing"

	types "github.com/dell/goscaleio/types/v1"
	"github.com/stretchr/testify/assert"
)

func TestGetSystemLimits(t *testing.T) {
	type checkFn func(*testing.T, *types.QuerySystemLimitsResponse, error)
	check := func(fns ...checkFn) []checkFn { return fns }

	hasNoError := func(t *testing.T, _ *types.QuerySystemLimitsResponse, err error) {
		if err != nil {
			t.Fatalf("expected no error")
		}
	}

	hasError := func(t *testing.T, _ *types.QuerySystemLimitsResponse, err error) {
		if err == nil {
			t.Fatalf("expected error")
		}
	}

	checkLimitType := func(expectedType string) func(t *testing.T, syslimit *types.QuerySystemLimitsResponse, err error) {
		return func(t *testing.T, syslimit *types.QuerySystemLimitsResponse, err error) {
			if err == nil {
				// Add your custom assertions here to check the syslimit.Type.
				assert.Equal(t, expectedType, syslimit.SystemLimitEntryList[0].Type)
			}
		}
	}

	checkLimitDescription := func(expectedDescription string) func(t *testing.T, syslimit *types.QuerySystemLimitsResponse, err error) {
		return func(t *testing.T, syslimit *types.QuerySystemLimitsResponse, err error) {
			if err == nil {
				// Add your custom assertions here to check the syslimit.Description.
				assert.Equal(t, expectedDescription, syslimit.SystemLimitEntryList[0].Description)
			}
		}
	}

	checkLimitMaxVal := func(expectedMaxVal string) func(t *testing.T, syslimit *types.QuerySystemLimitsResponse, err error) {
		return func(t *testing.T, syslimit *types.QuerySystemLimitsResponse, err error) {
			if err == nil {
				// Add your custom assertions here to check the syslimit.MaxVal.
				assert.Equal(t, expectedMaxVal, syslimit.SystemLimitEntryList[0].MaxVal)
			}
		}
	}

	tests := map[string]func(t *testing.T) (*httptest.Server, []checkFn){
		"success": func(t *testing.T) (*httptest.Server, []checkFn) {
			href := "/api/instances/System/action/querySystemLimits"

			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method != http.MethodPost {
					t.Fatal(fmt.Errorf("wrong method. Expected %s; but got %s", http.MethodPost, r.Method))
				}

				if r.URL.Path != href {
					t.Fatal(fmt.Errorf("wrong path. Expected %s; but got %s", href, r.URL.Path))
				}

				// Simulate a successful response for GetSystemLimits.
				resp := types.QuerySystemLimitsResponse{
					SystemLimitEntryList: []types.SystemLimits{
						{
							Type:        "volumeSizeGb",
							Description: "Maximum volume size in GB",
							MaxVal:      "1024",
						},
					},
				}

				respData, err := json.Marshal(resp)
				if err != nil {
					t.Fatal(err)
				}
				fmt.Fprintln(w, string(respData))
			}))

			return ts, check(hasNoError, checkLimitType("volumeSizeGb"), checkLimitDescription("Maximum volume size in GB"), checkLimitMaxVal("1024"))
		},
		"not found": func(t *testing.T) (*httptest.Server, []checkFn) {
			href := "/api/instances/System/action/querySystemLimits"

			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method != http.MethodPost {
					t.Fatal(fmt.Errorf("wrong method. Expected %s; but got %s", http.MethodPost, r.Method))
				}

				if r.URL.Path != href {
					t.Fatal(fmt.Errorf("wrong path. Expected %s; but got %s", href, r.URL.Path))
				}

				http.Error(w, "nas not found", http.StatusNotFound)
			}))

			return ts, check(hasError)
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			ts, checkFns := tc(t)
			defer ts.Close()

			// Create a test client and call GetSystemLimits.
			// client := NewTestClient(ts.URL) // Replace with your own client creation logic.
			client, err := NewClientWithArgs(ts.URL, "", math.MaxInt64, true, false)
			client.configConnect.Version = "4.0"
			if err != nil {
				t.Fatal(err)
			}

			sys := System{
				client: client,
			}

			resp, err := sys.client.GetSystemLimits(context.Background())
			for _, checkFn := range checkFns {
				checkFn(t, resp, err)
			}
		})
	}
}
