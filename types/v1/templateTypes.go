// Copyright © 2024 Dell Inc. or its subsidiaries. All Rights Reserved.
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

import "time"

// TemplateDetails defines struct for Template
type TemplateDetails struct {
	ID                        string                    `json:"id,omitempty"`
	TemplateName              string                    `json:"templateName,omitempty"`
	TemplateDescription       string                    `json:"templateDescription,omitempty"`
	TemplateType              string                    `json:"templateType,omitempty"`
	TemplateVersion           string                    `json:"templateVersion,omitempty"`
	TemplateValid             TemplateValid             `json:"templateValid,omitempty"`
	OriginalTemplateID        any                       `json:"originalTemplateId,omitempty"`
	TemplateLocked            bool                      `json:"templateLocked,omitempty"`
	Draft                     bool                      `json:"draft,omitempty"`
	InConfiguration           bool                      `json:"inConfiguration,omitempty"`
	CreatedDate               time.Time                 `json:"createdDate,omitempty"`
	CreatedBy                 string                    `json:"createdBy,omitempty"`
	UpdatedDate               time.Time                 `json:"updatedDate,omitempty"`
	LastDeployedDate          time.Time                 `json:"lastDeployedDate,omitempty"`
	UpdatedBy                 string                    `json:"updatedBy,omitempty"`
	Components                []Components              `json:"components,omitempty"`
	Category                  string                    `json:"category,omitempty"`
	AllUsersAllowed           bool                      `json:"allUsersAllowed,omitempty"`
	AssignedUsers             []any                     `json:"assignedUsers,omitempty"`
	ManageFirmware            bool                      `json:"manageFirmware,omitempty"`
	UseDefaultCatalog         bool                      `json:"useDefaultCatalog,omitempty"`
	FirmwareRepository        FirmwareRepository        `json:"firmwareRepository,omitempty"`
	LicenseRepository         any                       `json:"licenseRepository,omitempty"`
	Configuration             any                       `json:"configuration,omitempty"`
	ServerCount               int                       `json:"serverCount,omitempty"`
	StorageCount              int                       `json:"storageCount,omitempty"`
	ClusterCount              int                       `json:"clusterCount,omitempty"`
	ServiceCount              int                       `json:"serviceCount,omitempty"`
	SwitchCount               int                       `json:"switchCount,omitempty"`
	VMCount                   int                       `json:"vmCount,omitempty"`
	SdnasCount                int                       `json:"sdnasCount,omitempty"`
	BrownfieldTemplateType    string                    `json:"brownfieldTemplateType,omitempty"`
	Networks                  []Networks                `json:"networks,omitempty"`
	BlockServiceOperationsMap BlockServiceOperationsMap `json:"blockServiceOperationsMap,omitempty"`
}

// TemplateValid defines struct for TemplateValid
type TemplateValid struct {
	Valid    bool     `json:"valid,omitempty"`
	Messages []string `json:"messages,omitempty"`
}

// ComponentValid defines struct for ComponentValid
type ComponentValid struct {
	Valid    bool     `json:"valid,omitempty"`
	Messages []string `json:"messages,omitempty"`
}

// RelatedComponents defines struct for RelatedComponents
type RelatedComponents struct {
	NAMING_FAILED string `json:",omitempty"`
}

// Components defines struct for RelatedComponents
type Components struct {
	ID                  string            `json:"id,omitempty"`
	ComponentID         string            `json:"componentID,omitempty"`
	Identifier          string            `json:"identifier,omitempty"`
	ComponentValid      ComponentValid    `json:"componentValid,omitempty"`
	PuppetCertName      string            `json:"puppetCertName,omitempty"`
	OsPuppetCertName    string            `json:"osPuppetCertName,omitempty"`
	Name                string            `json:"name,omitempty"`
	Type                string            `json:"type,omitempty"`
	SubType             string            `json:"subType,omitempty"`
	Teardown            bool              `json:"teardown,omitempty"`
	HelpText            any               `json:"helpText,omitempty"`
	ManagementIPAddress string            `json:"managementIpAddress,omitempty"`
	ConfigFile          any               `json:"configFile,omitempty"`
	SerialNumber        string            `json:"serialNumber,omitempty"`
	AsmGUID             string            `json:"asmGUID,omitempty"`
	RelatedComponents   RelatedComponents `json:"relatedComponents,omitempty"`
	Resources           []any             `json:"resources,omitempty"`
	RefID               string            `json:"refId,omitempty"`
	Cloned              bool              `json:"cloned,omitempty"`
	ClonedFromID        any               `json:"clonedFromId,omitempty"`
	ManageFirmware      bool              `json:"manageFirmware,omitempty"`
	Brownfield          bool              `json:"brownfield,omitempty"`
	Instances           int               `json:"instances,omitempty"`
	ClonedFromAsmGUID   string            `json:"clonedFromAsmGuid,omitempty"`
	IP                  string            `json:"ip,omitempty"`
}

// FirmwareRepository defines struct for FirmwareRepository
type FirmwareRepository struct {
	ID                      string `json:"id,omitempty"`
	Name                    string `json:"name,omitempty"`
	SourceLocation          any    `json:"sourceLocation,omitempty"`
	SourceType              any    `json:"sourceType,omitempty"`
	DiskLocation            any    `json:"diskLocation,omitempty"`
	Filename                any    `json:"filename,omitempty"`
	Md5Hash                 any    `json:"md5Hash,omitempty"`
	Username                any    `json:"username,omitempty"`
	Password                any    `json:"password,omitempty"`
	DownloadStatus          any    `json:"downloadStatus,omitempty"`
	CreatedDate             any    `json:"createdDate,omitempty"`
	CreatedBy               any    `json:"createdBy,omitempty"`
	UpdatedDate             any    `json:"updatedDate,omitempty"`
	UpdatedBy               any    `json:"updatedBy,omitempty"`
	DefaultCatalog          bool   `json:"defaultCatalog,omitempty"`
	Embedded                bool   `json:"embedded,omitempty"`
	State                   any    `json:"state,omitempty"`
	SoftwareComponents      []any  `json:"softwareComponents,omitempty"`
	SoftwareBundles         []any  `json:"softwareBundles,omitempty"`
	Deployments             []any  `json:"deployments,omitempty"`
	BundleCount             int    `json:"bundleCount,omitempty"`
	ComponentCount          int    `json:"componentCount,omitempty"`
	UserBundleCount         int    `json:"userBundleCount,omitempty"`
	Minimal                 bool   `json:"minimal,omitempty"`
	DownloadProgress        int    `json:"downloadProgress,omitempty"`
	ExtractProgress         int    `json:"extractProgress,omitempty"`
	FileSizeInGigabytes     any    `json:"fileSizeInGigabytes,omitempty"`
	SignedKeySourceLocation any    `json:"signedKeySourceLocation,omitempty"`
	Signature               any    `json:"signature,omitempty"`
	Custom                  bool   `json:"custom,omitempty"`
	NeedsAttention          bool   `json:"needsAttention,omitempty"`
	JobID                   any    `json:"jobId,omitempty"`
	Rcmapproved             bool   `json:"rcmapproved,omitempty"`
}

// StaticNetworkConfiguration defines struct for StaticNetworkConfiguration
type StaticNetworkConfiguration struct {
	Gateway      string    `json:"gateway,omitempty"`
	Subnet       string    `json:"subnet,omitempty"`
	PrimaryDNS   string    `json:"primaryDns,omitempty"`
	SecondaryDNS string    `json:"secondaryDns,omitempty"`
	DNSSuffix    string    `json:"dnsSuffix,omitempty"`
	IPRange      []IPRange `json:"ipRange,omitempty"`
	IPAddress    string    `json:"ipAddress,omitempty"`
	StaticRoute  string    `json:"staticRoute,omitempty"`
}

// Networks defines struct for Networks
type Networks struct {
	ID                         string                     `json:"id,omitempty"`
	Name                       string                     `json:"name,omitempty"`
	Description                string                     `json:"description,omitempty"`
	Type                       string                     `json:"type,omitempty"`
	VlanID                     int                        `json:"vlanId,omitempty"`
	Static                     bool                       `json:"static,omitempty"`
	StaticNetworkConfiguration StaticNetworkConfiguration `json:"staticNetworkConfiguration,omitempty"`
	DestinationIPAddress       string                     `json:"destinationIpAddress,omitempty"`
}

// IPRange defines struct for IPRange
type IPRange struct {
	ID         string `json:"id"`
	StartingIP string `json:"startingIp"`
	EndingIP   string `json:"endingIp"`
	Role       any    `json:"role"`
}

// BlockServiceOperationsMap defines struct for BlockServiceOperationsMap
type BlockServiceOperationsMap struct {
}

// TemplateDetailsFilter defines struct for nodepools
type TemplateDetailsFilter struct {
	TemplateDetails []TemplateDetails `json:"serviceTemplate"`
}
