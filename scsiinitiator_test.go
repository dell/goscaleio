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
	"errors"
	"fmt"
	"math"
	"net/http"
	"net/http/httptest"
	"testing"

	types "github.com/dell/goscaleio/types/v1"
	"github.com/google/uuid"
)

func TestGetScsiInitiator(t *testing.T) {
	systemID := uuid.NewString()
	type testCase struct {
		server      *httptest.Server
		expectedErr error
	}

	cases := map[string]testCase{
		"succeed": {
			server: httptest.NewServer(http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
				switch req.RequestURI {
				case fmt.Sprintf("/api/instances/System::%v/relationships/ScsiInitiator", systemID):
					resp.WriteHeader(http.StatusOK)
				default:
					resp.WriteHeader(http.StatusBadRequest)
					resp.Write([]byte(`{"message":"no route handled","httpStatusCode":400,"errorCode":0}`))
				}
			})),
			expectedErr: nil,
		},
		"error: bad request": {
			server: httptest.NewServer(http.HandlerFunc(func(resp http.ResponseWriter, _ *http.Request) {
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

		system := System{
			System: &types.System{
				ID: systemID,
			},
			client: client,
		}

		system.GetScsiInitiator()

		tc.server.Close()
	}
}
