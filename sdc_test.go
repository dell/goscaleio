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
	"errors"
	"fmt"
	"math"
	"net/http"
	"net/http/httptest"
	"testing"

	types "github.com/dell/goscaleio/types/v1"
	"github.com/stretchr/testify/assert"
)

func Test_FindVolumes(t *testing.T) {
	type checkFn func(*testing.T, []*Volume, error)
	check := func(fns ...checkFn) []checkFn { return fns }

	hasNoError := func(t *testing.T, _ []*Volume, err error) {
		if err != nil {
			t.Fatalf("expected no error")
		}
	}

	checkLength := func(length int) func(t *testing.T, vols []*Volume, err error) {
		return func(t *testing.T, vols []*Volume, _ error) {
			assert.Equal(t, length, len(vols))
		}
	}

	hasError := func(t *testing.T, _ []*Volume, err error) {
		if err == nil {
			t.Fatalf("expected error")
		}
	}

	tests := map[string]func(t *testing.T) (*httptest.Server, types.Sdc, []checkFn){
		"success": func(t *testing.T) (*httptest.Server, types.Sdc, []checkFn) {
			sdcID := "000001111a2222b"
			href := fmt.Sprintf("/api/instances/Sdc::%s/relationships/Volume", sdcID)
			sdc := types.Sdc{
				ID: sdcID,
				Links: []*types.Link{
					{
						Rel:  "/api/Sdc/relationship/Volume",
						HREF: href,
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
				vols := []types.Volume{{}, {}, {}}
				respData, err := json.Marshal(vols)
				if err != nil {
					t.Fatal(err)
				}
				fmt.Fprintln(w, string(respData))
			}))
			return ts, sdc, check(hasNoError, checkLength(3))
		},
		"error from GetVolume": func(t *testing.T) (*httptest.Server, types.Sdc, []checkFn) {
			sdcID := "someID"
			href := fmt.Sprintf("/api/instances/Sdc::%s/relationships/Volume", sdcID)
			sdc := types.Sdc{
				ID: sdcID,
				Links: []*types.Link{
					{
						Rel:  "/api/Sdc/relationship/Volume",
						HREF: href,
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
			return ts, sdc, check(hasError)
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			ts, sdc, checkFns := tc(t)
			defer ts.Close()

			client, err := NewClientWithArgs(ts.URL, "", math.MaxInt64, true, false)
			if err != nil {
				t.Fatal(err)
			}

			sdcClient := NewSdc(client, &sdc)
			vols, err := sdcClient.FindVolumes()
			for _, checkFn := range checkFns {
				checkFn(t, vols, err)
			}
		})
	}
}

func TestRenameSdc(t *testing.T) {
	type testCase struct {
		sdcID    string
		name     string
		expected error
	}
	cases := []testCase{
		{
			"c4270bf500000053",
			"worker-node-2345",
			nil,
		},
		{
			"c4270bf500000053",
			"",
			errors.New("Request message is not valid: The following parameter(s) must be part of the request body: sdcName"),
		},
		{
			"c4270bf500000053",
			" ",
			errors.New("The given name contains invalid characters. Use alphanumeric and punctuation characters only. Spaces are not allowed"),
		},
		{
			"worker-node-2",
			"c4270bf500000053",
			errors.New("id (worker-node-2) must be a hexadecimal number (unsigned long)"),
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

			// calling RenameSdc with mock value
			err = client.RenameSdc(tc.sdcID, tc.name)
			if err != nil {
				if tc.expected == nil {
					t.Errorf("Renaming sdc did not work as expected, \n\tgot: %s \n\twant: %v", err, tc.expected)
				} else {
					if err.Error() != tc.expected.Error() {
						t.Errorf("Renaming sdc did not work as expected, \n\tgot: %s \n\twant: %s", err, tc.expected)
					}
				}
			}
		})
	}
}

func TestApproveSdc(t *testing.T) {
	type checkFn func(*testing.T, *types.ApproveSdcByGUIDResponse, error)
	check := func(fns ...checkFn) []checkFn { return fns }

	hasNoError := func(t *testing.T, _ *types.ApproveSdcByGUIDResponse, err error) {
		if err != nil {
			t.Fatalf("expected no error")
		}
	}

	hasError := func(t *testing.T, _ *types.ApproveSdcByGUIDResponse, err error) {
		if err == nil {
			t.Fatalf("expected error")
		}
	}

	checkResp := func(sdcId string) func(t *testing.T, resp *types.ApproveSdcByGUIDResponse, err error) {
		return func(t *testing.T, resp *types.ApproveSdcByGUIDResponse, _ error) {
			assert.Equal(t, sdcId, resp.SdcID)
		}
	}

	tests := map[string]func(t *testing.T) (*httptest.Server, *types.System, []checkFn){
		"success": func(t *testing.T) (*httptest.Server, *types.System, []checkFn) {
			systemID := "0000aaabbbccc1111"
			href := fmt.Sprintf("/api/instances/System::%v/action/approveSdc", systemID)
			system := types.System{
				ID:                       systemID,
				RestrictedSdcModeEnabled: true,
				RestrictedSdcMode:        "Guid",
			}

			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method != http.MethodPost {
					t.Fatal(fmt.Errorf("wrong method. Expected %s; but got %s", http.MethodPost, r.Method))
				}

				if r.URL.Path != href {
					t.Fatal(fmt.Errorf("wrong path. Expected %s; but got %s", href, r.URL.Path))
				}

				resp := types.ApproveSdcByGUIDResponse{
					SdcID: "aab12340000000x",
				}

				respData, err := json.Marshal(resp)
				if err != nil {
					t.Fatal(err)
				}
				fmt.Fprintln(w, string(respData))
			}))
			return ts, &system, check(hasNoError, checkResp("aab12340000000x"))
		},
		"Already Approved err": func(t *testing.T) (*httptest.Server, *types.System, []checkFn) {
			systemID := "0000aaabbbccc1111"
			href := fmt.Sprintf("/api/instances/System::%v/action/approveSdc", systemID)
			system := types.System{
				ID:                       systemID,
				RestrictedSdcModeEnabled: true,
				RestrictedSdcMode:        "Guid",
			}

			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method != http.MethodPost {
					t.Fatal(fmt.Errorf("wrong method. Expected %s; but got %s", http.MethodPost, r.Method))
				}

				if r.URL.Path != href {
					t.Fatal(fmt.Errorf("wrong path. Expected %s; but got %s", href, r.URL.Path))
				}

				http.Error(w, "The SDC is already approved.", http.StatusInternalServerError)
			}))
			return ts, &system, check(hasError)
		},
		"Invalid guid err": func(t *testing.T) (*httptest.Server, *types.System, []checkFn) {
			systemID := "0000aaabbbccc1111"
			href := fmt.Sprintf("/api/instances/System::%v/action/approveSdc", systemID)
			system := types.System{
				ID:                       systemID,
				RestrictedSdcModeEnabled: true,
				RestrictedSdcMode:        "Guid",
			}

			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method != http.MethodPost {
					t.Fatal(fmt.Errorf("wrong method. Expected %s; but got %s", http.MethodPost, r.Method))
				}

				if r.URL.Path != href {
					t.Fatal(fmt.Errorf("wrong path. Expected %s; but got %s", href, r.URL.Path))
				}

				http.Error(w, "The given GUID is invalid. Please specify GUID in the following format: xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx", http.StatusInternalServerError)
			}))
			return ts, &system, check(hasError)
		},
	}

	testCaseGuids := map[string]string{
		"success":              "1aaabd94-9acd-11ed-a8fc-0242ac120002",
		"Already Approved err": "1aaabd94-9acd-11ed-a8fc-0242ac120002",
		"Invalid guid err":     "invald_guid",
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
				client: client,
				System: system,
			}

			resp, err := s.ApproveSdcByGUID(testCaseGuids[name])
			for _, checkFn := range checkFns {
				checkFn(t, resp, err)
			}
		})
	}
}

func TestSetRestrictedMode(t *testing.T) {
	type testCase struct {
		mode     string
		expected error
	}

	system := types.System{
		ID:                       "0000aaabbbccc1111",
		RestrictedSdcModeEnabled: true,
		RestrictedSdcMode:        "Guid",
	}

	cases := []testCase{
		{
			mode:     "Guid",
			expected: nil,
		},
		{
			mode:     "random",
			expected: errors.New("restrictedSdcMode should get one of the following values: None, Guid, ApprovedIp, but its value is random"),
		},
	}
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	}))
	defer svr.Close()

	for _, tc := range cases {
		tc := tc
		t.Run("", func(ts *testing.T) {
			client, err := NewClientWithArgs(svr.URL, "", math.MaxInt64, true, false)
			if err != nil {
				t.Fatal(err)
			}

			s := System{
				client: client,
				System: &system,
			}

			err2 := s.SetRestrictedMode(tc.mode)
			if err2 != nil {
				if tc.expected == nil {
					t.Errorf("Modifying restricted mode did not work as expected, \n\tgot: %s \n\twant: %v", err2, tc.expected)
				} else {
					if err2.Error() != tc.expected.Error() {
						t.Errorf("Modifying restricted mode did not work as expected, \n\tgot: %s \n\twant: %s", err2, tc.expected)
					}
				}
			}
		})
	}
}

func TestApproveSdcbyIP(t *testing.T) {
	type testCase struct {
		param    types.ApproveSdcParam
		expected error
	}

	system := types.System{
		ID:                       "0000aaabbbccc1111",
		RestrictedSdcModeEnabled: true,
		RestrictedSdcMode:        "Guid",
	}

	cases := []testCase{
		{
			param: types.ApproveSdcParam{
				SdcIP: "10.10.10.10",
			},
			expected: nil,
		},
		{
			param: types.ApproveSdcParam{
				Name: "sdc_test",
			},
			expected: errors.New("Request message is not valid: One of the parameter(s) must be part of the request body: sdcIp, sdcIps, sdcGuid"),
		},
	}
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	}))
	defer svr.Close()

	for _, tc := range cases {
		tc := tc
		t.Run("", func(ts *testing.T) {
			client, err := NewClientWithArgs(svr.URL, "", math.MaxInt64, true, false)
			if err != nil {
				t.Fatal(err)
			}

			s := System{
				client: client,
				System: &system,
			}

			_, err2 := s.ApproveSdc(&tc.param)
			if err2 != nil {
				if tc.expected == nil {
					t.Errorf("Approving SDC did not work as expected, \n\tgot: %s \n\twant: %v", err2, tc.expected)
				} else {
					if err2.Error() != tc.expected.Error() {
						t.Errorf("Approving SDC did not work as expected, \n\tgot: %s \n\twant: %s", err2, tc.expected)
					}
				}
			}
		})
	}

}

func TestSetApprovedIps(t *testing.T) {
	type testCase struct {
		SdcID    string
		SdcIps   []string
		expected error
	}

	system := types.System{
		ID:                       "0000aaabbbccc1111",
		RestrictedSdcModeEnabled: true,
		RestrictedSdcMode:        "Guid",
	}

	cases := []testCase{
		{
			SdcID:    "5e662b1700000000",
			SdcIps:   []string{"10.10.10.10"},
			expected: nil,
		},
		{
			SdcIps:   []string{"10.10.10.10"},
			expected: errors.New("Request message is not valid: The following parameter(s) must be part of the request body: sdcId"),
		},
	}
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	}))
	defer svr.Close()

	for _, tc := range cases {
		tc := tc
		t.Run("", func(ts *testing.T) {
			client, err := NewClientWithArgs(svr.URL, "", math.MaxInt64, true, false)
			if err != nil {
				t.Fatal(err)
			}

			s := System{
				client: client,
				System: &system,
			}

			err2 := s.SetApprovedIps(tc.SdcID, tc.SdcIps)
			if err2 != nil {
				if tc.expected == nil {
					t.Errorf("Approving SDC IPs did not work as expected, \n\tgot: %s \n\twant: %v", err2, tc.expected)
				} else {
					if err2.Error() != tc.expected.Error() {
						t.Errorf("Approving SDC IPs did not work as expected, \n\tgot: %s \n\twant: %s", err2, tc.expected)
					}
				}
			}
		})
	}

}
