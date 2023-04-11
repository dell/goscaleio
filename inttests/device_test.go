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

	"github.com/dell/goscaleio"
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
		outDevice := goscaleio.NewDeviceEx(C, &d)
		allDevice = append(allDevice, outDevice)
		// create a device via NewDevice for testing purposes
		tempDevice := goscaleio.NewDevice(C)
		tempDevice.Device = &d
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
		outDevice := goscaleio.NewDeviceEx(C, &d)
		allDevice = append(allDevice, outDevice)
		// create a device via NewDevice for testing purposes
		tempDevice := goscaleio.NewDevice(C)
		tempDevice.Device = &d
		assert.Equal(t, outDevice.Device.ID, tempDevice.Device.ID)
	}
	return allDevice

}

// TestGetDevices will return all Device instances
func TestGetDevices(t *testing.T) {
	getAllDevices(t)
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

// TestAddDeviceInvalid will attempt to add an invalid device to an invalid SDS
func TestAddDeviceInvalid(t *testing.T) {
	pool := getStoragePool(t)
	assert.NotNil(t, pool)

	deviceID, err := pool.AttachDevice("/invalidPath/invalidDevice", invalidIdentifier)
	assert.NotNil(t, err)
	assert.Equal(t, "", deviceID)

}

func TestGetDeviceByDeviceID(t *testing.T) {

	d := goscaleio.NewDevice(C)
	d.Device.ID = "c7fc68a200000000"

	device, _ := d.GetDevice()
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
