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

func (c *Client) GetNodePoolByID(id int) (*types.NodePoolDetails, error) {
	defer TimeSpent("GetNodeByID", time.Now())

	path := fmt.Sprintf("/Api/V1/nodepool/%v", id)

	var nodePool types.NodePoolDetails
	err := c.getJSONWithRetry(http.MethodGet, path, nil, &nodePool)
	if err != nil {
		return nil, err
	}
	return &nodePool, nil
}
