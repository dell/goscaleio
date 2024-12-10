/*
 *
 * Copyright © 2020-2024 Dell Inc. or its subsidiaries. All Rights Reserved.
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
	"context"
	"errors"
	"fmt"
	"os"
	"strconv"
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
	targetSystem           *goscaleio.System
	sourceProtectionDomain *goscaleio.ProtectionDomain
	sourceStoragePool      *goscaleio.StoragePool
	targetProtectionDomain *goscaleio.ProtectionDomain
	targetStoragePool      *goscaleio.StoragePool
	targetVolume           *siotypes.Volume
	rcg                    *goscaleio.ReplicationConsistencyGroup
	rcgID                  string
	pair                   *goscaleio.ReplicationPair
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
	ctx := context.Background()
	srcpeers, err := C.GetPeerMDMs(ctx)
	assert.Nil(t, err)

	var sourceSystemID string
	for i := 0; i < len(srcpeers); i++ {
		sourceSystemID = srcpeers[i].SystemID
	}

	if C2 == nil {
		t.Skip("no client connection to replication target system")
	}

	tgtpeers, err := C2.GetPeerMDMs(ctx)
	assert.Nil(t, err)

	var targetSystemID string
	for i := 0; i < len(tgtpeers); i++ {
		fmt.Printf("Peer %d, %+v", i, tgtpeers[i])
		targetSystemID = tgtpeers[i].SystemID
	}

	// Test systems are validly paired
	found := false
	for i := 0; i < len(srcpeers); i++ {
		if srcpeers[i].PeerSystemID == targetSystemID {
			if srcpeers[i].CouplingRC != "SUCCESS" {
				t.Error(fmt.Printf("PeerMDM %s expected couplingRC SUCCESS but status was %s", srcpeers[i].PeerSystemID, srcpeers[i].CouplingRC))
			}

			found = true
			peer := goscaleio.NewPeerMDM(C, srcpeers[i])
			t.Logf("PeerMDMID %s", peer.PeerMDM.ID)
			break
		}
	}

	if !found {
		t.Error("Didn't find target MDM peer")
	}

	found = false
	for i := 0; i < len(tgtpeers); i++ {
		if tgtpeers[i].PeerSystemID == sourceSystemID {
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

// Get Specific Peer System
func TestGetPeerSystem(t *testing.T) {
	srcpeer, err := getPeerMdm()
	assert.Nil(t, err)

	_, getErr := C.GetPeerMDM(context.Background(), srcpeer.ID)
	assert.Nil(t, getErr)
}

// Remove a Peer System
func TestRemovePeerSystem(t *testing.T) {
	srcpeer, err := getPeerMdm()
	assert.Nil(t, err)

	removeErr := C.RemovePeerMdm(context.Background(), srcpeer.ID)
	assert.Nil(t, removeErr)
}

func TestAddPeerSystem(t *testing.T) {
	if C2 == nil {
		t.Skip("no client connection to replication target system")
	}
	system := getTargetSystem()
	peerPayload := &siotypes.AddPeerMdm{
		PeerSystemID:  system.System.ID,
		PeerSystemIps: system.System.MdmManagementIPList,
		Port:          "7611",
		Name:          "PeerSystemTestName",
	}
	// Add a Peer System
	_, err := C.AddPeerMdm(context.Background(), peerPayload)
	assert.Nil(t, err)
}

// Modify a Peer System name
func TestModifyPeerSystemName(t *testing.T) {
	srcpeer, err := getPeerMdm()
	assert.Nil(t, err)
	ctx := context.Background()

	// Modify name
	modifyErr := C.ModifyPeerMdmName(ctx, srcpeer.ID, &siotypes.ModifyPeerMDMNameParam{
		NewName: "testName",
	})
	assert.Nil(t, modifyErr)

	// Modify Name back to original
	modifyErr = C.ModifyPeerMdmName(ctx, srcpeer.ID, &siotypes.ModifyPeerMDMNameParam{
		NewName: srcpeer.Name,
	})
}

// Modify a Peer System name
func TestModifyPeerSystemPort(t *testing.T) {
	srcpeer, err := getPeerMdm()
	assert.Nil(t, err)
	ctx := context.Background()

	// Modify port
	modifyErr := C.ModifyPeerMdmPort(ctx, srcpeer.ID, &siotypes.ModifyPeerMDMPortParam{
		NewPort: "7612",
	})
	assert.Nil(t, modifyErr)

	// Modify Name back to original
	modifyErr = C.ModifyPeerMdmPort(ctx, srcpeer.ID, &siotypes.ModifyPeerMDMPortParam{
		NewPort: fmt.Sprint(srcpeer.Port),
	})
}

// Modify a Peer System Performance Parameters
func TestModifyPeerSystemPerformanceParameters(t *testing.T) {
	srcpeer, err := getPeerMdm()
	assert.Nil(t, err)
	ctx := context.Background()

	// Modify port
	modifyErr := C.ModifyPeerMdmPerformanceParameters(ctx, srcpeer.ID, &siotypes.ModifyPeerMdmPerformanceParametersParam{
		NewPreformanceProfile: "Compact",
	})
	assert.Nil(t, modifyErr)

	// Modify Name back to original
	modifyErr = C.ModifyPeerMdmPerformanceParameters(ctx, srcpeer.ID, &siotypes.ModifyPeerMdmPerformanceParametersParam{
		NewPreformanceProfile: srcpeer.PerfProfile,
	})
}

// Modify a Peer System Ips
func TestModifyPeerSystemIps(t *testing.T) {
	srcpeer, err := getPeerMdm()
	assert.Nil(t, err)

	var ips []string
	for _, ip := range srcpeer.IPList {
		ips = append(ips, ip.IP)
	}

	// Modify ips
	modifyErr := C.ModifyPeerMdmIP(context.Background(), srcpeer.ID, ips)
	assert.Nil(t, modifyErr)
}

// Make it easier to get a Peer System to run the tests against
func getPeerMdm() (*siotypes.PeerMDM, error) {
	srcpeers, err := C.GetPeerMDMs(context.Background())

	if err != nil || len(srcpeers) == 0 {
		return nil, fmt.Errorf("no peer systems found")
	}
	return srcpeers[0], nil
}

// Get the Target System
func getTargetSystem() *goscaleio.System {
	system := goscaleio.NewSystem(C2)
	targetSystems, _ := C2.GetSystems(context.Background())

	if len(targetSystems) > 0 {
		system.System = targetSystems[0]
	}

	return system
}

func TestGetTargetSystem(t *testing.T) {
	if C2 == nil {
		t.Skip("no client connection to replication target system")
	}

	rep.targetSystem = getTargetSystem()
	fmt.Printf("Target: %+v", rep.targetSystem.System)
	assert.NotNil(t, rep.targetSystem)
}

// Test getProtectionDomain
func TestGetProtectionDomain(t *testing.T) {
	rep.sourceProtectionDomain = getProtectionDomain(t)
	assert.NotNil(t, rep.sourceProtectionDomain)

	href := "/api/instances/ProtectionDomain::" + rep.sourceProtectionDomain.ProtectionDomain.ID
	_, err := getSystem().GetProtectionDomain(context.Background(), href)
	assert.Nil(t, err)
}

// Get the Target Protection Domain
func getTargetProtectionDomain() *goscaleio.ProtectionDomain {
	targetProtectionDomainName := os.Getenv(targetProtectionDomain)
	protectionDomains, _ := rep.targetSystem.GetProtectionDomain(context.Background(), "")
	for i := 0; i < len(protectionDomains); i++ {
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
	storagePools, _ := rep.targetProtectionDomain.GetStoragePool(context.Background(), "")
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

	rep.targetStoragePool = getTargetStoragePool()
	assert.NotNil(t, rep.targetStoragePool)
}

// Locate the volumes to be replicated together
func TestLocateVolumesToBeReplicated(t *testing.T) {
	if C2 == nil {
		t.Skip("no client connection to replication target system")
	}

	ctx := context.Background()

	srcName := os.Getenv(sourceVolume)
	assert.NotNil(t, srcName)
	sourceVolumes, err := rep.sourceStoragePool.GetVolume(ctx, "", "", "", srcName, false)
	if err != nil {
		t.Log(err)
	}

	var sourceVolume *siotypes.Volume
	if len(sourceVolumes) > 0 {
		sourceVolume = sourceVolumes[0]
	}
	assert.NotNil(t, sourceVolume)

	dstName := os.Getenv(targetVolume)
	assert.NotNil(t, dstName)
	targetVolumes, err := rep.targetStoragePool.GetVolume(ctx, "", "", "", dstName, false)
	if err != nil {
		t.Log(err)
	}
	if len(targetVolumes) > 0 {
		rep.targetVolume = targetVolumes[0]
	}
	assert.NotNil(t, rep.targetVolume)

	t.Logf("SourceVolume %s, TargetVolume %s", sourceVolume.Name, rep.targetVolume.Name)
}

// Test createReplicationConsistencyGroup
func TestCreateReplicationConsistencyGroup(t *testing.T) {
	if C2 == nil {
		t.Skip("no client connection to replication target system")
	}

	rcgPayload := &siotypes.ReplicationConsistencyGroupCreatePayload{
		Name:                     "inttestrcg",
		RpoInSeconds:             "60",
		ProtectionDomainID:       rep.sourceProtectionDomain.ProtectionDomain.ID,
		RemoteProtectionDomainID: rep.targetProtectionDomain.ProtectionDomain.ID,
		DestinationSystemID:      rep.targetSystem.System.ID,
	}

	rcgResp, err := C.CreateReplicationConsistencyGroup(context.Background(), rcgPayload)
	assert.Nil(t, err)

	log.Debugf("RCG ID: %s", rcgResp.ID)
	rep.rcgID = rcgResp.ID

	time.Sleep(5 * time.Second)
}

// Test GetReplicationConsistencyGroups
func TestGetReplicationConsistencyGroups(t *testing.T) {
	if C2 == nil {
		t.Skip("no client connection to replication target system")
	}

	rcgs, err := C.GetReplicationConsistencyGroups(context.Background())
	assert.Nil(t, err)
	for i := 0; i < len(rcgs); i++ {
		assert.Nil(t, err)

		if rcgs[i].Name == "inttestrcg" {
			rcg := goscaleio.NewReplicationConsistencyGroup(C)
			rcg.ReplicationConsistencyGroup = rcgs[i]
			rep.rcg = rcg
		}
	}
	assert.NotNil(t, rep.rcg)
}

// Add Replication Pair
func TestAddReplicationPair(t *testing.T) {
	if C2 == nil {
		t.Skip("no client connection to replication target system")
	}
	ctx := context.Background()

	srcName := os.Getenv(sourceVolume)
	assert.NotNil(t, srcName)

	localVolumeID, err := C.FindVolumeID(ctx, srcName)
	assert.Nil(t, err)

	t.Logf("[TestAddReplicationPair] Local Volume ID: %s", localVolumeID)

	dstName := os.Getenv(targetVolume)
	assert.NotNil(t, dstName)

	remoteVolumeID, err := C2.FindVolumeID(ctx, dstName)
	assert.Nil(t, err)

	t.Logf("[TestAddReplicationPair] Remote Volume ID: %s", remoteVolumeID)

	_, err = C.GetVolume(ctx, "", strings.TrimSpace(localVolumeID), "", "", false)
	assert.Nil(t, err)

	rpPayload := &siotypes.QueryReplicationPair{
		Name:                          "inttestrp",
		SourceVolumeID:                localVolumeID,
		DestinationVolumeID:           remoteVolumeID,
		ReplicationConsistencyGroupID: rep.rcgID,
		CopyType:                      "OnlineCopy",
	}

	rpResp, err := C.CreateReplicationPair(ctx, rpPayload)
	assert.Nil(t, err)

	t.Logf("ReplicationPairID: %s", rpResp.ID)
	replicationPair := goscaleio.NewReplicationPair(C)
	replicationPair.ReplicaitonPair = rpResp
	rep.pair = replicationPair

	// Ensure array knows pair exists.
	time.Sleep(5 * time.Second)
}

// Query Replication Pair
func TestQueryReplicationPairs(t *testing.T) {
	if C2 == nil {
		t.Skip("no client connection to replication target system")
	}

	pairs, err := rep.rcg.GetReplicationPairs(context.Background())
	assert.Nil(t, err)

	for _, pair := range pairs {
		replicationPair := goscaleio.NewReplicationPair(C)
		replicationPair.ReplicaitonPair = pair
		rep.pair = replicationPair
	}
}

// Query Specific Replication Pair
func TestQueryReplicationPair(t *testing.T) {
	if C2 == nil {
		t.Skip("no client connection to replication target system")
	}

	pair, err := C.GetReplicationPair(context.Background(), rep.pair.ReplicaitonPair.ID)
	assert.Nil(t, err)
	assert.NotNil(t, pair)
}

// Pause and Resume Replication Pair
func TestPauseAndResumeReplicationPair(t *testing.T) {
	if C2 == nil {
		t.Skip("no client connection to replication target system")
	}

	ctx := context.Background()

	// Pause
	pairP, err := C.PausePairInitialCopy(ctx, rep.pair.ReplicaitonPair.ID)
	assert.Nil(t, err)
	assert.NotNil(t, pairP)

	// Resume
	pairR, err := C.ResumePairInitialCopy(ctx, rep.pair.ReplicaitonPair.ID)
	assert.Nil(t, err)
	assert.NotNil(t, pairR)
}

// Query Replication Pair Statistics
func TestQueryReplicationPairsStatistics(t *testing.T) {
	if C2 == nil {
		t.Skip("no client connection to replication target system")
	}

	assert.NotNil(t, rep.pair)
	ctx := context.Background()

	t.Logf("Waiting for Replication Pair %s to be complete.", rep.pair.ReplicaitonPair.Name)
	for i := 0; i < 30; i++ {
		rpResp, err := rep.pair.GetReplicationPairStatistics(ctx)
		assert.Nil(t, err)

		t.Logf("Copied %f", rpResp.InitialCopyProgress)

		group, err := C.GetReplicationConsistencyGroupByID(ctx, rep.rcgID)
		assert.Nil(t, err)

		// Check if complete
		if rpResp.InitialCopyProgress == 1 && group.CurrConsistMode == "Consistent" {
			t.Logf("Copy Complete: %f", rpResp.InitialCopyProgress)
			break
		}

		time.Sleep(10 * time.Second)
	}
}

// Test CreateReplicationConsistencyGroupSnapshot
func TestCreateReplicationConsistencyGroupSnapshot(t *testing.T) {
	if C2 == nil {
		t.Skip("no client connection to replication target system")
	}

	resp, err := rep.rcg.CreateReplicationConsistencyGroupSnapshot(context.Background())
	assert.Nil(t, err)

	t.Logf("Consistency Group Snapshot ID: %s", resp.SnapshotGroupID)
	rep.snapshotGroupID = resp.SnapshotGroupID
}

// Test SnapshotRetrieval
func TestSnapshotRetrieval(t *testing.T) {
	if C2 == nil {
		t.Skip("no client connection to replication target system")
	}
	ctx := context.Background()

	pairs, err := rep.rcg.GetReplicationPairs(ctx)
	assert.Nil(t, err)

	var vols []string
	for _, pair := range pairs {
		t.Logf("Remote Pair Volume: %s\n", pair.RemoteVolumeID)
		vols = append(vols, pair.RemoteVolumeID)
	}

	actionAttributes := make(map[string]string)
	for _, vol := range vols {
		result, err := C2.GetVolume(ctx, "", "", vol, "", false)
		if err != nil {
			t.Errorf("Get Vols Error: %s\n", err.Error())
		} else {
			for _, snap := range result {
				if rep.snapshotGroupID == snap.ConsistencyGroupID {
					actionAttributes[snap.AncestorVolumeID] = snap.ID
				}
			}
		}
	}

	t.Logf("Action Attributes Result: %+v\n", actionAttributes)
}

// Test ExecuteFailoverOnReplicationGroup
func TestExecuteFailoverOnReplicationGroup(t *testing.T) {
	if C2 == nil {
		t.Skip("no client connection to replication target system")
	}

	err := rep.rcg.ExecuteFailoverOnReplicationGroup(context.Background())
	assert.Nil(t, err)
}

// Test ExecuteRestoreOnReplicationGroup
func TestExecuteRestoreOnReplicationGroup(t *testing.T) {
	if C2 == nil {
		t.Skip("no client connection to replication target system")
	}

	err := ensureFailover(t)
	assert.Nil(t, err)

	err = rep.rcg.ExecuteRestoreOnReplicationGroup(context.Background())
	assert.Nil(t, err)
}

// Test ExecuteSwitchoverOnReplicationGroup
func TestExecuteSwitchoverOnReplicationGroup(t *testing.T) {
	if C2 == nil {
		t.Skip("no client connection to replication target system")
	}

	err := waitForConsistency(t)
	assert.Nil(t, err)

	err = rep.rcg.ExecuteSwitchoverOnReplicationGroup(context.Background(), false)
	assert.Nil(t, err)
}

// Test ExecuteReverseOnReplicationGroup
func TestExecuteReverseOnReplicationGroup(t *testing.T) {
	if C2 == nil {
		t.Skip("no client connection to replication target system")
	}

	err := ensureFailover(t)
	assert.Nil(t, err)

	err = rep.rcg.ExecuteReverseOnReplicationGroup(context.Background())
	assert.Nil(t, err)
}

// Test ExecutePauseOnReplicationGroup
func TestExecutePauseOnReplicationGroup(t *testing.T) {
	if C2 == nil {
		t.Skip("no client connection to replication target system")
	}

	err := waitForConsistency(t)
	assert.Nil(t, err)

	err = rep.rcg.ExecutePauseOnReplicationGroup(context.Background())
	assert.Nil(t, err)
}

// Test ExecuteResumeOnReplicationGroup
func TestExecuteResumeOnReplicationGroup(t *testing.T) {
	if C2 == nil {
		t.Skip("no client connection to replication target system")
	}

	err := rep.rcg.ExecuteResumeOnReplicationGroup(context.Background())
	assert.Nil(t, err)
}

// Test TestSetRPOOnReplicationGroup
func TestSetRPOOnReplicationGroup(t *testing.T) {
	// Set the RCG context
	TestGetReplicationConsistencyGroups(t)
	// Update the RPO
	err := rep.rcg.SetRPOOnReplicationGroup(context.Background(), siotypes.SetRPOReplicationConsistencyGroup{RpoInSeconds: "60"})
	assert.Nil(t, err)
}

// Test TestSetTargetVolumeAccessModeOnReplicationGroup
func TestSetTargetVolumeAccessModeOnReplicationGroup(t *testing.T) {
	// Set the RCG context
	TestGetReplicationConsistencyGroups(t)
	err := rep.rcg.SetTargetVolumeAccessModeOnReplicationGroup(context.Background(), siotypes.SetTargetVolumeAccessModeOnReplicationGroup{TargetVolumeAccessMode: "ReadOnly"})
	assert.Nil(t, err)
}

// Test TestSetNewNameOnReplicationGroup
func TestSetNewNameOnReplicationGroup(t *testing.T) {
	// Set the RCG context
	TestGetReplicationConsistencyGroups(t)
	ctx := context.Background()

	err := rep.rcg.SetNewNameOnReplicationGroup(ctx, siotypes.SetNewNameOnReplicationGroup{NewName: "UpdatedNameRCG"})
	assert.Nil(t, err)
	// Sleep for 10 to make sure the name is updated, then update it back to the original name
	time.Sleep(10 * time.Second)
	err = rep.rcg.SetNewNameOnReplicationGroup(ctx, siotypes.SetNewNameOnReplicationGroup{NewName: "inttestrcg"})
	assert.Nil(t, err)
}

// Test TestExecuteInconsistentOnReplicationGroup
func TestExecuteInconsistentOnReplicationGroup(t *testing.T) {
	// Set the RCG context
	TestGetReplicationConsistencyGroups(t)
	err := rep.rcg.ExecuteInconsistentOnReplicationGroup(context.Background())
	assert.Nil(t, err)
}

// Test TestExecuteConsistentOnReplicationGroup
func TestExecuteConsistentOnReplicationGroup(t *testing.T) {
	// Set the RCG context
	TestGetReplicationConsistencyGroups(t)
	err := rep.rcg.ExecuteConsistentOnReplicationGroup(context.Background())
	assert.Nil(t, err)
}

// Test TestExecuteTerminateOnReplicationGroup
func TestExecuteTerminateOnReplicationGroup(t *testing.T) {
	// Set the RCG context
	TestGetReplicationConsistencyGroups(t)
	err := rep.rcg.ExecuteTerminateOnReplicationGroup(context.Background())
	assert.Nil(t, err)
}

// Test TestExecuteActivateOnReplicationGroup
func TestExecuteActivateOnReplicationGroup(t *testing.T) {
	// Set the RCG context
	TestGetReplicationConsistencyGroups(t)
	err := rep.rcg.ExecuteActivateOnReplicationGroup(context.Background())
	assert.Nil(t, err)
}

// Test ResizeReplicationPair
func TestResizeReplicationPair(t *testing.T) {
	if C2 == nil {
		t.Skip("no client connection to replication target system")
	}

	t.Logf("[TestResizeReplicationPair]  Failing over to get DR in Neutral state..")
	TestExecuteFailoverOnReplicationGroup(t)
	ctx := context.Background()

	err := ensureFailover(t)
	assert.Nil(t, err)

	srcName := os.Getenv(sourceVolume)
	assert.NotNil(t, srcName)

	localVolumeID, err := C.FindVolumeID(ctx, srcName)
	assert.Nil(t, err)

	dstName := os.Getenv(targetVolume)
	assert.NotNil(t, dstName)

	remoteVolumeID, err := C2.FindVolumeID(ctx, dstName)
	assert.Nil(t, err)

	sourceVol, err := C.GetVolume(ctx, "", strings.TrimSpace(localVolumeID), "", "", false)
	assert.Nil(t, err)
	assert.NotNil(t, sourceVol)

	destVol, err := C2.GetVolume(ctx, "", strings.TrimSpace(remoteVolumeID), "", "", false)
	assert.Nil(t, err)
	assert.NotNil(t, destVol)

	// Resize destination volume first...
	volume := goscaleio.NewVolume(C2)
	volume.Volume = destVol[0]
	existingSizeGB := volume.Volume.SizeInKb / (1024 * 1024)
	newSize := existingSizeGB * 2
	err = volume.SetVolumeSize(ctx, strconv.Itoa(int(newSize)))
	assert.Nil(t, err)

	// Delay to ensure that the destination syncs up...
	time.Sleep(10 * time.Second)

	volume = goscaleio.NewVolume(C)
	volume.Volume = sourceVol[0]
	existingSizeGB = volume.Volume.SizeInKb / (1024 * 1024)
	newSize = existingSizeGB * 2
	// double the szie of the volume
	err = volume.SetVolumeSize(ctx, strconv.Itoa(int(newSize)))
	assert.Nil(t, err)

	// Restart the initial copy process?
	TestExecuteRestoreOnReplicationGroup(t)
	err = waitForConsistency(t)
	assert.Nil(t, err)
}

// Test RemoveReplicationPair
func TestRemoveReplicationPairFromVolume(t *testing.T) {
	if C2 == nil {
		t.Skip("no client connection to replication target system")
	}
	ctx := context.Background()

	pairs, err := C.GetAllReplicationPairs(ctx)
	assert.Nil(t, err)

	var replicationPairID string
	for _, pair := range pairs {
		if rep.pair.ReplicaitonPair.LocalVolumeID == pair.LocalVolumeID {
			replicationPairID = pair.ID
			break
		}
	}

	if replicationPairID == "" {
		t.Logf("replication pair for that volume not found")
		assert.NotNil(t, replicationPairID)
	}

	_, err = rep.pair.RemoveReplicationPair(ctx, true)
	assert.Nil(t, err)

	t.Logf("[TestRemoveReplicationPairFromVolume] Removed the following pair %s", rep.pair.ReplicaitonPair.Name)

	// Delay to verify on the UI.
	time.Sleep(5 * time.Second)
}

// Test Freeze Replication Group
func TestFreezeReplcationGroup(t *testing.T) {
	if C2 == nil {
		t.Skip("no client connection to replication target system")
	}

	err := rep.rcg.FreezeReplicationConsistencyGroup(context.Background(), rep.rcgID)
	assert.Nil(t, err)

	// Delay to verify on the UI.
	time.Sleep(2 * time.Second)
}

// Test TestUnfreezeReplcationGroup
func TestUnfreezeReplcationGroup(t *testing.T) {
	if C2 == nil {
		t.Skip("no client connection to replication target system")
	}
	TestGetReplicationConsistencyGroups(t)
	assert.NotNil(t, rep.rcg)

	err := rep.rcg.UnfreezeReplicationConsistencyGroup(context.Background())
	assert.Nil(t, err)
}

// Test RemoveReplicatonConsistencyGroup
func TestRemoveReplicationConsistencyGroup(t *testing.T) {
	if C2 == nil {
		t.Skip("no client connection to replication target system")
	}
	TestGetReplicationConsistencyGroups(t)
	assert.NotNil(t, rep.rcg)

	err := rep.rcg.RemoveReplicationConsistencyGroup(context.Background(), false)
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
		group, err := C.GetReplicationConsistencyGroupByID(context.Background(), rep.rcgID)
		if err != nil {
			continue
		}

		if group.CurrConsistMode == "Consistent" && group.FailoverType == "None" {
			t.Logf("Consistency Group %s - Reached Consistency.", rep.rcgID)
			return nil
		}

		time.Sleep(5 * time.Second)
	}
	return errors.New("consistency group did not reach consistency")
}

func ensureFailover(t *testing.T) error {
	for i := 0; i < 30; i++ {
		group, err := C.GetReplicationConsistencyGroupByID(context.Background(), rep.rcgID)
		if err != nil {
			return errors.New("No replication consistency groups found: %")
		}

		if group.FailoverType != "None" && group.FailoverState == "Done" && group.DisasterRecoveryState == "Neutral" && group.RemoteDisasterRecoveryState == "Neutral" {
			t.Logf("Consistency Group is in %s", group.FailoverType)
			time.Sleep(1 * time.Second)
			return nil
		}

		time.Sleep(1 * time.Second)
	}

	return errors.New("unable to reach failover consistency")
}
