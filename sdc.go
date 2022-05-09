package goscaleio

import (
	"errors"
	"fmt"
	"net/http"
	"os/exec"
	"reflect"
	"strings"
	"time"

	types "github.com/dell/goscaleio/types/v1"
)


// Sdc defines struct for Sdc
type Sdc struct {
	Sdc    *types.Sdc
	client *Client
}

// NewSdc returns a new Sdc
func NewSdc(client *Client, sdc *types.Sdc) *Sdc {
	return &Sdc{
		Sdc:    sdc,
		client: client,
	}
}

// GetSdc returns a Sdc
func (s *System) GetSdc() ([]types.Sdc, error) {
	defer TimeSpent("GetSdc", time.Now())

	path := fmt.Sprintf("/api/instances/System::%v/relationships/Sdc",
		s.System.ID)

	var sdcs []types.Sdc
	err := s.client.getJSONWithRetry(
		http.MethodGet, path, nil, &sdcs)
	if err != nil {
		return nil, err
	}

	return sdcs, nil
}

// FindSdc returns a Sdc
func (s *System) FindSdc(field, value string) (*Sdc, error) {
	defer TimeSpent("FindSdc", time.Now())

	sdcs, err := s.GetSdc()
	if err != nil {
		return nil, err
	}

	for i, sdc := range sdcs {
		valueOf := reflect.ValueOf(sdc)
		switch {
		case reflect.Indirect(valueOf).FieldByName(field).String() == value:
			return NewSdc(s.client, &sdcs[i]), nil
		}
	}

	return nil, errors.New("Couldn't find SDC")
}

// GetStatistics returns a Sdc statistcs 
func (sdc *Sdc) GetStatistics() (*types.SdcStatistics, error) {
	defer TimeSpent("GetStatistics", time.Now())

	link, err := GetLink(sdc.Sdc.Links, "/api/Sdc/relationship/Statistics")
	if err != nil {
		return nil, err
	}

	var stats types.SdcStatistics
	err = sdc.client.getJSONWithRetry(
		http.MethodGet, link.HREF, nil, &stats)
	if err != nil {
		return nil, err
	}

	return &stats, nil
}

// GetVolume returns a volume 
func (sdc *Sdc) GetVolume() ([]*types.Volume, error) {
	defer TimeSpent("GetVolume", time.Now())

	link, err := GetLink(sdc.Sdc.Links, "/api/Sdc/relationship/Volume")
	if err != nil {
		return nil, err
	}

	var vols []*types.Volume
	err = sdc.client.getJSONWithRetry(
		http.MethodGet, link.HREF, nil, &vols)
	if err != nil {
		return nil, err
	}

	return vols, nil
}

// FindVolumes returns volumes
func (sdc *Sdc) FindVolumes() ([]*Volume, error) {
	defer TimeSpent("FindVolumes", time.Now())

	var rlt []*Volume
	vols, err := sdc.GetVolume()
	if err != nil {
		return nil, err
	}

	for _, v := range vols {
		volClient := NewVolume(sdc.client)
		volClient.Volume = v
		rlt = append(rlt, volClient)
	}

	return rlt, nil
}

// GetSdcLocalGUID returns GUID
func GetSdcLocalGUID() (string, error) {
	defer TimeSpent("GetSdcLocalGUID", time.Now())

	// get sdc kernel guid
	// /bin/emc/scaleio/drv_cfg --query_guid
	// sdcKernelGuid := "271bad82-08ee-44f2-a2b1-7e2787c27be1"

	out, err := exec.Command("/opt/emc/scaleio/sdc/bin/drv_cfg", "--query_guid").Output()
	if err != nil {
		return "", fmt.Errorf("GetSdcLocalGUID: query vols failed: %v", err)
	}

	sdcGUID := strings.Replace(string(out), "\n", "", -1)

	return sdcGUID, nil
}

// MapVolumeSdc maps a volume to Sdc
func (v *Volume) MapVolumeSdc(
	mapVolumeSdcParam *types.MapVolumeSdcParam) error {
	defer TimeSpent("MapVolumeSdc", time.Now())

	path := fmt.Sprintf("/api/instances/Volume::%s/action/addMappedSdc",
		v.Volume.ID)

	err := v.client.getJSONWithRetry(
		http.MethodPost, path, mapVolumeSdcParam, nil)
	if err != nil {
		return err
	}

	return nil
}

// UnmapVolumeSdc unmaps a volume from Sdc
func (v *Volume) UnmapVolumeSdc(
	unmapVolumeSdcParam *types.UnmapVolumeSdcParam) error {
	defer TimeSpent("UnmapVolumeSdc", time.Now())

	path := fmt.Sprintf("/api/instances/Volume::%s/action/removeMappedSdc",
		v.Volume.ID)

	err := v.client.getJSONWithRetry(
		http.MethodPost, path, unmapVolumeSdcParam, nil)
	if err != nil {
		return err
	}

	return nil
}

// SetMappedSdcLimits sets Sdc mapped limits 
func (v *Volume) SetMappedSdcLimits(
	setMappedSdcLimitsParam *types.SetMappedSdcLimitsParam) error {
	defer TimeSpent("SetMappedSdcLimits", time.Now())

	path := fmt.Sprintf(
		"/api/instances/Volume::%s/action/setMappedSdcLimits",
		v.Volume.ID)

	err := v.client.getJSONWithRetry(
		http.MethodPost, path, setMappedSdcLimitsParam, nil)
	if err != nil {
		return err
	}

	return nil
}
