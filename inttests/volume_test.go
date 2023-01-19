// Copyright © 2021 - 2022 Dell Inc. or its subsidiaries. All Rights Reserved.
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

package inttests

import (
	"fmt"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/dell/goscaleio"
	siotypes "github.com/dell/goscaleio/types/v1"
	"github.com/stretchr/testify/assert"
)

var (
	createdVolumes = make([]string, 0)
)

func getVolByID(id string) (*siotypes.Volume, error) {

	// The `GetVolume` API returns a slice of volumes, but when only passing
	// in a volume ID, the response will be just the one volume
	vols, err := C.GetVolume("", strings.TrimSpace(id), "", "", false)
	if err != nil {
		return nil, err
	}
	return vols[0], nil
}

func createVolume(t *testing.T, useName string) (string, error) {
	pool := getStoragePool(t)
	if pool == nil {
		return "", fmt.Errorf("Error when getting storagepool")
	}

	name := useName
	if useName == "" {
		name = getUniqueName()
	}

	// Create a volume
	volumeParam := &siotypes.VolumeParam{
		Name:           name,
		VolumeSizeInKb: fmt.Sprintf("%d", defaultVolumeSize),
		VolumeType:     "ThinProvisioned",
	}
	createResp, err := pool.CreateVolume(volumeParam)
	if err != nil {
		return "", fmt.Errorf("error when creating volume %s storagepool %s: %s", name, pool.StoragePool.Name, err.Error())
	}
	// add this voluem to slice of created volumes
	createdVolumes = append(createdVolumes, createResp.ID)
	return createResp.ID, nil
}

func deleteVolume(t *testing.T, volID string) error {
	existingVol, err := getVolByID(volID)
	if err != nil {
		return err
	}
	vol := goscaleio.NewVolume(C)
	vol.Volume = existingVol
	// by default, Remove volume will remove "ONLY_ME"
	err = vol.RemoveVolume("")
	if err != nil {
		return err
	}
	// remove the volume from the created volume slice
	existingVols := make([]string, 0)
	for _, v := range createdVolumes {
		if v != volID {
			existingVols = append(existingVols, v)
		}
	}
	createdVolumes = existingVols
	return nil
}

func deleteAllVolumes(t *testing.T) error {
	for _, v := range createdVolumes {
		deleteVolume(t, v)
	}
	return nil
}

func TestGetVolumes(t *testing.T) {
	volID, err := createVolume(t, "")
	assert.Nil(t, err)
	newVolume, err := getVolByID(volID)
	assert.Nil(t, err)

	// now make a snapshot
	snapshotDefs := make([]*siotypes.SnapshotDef, 0)
	snapname := fmt.Sprintf("%s-%s", newVolume.Name, "snap")

	snapDef := &siotypes.SnapshotDef{
		VolumeID:     volID,
		SnapshotName: snapname,
	}
	snapshotDefs = append(snapshotDefs, snapDef)
	snapParam := &siotypes.SnapshotVolumesParam{
		SnapshotDefs: snapshotDefs,
	}

	system := getSystem()
	assert.NotNil(t, system)

	// Create snapshot
	snapResponse, err := system.CreateSnapshotConsistencyGroup(snapParam)
	assert.Nil(t, err)
	assert.NotZero(t, len(snapResponse.VolumeIDList))

	pool := getStoragePool(t)
	volumes, err := pool.GetVolume("", volID, "", "", true)
	assert.Nil(t, err)
	assert.NotNil(t, volumes)

	// get volume by name
	volumes, err = pool.GetVolume("", "", "", newVolume.Name, true)
	assert.Nil(t, err)
	assert.NotNil(t, volumes)

	// get volume by ID
	volumes, err = pool.GetVolume("", newVolume.ID, "", "", true)
	assert.Nil(t, err)
	assert.NotNil(t, volumes)

	// get volume by href
	href := fmt.Sprintf("/api/instances/Volume::%s", newVolume.ID)
	volumes, err = pool.GetVolume(href, "", "", "", true)
	assert.Nil(t, err)
	assert.NotNil(t, volumes)

	// get volume by ancestor ID
	volumes, err = pool.GetVolume("", "", newVolume.ID, "", true)
	assert.Nil(t, err)
	assert.NotNil(t, volumes)

	// get the snapshots
	volumes, err = pool.GetVolume("", snapResponse.VolumeIDList[0], "", "", true)
	assert.Nil(t, err)
	assert.NotNil(t, volumes)

	// delete the snapshots
	for _, s := range snapResponse.VolumeIDList {
		deleteVolume(t, s)
	}

	err = deleteVolume(t, volID)
	assert.Nil(t, err)
	deleteAllVolumes(t)
}

func TestFindVolumeID(t *testing.T) {
	name := fmt.Sprintf("%s-%s", testPrefix, "getByID")
	// create a volume
	volID, err := createVolume(t, name)
	assert.Nil(t, err)
	assert.NotNil(t, volID)

	found, err := C.FindVolumeID(name)
	assert.Nil(t, err)
	assert.NotNil(t, found)
	assert.Equal(t, volID, found)

	deleteVolume(t, volID)
	deleteAllVolumes(t)
}

func TestCreateDeleteVolume(t *testing.T) {
	volID, err := createVolume(t, "")
	assert.Nil(t, err)
	err = deleteVolume(t, volID)
	assert.Nil(t, err)
	deleteAllVolumes(t)
}

func TestCreateVolumeExistingName(t *testing.T) {
	name := fmt.Sprintf("%s-%s", testPrefix, "existingVol")
	// create a volume
	_, err := createVolume(t, name)
	assert.Nil(t, err)
	// attempt create another volume with that name, it should fail
	_, err = createVolume(t, name)
	assert.NotNil(t, err)

	deleteAllVolumes(t)
}

func TestDeleteNonExistentVolume(t *testing.T) {
	// create a volume
	volID, err := createVolume(t, "")
	assert.Nil(t, err)
	// delete the volume, it should complete
	err = deleteVolume(t, volID)
	assert.Nil(t, err)
	// delete the volume again, it should fail
	err = deleteVolume(t, volID)
	assert.NotNil(t, err)

	deleteAllVolumes(t)
}

func TestCreateDeleteSnapshot(t *testing.T) {
	name := fmt.Sprintf("%s-%s", testPrefix, "toBeSnapped")
	volID, err := createVolume(t, "")
	assert.Nil(t, err)

	vol, err := getVolByID(volID)
	assert.Nil(t, err)
	assert.NotNil(t, vol)

	volume := goscaleio.NewVolume(C)
	volume.Volume = vol

	snapshotDefs := make([]*siotypes.SnapshotDef, 0)
	snapname := fmt.Sprintf("%s-%s", name, "snap")

	snapDef := &siotypes.SnapshotDef{
		VolumeID:     volID,
		SnapshotName: snapname,
	}
	snapshotDefs = append(snapshotDefs, snapDef)
	snapParam := &siotypes.SnapshotVolumesParam{
		SnapshotDefs: snapshotDefs,
	}

	system := getSystem()
	assert.NotNil(t, system)

	// Create snapshot
	snapResponse, err := system.CreateSnapshotConsistencyGroup(snapParam)
	assert.Nil(t, err)
	assert.NotZero(t, len(snapResponse.VolumeIDList))
	// delete the snapshots
	for _, s := range snapResponse.VolumeIDList {
		deleteVolume(t, s)
	}

	deleteAllVolumes(t)
}

func TestGetVolumeVtree(t *testing.T) {
	volID, err := createVolume(t, "")
	assert.Nil(t, err)

	vol, err := getVolByID(volID)
	assert.Nil(t, err)
	assert.NotNil(t, vol)

	volume := goscaleio.NewVolume(C)
	volume.Volume = vol

	// get a valid vtree
	vtree, err := volume.GetVTree()
	assert.Nil(t, err)
	assert.NotNil(t, vtree)
	assert.Equal(t, volume.Volume.Name, vtree.Name)

	// attempt to get the VTree again, this time with a non-existent volume
	badVolume := goscaleio.NewVolume(C)
	vtree, err = badVolume.GetVTree()
	assert.NotNil(t, err)
	assert.Nil(t, vtree)

	err = deleteVolume(t, volID)
	assert.Nil(t, err)

	deleteAllVolumes(t)
}

func TestGetVolumeStatistics(t *testing.T) {
	volID, err := createVolume(t, "")
	assert.Nil(t, err)

	vol, err := getVolByID(volID)
	assert.Nil(t, err)
	assert.NotNil(t, vol)

	volume := goscaleio.NewVolume(C)
	volume.Volume = vol
	stats, err := volume.GetVolumeStatistics()
	assert.Nil(t, err)
	assert.NotNil(t, stats)

	// attempt to get the statistics again, this time with a non-existent volume
	badVolume := goscaleio.NewVolume(C)
	stats, err = badVolume.GetVolumeStatistics()
	assert.NotNil(t, err)
	assert.Nil(t, stats)

	err = deleteVolume(t, volID)
	assert.Nil(t, err)
	deleteAllVolumes(t)
}

func TestResizeVolume(t *testing.T) {
	volID, err := createVolume(t, "")
	assert.Nil(t, err)

	vol, err := getVolByID(volID)
	assert.Nil(t, err)
	assert.NotNil(t, vol)

	volume := goscaleio.NewVolume(C)
	volume.Volume = vol
	existingSizeGB := volume.Volume.SizeInKb / (1024 * 1024)
	newSize := existingSizeGB * 2
	// double the szie of the volume
	err = volume.SetVolumeSize(strconv.Itoa(int(newSize)))

	volumeTemp, err := getVolByID(volID)
	assert.Nil(t, err)
	assert.NotNil(t, volumeTemp)
	volumeResized := goscaleio.NewVolume(C)
	volumeResized.Volume = volumeTemp
	assert.Equal(t, existingSizeGB*2, volumeResized.Volume.SizeInKb/(1024*1024))

	err = deleteVolume(t, volID)
	assert.Nil(t, err)
	deleteAllVolumes(t)
}

func TestMapQueryUnmapVolume(t *testing.T) {
	volID, err := createVolume(t, "")
	assert.Nil(t, err)

	vol, err := getVolByID(volID)
	assert.Nil(t, err)
	assert.NotNil(t, vol)
	volume := goscaleio.NewVolume(C)
	volume.Volume = vol

	// get the SDCs and pick one...
	sdcs := getAllSdc(t)
	assert.NotEqual(t, 0, len(sdcs))

	chosenSDC := sdcs[0]

	mapVolumeSdcParam := &siotypes.MapVolumeSdcParam{
		SdcID:                 chosenSDC.Sdc.ID,
		AllowMultipleMappings: "FALSE",
		AllSdcs:               "",
	}
	volume.MapVolumeSdc(mapVolumeSdcParam)

	stats, err := volume.GetVolumeStatistics()
	assert.Nil(t, err)
	assert.NotNil(t, stats)

	unmapVolumeSdcParam := &siotypes.UnmapVolumeSdcParam{
		SdcID:   chosenSDC.Sdc.ID,
		AllSdcs: "",
	}

	err = volume.UnmapVolumeSdc(unmapVolumeSdcParam)
	assert.Nil(t, err)

	err = deleteVolume(t, volID)
	assert.Nil(t, err)
	deleteAllVolumes(t)
}

func TestMapQueryUnmapSnapshot(t *testing.T) {
	volID, err := createVolume(t, "")
	assert.Nil(t, err)
	newVolume, err := getVolByID(volID)
	assert.Nil(t, err)

	// now make a snapshot
	snapshotDefs := make([]*siotypes.SnapshotDef, 0)
	snapname := fmt.Sprintf("%s-%s", newVolume.Name, "snap")

	snapDef := &siotypes.SnapshotDef{
		VolumeID:     volID,
		SnapshotName: snapname,
	}
	snapshotDefs = append(snapshotDefs, snapDef)
	snapParam := &siotypes.SnapshotVolumesParam{
		SnapshotDefs: snapshotDefs,
	}

	system := getSystem()
	assert.NotNil(t, system)

	// Create snapshot
	snapResponse, err := system.CreateSnapshotConsistencyGroup(snapParam)
	assert.Nil(t, err)
	assert.NotZero(t, len(snapResponse.VolumeIDList))

	// Get StoragePool
	pool := getStoragePool(t)
	volumes, err := pool.GetVolume("", volID, "", "", true)
	assert.Nil(t, err)
	assert.NotNil(t, volumes)

	// Get Snapshot
	volumes, err = pool.GetVolume("", snapResponse.VolumeIDList[0], "", "", true)
	assert.Nil(t, err)
	assert.NotNil(t, volumes)

	// Lock Snapshot
	snap, err := getVolByID(volumes[0].ID)
	assert.Nil(t, err)
	sr := goscaleio.NewVolume(C)
	sr.Volume = snap
	err = sr.SetVolumeAccessModeLimit("ReadWrite")
	assert.Nil(t, err)
	// testing invalid case
	err = sr.SetVolumeAccessModeLimit(invalidIdentifier)
	assert.NotNil(t, err)

	// get the SDCs and pick one...
	sdcs := getAllSdc(t)
	assert.NotEqual(t, 0, len(sdcs))
	chosenSDC := sdcs[0]
	mapVolumeSdcParam := &siotypes.MapVolumeSdcParam{
		SdcID:                 chosenSDC.Sdc.ID,
		AllowMultipleMappings: "FALSE",
		AllSdcs:               "",
	}
	err = sr.MapVolumeSdc(mapVolumeSdcParam)
	assert.Nil(t, err)

	unmapVolumeSdcParam := &siotypes.UnmapVolumeSdcParam{
		SdcID:   chosenSDC.Sdc.ID,
		AllSdcs: "",
	}
	sr.UnmapVolumeSdc(unmapVolumeSdcParam)
	assert.Nil(t, err)

	// Delete Snapshot and Volume
	err = deleteVolume(t, sr.Volume.ID)
	assert.Nil(t, err)
	err = deleteVolume(t, volID)
	assert.Nil(t, err)

}

func TestCreateInstanceVolume(t *testing.T) {
	name := fmt.Sprintf("%s-%s", testPrefix, "instanceCreated")

	poolName := getStoragePoolName(t)
	assert.NotNil(t, poolName)

	size := fmt.Sprintf("%d", defaultVolumeSize)

	volParams := siotypes.VolumeParam{
		VolumeSizeInKb: size,
		VolumeType:     "ThinProvisioned",
		Name:           name,
	}

	volResp, err := C.CreateVolume(&volParams, poolName, "")
	assert.Nil(t, err)
	assert.NotNil(t, volResp)

	deleteVolume(t, volResp.ID)
}

func TestCreateInstanceVolumeInvalidSize(t *testing.T) {
	name := fmt.Sprintf("%s-%s", testPrefix, "instanceCreated")

	poolName := getStoragePoolName(t)
	assert.NotNil(t, poolName)

	volParams := siotypes.VolumeParam{
		VolumeSizeInKb: "0",
		VolumeType:     "ThinProvisioned",
		Name:           name,
	}

	volResp, err := C.CreateVolume(&volParams, poolName, "")
	assert.NotNil(t, err)
	assert.Nil(t, volResp)
}

func TestGetInstanceVolume(t *testing.T) {
	volID, err := createVolume(t, "")
	assert.Nil(t, err)

	thisVolume, err := getVolByID(volID)
	assert.Nil(t, err)

	// get by ID
	volume, err := C.GetVolume("", volID, "", "", true)
	assert.Nil(t, err)
	assert.NotNil(t, volume)

	// Find by name
	volume, err = C.GetVolume("", "", "", thisVolume.Name, true)
	assert.Nil(t, err)
	assert.NotNil(t, volume)

	// Find by href
	href := fmt.Sprintf("/api/instances/Volume::%s", volID)
	volume, err = C.GetVolume(href, "", "", "", true)
	assert.Nil(t, err)
	assert.NotNil(t, volume)

	// Find with invalid name
	volume, err = C.GetVolume("", "", "", invalidIdentifier, true)
	assert.NotNil(t, err)
	assert.Nil(t, volume)

	// Find with invalid ID
	volume, err = C.GetVolume(invalidIdentifier, "", "", "", true)
	assert.NotNil(t, err)
	assert.Nil(t, volume)

	// Find with an invalid href
	href = fmt.Sprintf("/api/BAD/instances/Volume::%s", volID)
	volume, err = C.GetVolume(href, "", "", "", true)
	assert.NotNil(t, err)
	assert.Nil(t, volume)

	deleteAllVolumes(t)
}

// TestSetMappedSdcLimitsInvalid will attempt to set SDC limits against an invalid SDC
func TestSetMappedSdcLimitsInvalid(t *testing.T) {
	volID, err := createVolume(t, "")
	assert.Nil(t, err)

	typeVolume, err := getVolByID(volID)
	assert.Nil(t, err)

	thisVolume := goscaleio.NewVolume(C)
	thisVolume.Volume = typeVolume

	settings := siotypes.SetMappedSdcLimitsParam{
		SdcID:                invalidIdentifier,
		BandwidthLimitInKbps: "0",
		IopsLimit:            "0",
	}

	err = thisVolume.SetMappedSdcLimits(&settings)
	assert.NotNil(t, err)

	err = deleteVolume(t, volID)
	assert.Nil(t, err)
	deleteAllVolumes(t)
}

// Testing TestLockUnlockAutoSnapshot will attempting locking the auto snapshot and unlocking the auto snapshot
func TestLockUnlockAutoSnapshot(t *testing.T) {
	volID, err := createVolume(t, "")
	assert.Nil(t, err)
	newVolume, err := getVolByID(volID)
	assert.Nil(t, err)

	// now make a snapshot
	snapshotDefs := make([]*siotypes.SnapshotDef, 0)
	snapname := fmt.Sprintf("%s-%s", newVolume.Name, "snap")

	snapDef := &siotypes.SnapshotDef{
		VolumeID:     volID,
		SnapshotName: snapname,
	}
	snapshotDefs = append(snapshotDefs, snapDef)
	snapParam := &siotypes.SnapshotVolumesParam{
		SnapshotDefs: snapshotDefs,
	}

	system := getSystem()
	assert.NotNil(t, system)

	// Create snapshot
	snapResponse, err := system.CreateSnapshotConsistencyGroup(snapParam)
	assert.Nil(t, err)
	assert.NotZero(t, len(snapResponse.VolumeIDList))

	// Get StoragePool
	pool := getStoragePool(t)
	volumes, err := pool.GetVolume("", volID, "", "", true)
	assert.Nil(t, err)
	assert.NotNil(t, volumes)

	// Get Snapshot
	volumes, err = pool.GetVolume("", snapResponse.VolumeIDList[0], "", "", true)
	assert.Nil(t, err)
	assert.NotNil(t, volumes)

	// Lock Snapshot
	snap, err := getVolByID(volumes[0].ID)
	assert.Nil(t, err)
	sr := goscaleio.NewVolume(C)
	sr.Volume = snap
	err = sr.LockAutoSnapshot()
	assert.NotNil(t, err)
	err = sr.UnlockAutoSnapshot()
	assert.NotNil(t, err)

	// Delete Snapshot and Volume
	err = deleteVolume(t, sr.Volume.ID)
	assert.Nil(t, err)
	err = deleteVolume(t, volID)
	assert.Nil(t, err)
	deleteAllVolumes(t)
}

// Testing TestSetVolumeAccessModeLimit will be attempting set access mode of volume
func TestSetVolumeAccessModeLimit(t *testing.T) {
	volID, err := createVolume(t, "")
	assert.Nil(t, err)
	vol, err := getVolByID(volID)
	assert.Nil(t, err)
	vr := goscaleio.NewVolume(C)
	vr.Volume = vol
	err = vr.SetVolumeAccessModeLimit("ReadOnly")
	assert.Nil(t, err)
	// testing invalid case
	err = vr.SetVolumeAccessModeLimit(invalidIdentifier)
	assert.NotNil(t, err)

	err = deleteVolume(t, volID)
	assert.Nil(t, err)
}

// TestSetSnapshotSecurity will be attemting to set the snapshot security for a snapshot
func TestSetSnapshotSecurity(t *testing.T) {
	volID, err := createVolume(t, "")
	assert.Nil(t, err)
	newVolume, err := getVolByID(volID)
	assert.Nil(t, err)

	// now make a snapshot
	snapshotDefs := make([]*siotypes.SnapshotDef, 0)
	snapname := fmt.Sprintf("%s-%s", newVolume.Name, "snap2")

	snapDef := &siotypes.SnapshotDef{
		VolumeID:     volID,
		SnapshotName: snapname,
	}
	snapshotDefs = append(snapshotDefs, snapDef)
	snapParam := &siotypes.SnapshotVolumesParam{
		SnapshotDefs: snapshotDefs,
	}

	system := getSystem()
	assert.NotNil(t, system)

	// Create snapshot
	snapResponse, err := system.CreateSnapshotConsistencyGroup(snapParam)
	assert.Nil(t, err)
	assert.NotZero(t, len(snapResponse.VolumeIDList))
	// Get StoragePool
	pool := getStoragePool(t)
	volumes, err := pool.GetVolume("", volID, "", "", true)
	assert.Nil(t, err)
	assert.NotNil(t, volumes)

	// Get Snapshot
	volumes, err = pool.GetVolume("", snapResponse.VolumeIDList[0], "", "", true)
	assert.Nil(t, err)
	assert.NotNil(t, volumes)

	// Set a new retention period for the given snapshot
	snap, err := getVolByID(volumes[0].ID)
	assert.Nil(t, err)
	sr := goscaleio.NewVolume(C)
	sr.Volume = snap
	err = sr.SetSnapshotSecurity("0")
	assert.Nil(t, err)
	// testing invalid case
	err = sr.SetSnapshotSecurity(invalidIdentifier)
	assert.NotNil(t, err)
	// Delete Snapshot and Volume
	fmt.Println("Will wait for 60 sec so that the retention period expires and snapshot can be deleted")
	time.Sleep(60 * time.Second)
	err = deleteVolume(t, sr.Volume.ID)
	assert.Nil(t, err)
	err = deleteVolume(t, volID)
	assert.Nil(t, err)
}

// TestSetVolumeMappingAccessMode will be attemting to set the access mode on mapped sdc
func TestSetVolumeMappingAccessMode(t *testing.T) {
	volID, err := createVolume(t, "")
	assert.Nil(t, err)
	newVolume, err := getVolByID(volID)
	assert.Nil(t, err)

	// now make a snapshot
	snapshotDefs := make([]*siotypes.SnapshotDef, 0)
	snapname := fmt.Sprintf("%s-%s", newVolume.Name, "snap3")

	snapDef := &siotypes.SnapshotDef{
		VolumeID:     volID,
		SnapshotName: snapname,
	}
	snapshotDefs = append(snapshotDefs, snapDef)
	snapParam := &siotypes.SnapshotVolumesParam{
		SnapshotDefs: snapshotDefs,
	}

	system := getSystem()
	assert.NotNil(t, system)

	// Create snapshot
	snapResponse, err := system.CreateSnapshotConsistencyGroup(snapParam)
	assert.Nil(t, err)
	assert.NotZero(t, len(snapResponse.VolumeIDList))
	// Get StoragePool
	pool := getStoragePool(t)
	volumes, err := pool.GetVolume("", volID, "", "", true)
	assert.Nil(t, err)
	assert.NotNil(t, volumes)

	// Get Snapshot
	volumes, err = pool.GetVolume("", snapResponse.VolumeIDList[0], "", "", true)
	assert.Nil(t, err)
	assert.NotNil(t, volumes)

	// Set a access mode for the given snapshot
	snap, err := getVolByID(volumes[0].ID)
	assert.Nil(t, err)
	sr := goscaleio.NewVolume(C)
	sr.Volume = snap
	pfmvsp := &siotypes.MapVolumeSdcParam{
		SdcID:                 "c423b09800000003",
		AllowMultipleMappings: "true",
	}
	sr.MapVolumeSdc(pfmvsp)
	err = sr.SetVolumeMappingAccessMode("ReadWrite", "c423b09800000003")
	assert.Nil(t, err)
	// testing invalid case
	err = sr.SetVolumeMappingAccessMode(invalidIdentifier, invalidIdentifier)
	assert.NotNil(t, err)
	// Delete Snapshot and Volume
	sr.UnmapVolumeSdc(
		&siotypes.UnmapVolumeSdcParam{
			SdcID: "c423b09800000003",
		},
	)
	err = deleteVolume(t, sr.Volume.ID)
	assert.Nil(t, err)
	err = deleteVolume(t, volID)
	assert.Nil(t, err)
}

// Testing TestSetVolumeUseRmCache will be attempting set use rm cache
func TestSetVolumeUseRmCache(t *testing.T) {
	volID, err := createVolume(t, "")
	assert.Nil(t, err)
	vol, err := getVolByID(volID)
	assert.Nil(t, err)
	vr := goscaleio.NewVolume(C)
	vr.Volume = vol
	err = vr.SetVolumeUseRmCache(true)
	assert.Nil(t, err)
	err = deleteVolume(t, volID)
	assert.Nil(t, err)
}

// Testing TestSetCompressionMethod will be attempting set compression method
func TestSetCompressionMethod(t *testing.T) {
	volID, err := createVolume(t, "")
	assert.Nil(t, err)
	vol, err := getVolByID(volID)
	assert.Nil(t, err)
	vr := goscaleio.NewVolume(C)
	vr.Volume = vol
	// set compression method will only get pass for snapshot with fine granularity
	err = vr.SetCompressionMethod("None")
	assert.NotNil(t,err)
	// testing invalid case
	err = vr.SetCompressionMethod(invalidIdentifier)
	assert.NotNil(t, err)
	err = deleteVolume(t, volID)
	assert.Nil(t, err)
}
