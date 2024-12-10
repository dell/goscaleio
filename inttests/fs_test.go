// Copyright © 2021 - 2023 Dell Inc. or its subsidiaries. All Rights Reserved.
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
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/dell/goscaleio"
	// log "github.com/sirupsen/logrus"
	types "github.com/dell/goscaleio/types/v1"
	"github.com/stretchr/testify/assert"
)

// getFileSystemName returns GOSCALEIO_FILESYSTEM, if set
// if not set, returns the file system domain found
func getFileSystemName(t *testing.T) string {
	if os.Getenv("GOSCALEIO_FILESYSTEM") != "" {
		return os.Getenv("GOSCALEIO_FILESYSTEM")
	}
	system := getSystem()
	assert.NotNil(t, system)
	filesystems, _ := system.GetAllFileSystems(context.Background())
	assert.NotNil(t, filesystems)
	if filesystems == nil {
		return ""
	}
	fmt.Printf("filesystems[0].Name: %v", filesystems[0].Name)
	return filesystems[0].Name
}

// getAllFileSystems will return all file system instances
func getAllFileSystems(t *testing.T) []*goscaleio.FileSystem {
	system := getSystem()
	assert.NotNil(t, system)
	if system == nil {
		return nil
	}

	var allFs []*goscaleio.FileSystem
	fs, err := system.GetAllFileSystems(context.Background())
	assert.Nil(t, err)
	assert.NotZero(t, len(fs))
	for _, f := range fs {
		outFs := goscaleio.NewFileSystem(C, &f) // #nosec G601
		allFs = append(allFs, outFs)
	}
	return allFs
}

// TestGetAllFileSystems will return all file system instances
func TestGetAllFileSystems(t *testing.T) {
	filesystems := getAllFileSystems(t)
	assert.NotNil(t, filesystems)
	assert.NotZero(t, len(filesystems))
}

func TestCreateFileSystemSnapshot(t *testing.T) {
	system := getSystem()
	assert.NotNil(t, system)

	fsName := fmt.Sprintf("%s-%s", "FS", testPrefix)

	// get protection domain
	pd := getProtectionDomain(t)
	assert.NotNil(t, pd)

	ctx := context.Background()

	// get storage pool
	pool := getStoragePool(t)
	assert.NotNil(t, pool)
	var spID string
	if pd != nil && pool != nil {
		sp, _ := pd.FindStoragePool(ctx, pool.StoragePool.ID, "", "")
		assert.NotNil(t, sp)
		spID = sp.ID
	}

	// get NAS server ID
	var nasServerName string
	if os.Getenv("GOSCALEIO_NASSERVER") != "" {
		nasServerName = os.Getenv("GOSCALEIO_NASSERVER")
	}
	nasServer, err := system.GetNASByIDName(ctx, "", nasServerName)
	assert.NotNil(t, nasServer)
	assert.Nil(t, err)

	fs := &types.FsCreate{
		Name:          fsName,
		SizeTotal:     16106127360,
		StoragePoolID: spID,
		NasServerID:   nasServer.ID,
	}

	// create the file system
	filesystem, err := system.CreateFileSystem(ctx, fs)
	fsID := filesystem.ID
	assert.Nil(t, err)
	assert.NotNil(t, fsID)

	snapResp, err := system.CreateFileSystemSnapshot(ctx, &types.CreateFileSystemSnapshotParam{
		Name: "test-snapshot",
	}, fsID)

	assert.NotNil(t, snapResp)
	assert.Nil(t, err)

	snap, err := system.GetFileSystemByIDName(ctx, snapResp.ID, "")
	assert.NotNil(t, snap)
	assert.Nil(t, err)

	err = system.DeleteFileSystem(ctx, snap.Name)
	assert.Nil(t, err)
	err = system.DeleteFileSystem(ctx, fs.Name)
	assert.Nil(t, err)
}

func TestRestoreFileSystemSnapshot(t *testing.T) {
	system := getSystem()
	assert.NotNil(t, system)

	fsName := fmt.Sprintf("%s-%s", "FS", testPrefix)

	// get protection domain
	pd := getProtectionDomain(t)
	assert.NotNil(t, pd)

	ctx := context.Background()

	// get storage pool
	pool := getStoragePool(t)
	assert.NotNil(t, pool)
	var spID string
	if pd != nil && pool != nil {
		sp, _ := pd.FindStoragePool(ctx, pool.StoragePool.ID, "", "")
		assert.NotNil(t, sp)
		spID = sp.ID
	}

	// get NAS server ID
	var nasServerName string
	if os.Getenv("GOSCALEIO_NASSERVER") != "" {
		nasServerName = os.Getenv("GOSCALEIO_NASSERVER")
	}
	nasServer, err := system.GetNASByIDName(ctx, "", nasServerName)
	assert.NotNil(t, nasServer)
	assert.Nil(t, err)

	fs := &types.FsCreate{
		Name:          fsName,
		SizeTotal:     16106127360,
		StoragePoolID: spID,
		NasServerID:   nasServer.ID,
	}

	// create the file system
	filesystem, err := system.CreateFileSystem(ctx, fs)
	fsID := filesystem.ID
	assert.Nil(t, err)
	assert.NotNil(t, fsID)
	snapName := fmt.Sprintf("%s-%s", "SNAP", testPrefix)

	snapResp, err := system.CreateFileSystemSnapshot(ctx, &types.CreateFileSystemSnapshotParam{
		Name: snapName,
	}, fsID)

	assert.NotNil(t, snapResp)
	assert.Nil(t, err)

	snap, err := system.GetFileSystemByIDName(ctx, snapResp.ID, "")
	assert.NotNil(t, snap)
	assert.Nil(t, err)

	restoreSnap, err := system.RestoreFileSystemFromSnapshot(ctx, &types.RestoreFsSnapParam{
		SnapshotID: snap.ID,
	}, fsID)

	assert.Nil(t, restoreSnap)
	assert.Nil(t, err)

	restoreSnap, err = system.RestoreFileSystemFromSnapshot(ctx, &types.RestoreFsSnapParam{
		SnapshotID: snap.ID,
		CopyName:   "copy" + snapName,
	}, fsID)

	assert.NotNil(t, restoreSnap)
	assert.Nil(t, err)

	rSnap, err := system.GetFileSystemByIDName(ctx, restoreSnap.ID, "")

	assert.NotNil(t, rSnap)
	assert.Nil(t, err)

	err = system.DeleteFileSystem(ctx, rSnap.Name)
	assert.Nil(t, err)

	err = system.DeleteFileSystem(ctx, snap.Name)
	assert.Nil(t, err)

	err = system.DeleteFileSystem(ctx, fs.Name)
	assert.Nil(t, err)
}

func TestGetFsSnapshotsByVolumeID(t *testing.T) {
	system := getSystem()
	assert.NotNil(t, system)

	fsName := fmt.Sprintf("%s-%s", "FS", testPrefix)

	// get protection domain
	pd := getProtectionDomain(t)
	assert.NotNil(t, pd)

	ctx := context.Background()

	// get storage pool
	pool := getStoragePool(t)
	assert.NotNil(t, pool)
	var spID string
	if pd != nil && pool != nil {
		sp, _ := pd.FindStoragePool(ctx, pool.StoragePool.ID, "", "")
		assert.NotNil(t, sp)
		spID = sp.ID
	}

	// get NAS server ID
	var nasServerName string
	if os.Getenv("GOSCALEIO_NASSERVER") != "" {
		nasServerName = os.Getenv("GOSCALEIO_NASSERVER")
	}
	nasServer, err := system.GetNASByIDName(ctx, "", nasServerName)
	assert.NotNil(t, nasServer)
	assert.Nil(t, err)

	fs := &types.FsCreate{
		Name:          fsName,
		SizeTotal:     16106127360,
		StoragePoolID: spID,
		NasServerID:   nasServer.ID,
	}

	// create the file system
	filesystem, err := system.CreateFileSystem(ctx, fs)
	fsID := filesystem.ID
	assert.Nil(t, err)
	assert.NotNil(t, fsID)

	snapResp, err := system.CreateFileSystemSnapshot(ctx, &types.CreateFileSystemSnapshotParam{
		Name: "test-snapshot",
	}, fsID)

	assert.NotNil(t, snapResp)
	assert.Nil(t, err)

	snap, err := system.GetFileSystemByIDName(ctx, snapResp.ID, "")
	assert.NotNil(t, snap)
	assert.Nil(t, err)

	snapList, err := system.GetFsSnapshotsByVolumeID(ctx, fsID)
	assert.NotNil(t, snapList)
	assert.Nil(t, err)

	err = system.DeleteFileSystem(ctx, snap.Name)
	assert.Nil(t, err)
	err = system.DeleteFileSystem(ctx, fs.Name)
	assert.Nil(t, err)
}

// TestGetFileSystemByIDName will return specific filesystem by name or ID
func TestGetFileSystemByIDName(t *testing.T) {
	system := getSystem()
	assert.NotNil(t, system)

	fsName := getFileSystemName(t)
	assert.NotZero(t, len(fsName))

	ctx := context.Background()

	filesystem, err := system.GetFileSystemByIDName(ctx, "", fsName)
	assert.Nil(t, err)
	assert.Equal(t, fsName, filesystem.Name)

	if filesystem != nil {
		fs, err := system.GetFileSystemByIDName(ctx, filesystem.ID, "")
		assert.Nil(t, err)
		assert.Equal(t, filesystem.ID, fs.ID)
	}

	if len(fsName) > 0 {
		fs, err := system.GetFileSystemByIDName(ctx, "", fsName)
		assert.Nil(t, err)
		assert.Equal(t, fsName, fs.Name)
	}
}

// TestGetFileSystemByNameIDInvalid attempts to get a file system  that does not exist
func TestGetFileSystemByNameIDInvalid(t *testing.T) {
	system := getSystem()
	assert.NotNil(t, system)

	ctx := context.Background()

	fs, err := system.GetFileSystemByIDName(ctx, invalidIdentifier, "")
	assert.NotNil(t, err)
	assert.Nil(t, fs)

	filesystem, err := system.GetFileSystemByIDName(ctx, "", invalidIdentifier)
	assert.NotNil(t, err)
	assert.Nil(t, filesystem)
}

// TestCreateDeleteFileSystem attempts to create then delete a file system
func TestCreateDeleteFileSystem(t *testing.T) {
	system := getSystem()
	assert.NotNil(t, system)

	fsName := fmt.Sprintf("%s-%s", "FS", testPrefix)

	// get protection domain
	pd := getProtectionDomain(t)
	assert.NotNil(t, pd)

	ctx := context.Background()

	// get storage pool
	pool := getStoragePool(t)
	assert.NotNil(t, pool)
	var spID string
	if pd != nil && pool != nil {
		sp, _ := pd.FindStoragePool(ctx, pool.StoragePool.ID, "", "")
		assert.NotNil(t, sp)
		spID = sp.ID
	}

	// get NAS server ID
	var nasServerName string
	if os.Getenv("GOSCALEIO_NASSERVER") != "" {
		nasServerName = os.Getenv("GOSCALEIO_NASSERVER")
	}
	nasServer, err := system.GetNASByIDName(ctx, "", nasServerName)
	assert.NotNil(t, nasServer)

	fs := &types.FsCreate{
		Name:          fsName,
		SizeTotal:     16106127360,
		StoragePoolID: spID,
		NasServerID:   nasServer.ID,
	}

	// create the file system
	filesystem, err := system.CreateFileSystem(ctx, fs)
	fsID := filesystem.ID
	assert.Nil(t, err)
	assert.NotNil(t, fsID)

	// try to create a file system that exists
	filesystem, err = system.CreateFileSystem(ctx, fs)
	assert.NotNil(t, err)

	// delete the file system
	err = system.DeleteFileSystem(ctx, fsName)
	assert.Nil(t, err)

	// try to delete non-existent file system
	// delete the file system
	err = system.DeleteFileSystem(ctx, fsName)
	assert.NotNil(t, err)
}

// TestModifyFileSystem attempts to create then delete a file system
func TestModifyFileSystem(t *testing.T) {
	system := getSystem()
	assert.NotNil(t, system)

	fsName := fmt.Sprintf("%s-%s", "FS", testPrefix)

	// get protection domain
	pd := getProtectionDomain(t)
	assert.NotNil(t, pd)

	ctx := context.Background()

	// get storage pool
	pool := getStoragePool(t)
	assert.NotNil(t, pool)
	var spID string
	if pd != nil && pool != nil {
		sp, _ := pd.FindStoragePool(ctx, pool.StoragePool.ID, "", "")
		assert.NotNil(t, sp)
		spID = sp.ID
	}

	// get NAS server ID
	var nasServerName string
	if os.Getenv("GOSCALEIO_NASSERVER") != "" {
		nasServerName = os.Getenv("GOSCALEIO_NASSERVER")
	}
	nasServer, err := system.GetNASByIDName(ctx, "", nasServerName)
	assert.NotNil(t, nasServer)

	fs := &types.FsCreate{
		Name:          fsName,
		SizeTotal:     10737418240,
		StoragePoolID: spID,
		NasServerID:   nasServer.ID,
	}

	// create the file system
	filesystem, err := system.CreateFileSystem(ctx, fs)
	fsID := filesystem.ID
	assert.Nil(t, err)
	assert.NotNil(t, fsID)

	// modify-expand filesystem
	modifyParam := &types.FSModify{
		Size:        16106127360,
		Description: "filesystem size modified",
	}
	err = system.ModifyFileSystem(ctx, modifyParam, fsID)
	assert.Nil(t, err)

	// negative case
	err = system.ModifyFileSystem(ctx, modifyParam, "")
	assert.NotNil(t, err)
}
