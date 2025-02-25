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
	"math"
	"net/http"
	"net/http/httptest"
	"testing"

	types "github.com/dell/goscaleio/types/v1"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestCreateSds(t *testing.T) {
	pdID := uuid.NewString()
	pdName := "myPD"
	type testCase struct {
		ipList      []string
		server      *httptest.Server
		expectedErr error
	}

	cases := map[string]testCase{
		"succeed: single IP": {
			ipList: []string{"127.0.0.1"},
			server: httptest.NewServer(http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
				switch req.RequestURI {
				case "/api/types/Sds/instances":
					resp.WriteHeader(http.StatusOK)
				default:
					resp.WriteHeader(http.StatusBadRequest)
					resp.Write([]byte(`{"message":"no route handled","httpStatusCode":400,"errorCode":0}`))
				}
			})),
			expectedErr: nil,
		},
		"succeed: multiple IPs": {
			ipList: []string{"127.0.0.1", "127.0.0.2"},
			server: httptest.NewServer(http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
				switch req.RequestURI {
				case "/api/types/Sds/instances":
					resp.WriteHeader(http.StatusOK)
				default:
					resp.WriteHeader(http.StatusBadRequest)
					resp.Write([]byte(`{"message":"no route handled","httpStatusCode":400,"errorCode":0}`))
				}
			})),
			expectedErr: nil,
		},
		"error: no IPs": {
			ipList: []string{},
			server: httptest.NewServer(http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
				switch req.RequestURI {
				case "/api/types/Sds/instances":
					resp.WriteHeader(http.StatusOK)
				default:
					resp.WriteHeader(http.StatusBadRequest)
					resp.Write([]byte(`{"message":"no route handled","httpStatusCode":400,"errorCode":0}`))
				}
			})),
			expectedErr: fmt.Errorf("Must provide at least 1 SDS IP"),
		},
		"error: too many IPs": {
			ipList: []string{"127.0.0.1", "127.0.0.2", "127.0.0.3"},
			server: httptest.NewServer(http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
				switch req.RequestURI {
				case "/api/types/Sds/instances":
					resp.WriteHeader(http.StatusOK)
				default:
					resp.WriteHeader(http.StatusBadRequest)
					resp.Write([]byte(`{"message":"no route handled","httpStatusCode":400,"errorCode":0}`))
				}
			})),
			expectedErr: fmt.Errorf("Must explicitly provide IP role for more than 2 SDS IPs"),
		},
		"error: bad request": {
			ipList: []string{"127.0.0.1", "127.0.0.2"},
			server: httptest.NewServer(http.HandlerFunc(func(resp http.ResponseWriter, _ *http.Request) {
				resp.WriteHeader(http.StatusBadRequest)
				resp.Write([]byte(`{"message":"bad request","httpStatusCode":400,"errorCode":0}`))
			})),
			expectedErr: fmt.Errorf("bad request"),
		},
	}

	for _, tc := range cases {
		client, err := NewClientWithArgs(tc.server.URL, "3.6", math.MaxInt64, true, false)
		if err != nil {
			t.Fatal(err)
		}
		defer tc.server.Close()

		pd := ProtectionDomain{
			ProtectionDomain: &types.ProtectionDomain{
				ID: pdID,
			},
			client: client,
		}

		_, err = pd.CreateSds(pdName, tc.ipList)
		if err != nil {
			if tc.expectedErr.Error() != err.Error() {
				t.Fatal(err)
			}
		}
	}
}

func TestCreateSdsWithParams(t *testing.T) {
	pdID := uuid.NewString()
	sds := types.Sds{
		Name:               "mySds",
		ProtectionDomainID: uuid.NewString(),
		Port:               8080,
		RmcacheEnabled:     true,
		RmcacheSizeInKb:    1024,
		DrlMode:            "Volatile",
		FaultSetID:         uuid.NewString(),
	}
	type testCase struct {
		sdsIPs      []*types.SdsIP
		server      *httptest.Server
		expectedErr error
	}

	cases := map[string]testCase{
		"succeed: single IP": {
			sdsIPs: []*types.SdsIP{
				{
					Role: RoleAll,
					IP:   "127.0.0.1",
				},
			},
			server: httptest.NewServer(http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
				switch req.RequestURI {
				case "/api/types/Sds/instances":
					resp.WriteHeader(http.StatusOK)
				default:
					resp.WriteHeader(http.StatusBadRequest)
					resp.Write([]byte(`{"message":"no route handled","httpStatusCode":400,"errorCode":0}`))
				}
			})),
			expectedErr: nil,
		},
		"succeed: multiple IPs": {
			sdsIPs: []*types.SdsIP{
				{
					Role: RoleAll,
					IP:   "127.0.0.1",
				},
				{
					Role: RoleAll,
					IP:   "127.0.0.2",
				},
			},
			server: httptest.NewServer(http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
				switch req.RequestURI {
				case "/api/types/Sds/instances":
					resp.WriteHeader(http.StatusOK)
				default:
					resp.WriteHeader(http.StatusBadRequest)
					resp.Write([]byte(`{"message":"no route handled","httpStatusCode":400,"errorCode":0}`))
				}
			})),
			expectedErr: nil,
		},
		"error: no IPs": {
			sdsIPs: []*types.SdsIP{},
			server: httptest.NewServer(http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
				switch req.RequestURI {
				case "/api/types/Sds/instances":
					resp.WriteHeader(http.StatusOK)
				default:
					resp.WriteHeader(http.StatusBadRequest)
					resp.Write([]byte(`{"message":"no route handled","httpStatusCode":400,"errorCode":0}`))
				}
			})),
			expectedErr: fmt.Errorf("Must provide at least 1 SDS IP"),
		},
		"error: single IP, not role ALL": {
			sdsIPs: []*types.SdsIP{
				{
					Role: "InvalidRole",
					IP:   "127.0.0.1",
				},
			},
			server: httptest.NewServer(http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
				switch req.RequestURI {
				case "/api/types/Sds/instances":
					resp.WriteHeader(http.StatusOK)
				default:
					resp.WriteHeader(http.StatusBadRequest)
					resp.Write([]byte(`{"message":"no route handled","httpStatusCode":400,"errorCode":0}`))
				}
			})),
			expectedErr: fmt.Errorf("The only IP assigned to an SDS must be assigned \"%s\" role", RoleAll),
		},
		"error: multiple IP, no sdc": {
			sdsIPs: []*types.SdsIP{
				{
					Role: RoleSdsOnly,
					IP:   "127.0.0.1",
				},
				{
					Role: RoleSdsOnly,
					IP:   "127.0.0.2",
				},
			},
			server: httptest.NewServer(http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
				switch req.RequestURI {
				case "/api/types/Sds/instances":
					resp.WriteHeader(http.StatusOK)
				default:
					resp.WriteHeader(http.StatusBadRequest)
					resp.Write([]byte(`{"message":"no route handled","httpStatusCode":400,"errorCode":0}`))
				}
			})),
			expectedErr: fmt.Errorf("At least one IP must be assigned %s or %s role", RoleSdcOnly, RoleAll),
		},
		"error: multiple IP, no sds": {
			sdsIPs: []*types.SdsIP{
				{
					Role: RoleSdcOnly,
					IP:   "127.0.0.1",
				},
				{
					Role: RoleSdcOnly,
					IP:   "127.0.0.2",
				},
			},
			server: httptest.NewServer(http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
				switch req.RequestURI {
				case "/api/types/Sds/instances":
					resp.WriteHeader(http.StatusOK)
				default:
					resp.WriteHeader(http.StatusBadRequest)
					resp.Write([]byte(`{"message":"no route handled","httpStatusCode":400,"errorCode":0}`))
				}
			})),
			expectedErr: fmt.Errorf("At least one IP must be assigned %s or %s role", RoleSdsOnly, RoleAll),
		},
	}

	for _, tc := range cases {
		client, err := NewClientWithArgs(tc.server.URL, "3.6", math.MaxInt64, true, false)
		if err != nil {
			t.Fatal(err)
		}
		defer tc.server.Close()

		pd := ProtectionDomain{
			ProtectionDomain: &types.ProtectionDomain{
				ID: pdID,
			},
			client: client,
		}

		sds.IPList = tc.sdsIPs
		_, err = pd.CreateSdsWithParams(&sds)
		if err != nil {
			if tc.expectedErr.Error() != err.Error() {
				t.Fatal(err)
			}
		}
	}
}

func TestGetSds(t *testing.T) {
	pdID := uuid.NewString()
	type testCase struct {
		server      *httptest.Server
		expectedErr error
	}

	cases := map[string]testCase{
		"succeed": {
			server: httptest.NewServer(http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
				switch req.RequestURI {
				case fmt.Sprintf("/api/instances/ProtectionDomain::%s/relationships/Sds", pdID):
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
			expectedErr: fmt.Errorf("bad request"),
		},
	}

	for _, tc := range cases {
		client, err := NewClientWithArgs(tc.server.URL, "3.6", math.MaxInt64, true, false)
		if err != nil {
			t.Fatal(err)
		}
		defer tc.server.Close()

		pd := ProtectionDomain{
			ProtectionDomain: &types.ProtectionDomain{
				ID: pdID,
			},
			client: client,
		}

		_, err = pd.GetSds()
		if err != nil {
			if tc.expectedErr.Error() != err.Error() {
				t.Fatal(err)
			}
		}
	}
}

func TestGetAllSds(t *testing.T) {
	type testCase struct {
		server      *httptest.Server
		expectedErr error
	}

	cases := map[string]testCase{
		"succeed": {
			server: httptest.NewServer(http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
				switch req.RequestURI {
				case "/api/types/Sds/instances":
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
			expectedErr: fmt.Errorf("bad request"),
		},
	}

	for _, tc := range cases {
		client, err := NewClientWithArgs(tc.server.URL, "3.6", math.MaxInt64, true, false)
		if err != nil {
			t.Fatal(err)
		}
		defer tc.server.Close()

		system := System{
			System: &types.System{
				ID: uuid.NewString(),
			},
			client: client,
		}

		_, err = system.GetAllSds()
		if err != nil {
			if tc.expectedErr.Error() != err.Error() {
				t.Fatal(err)
			}
		}
	}
}

func TestFindSds(t *testing.T) {
	pdID := uuid.NewString()
	searchSdsID := uuid.NewString()
	type testCase struct {
		server      *httptest.Server
		expectedErr error
	}

	cases := map[string]testCase{
		"succeed": {
			server: httptest.NewServer(http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
				switch req.RequestURI {
				case fmt.Sprintf("/api/instances/ProtectionDomain::%s/relationships/Sds", pdID):
					content, err := json.Marshal([]types.Sds{
						{
							ID: searchSdsID,
						},
					})
					if err != nil {
						t.Fatal(err)
					}

					resp.Write(content)
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
			expectedErr: fmt.Errorf("bad request"),
		},
		"error: could not find sds": {
			server: httptest.NewServer(http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
				switch req.RequestURI {
				case fmt.Sprintf("/api/instances/ProtectionDomain::%s/relationships/Sds", pdID):
					content, err := json.Marshal([]types.Sds{
						{
							ID: uuid.NewString(),
						},
					})
					if err != nil {
						t.Fatal(err)
					}

					resp.Write(content)
					resp.WriteHeader(http.StatusOK)
				default:
					resp.WriteHeader(http.StatusBadRequest)
					resp.Write([]byte(`{"message":"no route handled","httpStatusCode":400,"errorCode":0}`))
				}
			})),
			expectedErr: fmt.Errorf("Couldn't find SDS"),
		},
	}

	for _, tc := range cases {
		client, err := NewClientWithArgs(tc.server.URL, "3.6", math.MaxInt64, true, false)
		if err != nil {
			t.Fatal(err)
		}
		defer tc.server.Close()

		pd := ProtectionDomain{
			ProtectionDomain: &types.ProtectionDomain{
				ID: pdID,
			},
			client: client,
		}

		_, err = pd.FindSds("ID", searchSdsID)
		if err != nil {
			if tc.expectedErr.Error() != err.Error() {
				t.Fatal(err)
			}
		}
	}
}

func TestFindSdsSystem(t *testing.T) {
	pdID := uuid.NewString()
	searchSdsID := uuid.NewString()
	type testCase struct {
		server      *httptest.Server
		expectedErr error
	}

	cases := map[string]testCase{
		"succeed": {
			server: httptest.NewServer(http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
				switch req.RequestURI {
				case "/api/types/Sds/instances":
					content, err := json.Marshal([]types.Sds{
						{
							ID: searchSdsID,
						},
					})
					if err != nil {
						t.Fatal(err)
					}

					resp.Write(content)
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
			expectedErr: fmt.Errorf("bad request"),
		},
		"error: could not find sds": {
			server: httptest.NewServer(http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
				switch req.RequestURI {
				case "/api/types/Sds/instances":
					content, err := json.Marshal([]types.Sds{
						{
							ID: uuid.NewString(),
						},
					})
					if err != nil {
						t.Fatal(err)
					}

					resp.Write(content)
					resp.WriteHeader(http.StatusOK)
				default:
					resp.WriteHeader(http.StatusBadRequest)
					resp.Write([]byte(`{"message":"no route handled","httpStatusCode":400,"errorCode":0}`))
				}
			})),
			expectedErr: fmt.Errorf("Couldn't find SDS"),
		},
	}

	for _, tc := range cases {
		client, err := NewClientWithArgs(tc.server.URL, "3.6", math.MaxInt64, true, false)
		if err != nil {
			t.Fatal(err)
		}
		defer tc.server.Close()

		system := System{
			System: &types.System{
				ID: pdID,
			},
			client: client,
		}

		_, err = system.FindSds("ID", searchSdsID)
		if err != nil {
			if tc.expectedErr.Error() != err.Error() {
				t.Fatal(err)
			}
		}
	}
}

func TestGetSdsByID(t *testing.T) {
	systemID := uuid.NewString()
	searchSdsID := uuid.NewString()
	type testCase struct {
		server      *httptest.Server
		expectedErr error
	}

	cases := map[string]testCase{
		"succeed": {
			server: httptest.NewServer(http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
				switch req.RequestURI {
				case fmt.Sprintf("/api/instances/Sds::%s", searchSdsID):
					content, err := json.Marshal(types.Sds{
						ID: searchSdsID,
					})
					if err != nil {
						t.Fatal(err)
					}

					resp.Write(content)
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
			expectedErr: fmt.Errorf("bad request"),
		},
	}

	for _, tc := range cases {
		client, err := NewClientWithArgs(tc.server.URL, "3.6", math.MaxInt64, true, false)
		if err != nil {
			t.Fatal(err)
		}
		defer tc.server.Close()

		pd := System{
			System: &types.System{
				ID: systemID,
			},
			client: client,
		}

		_, err = pd.GetSdsByID(searchSdsID)
		if err != nil {
			if tc.expectedErr.Error() != err.Error() {
				t.Fatal(err)
			}
		}
	}
}

func TestSdsActionsSuccess(t *testing.T) {
	pdID := uuid.NewString()
	sdsID := uuid.NewString()

	server := httptest.NewServer(http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
		switch req.RequestURI {
		case fmt.Sprintf("/api/instances/Sds::%s/action/removeSds", sdsID):
		case fmt.Sprintf("/api/instances/Sds::%s/action/addSdsIp", sdsID):
		case fmt.Sprintf("/api/instances/Sds::%s/action/setSdsIpRole", sdsID):
		case fmt.Sprintf("/api/instances/Sds::%s/action/removeSdsIp", sdsID):
		case fmt.Sprintf("/api/instances/Sds::%s/action/setSdsName", sdsID):
		case fmt.Sprintf("/api/instances/Sds::%s/action/setSdsPort", sdsID):
		case fmt.Sprintf("/api/instances/Sds::%s/action/setDrlMode", sdsID):
		case fmt.Sprintf("/api/instances/Sds::%s/action/enableRfcache", sdsID):
		case fmt.Sprintf("/api/instances/Sds::%s/action/disableRfcache", sdsID):
		case fmt.Sprintf("/api/instances/Sds::%s/action/setSdsRmcacheEnabled", sdsID):
		case fmt.Sprintf("/api/instances/Sds::%s/action/setSdsRmcacheSize", sdsID):
		case fmt.Sprintf("/api/instances/Sds::%s/action/setSdsPerformanceParameters", sdsID):
			resp.WriteHeader(http.StatusOK)
		default:
			resp.WriteHeader(http.StatusBadRequest)
			resp.Write([]byte(`{"message":"no route handled","httpStatusCode":400,"errorCode":0}`))
		}
	}))

	client, err := NewClientWithArgs(server.URL, "3.6", math.MaxInt64, true, false)
	if err != nil {
		t.Fatal(err)
	}
	defer server.Close()

	pd := ProtectionDomain{
		ProtectionDomain: &types.ProtectionDomain{
			ID: pdID,
		},
		client: client,
	}

	err = pd.DeleteSds(sdsID)
	assert.Nil(t, err)

	err = pd.AddSdSIP(sdsID, "127.0.0.3", RoleSdsOnly)
	assert.Nil(t, err)

	err = pd.SetSDSIPRole(sdsID, "127.0.0.1", RoleAll)
	assert.Nil(t, err)

	err = pd.RemoveSDSIP(sdsID, "127.0.0.1")
	assert.Nil(t, err)

	err = pd.SetSdsName(sdsID, "newSdsName")
	assert.Nil(t, err)

	err = pd.SetSdsPort(sdsID, 8081)
	assert.Nil(t, err)

	err = pd.SetSdsDrlMode(sdsID, "Volatile")
	assert.Nil(t, err)

	err = pd.SetSdsRfCache(sdsID, true)
	assert.Nil(t, err)

	err = pd.SetSdsRfCache(sdsID, false)
	assert.Nil(t, err)

	err = pd.SetSdsRmCache(sdsID, true)
	assert.Nil(t, err)

	err = pd.SetSdsRmCacheSize(sdsID, 1024)
	assert.Nil(t, err)

	err = pd.SetSdsPerformanceProfile(sdsID, "PowerProfile")
	assert.Nil(t, err)
}

func TestNewSdsEx(t *testing.T) {
	assert.NotNil(t, NewSdsEx(nil, nil))
}
