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
	"context"
	"math"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_CheckPfmpVersion(t *testing.T) {
	type checkFn func(*testing.T, int, error)
	check := func(fns ...checkFn) []checkFn { return fns }

	hasNoError := func(t *testing.T, _ int, err error) {
		if err != nil {
			t.Fatalf("expected no error")
		}
	}

	checkVersionEqual := func(t *testing.T, result int, _ error) {
		assert.Equal(t, result, 0)
	}

	checkVersionGreaterThan := func(t *testing.T, result int, _ error) {
		assert.Equal(t, result, 1)
	}

	checkVersionLessThan := func(t *testing.T, result int, _ error) {
		assert.Equal(t, result, -1)
	}

	hasError := func(t *testing.T, _ int, err error) {
		if err == nil {
			t.Fatalf("expected error")
		}
	}

	tests := map[string]func(t *testing.T) (*httptest.Server, string, []checkFn){
		"PFMP version == 4.6": func(t *testing.T) (*httptest.Server, string, []checkFn) {
			url := "/Api/V1/corelcm/status"
			responseJSON := `{
				"lcmStatus": "READY",
				"clusterVersion": "4.6.0.0",
				"clusterBuild": "1258"
			}`
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method == http.MethodGet && strings.EqualFold(r.URL.Path, url) {
					w.WriteHeader(http.StatusOK)
					_, err := w.Write([]byte(responseJSON))
					if err != nil {
						t.Fatalf("Unexpected error: %v", err)
					}
					return
				}
				http.NotFound(w, r)
			}))
			return server, "4.6", check(hasNoError, checkVersionEqual)
		},
		"PFMP version > 4.6": func(t *testing.T) (*httptest.Server, string, []checkFn) {
			url := "/Api/V1/corelcm/status"
			responseJSON := `{
				"lcmStatus": "READY",
				"clusterVersion": "4.7.0.0",
				"clusterBuild": "1258"
			}`
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method == http.MethodGet && strings.EqualFold(r.URL.Path, url) {
					w.WriteHeader(http.StatusOK)
					_, err := w.Write([]byte(responseJSON))
					if err != nil {
						t.Fatalf("Unexpected error: %v", err)
					}
					return
				}
				http.NotFound(w, r)
			}))
			return server, "4.6.0.0", check(hasNoError, checkVersionGreaterThan)
		},
		"PFMP version < 4.6": func(t *testing.T) (*httptest.Server, string, []checkFn) {
			url := "/Api/V1/corelcm/status"
			responseJSON := `{
				"lcmStatus": "READY",
				"clusterVersion": "4.5.0.0",
				"clusterBuild": "1258"
			}`
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method == http.MethodGet && strings.EqualFold(r.URL.Path, url) {
					w.WriteHeader(http.StatusOK)
					_, err := w.Write([]byte(responseJSON))
					if err != nil {
						t.Fatalf("Unexpected error: %v", err)
					}
					return
				}
				http.NotFound(w, r)
			}))
			return server, "4.6", check(hasNoError, checkVersionLessThan)
		},
		"wrong version": func(t *testing.T) (*httptest.Server, string, []checkFn) {
			url := "/Api/V1/corelcm/status"
			responseJSON := `{
				"lcmStatus": "READY",
				"clusterVersion": "4.5.0.0",
				"clusterBuild": "1258"
			}`
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method == http.MethodGet && strings.EqualFold(r.URL.Path, url) {
					w.WriteHeader(http.StatusOK)
					_, err := w.Write([]byte(responseJSON))
					if err != nil {
						t.Fatalf("Unexpected error: %v", err)
					}
					return
				}
				http.NotFound(w, r)
			}))
			return server, "4.a.b.c", check(hasError)
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			server, version, checkFns := tc(t)
			defer server.Close()

			client, _ := NewClientWithArgs(server.URL, "", math.MaxInt64, true, false)
			result, err := CheckPfmpVersion(context.Background(), client, version)

			for _, checkFn := range checkFns {
				checkFn(t, result, err)
			}
		})
	}
}
