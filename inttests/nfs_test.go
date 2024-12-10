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
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

// getNasName returns GOSCALEIO_NASSERVER, if set.
func getNasName(t *testing.T) string {
	if os.Getenv("GOSCALEIO_NASSERVER") != "" {
		return os.Getenv("GOSCALEIO_NASSERVER")
	}
	system := getSystem()
	assert.NotNil(t, system)
	nasServer, _ := system.GetNASByIDName(context.Background(), "", "")
	assert.NotNil(t, nasServer)
	if nasServer == nil {
		return ""
	}
	fmt.Printf("nas server[0].Name: %v", nasServer.Name)
	return nasServer.Name
}

// TestGetNasByIDName gets a single specific nas server by Name or ID
func TestGetNASByIDName(t *testing.T) {
	system := getSystem()
	assert.NotNil(t, system)

	nasName := getNasName(t)
	assert.NotZero(t, len(nasName))

	ctx := context.Background()

	nasserver, err := system.GetNASByIDName(ctx, "", nasName)
	assert.Nil(t, err)
	assert.Equal(t, nasName, nasserver.Name)

	if nasserver != nil {
		nas, err := system.GetNASByIDName(ctx, nasserver.ID, "")
		assert.Nil(t, err)
		assert.Equal(t, nasserver.ID, nas.ID)
	}

	if len(nasName) > 0 {
		nas, err := system.GetNASByIDName(ctx, "", nasName)
		assert.Nil(t, err)
		assert.Equal(t, nasName, nas.Name)
	}
}

// TestGetNasByIDNameInvalid attempts to get a file system that does not exist
func TestGetNasByIDNameInvalid(t *testing.T) {
	system := getSystem()
	assert.NotNil(t, system)

	ctx := context.Background()

	nas, err := system.GetNASByIDName(ctx, invalidIdentifier, "")
	assert.NotNil(t, err)
	assert.Nil(t, nas)

	nasName, err := system.GetNASByIDName(ctx, "", invalidIdentifier)
	assert.NotNil(t, err)
	assert.Nil(t, nasName)
}

func TestCreateDeleteNAS(t *testing.T) {
	system := getSystem()
	assert.NotNil(t, system)

	nasName := fmt.Sprintf("%s-%s", testPrefix, "twee2")

	ctx := context.Background()

	// get protection domain
	pd := getProtectionDomain(t)
	assert.NotNil(t, pd)

	var pdID string

	if pd != nil {
		pDomain, _ := system.FindProtectionDomain(ctx, pd.ProtectionDomain.ID, "", "")
		assert.NotNil(t, pDomain)
		pdID = pDomain.ID
	}

	// create the NAS Server
	nasID, err := system.CreateNAS(ctx, nasName, pdID)
	assert.Nil(t, err)
	assert.NotNil(t, nasID.ID)

	// try to create a NAS Server that exists
	_, err = system.CreateNAS(ctx, nasName, pdID)
	assert.NotNil(t, err)

	// delete the NAS Server
	err = system.DeleteNAS(ctx, nasID.ID)
	assert.Nil(t, err)

	// try to delete non-existent NAS Server
	err = system.DeleteNAS(ctx, nasID.ID)
	assert.NotNil(t, err)
}

func TestGetFileInterfaceById(t *testing.T) {
	system := getSystem()
	nasName := getNasName(t)
	assert.NotZero(t, len(nasName))

	ctx := context.Background()

	nasserver, err := system.GetNASByIDName(ctx, "", nasName)
	assert.Nil(t, err)
	assert.Equal(t, nasName, nasserver.Name)

	if nasserver != nil {
		fileInterface, err := system.GetFileInterface(ctx, nasserver.CurrentPreferredIPv4InterfaceID)
		assert.Nil(t, err)
		assert.NotNil(t, fileInterface)
	}
}
