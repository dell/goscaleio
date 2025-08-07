// Copyright © 2020 - 2024 Dell Inc. or its subsidiaries. All Rights Reserved.
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
	"fmt"
	"math"
	"net/http"
	"net/http/httptest"
	"testing"

	types "github.com/dell/goscaleio/types/v1"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func Test_FindVolumes(t *testing.T) {
	type checkFn func(*testing.T, []*Volume, error)
	check := func(fns ...checkFn) []checkFn { return fns }

	hasNoError := func(t *testing.T, _ []*Volume, err error) {
		if err != nil {
			t.Fatalf("expected no error")
		}
	}

	checkLength := func(length int) func(t *testing.T, vols []*Volume, err error) {
		return func(t *testing.T, vols []*Volume, _ error) {
			assert.Equal(t, length, len(vols))
		}
	}

	hasError := func(t *testing.T, _ []*Volume, err error) {
		if err == nil {
			t.Fatalf("expected error")
		}
	}

	tests := map[string]func(t *testing.T) (*httptest.Server, types.Sdc, []checkFn){
		"success": func(t *testing.T) (*httptest.Server, types.Sdc, []checkFn) {
			sdcID := "000001111a2222b"
			href := fmt.Sprintf("/api/instances/Sdc::%s/relationships/Volume", sdcID)
			sdc := types.Sdc{
				ID: sdcID,
				Links: []*types.Link{
					{
						Rel:  "/api/Sdc/relationship/Volume",
						HREF: href,
					},
				},
			}
			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method != http.MethodGet {
					t.Fatal(fmt.Errorf("wrong method. Expected %s; but got %s", http.MethodGet, r.Method))
				}
				if r.URL.Path != href {
					t.Fatal(fmt.Errorf("wrong path. Expected %s; but got %s", href, r.URL.Path))
				}
				vols := []types.Volume{{}, {}, {}}
				respData, err := json.Marshal(vols)
				if err != nil {
					t.Fatal(err)
				}
				fmt.Fprintln(w, string(respData))
			}))
			return ts, sdc, check(hasNoError, checkLength(3))
		},
		"error from GetVolume": func(t *testing.T) (*httptest.Server, types.Sdc, []checkFn) {
			sdcID := "someID"
			href := fmt.Sprintf("/api/instances/Sdc::%s/relationships/Volume", sdcID)
			sdc := types.Sdc{
				ID: sdcID,
				Links: []*types.Link{
					{
						Rel:  "/api/Sdc/relationship/Volume",
						HREF: href,
					},
				},
			}
			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method != http.MethodGet {
					t.Fatal(fmt.Errorf("wrong method. Expected %s; but got %s", http.MethodGet, r.Method))
				}
				if r.URL.Path != href {
					t.Fatal(fmt.Errorf("wrong path. Expected %s; but got %s", href, r.URL.Path))
				}
				http.NotFound(w, r)
			}))
			return ts, sdc, check(hasError)
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			ts, sdc, checkFns := tc(t)
			defer ts.Close()

			client, err := NewClientWithArgs(ts.URL, "", math.MaxInt64, true, false)
			if err != nil {
				t.Fatal(err)
			}

			sdcClient := NewSdc(client, &sdc)
			vols, err := sdcClient.FindVolumes()
			for _, checkFn := range checkFns {
				checkFn(t, vols, err)
			}
		})
	}
}

func TestRenameSdc(t *testing.T) {
	type testCase struct {
		sdcID    string
		name     string
		expected error
	}
	cases := []testCase{
		{
			"c4270bf500000053",
			"worker-node-2345",
			nil,
		},
		{
			"c4270bf500000053",
			"",
			errors.New("Request message is not valid: The following parameter(s) must be part of the request body: sdcName"),
		},
		{
			"c4270bf500000053",
			" ",
			errors.New("The given name contains invalid characters. Use alphanumeric and punctuation characters only. Spaces are not allowed"),
		},
		{
			"worker-node-2",
			"c4270bf500000053",
			errors.New("id (worker-node-2) must be a hexadecimal number (unsigned long)"),
		},
	}

	// mock a powerflex endpoint
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

			// calling RenameSdc with mock value
			err = client.RenameSdc(tc.sdcID, tc.name)
			if err != nil {
				if tc.expected == nil {
					t.Errorf("Renaming sdc did not work as expected, \n\tgot: %s \n\twant: %v", err, tc.expected)
				} else {
					if err.Error() != tc.expected.Error() {
						t.Errorf("Renaming sdc did not work as expected, \n\tgot: %s \n\twant: %s", err, tc.expected)
					}
				}
			}
		})
	}
}

func TestApproveSDC(t *testing.T) {
	type testCase struct {
		param    types.ApproveSdcParam
		expected error
	}

	system := types.System{
		ID:                       "0000aaabbbccc1111",
		RestrictedSdcModeEnabled: true,
		RestrictedSdcMode:        "Guid",
	}

	cases := []testCase{
		{
			param: types.ApproveSdcParam{
				SdcGUID: "UT3A9C7B2E-7F2D-4A1B-9C3E-8A1FDE9B1234",
				SdcIP:   "10.10.10.10",
				SdcIps:  []string{"10.10.10.10"},
			},
			expected: nil,
		},
		{
			param: types.ApproveSdcParam{
				Name: "sdc_test",
			},
			expected: errors.New("Request message is not valid: One of the parameter(s) must be part of the request body: sdcIp, sdcIps, sdcGuid"),
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
				System: &system,
			}

			_, err2 := s.ApproveSdc(&tc.param)
			if err2 != nil {
				if tc.expected == nil {
					t.Errorf("Approving SDC did not work as expected, \n\tgot: %s \n\twant: %v", err2, tc.expected)
				} else {
					if err2.Error() != tc.expected.Error() {
						t.Errorf("Approving SDC did not work as expected, \n\tgot: %s \n\twant: %s", err2, tc.expected)
					}
				}
			}
		})
	}
}

func TestSetRestrictedMode(t *testing.T) {
	type testCase struct {
		mode     string
		expected error
	}

	system := types.System{
		ID:                       "0000aaabbbccc1111",
		RestrictedSdcModeEnabled: true,
		RestrictedSdcMode:        "Guid",
	}

	cases := []testCase{
		{
			mode:     "Guid",
			expected: nil,
		},
		{
			mode:     "random",
			expected: errors.New("restrictedSdcMode should get one of the following values: None, Guid, ApprovedIp, but its value is random"),
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
				System: &system,
			}

			err2 := s.SetRestrictedMode(tc.mode)
			if err2 != nil {
				if tc.expected == nil {
					t.Errorf("Modifying restricted mode did not work as expected, \n\tgot: %s \n\twant: %v", err2, tc.expected)
				} else {
					if err2.Error() != tc.expected.Error() {
						t.Errorf("Modifying restricted mode did not work as expected, \n\tgot: %s \n\twant: %s", err2, tc.expected)
					}
				}
			}
		})
	}
}

func TestSetApprovedIps(t *testing.T) {
	type testCase struct {
		SdcID    string
		SdcIps   []string
		expected error
	}

	system := types.System{
		ID:                       "0000aaabbbccc1111",
		RestrictedSdcModeEnabled: true,
		RestrictedSdcMode:        "Guid",
	}

	cases := []testCase{
		{
			SdcID:    "5e662b1700000000",
			SdcIps:   []string{"10.10.10.10"},
			expected: nil,
		},
		{
			SdcIps:   []string{"10.10.10.10"},
			expected: errors.New("Request message is not valid: The following parameter(s) must be part of the request body: sdcId"),
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
				System: &system,
			}

			err2 := s.SetApprovedIps(tc.SdcID, tc.SdcIps)
			if err2 != nil {
				if tc.expected == nil {
					t.Errorf("Approving SDC IPs did not work as expected, \n\tgot: %s \n\twant: %v", err2, tc.expected)
				} else {
					if err2.Error() != tc.expected.Error() {
						t.Errorf("Approving SDC IPs did not work as expected, \n\tgot: %s \n\twant: %s", err2, tc.expected)
					}
				}
			}
		})
	}
}

func TestDeleteSdc(t *testing.T) {
	type testCase struct {
		ID       string
		expected error
	}

	system := types.System{
		ID:                       "0000aaabbbccc1111",
		RestrictedSdcModeEnabled: true,
		RestrictedSdcMode:        "Guid",
	}

	cases := []testCase{
		{
			ID:       "127.0.0.1",
			expected: nil,
		},
		{
			ID:       "10.10.10.10",
			expected: errors.New("Request message is not valid: The following parameter(s) must be part of the request body: sdcId"),
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
				System: &system,
			}

			err2 := s.DeleteSdc(tc.ID)
			if err2 != nil {
				if tc.expected == nil {
					t.Errorf("Approving SDC IPs did not work as expected, \n\tgot: %s \n\twant: %v", err2, tc.expected)
				} else {
					if err2.Error() != tc.expected.Error() {
						t.Errorf("Approving SDC IPs did not work as expected, \n\tgot: %s \n\twant: %s", err2, tc.expected)
					}
				}
			}
		})
	}
}

func TestGetSdcIDByIP(t *testing.T) {
	ip := "127.0.0.1"
	systemID := uuid.NewString()
	type testCase struct {
		server        *httptest.Server
		expectedErr   error
		finalResponse string
	}

	cases := map[string]testCase{
		"succeed": {
			finalResponse: "123",
			server: httptest.NewServer(http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
				switch req.RequestURI {
				case "/api/types/Sdc/instances/action/queryIdByKey":
					resp.Write([]byte("123"))
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
				ID: systemID,
			},
			client: client,
		}

		response, err := system.GetSdcIDByIP(ip)
		if err != nil {
			if tc.expectedErr.Error() != err.Error() {
				t.Fatal(err)
			}
		} else {
			if tc.finalResponse != response {
				t.Fatal(err)
			}
		}
	}
}

func TestFindSdc(t *testing.T) {
	systemID := uuid.NewString()
	searchSdcID := uuid.NewString()
	testSdc := []types.Sdc{
		{
			Name: "FirstTest",
			ID:   searchSdcID,
		},
		{
			Name: "SecondTest",
			ID:   searchSdcID,
		},
	}
	type testCase struct {
		server      *httptest.Server
		expectedErr error
	}

	cases := map[string]testCase{
		"succeed": {
			server: httptest.NewServer(http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
				switch req.RequestURI {
				case fmt.Sprintf("/api/instances/System::%v/relationships/Sdc", systemID):
					content, err := json.Marshal(testSdc)
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
		"error: could not find sdc": {
			server: httptest.NewServer(http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
				switch req.RequestURI {
				case fmt.Sprintf("/api/instances/System::%v/relationships/Sdc", systemID):
					content, err := json.Marshal([]types.Sdc{
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
			expectedErr: fmt.Errorf("Couldn't find SDC"),
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
				ID: systemID,
			},
			client: client,
		}

		res, err := system.FindSdc("ID", searchSdcID)
		if err != nil {
			if tc.expectedErr.Error() != err.Error() {
				t.Fatal(err)
			}
		} else {
			if res.Sdc.ID != searchSdcID {
				t.Fatal(err)
			}
		}
	}
}

func TestGetSdcByID(t *testing.T) {
	systemID := uuid.NewString()
	searchSdcID := uuid.NewString()
	testSdc := types.Sdc{
		ID: searchSdcID,
	}
	type testCase struct {
		server      *httptest.Server
		expectedErr error
	}

	cases := map[string]testCase{
		"succeed": {
			server: httptest.NewServer(http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
				switch req.RequestURI {
				case fmt.Sprintf("/api/instances/Sdc::%v", searchSdcID):
					content, err := json.Marshal(testSdc)
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

		system := System{
			System: &types.System{
				ID: systemID,
			},
			client: client,
		}

		res, err := system.GetSdcByID(searchSdcID)
		if err != nil {
			if tc.expectedErr.Error() != err.Error() {
				t.Fatal(err)
			}
		} else {
			if res.Sdc.ID != searchSdcID {
				t.Fatal(err)
			}
		}
	}
}

func TestChangeSdcName(t *testing.T) {
	systemID := uuid.NewString()
	sdcName := uuid.NewString()
	searchSdcID := uuid.NewString()
	testSdc := types.Sdc{
		Name: sdcName,
		ID:   searchSdcID,
	}
	type testCase struct {
		server      *httptest.Server
		expectedErr error
	}

	cases := map[string]testCase{
		"succeed": {
			server: httptest.NewServer(http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
				switch req.RequestURI {
				case fmt.Sprintf("/api/instances/Sdc::%v/action/setSdcName", searchSdcID):
					content, err := json.Marshal(testSdc)
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

		system := System{
			System: &types.System{
				ID: systemID,
			},
			client: client,
		}

		res, err := system.ChangeSdcName(searchSdcID, "NameSDC")
		if err != nil {
			if tc.expectedErr.Error() != err.Error() {
				t.Fatal(err)
			}
		} else {
			if res.Sdc.Name != sdcName {
				t.Fatal(err)
			}
		}
	}
}

func TestChangeSdcPerfProfile(t *testing.T) {
	systemID := uuid.NewString()
	searchsdcID := uuid.NewString()
	perfProfile := uuid.NewString()
	testSdc := types.Sdc{
		ID:          searchsdcID,
		PerfProfile: perfProfile,
	}
	type testCase struct {
		server      *httptest.Server
		expectedErr error
	}

	cases := map[string]testCase{
		"succeed": {
			server: httptest.NewServer(http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
				switch req.RequestURI {
				case fmt.Sprintf("/api/instances/Sdc::%v/action/setSdcPerformanceParameters", searchsdcID):
					content, err := json.Marshal(testSdc)
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

		system := System{
			System: &types.System{
				ID: systemID,
			},
			client: client,
		}

		res, err := system.ChangeSdcPerfProfile(searchsdcID, "Perf Profile")
		if err != nil {
			if tc.expectedErr.Error() != err.Error() {
				t.Fatal(err)
			}
		} else {
			if res.Sdc.PerfProfile != perfProfile {
				t.Fatal(err)
			}
		}
	}
}

func TestGetSatistics(t *testing.T) {
	sdcID := uuid.NewString()
	mapVolumes := 3
	testSdc := types.SdcStatistics{
		NumOfMappedVolumes: mapVolumes,
	}
	type testCase struct {
		server      *httptest.Server
		expectedErr error
		sdc         types.Sdc
	}

	cases := map[string]testCase{
		"succeed": {
			sdc: types.Sdc{
				ID: uuid.NewString(),
				Links: []*types.Link{
					{
						Rel:  "/api/Sdc/relationship/Statistics",
						HREF: fmt.Sprintf("/api/instances/Sdc::%s/relationships/Statistics", sdcID),
					},
				},
			},
			server: httptest.NewServer(http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
				switch req.RequestURI {
				case fmt.Sprintf("/api/instances/Sdc::%s/relationships/Statistics", sdcID):
					resp.WriteHeader(http.StatusOK)
					content, err := json.Marshal(testSdc)
					if err != nil {
						t.Fatalf("failed to marshal volume: %v", err)
					}
					resp.Write(content)
				default:
					resp.WriteHeader(http.StatusBadRequest)
					resp.Write([]byte(`{"message":"no route handled","httpStatusCode":400,"errorCode":0}`))
				}
			})),
			expectedErr: nil,
		},
		"error: bad request": {
			sdc: types.Sdc{
				ID: uuid.NewString(),
				Links: []*types.Link{
					{
						Rel:  "/api/Sdc/relationship/Statistics",
						HREF: fmt.Sprintf("/api/instances/Sdc::%s/relationships/Statistics", sdcID),
					},
				},
			},
			server: httptest.NewServer(http.HandlerFunc(func(resp http.ResponseWriter, _ *http.Request) {
				resp.WriteHeader(http.StatusBadRequest)
				resp.Write([]byte(`{"message":"bad request","httpStatusCode":400,"errorCode":0}`))
			})),
			expectedErr: fmt.Errorf("bad request"),
		},
		"error: bad link": {
			sdc: types.Sdc{
				ID: uuid.NewString(),
			},
			server: httptest.NewServer(http.HandlerFunc(func(resp http.ResponseWriter, _ *http.Request) {
				resp.WriteHeader(http.StatusBadRequest)
				resp.Write([]byte(`{"message":"bad request","httpStatusCode":400,"errorCode":0}`))
			})),
			expectedErr: fmt.Errorf("Error: problem finding link"),
		},
	}

	for _, tc := range cases {
		client, err := NewClientWithArgs(tc.server.URL, "3.6", math.MaxInt64, true, false)
		if err != nil {
			t.Fatal(err)
		}
		defer tc.server.Close()

		sdc := Sdc{
			Sdc:    &tc.sdc,
			client: client,
		}

		res, err := sdc.GetStatistics()
		if err != nil {
			if tc.expectedErr.Error() != err.Error() {
				t.Fatal(err)
			}
		} else {
			if res.NumOfMappedVolumes != mapVolumes {
				t.Fatal(err)
			}
		}
	}
}

func TestMapVolumeSdc(t *testing.T) {
	volumeID := uuid.NewString()
	mapVolumeSdcParam := &types.MapVolumeSdcParam{}
	type testCase struct {
		server      *httptest.Server
		expectedErr error
	}

	cases := map[string]testCase{
		"succeed": {
			server: httptest.NewServer(http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
				switch req.RequestURI {
				case fmt.Sprintf("/api/instances/Volume::%s/action/addMappedSdc", volumeID):
					content, err := json.Marshal(types.Sdc{
						ID: volumeID,
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
		"error: could not find sdc": {
			server: httptest.NewServer(http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
				switch req.RequestURI {
				case fmt.Sprintf("/api/instances/Volume::%s/action/addMappedSdc", volumeID):
					content, err := json.Marshal([]types.Sdc{
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
			expectedErr: fmt.Errorf("Couldn't find SDC"),
		},
	}

	for _, tc := range cases {
		client, err := NewClientWithArgs(tc.server.URL, "3.6", math.MaxInt64, true, false)
		if err != nil {
			t.Fatal(err)
		}
		defer tc.server.Close()

		volume := Volume{
			Volume: &types.Volume{
				ID: volumeID,
			},
			client: client,
		}

		err = volume.MapVolumeSdc(mapVolumeSdcParam)
		if err != nil {
			if tc.expectedErr.Error() != err.Error() {
				t.Fatal(err)
			}
		}
	}
}

func TestUnMapVolumeSdc(t *testing.T) {
	volumeID := uuid.NewString()
	unMapVolumeSdcParam := &types.UnmapVolumeSdcParam{}
	type testCase struct {
		server      *httptest.Server
		expectedErr error
	}

	cases := map[string]testCase{
		"succeed": {
			server: httptest.NewServer(http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
				switch req.RequestURI {
				case fmt.Sprintf("/api/instances/Volume::%s/action/removeMappedSdc", volumeID):
					content, err := json.Marshal(types.Sdc{
						ID: volumeID,
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

		volume := Volume{
			Volume: &types.Volume{
				ID: volumeID,
			},
			client: client,
		}

		err = volume.UnmapVolumeSdc(unMapVolumeSdcParam)
		if err != nil {
			if tc.expectedErr.Error() != err.Error() {
				t.Fatal(err)
			}
		}
	}
}

func TestSetMappedSdcLimits(t *testing.T) {
	volumeID := uuid.NewString()
	setMapVolumeSdcParam := &types.SetMappedSdcLimitsParam{}
	type testCase struct {
		server      *httptest.Server
		expectedErr error
	}

	cases := map[string]testCase{
		"succeed": {
			server: httptest.NewServer(http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
				switch req.RequestURI {
				case fmt.Sprintf("/api/instances/Volume::%s/action/setMappedSdcLimits", volumeID):
					content, err := json.Marshal(types.Sdc{
						ID: volumeID,
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
		"error: could not find sdc": {
			server: httptest.NewServer(http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
				switch req.RequestURI {
				case fmt.Sprintf("/api/instances/Volume::%s/action/setMappedSdcLimits", volumeID):
					content, err := json.Marshal([]types.Sdc{
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
			expectedErr: fmt.Errorf("Couldn't find SDC"),
		},
	}

	for _, tc := range cases {
		client, err := NewClientWithArgs(tc.server.URL, "3.6", math.MaxInt64, true, false)
		if err != nil {
			t.Fatal(err)
		}
		defer tc.server.Close()

		volume := Volume{
			Volume: &types.Volume{
				ID: volumeID,
			},
			client: client,
		}

		err = volume.SetMappedSdcLimits(setMapVolumeSdcParam)
		if err != nil {
			if tc.expectedErr.Error() != err.Error() {
				t.Fatal(err)
			}
		}
	}
}

func TestGetSdcLocalGUID(t *testing.T) {
	// Test case: successful execution
	expectedGUID := "271bad82-08ee-44f2-a2b1-7e2787c27be1"
	out := []byte(expectedGUID + "\n")
	execCmdOriginal := execCmd
	execCmd = func(_ string, _ ...string) ([]byte, error) {
		return out, nil
	}
	defer func() { execCmd = execCmdOriginal }()

	guid, err := GetSdcLocalGUID()
	assert.NoError(t, err)
	assert.Equal(t, expectedGUID, guid)

	// Test case: error in exec.Command
	expectedErr := errors.New("exec.Command failed")
	execCmd = func(_ string, _ ...string) ([]byte, error) {
		return nil, expectedErr
	}

	_, err = GetSdcLocalGUID()
	assert.Error(t, err)
}

func TestGetVolumeMetrics(t *testing.T) {
	// Test case 1: Successful retrieval of volume metrics
	t.Run("successful retrieval", func(t *testing.T) {
		// Create a test server with a successful response
		svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodPost {
				t.Errorf("expected POST method, got %s", r.Method)
			}

			if r.URL.Path != "/api/instances/Sdc::ed10ad4300000031/action/queryVolumeSdcBwc" {
				t.Errorf("expected URL path /api/instances/Sdc::ed10ad4300000031/action/queryVolumeSdcBwc, got %s", r.URL.Path)
			}

			w.WriteHeader(http.StatusOK)
			jsonData := []byte(`
				[
				  {
					"readLatencyBwc": {
					  "numSeconds": 0,
					  "totalWeightInKb": 0,
					  "numOccured": 0
					},
					"volumeId": "9d12552300000039",
					"trimBwc": {
					  "numSeconds": 0,
					  "totalWeightInKb": 0,
					  "numOccured": 0
					},
					"trimLatencyBwc": {
					  "numSeconds": 0,
					  "totalWeightInKb": 0,
					  "numOccured": 0
					},
					"readBwc": {
					  "numSeconds": 0,
					  "totalWeightInKb": 0,
					  "numOccured": 0
					},
					"writeLatencyBwc": {
					  "numSeconds": 0,
					  "totalWeightInKb": 0,
					  "numOccured": 0
					},
					"sdcId": "ed10ad4300000031",
					"writeBwc": {
					  "numSeconds": 0,
					  "totalWeightInKb": 0,
					  "numOccured": 0
					}
				  }
				]
			`)
			w.Write(jsonData)
		}))
		defer svr.Close()

		client, err := NewClientWithArgs(svr.URL, "3.6", math.MaxInt64, true, false)
		if err != nil {
			t.Fatal(err)
		}

		// Create a new Sdc instance with the test server URL
		sdc := &Sdc{
			Sdc: &types.Sdc{
				ID: "ed10ad4300000031",
			},
			client: client,
		}

		// Call the GetVolumeMetrics method
		metrics, err := sdc.GetVolumeMetrics()
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}

		if len(metrics) != 1 {
			t.Errorf("expected 1 volume info, got %d", len(metrics))
		}

		if metrics[0].VolumeID != "9d12552300000039" {
			t.Errorf("expected volume ID 9d12552300000039, got %s", metrics[0].VolumeID)
		}
	})

	// Test case 2: Error during retrieval of volume metrics
	t.Run("error during retrieval", func(t *testing.T) {
		// Create a test server with an error response
		svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(`{"error": "mock error"}`))
		}))
		defer svr.Close()

		client, err := NewClientWithArgs(svr.URL, "3.6", math.MaxInt64, true, false)
		if err != nil {
			t.Fatal(err)
		}

		// Create a new Sdc instance with the test server URL
		sdc := &Sdc{
			Sdc: &types.Sdc{
				ID: "ed10ad4300000031",
			},
			client: client,
		}

		// Call the GetVolumeMetrics method
		_, err = sdc.GetVolumeMetrics()
		if err == nil {
			t.Errorf("expected an error, got none")
		}
	})

	// Test case 3: Invalid JSON response
	t.Run("invalid JSON response", func(t *testing.T) {
		// Create a test server with an invalid JSON response
		svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`invalid JSON`))
		}))
		defer svr.Close()

		client, err := NewClientWithArgs(svr.URL, "3.6", math.MaxInt64, true, false)
		if err != nil {
			t.Fatal(err)
		}

		// Create a new Sdc instance with the test server URL
		sdc := &Sdc{
			Sdc: &types.Sdc{
				ID: "ed10ad4300000031",
			},
			client: client,
		}

		// Call the GetVolumeMetrics method
		_, err = sdc.GetVolumeMetrics()
		if err == nil {
			t.Errorf("expected an error, got none")
		}
	})
}
