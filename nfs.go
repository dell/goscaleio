package goscaleio

import (
	"fmt"
	"net/http"
	"time"

	types "github.com/dell/goscaleio/types/v1"
)

func (c *Client) GetNfsExport() ([]*types.NFSExport, error) {
	defer TimeSpent("GetNfsExport", time.Now())

	path := fmt.Sprintf("rest/v1/nfs-exports?select=*")

	var nfsList []*types.NFSExport
	err := c.getJSONWithRetry(
		http.MethodGet, path, nil, &nfsList)
	if err != nil {
		return nil, err
	}

	return nfsList, nil
}
