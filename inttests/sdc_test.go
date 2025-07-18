// Copyright © 2021 - 2024 Dell Inc. or its subsidiaries. All Rights Reserved.
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
	"testing"

	"github.com/dell/goscaleio"
	types "github.com/dell/goscaleio/types/v1"
	"github.com/stretchr/testify/assert"
)

// getAllSdc will return all Sdc instances
func getAllSdc(t *testing.T) []*goscaleio.Sdc {
	system := getSystem()
	assert.NotNil(t, system)
	if system == nil {
		return nil
	}

	var allSdc []*goscaleio.Sdc
	sdc, err := system.GetSdc()
	assert.Nil(t, err)
	assert.NotZero(t, len(sdc))
	for _, s := range sdc {
		outSdc := goscaleio.NewSdc(C, &s) // #nosec G601
		allSdc = append(allSdc, outSdc)
	}
	return allSdc
}

// TestGetSdcs will return all Sdc instances
func TestGetSdcs(t *testing.T) {
	getAllSdc(t)
}

// TestGetSdcByAttribute gets a single specific Sdc by attribute
func TestGetSdcByAttribute(t *testing.T) {
	system := getSystem()
	assert.NotNil(t, system)
	if system == nil {
		return
	}

	Sdc := getAllSdc(t)
	assert.NotNil(t, Sdc)
	if Sdc == nil {
		return
	}

	found, err := system.FindSdc("Name", Sdc[0].Sdc.Name)
	assert.Nil(t, err)
	assert.NotNil(t, found)
	assert.Equal(t, Sdc[0].Sdc.Name, found.Sdc.Name)

	found, err = system.FindSdc("ID", Sdc[0].Sdc.ID)
	assert.Nil(t, err)
	assert.NotNil(t, found)
	assert.Equal(t, Sdc[0].Sdc.ID, found.Sdc.ID)

	found, err = system.FindSdc("SdcGUID", Sdc[0].Sdc.SdcGUID)
	assert.Nil(t, err)
	assert.NotNil(t, found)
	assert.Equal(t, Sdc[0].Sdc.SdcGUID, found.Sdc.SdcGUID)
}

// TestGetSdcByAttributeInvalid fails to get a single specific Sdc by attribute
func TestGetSdcByAttributeInvalid(t *testing.T) {
	system := getSystem()
	assert.NotNil(t, system)
	if system == nil {
		return
	}
	Sdc := getAllSdc(t)
	assert.NotNil(t, Sdc)
	if Sdc == nil {
		return
	}

	found, err := system.FindSdc("Name", invalidIdentifier)
	assert.NotNil(t, err)
	assert.Nil(t, found)

	found, err = system.FindSdc("ID", invalidIdentifier)
	assert.NotNil(t, err)
	assert.Nil(t, found)
}

// TestGetSdcStatistics
func TestGetSdcStatistics(t *testing.T) {
	Sdc := getAllSdc(t)
	assert.NotNil(t, Sdc)
	if Sdc == nil {
		return
	}

	for _, s := range Sdc {
		stats, err := s.GetStatistics()
		assert.Nil(t, err)
		assert.NotNil(t, stats)
	}
}

// TestGetSdcVolumes
func TestGetSdcVolumes(t *testing.T) {
	Sdc := getAllSdc(t)
	assert.NotNil(t, Sdc)
	if Sdc == nil {
		return
	}

	for _, s := range Sdc {
		_, err := s.GetVolume()
		assert.Nil(t, err)
	}
}

// TestFindSdcVolumes
func TestFindSdcVolumes(t *testing.T) {
	Sdc := getAllSdc(t)
	assert.NotNil(t, Sdc)
	if Sdc == nil {
		return
	}

	for _, s := range Sdc {
		_, err := s.FindVolumes()
		assert.Nil(t, err)
	}
}

// TestChangeSdcName function tests Change name functionality of SDC.
func TestChangeSdcName(t *testing.T) {
	system := getSystem()
	assert.NotNil(t, system)

	sdc, err := system.GetSdc()
	assert.Nil(t, err)
	firstSdc := sdc[0]

	baseName := firstSdc.Name
	nameChng, err := system.ChangeSdcName(firstSdc.ID, randString(10))
	assert.Nil(t, err)
	assert.NotNil(t, nameChng)
	nameChngBack, err := system.ChangeSdcName(firstSdc.ID, baseName)
	assert.Nil(t, err)
	assert.NotNil(t, nameChngBack)
}

// TestChangeSdcPerfProfile function tests Change PerfProfile functionality of SDC.
func TestChangeSdcPerfProfile(t *testing.T) {
	system := getSystem()
	assert.NotNil(t, system)

	sdc, err := system.GetSdc()
	assert.Nil(t, err)
	firstSdc := sdc[0]

	basePerfProgile := firstSdc.PerfProfile
	ppChng, err := system.ChangeSdcPerfProfile(firstSdc.ID, "Compact")
	assert.Nil(t, err)
	assert.NotNil(t, ppChng)
	nameChngBack, err := system.ChangeSdcPerfProfile(firstSdc.ID, basePerfProgile)
	assert.Nil(t, err)
	assert.NotNil(t, nameChngBack)
}

// TestDeleteSdc will attempt to delete an SDS, which results in faliure
func TestDeleteSdc(t *testing.T) {
	system := getSystem()
	assert.NotNil(t, system)

	sdc, err := system.GetSdc()
	assert.Nil(t, err)
	firstSdc := sdc[0]

	sdsID := firstSdc.ID
	err = system.DeleteSdc(sdsID)
	assert.NotNil(t, err)
}

// GetSdcIdByIP will attempt to get SDC ID By IP Address
func TestGetSdcIdByIP(t *testing.T) {
	system := getSystem()
	assert.NotNil(t, system)

	sdc, err := system.GetSdc()
	assert.Nil(t, err)
	firstSdc := sdc[0]

	sdsIP := firstSdc.SdcIP
	sdcID, err := system.GetSdcIDByIP(sdsIP)
	assert.NotNil(t, sdcID)
	assert.Nil(t, err)
}

func TestSdcRestrictedMode(t *testing.T) {
	system := getSystem()
	assert.NotNil(t, system)

	err := system.SetRestrictedMode("Guid")
	assert.Nil(t, err)

	param := &types.ApproveSdcParam{
		SdcIP: "10.10.10.10",
	}

	_, err = system.ApproveSdc(param)
	assert.NotNil(t, err)

	err = system.SetApprovedIps("62276a432d28538a", []string{"10.10.10.10"})
	assert.NotNil(t, err)

	err = system.SetRestrictedMode("None")
	assert.Nil(t, err)
}

// TestGetVolumeMetrics
func TestGetVolumeMetrics(t *testing.T) {
	Sdc := getAllSdc(t)
	assert.NotNil(t, Sdc)
	if Sdc == nil {
		return
	}

	for _, s := range Sdc {
		_, err := s.GetVolumeMetrics()
		assert.Nil(t, err)
	}
}
