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
	"testing"

	"github.com/AnshumanPradipPatil1506/goscaleio"
	"github.com/stretchr/testify/assert"
)

// getAllSds will return all SDS instances
func getAllSds(t *testing.T) []*goscaleio.Sds {
	system := getSystem()
	assert.NotNil(t, system)
	if system == nil {
		return nil
	}

	pd := getProtectionDomain(t)
	assert.NotNil(t, pd)
	if pd == nil {
		return nil
	}

	if pd != nil {
		var allSds []*goscaleio.Sds
		sds, err := pd.GetSds()
		assert.Nil(t, err)
		assert.NotZero(t, len(sds))
		for _, s := range sds {
			// create an SDS via NewSdsEx to the caller (appending to the allSds slice)
			outSDS := goscaleio.NewSdsEx(C, &s)
			allSds = append(allSds, outSDS)
			// create an SDS via NewSds that we will through away
			tempSDS := goscaleio.NewSds(C)
			tempSDS.Sds = &s
			assert.Equal(t, outSDS.Sds.Name, tempSDS.Sds.Name)
		}
		return allSds
	}

	return nil
}

// TestGetSDSs will return all SDS instances
func TestGetSDSs(t *testing.T) {
	getAllSds(t)
}

// TestGetSDSByAttribute gets a single specific SDS by attribute
func TestGetSDSByAttribute(t *testing.T) {
	pd := getProtectionDomain(t)
	assert.NotNil(t, pd)
	if pd == nil {
		return
	}

	sds := getAllSds(t)
	assert.NotNil(t, sds)
	if sds == nil {
		return
	}

	found, err := pd.FindSds("Name", sds[0].Sds.Name)
	assert.Nil(t, err)
	assert.NotNil(t, found)
	assert.Equal(t, sds[0].Sds.Name, found.Name)

	found, err = pd.FindSds("ID", sds[0].Sds.ID)
	assert.Nil(t, err)
	assert.NotNil(t, found)
	assert.Equal(t, sds[0].Sds.ID, found.ID)
}

// TestGetSDSByAttributeInvalid fails to get a single specific SDS by attribute
func TestGetSDSByAttributeInvalid(t *testing.T) {
	pd := getProtectionDomain(t)
	assert.NotNil(t, pd)
	if pd == nil {
		return
	}

	sds := getAllSds(t)
	assert.NotNil(t, sds)
	if sds == nil {
		return
	}

	found, err := pd.FindSds("Name", invalidIdentifier)
	assert.NotNil(t, err)
	assert.Nil(t, found)

	found, err = pd.FindSds("ID", invalidIdentifier)
	assert.NotNil(t, err)
	assert.Nil(t, found)
}

// TestCreateSdsInvalid will attempt to add an SDS, which results in failure
func TestCreateSdsInvalid(t *testing.T) {
	pd := getProtectionDomain(t)
	assert.NotNil(t, pd)

	// attempt to create an SDS with a number of invalid IPs
	// this is done, in a failure mode, to prevent changing the Protection Domain used for testing
	sdsName := "invalid"
	sdsIPList := []string{"0.1.1.1", "0.2.2.2"}
	sdsID, err := pd.CreateSds(sdsName, sdsIPList)
	assert.NotNil(t, err)
	assert.Equal(t, "", sdsID)

}
