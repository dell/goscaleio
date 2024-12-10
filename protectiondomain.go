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
	"context"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	types "github.com/dell/goscaleio/types/v1"
)

// ProtectionDomain defines a struct for ProtectionDomain
type ProtectionDomain struct {
	ProtectionDomain *types.ProtectionDomain
	client           *Client
}

// NewProtectionDomain returns a new ProtectionDomain
func NewProtectionDomain(client *Client) *ProtectionDomain {
	return &ProtectionDomain{
		ProtectionDomain: &types.ProtectionDomain{},
		client:           client,
	}
}

// NewProtectionDomainEx returns a new ProtectionDomain
func NewProtectionDomainEx(client *Client, pd *types.ProtectionDomain) *ProtectionDomain {
	return &ProtectionDomain{
		ProtectionDomain: pd,
		client:           client,
	}
}

// CreateProtectionDomain creates a ProtectionDomain
func (s *System) CreateProtectionDomain(ctx context.Context, name string) (string, error) {
	defer TimeSpent("CreateProtectionDomain", time.Now())

	protectionDomainParam := &types.ProtectionDomainParam{
		Name: name,
	}

	path := fmt.Sprintf("/api/types/ProtectionDomain/instances")

	pd := types.ProtectionDomainResp{}
	err := s.client.getJSONWithRetry(ctx,
		http.MethodPost, path, protectionDomainParam, &pd)
	if err != nil {
		return "", err
	}

	return pd.ID, nil
}

// GetProtectionDomainEx fetches a ProtectionDomain by ID with embedded client
func (s *System) GetProtectionDomainEx(ctx context.Context, id string) (*ProtectionDomain, error) {
	defer TimeSpent("GetProtectionDomainEx", time.Now())
	pdResp, err := s.FindProtectionDomainByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return NewProtectionDomainEx(s.client, pdResp), nil
}

// DeleteProtectionDomain will delete a protection domain
func (s *System) DeleteProtectionDomain(ctx context.Context, name string) error {
	// get the protection domain
	domain, err := s.FindProtectionDomain(ctx, "", name, "")
	if err != nil {
		return err
	}

	link, err := GetLink(domain.Links, "self")
	if err != nil {
		return err
	}

	protectionDomainParam := &types.EmptyPayload{}

	path := fmt.Sprintf("%v/action/removeProtectionDomain", link.HREF)

	err = s.client.getJSONWithRetry(ctx,
		http.MethodPost, path, protectionDomainParam, nil)
	if err != nil {
		return err
	}

	return nil
}

// Delete (ProtectionDomain) will delete a protection domain
func (pd *ProtectionDomain) Delete(ctx context.Context) error {
	link, err := GetLink(pd.ProtectionDomain.Links, "self")
	if err != nil {
		return err
	}

	protectionDomainParam := &types.EmptyPayload{}

	path := fmt.Sprintf("%v/action/removeProtectionDomain", link.HREF)

	err = pd.client.getJSONWithRetry(ctx,
		http.MethodPost, path, protectionDomainParam, nil)
	if err != nil {
		return err
	}

	return nil
}

// GetProtectionDomain returns a ProtectionDomain
func (s *System) GetProtectionDomain(ctx context.Context,
	pdhref string,
) ([]*types.ProtectionDomain, error) {
	defer TimeSpent("GetprotectionDomain", time.Now())

	var (
		err error
		pd  = &types.ProtectionDomain{}
		pds []*types.ProtectionDomain
	)

	if pdhref == "" {
		var link *types.Link
		link, err = GetLink(
			s.System.Links,
			"/api/System/relationship/ProtectionDomain")
		if err != nil {
			return nil, err
		}

		err = s.client.getJSONWithRetry(ctx,
			http.MethodGet, link.HREF, nil, &pds)
	} else {
		err = s.client.getJSONWithRetry(ctx,
			http.MethodGet, pdhref, nil, pd)
	}
	if err != nil {
		return nil, err
	}

	if pdhref != "" {
		pds = append(pds, pd)
	}
	return pds, nil
}

// FindProtectionDomain returns a ProtectionDomain
func (s *System) FindProtectionDomain(ctx context.Context,
	id, name, href string,
) (*types.ProtectionDomain, error) {
	defer TimeSpent("FindProtectionDomain", time.Now())

	pds, err := s.GetProtectionDomain(ctx, href)
	if err != nil {
		return nil, fmt.Errorf("Error getting protection domains %s", err)
	}

	for _, pd := range pds {
		if pd.ID == id || pd.Name == name || href != "" {
			return pd, nil
		}
	}

	return nil, errors.New("Couldn't find protection domain")
}

// FindProtectionDomainByID returns the ProtectionDomain having a particular ID
func (s *System) FindProtectionDomainByID(ctx context.Context, id string) (*types.ProtectionDomain, error) {
	defer TimeSpent("FindProtectionDomainByID", time.Now())

	href := fmt.Sprintf("/api/instances/ProtectionDomain::%s", id)
	pds, err := s.GetProtectionDomain(ctx, href)
	if err != nil {
		return nil, fmt.Errorf("error getting protection domain by id: %s", err)
	}
	if len(pds) == 0 {
		return nil, fmt.Errorf("no protection domain found having id=%s", id)
	}
	return pds[0], nil
}

// FindProtectionDomainByName returns the ProtectionDomain having a particular name
func (s *System) FindProtectionDomainByName(ctx context.Context, name string) (*types.ProtectionDomain, error) {
	defer TimeSpent("FindProtectionDomainByName", time.Now())

	var id string
	path := "/api/types/ProtectionDomain/instances/action/queryIdByKey"
	body := map[string]string{
		"name": name,
	}
	err := s.client.getJSONWithRetry(ctx, http.MethodPost, path, body, &id)
	if err != nil {
		return nil, fmt.Errorf("error getting protection domain by name: %s", err)
	}
	return s.FindProtectionDomainByID(ctx, id)
}

// SetName sets the name of the pd
func (pd *ProtectionDomain) SetName(ctx context.Context, name string) error {
	path := "/api/instances/ProtectionDomain::%s/action/setProtectionDomainName"
	nameParam := types.ProtectionDomainParam{
		Name: name,
	}
	return pd.setParam(ctx, path, nameParam)
}

// Refresh reads and stores current values of the pd
func (pd *ProtectionDomain) Refresh(ctx context.Context) error {
	defer TimeSpent("Refresh Protection Domain", time.Now())

	path := fmt.Sprintf("/api/instances/ProtectionDomain::%s", pd.ProtectionDomain.ID)

	pdResp := types.ProtectionDomain{}
	err := pd.client.getJSONWithRetry(ctx,
		http.MethodGet, path, &types.EmptyPayload{}, &pdResp)
	if err != nil {
		return err
	}
	pd.ProtectionDomain = &pdResp
	return nil
}

// SetRfcacheParams sets the Read Flash Cache params of the pd
func (pd *ProtectionDomain) SetRfcacheParams(ctx context.Context, params types.PDRfCacheParams) error {
	path := "/api/instances/ProtectionDomain::%s/action/setRfcacheParameters"
	return pd.setParam(ctx, path, params)
}

// SetSdsNetworkLimits sets IOPS limits on all SDS under the pd
func (pd *ProtectionDomain) SetSdsNetworkLimits(ctx context.Context, params types.SdsNetworkLimitParams) error {
	path := "/api/instances/ProtectionDomain::%s/action/setSdsNetworkLimits"
	return pd.setParam(ctx, path, params)
}

func (pd *ProtectionDomain) setParam(ctx context.Context, path string, param any) error {
	link := fmt.Sprintf(path, pd.ProtectionDomain.ID)
	return pd.client.getJSONWithRetry(ctx, http.MethodPost, link, param, nil)
}

// Activate activates the Protection domain
func (pd *ProtectionDomain) Activate(ctx context.Context, forceActivate bool) error {
	path := "/api/instances/ProtectionDomain::%s/action/activateProtectionDomain"
	return pd.setParam(ctx, path, map[string]string{
		"forceActivate": types.GetBoolType(forceActivate),
	})
}

// InActivate disables the Protection domain
func (pd *ProtectionDomain) InActivate(ctx context.Context, forceShutDown bool) error {
	path := "/api/instances/ProtectionDomain::%s/action/inactivateProtectionDomain"
	return pd.setParam(ctx, path, map[string]string{
		"forceShutdown": types.GetBoolType(forceShutDown),
	})
}

// EnableRfcache enables SDS Read Flash cache for entire Protection Domain
func (pd *ProtectionDomain) EnableRfcache(ctx context.Context) error {
	path := "/api/instances/ProtectionDomain::%s/action/enableSdsRfcache"
	return pd.setParam(ctx, path, &types.EmptyPayload{})
}

// DisableRfcache disables SDS Read Flash cache for entire Protection Domain
func (pd *ProtectionDomain) DisableRfcache(ctx context.Context) error {
	path := "/api/instances/ProtectionDomain::%s/action/disableSdsRfcache"
	return pd.setParam(ctx, path, &types.EmptyPayload{})
}

// DisableFGLMcache disables Fine Granularity Metadata cache for the Protection Domain
func (pd *ProtectionDomain) DisableFGLMcache(ctx context.Context) error {
	path := "/api/instances/ProtectionDomain::%s/action/disableFglMetadataCache"
	return pd.setParam(ctx, path, &types.EmptyPayload{})
}

// EnableFGLMcache enables Fine Granularity Metadata cache for the Protection Domain
func (pd *ProtectionDomain) EnableFGLMcache(ctx context.Context) error {
	path := "/api/instances/ProtectionDomain::%s/action/enableFglMetadataCache"
	return pd.setParam(ctx, path, &types.EmptyPayload{})
}

// SetDefaultFGLMcacheSize sets the default FGL Metadata for all SDSs under the Protection Domain
func (pd *ProtectionDomain) SetDefaultFGLMcacheSize(ctx context.Context, cacheSizeInMB int) error {
	path := "/api/instances/ProtectionDomain::%s/action/setDefaultFglMetadataCacheSize"
	return pd.setParam(ctx, path, map[string]string{
		"cacheSizeInMB": strconv.Itoa(cacheSizeInMB),
	})
}
