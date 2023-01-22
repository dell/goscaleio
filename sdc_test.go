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
	"net/http"
	"net/http/httptest"
	"testing"

	types "github.com/dell/goscaleio/types/v1"
	"github.com/stretchr/testify/assert"
)

func Test_FindVolumes(t *testing.T) {
	type checkFn func(*testing.T, []*Volume, error)
	check := func(fns ...checkFn) []checkFn { return fns }

	hasNoError := func(t *testing.T, vols []*Volume, err error) {
		if err != nil {
			t.Fatalf("expected no error")
		}
	}

	checkLength := func(length int) func(t *testing.T, vols []*Volume, err error) {
		return func(t *testing.T, vols []*Volume, err error) {
			assert.Equal(t, length, len(vols))
		}
	}

	hasError := func(t *testing.T, vols []*Volume, err error) {
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

			client, err := NewClientWithArgs(ts.URL, "", true, false)
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

	//mock a powerflex endpoint
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	}))
	defer svr.Close()

	for _, tc := range cases {
		tc := tc
		t.Run("", func(ts *testing.T) {
			client, err := NewClientWithArgs(svr.URL, "", true, false)
			if err != nil {
				t.Fatal(err)
			}

			//calling RenameSdc with mock value
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
