// Copyright © 2019 - 2023 Dell Inc. or its subsidiaries. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//      http://www.apache.org/licenses/LICENSE-2.0
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package goscaleio

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	types "github.com/dell/goscaleio/types/v1"
)

// FileSystem defines struct for file system
type FileSystem struct {
	FileSystem *types.FileSystem
	client     *Client
}

// NewFileSystem returns a new file system
func NewFileSystem(client *Client, fs *types.FileSystem) *FileSystem {
	return &FileSystem{
		FileSystem: fs,
		client:     client,
	}
}

// GetAllFileSystems returns a file system
func (s *System) GetAllFileSystems() ([]types.FileSystem, error) {
	defer TimeSpent("GetAllFileSystems", time.Now())

	path := fmt.Sprintf("/rest/v1/file-systems?select=*")
	var fs []types.FileSystem
	err := s.client.getJSONWithRetry(
		http.MethodGet, path, nil, &fs)
	if err != nil {
		return nil, err
	}

	return fs, nil
}

// GetFileSystemByID returns a file system by ID
func (s *System) GetFileSystemByID(id string) (*types.FileSystem, error) {
	defer TimeSpent("GetFileSystemByID", time.Now())

	path := fmt.Sprintf("/rest/v1/file-systems/%v?select=*", id)
	var fs types.FileSystem
	err := s.client.getJSONWithRetry(
		http.MethodGet, path, nil, &fs)
	if err != nil {
		return nil, err
	}

	return &fs, nil
}

// GetFileSystemByName returns a file system by Name
func (s *System) GetFileSystemByName(name string) (*types.FileSystem, error) {
	defer TimeSpent("GetFileSystemByID", time.Now())

	filesystems, err := s.GetAllFileSystems()
	if err != nil {
		return nil, err
	}

	for _, fs := range filesystems {
		if fs.Name == name {
			return &fs, nil
		}
	}

	return nil, errors.New("Couldn't find storage pool")
}

// CreateFileSystem creates a file system
func (s *System) CreateFileSystem(fs *types.FsCreate) (string, error) {
	defer TimeSpent("CreateFileSystem", time.Now())

	path := fmt.Sprintf("/rest/v1/file-systems")
	fsResponse := types.FileSystemResp{}
	err := s.client.getJSONWithRetry(
		http.MethodPost, path, fs, &fsResponse)
	if err != nil {
		return " ", err
	}

	return fsResponse.ID, nil
}

// DeleteFileSystem deletes a file system
func (s *System) DeleteFileSystem(name string) error {
	defer TimeSpent("DeleteFileSystem", time.Now())

	fs, err := s.GetFileSystemByName(name)
	if err != nil {
		return err
	}

	path := fmt.Sprintf("/rest/v1/file-systems/%v", fs.ID)

	err = s.client.getJSONWithRetry(
		http.MethodDelete, path, nil, nil)
	if err != nil {
		return err
	}

	return nil
}
