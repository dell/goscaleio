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
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestMoveToNextPhase(t *testing.T) {
	type testCase struct {
		expected error
	}

	cases := []testCase{
		{
			nil,
		},
	}

	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	}))
	defer svr.Close()

	for _, tc := range cases {
		tc := tc
		t.Run("", func(ts *testing.T) {
			GC, err := NewGateway(svr.URL, "", "", true, true)
			if err != nil {
				t.Fatal(err)
			}

			_, err = GC.MoveToNextPhase()
			if err != nil {
				if tc.expected == nil {
					t.Errorf("Move to Next Phase did not work as expected, \n\tgot: %s \n\twant: %v", err, tc.expected)
				} else {
					if err.Error() != tc.expected.Error() {
						t.Errorf("Move to Next Phase did not work as expected, \n\tgot: %s \n\twant: %s", err, tc.expected)
					}
				}
			}
		})
	}
}

func TestUninstallCluster(t *testing.T) {
	type testCase struct {
		jsonInput   string
		username    string
		mdmPassword string
		liaPassword string
		expected    error
	}

	cases := []testCase{
		{
			"",
			"test",
			"123",
			"123",
			errors.New("unexpected end of JSON input"),
		},
	}

	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	}))
	defer svr.Close()

	for _, tc := range cases {
		tc := tc
		t.Run("", func(ts *testing.T) {
			GC, err := NewGateway(svr.URL, "", "", true, true)
			if err != nil {
				t.Fatal(err)
			}

			_, err = GC.UninstallCluster(tc.jsonInput, tc.username, tc.mdmPassword, tc.liaPassword, true, true, false, true)
			if err != nil {
				if tc.expected == nil {
					t.Errorf("Uninstalling Cluster did not work as expected, \n\tgot: %s \n\twant: %v", err, tc.expected)
				} else {
					if err.Error() != tc.expected.Error() {
						t.Errorf("Uninstalling Cluster did not work as expected, \n\tgot: %s \n\twant: %s", err, tc.expected)
					}
				}
			}
		})
	}
}

func TestGetClusterDetails(t *testing.T) {
	type testCase struct {
		mdmIP       string
		mdmPassword string
		expected    error
	}

	cases := []testCase{
		{
			"",
			"",
			errors.New("Error Getting Cluster Details"),
		},
	}

	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	}))
	defer svr.Close()

	for _, tc := range cases {
		tc := tc
		t.Run("", func(ts *testing.T) {
			GC, err := NewGateway(svr.URL, "", "", true, true)
			if err != nil {
				t.Fatal(err)
			}

			clusterData := map[string]interface{}{
				"mdmUser":     "admin",
				"mdmPassword": tc.mdmPassword,
			}
			clusterData["mdmIps"] = []string{tc.mdmIP}

			secureData := map[string]interface{}{
				"allowNonSecureCommunicationWithMdm": true,
				"allowNonSecureCommunicationWithLia": true,
				"disableNonMgmtComponentsAuth":       false,
			}
			clusterData["securityConfiguration"] = secureData

			jsonres, _ := json.Marshal(clusterData)

			_, err = GC.GetClusterDetails(jsonres, false)
			if err != nil {
				if tc.expected == nil {
					t.Errorf("Uninstalling Cluster did not work as expected, \n\tgot: %s \n\twant: %v", err, tc.expected)
				} else {
					if err.Error() != tc.expected.Error() {
						t.Errorf("Uninstalling Cluster did not work as expected, \n\tgot: %s \n\twant: %s", err, tc.expected)
					}
				}
			}
		})
	}
}
