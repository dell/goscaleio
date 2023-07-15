package goscaleio

import (
	"fmt"
	"net/http"

	types "github.com/dell/goscaleio/types/v1"
)

// CreateNFSExport create an NFS Export for a File System.
func (c *Client) CreateTreeQuota(createParams *types.TreeQuotaCreate) (respnfs *types.TreeQuotaCreateResponse, err error) {
	path := fmt.Sprintf("/rest/v1//file-tree-quotas")

	var body *types.TreeQuotaCreate = createParams
	err = c.getJSONWithRetry(http.MethodPost, path, body, &respnfs)
	if err != nil {
		return nil, err
	}

	return respnfs, nil
}
