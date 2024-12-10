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

package goscaleio

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	types "github.com/dell/goscaleio/types/v1"
)

// GetTreeQuota gets list of tree Quota
func (s *System) GetTreeQuota(ctx context.Context) (treeQuotaList []types.TreeQuota, err error) {
	defer TimeSpent("GetTreeQuota", time.Now())
	path := fmt.Sprintf("/rest/v1/file-tree-quotas?select=*")

	err = s.client.getJSONWithRetry(ctx,
		http.MethodGet, path, nil, &treeQuotaList)
	if err != nil {
		return nil, err
	}

	return treeQuotaList, nil
}

// GetTreeQuotaByID gets a specific tree quota by ID
func (s *System) GetTreeQuotaByID(ctx context.Context, id string) (treeQuota *types.TreeQuota, err error) {
	defer TimeSpent("GetTreeQuota", time.Now())
	path := fmt.Sprintf("/rest/v1/file-tree-quotas/%s?select=*", id)

	err = s.client.getJSONWithRetry(ctx,
		http.MethodGet, path, nil, &treeQuota)
	if err != nil {
		return nil, err
	}

	return treeQuota, nil
}

// CreateTreeQuota create an tree quota for a File System.
func (s *System) CreateTreeQuota(ctx context.Context, createParams *types.TreeQuotaCreate) (resp *types.TreeQuotaCreateResponse, err error) {
	path := fmt.Sprintf("/rest/v1/file-tree-quotas")

	var body *types.TreeQuotaCreate = createParams
	err = s.client.getJSONWithRetry(ctx, http.MethodPost, path, body, &resp)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

// ModifyTreeQuota modifies a tree quota
func (s *System) ModifyTreeQuota(ctx context.Context, ModifyParams *types.TreeQuotaModify, id string) (err error) {
	path := fmt.Sprintf("/rest/v1/file-tree-quotas/%s", id)

	var body *types.TreeQuotaModify = ModifyParams
	err = s.client.getJSONWithRetry(ctx, http.MethodPatch, path, body, nil)
	if err != nil {
		return err
	}

	return nil
}

// DeleteTreeQuota delete a tree quota by ID
func (s *System) DeleteTreeQuota(ctx context.Context, id string) error {
	defer TimeSpent("DeleteTreeQuota", time.Now())
	path := fmt.Sprintf("/rest/v1/file-tree-quotas/%s", id)

	err := s.client.getJSONWithRetry(ctx,
		http.MethodDelete, path, nil, nil)
	if err != nil {
		return err
	}

	return nil
}

// GetTreeQuotaByFSID gets a specific tree quota by filesystem ID
func (s *System) GetTreeQuotaByFSID(ctx context.Context, id string) (*types.TreeQuota, error) {
	defer TimeSpent("GetTreeQuotaByFSID", time.Now())
	treeQuotaList, err := s.GetTreeQuota(ctx)
	if err != nil {
		return nil, err
	}
	for _, treeQuota := range treeQuotaList {
		if treeQuota.FileSysytemID == id {
			return &treeQuota, nil
		}
	}
	return nil, errors.New("couldn't find tree quota by filesystem ID")
}
