// Copyright © 2021 - 2023 Dell Inc. or its subsidiaries. All Rights Reserved.
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

package inttests

import (
	"fmt"
	"os"
	"testing"

	//log "github.com/sirupsen/logrus"
	//types "github.com/dell/goscaleio/types/v1"

	"github.com/stretchr/testify/assert"
)

// getNasName returns GOSCALEIO_NASSERVER, if set.
func getNasName(t *testing.T) string {
	if os.Getenv("GOSCALEIO_NASSERVER") != "" {
		return os.Getenv("GOSCALEIO_NASSERVER")
	}
	system := getSystem()
	assert.NotNil(t, system)
	nasServer, _ := system.GetNASByName("")
	assert.NotNil(t, nasServer)
	if nasServer == nil {
		return ""
	}
	fmt.Printf("nas server[0].Name: %v", nasServer.Name)
	return nasServer.Name
}

// TestGetNasByName gets a single specific nas server by Name
func TestGetNasByName(t *testing.T) {
	system := getSystem()
	assert.NotNil(t, system)

	nasName := getNasName(t)
	assert.NotZero(t, len(nasName))

	if len(nasName) > 0 {
		nas, err := system.GetNASByName(nasName)
		assert.Nil(t, err)
		assert.Equal(t, nasName, nas.Name)
	}
}

// TestGetNasByNameInvalid attempts to get a nas server that does not exist
func TestGetNasByNameInvalid(t *testing.T) {
	system := getSystem()
	assert.NotNil(t, system)

	nas, err := system.GetNAS(invalidIdentifier)
	assert.NotNil(t, err)
	assert.Nil(t, nas)
}

// TestGetNasByID will return nas server by id
func TestGetNasByID(t *testing.T) {
	system := getSystem()
	assert.NotNil(t, system)

	nasName := getNasName(t)
	assert.NotZero(t, len(nasName))

	nasserver, err := system.GetNASByName(nasName)
	assert.Nil(t, err)
	assert.Equal(t, nasName, nasserver.Name)

	if nasserver != nil {
		nas, err := system.GetNAS(nasserver.ID)
		assert.Nil(t, err)
		assert.Equal(t, nasserver.ID, nas.ID)
	}
}

// TestGetNasByIDInvalid attempts to get a file system that does not exist
func TestGetNasByIDInvalid(t *testing.T) {
	system := getSystem()
	assert.NotNil(t, system)

	nas, err := system.GetNAS(invalidIdentifier)
	assert.NotNil(t, err)
	assert.Nil(t, nas)
}

func TestCreateDeleteNAS(t *testing.T) {
	system := getSystem()
	assert.NotNil(t, system)

	nasName := fmt.Sprintf("%s-%s", testPrefix, "twee2")

	// get protection domain
	pd := getProtectionDomain(t)
	assert.NotNil(t, pd)

	var pdID string

	if pd != nil {
		pDomain, _ := system.FindProtectionDomain(pd.ProtectionDomain.ID, "", "")
		assert.NotNil(t, pDomain)
		pdID = pDomain.ID
	}

	// create the NAS Server
	nasID, err := system.CreateNAS(nasName, pdID)
	assert.Nil(t, err)
	assert.NotNil(t, nasID.ID)

	// try to create a NAS Server that exists
	_, err = system.CreateNAS(nasName, pdID)
	assert.NotNil(t, err)

	// delete the NAS Server
	err = system.DeleteNAS(nasID.ID)
	assert.Nil(t, err)

	// try to delete non-existent NAS Server
	err = system.DeleteNAS(nasID.ID)
	assert.NotNil(t, err)

}
