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
	"fmt"
	"net/http"

	types "github.com/dell/goscaleio/types/v1"
)

func (c *Client) GetSSOUser(userID string) (*types.SSOUserDetails, error) {
	path := fmt.Sprintf("/rest/v1/users/%s", userID)
	user := &types.SSOUserDetails{}
	err := c.getJSONWithRetry(http.MethodGet, path, nil, &user)

	if err != nil {
		return nil, err
	}

	return user, nil
}

func (c *Client) CreateSSOUser(userParam *types.SSOUserCreateParam) (*types.SSOUserDetails, error) {
	userResp := &types.SSOUserDetails{}
	err := c.getJSONWithRetry(http.MethodPost, "/rest/v1/users", userParam, &userResp)
	if err != nil {
		return nil, err
	}
	return userResp, nil
}

func (c *Client) ModifySSOUser(userID string, userParam *types.SSOUserModifyParam) (*types.SSOUserDetails, error) {
	path := fmt.Sprintf("/rest/v1/users/%s", userID)
	err := c.getJSONWithRetry(http.MethodPatch, path, userParam, nil)
	if err != nil {
		return nil, err
	}
	return c.GetSSOUser(userID)
}

func (c *Client) ResetSSOUserPassword(userID string, userParam *types.SSOUserModifyParam) error {
	path := fmt.Sprintf("/rest/v1/users/%s/reset-password", userID)
	err := c.getJSONWithRetry(http.MethodPost, path, userParam, nil)
	if err != nil {
		return err
	}
	return nil
}

func (c *Client) DeleteSSOUser(userID string) error {
	path := fmt.Sprintf("/rest/v1/users/%s", userID)
	err := c.getJSONWithRetry(http.MethodDelete, path, nil, nil)
	if err != nil {
		return err
	}
	return nil
}
