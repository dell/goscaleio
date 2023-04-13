package goscaleio

import (
	"fmt"
	"net/http"
	"time"

	types "github.com/dell/goscaleio/types/v1"
)

func (c *Client) GetNfsExport() (nfsList []types.NFSExport, err error) {
	defer TimeSpent("GetNfsExport", time.Now())

	fmt.Printf("nfsssss")
	path := fmt.Sprintf("/rest/v1/nfs-exports/?select=*")

	err = c.getJSONWithRetry(
		http.MethodGet, path, nil, &nfsList)
	if err != nil {
		fmt.Println("error1:", err)
		return nil, err
	}

	return nfsList, nil
}

func (c *Client) CreateNFSExport(createParams *types.NFSExportCreate) (resp *types.CreateResponse, err error) {
	path := fmt.Sprintf("/rest/v1/nfs-exports")

	var body *types.NFSExportCreate
	body = createParams
	err = c.getJSONWithRetry(http.MethodPost, path, body, &resp)
	if err != nil {
		return nil, err
	}

	return resp, nil

}

// func (c *Client) GetNAS(id string) ([]types.NAS, error) {
// 	fmt.Println("Inside GetNAS- line 36")
// 	path := fmt.Sprintf("/rest/v1/nas-servers/%s?select=*", id)
// 	fmt.Println("Inside GetNAS- line 38")
// 	var resp []types.NAS
// 	err := c.getJSONWithRetry(
// 		http.MethodGet, path, nil, &resp)
// 	if err != nil {
// 		return nil, err
// 	}
// 	fmt.Println("Inside GetNAS- line 44")
// 	return resp, nil
// }

func (c *Client) GetNfsExportById(id string) (resp *types.NFSExport, err error) {
	defer TimeSpent("GetNfsExport", time.Now())

	fmt.Printf("nfsssssid")
	path := fmt.Sprintf("/rest/v1/nfs-exports/%s?select=*", id)

	err = c.getJSONWithRetry(
		http.MethodGet, path, nil, &resp)
	if err != nil {
		fmt.Println("error1:", err)
		return nil, err
	}
	return resp, nil
}
