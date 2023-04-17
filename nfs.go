package goscaleio

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	types "github.com/dell/goscaleio/types/v1"
)

// type NFSExportCreate struct {
// 	// NFS Export name
// 	Name string `json:"name"`
// 	// NFS Export description
// 	// Unique identifier of the file system on which the NFS Export was created
// 	FileSystemID string `json:"file_system_id"`
// 	// Local path to a location within the file system.
// 	Path string `json:"path"`
// }

// NGetNFSExport lists NFS Exports.
func (c *Client) GetNFSExport() (nfsList []types.NFSExport, err error) {
	defer TimeSpent("GetNfsExport", time.Now())
	path := fmt.Sprintf("/rest/v1/nfs-exports?select=*")

	err = c.getJSONWithRetry(
		http.MethodGet, path, nil, &nfsList)
	if err != nil {
		return nil, err
	}

	return nfsList, nil
}

// CreateNFSExport create an NFS Export for a File System.
func (c *Client) CreateNFSExport(createParams *types.NFSExportCreate) (respnfs *types.CreateResponse, err error) {
	path := fmt.Sprintf("/rest/v1/nfs-exports")

	var body *types.NFSExportCreate = createParams
	err = c.getJSONWithRetry(http.MethodPost, path, body, &respnfs)
	if err != nil {
		return nil, err
	}

	return respnfs, nil

}

// GetNFSExportById returns NFS Export properties by ID
func (c *Client) GetNFSExportById(id string) (respnfs *types.NFSExport, err error) {
	defer TimeSpent("GetNfsExport", time.Now())
	path := fmt.Sprintf("/rest/v1/nfs-exports/%s?select=*", id)

	err = c.getJSONWithRetry(
		http.MethodGet, path, nil, &respnfs)
	if err != nil {
		return nil, err
	}
	return respnfs, nil
}

// GetNFSExportByName returns NFS Export properties by name
func (c *Client) GetNFSExportByName(name string) (*types.NFSExport, error) {
	defer TimeSpent("GetFileSystemByName", time.Now())

	nfsList, err := c.GetNFSExport()
	if err != nil {
		return nil, err
	}

	for _, nfs := range nfsList {
		if nfs.Name == name {
			return &nfs, nil
		}
	}

	return nil, errors.New("Couldn't find file system")
}

// DeleteNFSExport deletes the NFS export
func (c *Client) DeleteNFSExport(id string) error {
	defer TimeSpent("DeleteNFSExport", time.Now())
	path := fmt.Sprintf("/rest/v1/nfs-exports/%s", id)

	err := c.getJSONWithRetry(
		http.MethodDelete, path, nil, nil)
	if err != nil {
		return err
	}

	return nil
}

// ModifyNFSExport modifies the NFS export properties
func (c *Client) ModifyNFSExport(ModifyParams *types.NFSExportModify, id string) (err error) {
	path := fmt.Sprintf("/rest/v1/nfs-exports/%s", id)

	var body *types.NFSExportModify = ModifyParams
	err = c.getJSONWithRetry(http.MethodPatch, path, body, nil)
	if err != nil {
		return err
	}

	return nil
}
