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
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	types "github.com/dell/goscaleio/types/v1"
	"github.com/google/uuid"
)

// This test can be checked when NewGateway() function is fixed
func TestUploadCompliance(t *testing.T) {
	type testCase struct {
		server      *httptest.Server
		expectedErr error
	}

	cases := map[string]testCase{
		"success": {
			server: httptest.NewServer(http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
				switch req.RequestURI {
				case "/api/version":
					resp.WriteHeader(http.StatusOK)
					fmt.Fprintln(resp, "4.0")
				case "/rest/auth/login":
					resp.WriteHeader(http.StatusOK)
					fmt.Fprintln(resp, `{"access_token":"mock_access_token"}`)
				case "/Api/V1/FirmwareRepository":
					resp.WriteHeader(http.StatusCreated)
					content, err := json.Marshal(types.UploadComplianceTopologyDetails{
						ID: uuid.NewString(),
					})
					if err != nil {
						t.Fatal(err)
					}

					resp.Write(content)
				default:
					resp.WriteHeader(http.StatusBadRequest)
					resp.Write([]byte(`{"message":"no route handled","httpStatusCode":400,"errorCode":0}`))
				}
			})),
			expectedErr: nil,
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			defer tc.server.Close()

			gc, err := NewGateway(tc.server.URL, "test_username", "test_password", false, false)
			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}

			uploadParams := types.UploadComplianceParam{
				Username: "user",
				Password: "password",
			}

			_, err = gc.UploadCompliance(&uploadParams)
			if err != nil {
				if tc.expectedErr.Error() != err.Error() {
					t.Fatal(err)
				}
			}
		})
	}
}
