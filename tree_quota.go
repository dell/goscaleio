package goscaleio

import (
	"fmt"
	"net/http"
	"time"

	types "github.com/dell/goscaleio/types/v1"
)

func (s *System) GetTreeQuota() (treeQuotaList []types.TreeQuota, err error) {
	defer TimeSpent("GetTreeQuota", time.Now())
	path := fmt.Sprintf("/rest/v1/file-tree-quotas?select=*")

	err = s.client.getJSONWithRetry(
		http.MethodGet, path, nil, &treeQuotaList)
	if err != nil {
		return nil, err
	}

	return treeQuotaList, nil
}

func (s *System) GetTreeQuotaByID(id string) (treeQuota *types.TreeQuota, err error) {
	defer TimeSpent("GetTreeQuota", time.Now())
	path := fmt.Sprintf("/rest/v1/file-tree-quotas/%s?select=*", id)

	err = s.client.getJSONWithRetry(
		http.MethodGet, path, nil, &treeQuota)
	if err != nil {
		return nil, err
	}

	return treeQuota, nil
}

// CreateNFSExport create an NFS Export for a File System.
func (s *System) CreateTreeQuota(createParams *types.TreeQuotaCreate) (resp *types.TreeQuotaCreateResponse, err error) {
	path := fmt.Sprintf("/rest/v1/file-tree-quotas")

	var body *types.TreeQuotaCreate = createParams
	err = s.client.getJSONWithRetry(http.MethodPost, path, body, &resp)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (s *System) ModifyTreeQuota(ModifyParams *types.TreeQuotaModify, id string) (err error) {
	path := fmt.Sprintf("/rest/v1/file-tree-quotas/%s", id)

	var body *types.TreeQuotaModify = ModifyParams
	err = s.client.getJSONWithRetry(http.MethodPatch, path, body, nil)
	if err != nil {
		return err
	}

	return nil
}

func (s *System) DeleteTreeQuota(id string) error {
	defer TimeSpent("DeleteTreeQuota", time.Now())
	path := fmt.Sprintf("/rest/v1/file-tree-quotas/%s", id)

	err := s.client.getJSONWithRetry(
		http.MethodDelete, path, nil, nil)
	if err != nil {
		return err
	}

	return nil
}
