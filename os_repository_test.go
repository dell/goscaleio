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
	"errors"
	"fmt"
	"math"
	"net/http"
	"net/http/httptest"
	"testing"

	types "github.com/dell/goscaleio/types/v1"
)

func TestGetAllOSRepositories(t *testing.T) {
	type testCase struct {
		server      *httptest.Server
		expectedErr error
	}

	cases := []testCase{
		{
			server: httptest.NewServer(http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
				switch req.RequestURI {
				case "/api/v1/OSRepository":
					resp.WriteHeader(http.StatusOK)
				default:
					resp.WriteHeader(http.StatusNoContent)
				}
			})),
			expectedErr: nil,
		},
		{
			server: httptest.NewServer(http.HandlerFunc(func(resp http.ResponseWriter, _ *http.Request) {
				resp.WriteHeader(http.StatusBadRequest)
				resp.Write([]byte(`{"message":"bad request","httpStatusCode":400,"errorCode":0}`))
			})),
			expectedErr: errors.New("bad request"),
		},
	}

	for _, tc := range cases {
		client, err := NewClientWithArgs(tc.server.URL, "", math.MaxInt64, true, false)
		client.configConnect.Version = "3.6"
		if err != nil {
			t.Fatal(err)
		}

		s := System{
			client: client,
		}

		_, err = s.GetAllOSRepositories()
		if err != nil {
			if tc.expectedErr.Error() != err.Error() {
				t.Fatal(err)
			}
		}

		tc.server.Close()
	}
}

func TestGetOSRepositoryByID(t *testing.T) {
	type testCase struct {
		id          string
		server      *httptest.Server
		expectedErr error
	}

	validID := "12345678-1234-1234-1234-123456789012"

	cases := []testCase{
		{
			id: validID,
			server: httptest.NewServer(http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
				switch req.RequestURI {
				case fmt.Sprintf("/api/v1/OSRepository/%s", validID):
					resp.WriteHeader(http.StatusOK)
				default:
					resp.WriteHeader(http.StatusNoContent)
				}
			})),
			expectedErr: nil,
		},
		{
			id: "invalid-id",
			server: httptest.NewServer(http.HandlerFunc(func(resp http.ResponseWriter, _ *http.Request) {
				resp.WriteHeader(http.StatusBadRequest)
				resp.Write([]byte(`{"message":"invalid id passed","httpStatusCode":400,"errorCode":0}`))
			})),
			expectedErr: errors.New("invalid id passed"),
		},
	}

	for _, tc := range cases {
		client, err := NewClientWithArgs(tc.server.URL, "", math.MaxInt64, true, false)
		client.configConnect.Version = "3.6"
		if err != nil {
			t.Fatal(err)
		}

		s := System{
			client: client,
		}

		_, err = s.GetOSRepositoryByID(tc.id)
		if err != nil {
			if tc.expectedErr.Error() != err.Error() {
				t.Fatal(err)
			}
		}

		tc.server.Close()
	}
}

func TestCreateOSRepository(t *testing.T) {
	type testCase struct {
		repoContent *types.OSRepository
		server      *httptest.Server
		expectedErr error
	}

	cases := []testCase{
		{
			repoContent: &types.OSRepository{
				Name:       "my-repo",
				RepoType:   "S3",
				SourcePath: "source",
				ImageType:  "ISO",
			},
			server: httptest.NewServer(http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
				switch req.RequestURI {
				case "/api/v1/OSRepository":
					resp.WriteHeader(http.StatusOK)
				default:
					resp.WriteHeader(http.StatusNoContent)
				}
			})),
			expectedErr: nil,
		},
		{
			repoContent: nil,
			server:      httptest.NewServer(http.HandlerFunc(func(_ http.ResponseWriter, _ *http.Request) {})),
			expectedErr: errors.New("createOSRepository cannot be nil"),
		},
		{
			repoContent: &types.OSRepository{
				Name:       "my-repo",
				RepoType:   "S3",
				SourcePath: "source",
				ImageType:  "ISO",
			},
			server: httptest.NewServer(http.HandlerFunc(func(resp http.ResponseWriter, _ *http.Request) {
				resp.WriteHeader(http.StatusBadRequest)
				resp.Write([]byte(`{"message":"invalid create repo request","httpStatusCode":400,"errorCode":0}`))
			})),
			expectedErr: errors.New("invalid create repo request"),
		},
	}

	for _, tc := range cases {
		client, err := NewClientWithArgs(tc.server.URL, "", math.MaxInt64, true, false)
		client.configConnect.Version = "3.6"
		if err != nil {
			t.Fatal(err)
		}

		s := System{
			client: client,
		}

		_, err = s.CreateOSRepository(tc.repoContent)
		if err != nil {
			if tc.expectedErr.Error() != err.Error() {
				t.Fatal(err)
			}
		}

		tc.server.Close()
	}
}

func TestRemoveOSRepository(t *testing.T) {
	type testCase struct {
		id          string
		server      *httptest.Server
		expectedErr error
	}

	validID := "12345678-1234-1234-1234-123456789012"

	cases := []testCase{
		{
			id: validID,
			server: httptest.NewServer(http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
				switch req.RequestURI {
				case fmt.Sprintf("/api/v1/OSRepository/%s", validID):
					resp.WriteHeader(http.StatusOK)
				default:
					resp.WriteHeader(http.StatusNoContent)
				}
			})),
			expectedErr: nil,
		},
		{
			id: "invalid-id",
			server: httptest.NewServer(http.HandlerFunc(func(resp http.ResponseWriter, _ *http.Request) {
				resp.WriteHeader(http.StatusBadRequest)
				resp.Write([]byte(`{"message":"invalid id passed","httpStatusCode":400,"errorCode":0}`))
			})),
			expectedErr: errors.New("invalid id passed"),
		},
	}

	for _, tc := range cases {
		client, err := NewClientWithArgs(tc.server.URL, "", math.MaxInt64, true, false)
		client.configConnect.Version = "3.6"
		if err != nil {
			t.Fatal(err)
		}

		s := System{
			client: client,
		}

		err = s.RemoveOSRepository(tc.id)
		if err != nil {
			if tc.expectedErr.Error() != err.Error() {
				t.Fatal(err)
			}
		}

		tc.server.Close()
	}
}
