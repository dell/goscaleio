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
	"os"
	"testing"

	"github.com/dell/goscaleio"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

// getProtectionDomainName returns GOSCALEIO_PROTECTIONDOMAIN, if set
// if not set, returns the first protection domain found
func getProtectionDomainName(t *testing.T) string {
	if os.Getenv("GOSCALEIO_PROTECTIONDOMAIN") != "" {
		return os.Getenv("GOSCALEIO_PROTECTIONDOMAIN")
	}
	system := getSystem()
	assert.NotNil(t, system)
	pd, _ := system.GetProtectionDomain("")
	assert.NotNil(t, pd)
	if pd == nil {
		return ""
	}
	return pd[0].Name
}

// getProtectionDomain returns the ProtectionDomain with the name retured by getProtectionDomainName
func getProtectionDomain(t *testing.T) *goscaleio.ProtectionDomain {
	system := getSystem()
	assert.NotNil(t, system)

	name := getProtectionDomainName(t)
	assert.NotEqual(t, name, "")
	pd, err := system.FindProtectionDomain("", name, "")
	assert.Nil(t, err)
	assert.NotNil(t, pd)
	if pd == nil {
		return nil
	}

	outPD := goscaleio.NewProtectionDomain(C)
	outPD.ProtectionDomain = pd
	return outPD
}

// getAllProtectionDomains returns all ProtectionDomains found
func getAllProtectionDomains(t *testing.T) []*goscaleio.ProtectionDomain {
	system := getSystem()
	assert.NotNil(t, system)

	log.SetLevel(log.DebugLevel)
	pd, err := system.GetProtectionDomain("")
	assert.Nil(t, err)
	assert.NotZero(t, len(pd))
	log.SetLevel(log.InfoLevel)

	var allDomains []*goscaleio.ProtectionDomain

	for _, domain := range pd {
		// create the PD to return
		outPD := goscaleio.NewProtectionDomainEx(C, domain)
		allDomains = append(allDomains, outPD)
		// create another PD for testng purposes (via NewProtectionDomain)
		tempPD := goscaleio.NewProtectionDomain(C)
		tempPD.ProtectionDomain = domain
		assert.Equal(t, outPD.ProtectionDomain.ID, tempPD.ProtectionDomain.ID)
	}
	return allDomains
}

// TestGetProtectionDomains gets all protection domains
func TestGetProtectionDomains(t *testing.T) {
	domains := getAllProtectionDomains(t)
	assert.NotNil(t, domains)
	assert.NotZero(t, len(domains))
}

// TestGetProtectionDomainByName gets a single specific ProtectionDomain by Name
func TestGetProtectionDomainByName(t *testing.T) {
	system := getSystem()
	assert.NotNil(t, system)

	pd := getProtectionDomain(t)
	assert.NotNil(t, pd)

	if pd != nil {
		prot, err := system.FindProtectionDomain("", pd.ProtectionDomain.Name, "")
		assert.Nil(t, err)
		assert.Equal(t, pd.ProtectionDomain.Name, prot.Name)
	}
}

// TestGetProtectionDomainByID gets a single specific ProtectionDomain by ID
func TestGetProtectionDomainByID(t *testing.T) {
	system := getSystem()
	assert.NotNil(t, system)

	pd := getProtectionDomain(t)
	assert.NotNil(t, pd)

	if pd != nil {
		prot, err := system.FindProtectionDomain(pd.ProtectionDomain.ID, "", "")
		assert.Nil(t, err)
		assert.Equal(t, pd.ProtectionDomain.ID, prot.ID)
	}
}

// TestGetProtectionDomainByNameInvalid attempts to get a ProtectionDomain that does not exist
func TestGetProtectionDomainByNameInvalid(t *testing.T) {
	system := getSystem()
	assert.NotNil(t, system)

	pd, err := system.FindProtectionDomain("", invalidIdentifier, "")
	assert.NotNil(t, err)
	assert.Nil(t, pd)
}

// TestGetProtectionDomainByIDInvalid attempts to get a ProtectionDomain that does not exist
func TestGetProtectionDomainByIDInvalid(t *testing.T) {
	system := getSystem()
	assert.NotNil(t, system)

	pd, err := system.FindProtectionDomain(invalidIdentifier, "", "")
	assert.NotNil(t, err)
	assert.Nil(t, pd)
}

func TestCreateDeleteProtectionDomain(t *testing.T) {
	system := getSystem()
	assert.NotNil(t, system)

	domainName := fmt.Sprintf("%s-%s", testPrefix, "Domain")

	// create the pool
	domainID, err := system.CreateProtectionDomain(domainName)
	assert.Nil(t, err)
	assert.NotNil(t, domainID)

	// try to create a pool that exists
	domainID, err = system.CreateProtectionDomain(domainName)
	assert.NotNil(t, err)
	assert.Equal(t, "", domainID)

	// delete the pool
	err = system.DeleteProtectionDomain(domainName)
	assert.Nil(t, err)

	// try to delete non-existent storage pool
	// delete the pool
	err = system.DeleteProtectionDomain(domainName)
	assert.NotNil(t, err)

}
