// Copyright © 2023 Dell Inc. or its subsidiaries. All Rights Reserved.
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
	"math"
	"net/http"
	"net/http/httptest"
	"testing"

	types "github.com/dell/goscaleio/types/v1"
	"github.com/stretchr/testify/assert"
)

var (
	ID    string
	Name  string
	errFs error
)

func TestCreateFaultSet(t *testing.T) {
	type testCase struct {
		fs       *types.FaultSetParam
		expected error
	}
	cases := []testCase{
		{
			fs: &types.FaultSetParam{
				Name:               "testFaultSet",
				ProtectionDomainID: "202a046600000000",
			},
			expected: nil,
		},
		{
			fs: &types.FaultSetParam{
				Name:               "testFaultSet",
				ProtectionDomainID: "202a0466000000",
			},
			expected: errors.New("invalid Protection Domain. Please try again with the correct ID or name"),
		},
	}
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

			p := NewProtectionDomain(client)

			ID, errFs = p.CreateFaultSet(tc.fs)
			if errFs != nil {
				if tc.expected == nil {
					t.Errorf("Creating fault set did not work as expected, \n\tgot: %s \n\twant: %v", errFs, tc.expected)
				} else {
					if errFs.Error() != tc.expected.Error() {
						t.Errorf("Creating fault set did not work as expected, \n\tgot: %s \n\twant: %s", errFs, tc.expected)
					}
				}
			}
		})
	}
}

func TestGetFaultByID(t *testing.T) {
	type testCase struct {
		id       string
		expected error
	}
	cases := []testCase{
		{
			id:       ID,
			expected: nil,
		},
		{
			id:       "1234",
			expected: errors.New("invalid Fault Set. Please try again with the correct ID or name"),
		},
	}
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

			s := System{
				client: client,
			}
			_, err2 := s.GetFaultSetByID(tc.id)
			if err2 != nil {
				if tc.expected == nil {
					t.Errorf("Fetching fault set did not work as expected, \n\tgot: %s \n\twant: %v", err2, tc.expected)
				} else {
					if err2.Error() != tc.expected.Error() {
						t.Errorf("Fetching fault set did not work as expected, \n\tgot: %s \n\twant: %s", err2, tc.expected)
					}
				}
			}
		})
	}
}

func TestModifyFaultSetName(t *testing.T) {
	type testCase struct {
		name     string
		id       string
		expected error
	}
	cases := []testCase{
		{
			id:       ID,
			name:     "renameFaultSet",
			expected: nil,
		},
		{
			id:       "1234",
			name:     "renameFaultSet",
			expected: errors.New("invalid Fault Set. Please try again with the correct ID or name"),
		},
	}
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

			p := ProtectionDomain{
				client: client,
			}
			err2 := p.ModifyFaultSetName(tc.id, tc.name)
			if err2 != nil {
				if tc.expected == nil {
					t.Errorf("Modifying fault set did not work as expected, \n\tgot: %s \n\twant: %v", err2, tc.expected)
				} else {
					if err2.Error() != tc.expected.Error() {
						t.Errorf("Modifying fault set did not work as expected, \n\tgot: %s \n\twant: %s", err2, tc.expected)
					}
				}
			}
		})
	}
}

func TestModifyFaultPerfProfile(t *testing.T) {
	type testCase struct {
		perfProfile string
		id          string
		expected    error
	}
	cases := []testCase{
		{
			id:          ID,
			perfProfile: "Compact",
			expected:    nil,
		},
		{
			id:          ID,
			perfProfile: "HighPerformance",
			expected:    nil,
		},
		{
			id:          ID,
			perfProfile: "Invalid",
			expected:    errors.New("perfProfile should get one of the following values: Compact, HighPerformance, but its value is Invalid"),
		},
	}
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

			p := ProtectionDomain{
				client: client,
			}
			err2 := p.ModifyFaultSetPerfProfile(tc.id, tc.perfProfile)
			if err2 != nil {
				if tc.expected == nil {
					t.Errorf("Modifying fault set did not work as expected, \n\tgot: %s \n\twant: %v", err2, tc.expected)
				} else {
					if err2.Error() != tc.expected.Error() {
						t.Errorf("Modifying fault set did not work as expected, \n\tgot: %s \n\twant: %s", err2, tc.expected)
					}
				}
			}
		})
	}
}

func TestDeleteFaultSet(t *testing.T) {
	type testCase struct {
		id       string
		expected error
	}
	cases := []testCase{
		{
			id:       ID,
			expected: nil,
		},
		{
			id:       "1234",
			expected: errors.New("invalid Fault Set. Please try again with the correct ID or name"),
		},
	}
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

			p := ProtectionDomain{
				client: client,
			}
			err2 := p.DeleteFaultSet(tc.id)
			if err2 != nil {
				if tc.expected == nil {
					t.Errorf("Removing fault set did not work as expected, \n\tgot: %s \n\twant: %v", err2, tc.expected)
				} else {
					if err2.Error() != tc.expected.Error() {
						t.Errorf("Removing fault set did not work as expected, \n\tgot: %s \n\twant: %s", err2, tc.expected)
					}
				}
			}
		})
	}
}

func TestGetAllFaultSets(t *testing.T) {
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	}))
	defer svr.Close()

	client, err := NewClientWithArgs(svr.URL, "", math.MaxInt64, true, false)
	if err != nil {
		t.Fatal(err)
	}

	s := System{
		client: client,
	}

	faultsets, err := s.GetAllFaultSets()
	assert.Equal(t, len(faultsets), 0)
	assert.Nil(t, err)
}

func TestGetAllFaultSetsSds(t *testing.T) {
	type testCase struct {
		id       string
		expected error
	}
	cases := []testCase{
		{
			id:       "6b2b5ce800000000",
			expected: nil,
		},
		{
			id:       "6b2b5ce800000001",
			expected: errors.New("Error in get relationship Sds"),
		},
	}

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

			s := System{
				client: client,
			}
			_, err2 := s.GetAllSDSByFaultSetID(tc.id)
			if err2 != nil {
				if tc.expected == nil {
					t.Errorf("Getting sds related with fault set did not work as expected, \n\tgot: %s \n\twant: %v", err2, tc.expected)
				} else {
					if err2.Error() != tc.expected.Error() {
						t.Errorf("Getting sds related with fault set did not work as expected, \n\tgot: %s \n\twant: %s", err2, tc.expected)
					}
				}
			}
		})
	}
}

func TestGetFaultSetByName(t *testing.T) {
	type testCase struct {
		server      *httptest.Server
		expectedErr error
	}

	cases := map[string]testCase{
		"success": {
			server: httptest.NewServer(http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
				switch req.RequestURI {
				case "/api/types/FaultSet/instances":
					resp.WriteHeader(http.StatusOK)
					response := []types.FaultSet{
						{ID: "mock-fault-set-id", Name: "mock-fault-set-name"},
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

	for id, tc := range cases {
		t.Run(id, func(t *testing.T) {
			client, err := NewClientWithArgs(tc.server.URL, "", math.MaxInt64, true, false)
			client.configConnect.Version = "4.0"
			if err != nil {
				t.Fatal(err)
			}

			s := System{
				client: client,
			}
			_, err = s.GetFaultSetByName(context.Background(), "mock-fault-set-name")
			if err != nil {
				if tc.expectedErr.Error() != err.Error() {
					t.Fatal(err)
				}
			}
		})
	}
}
