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
