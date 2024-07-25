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

package inttests

import (
	"testing"
	"time"

	types "github.com/dell/goscaleio/types/v1"
	"github.com/stretchr/testify/assert"
)

// TestOSRepositoryGetAll tests get All OS Repositories
func TestOSRepositoryGetAll(t *testing.T) {
	system := getSystem()
	osRepositories, err := system.GetAllOSRepositories()
	assert.Nil(t, err)
	assert.NotNil(t, osRepositories)
}

// TestOSRepositoryGetByID tests Get OS Repository by Id
func TestOSRepositoryGetByID(t *testing.T) {
	system := getSystem()
	osRepository, err := system.GetOSRepositoryByID("8aaa80458fca6913018fce6449f50e81")
	assert.Nil(t, err)
	assert.NotNil(t, osRepository)
}

// TestOSRepositoryGetByIDFail tests the negative case for Get OS Repository by Id
func TestOSRepositoryGetByIDFail(t *testing.T) {
	system := getSystem()
	_, err := system.GetOSRepositoryByID("Invalid")
	assert.NotNil(t, err)
}

// TestGetOSRepositoryByIDFail tests the negative case for Get OS Repository by Id
func TestOSRepositoryCreateFail(t *testing.T) {
	system := getSystem()
	_, err := system.CreateOSRepository(nil)
	assert.NotNil(t, err)
}

// TestOSRepositoryDeleteFail tests the negative case for Delete OS Repository
func TestOSRepositoryDeleteFail(t *testing.T) {
	system := getSystem()
	// Delete
	err := system.RemoveOSRepository("invalid")
	assert.NotNil(t, err)
}

// TestOSRepositoryCreateAndDelete tests create and delete operations for OS Repository
func TestOSRepositoryCreateAndDelete(t *testing.T) {
	system := getSystem()
	createOSRepository := &types.OSRepository{
		Name:       "Test-OS-Repository",
		RepoType:   "ISO",
		SourcePath: "https://100.65.27.72/artifactory/vxfm-yum-release/pfmp20/RCM/Denver/RCMs/esxi/ESXi-8.0.0-20513097-3.8.0.0_Dell.iso",
		ImageType:  "vmware_esxi",
	}
	// Create
	osRepository, err := system.CreateOSRepository(createOSRepository)
	assert.Nil(t, err)
	assert.NotNil(t, osRepository)
	// We will wait for Repository to be unpacked and created
	time.Sleep(240 * time.Second)

	var repoID string
	osRepositories, err := system.GetAllOSRepositories()
	assert.Nil(t, err)
	assert.NotNil(t, osRepositories)
	for _, repo := range osRepositories {
		if repo.Name == "Test-OS-Repository" {
			repoID = repo.ID
			break
		}
	}

	// Delete
	err = system.RemoveOSRepository(repoID)
	assert.Nil(t, err)
}
