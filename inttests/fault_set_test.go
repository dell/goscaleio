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
	"time"

	types "github.com/dell/goscaleio/types/v1"
	"github.com/stretchr/testify/assert"
)

func TestCreateModifyDeleteFaultSet(t *testing.T) {
	domain := getProtectionDomain(t)
	assert.NotNil(t, domain)
	fsName := fmt.Sprintf("%s-%s", testPrefix, "FaultSet")

	fs := &types.FaultSetParam{
		Name:      fsName,
		ProtectionDomainID: domain.ProtectionDomain.ID,
	}

	// create the fault set
	fsID, err := domain.CreateFaultSet(fs)
	assert.Nil(t, err)
	assert.NotNil(t, fsID)

	// create a fault set that exists
	fsID2, err2 := domain.CreateFaultSet(fs)
	assert.NotNil(t, err2)
	assert.Equal(t, "", fsID2)

	// modify fault set name
	err = domain.ModifyFaultSetName(fsID, "faultSetRenamed")
	assert.Nil(t, err)
	
	// modify fault set performance profile
	err = domain.ModifyFaultSetPerFrofile(fsID, "Compact")
	assert.Nil(t, err)
	time.Sleep(5 * time.Second)

	// read the fault set
	fsr, err := domain.ReadFaultSet(fsID)
	assert.Equal(t, "faultSetRenamed", fsr.Name)

	// delete the fault set
	err = domain.DeleteFaultSet(fsID)
	assert.Nil(t, err)

	// try to delete non-existent fault set
	err3 := domain.DeleteFaultSet(invalidIdentifier)
	assert.NotNil(t, err3)
}

