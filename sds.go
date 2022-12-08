// Copyright © 2019 - 2022 Dell Inc. or its subsidiaries. All Rights Reserved.
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
	"errors"
	"fmt"
	"net/http"
	"reflect"
	"time"

	types "github.com/dell/goscaleio/types/v1"
)

// Sds defines struct for Sds
type Sds struct {
	Sds    *types.Sds
	client *Client
}

// NewSds returns a new Sds
func NewSds(client *Client) *Sds {
	return &Sds{
		Sds:    &types.Sds{},
		client: client,
	}
}

// NewSdsEx returns a new SdsEx
func NewSdsEx(client *Client, sds *types.Sds) *Sds {
	return &Sds{
		Sds:    sds,
		client: client,
	}
}

// CreateSds creates a new Sds
func (pd *ProtectionDomain) CreateSds(
	name string, ipList []string) (string, error) {
	defer TimeSpent("CreateSds", time.Now())

	sdsParam := &types.SdsParam{
		Name:               name,
		ProtectionDomainID: pd.ProtectionDomain.ID,
	}

	if len(ipList) == 0 {
		return "", fmt.Errorf("Must provide at least 1 SDS IP")
	} else if len(ipList) == 1 {
		sdsIP := types.SdsIP{IP: ipList[0], Role: "all"}
		sdsIPList := &types.SdsIPList{SdsIP: sdsIP}
		sdsParam.IPList = append(sdsParam.IPList, sdsIPList)
	} else if len(ipList) >= 2 {
		sdsIP1 := types.SdsIP{IP: ipList[0], Role: "sdcOnly"}
		sdsIP2 := types.SdsIP{IP: ipList[1], Role: "sdsOnly"}
		sdsIPList1 := &types.SdsIPList{SdsIP: sdsIP1}
		sdsIPList2 := &types.SdsIPList{SdsIP: sdsIP2}
		sdsParam.IPList = append(sdsParam.IPList, sdsIPList1)
		sdsParam.IPList = append(sdsParam.IPList, sdsIPList2)
	}

	path := fmt.Sprintf("/api/types/Sds/instances")

	sds := types.SdsResp{}
	err := pd.client.getJSONWithRetry(
		http.MethodPost, path, sdsParam, &sds)
	if err != nil {
		return "", err
	}

	return sds.ID, nil
}

// GetSds returns a Sds
func (pd *ProtectionDomain) GetSds() ([]types.Sds, error) {
	defer TimeSpent("GetSds", time.Now())
	path := fmt.Sprintf("/api/instances/ProtectionDomain::%v/relationships/Sds",
		pd.ProtectionDomain.ID)

	var sdss []types.Sds
	err := pd.client.getJSONWithRetry(
		http.MethodGet, path, nil, &sdss)
	if err != nil {
		return nil, err
	}

	return sdss, nil
}

// FindSds returns a Sds
func (pd *ProtectionDomain) FindSds(
	field, value string) (*types.Sds, error) {
	defer TimeSpent("FindSds", time.Now())

	sdss, err := pd.GetSds()
	if err != nil {
		return nil, err
	}

	for _, sds := range sdss {
		valueOf := reflect.ValueOf(sds)
		switch {
		case reflect.Indirect(valueOf).FieldByName(field).String() == value:
			return &sds, nil
		}
	}

	return nil, errors.New("Couldn't find SDS")
}

// DeleteSds deleted a Sds against Id
func (pd *ProtectionDomain) DeleteSds(id string) error {
	defer TimeSpent("DeleteSds", time.Now())

	path := fmt.Sprintf("/api/instances/Sds::%v/action/removeSds", id)

	err := pd.client.getJSONWithRetry(http.MethodPost, path, nil, nil)
	if err != nil {
		return err
	}

	return nil
}

// SetSDSIPRole sets IP and Role of SDS
func (pd *ProtectionDomain) SetSDSIPRole(id, ip, role string) error {
	defer TimeSpent("SetSDSIPRole", time.Now())

	sdsParam := &types.SdsIPRole{
		SdsIPToSet: ip,
		NewRole:    role,
	}

	path := fmt.Sprintf("/api/instances/Sds::%v/action/setSdsIpRole", id)

	err := pd.client.getJSONWithRetry(http.MethodPost, path, sdsParam, nil)
	if err != nil {
		return err
	}

	return nil
}

// RemoveSDSIP sets IP and Role of SDS
func (pd *ProtectionDomain) RemoveSDSIP(id, ip string) error {
	defer TimeSpent("RemoveSDSIP", time.Now())

	sdsParam := &types.SdsIP{
		IP: ip,
	}

	path := fmt.Sprintf("/api/instances/Sds::%v/action/removeSdsIp", id)

	err := pd.client.getJSONWithRetry(http.MethodPost, path, sdsParam, nil)
	if err != nil {
		return err
	}

	return nil
}

// SetSdsName sets sds name
func (pd *ProtectionDomain) SetSdsName(id, name string) error {
	defer TimeSpent("SetSdsName", time.Now())

	sdsParam := &types.SdsName{
		Name: name,
	}

	path := fmt.Sprintf("/api/instances/Sds::%v/action/setSdsName", id)

	err := pd.client.getJSONWithRetry(http.MethodPost, path, sdsParam, nil)
	if err != nil {
		return err
	}

	return nil
}

// SetSdsPort sets sds name
func (pd *ProtectionDomain) SetSdsPort(id, port string) error {
	defer TimeSpent("SetSdsPort", time.Now())

	sdsParam := &types.SdsPort{
		SdsPort: port,
	}

	path := fmt.Sprintf("/api/instances/Sds::%v/action/setSdsPort", id)

	err := pd.client.getJSONWithRetry(http.MethodPost, path, sdsParam, nil)
	if err != nil {
		return err
	}

	return nil
}
