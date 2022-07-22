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
	"github.com/stretchr/testify/assert"
	"github.com/dell/goscaleio"
	siotypes "github.com/dell/goscaleio/types/v1"
	"time"
)

// Replication global variables used to set up the replication relationships
type replication struct {
	sourcePeerMDM *goscaleio.PeerMDM
	targetPeerMDM *goscaleio.PeerMDM
	sourceSystem *goscaleio.System
	targetSystem *goscaleio.System
	sourceSystemID string
	sourceProtectionDomainID string
	sourceProtectionDomain *goscaleio.ProtectionDomain
	sourceStoragePool *goscaleio.StoragePool
	sourceVolume *siotypes.Volume
	targetSystemID string
	targetProtectionDomain *goscaleio.ProtectionDomain
	targetStoragePool *goscaleio.StoragePool
	targetVolume *siotypes.Volume
	rcg *goscaleio.ReplicationConsistencyGroup
}
var rep replication

// Test GetPeerMDMs
func TestGetPeerMDMs(t *testing.T) {
	srcpeers, err := C.GetPeerMDMs()
	assert.Nil(t, err)
	for i:=0; i < len(srcpeers); i++ {
		t.Logf("Source PeerMDM: %+v", srcpeers[i])
		rep.sourceSystemID = srcpeers[i].SystemID
	}

	if C2 == nil {
		t.Skip("no client connection to replication target system")
	}
	tgtpeers, err := C2.GetPeerMDMs()
	assert.Nil(t, err)
	for i:=0; i < len(tgtpeers); i++ {
		t.Logf("Target PeerMDM: %+v", tgtpeers[i])
		rep.targetSystemID = tgtpeers[i].SystemID
	}

	// Test systems are validly paired
	foundTarget := false
	for i:=0; i < len(srcpeers); i++ {
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
	for i:=0; i < len(tgtpeers); i++ {
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
	for _,pd := range protectionDomains {
		t.Logf("source protection domain %+v", pd)
	}
}

// Get the Target Protection Domain
func getTargetProtectionDomain() *goscaleio.ProtectionDomain {
	TargetProtectionDomainName := os.Getenv("GOSCALEIO_PROTECTIONDOMAIN2")
	protectionDomains, _ := rep.targetSystem.GetProtectionDomain("")
	for i:=0; i < len(protectionDomains); i++ {
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
	for i:=0; i < len(storagePools); i++ {
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

// TestDelayBeforeRCGCreation
func TestDelayBeforeRCGCreation(t *testing.T) {
	t.Logf("WAITING 30 SECONDS BEFORE ATTEMPTING RCG CREATE")
	time.Sleep(1 * time.Second)
}


// Test createReplicationConsistencyGroup
func TestCreateReplicationConsistencyGroup(t *testing.T) {
	var err error
	if C2 == nil {
		t.Skip("no client connection to replication target system")
	}
	rcgPayload := &siotypes.ReplicationConsistencyGroupCreatePayload {
		Name: "inttestrcg",
		RpoInSeconds: "60",
		ProtectionDomainId: rep.sourceProtectionDomain.ProtectionDomain.ID,
		RemoteProtectionDomainId: rep.targetProtectionDomain.ProtectionDomain.ID,
		PeerMdmId: rep.sourcePeerMDM.PeerMDM.ID, 
		//DestinationSystemId: rep.targetSystem.System.ID,
	}
	t.Logf("rcgPayload %+v", rcgPayload)
	rep.rcg, err = C2.CreateReplicationConsistencyGroup(rcgPayload)
	if err != nil {
		t.Logf("Error creating RCG: %s", err.Error())
	}
	assert.Nil(t, err)
}

// TestDelayAfterRCGCreation
func TestDelayAfterRCGCreation(t *testing.T) {
	t.Logf("WAITING 30 SECONDS AFTER ATTEMPTING RCG CREATE")
	time.Sleep(5 * time.Second)
}

// Test GetReplicationConsistencyGroups
func TestGetReplicationConsistencyGroups(t *testing.T) {
	if C2 == nil {
		t.Skip("no client connection to replication target system")
	}
	rcgs, err := C.GetReplicationConsistencyGroups()
	assert.Nil(t, err)
	for i:=0; i < len(rcgs); i++ {
		t.Logf("RCG: %+v", rcgs[i])
		pairs, err := C.GetReplicationPairs(rcgs[i].ID)
		assert.Nil(t, err)
		for j:=0; j < len(pairs); j++ {
			t.Logf("ReplicationPair: %+v", pairs[j])
		}

	}
}

// Test
