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

	types "github.com/dell/goscaleio/types/v1"
)

// StoragePool struct defines struct for StoragePool
type StoragePool struct {
	StoragePool *types.StoragePool
	client      *Client
}

// NewStoragePool returns a new StoragePool
func NewStoragePool(client *Client) *StoragePool {
	return &StoragePool{
		StoragePool: &types.StoragePool{},
		client:      client,
	}
}

// NewStoragePoolEx returns a new StoragePoolEx
func NewStoragePoolEx(client *Client, pool *types.StoragePool) *StoragePool {
	return &StoragePool{
		StoragePool: pool,
		client:      client,
	}
}

// CreateStoragePool creates a storage pool
func (pd *ProtectionDomain) CreateStoragePool(sp *types.StoragePoolParam) (string, error) {
	path := fmt.Sprintf("/api/types/StoragePool/instances")
	sp.ProtectionDomainID = pd.ProtectionDomain.ID
	spResponse := types.StoragePoolResp{}
	err := pd.client.getJSONWithRetry(
		http.MethodPost, path, sp, &spResponse)
	if err != nil {
		return "", err
	}
	return spResponse.ID, nil
}

// ModifyStoragePoolName Modifies Storagepool Name
func (pd *ProtectionDomain) ModifyStoragePoolName(ID, name string) (string, error) {

	storagePoolParam := &types.ModifyStoragePoolName{
		Name: name,
	}

	path := fmt.Sprintf("/api/instances/StoragePool::%v/action/setStoragePoolName", ID)

	spResp := types.StoragePoolResp{}
	err := pd.client.getJSONWithRetry(
		http.MethodPost, path, storagePoolParam, &spResp)
	if err != nil {
		return "", err
	}

	return spResp.ID, nil
}

// ModifyStoragePoolMedia Modifies Storagepool Media Type
func (pd *ProtectionDomain) ModifyStoragePoolMedia(ID, mediaType string) (string, error) {

	storagePool := &types.StoragePoolMediaType{
		MediaType: mediaType,
	}

	path := fmt.Sprintf("/api/instances/StoragePool::%v/action/setMediaType", ID)

	spResp := types.StoragePoolResp{}
	err := pd.client.getJSONWithRetry(
		http.MethodPost, path, storagePool, &spResp)
	if err != nil {
		return "", err
	}

	return spResp.ID, nil
}

// ModifyRMCache Sets Read RAM Cache
func (sp *StoragePool) ModifyRMCache(useRmcache string) error {

	link, err := GetLink(sp.StoragePool.Links, "self")
	if err != nil {
		return err
	}
	path := fmt.Sprintf("%v/action/setUseRmcache", link.HREF)
	payload := &types.StoragePoolUseRmCache{
		UseRmcache: useRmcache,
	}
	err = sp.client.getJSONWithRetry(
		http.MethodPost, path, payload, nil)
	return err
}

// EnableRFCache Enables RFCache
func (pd *ProtectionDomain) EnableRFCache(ID string) (string, error) {

	storagePoolParam := &types.StoragePoolUseRfCache{}

	path := fmt.Sprintf("/api/instances/StoragePool::%v/action/enableRfcache", ID)

	spResp := types.StoragePoolResp{}
	err := pd.client.getJSONWithRetry(
		http.MethodPost, path, storagePoolParam, &spResp)
	if err != nil {
		return "", err
	}

	return spResp.ID, nil
}

// DisableRFCache Disables RFCache
func (pd *ProtectionDomain) DisableRFCache(ID string) (string, error) {

	payload := &types.StoragePoolUseRfCache{}

	path := fmt.Sprintf("/api/instances/StoragePool::%v/action/disableRfcache", ID)

	spResp := types.StoragePoolResp{}
	err := pd.client.getJSONWithRetry(

		http.MethodPost, path, payload, &spResp)
	if err != nil {
		return "", err
	}

	return spResp.ID, nil
}

// DeleteStoragePool will delete a storage pool
func (pd *ProtectionDomain) DeleteStoragePool(name string) error {
	// get the storage pool name
	pool, err := pd.client.FindStoragePool("", name, "", "")
	if err != nil {
		return err
	}

	link, err := GetLink(pool.Links, "self")
	if err != nil {
		return err
	}

	storagePoolParam := &types.EmptyPayload{}

	path := fmt.Sprintf("%v/action/removeStoragePool", link.HREF)

	err = pd.client.getJSONWithRetry(
		http.MethodPost, path, storagePoolParam, nil)
	if err != nil {
		return err
	}

	return nil
}

// GetStoragePool returns a storage pool
func (pd *ProtectionDomain) GetStoragePool(
	storagepoolhref string) ([]*types.StoragePool, error) {

	var (
		err error
		sp  = &types.StoragePool{}
		sps []*types.StoragePool
	)

	if storagepoolhref == "" {
		var link *types.Link
		link, err := GetLink(pd.ProtectionDomain.Links,
			"/api/ProtectionDomain/relationship/StoragePool")
		if err != nil {
			return nil, err
		}
		err = pd.client.getJSONWithRetry(
			http.MethodGet, link.HREF, nil, &sps)
	} else {
		err = pd.client.getJSONWithRetry(
			http.MethodGet, storagepoolhref, nil, sp)
	}
	if err != nil {
		return nil, err
	}

	if storagepoolhref != "" {
		sps = append(sps, sp)
	}
	return sps, nil
}

// FindStoragePool returns a storagepool based on id or name
func (pd *ProtectionDomain) FindStoragePool(
	id, name, href string) (*types.StoragePool, error) {

	sps, err := pd.GetStoragePool(href)
	if err != nil {
		return nil, fmt.Errorf("Error getting protection domains %s", err)
	}

	for _, sp := range sps {
		if sp.ID == id || sp.Name == name || href != "" {
			return sp, nil
		}
	}

	return nil, errors.New("Couldn't find storage pool")

}

// GetStatistics returns statistics
func (sp *StoragePool) GetStatistics() (*types.Statistics, error) {

	link, err := GetLink(sp.StoragePool.Links,
		"/api/StoragePool/relationship/Statistics")
	if err != nil {
		return nil, err
	}

	stats := types.Statistics{}
	err = sp.client.getJSONWithRetry(
		http.MethodGet, link.HREF, nil, &stats)
	if err != nil {
		return nil, err
	}

	return &stats, nil
}

// GetSDSStoragePool return SDS instances associated with storage pool
func (sp *StoragePool) GetSDSStoragePool() ([]types.Sds, error) {
	link, err := GetLink(sp.StoragePool.Links,
		"/api/StoragePool/relationship/SpSds")
	if err != nil {
		return nil, err
	}

	sds := []types.Sds{}
	err = sp.client.getJSONWithRetry(
		http.MethodGet, link.HREF, nil, &sds)
	if err != nil {
		return nil, err
	}

	return sds, nil
}
