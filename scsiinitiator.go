package goscaleio

import (
	"fmt"
	"net/http"
	"time"

	types "github.com/dell/goscaleio/types/v1"
)

func (s *System) GetScsiInitiator() ([]types.ScsiInitiator, error) {
	defer TimeSpent("GetScsiInitiator", time.Now())

	path := fmt.Sprintf(
		"/api/instances/System::%v/relationships/ScsiInitiator",
		s.System.ID)

	var si []types.ScsiInitiator
	err := s.client.getJSONWithRetry(
		http.MethodGet, path, nil, &si)
	if err != nil {
		return nil, err
	}

	return si, nil
}
