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

func TestGetPeerMdm(t *testing.T) {
	searchID := uuid.NewString()
	peerMdms := []*types.PeerMDM{
		{
			Name: "firstPeer",
			ID:   searchID,
		},
		{
			Name: "secondPeer",
			ID:   uuid.NewString(),
		},
	}

	server := httptest.NewServer(http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
		switch req.RequestURI {
		case "/api/types/PeerMdm/instances":
			resp.WriteHeader(http.StatusOK)
			content, err := json.Marshal(peerMdms)
			if err != nil {
				t.Fatal(err)
			}

			resp.Write(content)
		case fmt.Sprintf("/api/instances/PeerMdm::%s", searchID):
			resp.WriteHeader(http.StatusOK)
			var peerMdm types.PeerMDM
			for _, val := range peerMdms {
				if val.ID == searchID {
					peerMdm = *val
				}
			}

			content, err := json.Marshal(peerMdm)
			if err != nil {
				t.Fatal(err)
			}

			resp.Write(content)
		default:
			resp.WriteHeader(http.StatusNoContent)
		}
	}))
	defer server.Close()

	client, err := NewClientWithArgs(server.URL, "3.6", math.MaxInt64, true, false)
	if err != nil {
		t.Fatal(err)
	}

	response, err := client.GetPeerMDMs()
	if err != nil {
		t.Fatal(err)
	}

	if len(response) != 2 {
		t.Errorf("expected %d, got %d", 2, len(response))
	}

	res, err := client.GetPeerMDM(searchID)
	if err != nil || res == nil {
		t.Fatal(err)
	}
}

func TestModifyPeerMdm(t *testing.T) {
	searchID := uuid.NewString()

	server := httptest.NewServer(http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
		switch req.RequestURI {
		case fmt.Sprintf("/api/instances/PeerMdm::%s/action/modifyPeerMdmIp", searchID):
			fmt.Printf("modifyPeerMdmIp for %s\n", searchID)
			resp.WriteHeader(http.StatusOK)
		case fmt.Sprintf("/api/instances/PeerMdm::%s/action/modifyPeerMdmName", searchID):
			fmt.Printf("modifyPeerMdmName for %s\n", searchID)
			resp.WriteHeader(http.StatusOK)
		case fmt.Sprintf("/api/instances/PeerMdm::%s/action/modifyPeerMdmPort", searchID):
			fmt.Printf("modifyPeerMdmPort for %s\n", searchID)
			resp.WriteHeader(http.StatusOK)
		case fmt.Sprintf("/api/instances/PeerMdm::%s/action/setPeerMdmPerformanceParameters", searchID):
			fmt.Printf("setPeerMdmPerformanceParameters for %s\n", searchID)
			resp.WriteHeader(http.StatusOK)
		default:
			resp.WriteHeader(http.StatusNoContent)
		}
	}))
	defer server.Close()

	client, err := NewClientWithArgs(server.URL, "3.6", math.MaxInt64, true, false)
	if err != nil {
		t.Fatal(err)
	}

	ips := []string{"127.0.0.1", "127.0.0.2"}
	err = client.ModifyPeerMdmIP(searchID, ips)
	if err != nil {
		t.Fatal(err)
	}

	modifyName := types.ModifyPeerMDMNameParam{
		NewName: "newPeerName",
	}
	err = client.ModifyPeerMdmName(searchID, &modifyName)
	if err != nil {
		t.Fatal(err)
	}

	modifyPort := types.ModifyPeerMDMPortParam{
		NewPort: "newPort",
	}
	err = client.ModifyPeerMdmPort(searchID, &modifyPort)
	if err != nil {
		t.Fatal(err)
	}

	modifyPerf := types.ModifyPeerMdmPerformanceParametersParam{
		NewPreformanceProfile: "Compact",
	}
	err = client.ModifyPeerMdmPerformanceParameters(searchID, &modifyPerf)
	if err != nil {
		t.Fatal(err)
	}
}

func TestAddPeerMdm(t *testing.T) {
	type testCase struct {
		addPeerMdm  *types.AddPeerMdm
		server      *httptest.Server
		expectedErr error
	}

	cases := map[string]testCase{
		"succeed": {
			addPeerMdm: &types.AddPeerMdm{
				PeerSystemID:  uuid.NewString(),
				PeerSystemIps: []string{"127.0.0.1", "127.0.0.2"},
			},
			server: httptest.NewServer(http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
				switch req.RequestURI {
				case "/api/types/PeerMdm/instances":
					resp.WriteHeader(http.StatusOK)
				default:
					resp.WriteHeader(http.StatusNoContent)
				}
			})),
			expectedErr: nil,
		},
		"error: no peer system ID": {
			addPeerMdm: &types.AddPeerMdm{},
			server: httptest.NewServer(http.HandlerFunc(func(resp http.ResponseWriter, _ *http.Request) {
				resp.WriteHeader(http.StatusNoContent)
			})),
			expectedErr: errors.New("PeerSystemID and PeerSystemIps are required"),
		},
		"error: bad request": {
			addPeerMdm: &types.AddPeerMdm{
				PeerSystemID:  uuid.NewString(),
				PeerSystemIps: []string{"127.0.0.1", "127.0.0.2"},
			},
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

		_, err = client.AddPeerMdm(tc.addPeerMdm)
		if err != nil {
			if tc.expectedErr.Error() != err.Error() {
				t.Fatal(err)
			}
		}

		tc.server.Close()
	}
}

func TestRemovePeerMdm(t *testing.T) {
	peerID := uuid.NewString()
	type testCase struct {
		server      *httptest.Server
		expectedErr error
	}

	cases := map[string]testCase{
		"succeed": {
			server: httptest.NewServer(http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
				switch req.RequestURI {
				case fmt.Sprintf("/api/instances/PeerMdm::%s/action/removePeerMdm", peerID):
					resp.WriteHeader(http.StatusOK)
				default:
					resp.WriteHeader(http.StatusNoContent)
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

		err = client.RemovePeerMdm(peerID)
		if err != nil {
			if tc.expectedErr.Error() != err.Error() {
				t.Fatal(err)
			}
		}

		tc.server.Close()
	}
}

func TestGetReplicationConsistencyGroup(t *testing.T) {
	searchID := uuid.NewString()
	replicationGroups := []*types.ReplicationConsistencyGroup{
		{
			Name: "firstReplicationGroup",
			ID:   searchID,
		},
		{
			Name: "secondReplicationGroup",
			ID:   uuid.NewString(),
		},
	}

	server := httptest.NewServer(http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
		switch req.RequestURI {
		case "/api/types/ReplicationConsistencyGroup/instances":
			resp.WriteHeader(http.StatusOK)
			content, err := json.Marshal(replicationGroups)
			if err != nil {
				t.Fatal(err)
			}

			resp.Write(content)
		case fmt.Sprintf("/api/instances/ReplicationConsistencyGroup::%s", searchID):
			resp.WriteHeader(http.StatusOK)
			var group types.ReplicationConsistencyGroup
			for _, val := range replicationGroups {
				if val.ID == searchID {
					group = *val
				}
			}

			content, err := json.Marshal(group)
			if err != nil {
				t.Fatal(err)
			}

			resp.Write(content)
		default:
			resp.WriteHeader(http.StatusNoContent)
		}
	}))
	defer server.Close()

	client, err := NewClientWithArgs(server.URL, "3.6", math.MaxInt64, true, false)
	if err != nil {
		t.Fatal(err)
	}

	groupsResponse, err := client.GetReplicationConsistencyGroups()
	if err != nil {
		t.Fatal(err)
	}

	if len(groupsResponse) != 2 {
		t.Errorf("expected %d, got %d", 2, len(groupsResponse))
	}

	res, err := client.GetReplicationConsistencyGroupByID(searchID)
	if err != nil || res == nil {
		t.Fatal(err)
	}
}

func TestCreateReplicationConsistencyGroup(t *testing.T) {
	type testCase struct {
		group       *types.ReplicationConsistencyGroupCreatePayload
		server      *httptest.Server
		expectedErr error
	}

	cases := map[string]testCase{
		"succeed": {
			group: &types.ReplicationConsistencyGroupCreatePayload{
				Name:                     "myReplicationGroup",
				RpoInSeconds:             "60",
				ProtectionDomainID:       uuid.NewString(),
				RemoteProtectionDomainID: uuid.NewString(),
				DestinationSystemID:      uuid.NewString(),
				PeerMdmID:                uuid.NewString(),
			},
			server: httptest.NewServer(http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
				switch req.RequestURI {
				case "/api/types/ReplicationConsistencyGroup/instances":
					resp.WriteHeader(http.StatusOK)
					content, err := json.Marshal(&types.ReplicationConsistencyGroupResp{
						ID: uuid.NewString(),
					})
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
		"error: missing protection domain": {
			group: &types.ReplicationConsistencyGroupCreatePayload{},
			server: httptest.NewServer(http.HandlerFunc(func(resp http.ResponseWriter, _ *http.Request) {
				resp.WriteHeader(http.StatusNoContent)
			})),
			expectedErr: errors.New("RpoInSeconds, ProtectionDomainId, and RemoteProtectionDomainId are required"),
		},
		"error: missing destination system id": {
			group: &types.ReplicationConsistencyGroupCreatePayload{
				Name:                     "myReplicationGroup",
				RpoInSeconds:             "60",
				ProtectionDomainID:       uuid.NewString(),
				RemoteProtectionDomainID: uuid.NewString(),
			},
			server: httptest.NewServer(http.HandlerFunc(func(resp http.ResponseWriter, _ *http.Request) {
				resp.WriteHeader(http.StatusNoContent)
			})),
			expectedErr: errors.New("either DestinationSystemId or PeerMdmId are required"),
		},
		"error: bad request": {
			group: &types.ReplicationConsistencyGroupCreatePayload{
				Name:                     "myReplicationGroup",
				RpoInSeconds:             "60",
				ProtectionDomainID:       uuid.NewString(),
				RemoteProtectionDomainID: uuid.NewString(),
				DestinationSystemID:      uuid.NewString(),
				PeerMdmID:                uuid.NewString(),
			},
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

		_, err = client.CreateReplicationConsistencyGroup(tc.group)
		if err != nil {
			if tc.expectedErr.Error() != err.Error() {
				t.Fatal(err)
			}
		}

		tc.server.Close()
	}
}

func TestRemoveReplicationConsistencyGroup(t *testing.T) {
	ref := "localhost"
	type testCase struct {
		group       *types.ReplicationConsistencyGroup
		server      *httptest.Server
		expectedErr error
	}

	cases := map[string]testCase{
		"succeed": {
			group: &types.ReplicationConsistencyGroup{
				Name: "myReplicationGroup",
				Links: []*types.Link{
					{Rel: "self", HREF: ref},
				},
			},
			server: httptest.NewServer(http.HandlerFunc(func(resp http.ResponseWriter, _ *http.Request) {
				resp.WriteHeader(http.StatusNoContent)
			})),
			expectedErr: nil,
		},
		"error: bad link": {
			group: &types.ReplicationConsistencyGroup{
				Name: "myReplicationGroup",
			},
			server: httptest.NewServer(http.HandlerFunc(func(resp http.ResponseWriter, _ *http.Request) {
				resp.WriteHeader(http.StatusNoContent)
			})),
			expectedErr: errors.New("Error: problem finding link"),
		},
	}

	for _, tc := range cases {
		client, err := NewClientWithArgs(tc.server.URL, "3.6", math.MaxInt64, true, false)
		if err != nil {
			t.Fatal(err)
		}

		rcg := NewReplicationConsistencyGroup(client)
		rcg.ReplicationConsistencyGroup = tc.group

		err = rcg.RemoveReplicationConsistencyGroup(true)
		if err != nil {
			if tc.expectedErr.Error() != err.Error() {
				t.Fatal(err)
			}
		}

		tc.server.Close()
	}
}

func TestCreateReplicationPair(t *testing.T) {
	type testCase struct {
		replicationPair *types.QueryReplicationPair
		server          *httptest.Server
		expectedErr     error
	}

	cases := map[string]testCase{
		"succeed": {
			replicationPair: &types.QueryReplicationPair{
				Name:                          "myReplicationPair",
				SourceVolumeID:                uuid.NewString(),
				DestinationVolumeID:           uuid.NewString(),
				ReplicationConsistencyGroupID: uuid.NewString(),
				CopyType:                      "Remote",
			},
			server: httptest.NewServer(http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
				switch req.RequestURI {
				case "/api/types/ReplicationPair/instances":
					resp.WriteHeader(http.StatusOK)
				default:
					resp.WriteHeader(http.StatusNoContent)
				}
			})),
			expectedErr: nil,
		},
		"error: missing params": {
			replicationPair: &types.QueryReplicationPair{
				Name: "myReplicationPair",
			},
			server: httptest.NewServer(http.HandlerFunc(func(resp http.ResponseWriter, _ *http.Request) {
				resp.WriteHeader(http.StatusBadRequest)
			})),
			expectedErr: errors.New("CopyType, SourceVolumeID, DestinationVolumeID, and ReplicationConsistencyGroupID are required"),
		},
		"error: bad request": {
			replicationPair: &types.QueryReplicationPair{
				Name:                          "myReplicationPair",
				SourceVolumeID:                uuid.NewString(),
				DestinationVolumeID:           uuid.NewString(),
				ReplicationConsistencyGroupID: uuid.NewString(),
				CopyType:                      "Remote",
			},
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

		_, err = client.CreateReplicationPair(tc.replicationPair)
		if err != nil {
			if tc.expectedErr.Error() != err.Error() {
				t.Fatal(err)
			}
		}

		tc.server.Close()
	}
}

func TestRemoveReplicationPair(t *testing.T) {
	rpID := uuid.NewString()
	type testCase struct {
		replicationPair *types.ReplicationPair
		server          *httptest.Server
		expectedErr     error
	}

	cases := map[string]testCase{
		"succeed": {
			replicationPair: &types.ReplicationPair{
				Name: "myReplicationPair",
				ID:   rpID,
			},
			server: httptest.NewServer(http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
				switch req.RequestURI {
				case fmt.Sprintf("/api/instances/ReplicationPair::%s/action/removeReplicationPair", rpID):
					resp.WriteHeader(http.StatusOK)
				default:
					resp.WriteHeader(http.StatusNoContent)
				}
			})),
			expectedErr: nil,
		},
		"error: bad request": {
			replicationPair: &types.ReplicationPair{
				Name: "myReplicationPair",
			},
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

		rp := NewReplicationPair(client)
		rp.ReplicaitonPair = tc.replicationPair

		_, err = rp.RemoveReplicationPair(true)
		if err != nil {
			if tc.expectedErr.Error() != err.Error() {
				t.Fatal(err)
			}
		}

		tc.server.Close()
	}
}

func TestReplicationPairActions(t *testing.T) {
	rpID := uuid.NewString()

	replicationPair := &types.ReplicationPair{
		Name: "myReplicationPair",
		ID:   rpID,
	}

	server := httptest.NewServer(http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
		switch req.RequestURI {
		case fmt.Sprintf("/api/instances/ReplicationPair::%s/relationships/Statistics", rpID):
			resp.WriteHeader(http.StatusOK)
		case fmt.Sprintf("/api/instances/ReplicationPair::%s/action/pausePairInitialCopy", rpID):
			resp.WriteHeader(http.StatusOK)
		case fmt.Sprintf("/api/instances/ReplicationPair::%s/action/resumePairInitialCopy", rpID):
			resp.WriteHeader(http.StatusOK)
		case "/api/types/ReplicationPair/instances":
			resp.WriteHeader(http.StatusOK)
		case fmt.Sprintf("/api/instances/ReplicationPair::%s", rpID):
			resp.WriteHeader(http.StatusOK)
		default:
			resp.WriteHeader(http.StatusBadRequest)
		}
	}))
	defer server.Close()

	client, err := NewClientWithArgs(server.URL, "3.6", math.MaxInt64, true, false)
	if err != nil {
		t.Fatal(err)
	}

	rp := NewReplicationPair(client)
	rp.ReplicaitonPair = replicationPair

	_, err = rp.GetReplicationPairStatistics()
	if err != nil {
		t.Fatal(err)
	}

	_, err = client.PausePairInitialCopy(rpID)
	if err != nil {
		t.Fatal(err)
	}

	_, err = client.ResumePairInitialCopy(rpID)
	if err != nil {
		t.Fatal(err)
	}

	_, err = client.GetAllReplicationPairs()
	if err != nil {
		t.Fatal(err)
	}

	_, err = client.GetReplicationPair(rpID)
	if err != nil {
		t.Fatal(err)
	}
}

func TestReplicationConsistencyGroupAction(t *testing.T) {
	groupID := uuid.NewString()

	server := httptest.NewServer(http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
		switch req.RequestURI {
		case fmt.Sprintf("/api/instances/ReplicationConsistencyGroup::%s/action/freezeApplyReplicationConsistencyGroup", groupID):
		case fmt.Sprintf("/api/instances/ReplicationConsistencyGroup::%s/relationships/ReplicationPair", groupID):
		case fmt.Sprintf("/api/instances/ReplicationConsistencyGroup::%s/action/unfreezeApplyReplicationConsistencyGroup", groupID):
		case fmt.Sprintf("/api/instances/ReplicationConsistencyGroup::%s/action/failoverReplicationConsistencyGroup", groupID):
		case fmt.Sprintf("/api/instances/ReplicationConsistencyGroup::%s/action/switchoverReplicationConsistencyGroup", groupID):
		case fmt.Sprintf("/api/instances/ReplicationConsistencyGroup::%s/action/restoreReplicationConsistencyGroup", groupID):
		case fmt.Sprintf("/api/instances/ReplicationConsistencyGroup::%s/action/reverseReplicationConsistencyGroup", groupID):
		case fmt.Sprintf("/api/instances/ReplicationConsistencyGroup::%s/action/pauseReplicationConsistencyGroup", groupID):
		case fmt.Sprintf("/api/instances/ReplicationConsistencyGroup::%s/action/resumeReplicationConsistencyGroup", groupID):
		case fmt.Sprintf("/api/instances/ReplicationConsistencyGroup::%s/action/syncNowReplicationConsistencyGroup", groupID):
		case fmt.Sprintf("/api/instances/ReplicationConsistencyGroup::%s/action/ModifyReplicationConsistencyGroupRpo", groupID):
		case fmt.Sprintf("/api/instances/ReplicationConsistencyGroup::%s/action/modifyReplicationConsistencyGroupTargetVolumeAccessMode", groupID):
		case fmt.Sprintf("/api/instances/ReplicationConsistencyGroup::%s/action/renameReplicationConsistencyGroup", groupID):
		case fmt.Sprintf("/api/instances/ReplicationConsistencyGroup::%s/action/setReplicationConsistencyGroupConsistent", groupID):
		case fmt.Sprintf("/api/instances/ReplicationConsistencyGroup::%s/action/setReplicationConsistencyGroupInconsistent", groupID):
		case fmt.Sprintf("/api/instances/ReplicationConsistencyGroup::%s/action/activateReplicationConsistencyGroup", groupID):
		case fmt.Sprintf("/api/instances/ReplicationConsistencyGroup::%s/action/terminateReplicationConsistencyGroup", groupID):
		case fmt.Sprintf("/api/instances/ReplicationConsistencyGroup::%s/action/querySyncNowReplicationConsistencyGroup", groupID):
			resp.WriteHeader(http.StatusOK)
		case fmt.Sprintf("/api/instances/ReplicationConsistencyGroup::%s/action/createReplicationConsistencyGroupSnapshots", groupID):
			fmt.Printf("createReplicationConsistencyGroupSnapshots for %s\n", groupID)
			resp.WriteHeader(http.StatusOK)
			content, err := json.Marshal(&types.CreateReplicationConsistencyGroupSnapshotResp{
				SnapshotGroupID: uuid.NewString(),
			})
			if err != nil {
				t.Fatal(err)
			}
			resp.Write(content)
		default:
			fmt.Println("Not handled. Add route.")
			resp.WriteHeader(http.StatusBadRequest)
			resp.Write([]byte(`{"message":"no route handled","httpStatusCode":400,"errorCode":0}`))
		}
	}))
	defer server.Close()

	client, err := NewClientWithArgs(server.URL, "3.6", math.MaxInt64, true, false)
	if err != nil {
		t.Fatal(err)
	}

	rcg := NewReplicationConsistencyGroup(client)
	rcg.ReplicationConsistencyGroup = &types.ReplicationConsistencyGroup{
		Name: "myReplicationGroup",
		ID:   groupID,
	}

	_, err = rcg.GetReplicationPairs()
	if err != nil {
		t.Fatal(err)
	}

	err = rcg.FreezeReplicationConsistencyGroup(groupID)
	if err != nil {
		t.Fatal(err)
	}

	err = rcg.UnfreezeReplicationConsistencyGroup()
	if err != nil {
		t.Fatal(err)
	}

	_, err = rcg.CreateReplicationConsistencyGroupSnapshot()
	if err != nil {
		t.Fatal(err)
	}

	err = rcg.ExecuteFailoverOnReplicationGroup()
	if err != nil {
		t.Fatal(err)
	}

	err = rcg.ExecuteSwitchoverOnReplicationGroup(true)
	if err != nil {
		t.Fatal(err)
	}

	err = rcg.ExecuteRestoreOnReplicationGroup()
	if err != nil {
		t.Fatal(err)
	}

	err = rcg.ExecuteReverseOnReplicationGroup()
	if err != nil {
		t.Fatal(err)
	}

	err = rcg.ExecutePauseOnReplicationGroup()
	if err != nil {
		t.Fatal(err)
	}

	err = rcg.ExecuteResumeOnReplicationGroup()
	if err != nil {
		t.Fatal(err)
	}

	_, err = rcg.ExecuteSyncOnReplicationGroup()
	if err != nil {
		t.Fatal(err)
	}

	err = rcg.SetRPOOnReplicationGroup(types.SetRPOReplicationConsistencyGroup{
		RpoInSeconds: "90",
	})
	if err != nil {
		t.Fatal(err)
	}

	err = rcg.SetTargetVolumeAccessModeOnReplicationGroup(types.SetTargetVolumeAccessModeOnReplicationGroup{
		TargetVolumeAccessMode: "ReadOnly",
	})
	if err != nil {
		t.Fatal(err)
	}

	err = rcg.SetNewNameOnReplicationGroup(types.SetNewNameOnReplicationGroup{
		NewName: "newReplicationGroup",
	})
	if err != nil {
		t.Fatal(err)
	}

	err = rcg.ExecuteConsistentOnReplicationGroup()
	if err != nil {
		t.Fatal(err)
	}

	err = rcg.ExecuteInconsistentOnReplicationGroup()
	if err != nil {
		t.Fatal(err)
	}

	err = rcg.ExecuteActivateOnReplicationGroup()
	if err != nil {
		t.Fatal(err)
	}

	err = rcg.ExecuteTerminateOnReplicationGroup()
	if err != nil {
		t.Fatal(err)
	}

	err = rcg.GetSyncStateOnReplicationGroup("syncKeyVal")
	if err != nil {
		t.Fatal(err)
	}
}

func TestNewPeerMDM(t *testing.T) {
	assert.NotNil(t, NewPeerMDM(nil, nil))
}
