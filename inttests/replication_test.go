/*
 *
 * Copyright © 2020 Dell Inc. or its subsidiaries. All Rights Reserved.
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

package inttests

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/dell/goscaleio"
	siotypes "github.com/dell/goscaleio/types/v1"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

// Replication global variables used to set up the replication relationships
type replication struct {
	sourcePeerMDM          *goscaleio.PeerMDM
	targetSystem           *goscaleio.System
	sourceSystemID         string
	sourceProtectionDomain *goscaleio.ProtectionDomain
	sourceStoragePool      *goscaleio.StoragePool
	sourceVolume           *siotypes.Volume
	targetSystemID         string
	targetProtectionDomain *goscaleio.ProtectionDomain
	targetStoragePool      *goscaleio.StoragePool
	targetVolume           *siotypes.Volume
	rcg                    *goscaleio.ReplicationConsistencyGroup
	rcgID                  string
	replicationPair        *siotypes.ReplicationPair
	snapshotGroupID        string
}

var rep replication

const (
	sourceVolume           = "GOSCALEIO_REPLICATION_SOURCE_NAME"
	targetVolume           = "GOSCALEIO_REPLICATION_TARGET_NAME"
	targetProtectionDomain = "GOSCALEIO_PROTECTIONDOMAIN2"
	targetStoragePool      = "GOSCALEIO_STORAGEPOOL2"
)

// Test GetPeerMDMs
func TestGetPeerMDMs(t *testing.T) {
	srcpeers, err := C.GetPeerMDMs()
	assert.Nil(t, err)
	for i := 0; i < len(srcpeers); i++ {
		t.Logf("Source PeerMDM: %+v", srcpeers[i])
		rep.sourceSystemID = srcpeers[i].SystemID
	}

	if C2 == nil {
		t.Skip("no client connection to replication target system")
	}

	log.SetLevel(log.DebugLevel)
	tgtpeers, err := C2.GetPeerMDMs()
	assert.Nil(t, err)
	log.SetLevel(log.InfoLevel)
	for i := 0; i < len(tgtpeers); i++ {
		t.Logf("Target PeerMDM: %+v", tgtpeers[i])
		rep.targetSystemID = tgtpeers[i].SystemID
	}

	// Test systems are validly paired
	found := false
	for i := 0; i < len(srcpeers); i++ {
		if srcpeers[i].PeerSystemID == rep.targetSystemID {
			if srcpeers[i].CouplingRC != "SUCCESS" {
				t.Error(fmt.Printf("PeerMDM %s expected couplingRC SUCCESS but status was %s", srcpeers[i].PeerSystemID, srcpeers[i].CouplingRC))
			}

			found = true
			rep.sourcePeerMDM = goscaleio.NewPeerMDM(C, srcpeers[i])
			t.Logf("PeerMDMID %s", rep.sourcePeerMDM.PeerMDM.ID)
			break
		}
	}

	if !found {
		t.Error("Didn't find target MDM peer")
	}

	found = false
	for i := 0; i < len(tgtpeers); i++ {
		if tgtpeers[i].PeerSystemID == rep.sourceSystemID {
			if tgtpeers[i].CouplingRC != "SUCCESS" {
				t.Error(fmt.Printf("PeerMDM %s expected couplingRC SUCCESS but status was %s", tgtpeers[i].PeerSystemID, tgtpeers[i].CouplingRC))
			}

			found = true
			break
		}
	}
	if !found {
		t.Error("Didn't find source MDM peer")
	}
}

// Get the Target System
func getTargetSystem() *goscaleio.System {
	system := goscaleio.NewSystem(C2)
	targetSystems, _ := C2.GetSystems()
	if len(targetSystems) > 0 {
		system.System = targetSystems[0]
	}
	rep.targetSystem = system
	return rep.targetSystem
}

func TestGetTargetSystem(t *testing.T) {
	if C2 == nil {
		t.Skip("no client connection to replication target system")
	}

	getTargetSystem()
	assert.NotNil(t, rep.targetSystem)
}

// Test getProtectionDomain
func TestGetProtectionDomain(t *testing.T) {
	rep.sourceProtectionDomain = getProtectionDomain(t)
	assert.NotNil(t, rep.sourceProtectionDomain)

	t.Logf("source protction domain: %+v", rep.sourceProtectionDomain.ProtectionDomain)
	href := "/api/instances/ProtectionDomain::" + rep.sourceProtectionDomain.ProtectionDomain.ID
	protectionDomains, err := getSystem().GetProtectionDomain(href)
	assert.Nil(t, err)

	for _, pd := range protectionDomains {
		t.Logf("source protection domain %+v", pd)
	}
}

// Get the Target Protection Domain
func getTargetProtectionDomain() *goscaleio.ProtectionDomain {
	targetProtectionDomainName := os.Getenv(targetProtectionDomain)
	protectionDomains, _ := rep.targetSystem.GetProtectionDomain("")
	for i := 0; i < len(protectionDomains); i++ {
		fmt.Printf("target protection domain %+v", protectionDomains[i])
		if protectionDomains[i].Name == targetProtectionDomainName {
			return goscaleio.NewProtectionDomainEx(C2, protectionDomains[i])
		}
	}
	return nil
}

func TestTargetProtectionDomain(t *testing.T) {
	if C2 == nil {
		t.Skip("no client connection to replication target system")
	}

	rep.targetProtectionDomain = getTargetProtectionDomain()
	assert.NotNil(t, rep.targetProtectionDomain)
}

// Get the Target StoragePool
func getTargetStoragePool() *goscaleio.StoragePool {
	targetStoragePoolName := os.Getenv(targetStoragePool)
	storagePools, _ := rep.targetProtectionDomain.GetStoragePool("")
	for i := 0; i < len(storagePools); i++ {
		if storagePools[i].Name == targetStoragePoolName {
			return goscaleio.NewStoragePoolEx(C2, storagePools[i])
		}
	}
	return nil
}

func TestStoragePools(t *testing.T) {
	if C2 == nil {
		t.Skip("no client connection to replication target system")
	}
	rep.sourceStoragePool = getStoragePool(t)
	assert.NotNil(t, rep.sourceStoragePool)
	t.Logf("sourceStoragePool %s", rep.sourceStoragePool.StoragePool.Name)

	rep.targetStoragePool = getTargetStoragePool()
	assert.NotNil(t, rep.targetStoragePool)
	t.Logf("targetStoragePool %s", rep.targetStoragePool.StoragePool.Name)
}

// Locate the volumes to be replicated together
func TestLocateVolumesToBeReplicated(t *testing.T) {
	if C2 == nil {
		t.Skip("no client connection to replication target system")
	}

	srcName := os.Getenv(sourceVolume)
	assert.NotNil(t, srcName)
	sourceVolumes, err := rep.sourceStoragePool.GetVolume("", "", "", srcName, false)
	if err != nil {
		t.Log(err)
	}
	if len(sourceVolumes) > 0 {
		rep.sourceVolume = sourceVolumes[0]
	}
	assert.NotNil(t, rep.sourceVolume)

	dstName := os.Getenv(targetVolume)
	assert.NotNil(t, dstName)
	targetVolumes, err := rep.targetStoragePool.GetVolume("", "", "", dstName, false)
	if err != nil {
		t.Log(err)
	}
	if len(targetVolumes) > 0 {
		rep.targetVolume = targetVolumes[0]
	}
	assert.NotNil(t, rep.targetVolume)

	t.Logf("SourceVolume %s, TargetVolume %s", rep.sourceVolume.Name, rep.targetVolume.Name)
}

// Test createReplicationConsistencyGroup
func TestCreateReplicationConsistencyGroup(t *testing.T) {
	if C2 == nil {
		t.Skip("no client connection to replication target system")
	}

	rcgPayload := &siotypes.ReplicationConsistencyGroupCreatePayload{
		Name:                     "inttestrcg",
		RpoInSeconds:             "60",
		ProtectionDomainId:       rep.sourceProtectionDomain.ProtectionDomain.ID,
		RemoteProtectionDomainId: rep.targetProtectionDomain.ProtectionDomain.ID,
		//PeerMdmId:                rep.sourcePeerMDM.PeerMDM.ID,
		DestinationSystemId: rep.targetSystem.System.ID,
	}

	rcgResp, err := C.CreateReplicationConsistencyGroup(rcgPayload)
	assert.Nil(t, err)

	log.Infof("RCG ID: %s", rcgResp.ID)
	rep.rcgID = rcgResp.ID

	time.Sleep(5 * time.Second)
}

// Add Replication Pair
func TestAddReplicationPair(t *testing.T) {
	if C2 == nil {
		t.Skip("no client connection to replication target system")
	}

	srcName := os.Getenv(sourceVolume)
	assert.NotNil(t, srcName)
	t.Logf("Source Volume %s", srcName)

	localVolumeID, err := C.FindVolumeID(srcName)
	assert.Nil(t, err)

	t.Logf("[TestAddReplicationPair] Local Volume ID: %s", localVolumeID)

	dstName := os.Getenv(targetVolume)
	assert.NotNil(t, dstName)
	t.Logf("Target Volume %s", dstName)

	remoteVolumeID, err := C2.FindVolumeID(dstName)
	assert.Nil(t, err)

	t.Logf("[TestAddReplicationPair] Remote Volume ID: %s", remoteVolumeID)

	vol, err := C.GetVolume("", strings.TrimSpace(localVolumeID), "", "", false)
	assert.Nil(t, err)

	t.Logf("Source Volume Content %+v", vol[0])

	rpPayload := &siotypes.QueryReplicationPair{
		Name:                          "inttestrp",
		SourceVolumeID:                localVolumeID,
		DestinationVolumeID:           remoteVolumeID,
		ReplicationConsistencyGroupID: rep.rcgID,
		CopyType:                      "OnlineCopy",
	}

	rpResp, err := C.CreateReplicationPair(rpPayload)
	assert.Nil(t, err)

	t.Logf("ReplicationPairID: %s", rpResp.ID)
	rep.replicationPair = rpResp
}

// Query Replication Pair
func TestQueryReplicationPairs(t *testing.T) {
	if C2 == nil {
		t.Skip("no client connection to replication target system")
	}

	pairs, err := C.GetReplicationPairs(rep.rcgID)
	assert.Nil(t, err)

	for i, pair := range pairs {
		t.Logf("%d, ReplicationPair: %+v", i, pair)
		rep.replicationPair = pair
	}
}

// Query Replication Pair Statistics
func TestQueryReplicationPairsStatistics(t *testing.T) {
	if C2 == nil {
		t.Skip("no client connection to replication target system")
	}

	assert.NotNil(t, rep.replicationPair)

	t.Logf("Waiting for Replication Pair %s to be complete.", rep.replicationPair.Name)
	for i := 0; i < 30; i++ {
		rpResp, err := C.GetReplicationPairStatistics(rep.replicationPair.ID)
		assert.Nil(t, err)

		t.Logf("Copied %f", rpResp.InitialCopyProgress)

		group, err := C.GetReplicationConsistencyGroupById(rep.rcgID)
		assert.Nil(t, err)

		// Check if complete
		if rpResp.InitialCopyProgress == 1 && group.CurrConsistMode == "Consistent" {
			t.Logf("Copy Complete: %f", rpResp.InitialCopyProgress)
			break
		}

		time.Sleep(10 * time.Second)
	}
}

// Test GetReplicationConsistencyGroups
func TestGetReplicationConsistencyGroups(t *testing.T) {
	if C2 == nil {
		t.Skip("no client connection to replication target system")
	}

	rcgs, err := C.GetReplicationConsistencyGroups()
	assert.Nil(t, err)
	for i := 0; i < len(rcgs); i++ {
		t.Logf("RCG: %+v\n\n", rcgs[i])
		assert.Nil(t, err)

		parseLinks(rcgs[i].Links, t)

		if rcgs[i].Name == "inttestrcg" {
			rcg := goscaleio.NewReplicationConsistencyGroup(C)
			rcg.ReplicationConsistencyGroup = rcgs[i]
			rep.rcg = rcg
		}
	}
	assert.NotNil(t, rep.rcg)
}

// Test CreateReplicationConsistencyGroupSnapshot
func TestCreateReplicationConsistencyGroupSnapshot(t *testing.T) {
	if C2 == nil {
		t.Skip("no client connection to replication target system")
	}

	resp, err := C.CreateReplicationConsistencyGroupSnapshot(rep.rcgID, false)
	assert.Nil(t, err)
	rep.snapshotGroupID = resp.SnapshotGroupID
}

// Test SnapshotRetrieval
func TestSnapshotRetrieval(t *testing.T) {
	if C2 == nil {
		t.Skip("no client connection to replication target system")
	}

	pairs, err := C.GetReplicationPairs(rep.rcgID)
	assert.Nil(t, err)

	var vols []string
	for _, pair := range pairs {
		fmt.Printf("Remote Pair Volume: %s\n", pair.RemoteVolumeID)
		vols = append(vols, pair.RemoteVolumeID)
	}

	actionAttributes := make(map[string]string)
	for _, vol := range vols {
		result, err := C2.GetVolume("", "", vol, "", false)
		if err != nil {
			fmt.Printf("Get Vols Error: %s\n", err.Error())
		} else {
			for _, snap := range result {
				fmt.Printf("Get Vols Content: %+v\n", snap)
				if rep.snapshotGroupID == snap.ConsistencyGroupID {
					actionAttributes[snap.AncestorVolumeID] = snap.ID
				}
			}
		}
	}

	fmt.Printf("Action Attributes Result: %+v\n", actionAttributes)
}

// Test ExecuteFailoverOnReplicationGroup
func TestExecuteFailoverOnReplicationGroup(t *testing.T) {
	if C2 == nil {
		t.Skip("no client connection to replication target system")
	}

	err := C.ExecuteFailoverOnReplicationGroup(rep.rcgID)
	assert.Nil(t, err)

}

// Test ExecuteRestoreOnReplicationGroup
func TestExecuteRestoreOnReplicationGroup(t *testing.T) {
	if C2 == nil {
		t.Skip("no client connection to replication target system")
	}

	err := ensureFailover(t)
	assert.Nil(t, err)

	err = C.ExecuteRestoreOnReplicationGroup(rep.rcgID)
	assert.Nil(t, err)
}

// Test ExecuteSwitchoverOnReplicationGroup
func TestExecuteSwitchoverOnReplicationGroup(t *testing.T) {
	if C2 == nil {
		t.Skip("no client connection to replication target system")
	}

	err := waitForConsistency(t)
	assert.Nil(t, err)

	err = C.ExecuteSwitchoverOnReplicationGroup(rep.rcgID, false)
	assert.Nil(t, err)
}

// Test ExecuteReverseOnReplicationGroup
func TestExecuteReverseOnReplicationGroup(t *testing.T) {
	if C2 == nil {
		t.Skip("no client connection to replication target system")
	}

	err := ensureFailover(t)
	assert.Nil(t, err)

	err = C.ExecuteReverseOnReplicationGroup(rep.rcgID)
	assert.Nil(t, err)
}

// Test ExecutePauseOnReplicationGroup
func TestExecutePauseOnReplicationGroup(t *testing.T) {
	if C2 == nil {
		t.Skip("no client connection to replication target system")
	}

	err := waitForConsistency(t)
	assert.Nil(t, err)

	err = C.ExecutePauseOnReplicationGroup(rep.rcgID, siotypes.ONLY_TRACK_CHANGES)
	assert.Nil(t, err)
}

// Test ExecuteResumeOnReplicationGroup
func TestExecuteResumeOnReplicationGroup(t *testing.T) {
	if C2 == nil {
		t.Skip("no client connection to replication target system")
	}

	err := C.ExecuteResumeOnReplicationGroup(rep.rcgID)
	assert.Nil(t, err)
}

// Test RemoveReplicationPair
func TestRemoveReplicationPairFromVolume(t *testing.T) {
	if C2 == nil {
		t.Skip("no client connection to replication target system")
	}

	pairs, err := C.GetReplicationPairs("")
	assert.Nil(t, err)

	var replicationPairId string
	for i, pair := range pairs {
		t.Logf("%d, ReplicationPair: %+v", i, pair)

		if rep.replicationPair.LocalVolumeID == pair.LocalVolumeID {
			replicationPairId = pair.ID
			break
		}
	}

	if replicationPairId == "" {
		t.Logf("replication pair for that volume not found")
		assert.NotNil(t, replicationPairId)
	}

	_, err = C.RemoveReplicationPair(replicationPairId, true)
	assert.Nil(t, err)

	t.Logf("[TestRemoveReplicationPairFromVolume] Removed the following pair %s", rep.replicationPair.Name)

	// Delay to verify on the UI.
	time.Sleep(5 * time.Second)
}

// Test Freeze Replication Group
func TestFreezeReplcationGroup(t *testing.T) {
	if C2 == nil {
		t.Skip("no client connection to replication target system")
	}

	t.Logf("[TestFreezeReplcationGroup] Freezing replication group: %s", rep.rcgID)

	err := rep.rcg.FreezeReplicationConsistencyGroup(rep.rcgID)
	assert.Nil(t, err)

	t.Logf("[TestFreezeReplcationGroup] Froze replication pair, check UI")

	// Delay to verify on the UI.
	time.Sleep(2 * time.Second)
}

// Test RemoveReplicatonConsistencyGroup
func TestRemoveReplicationConsistencyGroup(t *testing.T) {
	if C2 == nil {
		t.Skip("no client connection to replication target system")
	}

	assert.NotNil(t, rep.rcg)

	err := rep.rcg.RemoveReplicationConsistencyGroup(false)
	assert.Nil(t, err)
}

func parseLinks(links []*siotypes.Link, t *testing.T) {
	length := len(links)

	if length == 0 {
		t.Logf("No links found in the RCG")
		return
	}

	for _, link := range links {
		t.Logf("Rel: %s\nHREF: %s\n", link.Rel, link.HREF)
	}
}

func waitForConsistency(t *testing.T) error {
	for i := 0; i < 15; i++ {
		group, err := C.GetReplicationConsistencyGroupById(rep.rcgID)
		if err != nil {
			continue
		}

		if group.CurrConsistMode == "Consistent" && group.FailoverType == "None" {
			t.Logf("Consistency Group %s - Reached Consistency.", rep.rcgID)
			return nil
		}

		time.Sleep(5 * time.Second)
	}
	return errors.New("consistency group did not reach consistency.")
}

func ensureFailover(t *testing.T) error {
	for i := 0; i < 30; i++ {
		group, err := C.GetReplicationConsistencyGroupById(rep.rcgID)
		if err != nil {
			return errors.New("No replication consistency groups found: %")
		}

		t.Logf("[ensureFailover] - %+v", group)

		if group.FailoverType != "None" && group.FailoverState == "Done" && group.DisasterRecoveryState == "Neutral" && group.RemoteDisasterRecoveryState == "Neutral" {
			t.Logf("Consistency Group is in %s", group.FailoverType)
			return nil
		}

		time.Sleep(1 * time.Second)
	}

	return errors.New("unable to reach failover consistency")
}
