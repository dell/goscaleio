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
	"fmt"
	"net/http"
	"time"

	types "github.com/dell/goscaleio/types/v1"
)

// GetUser returns user
func (s *System) GetUser() ([]types.User, error) {
	defer TimeSpent("GetUser", time.Now())

	path := fmt.Sprintf("/api/instances/System::%v/relationships/User",
		s.System.ID)

	var user []types.User
	err := s.client.getJSONWithRetry(
		http.MethodGet, path, nil, &user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

// GetUserByID returns a specific user based on it's user id
func (s *System) GetUserByID(userId string) (*types.User, error) {
	path := fmt.Sprintf("/api/instances/User::%v",
		userId)

	user := &types.User{}
	err := s.client.getJSONWithRetry(
		http.MethodGet, path, nil, &user)
	if err != nil {
		return nil, err
	}
	return user, nil
}

// CreateUser creates a new user with some role.
func (s *System) CreateUser(userParam *types.UserParam) (*types.UserResp, error) {
	userResp := &types.UserResp{}
	err := s.client.getJSONWithRetry(
		http.MethodPost, "/api/types/User/instances", userParam, &userResp)
	if err != nil {
		return nil, err
	}
	return userResp, nil
}

// RemoveUser removes a particular user.
func (s *System) RemoveUser(userId string) error {
	path := fmt.Sprintf("/api/instances/User::%v/action/removeUser", userId)
	empty := &types.EmptyPayload{}
	err := s.client.getJSONWithRetry(
		http.MethodPost, path, empty, nil)
	return err
}

// SetUserRole sets a new role for a particular user.
func (s *System) SetUserRole(userRole *types.UserRoleParam, userId string) error {
	path := fmt.Sprintf("/api/instances/User::%v/action/setUserRole", userId)
	err := s.client.getJSONWithRetry(
		http.MethodPost, path, userRole, nil)
	return err
}
