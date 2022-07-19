package goscaleio

import (
	"encoding/json"
	"errors"
	"fmt"
	types "github.com/dell/goscaleio/types/v1"
	"net/http"
	"time"
)

// PeerMDM encpsulates a PeerMDM type and a client.
type PeerMDM struct {
	PeerMDM *types.PeerMDM
	client  *Client
}

// NewPeerMDM creates a PeerMDM from a types.PeerMDM and a client.
func NewPeerMDM(client *Client, peerMDM *types.PeerMDM) *PeerMDM {
	newPeerMDM := &PeerMDM{
		client:  client,
		PeerMDM: peerMDM,
	}
	return newPeerMDM
}

// GetPeerMDMs returns a list of peer MDMs know to the System
func (c *Client) GetPeerMDMs() ([]*types.PeerMDM, error) {
	path := "/api/types/PeerMdm/instances"
	var peerMdms []*types.PeerMDM
	var err error
	defer TimeSpent("GetPeerMDMs", time.Now())

	err = c.getJSONWithRetry(http.MethodGet, path, nil, &peerMdms)
	return peerMdms, err
}

// ReplicationConsistencyGroup encpsulates a types.ReplicationConsistencyGroup and a client.
type ReplicationConsistencyGroup struct {
	ReplicationConsistencyGroup *types.ReplicationConsistencyGroup
	client                      *Client
}

// NewReplicationConsistencyGroup creates a new ReplicationConsistencyGroup.
func NewReplicationConsistencyGroup(client *Client) *ReplicationConsistencyGroup {
	rcg := &ReplicationConsistencyGroup{
		client:                      client,
		ReplicationConsistencyGroup: &types.ReplicationConsistencyGroup{},
	}
	return rcg
}

// GetReplicationConsistencyGroups returns a list of the ReplicationConsistencyGroups
func (c *Client) GetReplicationConsistencyGroups() ([]*types.ReplicationConsistencyGroup, error) {
	path := "/api/types/ReplicationConsistencyGroup/instances"
	var rcgs []*types.ReplicationConsistencyGroup
	var err error
	defer TimeSpent("GetReplicationConsistencyGroups", time.Now())

	err = c.getJSONWithRetry(http.MethodGet, path, nil, &rcgs)
	return rcgs, err
}

// CreateReplicationConsistencyGroup
func (c *Client) CreateReplicationConsistencyGroup(rcg *types.ReplicationConsistencyGroup) (*ReplicationConsistencyGroup, error) {
	if rcg.RpoInSeconds == 0 || rcg.ProtectionDomainId == "" || rcg.RemoteProtectionDomainId == "" {
		return nil, errors.New("RpoInSeconds, ProtectionDomainId, and RemoteProtectionDomainId are required")
	}
	if rcg.DestinationSystemId == "" && rcg.PeerMdmId == "" {
		return nil, errors.New("Either DestinationSystemId or PeerMdmId are required")
	}
	bytes, err := json.Marshal(rcg)
	if err != nil {
		fmt.Printf("Marshal error: %s\n", err)
	}
	fmt.Printf("Marshal output: %s\n", string(bytes))
	defer TimeSpent("CreateReplicationConsistencyGroup", time.Now())
	path := "/api/types/ReplicationConsistencyGroup/instances"
	rcgResp := NewReplicationConsistencyGroup(c)

	err = c.getJSONWithRetry(http.MethodPost, path, rcg, &rcgResp.ReplicationConsistencyGroup)
	if err != nil {
		return nil, err
	}
	return rcgResp, nil
}

// GetReplicationPairs returns a list of ReplicationPair objects. If a ReplicationConsistencyGroupId is specified, will be limited to paris of that RCG.
func (c *Client) GetReplicationPairs(RCGId string) ([]*types.ReplicationPair, error) {
	path := "/api/types/ReplicationPair/instances"
	var err error
	var pairs []*types.ReplicationPair
	defer TimeSpent("GetReplicationPairs", time.Now())
	if RCGId != "" {
		path = "/api/instances/ReplicationConsistencyGroup::" + RCGId + "/relationships/ReplicationPair"
	}
	err = c.getJSONWithRetry(http.MethodGet, path, nil, &pairs)
	return pairs, err
}
