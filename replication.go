package goscaleio

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	types "github.com/dell/goscaleio/types/v1"
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
func (c *Client) CreateReplicationConsistencyGroup(rcg *types.ReplicationConsistencyGroupCreatePayload) (*types.ReplicationConsistencyGroupResp, error) {
	debug = true
	showHTTP = true
	if rcg.RpoInSeconds == "" || rcg.ProtectionDomainId == "" || rcg.RemoteProtectionDomainId == "" {
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
	rcgResp := &types.ReplicationConsistencyGroupResp{}

	err = c.getJSONWithRetry(http.MethodPost, path, rcg, rcgResp)
	if err != nil {
		fmt.Printf("c.getJSONWithRetry(http.MethodPost, path, rcg, rcgResp) returned %s", err)
		return nil, err
	}
	return rcgResp, nil
}

// RemoveReplicationConsistencyGroup removes a replication consistency group
// At this point I don't know when forceIgnoreConsistency might be required.
func (rcg *ReplicationConsistencyGroup) RemoveReplicationConsistencyGroup(forceIgnoreConsistency bool) error {
	defer TimeSpent("RemoveReplicationConsistencyGroup", time.Now())

	link, err := GetLink(rcg.ReplicationConsistencyGroup.Links, "self")
	if err != nil {
		return err
	}
	path := fmt.Sprintf("%v/action/removeReplicationConsistencyGroup", link.HREF)

	removeRCGParam := &types.RemoveReplicationConsistencyGroupParam{}
	if forceIgnoreConsistency {
		removeRCGParam.ForceIgnoreConsistency = "TRUE"
	}

	err = rcg.client.getJSONWithRetry(http.MethodPost, path, removeRCGParam, nil)
	return err
}

func (c *Client) CreateReplicationPair(rp *types.QueryReplicationPair) (*types.ReplicationPair, error) {
	debug = true
	showHTTP = true
	if rp.CopyType == "" || rp.SourceVolumeID == "" || rp.DestinationVolumeID == "" || rp.ReplicationConsistencyGroupID == "" {
		return nil, errors.New("CopyType, SourceVolumeID, DestinationVolumeID, and ReplicationConsistencyGroupID are required")
	}

	bytes, err := json.Marshal(rp)
	if err != nil {
		fmt.Printf("Marshal error: %s\n", err)
	}
	fmt.Printf("Marshal output: %s\n", string(bytes))
	defer TimeSpent("CreateReplicationPair", time.Now())

	path := "/api/types/ReplicationPair/instances"
	rpResp := &types.ReplicationPair{}

	err = c.getJSONWithRetry(http.MethodPost, path, rp, rpResp)
	if err != nil {
		fmt.Printf("c.getJSONWithRetry(http.MethodPost, path, rp, rpResp) returned %s", err)
		return nil, err
	}
	return rpResp, nil
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
