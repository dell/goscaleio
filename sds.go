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
	"strconv"
	"time"

	types "github.com/dell/goscaleio/types/v1"
)

// SDS IPs can have 3 roles.
const (
	RoleSdsOnly = "sdsOnly"
	RoleSdcOnly = "sdcOnly"
	RoleAll     = "all"
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

// CreateSds creates a new Sds with automatically assigned roles to IPs
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
		sdsIP := types.SdsIP{IP: ipList[0], Role: RoleAll}
		sdsIPList := &types.SdsIPList{SdsIP: sdsIP}
		sdsParam.IPList = append(sdsParam.IPList, sdsIPList)
	} else if len(ipList) == 2 {
		sdsIP1 := types.SdsIP{IP: ipList[0], Role: RoleSdcOnly}
		sdsIP2 := types.SdsIP{IP: ipList[1], Role: RoleSdsOnly}
		sdsIPList1 := &types.SdsIPList{SdsIP: sdsIP1}
		sdsIPList2 := &types.SdsIPList{SdsIP: sdsIP2}
		sdsParam.IPList = append(sdsParam.IPList, sdsIPList1)
		sdsParam.IPList = append(sdsParam.IPList, sdsIPList2)
	} else {
		return "", fmt.Errorf("Must explicitly provide IP role for more than 2 SDS IPs")
	}

	return pd.createSds(sdsParam)
}

func getNonZeroIntType(i int) string {
	if i == 0 {
		return ""
	}
	return strconv.Itoa(i)
}

// CreateSdsWithParams creates a new Sds with user defined SdsParam struct
func (pd *ProtectionDomain) CreateSdsWithParams(sds *types.Sds) (string, error) {
	defer TimeSpent("CreateSdsWithParams", time.Now())

	sdsParam := &types.SdsParam{
		Name:               sds.Name,
		ProtectionDomainID: pd.ProtectionDomain.ID,
		Port:               getNonZeroIntType(sds.Port),
		RmcacheEnabled:     types.GetBoolType(sds.RmcacheEnabled),
		RmcacheSizeInKb:    getNonZeroIntType(sds.RmcacheSizeInKb),
		FaultSetID:         sds.FaultSetID,
		NumOfIoBuffers:     getNonZeroIntType(sds.NumOfIoBuffers),
		DrlMode:            sds.DrlMode,
		IPList:             make([]*types.SdsIPList, 0),
	}

	ipList := sds.IPList

	if len(ipList) == 0 {
		return "", fmt.Errorf("Must provide at least 1 SDS IP")
	} else if len(ipList) == 1 {
		if ipList[0].Role != RoleAll {
			return "", fmt.Errorf("The only IP assigned to an SDS must be assigned \"%s\" role", RoleAll)
		}
		sdsParam.IPList = append(sdsParam.IPList, &types.SdsIPList{SdsIP: *ipList[0]})
	} else if len(ipList) >= 2 {
		nSdsOnly, nSdcOnly := 0, 0
		for i, ip := range ipList {
			if ip.Role == RoleAll || ip.Role == RoleSdcOnly {
				nSdcOnly++
			}
			if ip.Role == RoleAll || ip.Role == RoleSdsOnly {
				nSdsOnly++
			}
			sdsParam.IPList = append(sdsParam.IPList, &types.SdsIPList{SdsIP: *ipList[i]})
		}
		if nSdsOnly < 1 {
			return "", fmt.Errorf("At least one IP must be assigned %s or %s role", RoleSdsOnly, RoleAll)
		}
		if nSdcOnly < 1 {
			return "", fmt.Errorf("At least one IP must be assigned %s or %s role", RoleSdcOnly, RoleAll)
		}
	}

	return pd.createSds(sdsParam)
}

func (pd *ProtectionDomain) createSds(sdsParam *types.SdsParam) (string, error) {
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

// DeleteSds deletes a Sds against Id
func (pd *ProtectionDomain) DeleteSds(id string) error {
	defer TimeSpent("DeleteSds", time.Now())

	path := fmt.Sprintf("/api/instances/Sds::%v/action/removeSds", id)

	sdsParam := &types.EmptyPayload{}
	err := pd.client.getJSONWithRetry(http.MethodPost, path, sdsParam, nil)
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
