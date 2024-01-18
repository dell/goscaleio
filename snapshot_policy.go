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

package goscaleio

import (
	"fmt"
	"net/http"
	"time"

	types "github.com/dell/goscaleio/types/v1"
)

// CreateSnapshotPolicy creates a snapshot policy on the PowerFlex array
func (system *System) CreateSnapshotPolicy(snapPolicy *types.SnapshotPolicyCreateParam) (string, error) {
	defer TimeSpent("crate snapshot policy", time.Now())

	path := fmt.Sprintf("/api/types/SnapshotPolicy/instances")
	snapResp := types.SnapShotPolicyCreateResp{}
	 err := system.client.getJSONWithRetry(
		http.MethodPost, path, snapPolicy, &snapResp)
	if err != nil {
		return "", err
	}
	return snapResp.ID, nil
}

// RemoveSnapshotPolicy removes a snapshot policy from the PowerFlex array
func (system *System) RemoveSnapshotPolicy(id string) error {
	path := fmt.Sprintf("/api/instances/SnapshotPolicy::%v/action/removeSnapshotPolicy", id)
	removeParam := &types.EmptyPayload{}
	 err := system.client.getJSONWithRetry(
		http.MethodPost, path, removeParam, nil)
	if err != nil {
		return err
	}
	return nil
}

// RenameSnapshotPolicy renames a snapshot policy 
func (system *System) RenameSnapshotPolicy(id, name string) error {
	path := fmt.Sprintf("/api/instances/SnapshotPolicy::%v/action/renameSnapshotPolicy", id)
	renameSnap := &types.SnapshotPolicyRenameParam{
		NewName: name,
	}
	 err := system.client.getJSONWithRetry(
		http.MethodPost, path, renameSnap, nil)
	if err != nil {
		return err
	}
	return nil
}

// ModifySnapshotPolicy modifies a snapshot policy
func (system *System) ModifySnapshotPolicy(modifysnapPolicy *types.SnapshotPolicyModifyParam, id string) error {
	path := fmt.Sprintf("/api/instances/SnapshotPolicy::%v/action/modifySnapshotPolicy", id)
	 err := system.client.getJSONWithRetry(
		http.MethodPost, path, modifysnapPolicy, nil)
	if err != nil {
		return err
	}
	return nil
}

// AssignVolumeToSnapshotPolicy assigns volume to a snapshot policy
func (system *System) AssignVolumeToSnapshotPolicy(assignVoltoSnap *types.AssignVolumeToSnapshotPolicyParam, id string) error {
	path := fmt.Sprintf("/api/instances/SnapshotPolicy::%v/action/addSourceVolumeToSnapshotPolicy", id)
	 err := system.client.getJSONWithRetry(
		http.MethodPost, path, assignVoltoSnap, nil)
	if err != nil {
		return err
	}
	return nil
}

// UnassignVolumeFromSnapshotPolicy unassigns volume from a snapshot policy
func (system *System) UnassignVolumeFromSnapshotPolicy(UnassignVolFromSnap *types.AssignVolumeToSnapshotPolicyParam, id string) error {
	path := fmt.Sprintf("/api/instances/SnapshotPolicy::%v/action/removeSourceVolumeFromSnapshotPolicy", id)
	 err := system.client.getJSONWithRetry(
		http.MethodPost, path, UnassignVolFromSnap, nil)
	if err != nil {
		return err
	}
	return nil
}

// PauseSnapshotPolicy pause a snapshot policy
func (system *System) PauseSnapshotPolicy(id string) error {
	path := fmt.Sprintf("/api/instances/SnapshotPolicy::%v/action/pauseSnapshotPolicy", id)
	pauseParam := &types.EmptyPayload{}
	 err := system.client.getJSONWithRetry(
		http.MethodPost, path, pauseParam, nil)
	if err != nil {
		return err
	}
	return nil
}

// ResumeSnapshotPolicy resume a snapshot policy which was paused
func (system *System) ResumeSnapshotPolicy(id string) error {
	path := fmt.Sprintf("/api/instances/SnapshotPolicy::%v/action/resumeSnapshotPolicy", id)
	resumeParam := &types.EmptyPayload{}
	 err := system.client.getJSONWithRetry(
		http.MethodPost, path, resumeParam, nil)
	if err != nil {
		return err
	}
	return nil
}