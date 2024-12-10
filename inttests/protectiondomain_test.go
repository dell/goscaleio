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
	"os"
	"testing"

	"github.com/dell/goscaleio"
	types "github.com/dell/goscaleio/types/v1"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

// getProtectionctx, returns GOSCALEIO_PROTECTIONDOMAIN, if set
// if not set, returns the first protection domain found
func getProtectionDomainName(t *testing.T) string {
	if os.Getenv("GOSCALEIO_PROTECTIONDOMAIN") != "" {
		return os.Getenv("GOSCALEIO_PROTECTIONDOMAIN")
	}
	system := getSystem()
	assert.NotNil(t, system)
	pd, _ := system.GetProtectionDomain(context.Background(), "")
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
	pd, err := system.FindProtectionDomain(context.Background(), "", name, "")
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
	pd, err := system.GetProtectionDomain(context.Background(), "")
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
		prot, err := system.FindProtectionDomain(context.Background(), "", pd.ProtectionDomain.Name, "")
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
		prot, err := system.FindProtectionDomain(context.Background(), pd.ProtectionDomain.ID, "", "")
		assert.Nil(t, err)
		assert.Equal(t, pd.ProtectionDomain.ID, prot.ID)
	}
}

// TestGetProtectionDomainByNameInvalid attempts to get a ProtectionDomain that does not exist
func TestGetProtectionDomainByNameInvalid(t *testing.T) {
	system := getSystem()
	assert.NotNil(t, system)

	pd, err := system.FindProtectionDomain(context.Background(), "", invalidIdentifier, "")
	assert.NotNil(t, err)
	assert.Nil(t, pd)
}

// TestGetProtectionDomainByIDInvalid attempts to get a ProtectionDomain that does not exist
func TestGetProtectionDomainByIDInvalid(t *testing.T) {
	system := getSystem()
	assert.NotNil(t, system)

	pd, err := system.FindProtectionDomain(context.Background(), invalidIdentifier, "", "")
	assert.NotNil(t, err)
	assert.Nil(t, pd)
}

func TestCreateDeleteProtectionDomain(t *testing.T) {
	system := getSystem()
	assert.NotNil(t, system)

	domainName := fmt.Sprintf("%s-%s", testPrefix, "Domain")
	ctx := context.Background()

	// create the pool
	domainID, err := system.CreateProtectionDomain(ctx, domainName)
	assert.Nil(t, err)
	assert.NotNil(t, domainID)

	// try to create a pool that exists
	domainID, err = system.CreateProtectionDomain(ctx, domainName)
	assert.NotNil(t, err)
	assert.Equal(t, "", domainID)

	// delete the pool
	err = system.DeleteProtectionDomain(ctx, domainName)
	assert.Nil(t, err)

	// try to delete non-existent storage pool
	// delete the pool
	err = system.DeleteProtectionDomain(ctx, domainName)
	assert.NotNil(t, err)
}

func TestCRUDProtectionDomain(t *testing.T) {
	system := getSystem()
	assert.NotNil(t, system)

	domainName := fmt.Sprintf("%s-%s", testPrefix, "Domain")
	ctx := context.Background()

	// create the pd
	domainID, err := system.CreateProtectionDomain(ctx, domainName)
	assert.Nil(t, err)
	assert.NotNil(t, domainID)

	pd, err2 := system.GetProtectionDomainEx(ctx, domainID)
	assert.Nil(t, err2)

	// change name of pd
	newName := fmt.Sprintf("%s2-%s", testPrefix, "Domain")
	err = pd.SetName(ctx, newName)
	assert.Nil(t, err)

	// inactivate pd
	err = pd.InActivate(ctx, false)
	assert.Nil(t, err)
	err = pd.Refresh(ctx)
	assert.Nil(t, err)
	assert.Equal(t, pd.ProtectionDomain.Name, newName)
	assert.Equal(t, pd.ProtectionDomain.ProtectionDomainState, "Inactive")

	testRfCacheProtectionDomain(t, pd)

	testNwLimitsProtectionDomain(t, pd)

	// activate pd
	err = pd.Activate(ctx, true)
	assert.Nil(t, err)
	err = pd.Refresh(ctx)
	assert.Nil(t, err)
	assert.Equal(t, pd.ProtectionDomain.ProtectionDomainState, "Active")

	testFGLMCacheProtectionDomain(t, pd)

	// check that finding pd by name yields same struct as refreshing
	pdByName, err3 := system.FindProtectionDomainByName(ctx, newName)
	assert.Nil(t, err3)
	assert.Equal(t, pd.ProtectionDomain, pdByName)

	// delete pd
	err = pd.Delete(ctx)
	assert.Nil(t, err)
}

func testFGLMCacheProtectionDomain(t *testing.T, pd *goscaleio.ProtectionDomain) {
	ctx := context.Background()
	err := pd.DisableFGLMcache(ctx)
	assert.Nil(t, err)
	err = pd.Refresh(ctx)
	assert.Nil(t, err)
	assert.Equal(t, pd.ProtectionDomain.FglMetadataCacheEnabled, false)

	err = pd.SetDefaultFGLMcacheSize(ctx, 1024)
	assert.Nil(t, err)
	err = pd.EnableFGLMcache(ctx)
	assert.Nil(t, err)
	err = pd.Refresh(ctx)
	assert.Nil(t, err)
	assert.Equal(t, pd.ProtectionDomain.FglMetadataCacheEnabled, true)
	assert.Equal(t, pd.ProtectionDomain.FglDefaultMetadataCacheSize, 1024)
}

func testRfCacheProtectionDomain(t *testing.T, pd *goscaleio.ProtectionDomain) {
	ctx := context.Background()
	err := pd.DisableRfcache(ctx)
	assert.Nil(t, err)
	err = pd.Refresh(ctx)
	assert.Nil(t, err)

	p := types.PDRCModeRead
	err = pd.SetRfcacheParams(ctx, types.PDRfCacheParams{
		RfCacheOperationalMode: p,
		RfCachePageSizeKb:      16,
		RfCacheMaxIoSizeKb:     64,
	})
	assert.Nil(t, err)
	err = pd.Refresh(ctx)
	assert.Nil(t, err)
	assert.Equal(t, pd.ProtectionDomain.RfCacheEnabled, false)
	assert.Equal(t, pd.ProtectionDomain.RfCacheOperationalMode, p)
	assert.Equal(t, pd.ProtectionDomain.RfCachePageSizeKb, 16)
	assert.Equal(t, pd.ProtectionDomain.RfCacheMaxIoSizeKb, 64)

	err = pd.SetRfcacheParams(ctx, types.PDRfCacheParams{
		RfCacheOperationalMode: p,
		RfCachePageSizeKb:      15,
		RfCacheMaxIoSizeKb:     65,
	})
	assert.NotNil(t, err)
	err = pd.Refresh(ctx)
	assert.Nil(t, err)
	assert.Equal(t, pd.ProtectionDomain.RfCacheEnabled, false)
	assert.Equal(t, pd.ProtectionDomain.RfCacheOperationalMode, p)
	assert.Equal(t, pd.ProtectionDomain.RfCachePageSizeKb, 16)
	assert.Equal(t, pd.ProtectionDomain.RfCacheMaxIoSizeKb, 64)

	err = pd.EnableRfcache(ctx)
	assert.Nil(t, err)
	err = pd.SetRfcacheParams(ctx, types.PDRfCacheParams{RfCachePageSizeKb: 4, RfCacheMaxIoSizeKb: 0})
	assert.Nil(t, err)
	err = pd.Refresh(ctx)
	assert.Nil(t, err)
	assert.Equal(t, pd.ProtectionDomain.RfCacheEnabled, true)
	assert.Equal(t, pd.ProtectionDomain.RfCacheOperationalMode, p)
	assert.Equal(t, pd.ProtectionDomain.RfCachePageSizeKb, 4)
	assert.Equal(t, pd.ProtectionDomain.RfCacheMaxIoSizeKb, 64)

	err = pd.SetRfcacheParams(ctx, types.PDRfCacheParams{RfCachePageSizeKb: 8, RfCacheMaxIoSizeKb: 32})
	assert.Nil(t, err)
	err = pd.Refresh(ctx)
	assert.Nil(t, err)
	assert.Equal(t, pd.ProtectionDomain.RfCacheEnabled, true)
	assert.Equal(t, pd.ProtectionDomain.RfCacheOperationalMode, p)
	assert.Equal(t, pd.ProtectionDomain.RfCachePageSizeKb, 8)
	assert.Equal(t, pd.ProtectionDomain.RfCacheMaxIoSizeKb, 32)

	err = pd.SetRfcacheParams(ctx, types.PDRfCacheParams{RfCachePageSizeKb: 16, RfCacheMaxIoSizeKb: 32})
	assert.Nil(t, err)
	err = pd.Refresh(ctx)
	assert.Nil(t, err)
	assert.Equal(t, pd.ProtectionDomain.RfCacheEnabled, true)
	assert.Equal(t, pd.ProtectionDomain.RfCacheOperationalMode, p)
	assert.Equal(t, pd.ProtectionDomain.RfCachePageSizeKb, 16)
	assert.Equal(t, pd.ProtectionDomain.RfCacheMaxIoSizeKb, 32)
}

func testNwLimitsProtectionDomain(t *testing.T, pd *goscaleio.ProtectionDomain) {
	oldPd := *pd.ProtectionDomain
	ctx := context.Background()
	a, b, c := 10*1024, 16*1024, 0
	err := pd.SetSdsNetworkLimits(ctx, types.SdsNetworkLimitParams{
		VTreeMigrationNetworkThrottlingInKbps:           &a,
		ProtectedMaintenanceModeNetworkThrottlingInKbps: &b,
		OverallIoNetworkThrottlingInKbps:                &c,
	})
	assert.Nil(t, err)
	err = pd.Refresh(ctx)
	assert.Nil(t, err)
	assert.Equal(t, pd.ProtectionDomain.RebuildNetworkThrottlingInKbps, oldPd.RebuildNetworkThrottlingInKbps)
	assert.Equal(t, pd.ProtectionDomain.RebalanceNetworkThrottlingInKbps, oldPd.RebalanceNetworkThrottlingInKbps)
	assert.Equal(t, pd.ProtectionDomain.VTreeMigrationNetworkThrottlingInKbps, a)
	assert.Equal(t, pd.ProtectionDomain.ProtectedMaintenanceModeNetworkThrottlingInKbps, b)
	assert.Equal(t, pd.ProtectionDomain.OverallIoNetworkThrottlingInKbps, c)

	a1, c1 := 64*1024, 100*1024
	err = pd.SetSdsNetworkLimits(ctx, types.SdsNetworkLimitParams{
		RebuildNetworkThrottlingInKbps:        &a1,
		RebalanceNetworkThrottlingInKbps:      &a1,
		VTreeMigrationNetworkThrottlingInKbps: &a1,
		OverallIoNetworkThrottlingInKbps:      &c1,
	})
	assert.Nil(t, err)
	err = pd.Refresh(ctx)
	assert.Nil(t, err)
	assert.Equal(t, pd.ProtectionDomain.RebuildNetworkThrottlingInKbps, a1)
	assert.Equal(t, pd.ProtectionDomain.RebalanceNetworkThrottlingInKbps, a1)
	assert.Equal(t, pd.ProtectionDomain.VTreeMigrationNetworkThrottlingInKbps, a1)
	assert.Equal(t, pd.ProtectionDomain.ProtectedMaintenanceModeNetworkThrottlingInKbps, b)
	assert.Equal(t, pd.ProtectionDomain.OverallIoNetworkThrottlingInKbps, c1)
}
