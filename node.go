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
	"errors"
	"fmt"
	"net/http"
	"time"

	types "github.com/dell/goscaleio/types/v1"
)

// GetNodeByID gets the node details based on ID
func (c *Client) GetNodeByID(id string) (*types.NodeDetails, error) {
	defer TimeSpent("GetNodeByID", time.Now())

	path := fmt.Sprintf("/Api/V1/ManagedDevice/%v", id)

	var node types.NodeDetails
	err := c.getJSONWithRetry(http.MethodGet, path, nil, &node)
	if err != nil {
		return nil, err
	}
	return &node, nil
}

// GetAllNodes gets all the node details
func (c *Client) GetAllNodes() ([]types.NodeDetails, error) {
	defer TimeSpent("GetNodeByID", time.Now())

	path := fmt.Sprintf("/Api/V1/ManagedDevice")

	var nodes []types.NodeDetails
	err := c.getJSONWithRetry(http.MethodGet, path, nil, &nodes)
	if err != nil {
		return nil, err
	}
	return nodes, nil
}

// GetNodeByFilters gets the node details based on the provided filter
func (c *Client) GetNodeByFilters(key string, value string) ([]types.NodeDetails, error) {
	defer TimeSpent("GetNodeByFilters", time.Now())

	path := fmt.Sprintf("/Api/V1/ManagedDevice?filter=eq,%v,%v", key, value)

	var nodes []types.NodeDetails
	err := c.getJSONWithRetry(http.MethodGet, path, nil, &nodes)
	if err != nil {
		return nil, err
	}

	if len(nodes) == 0 {
		return nil, errors.New("Couldn't find nodes with the given filter")
	}
	return nodes, nil
}

// GetNodePoolByID gets the nodepool details based on ID
func (c *Client) GetNodePoolByID(id int) (*types.NodePoolDetails, error) {
	defer TimeSpent("GetNodePoolByID", time.Now())

	path := fmt.Sprintf("/Api/V1/nodepool/%v", id)

	var nodePool types.NodePoolDetails
	err := c.getJSONWithRetry(http.MethodGet, path, nil, &nodePool)
	if err != nil {
		return nil, err
	}

	return &nodePool, nil
}

// GetNodePoolByName gets the nodepool details based on name
func (c *Client) GetNodePoolByName(name string) (*types.NodePoolDetails, error) {
	defer TimeSpent("GetNodePoolByName", time.Now())

	nodePools, err := c.GetAllNodePools()
	if err != nil {
		return nil, err
	}

	for _, nodePool := range nodePools.NodePoolDetails {
		if nodePool.GroupName == name {
			return c.GetNodePoolByID(nodePool.GroupSeqID)
		}
	}
	return nil, errors.New("no node pool found with name " + name)
}

// GetAllNodePools gets all the nodepool details
func (c *Client) GetAllNodePools() (*types.NodePoolDetailsFilter, error) {
	defer TimeSpent("GetAllNodePools", time.Now())

	path := fmt.Sprintf("/Api/V1/nodepool")

	var nodePools types.NodePoolDetailsFilter
	err := c.getJSONWithRetry(http.MethodGet, path, nil, &nodePools)
	if err != nil {
		return nil, err
	}

	return &nodePools, nil
}
