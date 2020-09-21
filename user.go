package goscaleio

import (
	"fmt"
	"net/http"
	"time"

	types "github.com/dell/goscaleio/types/v1"
)

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
