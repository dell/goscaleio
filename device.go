package goscaleio

import (
	"errors"
	"fmt"
	"net/http"
	"reflect"
	"time"

	types "github.com/dell/goscaleio/types/v1"
)

// Device defines struct for Device
type Device struct {
	Device *types.Device
	client *Client
}

// NewDevice returns a new Device
func NewDevice(client *Client) *Device {
	return &Device{
		Device: &types.Device{},
		client: client,
	}
}

// NewDeviceEx returns a new Device
func NewDeviceEx(client *Client, device *types.Device) *Device {
	return &Device{
		Device: device,
		client: client,
	}
}

// AttachDevice attaches a device
func (sp *StoragePool) AttachDevice(
	path string,
	sdsID string) (string, error) {
	defer TimeSpent("AttachDevice", time.Now())

	deviceParam := &types.DeviceParam{
		Name:                  path,
		DeviceCurrentPathname: path,
		StoragePoolID:         sp.StoragePool.ID,
		SdsID:                 sdsID,
		TestMode:              "testAndActivate"}

	dev := types.DeviceResp{}
	err := sp.client.getJSONWithRetry(
		http.MethodPost, "/api/types/Device/instances",
		deviceParam, &dev)
	if err != nil {
		return "", err
	}

	return dev.ID, nil
}

// GetDevice returns a device
func (sp *StoragePool) GetDevice() ([]types.Device, error) {
	defer TimeSpent("GetDevice", time.Now())

	path := fmt.Sprintf(
		"/api/instances/StoragePool::%v/relationships/Device",
		sp.StoragePool.ID)

	var devices []types.Device
	err := sp.client.getJSONWithRetry(
		http.MethodGet, path, nil, &devices)
	if err != nil {
		return nil, err
	}

	return devices, nil
}

// FindDevice returns a Device
func (sp *StoragePool) FindDevice(
	field, value string) (*types.Device, error) {
	defer TimeSpent("FindDevice", time.Now())

	devices, err := sp.GetDevice()
	if err != nil {
		return nil, err
	}

	for _, device := range devices {
		valueOf := reflect.ValueOf(device)
		switch {
		case reflect.Indirect(valueOf).FieldByName(field).String() == value:
			return &device, nil
		}
	}

	return nil, errors.New("Couldn't find DEV")
}
