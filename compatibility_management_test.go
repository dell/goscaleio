/*
 *
 * Copyright © 2021-2024 Dell Inc. or its subsidiaries. All Rights Reserved.
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
	"encoding/json"
	"math"
	"net/http"
	"net/http/httptest"
	"testing"

	types "github.com/dell/goscaleio/types/v1"
)

func mockCompatibilityTestServerHandler(resp http.ResponseWriter, req *http.Request) {
	switch req.RequestURI {
	case "/api/v1/Compatibility":
		if req.Method == http.MethodGet {
			resp.WriteHeader(http.StatusOK)
			response := types.CompatibilityManagement{
				ID: "mock-compatibility-system-id",
			}
			content, err := json.Marshal(response)
			if err != nil {
				http.Error(resp, err.Error(), http.StatusNotFound)
			}
			resp.Write(content)
		} else if req.Method == http.MethodPost {
			resp.WriteHeader(http.StatusOK)
			response := types.CompatibilityManagement{
				ID: "mock-compatibility-management-id",
			}
			content, err := json.Marshal(response)
			if err != nil {
				http.Error(resp, err.Error(), http.StatusNotFound)
			}
			resp.Write(content)
		}
	}
}

func TestGetCompatibility(t *testing.T) {
	mockServer := httptest.NewServer(http.HandlerFunc(mockCompatibilityTestServerHandler))
	defer mockServer.Close()
	tests := []struct {
		name  string
		error string
	}{
		{
			name:  "compatibility management success",
			error: "",
		},
	}

	for _, tc := range tests {
		client, err := NewClientWithArgs(mockServer.URL, "3.6", math.MaxInt64, true, false)
		if err != nil {
			t.Fatal(err)
		}

		s := System{
			client: client,
		}

		_, err = s.GetCompatibilityManagement()
		if err != nil {
			if tc.error != err.Error() {
				t.Fatal(err)
			}
		}
	}
}

func TestSetCompatibility(t *testing.T) {
	mockServer := httptest.NewServer(http.HandlerFunc(mockCompatibilityTestServerHandler))
	defer mockServer.Close()
	tests := []struct {
		name                    string
		compatibilitymanagement types.CompatibilityManagementPost
		error                   string
	}{
		{
			name: "compatibility management success",
			compatibilitymanagement: types.CompatibilityManagementPost{
				ID: "mock-compatibility-management-id",
			},
			error: "",
		},
	}

	for _, tc := range tests {
		client, err := NewClientWithArgs(mockServer.URL, "4.0", math.MaxInt64, true, false)
		if err != nil {
			t.Fatal(err)
		}

		s := System{
			client: client,
		}

		_, err = s.SetCompatibilityManagement(&tc.compatibilitymanagement)
		if err != nil {
			if tc.error != err.Error() {
				t.Fatal(err)
			}
		}
	}
}
