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
func (pd *ProtectionDomain) CreateStoragePool(name string, mediaType string) (string, error) {

	if mediaType == "" {
		mediaType = "HDD"
	}
	storagePoolParam := &types.StoragePoolParam{
		Name:               name,
		ProtectionDomainID: pd.ProtectionDomain.ID,
		MediaType:          mediaType,
	}

	path := fmt.Sprintf("/api/types/StoragePool/instances")

	sp := types.StoragePoolResp{}
	err := pd.client.getJSONWithRetry(
		http.MethodPost, path, storagePoolParam, &sp)
	if err != nil {
		return "", err
	}

	return sp.ID, nil
}

// DeleteStoragePool will delete a storage pool
func (pd *ProtectionDomain) DeleteStoragePool(name string) error {
	// get the storage pool name
	pool, err := pd.client.FindStoragePool("", name, "")
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
