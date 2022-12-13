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
	"github.com/stretchr/testify/assert"
)

// getStoragePoolName returns GOSCALEIO_STORAGEPOOL, if set
// if not set, returns the first storage pool found
func getStoragePoolName(t *testing.T) string {
	if os.Getenv("GOSCALEIO_STORAGEPOOL") != "" {
		return os.Getenv("GOSCALEIO_STORAGEPOOL")
	}

	system := getSystem()
	assert.NotNil(t, system)

	pd := getProtectionDomain(t)
	assert.NotNil(t, pd)

	if pd == nil {
		return ""
	}

	pools, err := pd.GetStoragePool("")
	assert.Nil(t, err)
	assert.NotZero(t, len(pools))
	if pools == nil {
		return ""
	}
	return pools[0].Name
}

// getStoragePool returns the StoragePool with the name retured by getStoragePool
func getStoragePool(t *testing.T) *goscaleio.StoragePool {
	pd := getProtectionDomain(t)
	assert.NotNil(t, pd)
	if pd == nil {
		return nil
	}

	name := getStoragePoolName(t)
	assert.NotEqual(t, name, "")

	pool, err := pd.FindStoragePool("", name, "")
	assert.Nil(t, err)
	assert.NotNil(t, pool)
	if pool == nil {
		return nil
	}

	// create a StoragePool instance to return
	outPool := goscaleio.NewStoragePoolEx(C, pool)

	// creare a storagePool via NewStoragePool to test
	tempPool := goscaleio.NewStoragePool(C)
	tempPool.StoragePool = pool
	assert.Equal(t, outPool.StoragePool.ID, tempPool.StoragePool.ID)

	return outPool
}

// TestGetStoragePools will return all storage pools
func TestGetStoragePools(t *testing.T) {
	system := getSystem()
	assert.NotNil(t, system)

	pd := getProtectionDomain(t)
	assert.NotNil(t, pd)

	if pd != nil {
		pools, err := pd.GetStoragePool("")
		assert.Nil(t, err)
		assert.NotZero(t, len(pools))
	}
}

// TestGetStoragePoolByName gets a single specific StoragePool by Name
func TestGetStoragePoolByName(t *testing.T) {
	pd := getProtectionDomain(t)
	assert.NotNil(t, pd)

	pool := getStoragePool(t)
	assert.NotNil(t, pool)

	if pd != nil && pool != nil {
		foundPool, err := pd.FindStoragePool("", pool.StoragePool.Name, "")
		assert.Nil(t, err)
		assert.Equal(t, foundPool.Name, pool.StoragePool.Name)
	}
}

// TestGetStoragePoolByID gets a single specific StoragePool by ID
func TestGetStoragePoolByID(t *testing.T) {
	pd := getProtectionDomain(t)
	assert.NotNil(t, pd)

	pool := getStoragePool(t)
	assert.NotNil(t, pool)

	if pd != nil && pool != nil {
		foundPool, err := pd.FindStoragePool(pool.StoragePool.ID, "", "")
		assert.Nil(t, err)
		assert.Equal(t, foundPool.ID, pool.StoragePool.ID)
	}
}

// TestGetStoragePoolByNameInvalid attempts to get a StoragePool that does not exist
func TestGetStoragePoolByNameInvalid(t *testing.T) {
	pd := getProtectionDomain(t)
	assert.NotNil(t, pd)

	pool, err := pd.FindStoragePool("", invalidIdentifier, "")
	assert.NotNil(t, err)
	assert.Nil(t, pool)
}

// TestGetStoragePoolByIDInvalid attempts to get a StoragePool that does not exist
func TestGetStoragePoolByIDInvalid(t *testing.T) {
	pd := getProtectionDomain(t)
	assert.NotNil(t, pd)

	pool, err := pd.FindStoragePool(invalidIdentifier, "", "")
	assert.NotNil(t, err)
	assert.Nil(t, pool)
}

// TestGetStoragePoolStatistics
func TestGetStoragePoolStatistics(t *testing.T) {
	pool := getStoragePool(t)
	assert.NotNil(t, pool)

	stats, err := pool.GetStatistics()
	assert.Nil(t, err)
	assert.NotNil(t, stats)
}

func TestGetInstanceStoragePool(t *testing.T) {
	name := getStoragePoolName(t)
	assert.NotNil(t, name)

	// Find by name
	pool, err := C.FindStoragePool("", name, "", "")
	assert.Nil(t, err)
	assert.NotNil(t, pool)

	// Find by ID
	pool, err = C.FindStoragePool(pool.ID, "", "", "")
	assert.Nil(t, err)
	assert.NotNil(t, pool)

	// Find by href
	href := fmt.Sprintf("/api/instances/StoragePool::%s", pool.ID)
	pool, err = C.FindStoragePool("", "", href, "")
	assert.Nil(t, err)
	assert.NotNil(t, pool)

	// Find with invalid name
	pool, err = C.FindStoragePool("", invalidIdentifier, "", "")
	assert.NotNil(t, err)
	assert.Nil(t, pool)

	// Find with invalid ID
	pool, err = C.FindStoragePool(invalidIdentifier, "", "", "")
	assert.NotNil(t, err)
	assert.Nil(t, pool)

	// Find with invalid href
	href = fmt.Sprintf("/api/badurl/willnotwork")
	pool, err = C.FindStoragePool("", "", href, "")
	assert.NotNil(t, err)
	assert.Nil(t, pool)

	//Find with name and Protection Domain ID
	pd := getProtectionDomain(t)
	assert.NotNil(t, pd)

	pool, err = C.FindStoragePool("", name, "", pd.ProtectionDomain.ID)
	assert.Nil(t, err)
	assert.NotNil(t, pool)

}

func TestCreateDeleteStoragePool(t *testing.T) {
	domain := getProtectionDomain(t)
	assert.NotNil(t, domain)

	poolName := fmt.Sprintf("%s-%s", testPrefix, "StoragePool")

	// create the pool
	poolID, err := domain.CreateStoragePool(poolName, "")
	assert.Nil(t, err)
	assert.NotNil(t, poolID)

	// try to create a pool that exists
	poolID, err = domain.CreateStoragePool(poolName, "")
	assert.NotNil(t, err)
	assert.Equal(t, "", poolID)

	// delete the pool
	err = domain.DeleteStoragePool(poolName)
	assert.Nil(t, err)

	// try to dleet non-existent storage pool
	// delete the pool
	err = domain.DeleteStoragePool(invalidIdentifier)
	assert.NotNil(t, err)

}

// TestGetSDSStoragePool gets the SDS instances associated with storage pool
func TestGetSDSStoragePool(t *testing.T) {
	pool := getStoragePool(t)
	assert.NotNil(t, pool)

	stats, err := pool.GetSDSStoragePool()
	assert.Nil(t, err)
	assert.NotNil(t, stats)
}

// Modify TestModifyStoragePoolName
func TestModifyStoragePoolName(t *testing.T) {
	domain := getProtectionDomain(t)
	assert.NotNil(t, domain)
	_, err := domain.ModifyStoragePoolName("Invalid", "STPnew")
	assert.NotNil(t, err)
}

// Modify TestStoragePoolMediaType
func TestStoragePoolMediaType(t *testing.T) {
	domain := getProtectionDomain(t)
	assert.NotNil(t, domain)
	_, err := domain.ModifyStoragePoolMedia("b9b0be6600000004", "SSD")
	assert.Nil(t, err)
}

// Modify TestEnableRFCache
func TestEnableRFCache(t *testing.T) {
	domain := getProtectionDomain(t)
	assert.NotNil(t, domain)
	_, err := domain.EnableRFCache("b9b0be6400000003")
	assert.Nil(t, err)
}

// Modify TestDisableRFCache
func TestDisableRFCache(t *testing.T) {
	domain := getProtectionDomain(t)
	assert.NotNil(t, domain)
	_, err := domain.DisableRFCache("b9b0be6400000003")
	assert.Nil(t, err)
}

// Set TestSetRmcache
func TestSetRmcache(t *testing.T) {
	pd := getProtectionDomain(t)
	name := getStoragePoolName(t)

	pool, _ := pd.FindStoragePool("", name, "")

	// create a StoragePool instance to return
	domain := goscaleio.NewStoragePoolEx(C, pool)

	// create a storagePool via NewStoragePool to test
	tempPool := goscaleio.NewStoragePool(C)
	tempPool.StoragePool = pool

	_ = domain.ModifyRMCache("true")
}
