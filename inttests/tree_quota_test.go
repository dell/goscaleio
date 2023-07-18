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
	"fmt"
	"os"
	"testing"

	types "github.com/dell/goscaleio/types/v1"
	"github.com/stretchr/testify/assert"
)

// TestTreeQuotaByID will return specific specific Tree Quota by ID
func TestTreeQuotaByID(t *testing.T) {
	system := getSystem()
	assert.NotNil(t, system)

	var quotaid string

	if os.Getenv("GOSCALEIO_TREEQUOTAID") != "" {
		quotaid = os.Getenv("GOSCALEIO_TREEQUOTAID")
	}
	fmt.Println("quotaid", quotaid)
	quota, err := system.GetTreeQuotaByID(quotaid)
	assert.Nil(t, err)
	assert.Equal(t, quotaid, quota.ID)

	if quota != nil {
		treequota, err := system.GetTreeQuotaByID(quota.ID)
		assert.Nil(t, err)
		assert.Equal(t, treequota.ID, quota.ID)
	}

}

// TestCreateModifyDeleteTreeQuota will create , modify and delete a tree quota
func TestCreateModifyDeleteTreeQuota(t *testing.T) {
	system := getSystem()
	assert.NotNil(t, system)

	treeQuotaModify := "Tree Quota modify testing"
	var filesystemname string
	if os.Getenv("GOSCALEIO_FILESYSTEM") != "" {
		filesystemname = os.Getenv("GOSCALEIO_FILESYSTEM")
		fmt.Println("name...", filesystemname)
	}
	filesystem, err := system.GetFileSystemByIDName("", filesystemname)
	treequota := &types.TreeQuotaCreate{
		FileSystemID: filesystem.ID,
		Path:         "/" + "fs123",
	}

	err = system.ModifyFileSystem(&types.FSModify{
		IsQuotaEnabled: true}, filesystem.ID)

	fmt.Println("iddd..", filesystem.ID)

	//create tree quota
	quota, err := system.CreateTreeQuota(treequota)
	quotaid := quota.ID
	fmt.Println("quotaid", quotaid)
	assert.Nil(t, err)
	assert.NotNil(t, quotaid)

	// try to create existing tree quota
	quota, err = system.CreateTreeQuota(treequota)
	assert.NotNil(t, err)

	//Modify Tree Quota
	quotaModify := &types.TreeQuotaModify{
		Description: treeQuotaModify,
		SoftLimit:   900,
	}

	err = system.ModifyTreeQuota(quotaModify, quotaid)
	assert.Nil(t, err)
	fmt.Println("errmod", err)

	// negative case
	err = system.ModifyTreeQuota(quotaModify, "")
	assert.NotNil(t, err)

	//Delete tree Quota
	err = system.DeleteTreeQuota(quotaid)
	assert.Nil(t, err)
}
