package goscaleio

import (
	"fmt"
	"net/http"
	"time"

	types "github.com/dell/goscaleio/types/v1"
)

type NFSExportCreate struct {
	// NFS Export name
	Name string `json:"name"`
	// NFS Export description
	// Unique identifier of the file system on which the NFS Export was created
	FileSystemID string `json:"file_system_id"`
	// Local path to a location within the file system.
	Path string `json:"path"`
}

type NFSExportModify struct {
	Description string `json:"description,omitempty"`
	// Default access level for all hosts that can access the Export
	// [ No_Access, Read_Only, Read_Write, Root, Read_Only_Root ]
	DefaultAccess string `json:"default_access,omitempty"`
	// Local path to a location within the file system.
	// With NFS, each export must have a unique local path.
	Path string `json:"path,omitempty"`
	// Read-Write hostsread_write_root_hosts
	ReadWriteHosts []string `json:"read_write_hosts,omitempty"`
	// Read-Only hosts
	ReadOnlyHosts []string `json:"read_only_hosts,omitempty"`
	// Read-Write, allow Root hosts
	ReadWriteRootHosts []string `json:"read_write_root_hosts,omitempty"`
	// Read-Only, allow Roots hosts
	ReadOnlyRootHosts []string `json:"read_only_root_hosts,omitempty"`
}

func (c *Client) GetNfsExport() (nfsList []types.NFSExport, err error) {
	defer TimeSpent("GetNfsExport", time.Now())

	fmt.Printf("nfsssss")
	path := fmt.Sprintf("/rest/v1/nfs-exports?select=*")

	err = c.getJSONWithRetry(
		http.MethodGet, path, nil, &nfsList)
	if err != nil {
		fmt.Println("error1:", err)
		return nil, err
	}

	return nfsList, nil
}

func (c *Client) CreateNFSExport(createParams *NFSExportCreate) (resp *types.CreateResponse, err error) {
	path := fmt.Sprintf("/rest/v1/nfs-exports")

	var body *NFSExportCreate = createParams
	err = c.getJSONWithRetry(http.MethodPost, path, body, &resp)
	if err != nil {
		fmt.Println("create err", err)
		return nil, err
	}

	return resp, nil

}

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

func (c *Client) DeleteNFSExport(id string) error {
	defer TimeSpent("DeleteNFSExport", time.Now())

	path := fmt.Sprintf("/rest/v1/nfs-exports/%s", id)

	err := c.getJSONWithRetry(
		http.MethodDelete, path, nil, nil)
	if err != nil {
		fmt.Println("errdelete", err)
		return err
	}

	return nil
}

func (c *Client) ModifyNFSExport(ModifyParams *NFSExportModify, id string) (err error) {
	path := fmt.Sprintf("/rest/v1/nfs-exports/%s", id)

	var body *NFSExportModify = ModifyParams
	err = c.getJSONWithRetry(http.MethodPatch, path, body, nil)
	if err != nil {
		fmt.Println("create err", err)
		return err
	}

	return nil

}
