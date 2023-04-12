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
	path := fmt.Sprintf("https://10.225.109.54/rest/v1/nfs-exports/?select=*")

	err = c.getJSONWithRetry(
		http.MethodGet, path, nil, &nfsList)
	if err != nil {
		fmt.Println("error1:", err)
		return nil, err
	}

	return nfsList, nil
}
