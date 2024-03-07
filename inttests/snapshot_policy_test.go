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
	"fmt"
	"testing"
	"time"

	types "github.com/dell/goscaleio/types/v1"
	"github.com/stretchr/testify/assert"
)

func TestCreateModifyDeleteSnapshotPolicy(t *testing.T) {
	system := getSystem()
	assert.NotNil(t, system)
	spName := fmt.Sprintf("%s-%s", testPrefix, "SnapPolicy")

	snap := &types.SnapshotPolicyCreateParam{
		Name:                             spName,
		AutoSnapshotCreationCadenceInMin: "5",
		NumOfRetainedSnapshotsPerLevel:   []string{"1"},
		SnapshotAccessMode:               "ReadOnly",
		Paused:                           "true",
	}

	// create the snapshot policy
	snapID, err := system.CreateSnapshotPolicy(snap)
	assert.Nil(t, err)
	assert.NotNil(t, snapID)
	time.Sleep(5 * time.Second)

	// create a snapsshot policy that exists
	_, err2 := system.CreateSnapshotPolicy(snap)
	assert.NotNil(t, err2)

	// modify snapshot policy name
	err = system.RenameSnapshotPolicy(snapID, "SnapshotPolicyRenamed")
	assert.Nil(t, err)

	// modify other parameters of snapshot policy
	snapModify := &types.SnapshotPolicyModifyParam{
		AutoSnapshotCreationCadenceInMin: "6",
		NumOfRetainedSnapshotsPerLevel:   []string{"2", "6"},
	}
	err = system.ModifySnapshotPolicy(snapModify, snapID)
	assert.Nil(t, err)

	volID, err := createVolume(t, "")
	assignVolume := &types.AssignVolumeToSnapshotPolicyParam{
		SourceVolumeID: volID,
	}

	// Assign and unassign volume to Snapshot Policy
	err = system.AssignVolumeToSnapshotPolicy(assignVolume, snapID)
	assert.Nil(t, err)

	vol, err2 := system.GetSourceVolume(snapID)
	assert.Nil(t, err2)
	assert.NotNil(t, vol)

	assignVolume = &types.AssignVolumeToSnapshotPolicyParam{
		SourceVolumeID: "Invalid",
	}
	err = system.AssignVolumeToSnapshotPolicy(assignVolume, snapID)
	assert.NotNil(t, err)

	unassignVolume := &types.AssignVolumeToSnapshotPolicyParam{
		SourceVolumeID:            volID,
		AutoSnapshotRemovalAction: "Remove",
	}
	err = system.UnassignVolumeFromSnapshotPolicy(unassignVolume, snapID)
	assert.Nil(t, err)

	unassignVolume = &types.AssignVolumeToSnapshotPolicyParam{
		SourceVolumeID:            volID,
		AutoSnapshotRemovalAction: "Invalid",
	}
	err = system.UnassignVolumeFromSnapshotPolicy(unassignVolume, snapID)
	assert.NotNil(t, err)

	// Resume and Pause the SnapshotPolicy
	err = system.ResumeSnapshotPolicy(snapID)
	assert.Nil(t, err)

	err = system.ResumeSnapshotPolicy("Invalid")
	assert.NotNil(t, err)

	err = system.PauseSnapshotPolicy(snapID)
	assert.Nil(t, err)

	err = system.PauseSnapshotPolicy("Invalid")
	assert.NotNil(t, err)

	// delete the snapshot policy
	err = system.RemoveSnapshotPolicy(snapID)
	assert.Nil(t, err)

	// try to delete non-existent snapsot policy
	err = system.RemoveSnapshotPolicy(invalidIdentifier)
	assert.NotNil(t, err)
}
