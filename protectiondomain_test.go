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
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"net/http"
	"net/http/httptest"
	"testing"

	types "github.com/dell/goscaleio/types/v1"
)

func TestCreateProtectionDomain(t *testing.T) {
	pdName := "myPD"
	type testCase struct {
		server      *httptest.Server
		expectedErr error
	}

	cases := map[string]testCase{
		"succeed": {
			server: httptest.NewServer(http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
				switch req.RequestURI {
				case "/api/types/ProtectionDomain/instances":
					resp.WriteHeader(http.StatusOK)
					response := types.ProtectionDomainResp{
						ID: "1",
					}

					content, err := json.Marshal(response)
					if err != nil {
						t.Fatal(err)
					}

					resp.Write(content)
				default:
					resp.WriteHeader(http.StatusNoContent)
				}
			})),
			expectedErr: nil,
		},
		"bad request": {
			server: httptest.NewServer(http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
				resp.WriteHeader(http.StatusBadRequest)
				resp.Write([]byte(`{"message":"bad request","httpStatusCode":400,"errorCode":0}`))
			})),
			expectedErr: errors.New("bad request"),
		},
	}

	for _, tc := range cases {
		client, err := NewClientWithArgs(tc.server.URL, "3.6", math.MaxInt64, true, false)
		if err != nil {
			t.Fatal(err)
		}

		s := System{
			client: client,
		}

		_, err = s.CreateProtectionDomain(context.Background(), pdName)
		if err != nil {
			if tc.expectedErr.Error() != err.Error() {
				t.Fatal(err)
			}
		}

		tc.server.Close()
	}
}

func TestGetProtectionDomainEx(t *testing.T) {
	pdId := "12345678-1234-1234-1234-123456789012"
	type testCase struct {
		server      *httptest.Server
		expectedErr error
	}

	cases := map[string]testCase{
		"succeed": {
			server: httptest.NewServer(http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
				switch req.RequestURI {
				case fmt.Sprintf("/api/instances/ProtectionDomain::%s", pdId):
					resp.WriteHeader(http.StatusOK)
					response := types.ProtectionDomain{
						ID:   pdId,
						Name: "domain1",
					}

					content, err := json.Marshal(response)
					if err != nil {
						t.Fatal(err)
					}

					resp.Write(content)
				default:
					resp.WriteHeader(http.StatusNoContent)
				}
			})),
			expectedErr: nil,
		},
		"bad request": {
			server: httptest.NewServer(http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
				resp.WriteHeader(http.StatusBadRequest)
				resp.Write([]byte(`{"message":"bad request","httpStatusCode":400,"errorCode":0}`))
			})),
			expectedErr: errors.New("error getting protection domain by id: bad request"),
		},
	}

	for _, tc := range cases {
		client, err := NewClientWithArgs(tc.server.URL, "3.6", math.MaxInt64, true, false)
		if err != nil {
			t.Fatal(err)
		}

		s := System{
			client: client,
		}

		_, err = s.GetProtectionDomainEx(context.Background(), pdId)
		if err != nil {
			if tc.expectedErr.Error() != err.Error() {
				t.Fatal(err)
			}
		}

		tc.server.Close()
	}
}
