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
	types "github.com/dell/goscaleio/types/v1"
	"github.com/stretchr/testify/assert"
)

// getAllDevices will return all Device instances
func getAllDevices(t *testing.T) []*goscaleio.Device {
	pool := getStoragePool(t)
	assert.NotNil(t, pool)
	if pool == nil {
		return nil
	}

	var allDevice []*goscaleio.Device
	devices, err := pool.GetDevice()
	assert.Nil(t, err)
	assert.NotZero(t, len(devices))
	for _, d := range devices {
		// create a device to return to the caller
		outDevice := goscaleio.NewDeviceEx(C, &d) // #nosec G601
		allDevice = append(allDevice, outDevice)
		// create a device via NewDevice for testing purposes
		tempDevice := goscaleio.NewDevice(C)
		tempDevice.Device = &d // #nosec G601
		assert.Equal(t, outDevice.Device.ID, tempDevice.Device.ID)
	}
	return allDevice
}

func getAllDevicesFromSystem(t *testing.T) []goscaleio.Device {
	system := getSystem()
	var allDevice []goscaleio.Device
	devices, err1 := system.GetAllDevice()
	assert.Nil(t, err1)
	assert.NotZero(t, len(devices))
	for _, d := range devices {
		// create a device to return to the caller
		outDevice := goscaleio.NewDeviceEx(C, &d) // #nosec G601
		allDevice = append(allDevice, *outDevice)
		// create a device via NewDevice for testing purposes
		tempDevice := goscaleio.NewDevice(C)
		tempDevice.Device = &d // #nosec G601
		assert.Equal(t, outDevice.Device.ID, tempDevice.Device.ID)
	}
	return allDevice
}

// getAllDevices will return all Device instances
func getAllSdsDevices(t *testing.T) []*goscaleio.Device {
	sds := getAllSds(t)[0]
	assert.NotNil(t, sds)
	if sds == nil {
		return nil
	}

	var allDevice []*goscaleio.Device
	devices, err := sds.GetDevice()
	assert.Nil(t, err)
	assert.NotZero(t, len(devices))
	for _, d := range devices {
		// create a device to return to the caller
		outDevice := goscaleio.NewDeviceEx(C, &d) // #nosec G601
		allDevice = append(allDevice, outDevice)
		// create a device via NewDevice for testing purposes
		tempDevice := goscaleio.NewDevice(C)
		tempDevice.Device = &d // #nosec G601
		assert.Equal(t, outDevice.Device.ID, tempDevice.Device.ID)
	}
	return allDevice
}

// TestGetDevices will return all Device instances
func TestGetDevices(t *testing.T) {
	getAllDevices(t)
	getAllDevicesFromSystem(t)
	TestGetDeviceByField(t)
}

// TestGetDeviceByAttribute gets a single specific Device by attribute
func TestGetDeviceByAttribute(t *testing.T) {
	pool := getStoragePool(t)
	assert.NotNil(t, pool)
	if pool == nil {
		return
	}

	devices := getAllDevices(t)
	assert.NotNil(t, devices)
	assert.NotZero(t, len(devices))
	if devices == nil {
		return
	}

	found, err := pool.FindDevice("Name", devices[0].Device.Name)
	assert.Nil(t, err)
	assert.NotNil(t, found)
	assert.Equal(t, devices[0].Device.Name, found.Name)

	found, err = pool.FindDevice("ID", devices[0].Device.ID)
	assert.Nil(t, err)
	assert.NotNil(t, found)
	assert.Equal(t, devices[0].Device.ID, found.ID)
}

// TestGetDeviceByAttributeInvalid fails to get a single specific Device by attribute
func TestGetDeviceByAttributeInvalid(t *testing.T) {
	pool := getStoragePool(t)
	assert.NotNil(t, pool)
	if pool == nil {
		return
	}

	devices := getAllDevices(t)
	assert.NotNil(t, devices)
	assert.NotZero(t, len(devices))
	if devices == nil {
		return
	}

	found, err := pool.FindDevice("Name", invalidIdentifier)
	assert.NotNil(t, err)
	assert.Nil(t, found)

	found, err = pool.FindDevice("ID", invalidIdentifier)
	assert.NotNil(t, err)
	assert.Nil(t, found)
}

// TestGetDeviceByAttribute gets a single specific Device by attribute
func TestGetDeviceByField(t *testing.T) {
	system := getSystem()
	devices, err1 := system.GetAllDevice()
	assert.NotNil(t, devices)
	assert.NotZero(t, len(devices))
	if devices == nil || err1 != nil {
		return
	}

	found, err := system.GetDeviceByField("Name", devices[0].Name)
	assert.Nil(t, err)
	assert.NotNil(t, found)
	assert.Equal(t, devices[0].Name, found[0].Name)

	found, err = system.GetDeviceByField("ID", devices[0].ID)
	assert.Nil(t, err)
	assert.NotNil(t, found)
	assert.Equal(t, devices[0].ID, found[0].ID)
}

// TestAddDeviceInvalid will attempt to add an invalid device to an invalid SDS
func TestAddDeviceInvalid(t *testing.T) {
	pool := getStoragePool(t)
	assert.NotNil(t, pool)

	dev := &types.DeviceParam{
		DeviceCurrentPathname: "/invalidPath/invalidDevice",
		SdsID:                 invalidIdentifier,
	}

	deviceID, err := pool.AttachDevice(dev)
	assert.NotNil(t, err)
	assert.Equal(t, "", deviceID)
}

func getSdsID() string {
	if os.Getenv("GOSCALEIO_SDSID") != "" {
		return os.Getenv("GOSCALEIO_SDSID")
	}

	return ""
}

// TestAddDeviceValid add/remove device to/from the storage pool
func TestAddDeviceValid(t *testing.T) {
	pool := getStoragePool(t)
	assert.NotNil(t, pool)

	sdsID := getSdsID()
	assert.NotEqual(t, sdsID, "")

	dev := &types.DeviceParam{
		DeviceCurrentPathname: "/dev/sdc",
		SdsID:                 sdsID,
	}

	deviceID, err := pool.AttachDevice(dev)
	assert.Nil(t, err)
	assert.NotNil(t, deviceID)

	err = pool.RemoveDevice(deviceID)
	assert.Nil(t, err)
}

// TestDeviceSetName sets the name of the device
func TestDeviceSetName(t *testing.T) {
	pool := getStoragePool(t)
	assert.NotNil(t, pool)

	sdsID := getSdsID()
	assert.NotEqual(t, sdsID, "")

	dev := &types.DeviceParam{
		DeviceCurrentPathname: "/dev/sdc",
		SdsID:                 sdsID,
	}

	deviceID, err := pool.AttachDevice(dev)
	assert.Nil(t, err)
	assert.NotNil(t, deviceID)

	err = pool.SetDeviceName(deviceID, "device_renamed")
	assert.Nil(t, err)

	system := getSystem()
	device, err1 := system.GetDevice(deviceID)
	assert.Nil(t, err1)
	assert.Equal(t, device.Name, "device_renamed")

	err = pool.RemoveDevice(deviceID)
	assert.Nil(t, err)
}

// TestDeviceMediaType modifies the media type of device
func TestDeviceMediaType(t *testing.T) {
	domain := getProtectionDomain(t)
	assert.NotNil(t, domain)

	poolName := fmt.Sprintf("%s-%s", testPrefix, "StoragePool")

	sp := &types.StoragePoolParam{
		Name:      poolName,
		MediaType: "HDD",
	}

	// create the pool
	poolID, err := domain.CreateStoragePool(sp)
	assert.Nil(t, err)
	assert.NotNil(t, poolID)

	poolID1, err1 := domain.ModifyStoragePoolMedia(poolID, "Transitional")
	assert.Nil(t, err1)
	assert.NotNil(t, poolID1)

	pool, err := domain.FindStoragePool(poolID, "", "")
	assert.Nil(t, err)
	assert.NotNil(t, pool)

	// create a StoragePool instance
	spInstance := goscaleio.NewStoragePoolEx(C, pool)

	sdsID := getSdsID()
	assert.NotEqual(t, sdsID, "")

	dev := &types.DeviceParam{
		DeviceCurrentPathname: "/dev/sdc",
		SdsID:                 sdsID,
		MediaType:             "HDD",
	}

	deviceID, err := spInstance.AttachDevice(dev)
	assert.Nil(t, err)
	assert.NotNil(t, deviceID)

	err = spInstance.SetDeviceMediaType(deviceID, "SSD")
	assert.Nil(t, err)

	system := getSystem()
	device, err1 := system.GetDevice(deviceID)
	assert.Nil(t, err1)
	assert.Equal(t, device.MediaType, "SSD")

	// Remove the device
	err = spInstance.RemoveDevice(deviceID)
	assert.Nil(t, err)

	// Delete the pool
	err = domain.DeleteStoragePool(poolName)
	assert.Nil(t, err)
}

// TestDeviceExternalAccelerationType modifies the device external acceleration type
func TestDeviceExternalAccelerationType(t *testing.T) {
	pool := getStoragePool(t)
	assert.NotNil(t, pool)

	sdsID := getSdsID()
	assert.NotEqual(t, sdsID, "")

	dev := &types.DeviceParam{
		DeviceCurrentPathname: "/dev/sdc",
		SdsID:                 sdsID,
	}

	deviceID, err := pool.AttachDevice(dev)
	assert.Nil(t, err)
	assert.NotNil(t, deviceID)

	err = pool.SetDeviceExternalAccelerationType(deviceID, "Read")
	assert.Nil(t, err)

	system := getSystem()
	device, err1 := system.GetDevice(deviceID)
	assert.Nil(t, err1)
	assert.Equal(t, device.ExternalAccelerationType, "Read")

	err = pool.RemoveDevice(deviceID)
	assert.Nil(t, err)
}

// TestDeviceCapacityLimit sets the device capacity limit
func TestDeviceCapacityLimit(t *testing.T) {
	pool := getStoragePool(t)
	assert.NotNil(t, pool)

	sdsID := getSdsID()
	assert.NotEqual(t, sdsID, "")

	dev := &types.DeviceParam{
		DeviceCurrentPathname: "/dev/sdc",
		SdsID:                 sdsID,
	}

	deviceID, err := pool.AttachDevice(dev)
	assert.Nil(t, err)
	assert.NotNil(t, deviceID)

	err = pool.SetDeviceCapacityLimit(deviceID, "300")
	assert.Nil(t, err)

	system := getSystem()
	device, err1 := system.GetDevice(deviceID)
	assert.Nil(t, err1)
	assert.Equal(t, device.CapacityLimitInKb, 314572800)

	err = pool.RemoveDevice(deviceID)
	assert.Nil(t, err)
}

// TestDeviceUpdateOriginalPathways updates device path if changed during server restart
func TestDeviceUpdateOriginalPathways(t *testing.T) {
	pool := getStoragePool(t)
	assert.NotNil(t, pool)

	sdsID := getSdsID()
	assert.NotEqual(t, sdsID, "")

	dev := &types.DeviceParam{
		DeviceCurrentPathname: "/dev/sdc",
		SdsID:                 sdsID,
	}

	deviceID, err := pool.AttachDevice(dev)
	assert.Nil(t, err)
	assert.NotNil(t, deviceID)

	err = pool.UpdateDeviceOriginalPathways(deviceID)
	assert.Nil(t, err)

	err = pool.RemoveDevice(deviceID)
	assert.Nil(t, err)
}

func TestGetDeviceByDeviceID(t *testing.T) {
	system := getSystem()

	device, _ := system.GetDevice("c7fc68a200000000")
	assert.NotNil(t, device)
	assert.NotNil(t, device.SdsID)
	assert.NotNil(t, device.StoragePoolID)
	if device == nil {
		return
	}
}

// TestGetDeviceByAttribute gets a single specific Device by attribute
func TestGetDeviceBySdsAttribute(t *testing.T) {
	sds := getAllSds(t)[0]
	assert.NotNil(t, sds)
	if sds == nil {
		return
	}

	devices := getAllSdsDevices(t)
	assert.NotNil(t, devices)
	assert.NotZero(t, len(devices))
	if devices == nil {
		return
	}

	found, err := sds.FindDevice("Name", devices[0].Device.Name)
	assert.Nil(t, err)
	assert.NotNil(t, found)
	assert.Equal(t, devices[0].Device.Name, found.Name)

	found, err = sds.FindDevice("ID", devices[0].Device.ID)
	assert.Nil(t, err)
	assert.NotNil(t, found)
	assert.Equal(t, devices[0].Device.ID, found.ID)
}

// TestGetDeviceByAttributeInvalid fails to get a single specific Device by attribute
func TestGetDeviceBySdsAttributeInvalid(t *testing.T) {
	sds := getAllSds(t)[0]
	assert.NotNil(t, sds)
	if sds == nil {
		return
	}

	devices := getAllSdsDevices(t)
	assert.NotNil(t, devices)
	assert.NotZero(t, len(devices))
	if devices == nil {
		return
	}

	found, err := sds.FindDevice("Name", invalidIdentifier)
	assert.NotNil(t, err)
	assert.Nil(t, found)

	found, err = sds.FindDevice("ID", invalidIdentifier)
	assert.NotNil(t, err)
	assert.Nil(t, found)
}
