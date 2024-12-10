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

	pools, err := pd.GetStoragePool(context.Background(), "")
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

	pool, err := pd.FindStoragePool(context.Background(), "", name, "")
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
		pools, err := pd.GetStoragePool(context.Background(), "")
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
		foundPool, err := pd.FindStoragePool(context.Background(), "", pool.StoragePool.Name, "")
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
		foundPool, err := pd.FindStoragePool(context.Background(), pool.StoragePool.ID, "", "")
		assert.Nil(t, err)
		assert.Equal(t, foundPool.ID, pool.StoragePool.ID)
	}
}

// TestGetStoragePoolByNameInvalid attempts to get a StoragePool that does not exist
func TestGetStoragePoolByNameInvalid(t *testing.T) {
	pd := getProtectionDomain(t)
	assert.NotNil(t, pd)

	pool, err := pd.FindStoragePool(context.Background(), "", invalidIdentifier, "")
	assert.NotNil(t, err)
	assert.Nil(t, pool)
}

// TestGetStoragePoolByIDInvalid attempts to get a StoragePool that does not exist
func TestGetStoragePoolByIDInvalid(t *testing.T) {
	pd := getProtectionDomain(t)
	assert.NotNil(t, pd)

	pool, err := pd.FindStoragePool(context.Background(), invalidIdentifier, "", "")
	assert.NotNil(t, err)
	assert.Nil(t, pool)
}

// TestGetStoragePoolStatistics
func TestGetStoragePoolStatistics(t *testing.T) {
	pool := getStoragePool(t)
	assert.NotNil(t, pool)

	stats, err := pool.GetStatistics(context.Background())
	assert.Nil(t, err)
	assert.NotNil(t, stats)
}

func TestGetInstanceStoragePool(t *testing.T) {
	ctx := context.Background()
	name := getStoragePoolName(t)
	assert.NotNil(t, name)

	// Find by name
	pool, err := C.FindStoragePool(ctx, "", name, "", "")
	assert.Nil(t, err)
	assert.NotNil(t, pool)

	// Find by ID
	pool, err = C.FindStoragePool(ctx, pool.ID, "", "", "")
	assert.Nil(t, err)
	assert.NotNil(t, pool)

	// Find by href
	href := fmt.Sprintf("/api/instances/StoragePool::%s", pool.ID)
	pool, err = C.FindStoragePool(ctx, "", "", href, "")
	assert.Nil(t, err)
	assert.NotNil(t, pool)

	// Find with invalid name
	pool, err = C.FindStoragePool(ctx, "", invalidIdentifier, "", "")
	assert.NotNil(t, err)
	assert.Nil(t, pool)

	// Find with invalid ID
	pool, err = C.FindStoragePool(ctx, invalidIdentifier, "", "", "")
	assert.NotNil(t, err)
	assert.Nil(t, pool)

	// Find with invalid href
	href = fmt.Sprintf("/api/badurl/willnotwork")
	pool, err = C.FindStoragePool(ctx, "", "", href, "")
	assert.NotNil(t, err)
	assert.Nil(t, pool)

	// Find with name and Protection Domain ID
	pd := getProtectionDomain(t)
	assert.NotNil(t, pd)

	pool, err = C.FindStoragePool(ctx, "", name, "", pd.ProtectionDomain.ID)
	assert.Nil(t, err)
	assert.NotNil(t, pool)
}

func TestCreateDeleteStoragePool(t *testing.T) {
	ctx := context.Background()
	domain := getProtectionDomain(t)
	assert.NotNil(t, domain)

	poolName := fmt.Sprintf("%s-%s", testPrefix, "StoragePool")

	sp := &types.StoragePoolParam{
		Name:      poolName,
		MediaType: "HDD",
	}

	// create the pool
	poolID, err := domain.CreateStoragePool(ctx, sp)
	assert.Nil(t, err)
	assert.NotNil(t, poolID)

	// try to create a pool that exists
	poolID, err = domain.CreateStoragePool(ctx, sp)
	assert.NotNil(t, err)
	assert.Equal(t, "", poolID)

	// delete the pool
	err = domain.DeleteStoragePool(ctx, poolName)
	assert.Nil(t, err)

	// try to dleet non-existent storage pool
	// delete the pool
	err = domain.DeleteStoragePool(ctx, invalidIdentifier)
	assert.NotNil(t, err)
}

// TestGetSDSStoragePool gets the SDS instances associated with storage pool
func TestGetSDSStoragePool(t *testing.T) {
	pool := getStoragePool(t)
	assert.NotNil(t, pool)

	stats, err := pool.GetSDSStoragePool(context.Background())
	assert.Nil(t, err)
	assert.NotNil(t, stats)
}

// Modify TestModifyStoragePoolName
func TestModifyStoragePoolName(t *testing.T) {
	domain := getProtectionDomain(t)
	assert.NotNil(t, domain)
	_, err := domain.ModifyStoragePoolName(context.Background(), "Invalid", "STPnew")
	assert.NotNil(t, err)
}

// Modify TestStoragePoolMediaType
func TestStoragePoolMediaType(t *testing.T) {
	ctx := context.Background()
	domain := getProtectionDomain(t)
	assert.NotNil(t, domain)

	poolName := fmt.Sprintf("%s-%s", testPrefix, "StoragePool")

	sp := &types.StoragePoolParam{
		Name:      poolName,
		MediaType: "HDD",
	}

	// create the storage pool
	poolID, err := domain.CreateStoragePool(ctx, sp)
	assert.Nil(t, err)
	assert.NotNil(t, poolID)
	_, err = domain.ModifyStoragePoolMedia(ctx, poolID, "SSD")
	assert.Nil(t, err)

	// delete the pool
	err = domain.DeleteStoragePool(ctx, poolName)
	assert.Nil(t, err)
}

// Modify TestEnableRFCache
func TestEnableRFCache(t *testing.T) {
	ctx := context.Background()
	domain := getProtectionDomain(t)
	assert.NotNil(t, domain)

	poolName := fmt.Sprintf("%s-%s", testPrefix, "StoragePool")

	sp := &types.StoragePoolParam{
		Name:      poolName,
		MediaType: "HDD",
	}

	// create the storage pool
	poolID, err := domain.CreateStoragePool(ctx, sp)
	assert.Nil(t, err)
	assert.NotNil(t, poolID)
	_, err = domain.EnableRFCache(ctx, poolID)
	assert.Nil(t, err)
	// delete the pool
	err = domain.DeleteStoragePool(ctx, poolName)
	assert.Nil(t, err)
}

// Test all the additional functionality for a storage pool
func TestStoragePoolAdditionalFunctionality(t *testing.T) {
	ctx := context.Background()
	// get the protection domain
	domain := getProtectionDomain(t)
	assert.NotNil(t, domain)

	poolName := fmt.Sprintf("%s-%s", testPrefix, "StoragePool")

	sp := &types.StoragePoolParam{
		Name:      poolName,
		MediaType: "HDD",
	}

	// create the storage pool
	poolID, err := domain.CreateStoragePool(ctx, sp)
	assert.Nil(t, err)
	assert.NotNil(t, poolID)

	// disable the padding
	err = domain.EnableOrDisableZeroPadding(ctx, poolID, "false")
	assert.Nil(t, err)
	pool, _ := domain.FindStoragePool(ctx, poolID, "", "")
	// check the value
	assert.Equal(t, pool.ZeroPaddingEnabled, false)

	// Now enable the padding
	err = domain.EnableOrDisableZeroPadding(ctx, poolID, "true")
	assert.Nil(t, err)
	pool, _ = domain.FindStoragePool(ctx, poolID, "", "")
	// check the value
	assert.Equal(t, pool.ZeroPaddingEnabled, true)

	// Modify Replication Journal Capacity to make it 36
	err = domain.SetReplicationJournalCapacity(ctx, poolID, "36")
	assert.Nil(t, err)
	pool, _ = domain.FindStoragePool(ctx, poolID, "", "")
	// check the value
	assert.Equal(t, pool.ReplicationCapacityMaxRatio, 36)

	// Again Modify Replication Journal Capacity to make it 0 else storage pool can't be deleted
	err = domain.SetReplicationJournalCapacity(ctx, poolID, "0")
	assert.Nil(t, err)
	pool, _ = domain.FindStoragePool(ctx, poolID, "", "")
	// again check the value
	assert.Equal(t, pool.ReplicationCapacityMaxRatio, 0)

	// set the capacity threshold for the storage pool
	capacityAlertThreshold := &types.CapacityAlertThresholdParam{
		CapacityAlertHighThresholdPercent: "68",
	}
	err = domain.SetCapacityAlertThreshold(ctx, poolID, capacityAlertThreshold)
	assert.Nil(t, err)
	pool, _ = domain.FindStoragePool(ctx, poolID, "", "")
	// check the value
	assert.Equal(t, pool.CapacityAlertHighThreshold, 68)

	// Set the protected maintenance mode
	protectedMaintenanceModeParam := &types.ProtectedMaintenanceModeParam{
		Policy:                      "favorAppIos",
		NumOfConcurrentIosPerDevice: "18",
	}
	err = domain.SetProtectedMaintenanceModeIoPriorityPolicy(ctx, poolID, protectedMaintenanceModeParam)
	assert.Nil(t, err)
	pool, _ = domain.FindStoragePool(ctx, poolID, "", "")
	// check the value
	assert.Equal(t, pool.ProtectedMaintenanceModeIoPriorityPolicy, "favorAppIos")
	assert.Equal(t, pool.ProtectedMaintenanceModeIoPriorityNumOfConcurrentIosPerDevice, 18)

	// set rebalance enablement value
	err = domain.SetRebalanceEnabled(ctx, poolID, "true")
	assert.Nil(t, err)
	pool, _ = domain.FindStoragePool(ctx, poolID, "", "")
	// check the value
	assert.Equal(t, pool.RebalanceEnabled, true)

	// Again set rebalance enablement value
	err = domain.SetRebalanceEnabled(ctx, poolID, "false")
	assert.Nil(t, err)
	pool, _ = domain.FindStoragePool(ctx, poolID, "", "")
	// check the value
	assert.Equal(t, pool.RebalanceEnabled, false)

	// set the rebalance IO priority policy for the storage pool
	protectedMaintenanceModeParam = &types.ProtectedMaintenanceModeParam{
		Policy:                      "limitNumOfConcurrentIos",
		NumOfConcurrentIosPerDevice: "13",
	}
	err = domain.SetRebalanceIoPriorityPolicy(ctx, poolID, protectedMaintenanceModeParam)
	assert.Nil(t, err)
	// check the value
	pool, _ = domain.FindStoragePool(ctx, poolID, "", "")
	assert.Equal(t, pool.RebalanceioPriorityPolicy, "limitNumOfConcurrentIos")
	assert.Equal(t, pool.RebalanceioPriorityNumOfConcurrentIosPerDevice, 13)
	assert.Nil(t, err)

	// Set vtree migration IO priority policy
	protectedMaintenanceModeParam = &types.ProtectedMaintenanceModeParam{
		Policy:                      "favorAppIos",
		NumOfConcurrentIosPerDevice: "12",
		BwLimitPerDeviceInKbps:      "1030",
	}
	err = domain.SetVTreeMigrationIOPriorityPolicy(ctx, poolID, protectedMaintenanceModeParam)
	assert.Nil(t, err)
	// check the value
	pool, _ = domain.FindStoragePool(ctx, poolID, "", "")
	assert.Equal(t, pool.VtreeMigrationIoPriorityPolicy, "favorAppIos")
	assert.Equal(t, pool.VtreeMigrationIoPriorityNumOfConcurrentIosPerDevice, 12)
	assert.Equal(t, pool.VtreeMigrationIoPriorityBwLimitPerDeviceInKbps, 1030)

	// set the spare percentage
	err = domain.SetSparePercentage(ctx, poolID, "67")
	assert.Nil(t, err)
	// check the value
	pool, _ = domain.FindStoragePool(ctx, poolID, "", "")
	assert.Equal(t, pool.SparePercentage, 67)

	// set the Rmcache write handling mode
	err = domain.SetRMcacheWriteHandlingMode(ctx, poolID, "Cached")
	assert.Nil(t, err)
	// check the value
	pool, _ = domain.FindStoragePool(ctx, poolID, "", "")
	assert.Equal(t, pool.RmCacheWriteHandlingMode, "Cached")

	// set the rebuild enablemenent value
	err = domain.SetRebuildEnabled(ctx, poolID, "false")
	assert.Nil(t, err)
	// check the value
	pool, _ = domain.FindStoragePool(ctx, poolID, "", "")
	assert.Equal(t, pool.RebuildEnabled, false)

	// set the number of parallel rebuild rebalance jobs per device
	err = domain.SetRebuildRebalanceParallelismParam(ctx, poolID, "9")
	assert.Nil(t, err)
	// check the value
	pool, _ = domain.FindStoragePool(ctx, poolID, "", "")
	assert.Equal(t, pool.NumofParallelRebuildRebalanceJobsPerDevice, 9)

	// enable fragmentation
	err = domain.Fragmentation(ctx, poolID, true)
	assert.Nil(t, err)
	pool, _ = domain.FindStoragePool(ctx, poolID, "", "")
	assert.Equal(t, pool.FragmentationEnabled, true)

	// disable fragmentation
	err = domain.Fragmentation(ctx, poolID, false)
	assert.Nil(t, err)
	pool, _ = domain.FindStoragePool(ctx, poolID, "", "")
	assert.Equal(t, pool.FragmentationEnabled, false)

	// finally after all the operations, now delete the pool
	err = domain.DeleteStoragePool(ctx, poolName)
	assert.Nil(t, err)
}

// Modify TestDisableRFCache
func TestDisableRFCache(t *testing.T) {
	ctx := context.Background()
	domain := getProtectionDomain(t)
	assert.NotNil(t, domain)

	poolName := fmt.Sprintf("%s-%s", testPrefix, "StoragePool")

	sp := &types.StoragePoolParam{
		Name:      poolName,
		MediaType: "HDD",
	}

	// create the storage pool
	poolID, err := domain.CreateStoragePool(ctx, sp)
	assert.Nil(t, err)
	assert.NotNil(t, poolID)
	_, err = domain.DisableRFCache(ctx, poolID)
	assert.Nil(t, err)
	// delete the pool
	err = domain.DeleteStoragePool(ctx, poolName)
	assert.Nil(t, err)
}

// Modify TestModifyRmCache
func TestModifyRmCache(t *testing.T) {
	ctx := context.Background()
	pd := getProtectionDomain(t)
	name := getStoragePoolName(t)

	pool, _ := pd.FindStoragePool(ctx, "", name, "")

	// create a StoragePool instance to return
	domain := goscaleio.NewStoragePoolEx(C, pool)

	// create a storagePool via NewStoragePool to test
	tempPool := goscaleio.NewStoragePool(C)
	tempPool.StoragePool = pool

	err := domain.ModifyRMCache(ctx, "true")
	assert.Nil(t, err)
}

// TestGetAllStoragePoolsApi gets all storage pools available on system
func TestGetAllStoragePoolsApi(t *testing.T) {
	// get system
	system := getSystem()
	assert.NotNil(t, system)

	// get all storagepools on the system
	storagepools, err := system.GetAllStoragePools(context.Background())
	assert.Nil(t, err)
	assert.NotNil(t, storagepools)
}

// TestGetStoragePoolByIDApi gets storage pool by ID
func TestGetStoragePoolByIDApi(t *testing.T) {
	ctx := context.Background()
	name := getStoragePoolName(t)
	assert.NotNil(t, name)

	// get system
	system := getSystem()
	assert.NotNil(t, system)

	// Find by name
	pool, err := C.FindStoragePool(ctx, "", name, "", "")
	assert.Nil(t, err)
	assert.NotNil(t, pool)

	// Find by ID
	pool1, err := system.GetStoragePoolByID(ctx, pool.ID)
	assert.Nil(t, err)
	assert.NotNil(t, pool1)

	// Find with invalid identifier
	pool, err = C.FindStoragePool(ctx, "", invalidIdentifier, "", "")
	assert.NotNil(t, err)
	assert.Nil(t, pool)

	// Find by ID
	pool1, err = system.GetStoragePoolByID(ctx, invalidIdentifier)
	assert.NotNil(t, err)
	assert.Nil(t, pool1)
}
