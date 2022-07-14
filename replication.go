package goscaleio

import(
	"net/http"
	"time"
	types "github.com/dell/goscaleio/types/v1"
)

// GetPeerMDMs returns a list of peer MDMs know to the System
func (c *Client) GetPeerMDMs()([]*types.PeerMDM, error) {
	path := "/api/types/PeerMdm/instances"
	var peerMdms []*types.PeerMDM
	var err error
	defer TimeSpent("GetPeerMDMs", time.Now())
	
	err = c.getJSONWithRetry(http.MethodGet, path, nil, &peerMdms)
	return peerMdms, err
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
