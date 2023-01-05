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
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/dell/goscaleio"
	siotypes "github.com/dell/goscaleio/types/v1"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

// Replication global variables used to set up the replication relationships
type replication struct {
	sourcePeerMDM            *goscaleio.PeerMDM
	targetPeerMDM            *goscaleio.PeerMDM
	sourceSystem             *goscaleio.System
	targetSystem             *goscaleio.System
	sourceSystemID           string
	sourceProtectionDomainID string
	sourceProtectionDomain   *goscaleio.ProtectionDomain
	sourceStoragePool        *goscaleio.StoragePool
	sourceVolume             *siotypes.Volume
	targetSystemID           string
	targetProtectionDomain   *goscaleio.ProtectionDomain
	targetStoragePool        *goscaleio.StoragePool
	targetVolume             *siotypes.Volume
	rcg                      *goscaleio.ReplicationConsistencyGroup
	rcgID                    string
	rpID                     string
}

var rep replication

const delay_s = 15

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
	tgtpeers, err := C2.GetPeerMDMs()
	assert.Nil(t, err)
	for i := 0; i < len(tgtpeers); i++ {
		t.Logf("Target PeerMDM: %+v", tgtpeers[i])
		rep.targetSystemID = tgtpeers[i].SystemID
	}

	// Test systems are validly paired
	foundTarget := false
	for i := 0; i < len(srcpeers); i++ {
		if srcpeers[i].PeerSystemID == rep.targetSystemID {
			foundTarget = true
			if srcpeers[i].CouplingRC != "SUCCESS" {
				t.Error(fmt.Printf("PeerMDM %s expected couplingRC SUCCESS but status was %s", srcpeers[i].PeerSystemID, srcpeers[i].CouplingRC))
			} else {
				rep.sourcePeerMDM = goscaleio.NewPeerMDM(C, srcpeers[i])
				t.Logf("PeerMDMID %s", rep.sourcePeerMDM.PeerMDM.ID)
			}
			break
		}
	}
	if !foundTarget {
		t.Error("Didn't find target MDM peer")
	}

	foundSource := false
	for i := 0; i < len(tgtpeers); i++ {
		if tgtpeers[i].PeerSystemID == rep.sourceSystemID {
			foundSource = true
			if tgtpeers[i].CouplingRC != "SUCCESS" {
				t.Error(fmt.Printf("PeerMDM %s expected couplingRC SUCCESS but status was %s", tgtpeers[i].PeerSystemID, tgtpeers[i].CouplingRC))
			}
			break
		}
	}
	if !foundSource {
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
	t.Logf("get ProtectionDomain href %s", href)
	protectionDomains, err := getSystem().GetProtectionDomain(href)
	assert.Nil(t, err)
	for _, pd := range protectionDomains {
		t.Logf("source protection domain %+v", pd)
	}
}

// Get the Target Protection Domain
func getTargetProtectionDomain() *goscaleio.ProtectionDomain {
	TargetProtectionDomainName := os.Getenv("GOSCALEIO_PROTECTIONDOMAIN2")
	protectionDomains, _ := rep.targetSystem.GetProtectionDomain("")
	for i := 0; i < len(protectionDomains); i++ {
		fmt.Printf("target protection domain %+v", protectionDomains[i])
		if protectionDomains[i].Name == TargetProtectionDomainName {
			rep.targetProtectionDomain = goscaleio.NewProtectionDomainEx(C2, protectionDomains[i])
		}
	}
	return rep.targetProtectionDomain
}

func TestTargetProtectionDomain(t *testing.T) {
	if C2 == nil {
		t.Skip("no client connection to replication target system")
	}
	getTargetProtectionDomain()
	assert.NotNil(t, rep.targetProtectionDomain)
	t.Logf("source protction domain: %+v", rep.targetProtectionDomain.ProtectionDomain)
}

// Get the Target StoragePool
func getTargetStoragePool() *goscaleio.StoragePool {
	TargetStoragePoolName := os.Getenv("GOSCALEIO_STORAGEPOOL2")
	storagePools, _ := rep.targetProtectionDomain.GetStoragePool("")
	for i := 0; i < len(storagePools); i++ {
		if storagePools[i].Name == TargetStoragePoolName {
			rep.targetStoragePool = goscaleio.NewStoragePoolEx(C2, storagePools[i])
			return rep.targetStoragePool
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
	getTargetStoragePool()
	assert.NotNil(t, rep.targetStoragePool)
	t.Logf("targetStoragePool %s", rep.targetStoragePool.StoragePool.Name)
}

// Locate the volumes to be replicated together
func TestLocateVolumesToBeReplicated(t *testing.T) {
	if C2 == nil {
		t.Skip("no client connection to replication target system")
	}
	srcName := os.Getenv("GOSCALEIO_REPLICATION_SOURCE_NAME")
	assert.NotNil(t, srcName)
	t.Logf("looking for source %s", srcName)
	sourceVolumes, err := rep.sourceStoragePool.GetVolume("", "", "", srcName, false)
	if err != nil {
		t.Log(err)
	}
	if len(sourceVolumes) > 0 {
		rep.sourceVolume = sourceVolumes[0]
	}
	assert.NotNil(t, rep.sourceVolume)

	dstName := os.Getenv("GOSCALEIO_REPLICATION_TARGET_NAME")
	assert.NotNil(t, dstName)
	t.Logf("looking for target %s", dstName)
	targetVolumes, err := rep.targetStoragePool.GetVolume("", "", "", dstName, false)
	if err != nil {
		t.Log(err)
	}
	if len(targetVolumes) > 0 {
		rep.targetVolume = targetVolumes[0]
	}
	assert.NotNil(t, rep.targetVolume)

	t.Logf("sourceVolume %s targetVolume %s", rep.sourceVolume.Name, rep.targetVolume.Name)
}

// Test createReplicationConsistencyGroup
func TestCreateReplicationConsistencyGroup(t *testing.T) {
	var err error
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
	t.Logf("rcgPayload %+v", rcgPayload)
	log.SetLevel(log.DebugLevel)
	rcgResp, err := C.CreateReplicationConsistencyGroup(rcgPayload)
	if err != nil {
		t.Logf("Error creating RCG: %s", err.Error())
	}
	log.SetLevel(log.InfoLevel)
	assert.Nil(t, err)
	log.Infof("RCG ID: %s", rcgResp.ID)
	rep.rcgID = rcgResp.ID
}

// TestDelayAfterRCGCreation
func TestDelayAfterRCGCreation(t *testing.T) {
	t.Logf("WAITING 30 SECONDS AFTER ATTEMPTING RCG CREATE")
	time.Sleep(5 * time.Second)
}

// Add Replication Pair
func TestAddReplicationPair(t *testing.T) {
	t.Logf("[TestAddReplicationPair] Start")

	var err error
	if C2 == nil {
		t.Skip("[TestAddReplicationPair] no client connection to replication target system")
	}

	// Todo: Create temp volume to test with
	localVolumeID, err := C.FindVolumeID("stkof-snaprw-5-snap-1")
	if err != nil {
		t.Skip("Error finding volume inttest-snap")
	}
	t.Logf("[TestAddReplicationPair] Local Volume ID: %s", localVolumeID)

	remoteVolumeID, err := C2.FindVolumeID("inttest-replication")
	if err != nil {
		t.Skip("Error finding volume inttest-snap")
	}
	t.Logf("[TestAddReplicationPair] Remote Volume ID: %s", remoteVolumeID)

	rpPayload := &siotypes.QueryReplicationPair{
		Name:                          "inttestrp",
		SourceVolumeID:                localVolumeID,
		DestinationVolumeID:           remoteVolumeID,
		ReplicationConsistencyGroupID: rep.rcgID,
		CopyType:                      "OnlineCopy",
	}

	rpResp, err := C.CreateReplicationPair(rpPayload)
	if err != nil {
		t.Logf("[TestAddReplicationPair] Error: %s", err.Error())
	} else {
		t.Logf("[TestAddReplicationPair] Response: %+v", rpResp)
	}

	rep.rpID = rpResp.ID

	t.Logf("[TestAddReplicationPair] End")
}

// Query Replication Pair
func TestQueryReplicationPairs(t *testing.T) {
	t.Logf("[TestQueryReplicationPairs] Start")

	var err error
	if C2 == nil {
		t.Skip("[TestQueryReplicationPairs] no client connection to replication target system")
	}

	pairs, err := C.GetReplicationPairs(rep.rcgID)

	if err != nil {
		t.Logf("[TestQueryReplicationPairs] Error: %s", err.Error())
		return
	}

	for i, pair := range pairs {
		t.Logf("%d, ReplicationPair: %+v", i, pair)
	}

	t.Logf("[TestQueryReplicationPairs] End")
}

// Query Replication Pair Statistics
func TestQueryReplicationPairsStatistics(t *testing.T) {
	t.Logf("[TestQueryReplicationPairsStatistics] Start")

	// var err error
	if C2 == nil {
		t.Skip("[TestQueryReplicationPairsStatistics] no client connection to replication target system")
	}

	for i := 0; i < 20; i++ {
		rpResp, err := C.GetReplicationPairStatistics(rep.rpID)
		if err != nil {
			t.Logf("[TestQueryReplicationPairsStatistics] Error: %s", err.Error())
			break
		}

		t.Logf("[TestQueryReplicationPairsStatistics] Response: %+v", rpResp)

		// Check if complete
		if rpResp.InitialCopyProgress == 1 {
			t.Logf("[TestQueryReplicationPairsStatistics] Copy Complete: %f", rpResp.InitialCopyProgress)
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
	rcgs, err := C.GetReplicationConsistencyGroups("")
	assert.Nil(t, err)
	for i := 0; i < len(rcgs); i++ {
		t.Logf("RCG: %+v\n\n", rcgs[i])
		// t.Logf("Links: %+v\n\n", rcgs[i].Links)
		parseLinks(rcgs[i].Links, t)
		pairs, err := C.GetReplicationPairs(rcgs[i].ID)
		assert.Nil(t, err)
		for j := 0; j < len(pairs); j++ {
			t.Logf("ReplicationPair: %+v", pairs[j])
		}
		if rcgs[i].Name == "inttestrcg" {
			rcg := goscaleio.NewReplicationConsistencyGroup(C)
			rcg.ReplicationConsistencyGroup = rcgs[i]
			rep.rcg = rcg
		}
	}
	assert.NotNil(t, rep.rcg)
}

// Test RemoveReplicationPair
func TestRemoveReplicationPair(t *testing.T) {
	if C2 == nil {
		t.Skip("no client connection to replication target system")
	}

	t.Logf("[TestRemoveReplicationPair] Removing replication pair: %s", rep.rpID)

	pair, err := C.RemoveReplicationPair(rep.rpID)
	assert.Nil(t, err)
	assert.NotNil(t, pair)

	t.Logf("[TestRemoveReplicationPair] Removed the following pair: %+v", pair)

	// Delay to verify on the UI.
	time.Sleep(10 * time.Second)
}

// Test RemoveReplicatonConsistencyGroup
func TestRemoveReplicationConsistencyGroup(t *testing.T) {
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
