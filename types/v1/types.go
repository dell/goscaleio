// Copyright © 2019 - 2022 Dell Inc. or its subsidiaries. All Rights Reserved.
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

package goscaleio

import (
	"fmt"
	"net/http"
	"sync"
)

const errorWithDetails = "Error with details"
const (
	trueType  = "TRUE"
	falseType = "FALSE"
)

// GetBoolType returns the true and false strings expected by the REST API
func GetBoolType(b bool) string {
	if b {
		return trueType
	}
	return falseType
}

// ErrorMessageDetails defines contents of an error msg
type ErrorMessageDetails struct {
	Error        string `json:"error"`
	ReturnCode   int    `json:"rc"`
	ErrorMessage string `json:"errorMessage"`
}

// Error struct defines the structure of an error
type Error struct {
	Message        string                `json:"message"`
	HTTPStatusCode int                   `json:"httpStatusCode"`
	ErrorCode      int                   `json:"errorCode"`
	ErrorDetails   []ErrorMessageDetails `json:"details"`
}

func (e Error) Error() string {
	if e.Message == errorWithDetails && len(e.ErrorDetails) > 0 {
		fmt.Printf("goscaleio.Error Error with details  %#v\n", e)
		if e.ErrorDetails[0].ErrorMessage != "" {
			e.Message = e.ErrorDetails[0].ErrorMessage
			return e.ErrorDetails[0].ErrorMessage
		}
		if e.ErrorDetails[0].Error != "" {
			translation := TranslateErrorCodeToErrorMessage(e.ErrorDetails[0].Error)
			if translation != "" {
				e.Message = translation
				e.Message = e.ErrorDetails[0].ErrorMessage
				return translation
			}
		}
		// No ErrorMessage or Error string, have to punt
	}
	return e.Message
}

// type session struct {
// 	Link []*types.Link `xml:"Link"`
// }

// System defines struct of PFlex array
type System struct {
	MdmMode                               string   `json:"mdmMode"`
	MdmClusterState                       string   `json:"mdmClusterState"`
	SecondaryMdmActorIPList               []string `json:"secondaryMdmActorIpList"`
	InstallID                             string   `json:"installId"`
	PrimaryActorIPList                    []string `json:"primaryMdmActorIpList"`
	SystemVersionName                     string   `json:"systemVersionName"`
	CapacityAlertHighThresholdPercent     int      `json:"capacityAlertHighThresholdPercent"`
	CapacityAlertCriticalThresholdPercent int      `json:"capacityAlertCriticalThresholdPercent"`
	RemoteReadOnlyLimitState              bool     `json:"remoteReadOnlyLimitState"`
	PrimaryMdmActorPort                   int      `json:"primaryMdmActorPort"`
	SecondaryMdmActorPort                 int      `json:"secondaryMdmActorPort"`
	TiebreakerMdmActorPort                int      `json:"tiebreakerMdmActorPort"`
	MdmManagementPort                     int      `json:"mdmManagementPort"`
	TiebreakerMdmIPList                   []string `json:"tiebreakerMdmIpList"`
	MdmManagementIPList                   []string `json:"mdmManagementIPList"`
	DefaultIsVolumeObfuscated             bool     `json:"defaultIsVolumeObfuscated"`
	RestrictedSdcModeEnabled              bool     `json:"restrictedSdcModeEnabled"`
	RestrictedSdcMode                     string   `json:"restrictedSdcMode"`
	Swid                                  string   `json:"swid"`
	DaysInstalled                         int      `json:"daysInstalled"`
	MaxCapacityInGb                       string   `json:"maxCapacityInGb"`
	CapacityTimeLeftInDays                string   `json:"capacityTimeLeftInDays"`
	EnterpriseFeaturesEnabled             bool     `json:"enterpriseFeaturesEnabled"`
	IsInitialLicense                      bool     `json:"isInitialLicense"`
	Name                                  string   `json:"name"`
	ID                                    string   `json:"id"`
	Links                                 []*Link  `json:"links"`
}

// Link defines struct of Link
type Link struct {
	Rel  string `json:"rel"`
	HREF string `json:"href"`
}

// BWC defines struct of BWC
type BWC struct {
	TotalWeightInKb int `json:"totalWeightInKb"`
	NumOccured      int `json:"numOccured"`
	NumSeconds      int `json:"numSeconds"`
}

// Statistics defines struct of Statistics for Pflex Array
type Statistics struct {
	PrimaryReadFromDevBwc                    BWC `json:"primaryReadFromDevBwc"`
	NumOfStoragePools                        int `json:"numOfStoragePools"`
	ProtectedCapacityInKb                    int `json:"protectedCapacityInKb"`
	MovingCapacityInKb                       int `json:"movingCapacityInKb"`
	SnapCapacityInUseOccupiedInKb            int `json:"snapCapacityInUseOccupiedInKb"`
	SnapCapacityInUseInKb                    int `json:"snapCapacityInUseInKb"`
	ActiveFwdRebuildCapacityInKb             int `json:"activeFwdRebuildCapacityInKb"`
	DegradedHealthyVacInKb                   int `json:"degradedHealthyVacInKb"`
	ActiveMovingRebalanceJobs                int `json:"activeMovingRebalanceJobs"`
	TotalReadBwc                             BWC `json:"totalReadBwc"`
	MaxCapacityInKb                          int `json:"maxCapacityInKb"`
	PendingBckRebuildCapacityInKb            int `json:"pendingBckRebuildCapacityInKb"`
	ActiveMovingOutFwdRebuildJobs            int `json:"activeMovingOutFwdRebuildJobs"`
	CapacityLimitInKb                        int `json:"capacityLimitInKb"`
	SecondaryVacInKb                         int `json:"secondaryVacInKb"`
	PendingFwdRebuildCapacityInKb            int `json:"pendingFwdRebuildCapacityInKb"`
	ThinCapacityInUseInKb                    int `json:"thinCapacityInUseInKb"`
	AtRestCapacityInKb                       int `json:"atRestCapacityInKb"`
	ActiveMovingInBckRebuildJobs             int `json:"activeMovingInBckRebuildJobs"`
	DegradedHealthyCapacityInKb              int `json:"degradedHealthyCapacityInKb"`
	NumOfScsiInitiators                      int `json:"numOfScsiInitiators"`
	NumOfUnmappedVolumes                     int `json:"numOfUnmappedVolumes"`
	FailedCapacityInKb                       int `json:"failedCapacityInKb"`
	SecondaryReadFromDevBwc                  BWC `json:"secondaryReadFromDevBwc"`
	NumOfVolumes                             int `json:"numOfVolumes"`
	SecondaryWriteBwc                        BWC `json:"secondaryWriteBwc"`
	ActiveBckRebuildCapacityInKb             int `json:"activeBckRebuildCapacityInKb"`
	FailedVacInKb                            int `json:"failedVacInKb"`
	PendingMovingCapacityInKb                int `json:"pendingMovingCapacityInKb"`
	ActiveMovingInRebalanceJobs              int `json:"activeMovingInRebalanceJobs"`
	PendingMovingInRebalanceJobs             int `json:"pendingMovingInRebalanceJobs"`
	BckRebuildReadBwc                        BWC `json:"bckRebuildReadBwc"`
	DegradedFailedVacInKb                    int `json:"degradedFailedVacInKb"`
	NumOfSnapshots                           int `json:"numOfSnapshots"`
	RebalanceCapacityInKb                    int `json:"rebalanceCapacityInKb"`
	FwdRebuildReadBwc                        BWC `json:"fwdRebuildReadBwc"`
	NumOfSdc                                 int `json:"numOfSdc"`
	ActiveMovingInFwdRebuildJobs             int `json:"activeMovingInFwdRebuildJobs"`
	NumOfVtrees                              int `json:"numOfVtrees"`
	ThickCapacityInUseInKb                   int `json:"thickCapacityInUseInKb"`
	ProtectedVacInKb                         int `json:"protectedVacInKb"`
	PendingMovingInBckRebuildJobs            int `json:"pendingMovingInBckRebuildJobs"`
	CapacityAvailableForVolumeAllocationInKb int `json:"capacityAvailableForVolumeAllocationInKb"`
	PendingRebalanceCapacityInKb             int `json:"pendingRebalanceCapacityInKb"`
	PendingMovingRebalanceJobs               int `json:"pendingMovingRebalanceJobs"`
	NumOfProtectionDomains                   int `json:"numOfProtectionDomains"`
	NumOfSds                                 int `json:"numOfSds"`
	CapacityInUseInKb                        int `json:"capacityInUseInKb"`
	BckRebuildWriteBwc                       BWC `json:"bckRebuildWriteBwc"`
	DegradedFailedCapacityInKb               int `json:"degradedFailedCapacityInKb"`
	NumOfThinBaseVolumes                     int `json:"numOfThinBaseVolumes"`
	PendingMovingOutFwdRebuildJobs           int `json:"pendingMovingOutFwdRebuildJobs"`
	SecondaryReadBwc                         BWC `json:"secondaryReadBwc"`
	PendingMovingOutBckRebuildJobs           int `json:"pendingMovingOutBckRebuildJobs"`
	RebalanceWriteBwc                        BWC `json:"rebalanceWriteBwc"`
	PrimaryReadBwc                           BWC `json:"primaryReadBwc"`
	NumOfVolumesInDeletion                   int `json:"numOfVolumesInDeletion"`
	NumOfDevices                             int `json:"numOfDevices"`
	RebalanceReadBwc                         BWC `json:"rebalanceReadBwc"`
	InUseVacInKb                             int `json:"inUseVacInKb"`
	UnreachableUnusedCapacityInKb            int `json:"unreachableUnusedCapacityInKb"`
	TotalWriteBwc                            BWC `json:"totalWriteBwc"`
	SpareCapacityInKb                        int `json:"spareCapacityInKb"`
	ActiveMovingOutBckRebuildJobs            int `json:"activeMovingOutBckRebuildJobs"`
	PrimaryVacInKb                           int `json:"primaryVacInKb"`
	NumOfThickBaseVolumes                    int `json:"numOfThickBaseVolumes"`
	BckRebuildCapacityInKb                   int `json:"bckRebuildCapacityInKb"`
	NumOfMappedToAllVolumes                  int `json:"numOfMappedToAllVolumes"`
	ActiveMovingCapacityInKb                 int `json:"activeMovingCapacityInKb"`
	PendingMovingInFwdRebuildJobs            int `json:"pendingMovingInFwdRebuildJobs"`
	ActiveRebalanceCapacityInKb              int `json:"activeRebalanceCapacityInKb"`
	RmcacheSizeInKb                          int `json:"rmcacheSizeInKb"`
	FwdRebuildCapacityInKb                   int `json:"fwdRebuildCapacityInKb"`
	FwdRebuildWriteBwc                       BWC `json:"fwdRebuildWriteBwc"`
	PrimaryWriteBwc                          BWC `json:"primaryWriteBwc"`
	NetUserDataCapacityInKb                  int `json:"netUserDataCapacityInKb"`
	NetUnusedCapacityInKb                    int `json:"netUnusedCapacityInKb"`
	VolumeAddressSpaceInKb                   int `json:"volumeAddressSpaceInKb"`
}

// SdcStatistics defines struct of Statistics for PFlex SDC
type SdcStatistics struct {
	UserDataReadBwc         BWC      `json:"userDataReadBwc"`
	UserDataWriteBwc        BWC      `json:"userDataWriteBwc"`
	UserDataTrimBwc         BWC      `json:"userDataTrimBwc"`
	UserDataSdcReadLatency  BWC      `json:"userDataSdcReadLatency"`
	UserDataSdcWriteLatency BWC      `json:"userDataSdcWriteLatency"`
	UserDataSdcTrimLatency  BWC      `json:"userDataSdcTrimLatency"`
	VolumeIds               []string `json:"volumeIds"`
	NumOfMappedVolumes      int      `json:"numOfMappedVolumes"`
}

// VolumeStatistics defines struct of Statistics for PFlex volume
type VolumeStatistics struct {
	UserDataReadBwc         BWC      `json:"userDataReadBwc"`
	UserDataWriteBwc        BWC      `json:"userDataWriteBwc"`
	UserDataTrimBwc         BWC      `json:"userDataTrimBwc"`
	UserDataSdcReadLatency  BWC      `json:"userDataSdcReadLatency"`
	UserDataSdcWriteLatency BWC      `json:"userDataSdcWriteLatency"`
	UserDataSdcTrimLatency  BWC      `json:"userDataSdcTrimLatency"`
	MappedSdcIds            []string `json:"mappedSdcIds"`
	NumOfMappedSdcs         int      `json:"numOfMappedSdcs"`
}

// User defines struct of User for PFlex array
type User struct {
	SystemID              string  `json:"systemId"`
	UserRole              string  `json:"userRole"`
	PasswordChangeRequire bool    `json:"passwordChangeRequired"`
	Name                  string  `json:"name"`
	ID                    string  `json:"id"`
	Links                 []*Link `json:"links"`
}

// ScsiInitiator defines struct for ScsiInitiator
type ScsiInitiator struct {
	Name     string  `json:"name"`
	IQN      string  `json:"iqn"`
	SystemID string  `json:"systemID"`
	Links    []*Link `json:"links"`
}

// PDRfCacheOpMode is an enum type for Protection Domain Rf Cache Operational Mode
type PDRfCacheOpMode string

// Available values for enum type PDRfCacheOpMode
const (
	PDRCModeRead         PDRfCacheOpMode = "Read"
	PDRCModeWrite        PDRfCacheOpMode = "Write"
	PDRCModeReadAndWrite PDRfCacheOpMode = "ReadAndWrite"
	PDRCModeWriteMiss    PDRfCacheOpMode = "WriteMiss"
)

// PDCounterWindow defines one window for a Protection Domain Failure Counter
type PDCounterWindow struct {
	Threshold       int `json:"threshold"`
	WindowSizeInSec int `json:"windowSizeInSec"`
}

// PDCounterParams defines all the windows for a Protection Domain Failure Counter
type PDCounterParams struct {
	ShortWindow  PDCounterWindow `json:"shortWindow"`
	MediumWindow PDCounterWindow `json:"mediumWindow"`
	LongWindow   PDCounterWindow `json:"longWindow"`
}

// PDConnInfo defines Protection Domain Connection information
type PDConnInfo struct {
	ClientServerConnStatus string  `json:"clientServerConnStatus"`
	DisconnectedClientID   *string `json:"disconnectedClientId"`
	DisconnectedClientName *string `json:"disconnectedClientName"`
	DisconnectedServerID   *string `json:"disconnectedServerId"`
	DisconnectedServerName *string `json:"disconnectedServerName"`
	DisconnectedServerIP   *string `json:"disconnectedServerIp"`
}

// ProtectionDomain defines struct for PFlex ProtectionDomain
type ProtectionDomain struct {
	SystemID                    string     `json:"systemId"`
	SdrSdsConnectivityInfo      PDConnInfo `json:"sdrSdsConnectivityInfo"`
	ReplicationCapacityMaxRatio *int       `json:"replicationCapacityMaxRatio"`

	// Network throttling params
	RebuildNetworkThrottlingInKbps                   int  `json:"rebuildNetworkThrottlingInKbps"`
	RebalanceNetworkThrottlingInKbps                 int  `json:"rebalanceNetworkThrottlingInKbps"`
	OverallIoNetworkThrottlingInKbps                 int  `json:"overallIoNetworkThrottlingInKbps"`
	VTreeMigrationNetworkThrottlingInKbps            int  `json:"vtreeMigrationNetworkThrottlingInKbps"`
	ProtectedMaintenanceModeNetworkThrottlingInKbps  int  `json:"protectedMaintenanceModeNetworkThrottlingInKbps"`
	OverallIoNetworkThrottlingEnabled                bool `json:"overallIoNetworkThrottlingEnabled"`
	RebuildNetworkThrottlingEnabled                  bool `json:"rebuildNetworkThrottlingEnabled"`
	RebalanceNetworkThrottlingEnabled                bool `json:"rebalanceNetworkThrottlingEnabled"`
	VTreeMigrationNetworkThrottlingEnabled           bool `json:"vtreeMigrationNetworkThrottlingEnabled"`
	ProtectedMaintenanceModeNetworkThrottlingEnabled bool `json:"protectedMaintenanceModeNetworkThrottlingEnabled"`

	// Fine Granularity Params
	FglDefaultNumConcurrentWrites int  `json:"fglDefaultNumConcurrentWrites"`
	FglMetadataCacheEnabled       bool `json:"fglMetadataCacheEnabled"`
	FglDefaultMetadataCacheSize   int  `json:"fglDefaultMetadataCacheSize"`

	// RfCache Params
	RfCacheEnabled         bool            `json:"rfcacheEnabled"`
	RfCacheAccpID          string          `json:"rfcacheAccpId"`
	RfCacheOperationalMode PDRfCacheOpMode `json:"rfcacheOpertionalMode"`
	RfCachePageSizeKb      int             `json:"rfcachePageSizeKb"`
	RfCacheMaxIoSizeKb     int             `json:"rfcacheMaxIoSizeKb"`

	// Counter Params
	SdsConfigurationFailureCP            PDCounterParams `json:"sdsConfigurationFailureCounter"`
	SdsDecoupledCP                       PDCounterParams `json:"sdsDecoupledCounterParameters"`
	MdmSdsNetworkDisconnectionsCP        PDCounterParams `json:"mdmSdsNetworkDisconnectionsCounterParameters"`
	SdsSdsNetworkDisconnectionsCP        PDCounterParams `json:"sdsSdsNetworkDisconnectionsCounterParameters"`
	SdsReceiveBufferAllocationFailuresCP PDCounterParams `json:"sdsReceiveBufferAllocationFailuresCounterParameters"`

	ProtectionDomainState string  `json:"protectionDomainState"`
	Name                  string  `json:"name"`
	ID                    string  `json:"id"`
	Links                 []*Link `json:"links"`
}

// ProtectionDomainParam defines struct for ProtectionDomainParam
type ProtectionDomainParam struct {
	Name string `json:"name"`
}

// ChangeSdcNameParam defines struct for passing parameters to changeSDCname endpoint
type ChangeSdcNameParam struct {
	SdcName string `json:"sdcName"`
}

type ApproveSdcParam struct {
	SdcGuid string `json:"sdcGuid"`
}

type ApproveSdcByGuidResponse struct {
	SdcId string `json:"id"`
}

// ProtectionDomainResp defines struct for ProtectionDomainResp
type ProtectionDomainResp struct {
	ID string `json:"id"`
}

// Sdc defines struct for PFlex Sdc
type Sdc struct {
	SystemID           string  `json:"systemId"`
	SdcApproved        bool    `json:"sdcApproved"`
	SdcIP              string  `json:"SdcIp"`
	OnVMWare           bool    `json:"onVmWare"`
	SdcGUID            string  `json:"sdcGuid"`
	MdmConnectionState string  `json:"mdmConnectionState"`
	Name               string  `json:"name"`
	ID                 string  `json:"id"`
	Links              []*Link `json:"links"`
}

// SdsIP defines struct for SdsIP
type SdsIP struct {
	IP   string `json:"ip"`
	Role string `json:"role,omitempty"`
}

// SdsIPList defines struct for SdsIPList
type SdsIPList struct {
	SdsIP SdsIP `json:"SdsIp"`
}

// SdsWindowType defines struct for SdsWindowType
type SdsWindowType struct {
	Threshold            int `json:"threshold,omitempty"`
	WindowSizeInSec      int `json:"windowSizeInSec,omitempty"`
	LastOscillationCount int `json:"lastOscillationCount,omitempty"`
	LastOscillationTime  int `json:"lastOscillationTime,omitempty"`
	MaxFailuresCount     int `json:"maxFailuresCount,omitempty"`
}

// SdsWindow defines struct for SdsWindow
type SdsWindow struct {
	ShortWindow  SdsWindowType `json:"shortWindow,omitempty"`
	MediumWindow SdsWindowType `json:"mediumWindow,omitempty"`
	LongWindow   SdsWindowType `json:"longWindow,omitempty"`
}

// RaidControllers defines struct for raid controllers
type RaidControllers struct {
	SerialNumber    string `json:"serialNumber,omitempty"`
	ModelName       string `json:"modelName,omitempty"`
	VendorName      string `json:"vendorName,omitempty"`
	FirmwareVersion string `json:"firmwareVersion,omitempty"`
	DriverVersion   string `json:"driverVersion,omitempty"`
	DriverName      string `json:"driverName,omitempty"`
	PciAddress      string `json:"pciAddress,omitempty"`
	Status          string `json:"status,omitempty"`
	BatteryStatus   string `json:"batteryStatus,omitempty"`
}

// CertificateInfo defines struct for certificate information
type CertificateInfo struct {
	Subject             string `json:"subject,omitempty"`
	Issuer              string `json:"issuer,omitempty"`
	ValidFrom           string `json:"validFrom,omitempty"`
	ValidTo             string `json:"validTo,omitempty"`
	Thumbprint          string `json:"thumbprint,omitempty"`
	ValidFromAsn1Format string `json:"validFromAsn1Format,omitempty"`
	ValidToAsn1Format   string `json:"validToAsn1Format,omitempty"`
}

// Sds defines struct for Sds
type Sds struct {
	ID                                          string            `json:"id"`
	Name                                        string            `json:"name,omitempty"`
	ProtectionDomainID                          string            `json:"protectionDomainId"`
	IPList                                      []*SdsIP          `json:"ipList"`
	Port                                        int               `json:"port,omitempty"`
	SdsState                                    string            `json:"sdsState"`
	MembershipState                             string            `json:"membershipState"`
	MdmConnectionState                          string            `json:"mdmConnectionState"`
	DrlMode                                     string            `json:"drlMode,omitempty"`
	RmcacheEnabled                              bool              `json:"rmcacheEnabled,omitempty"`
	RmcacheSizeInKb                             int               `json:"rmcacheSizeInKb,omitempty"`
	RmcacheFrozen                               bool              `json:"rmcacheFrozen,omitempty"`
	IsOnVMware                                  bool              `json:"isOnVmWare,omitempty"`
	FaultSetID                                  string            `json:"faultSetId,omitempty"`
	NumOfIoBuffers                              int               `json:"numOfIoBuffers,omitempty"`
	RmcacheMemoryAllocationState                string            `json:"RmcacheMemoryAllocationState,omitempty"`
	PerformanceProfile                          string            `json:"perfProfile,omitempty"`
	SoftwareVersionInfo                         string            `json:"softwareVersionInfo,omitempty"`
	ConfiguredDrlMode                           string            `json:"configuredDrlMode,omitempty"`
	RfcacheEnabled                              bool              `json:"rfcacheEnabled,omitempty"`
	MaintenanceState                            string            `json:"maintenanceState,omitempty"`
	MaintenanceType                             string            `json:"maintenanceType,omitempty"`
	RfcacheErrorLowResources                    bool              `json:"rfcacheErrorLowResources,omitempty"`
	RfcacheErrorAPIVersionMismatch              bool              `json:"rfcacheErrorApiVersionMismatch,omitempty"`
	RfcacheErrorInconsistentCacheConfiguration  bool              `json:"rfcacheErrorInconsistentCacheConfiguration,omitempty"`
	RfcacheErrorInconsistentSourceConfiguration bool              `json:"rfcacheErrorInconsistentSourceConfiguration,omitempty"`
	RfcacheErrorInvalidDriverPath               bool              `json:"rfcacheErrorInvalidDriverPath,omitempty"`
	RfcacheErrorDeviceDoesNotExist              bool              `json:"rfcacheErrorDeviceDoesNotExist,omitempty"`
	AuthenticationError                         string            `json:"authenticationError,omitempty"`
	FglNumConcurrentWrites                      int               `json:"fglNumConcurrentWrites,omitempty"`
	FglMetadataCacheState                       string            `json:"fglMetadataCacheState,omitempty"`
	FglMetadataCacheSize                        int               `json:"fglMetadataCacheSize,omitempty"`
	NumRestarts                                 int               `json:"numRestarts,omitempty"`
	LastUpgradeTime                             int               `json:"lastUpgradeTime,omitempty"`
	SdsDecoupled                                SdsWindow         `json:"sdsDecoupled,omitempty"`
	SdsConfigurationFailure                     SdsWindow         `json:"sdsConfigurationFailure,omitempty"`
	SdsReceiveBufferAllocationFailures          SdsWindow         `json:"sdsReceiveBufferAllocationFailures,omitempty"`
	RaidControllers                             []RaidControllers `json:"raidControllers,omitempty"`
	CertificateInfo                             CertificateInfo   `json:"certificateInfo,omitempty"`
	Links                                       []*Link           `json:"links"`
}

// DeviceInfo defines struct for DeviceInfo
type DeviceInfo struct {
	DevicePath    string `json:"devicePath"`
	StoragePoolID string `json:"storagePoolId"`
	DeviceName    string `json:"deviceName,omitempty"`
}

// Constants representing states of SDS
const (
	SdsDrlModeVolatile        = "Volatile"
	SdsDrlModeNonVolatile     = "NonVolatile"
	PerformanceProfileHigh    = "HighPerformance"
	PerformanceProfileCompact = "Compact"
)

// SdsParam defines struct for SdsParam
type SdsParam struct {
	Name               string        `json:"name,omitempty"`
	IPList             []*SdsIPList  `json:"sdsIpList"`
	Port               string        `json:"sdsPort,omitempty"`
	DrlMode            string        `json:"drlMode,omitempty"`
	RmcacheEnabled     string        `json:"rmcacheEnabled,omitempty"`
	RmcacheSizeInKb    string        `json:"rmcacheSizeInKb,omitempty"`
	RmcacheFrozen      bool          `json:"rmcacheFrozen,omitempty"`
	ProtectionDomainID string        `json:"protectionDomainId"`
	FaultSetID         string        `json:"faultSetId,omitempty"`
	NumOfIoBuffers     string        `json:"numOfIoBuffers,omitempty"`
	DeviceInfoList     []*DeviceInfo `json:"deviceInfoList,omitempty"`
	ForceClean         bool          `json:"forceClean,omitempty"`
	DeviceTestTimeSecs int           `json:"deviceTestTimeSecs ,omitempty"`
	DeviceTestMode     string        `json:"deviceTestMode,omitempty"`
}

// SdsResp defines struct for SdsResp
type SdsResp struct {
	ID string `json:"id"`
}

// SdsIPRole defines struct for Sds IP and Role
type SdsIPRole struct {
	SdsIPToSet string `json:"sdsIpToSet"`
	NewRole    string `json:"newRole"`
}

// SdsName defines struct for Sds Name
type SdsName struct {
	Name string `json:"name"`
}

// SdsPort defines struct for Sds Port
type SdsPort struct {
	SdsPort string `json:"sdsPort"`
}

// Device defines struct for Device
type Device struct {
	ID                     string `json:"id,omitempty"`
	Name                   string `json:"name,omitempty"`
	DeviceCurrentPathname  string `json:"deviceCurrentPathname"`
	DeviceOriginalPathname string `json:"deviceOriginalPathname,omitempty"`
	DeviceState            string `json:"deviceState,omitempty"`
	ErrorState             string `json:"errorState,omitempty"`
	CapacityLimitInKb      int    `json:"capacityLimitInKb,omitempty"`
	MaxCapacityInKb        int    `json:"maxCapacityInKb,omitempty"`
	StoragePoolID          string `json:"storagePoolId"`
	SdsID                  string `json:"sdsId"`
}

// DeviceParam defines struct for DeviceParam
type DeviceParam struct {
	Name                  string `json:"name,omitempty"`
	DeviceCurrentPathname string `json:"deviceCurrentPathname"`
	CapacityLimitInKb     int    `json:"capacityLimitInKb,omitempty"`
	StoragePoolID         string `json:"storagePoolId"`
	SdsID                 string `json:"sdsId"`
	TestTimeSecs          int    `json:"testTimeSecs,omitempty"`
	TestMode              string `json:"testMode,omitempty"`
}

// DeviceResp defines struct for DeviceParam
type DeviceResp struct {
	ID string `json:"id"`
}

// StoragePool defines struct for PFlex StoragePool
type StoragePool struct {
	ProtectionDomainID                                              string  `json:"protectionDomainId"`
	RebalanceioPriorityPolicy                                       string  `json:"rebalanceIoPriorityPolicy"`
	RebuildioPriorityPolicy                                         string  `json:"rebuildIoPriorityPolicy"`
	RebuildioPriorityBwLimitPerDeviceInKbps                         int     `json:"rebuildIoPriorityBwLimitPerDeviceInKbps"`
	RebuildioPriorityNumOfConcurrentIosPerDevice                    int     `json:"rebuildIoPriorityNumOfConcurrentIosPerDevice"`
	RebalanceioPriorityNumOfConcurrentIosPerDevice                  int     `json:"rebalanceIoPriorityNumOfConcurrentIosPerDevice"`
	RebalanceioPriorityBwLimitPerDeviceInKbps                       int     `json:"rebalanceIoPriorityBwLimitPerDeviceInKbps"`
	RebuildioPriorityAppIopsPerDeviceThreshold                      int     `json:"rebuildIoPriorityAppIopsPerDeviceThreshold"`
	RebalanceioPriorityAppIopsPerDeviceThreshold                    int     `json:"rebalanceIoPriorityAppIopsPerDeviceThreshold"`
	RebuildioPriorityAppBwPerDeviceThresholdInKbps                  int     `json:"rebuildIoPriorityAppBwPerDeviceThresholdInKbps"`
	RebalanceioPriorityAppBwPerDeviceThresholdInKbps                int     `json:"rebalanceIoPriorityAppBwPerDeviceThresholdInKbps"`
	RebuildioPriorityQuietPeriodInMsec                              int     `json:"rebuildIoPriorityQuietPeriodInMsec"`
	RebalanceioPriorityQuietPeriodInMsec                            int     `json:"rebalanceIoPriorityQuietPeriodInMsec"`
	ZeroPaddingEnabled                                              bool    `json:"zeroPaddingEnabled"`
	UseRmcache                                                      bool    `json:"useRmcache"`
	SparePercentage                                                 int     `json:"sparePercentage"`
	RmCacheWriteHandlingMode                                        string  `json:"rmcacheWriteHandlingMode"`
	RebuildEnabled                                                  bool    `json:"rebuildEnabled"`
	RebalanceEnabled                                                bool    `json:"rebalanceEnabled"`
	NumofParallelRebuildRebalanceJobsPerDevice                      int     `json:"numOfParallelRebuildRebalanceJobsPerDevice"`
	Name                                                            string  `json:"name"`
	ID                                                              string  `json:"id"`
	Links                                                           []*Link `json:"links"`
	BackgroundScannerBWLimitKBps                                    int     `json:"backgroundScannerBWLimitKBps"`
	ProtectedMaintenanceModeIoPriorityNumOfConcurrentIosPerDevice   int     `json:"protectedMaintenanceModeIoPriorityNumOfConcurrentIosPerDevice"`
	DataLayout                                                      string  `json:"dataLayout"`
	VtreeMigrationIoPriorityBwLimitPerDeviceInKbps                  int     `json:"vtreeMigrationIoPriorityBwLimitPerDeviceInKbps"`
	VtreeMigrationIoPriorityPolicy                                  string  `json:"vtreeMigrationIoPriorityPolicy"`
	AddressSpaceUsage                                               string  `json:"addressSpaceUsage"`
	ExternalAccelerationType                                        string  `json:"externalAccelerationType"`
	PersistentChecksumState                                         string  `json:"persistentChecksumState"`
	UseRfcache                                                      bool    `json:"useRfcache"`
	ChecksumEnabled                                                 bool    `json:"checksumEnabled"`
	CompressionMethod                                               string  `json:"compressionMethod"`
	FragmentationEnabled                                            bool    `json:"fragmentationEnabled"`
	CapacityUsageState                                              string  `json:"capacityUsageState"`
	CapacityUsageType                                               string  `json:"capacityUsageType"`
	AddressSpaceUsageType                                           string  `json:"addressSpaceUsageType"`
	BgScannerCompareErrorAction                                     string  `json:"bgScannerCompareErrorAction"`
	BgScannerReadErrorAction                                        string  `json:"bgScannerReadErrorAction"`
	ReplicationCapacityMaxRatio                                     int     `json:"replicationCapacityMaxRatio"`
	PersistentChecksumEnabled                                       bool    `json:"persistentChecksumEnabled"`
	PersistentChecksumBuilderLimitKb                                int     `json:"persistentChecksumBuilderLimitKb"`
	PersistentChecksumValidateOnRead                                bool    `json:"persistentChecksumValidateOnRead"`
	VtreeMigrationIoPriorityNumOfConcurrentIosPerDevice             int     `json:"vtreeMigrationIoPriorityNumOfConcurrentIosPerDevice"`
	ProtectedMaintenanceModeIoPriorityPolicy                        string  `json:"protectedMaintenanceModeIoPriorityPolicy"`
	BackgroundScannerMode                                           string  `json:"backgroundScannerMode"`
	MediaType                                                       string  `json:"mediaType"`
	CapacityAlertHighThreshold                                      int     `json:"capacityAlertHighThreshold"`
	CapacityAlertCriticalThreshold                                  int     `json:"capacityAlertCriticalThreshold"`
	VtreeMigrationIoPriorityAppIopsPerDeviceThreshold               int     `json:"vtreeMigrationIoPriorityAppIopsPerDeviceThreshold"`
	VtreeMigrationIoPriorityAppBwPerDeviceThresholdInKbps           int     `json:"vtreeMigrationIoPriorityAppBwPerDeviceThresholdInKbps"`
	VtreeMigrationIoPriorityQuietPeriodInMsec                       int     `json:"vtreeMigrationIoPriorityQuietPeriodInMsec"`
	FglAccpID                                                       string  `json:"fglAccpId"`
	FglExtraCapacity                                                int     `json:"fglExtraCapacity"`
	FglOverProvisioningFactor                                       int     `json:"fglOverProvisioningFactor"`
	FglWriteAtomicitySize                                           int     `json:"fglWriteAtomicitySize"`
	FglNvdimmWriteCacheSizeInMb                                     int     `json:"fglNvdimmWriteCacheSizeInMb"`
	FglNvdimmMetadataAmortizationX100                               int     `json:"fglNvdimmMetadataAmortizationX100"`
	FglPerfProfile                                                  string  `json:"fglPerfProfile"`
	ProtectedMaintenanceModeIoPriorityBwLimitPerDeviceInKbps        int     `json:"protectedMaintenanceModeIoPriorityBwLimitPerDeviceInKbps"`
	ProtectedMaintenanceModeIoPriorityAppIopsPerDeviceThreshold     int     `json:"protectedMaintenanceModeIoPriorityAppIopsPerDeviceThreshold"`
	ProtectedMaintenanceModeIoPriorityAppBwPerDeviceThresholdInKbps int     `json:"protectedMaintenanceModeIoPriorityAppBwPerDeviceThresholdInKbps"`
	ProtectedMaintenanceModeIoPriorityQuietPeriodInMsec             int     `json:"protectedMaintenanceModeIoPriorityQuietPeriodInMsec"`
}

// StoragePoolParam defines struct for StoragePoolParam
type StoragePoolParam struct {
	Name                     string `json:"name"`
	SparePercentage          int    `json:"sparePercentage,omitempty"`
	RebuildEnabled           bool   `json:"rebuildEnabled,omitempty"`
	RebalanceEnabled         bool   `json:"rebalanceEnabled,omitempty"`
	ProtectionDomainID       string `json:"protectionDomainId"`
	ZeroPaddingEnabled       bool   `json:"zeroPaddingEnabled,omitempty"`
	UseRmcache               string `json:"useRmcache,omitempty"`
	UseRfcache               string `json:"useRfcache,omitempty"`
	RmcacheWriteHandlingMode string `json:"rmcacheWriteHandlingMode,omitempty"`
	MediaType                string `json:"mediaType,omitempty"`
}

// ModifyStoragePoolName defines struct for ModifyStoragePoolName
type ModifyStoragePoolName struct {
	Name string `json:"name"`
}

// StoragePoolMediaType defines struct for StoragePoolMediaType
type StoragePoolMediaType struct {
	MediaType string `json:"mediaType"`
}

// StoragePoolUseRmCache defines struct for StoragePoolUseRmCache
type StoragePoolUseRmCache struct {
	UseRmcache string `json:"useRmcache"`
}

// StoragePoolUseRfCache defines struct for StoragePoolUseRfCache
type StoragePoolUseRfCache struct {
}

// StoragePoolResp defines struct for StoragePoolResp
type StoragePoolResp struct {
	ID string `json:"id"`
}

// MappedSdcInfo defines struct for MappedSdcInfo
type MappedSdcInfo struct {
	SdcID                 string `json:"sdcId"`
	SdcIP                 string `json:"sdcIp"`
	LimitIops             int    `json:"limitIops"`
	LimitBwInMbps         int    `json:"limitBwInMbps"`
	SdcName               string `json:"sdcName"`
	AccessMode            string `json:"accessMode"`
	IsDirectBufferMapping bool   `json:"isDirectBufferMapping"`
}

// Volume defines struct for Volume
type Volume struct {
	StoragePoolID                      string           `json:"storagePoolId"`
	UseRmCache                         bool             `json:"useRmcache"`
	MappingToAllSdcsEnabled            bool             `json:"mappingToAllSdcsEnabled"`
	MappedSdcInfo                      []*MappedSdcInfo `json:"mappedSdcInfo"`
	IsObfuscated                       bool             `json:"isObfuscated"`
	VolumeType                         string           `json:"volumeType"`
	ConsistencyGroupID                 string           `json:"consistencyGroupId"`
	VTreeID                            string           `json:"vtreeId"`
	AncestorVolumeID                   string           `json:"ancestorVolumeId"`
	MappedScsiInitiatorInfo            string           `json:"mappedScsiInitiatorInfo"`
	SizeInKb                           int              `json:"sizeInKb"`
	CreationTime                       int              `json:"creationTime"`
	Name                               string           `json:"name"`
	ID                                 string           `json:"id"`
	DataLayout                         string           `json:"dataLayout"`
	NotGenuineSnapshot                 bool             `json:"notGenuineSnapshot"`
	AccessModeLimit                    string           `json:"accessModeLimit"`
	SecureSnapshotExpTime              int              `json:"secureSnapshotExpTime"`
	ManagedBy                          string           `json:"managedBy"`
	LockedAutoSnapshot                 bool             `json:"lockedAutoSnapshot"`
	LockedAutoSnapshotMarkedForRemoval bool             `json:"lockedAutoSnapshotMarkedForRemoval"`
	CompressionMethod                  string           `json:"compressionMethod"`
	TimeStampIsAccurate                bool             `json:"timeStampIsAccurate"`
	OriginalExpiryTime                 int              `json:"originalExpiryTime"`
	VolumeReplicationState             string           `json:"volumeReplicationState"`
	ReplicationJournalVolume           bool             `json:"replicationJournalVolume"`
	ReplicationTimeStamp               int              `json:"replicationTimeStamp"`
	Links                              []*Link          `json:"links"`
}

// VolumeParam defines struct for VolumeParam
type VolumeParam struct {
	ProtectionDomainID string    `json:"protectionDomainId,omitempty"`
	StoragePoolID      string    `json:"storagePoolId,omitempty"`
	UseRmCache         string    `json:"useRmcache,omitempty"`
	VolumeType         string    `json:"volumeType,omitempty"`
	VolumeSizeInKb     string    `json:"volumeSizeInKb,omitempty"`
	Name               string    `json:"name,omitempty"`
	once               sync.Once // creates the metadata value once.
	metadata           http.Header
}

// MetaData returns the metadata headers.
func (vp *VolumeParam) MetaData() http.Header {
	vp.once.Do(func() {
		vp.metadata = make(http.Header)
	})
	return vp.metadata
}

// SetVolumeSizeParam defines struct for SetVolumeSizeParam
type SetVolumeSizeParam struct {
	SizeInGB string `json:"sizeInGB,omitempty"`
}

// SetVolumeNameParam defines struct for SetVolumeNameParam
type SetVolumeNameParam struct {
	NewName string `json:"newName,omitempty"`
}

// VolumeResp defines struct for SetVolumeNameParam
type VolumeResp struct {
	ID string `json:"id"`
}

// VolumeQeryIDByKeyParam defines struct for VolumeQeryIDByKeyParam
type VolumeQeryIDByKeyParam struct {
	Name string `json:"name"`
}

// VolumeQeryBySelectedIdsParam defines struct for VolumeQeryBySelectedIdsParam
type VolumeQeryBySelectedIdsParam struct {
	IDs []string `json:"ids"`
}

// MapVolumeSdcParam defines struct for MapVolumeSdcParam
type MapVolumeSdcParam struct {
	SdcID                 string `json:"sdcId,omitempty"`
	AllowMultipleMappings string `json:"allowMultipleMappings,omitempty"`
	AllSdcs               string `json:"allSdcs,omitempty"`
}

// UnmapVolumeSdcParam defines struct for UnmapVolumeSdcParam
type UnmapVolumeSdcParam struct {
	SdcID                string `json:"sdcId,omitempty"`
	IgnoreScsiInitiators string `json:"ignoreScsiInitiators,omitempty"`
	AllSdcs              string `json:"allSdcs,omitempty"`
}

// SetMappedSdcLimitsParam defines struct for SetMappedSdcLimitsParam
type SetMappedSdcLimitsParam struct {
	SdcID                string `json:"sdcId,omitempty"`
	BandwidthLimitInKbps string `json:"bandwidthLimitInKbps,omitempty"`
	IopsLimit            string `json:"iopsLimit,omitempty"`
}

// RenameSdcParam defines struct for RenameSdc
type RenameSdcParam struct {
	SdcName string `json:"sdcName,omitempty"`
}

// SnapshotDef defines struct for SnapshotDef
type SnapshotDef struct {
	VolumeID     string `json:"volumeId,omitempty"`
	SnapshotName string `json:"snapshotName,omitempty"`
}

// SnapshotVolumesParam defines struct for SnapshotVolumesParam
type SnapshotVolumesParam struct {
	SnapshotDefs         []*SnapshotDef `json:"snapshotDefs"`
	RetentionPeriodInMin string         `json:"retentionPeriodInMin,omitempty"`
	AccessMode           string         `json:"accessModeLimit,omitempty"`
	AllowOnExtManagedVol bool           `json:"allowOnExtManagedVol,omitempty"`
}

// SnapshotVolumesResp defines struct for SnapshotVolumesResp
type SnapshotVolumesResp struct {
	VolumeIDList    []string `json:"volumeIdList"`
	SnapshotGroupID string   `json:"snapshotGroupId"`
}

// VTree defines struct for VTree
type VTree struct {
	ID            string  `json:"id"`
	Name          string  `json:"name"`
	BaseVolumeID  string  `json:"baseVolumeId"`
	StoragePoolID string  `json:"storagePoolId"`
	Links         []*Link `json:"links"`
}

// RemoveVolumeParam defines struct for RemoveVolumeParam
type RemoveVolumeParam struct {
	RemoveMode string `json:"removeMode"`
}

// EmptyPayload defines struct for EmptyPayload
type EmptyPayload struct {
}

// SnapshotPolicy defines the struct for SnapshotPolicy
type SnapshotPolicy struct {
	SnapshotPolicyState                   string  `json:"snapshotPolicyState"`
	AutoSnapshotCreationCadenceInMin      int     `json:"autoSnapshotCreationCadenceInMin"`
	MaxVTreeAutoSnapshots                 int     `json:"maxVTreeAutoSnapshots"`
	NumOfSourceVolumes                    int     `json:"numOfSourceVolumes"`
	NumOfExpiredButLockedSnapshots        int     `json:"numOfExpiredButLockedSnapshots"`
	NumOfCreationFailures                 int     `json:"numOfCreationFailures"`
	NumOfRetainedSnapshotsPerLevel        []int   `json:"numOfRetainedSnapshotsPerLevel"`
	SnapshotAccessMode                    string  `json:"snapshotAccessMode"`
	SecureSnapshots                       bool    `json:"secureSnapshots"`
	TimeOfLastAutoSnapshot                int     `json:"timeOfLastAutoSnapshot"`
	NextAutoSnapshotCreationTime          int     `json:"nextAutoSnapshotCreationTime"`
	TimeOfLastAutoSnapshotCreationFailure int     `json:"timeOfLastAutoSnapshotCreationFailure"`
	LastAutoSnapshotCreationFailureReason string  `json:"lastAutoSnapshotCreationFailureReason"`
	LastAutoSnapshotFailureInFirstLevel   bool    `json:"lastAutoSnapshotFailureInFirstLevel"`
	NumOfAutoSnapshots                    int     `json:"numOfAutoSnapshots"`
	NumOfLockedSnapshots                  int     `json:"numOfLockedSnapshots"`
	SystemID                              string  `json:"systemId"`
	Name                                  string  `json:"name"`
	ID                                    string  `json:"id"`
	Links                                 []*Link `json:"links"`
}

// SnapshotPolicyQueryIDByKeyParam defines struct for SnapshotPolicyQueryIDByKeyParam
type SnapshotPolicyQueryIDByKeyParam struct {
	Name string `json:"name"`
}
