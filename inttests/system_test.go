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
	"fmt"
	"testing"

	"github.com/dell/goscaleio"
	"github.com/stretchr/testify/assert"
)

// get the system, this code does not check errors as it is a simple helper
// and there are specific tests got querying the Systems
func getSystem() *goscaleio.System {
	// first, get all of the systems
	allSystems, _ := C.GetSystems()

	// then try to get the first one returned, explicitly
	system, _ := C.FindSystem(allSystems[0].ID, "", "")

	return system
}

// TestGetSystems will detect all of the PowerFlex systems connected to this Gateway
// There should be EXACTLY ONE system, as two/more is not supported
func TestGetSystems(t *testing.T) {
	allSystems, err := C.GetSystems()
	assert.Nil(t, err)
	assert.Equal(t, 1, len(allSystems))
}

// TestGetSystem will search the systems for a specific one
func TestGetSingleSystemID(t *testing.T) {
	// first, get all of the systems
	allSystems, err := C.GetSystems()
	assert.Nil(t, err)

	// then try to get the first one returned, explicitly
	system, err := C.FindSystem(allSystems[0].ID, "", "")
	assert.Nil(t, err)
	assert.Equal(t, allSystems[0].ID, system.System.ID)
}

// TestGetSingleSystemByIDInvalid will search the for a system that will not be found
func TestGetSingleSystemByIDInvalid(t *testing.T) {
	// then try to get the first one returned, explicitly
	system, err := C.FindSystem(invalidIdentifier, "", "")
	fmt.Printf("system %v err %s\n", system, err)
	assert.NotNil(t, err)
	assert.Nil(t, system)
}

// TestGetSingleSystemByNameInvalid will search the for a system that will not be found
func TestGetSingleSystemByIDName(t *testing.T) {
	// then try to get the first one returned, explicitly
	system, err := C.FindSystem("", invalidIdentifier, "")
	assert.NotNil(t, err)
	assert.Nil(t, system)
}

// TestGetSystemStatistics will return System statistics
func TestGetSystemStatistics(t *testing.T) {
	// first, get all of the systems
	allSystems, err := C.GetSystems()
	assert.Nil(t, err)

	// then try to get the first one returned, explicitly
	system, err := C.FindSystem(allSystems[0].ID, "", "")
	assert.Nil(t, err)
	assert.Equal(t, allSystems[0].ID, system.System.ID)

	stats, err := system.GetStatistics()
	assert.Nil(t, err)
	assert.NotNil(t, stats)
}
