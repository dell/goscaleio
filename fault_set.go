// Copyright © 2023 Dell Inc. or its subsidiaries. All Rights Reserved.
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
	types "github.com/dell/goscaleio/types/v1"
	"net/http"
)

// CreateFaultSet creates a fault set
func (pd *ProtectionDomain) CreateFaultSet(fs *types.FaultSetParam) (string, error) {
	path := fmt.Sprintf("/api/types/FaultSet/instances")
	fs.ProtectionDomainID = pd.ProtectionDomain.ID
	fsResp := types.FaultSetResp{}
	err := pd.client.getJSONWithRetry(
		http.MethodPost, path, fs, &fsResp)
	if err != nil {
		return "", err
	}
	return fsResp.ID, nil
}

// DeleteFaultSet will delete a fault set
func (pd *ProtectionDomain) DeleteFaultSet(id string) error {
	path := fmt.Sprintf("/api/instances/FaultSet::%v/action/removeFaultSet", id)
	fsParam := &types.EmptyPayload{}
	err := pd.client.getJSONWithRetry(
		http.MethodPost, path, fsParam, nil)
	if err != nil {
		return err
	}
	return nil
}

// ModifyFaultSetName will modify the name of the fault set
func (pd *ProtectionDomain) ModifyFaultSetName(id, name string) error {
	fs := &types.FaultSetRename{}
	fs.NewName = name
	path := fmt.Sprintf("/api/instances/FaultSet::%v/action/setFaultSetName", id)

	err := pd.client.getJSONWithRetry(
		http.MethodPost, path, fs, nil)
	if err != nil {
		return err
	}
	return nil
}

// ModifyFaultSetPerfProfile will modify the performance profile of the fault set
func (pd *ProtectionDomain) ModifyFaultSetPerfProfile(id, perfProfile string) error {
	pp := &types.ChangeSdcPerfProfile{}
	pp.PerfProfile = perfProfile
	path := fmt.Sprintf("/api/instances/FaultSet::%v/action/setSdsPerformanceParameters", id)

	err := pd.client.getJSONWithRetry(
		http.MethodPost, path, pp, nil)
	if err != nil {
		return err
	}
	return nil
}

// GetFaultSetByID will read the fault set using the ID.
func (system *System) GetFaultSetByID(id string) (*types.FaultSet, error) {
	fs := &types.FaultSet{}
	path := fmt.Sprintf("/api/instances/FaultSet::%v", id)

	err := system.client.getJSONWithRetry(
		http.MethodGet, path, nil, fs)
	if err != nil {
		return nil, err
	}
	return fs, nil
}
