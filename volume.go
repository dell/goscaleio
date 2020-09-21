package goscaleio

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"

	types "github.com/dell/goscaleio/types/v1"
)

// For unit testing on Windows - dev is not in the same place
var FSDevDirectoryPrefix string = ""

type SdcMappedVolume struct {
	MdmID     string
	VolumeID  string
	SdcDevice string
	// Mounted   bool
	// MountPath bool
	// Mapped    bool
}

type Volume struct {
	Volume *types.Volume
	client *Client
}

func NewVolume(client *Client) *Volume {
	return &Volume{
		Volume: &types.Volume{},
		client: client,
	}
}

func (sp *StoragePool) GetVolume(
	volumehref, volumeid, ancestorvolumeid, volumename string,
	getSnapshots bool) ([]*types.Volume, error) {
	defer TimeSpent("GetVolume", time.Now())

	var (
		err     error
		path    string
		volume  = &types.Volume{}
		volumes []*types.Volume
	)

	if volumename != "" {
		volumeid, err = sp.FindVolumeID(volumename)
		if err != nil && err.Error() == "Not found" {
			return nil, nil
		}
		if err != nil {
			return nil, fmt.Errorf("Error: problem finding volume: %s", err)
		}
	}

	if volumeid != "" {
		path = fmt.Sprintf("/api/instances/Volume::%s", volumeid)
	} else if volumehref == "" {
		link, err := GetLink(sp.StoragePool.Links,
			"/api/StoragePool/relationship/Volume")
		if err != nil {
			return nil, err
		}
		path = link.HREF
	} else {
		path = volumehref
	}

	if volumehref == "" && volumeid == "" {
		err = sp.client.getJSONWithRetry(
			http.MethodGet, path, nil, &volumes)
	} else {
		err = sp.client.getJSONWithRetry(
			http.MethodGet, path, nil, volume)
	}
	if err != nil {
		return nil, err
	}

	if volumehref == "" && volumeid == "" {
		var volumesNew []*types.Volume
		for _, v := range volumes {
			if (!getSnapshots && v.AncestorVolumeID == ancestorvolumeid) || (getSnapshots && v.AncestorVolumeID != "") {
				volumesNew = append(volumesNew, v)
			}
		}
		volumes = volumesNew
	} else {
		volumes = append(volumes, volume)
	}
	return volumes, nil
}

func (sp *StoragePool) FindVolumeID(volumename string) (string, error) {
	defer TimeSpent("FindVolumeID", time.Now())

	volumeQeryIdByKeyParam := &types.VolumeQeryIdByKeyParam{
		Name: volumename,
	}

	path := fmt.Sprintf("/api/types/Volume/instances/action/queryIdByKey")

	volumeID, err := sp.client.getStringWithRetry(
		http.MethodPost, path, volumeQeryIdByKeyParam)
	if err != nil {
		return "", err
	}

	return volumeID, nil
}

// GetLocalVolumeMapByRegex will return the SdcMappedVolume entries matching supplied regex values
func GetLocalVolumeMapByRegex(systemIDRegex string, volumeIDRegex string) (mappedVolumes []*SdcMappedVolume, err error) {
	actualSystemRegex := systemIDRegex
	actualVolumeRegex := volumeIDRegex
	if actualSystemRegex == "" {
		// if an empty regex was supplied, change it to allow any word
		actualSystemRegex = `\w*`
	}
	if actualVolumeRegex == "" {
		// if an empty regex was supplied, change it to allow any word
		actualVolumeRegex = `\w*`
	}
	// look for a specific systems and volumes
	return getVolumeMapping(actualSystemRegex, actualVolumeRegex)
}

// GetLocalVolumeMap will return all SdcMappedVolume entries
func GetLocalVolumeMap() (mappedVolumes []*SdcMappedVolume, err error) {
	// pass in a regex value for any word in place of the system ID and volume ID
	return getVolumeMapping(`\w*`, `\w*`)
}

func getVolumeMapping(sysID string, volID string) (mappedVolumes []*SdcMappedVolume, err error) {
	defer TimeSpent("GetLocalVolumeMap", time.Now())

	mappedVolumesMap := make(map[string]*SdcMappedVolume)

	diskIDPath := FSDevDirectoryPrefix + "/dev/disk/by-id"
	files, _ := ioutil.ReadDir(diskIDPath)
	strRegex := fmt.Sprintf(`^emc-vol-%s-%s$`, sysID, volID)
	r, _ := regexp.Compile(strRegex)
	for _, f := range files {
		matched := r.MatchString(f.Name())
		if matched {
			split := strings.Split(f.Name(), "-")
			mdmVolumeID := fmt.Sprintf("%s-%s", split[2], split[3])
			devPath, _ := filepath.EvalSymlinks(fmt.Sprintf("%s/%s", diskIDPath, f.Name()))
			mappedVolumesMap[mdmVolumeID] = &SdcMappedVolume{MdmID: split[2], VolumeID: split[3], SdcDevice: devPath}
		}
	}

	keys := make([]string, 0, len(mappedVolumesMap))
	for key := range mappedVolumesMap {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	for _, key := range keys {
		mappedVolumes = append(mappedVolumes, mappedVolumesMap[key])
	}

	return mappedVolumes, nil
}

func (sp *StoragePool) CreateVolume(
	volume *types.VolumeParam) (*types.VolumeResp, error) {
	defer TimeSpent("CreateVolume", time.Now())

	path := "/api/types/Volume/instances"

	volume.StoragePoolID = sp.StoragePool.ID
	volume.ProtectionDomainID = sp.StoragePool.ProtectionDomainID

	volumeResp := &types.VolumeResp{}
	err := sp.client.getJSONWithRetry(
		http.MethodPost, path, volume, volumeResp)
	if err != nil {
		return nil, err
	}

	return volumeResp, nil
}

func (v *Volume) GetVTree() (*types.VTree, error) {
	defer TimeSpent("GetVTree", time.Now())

	link, err := GetLink(v.Volume.Links, "/api/parent/relationship/vtreeId")
	if err != nil {
		return nil, err
	}

	vtree := &types.VTree{}
	err = v.client.getJSONWithRetry(
		http.MethodGet, link.HREF, nil, vtree)
	if err != nil {
		return nil, err
	}

	return vtree, nil
}

func (v *Volume) GetVolumeStatistics() (*types.VolumeStatistics, error) {
	defer TimeSpent("GetStatistics", time.Now())

	link, err := GetLink(v.Volume.Links, "/api/Volume/relationship/Statistics")
	if err != nil {
		return nil, err
	}

	var stats types.VolumeStatistics
	err = v.client.getJSONWithRetry(
		http.MethodGet, link.HREF, nil, &stats)
	if err != nil {
		return nil, err
	}

	return &stats, nil
}

func (v *Volume) RemoveVolume(removeMode string) error {
	defer TimeSpent("RemoveVolume", time.Now())

	link, err := GetLink(v.Volume.Links, "self")
	if err != nil {
		return err
	}

	path := fmt.Sprintf("%v/action/removeVolume", link.HREF)

	if removeMode == "" {
		removeMode = "ONLY_ME"
	}
	removeVolumeParam := &types.RemoveVolumeParam{
		RemoveMode: removeMode,
	}

	err = v.client.getJSONWithRetry(
		http.MethodPost, path, removeVolumeParam, nil)
	return err
}

func (v *Volume) SetVolumeSize(sizeInGB string) error {

	link, err := GetLink(v.Volume.Links, "self")
	if err != nil {
		return err
	}
	path := fmt.Sprintf("%v/action/setVolumeSize", link.HREF)
	payload := &types.SetVolumeSizeParam{
		SizeInGB: sizeInGB,
	}
	err = v.client.getJSONWithRetry(
		http.MethodPost, path, payload, nil)
	return err
}
