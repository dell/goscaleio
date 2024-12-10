// Copyright © 2021 - 2022 Dell Inc. or its subsidiaries. All Rights Reserved.
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
	"testing"

	"github.com/dell/goscaleio"
	types "github.com/dell/goscaleio/types/v1"
	"github.com/stretchr/testify/assert"
)

// get the system, this code does not check errors as it is a simple helper
// and there are specific tests got querying the Systems
func getSystem() *goscaleio.System {
	ctx := context.Background()

	// first, get all of the systems
	allSystems, _ := C.GetSystems(ctx)

	// then try to get the first one returned, explicitly
	system, _ := C.FindSystem(ctx, allSystems[0].ID, "", "")

	return system
}

// TestGetSystems will detect all of the PowerFlex systems connected to this Gateway
// There should be EXACTLY ONE system, as two/more is not supported
func TestGetSystems(t *testing.T) {
	allSystems, err := C.GetSystems(context.Background())
	assert.Nil(t, err)
	assert.Equal(t, 1, len(allSystems))
}

// TestGetSystem will search the systems for a specific one
func TestGetSingleSystemID(t *testing.T) {
	ctx := context.Background()

	// first, get all of the systems
	allSystems, err := C.GetSystems(ctx)
	assert.Nil(t, err)

	// then try to get the first one returned, explicitly
	system, err := C.FindSystem(ctx, allSystems[0].ID, "", "")
	assert.Nil(t, err)
	assert.Equal(t, allSystems[0].ID, system.System.ID)
}

// TestGetSingleSystemByIDInvalid will search the for a system that will not be found
func TestGetSingleSystemByIDInvalid(t *testing.T) {
	// then try to get the first one returned, explicitly
	system, err := C.FindSystem(context.Background(), invalidIdentifier, "", "")
	fmt.Printf("system %v err %s\n", system, err)
	assert.NotNil(t, err)
	assert.Nil(t, system)
}

// TestGetSingleSystemByNameInvalid will search the for a system that will not be found
func TestGetSingleSystemByIDName(t *testing.T) {
	// then try to get the first one returned, explicitly
	system, err := C.FindSystem(context.Background(), "", invalidIdentifier, "")
	assert.NotNil(t, err)
	assert.Nil(t, system)
}

// TestGetSystemStatistics will return System statistics
func TestGetSystemStatistics(t *testing.T) {
	ctx := context.Background()

	// first, get all of the systems
	allSystems, err := C.GetSystems(ctx)
	assert.Nil(t, err)

	// then try to get the first one returned, explicitly
	system, err := C.FindSystem(ctx, allSystems[0].ID, "", "")
	assert.Nil(t, err)
	assert.Equal(t, allSystems[0].ID, system.System.ID)

	stats, err := system.GetStatistics(ctx)
	assert.Nil(t, err)
	assert.NotNil(t, stats)
}

// TestGetMDMClusterDetails will return MDM cluster details
func TestGetMDMClusterDetails(t *testing.T) {
	ctx := context.Background()

	// first, get all of the systems
	allSystems, err := C.GetSystems(ctx)
	assert.Nil(t, err)

	// then try to get the first one returned, explicitly
	system, err := C.FindSystem(ctx, allSystems[0].ID, "", "")
	assert.Nil(t, err)
	assert.Equal(t, allSystems[0].ID, system.System.ID)

	mdmDetails, err := system.GetMDMClusterDetails(ctx)
	assert.Nil(t, err)
	assert.NotNil(t, mdmDetails)
}

// TestMDMClusterPerformanceProfile modifies MDM cluster performace profile
func TestMDMClusterPerformanceProfileInvalid(t *testing.T) {
	ctx := context.Background()

	// first, get all of the systems
	allSystems, err := C.GetSystems(ctx)
	assert.Nil(t, err)

	// then try to get the first one returned, explicitly
	system, err := C.FindSystem(ctx, allSystems[0].ID, "", "")
	assert.Nil(t, err)
	assert.Equal(t, allSystems[0].ID, system.System.ID)

	err = system.ModifyPerformanceProfileMdmCluster(ctx, "High")
	assert.NotNil(t, err)
}

// TestAddStandByMDMInvalid attempts to add invalid standby MDM, results in failure
func TestAddStandByMDMInvalid(t *testing.T) {
	ctx := context.Background()

	allSystems, err := C.GetSystems(ctx)
	assert.Nil(t, err)

	// then try to get the first one returned, explicitly
	system, err := C.FindSystem(ctx, allSystems[0].ID, "", "")
	assert.Nil(t, err)
	assert.Equal(t, allSystems[0].ID, system.System.ID)

	payload := types.StandByMdm{
		IPs:  []string{"0.1.1.1", "0.2.2.2"},
		Role: "Manager",
	}

	_, err1 := system.AddStandByMdm(ctx, &payload)
	assert.NotNil(t, err1)
}

// TestRemoveStandByMDMInvalid attempts to remove invalid standby MDM, results in failure
func TestRemoveStandByMDMInvalid(t *testing.T) {
	ctx := context.Background()

	allSystems, err := C.GetSystems(ctx)
	assert.Nil(t, err)

	// then try to get the first one returned, explicitly
	system, err := C.FindSystem(ctx, allSystems[0].ID, "", "")
	assert.Nil(t, err)
	assert.Equal(t, allSystems[0].ID, system.System.ID)

	mdmDetails, err1 := system.GetMDMClusterDetails(ctx)
	assert.Nil(t, err1)
	assert.NotNil(t, mdmDetails)

	if len(mdmDetails.StandByMdm) > 0 {
		err = system.RemoveStandByMdm(ctx, mdmDetails.StandByMdm[0].ID)
		assert.Nil(t, err)
	}
}

// TestSwitchClusterModeInvalid attempts to switch cluster mode with invalid mode
func TestSwitchClusterModeInvalid(t *testing.T) {
	ctx := context.Background()

	allSystems, err := C.GetSystems(ctx)
	assert.Nil(t, err)

	// then try to get the first one returned, explicitly
	system, err := C.FindSystem(ctx, allSystems[0].ID, "", "")
	assert.Nil(t, err)
	assert.Equal(t, allSystems[0].ID, system.System.ID)

	mdmDetails, err1 := system.GetMDMClusterDetails(ctx)
	assert.Nil(t, err1)
	assert.NotNil(t, mdmDetails)

	payload := types.SwitchClusterMode{
		Mode:             "FourNodes",
		AddSecondaryMdms: []string{"invalid"},
	}

	err = system.SwitchClusterMode(ctx, &payload)
	assert.NotNil(t, err)
}

// TestRenameMdm modifies MDM name
func TestRenameMdm(t *testing.T) {
	ctx := context.Background()

	// first, get all of the systems
	allSystems, err := C.GetSystems(ctx)
	assert.Nil(t, err)

	// then try to get the first one returned, explicitly
	system, err := C.FindSystem(ctx, allSystems[0].ID, "", "")
	assert.Nil(t, err)

	mdmDetails, err1 := system.GetMDMClusterDetails(ctx)
	assert.Nil(t, err1)
	assert.NotNil(t, mdmDetails)

	payload := types.RenameMdm{
		ID:      mdmDetails.TiebreakerMdm[0].ID,
		NewName: "mdm_renamed",
	}

	err = system.RenameMdm(ctx, &payload)
	assert.Nil(t, err)

	payload = types.RenameMdm{
		ID:      mdmDetails.TiebreakerMdm[0].ID,
		NewName: "mdm_renamed",
	}

	err = system.RenameMdm(ctx, &payload)
	assert.NotNil(t, err)

	payload = types.RenameMdm{
		ID:      mdmDetails.TiebreakerMdm[0].ID,
		NewName: "tb_mdm",
	}

	err = system.RenameMdm(ctx, &payload)
	assert.Nil(t, err)
}

// TestChangeMDMOwnership modifies primary MDM
func TestChangeMDMOwnership(t *testing.T) {
	ctx := context.Background()

	// first, get all of the systems
	allSystems, err := C.GetSystems(ctx)
	assert.Nil(t, err)

	// then try to get the first one returned, explicitly
	system, err := C.FindSystem(ctx, allSystems[0].ID, "", "")
	assert.Nil(t, err)
	assert.Equal(t, allSystems[0].ID, system.System.ID)

	mdmDetails, err1 := system.GetMDMClusterDetails(ctx)
	assert.Nil(t, err1)
	assert.NotNil(t, mdmDetails)

	err = system.ChangeMdmOwnerShip(ctx, mdmDetails.SecondaryMDM[0].ID)
	assert.Nil(t, err)
}
