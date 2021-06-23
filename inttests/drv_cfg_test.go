/*
 *
 * Copyright © 2020 Dell Inc. or its subsidiaries. All Rights Reserved.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */

package inttests

import (
	"fmt"
	"os"
	"testing"

	"github.com/dell/goscaleio"
	"github.com/stretchr/testify/assert"
)

// TestGetDrvCfgGUID will return the SDC GUID for the locally installed SDC
func TestGetDrvCfgGUID(t *testing.T) {
	goscaleio.SDCDevice = goscaleio.IOCTLDevice
	guid, err := goscaleio.DrvCfgQueryGUID()

	// The response depends on the installation state of the SDC
	if goscaleio.DrvCfgIsSDCInstalled() {
		// SDC is installed, we should get back a GUID and no error
		assert.NotEmpty(t, guid)
		assert.Nil(t, err)
		assert.Equal(t, os.Getenv("GOSCALEIO_SDC_GUID"), guid)
	} else {
		// SDC is not installed, we should get an emptry string for GUID and an error
		assert.Empty(t, guid)
		assert.NotNil(t, err)
		t.Skip("PowerFlex SDC is not installed. Cannot validate DrvCfg functionality")
	}
}

// TestGetDrvCfgGUIDSDCNotInstalled will check the SDC GUID, when the SDC is not installed
func TestGetDrvCfgGUIDSDCNotInstalled(t *testing.T) {
	goscaleio.SDCDevice = "/fff/dddddd/dddd"
	guid, err := goscaleio.DrvCfgQueryGUID()

	assert.Empty(t, guid)
	assert.NotNil(t, err)
	if goscaleio.DrvCfgIsSDCInstalled() {
		assert.Equal(t, false, true, "Expected the SDC to report as not installed")
	}
}

// TestGetSrvCfgSystems will return the PowerFlex systems connected to tyhe local SDC
func TestGetDrvCfgSystems(t *testing.T) {
	goscaleio.SDCDevice = goscaleio.IOCTLDevice
	systems, err := goscaleio.DrvCfgQuerySystems()

	// The response depends on the installation state of the SDC
	if goscaleio.DrvCfgIsSDCInstalled() {
		// SDC is installed, should get at least one ConfiguredSystem
		assert.NotEmpty(t, systems)
		assert.Nil(t, err)
		bFoundSystem := false
		for _, s := range *systems {
			assert.NotEqual(t, "", s.SystemID)
			assert.NotEqual(t, "", s.SdcID)
			if s.SystemID == os.Getenv("GOSCALEIO_SYSTEMID") {
				bFoundSystem = true
			}
		}
		assert.Equal(t, true, bFoundSystem, "Unable to find correct MDM system of %s", os.Getenv("GOSCALEIO_SYSTEMID"))
		assert.Equal(t, os.Getenv("GOSCALEIO_NUMBER_SYSTEMS"), fmt.Sprintf("%d", len(*systems)))
	} else {
		// SDC is not installed, should get no ConfiguredSystems and an error
		assert.Empty(t, systems)
		assert.NotNil(t, err)
		t.Skip("PowerFlex SDC is not installed. Cannot validate DrvCfg functionality")
	}
}

// TestGetSrvCfgSystemsSDCNotInstalled will check the PowerFlex systems when the SDC is not installed
func TestGetDrvCfgSystemsSDCNotInstalled(t *testing.T) {
	goscaleio.SDCDevice = "/fff/dddddd/dddd"
	systems, err := goscaleio.DrvCfgQuerySystems()

	// SDC is not installed, should get no ConfiguredSystems and an error
	assert.Empty(t, systems)
	assert.NotNil(t, err)

	if goscaleio.DrvCfgIsSDCInstalled() {
		assert.Equal(t, false, true, "Expected the SDC to report as not installed")
	}
}

func TestGetDrvCfgQueryRescan(t *testing.T) {
	goscaleio.SDCDevice = goscaleio.IOCTLDevice
	rc, err := goscaleio.DrvCfgQueryRescan()

	// The response depends on the installation state of the SDC
	if goscaleio.DrvCfgIsSDCInstalled() {
		// SDC is installed, we should get back a GUID and no error
		assert.NotEmpty(t, rc)
		assert.Nil(t, err)
		assert.Equal(t, rc, "0", "Device Rescan successful")
	} else {
		// SDC is not installed, we should get an emptry string for GUID and an error
		assert.Empty(t, rc)
		assert.NotNil(t, err)
		t.Skip("PowerFlex SDC is not installed. Cannot validate DrvCfg functionality")
	}
}

// TestGetSrvCfgSystemsSDCNotInstalled will check the PowerFlex systems when the SDC is not installed
func TestGetDrvCfgQueryRescanSDCNotInstalled(t *testing.T) {
	goscaleio.SDCDevice = "/fff/dddddd/dddd"
	systems, err := goscaleio.DrvCfgQueryRescan()

	// SDC is not installed, should get no ConfiguredSystems and an error
	assert.Empty(t, systems)
	assert.NotNil(t, err)

	if goscaleio.DrvCfgIsSDCInstalled() {
		assert.Equal(t, false, true, "Expected the SDC to report as not installed")
	}
}
