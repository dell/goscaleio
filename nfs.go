// Copyright © 2019 - 2023 Dell Inc. or its subsidiaries. All Rights Reserved.
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

// GetNASByName gets a NAS server by name
func (s *System) GetNASByName(name string) (*types.NAS, error) {
	var nasList []types.NAS

	path := fmt.Sprintf("/rest/v1/nas-servers?select=*")
	err := s.client.getJSONWithRetry(
		http.MethodGet, path, nil, &nasList)
	if err != nil {
		return nil, err
	}

	for _, nas := range nasList {
		if nas.Name == name {
			return &nas, nil
		}
	}
	return nil, errors.New("couldn't find given NAS server name")
}

// GetNAS gets a NAS server by ID
func (s *System) GetNAS(id string) (*types.NAS, error) {
	path := fmt.Sprintf("/rest/v1/nas-servers/%s?select=*", id)

	var resp *types.NAS
	err := s.client.getJSONWithRetry(
		http.MethodGet, path, nil, &resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// GetNASByIDName gets a NAS server by name or ID
func (s *System) GetNASBYIDName(id string, name string) (*types.NAS, error) {
	var nasList []types.NAS
	if id != "" {
		path := fmt.Sprintf("/rest/v1/nas-servers/%s?select=*", id)

		var resp *types.NAS
		err := s.client.getJSONWithRetry(
			http.MethodGet, path, nil, &resp)
		if err != nil {
			return nil, err
		}
		return resp, nil

	}

	if name != "" {
		path := fmt.Sprintf("/rest/v1/nas-servers?select=*")
		err := s.client.getJSONWithRetry(
			http.MethodGet, path, nil, &nasList)
		if err != nil {
			return nil, err
		}

		for _, nas := range nasList {
			if nas.Name == name {
				return &nas, nil
			}
		}
		return nil, errors.New("couldn't find given NAS server name")

	}
	if name == "" && id == "" {
		return nil, errors.New("id and name cannot be empty")
	}
	return nil, errors.New("couldn't find NAS server")
}

// CreateNAS creates a NAS server
func (s *System) CreateNAS(name string, protectionDomainID string) (*types.CreateNASResponse, error) {
	var resp types.CreateNASResponse

	path := fmt.Sprintf("/rest/v1/nas-servers")

	var body types.CreateNASParam = types.CreateNASParam{
		Name:               name,
		ProtectionDomainID: protectionDomainID,
	}

	err := s.client.getJSONWithRetry(http.MethodPost, path, body, &resp)

	if err != nil {
		return nil, err
	}

	return &resp, nil
}

// DeleteNAS deletes a NAS server
func (s *System) DeleteNAS(id string) error {

	path := fmt.Sprintf("/rest/v1/nas-servers/%s", id)

	err := s.client.getJSONWithRetry(http.MethodDelete, path, nil, nil)

	if err != nil {
		fmt.Println("err", err)
		return err
	}

	return nil
}
