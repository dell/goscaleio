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
		"error: invalid http response": {
			server: httptest.NewServer(http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
				resp.WriteHeader(http.StatusBadRequest)
				resp.Write([]byte(`{"message":"bad request","httpStatusCode":400,"errorCode":0}`))
			})),
			expectedErr: errors.New("error getting protection domain by id: bad request"),
		},
	}

	for title, tc := range cases {
		fmt.Printf("running test case: %s\n", title)
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

func TestDeleteProtectionDomain(t *testing.T) {
	pdName := "myDomain"
	domainHost := "localhost"
	pdId := "12345678-1234-1234-1234-123456789012"
	type testCase struct {
		server      *httptest.Server
		expectedErr error
	}

	cases := map[string]testCase{
		"succeed": {
			server: httptest.NewServer(http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
				switch req.RequestURI {
				case "/api/System/relationship/ProtectionDomain":
					resp.WriteHeader(http.StatusOK)
					response := []types.ProtectionDomain{
						{
							ID:   pdId,
							Name: pdName,
							Links: []*types.Link{
								{Rel: "self", HREF: domainHost},
							},
						},
					}

					content, err := json.Marshal(response)
					if err != nil {
						t.Fatal(err)
					}

					resp.Write(content)
				case fmt.Sprintf("/%s/action/removeProtectionDomain", domainHost):
					resp.WriteHeader(http.StatusOK)
				default:
					resp.WriteHeader(http.StatusNoContent)
				}
			})),
			expectedErr: nil,
		},
		"error: invalid relationship protection domain call": {
			server: httptest.NewServer(http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
				switch req.RequestURI {
				case "/api/System/relationship/ProtectionDomain":
					resp.WriteHeader(http.StatusBadRequest)
					resp.Write([]byte(`{"message":"bad request","httpStatusCode":400,"errorCode":0}`))
				default:
					resp.WriteHeader(http.StatusNoContent)
				}
			})),
			expectedErr: errors.New("Error getting protection domains bad request"),
		},
		"error: unable to find protection domain": {
			server: httptest.NewServer(http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
				switch req.RequestURI {
				case "/api/System/relationship/ProtectionDomain":
					resp.WriteHeader(http.StatusOK)
					response := []types.ProtectionDomain{
						{
							ID:   "invalidID",
							Name: "invalidName",
						},
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
			expectedErr: errors.New("Couldn't find protection domain"),
		},
		"error: invalid domain links": {
			server: httptest.NewServer(http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
				switch req.RequestURI {
				case "/api/System/relationship/ProtectionDomain":
					resp.WriteHeader(http.StatusOK)
					response := []types.ProtectionDomain{
						{
							ID:   pdId,
							Name: pdName,
						},
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
			expectedErr: errors.New("Error: problem finding link"),
		},
		"error: unable to remove protection domain": {
			server: httptest.NewServer(http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
				switch req.RequestURI {
				case "/api/System/relationship/ProtectionDomain":
					resp.WriteHeader(http.StatusOK)
					response := []types.ProtectionDomain{
						{
							ID:   pdId,
							Name: pdName,
							Links: []*types.Link{
								{Rel: "self", HREF: domainHost},
							},
						},
					}

					content, err := json.Marshal(response)
					if err != nil {
						t.Fatal(err)
					}

					resp.Write(content)
				case fmt.Sprintf("/%s/action/removeProtectionDomain", domainHost):
					resp.WriteHeader(http.StatusBadRequest)
					resp.Write([]byte(`{"message":"unable to remove protection domain","httpStatusCode":400,"errorCode":0}`))
				default:
					resp.WriteHeader(http.StatusNoContent)
				}
			})),
			expectedErr: errors.New("unable to remove protection domain"),
		},
	}

	for title, tc := range cases {
		fmt.Printf("running test case: %s\n", title)
		client, err := NewClientWithArgs(tc.server.URL, "3.6", math.MaxInt64, true, false)
		if err != nil {
			t.Fatal(err)
		}

		s := System{
			System: &types.System{
				Links: []*types.Link{{Rel: "/api/System/relationship/ProtectionDomain", HREF: "/api/System/relationship/ProtectionDomain"}, {Rel: "self", HREF: "localhost"}},
			},
			client: client,
		}

		err = s.DeleteProtectionDomain(context.Background(), pdName)
		if err != nil {
			if tc.expectedErr.Error() != err.Error() {
				t.Fatal(err)
			}
		}

		tc.server.Close()
	}
}

func TestProtectionDomainDelete(t *testing.T) {
	domainHost := "localhost"
	type testCase struct {
		protectionDomain *types.ProtectionDomain
		server           *httptest.Server
		expectedErr      error
	}

	cases := map[string]testCase{
		"succeed": {
			protectionDomain: &types.ProtectionDomain{
				Links: []*types.Link{
					{Rel: "self", HREF: domainHost},
				},
			},
			server: httptest.NewServer(http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
				switch req.RequestURI {
				case fmt.Sprintf("/%s/action/removeProtectionDomain", domainHost):
					resp.WriteHeader(http.StatusOK)
				default:
					resp.WriteHeader(http.StatusNoContent)
				}
			})),
			expectedErr: nil,
		},
		"error: invalid domain links": {
			protectionDomain: &types.ProtectionDomain{},
			server: httptest.NewServer(http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
				resp.WriteHeader(http.StatusNoContent)
			})),
			expectedErr: errors.New("Error: problem finding link"),
		},
		"error: unable to remove protection domain": {
			protectionDomain: &types.ProtectionDomain{
				Links: []*types.Link{
					{Rel: "self", HREF: domainHost},
				},
			},
			server: httptest.NewServer(http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
				resp.WriteHeader(http.StatusBadRequest)
				resp.Write([]byte(`{"message":"unable to remove protection domain","httpStatusCode":400,"errorCode":0}`))
			})),
			expectedErr: errors.New("unable to remove protection domain"),
		},
	}

	for title, tc := range cases {
		fmt.Printf("running test case: %s\n", title)
		client, err := NewClientWithArgs(tc.server.URL, "3.6", math.MaxInt64, true, false)
		if err != nil {
			t.Fatal(err)
		}

		pdClient := ProtectionDomain{
			ProtectionDomain: tc.protectionDomain,
			client:           client,
		}

		err = pdClient.Delete(context.Background())
		if err != nil {
			if tc.expectedErr.Error() != err.Error() {
				t.Fatal(err)
			}
		}

		tc.server.Close()
	}
}

func TestFindProtectionDomainByName(t *testing.T) {
	domainName := "myDomain"
	domainID := "12345678-1234-1234-1234-123456789012"
	type testCase struct {
		server      *httptest.Server
		expectedErr error
	}

	cases := map[string]testCase{
		"succeed": {
			server: httptest.NewServer(http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
				switch req.RequestURI {
				case "/api/types/ProtectionDomain/instances/action/queryIdByKey":
					resp.WriteHeader(http.StatusOK)
					content, err := json.Marshal(domainID)
					if err != nil {
						t.Fatal(err)
					}

					resp.Write(content)
				case fmt.Sprintf("/api/instances/ProtectionDomain::%s", domainID):
					resp.WriteHeader(http.StatusOK)
					response := types.ProtectionDomain{
						ID:   domainID,
						Name: domainName,
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
		"error: bad request": {
			server: httptest.NewServer(http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
				switch req.RequestURI {
				case "/api/types/ProtectionDomain/instances/action/queryIdByKey":
					resp.WriteHeader(http.StatusBadRequest)
					resp.Write([]byte(`{"message":"bad request","httpStatusCode":400,"errorCode":0}`))
				default:
					resp.WriteHeader(http.StatusNoContent)
				}
			})),
			expectedErr: errors.New("error getting protection domain by name: bad request"),
		},
	}

	for title, tc := range cases {
		fmt.Printf("running test case: %s\n", title)
		client, err := NewClientWithArgs(tc.server.URL, "3.6", math.MaxInt64, true, false)
		if err != nil {
			t.Fatal(err)
		}

		s := System{
			client: client,
		}

		_, err = s.FindProtectionDomainByName(context.Background(), domainName)
		if err != nil {
			if tc.expectedErr.Error() != err.Error() {
				t.Fatal(err)
			}
		}

		tc.server.Close()
	}
}

func TestProtectionDomainRefresh(t *testing.T) {
	domainHost := "localhost"
	domainID := "12345678-1234-1234-1234-123456789012"
	type testCase struct {
		protectionDomain *types.ProtectionDomain
		server           *httptest.Server
		expectedErr      error
	}

	cases := map[string]testCase{
		"succeed": {
			protectionDomain: &types.ProtectionDomain{
				ID: domainID,
				Links: []*types.Link{
					{Rel: "self", HREF: domainHost},
				},
			},
			server: httptest.NewServer(http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
				switch req.RequestURI {
				case fmt.Sprintf("/api/instances/ProtectionDomain::%s", domainID):
					resp.WriteHeader(http.StatusOK)
				default:
					resp.WriteHeader(http.StatusNoContent)
				}
			})),
			expectedErr: nil,
		},
		"error: bad request": {
			protectionDomain: &types.ProtectionDomain{
				ID: domainID,
				Links: []*types.Link{
					{Rel: "self", HREF: domainHost},
				},
			},
			server: httptest.NewServer(http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
				switch req.RequestURI {
				case fmt.Sprintf("/api/instances/ProtectionDomain::%s", domainID):
					resp.WriteHeader(http.StatusBadRequest)
					resp.Write([]byte(`{"message":"bad request","httpStatusCode":400,"errorCode":0}`))
				default:
					resp.WriteHeader(http.StatusNoContent)
				}
			})),
			expectedErr: errors.New("bad request"),
		},
	}

	for title, tc := range cases {
		fmt.Printf("running test case: %s\n", title)
		client, err := NewClientWithArgs(tc.server.URL, "3.6", math.MaxInt64, true, false)
		if err != nil {
			t.Fatal(err)
		}

		pdClient := ProtectionDomain{
			ProtectionDomain: tc.protectionDomain,
			client:           client,
		}

		err = pdClient.Refresh(context.Background())
		if err != nil {
			if tc.expectedErr.Error() != err.Error() {
				t.Fatal(err)
			}
		}

		tc.server.Close()
	}
}

func TestProtectionDomainSetParamters(t *testing.T) {
	domainHost := "localhost"
	domainID := "12345678-1234-1234-1234-123456789012"
	ctx := context.Background()

	server := httptest.NewServer(http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
		switch req.RequestURI {
		case fmt.Sprintf("/api/instances/ProtectionDomain::%s/action/setProtectionDomainName", domainID):
			fmt.Println("handling setProtectionDomainName")
			resp.WriteHeader(http.StatusOK)
		case fmt.Sprintf("/api/instances/ProtectionDomain::%s/action/setRfcacheParameters", domainID):
			fmt.Println("handling setRfcacheParameters")
			resp.WriteHeader(http.StatusOK)
		case fmt.Sprintf("/api/instances/ProtectionDomain::%s/action/setSdsNetworkLimits", domainID):
			fmt.Println("handling setSdsNetworkLimits")
			resp.WriteHeader(http.StatusOK)
		case fmt.Sprintf("/api/instances/ProtectionDomain::%s/action/activateProtectionDomain", domainID):
			fmt.Println("handling activateProtectionDomain")
			resp.WriteHeader(http.StatusOK)
		case fmt.Sprintf("/api/instances/ProtectionDomain::%s/action/inactivateProtectionDomain", domainID):
			fmt.Println("handling inactivateProtectionDomain")
			resp.WriteHeader(http.StatusOK)
		case fmt.Sprintf("/api/instances/ProtectionDomain::%s/action/enableSdsRfcache", domainID):
			fmt.Println("handling enableSdsRfcache")
			resp.WriteHeader(http.StatusOK)
		case fmt.Sprintf("/api/instances/ProtectionDomain::%s/action/disableSdsRfcache", domainID):
			fmt.Println("handling disableSdsRfcache")
			resp.WriteHeader(http.StatusOK)
		case fmt.Sprintf("/api/instances/ProtectionDomain::%s/action/disableFglMetadataCache", domainID):
			fmt.Println("handling disableFglMetadataCache")
		case fmt.Sprintf("/api/instances/ProtectionDomain::%s/action/enableFglMetadataCache", domainID):
			fmt.Println("handling disableFglMetadataCache")
		case fmt.Sprintf("/api/instances/ProtectionDomain::%s/action/setDefaultFglMetadataCacheSize", domainID):
			fmt.Println("handling disableFglMetadataCache")
			resp.WriteHeader(http.StatusOK)
		default:
			resp.WriteHeader(http.StatusNoContent)
		}
	}))

	client, err := NewClientWithArgs(server.URL, "3.6", math.MaxInt64, true, false)
	if err != nil {
		t.Fatal(err)
	}

	pdClient := ProtectionDomain{
		ProtectionDomain: &types.ProtectionDomain{
			ID: domainID,
			Links: []*types.Link{
				{Rel: "self", HREF: domainHost},
			},
		},
		client: client,
	}

	var intPtr int = 1

	_ = pdClient.SetName(ctx, "myDomain")
	_ = pdClient.SetRfcacheParams(ctx, types.PDRfCacheParams{
		RfCachePageSizeKb:  16,
		RfCacheMaxIoSizeKb: 64,
	})
	_ = pdClient.SetSdsNetworkLimits(ctx, types.SdsNetworkLimitParams{
		RebuildNetworkThrottlingInKbps: &intPtr,
	})
	_ = pdClient.Activate(ctx, true)
	_ = pdClient.InActivate(ctx, true)
	_ = pdClient.EnableRfcache(ctx)
	_ = pdClient.DisableRfcache(ctx)
	_ = pdClient.DisableFGLMcache(ctx)
	_ = pdClient.EnableFGLMcache(ctx)
	_ = pdClient.SetDefaultFGLMcacheSize(ctx, 128)
}
