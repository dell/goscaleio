package goscaleio

import (
	"errors"
	"fmt"
	"net/http"
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
func (s *System) CreateProtectionDomain(name string) (string, error) {
	defer TimeSpent("CreateProtectionDomain", time.Now())

	protectionDomainParam := &types.ProtectionDomainParam{
		Name: name,
	}

	path := fmt.Sprintf("/api/types/ProtectionDomain/instances")

	pd := types.ProtectionDomainResp{}
	err := s.client.getJSONWithRetry(
		http.MethodPost, path, protectionDomainParam, &pd)
	if err != nil {
		return "", err
	}

	return pd.ID, nil
}

// DeleteProtectionDomain will delete a protection domain
func (s *System) DeleteProtectionDomain(name string) error {
	// get the protection domain
	domain, err := s.FindProtectionDomain("", name, "")
	if err != nil {
		return err
	}

	link, err := GetLink(domain.Links, "self")
	if err != nil {
		return err
	}

	protectionDomainParam := &types.EmptyPayload{}

	path := fmt.Sprintf("%v/action/removeProtectionDomain", link.HREF)

	err = s.client.getJSONWithRetry(
		http.MethodPost, path, protectionDomainParam, nil)
	if err != nil {
		return err
	}

	return nil
}

// GetProtectionDomain returns a ProtectionDomain
func (s *System) GetProtectionDomain(
	pdhref string) ([]*types.ProtectionDomain, error) {
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

		err = s.client.getJSONWithRetry(
			http.MethodGet, link.HREF, nil, &pds)
	} else {
		err = s.client.getJSONWithRetry(
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
func (s *System) FindProtectionDomain(
	id, name, href string) (*types.ProtectionDomain, error) {
	defer TimeSpent("FindProtectionDomain", time.Now())

	pds, err := s.GetProtectionDomain(href)
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
