// // Copyright © 2021 - 2023 Dell Inc. or its subsidiaries. All Rights Reserved.
// //
// // Licensed under the Apache License, Version 2.0 (the "License");
// // you may not use this file except in compliance with the License.
// // You may obtain a copy of the License at
// //      http://www.apache.org/licenses/LICENSE-2.0
// // Unless required by applicable law or agreed to in writing, software
// // distributed under the License is distributed on an "AS IS" BASIS,
// // WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// // See the License for the specific language governing permissions and
// // limitations under the License.

package inttests

import (
	"os"
	"testing"

	types "github.com/dell/goscaleio/types/v1"
	"github.com/stretchr/testify/assert"
)

func GetNFSExportbyName(t *testing.T) string {
	if os.Getenv("GOSCALEIO_NFSEXPORT") != "" {
		return os.Getenv("GOSCALEIO_NFSEXPORT")
	}

	nfsexport, _ := C.GetNFSExport()
	assert.NotNil(t, nfsexport)
	if nfsexport == nil {
		return ""
	}
	return nfsexport[0].Name
}

func TestNFSExportByName(t *testing.T) {
	fsName := GetNFSExportbyName(t)
	assert.NotZero(t, len(fsName))

	if len(fsName) > 0 {
		fs, err := C.GetNFSExportByName(fsName)
		assert.Nil(t, err)
		assert.Equal(t, fsName, fs.Name)
	}
}

func TestNFSExportByID(t *testing.T) {
	nfsName := GetNFSExportbyName(t)
	assert.NotZero(t, len(nfsName))

	nfs, err := C.GetNFSExportByName(nfsName)
	assert.Nil(t, err)
	assert.Equal(t, nfsName, nfs.Name)

	if nfs != nil {
		fs, err := C.GetNFSExportByID(nfs.ID)
		assert.Nil(t, err)
		assert.Equal(t, nfs.ID, fs.ID)
	}
}

func TestNFSExportByNameInvalid(t *testing.T) {
	fs, err := C.GetNFSExportByName(invalidIdentifier)
	assert.NotNil(t, err)
	assert.Nil(t, fs)
}

func TestNFSExportByIDInvalid(t *testing.T) {
	fs, err := C.GetNFSExportByID(invalidIdentifier)
	assert.NotNil(t, err)
	assert.Nil(t, fs)
}

func TestCreateModifyDeleteNFSExport(t *testing.T) {
	system := getSystem()
	assert.NotNil(t, system)

	nfsName := "NFS" + testPrefix + randString(8)
	nfsmodify := "NFS export modify testing"
	var filesystemname string
	if os.Getenv("GOSCALEIO_FILESYSTEM_NFSEXPORT") != "" {
		filesystemname = os.Getenv("GOSCALEIO_FILESYSTEM_NFSEXPORT")
	}
	filesystem, err := system.GetFileSystemByName(filesystemname)
	nfsexport := &types.NFSExportCreate{
		Name:         nfsName,
		FileSystemID: filesystem.ID,
		Path:         "/" + filesystemname,
	}

	//create nfs export
	nfs, err := C.CreateNFSExport(nfsexport)
	fsID := nfs.ID
	assert.Nil(t, err)
	assert.NotNil(t, fsID)

	// try to create existing nfs export
	nfs, err = C.CreateNFSExport(nfsexport)
	assert.NotNil(t, err)

	//Modify NFS export proprties
	nfsexportmodify := &types.NFSExportModify{
		Description:           nfsmodify,
		AddReadWriteRootHosts: []string{"192.168.100.10", "192.168.100.11"},
	}
	err = C.ModifyNFSExport(nfsexportmodify, fsID)

	//delete the NFS export
	err = C.DeleteNFSExport(fsID)
	assert.Nil(t, err)

}
