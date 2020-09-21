package goscaleio

import (
	"fmt"
	"net/http"
	"time"

	types "github.com/dell/goscaleio/types/v1"
)

type System struct {
	System *types.System
	client *Client
}

func NewSystem(client *Client) *System {
	return &System{
		System: &types.System{},
		client: client,
	}
}

func (c *Client) FindSystem(
	instanceID, name, href string) (*System, error) {
	defer TimeSpent("FindSystem", time.Now())

	systems, err := c.GetInstance(href)
	if err != nil {
		return nil, fmt.Errorf("err: problem getting instances: %s", err)
	}

	for _, system := range systems {
		if system.ID == instanceID || system.Name == name || href != "" {
			outSystem := NewSystem(c)
			outSystem.System = system
			return outSystem, nil
		}
	}
	return nil, fmt.Errorf("err: systemid or systemname not found")
}

func (s *System) GetStatistics() (*types.Statistics, error) {
	defer TimeSpent("GetStatistics", time.Now())

	link, err := GetLink(s.System.Links,
		"/api/System/relationship/Statistics")
	if err != nil {
		return nil, err
	}

	stats := types.Statistics{}
	err = s.client.getJSONWithRetry(
		http.MethodGet, link.HREF, nil, &stats)
	if err != nil {
		return nil, err
	}

	return &stats, nil
}

func (s *System) CreateSnapshotConsistencyGroup(
	snapshotVolumesParam *types.SnapshotVolumesParam) (*types.SnapshotVolumesResp, error) {
	defer TimeSpent("CreateSnapshotConsistencyGroup", time.Now())

	link, err := GetLink(s.System.Links, "self")
	if err != nil {
		return nil, err
	}

	path := fmt.Sprintf("%v/action/snapshotVolumes", link.HREF)

	snapResp := types.SnapshotVolumesResp{}
	err = s.client.getJSONWithRetry(
		http.MethodPost, path, snapshotVolumesParam, &snapResp)
	if err != nil {
		return nil, err
	}

	return &snapResp, nil
}
