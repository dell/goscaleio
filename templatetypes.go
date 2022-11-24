package goscaleio

type ParameterHelper struct {
}

func StringPtr(v string) *string {
	return &v
}

// TemplateParam defines struct for Payload send to create Template
type TemplateParam struct {
	TemplateName              string                    `json:"templateName"`
	TemplateDescription       string                    `json:"templateDescription"`
	TemplateType              string                    `json:"templateType"`
	Draft                     bool                      `json:"draft"`
	Components                []TemplateComponent       `json:"components"`
	Category                  *string                   `json:"category"`
	AllUsersAllowed           bool                      `json:"allUsersAllowed"`
	AssignedUsers             *string                   `json:"assignedUsers"`
	ManageFirmware            bool                      `json:"manageFirmware"`
	UseDefaultCatalog         bool                      `json:"useDefaultCatalog"`
	FirmwareRepository        FirmwareRepository        `json:"firmwareRepository"`
	LicenseRepository         *string                   `json:"licenseRepository"`
	Configuration             *string                   `json:"configuration"`
	ServerCount               *string                   `json:"serverCount"`
	StorageCount              *string                   `json:"storageCount"`
	ClusterCount              *string                   `json:"clusterCount"`
	ServiceCount              *string                   `json:"serviceCount"`
	SwitchCount               *string                   `json:"switchCount"`
	VmCount                   *string                   `json:"vmCount"`
	SdnasCount                *string                   `json:"sdnasCount"`
	BrownfieldTemplateType    string                    `json:"brownfieldTemplateType"`
	Networks                  *string                   `json:"networks"`
	BlockServiceOperationsMap BlockServiceOperationsMap `json:"blockServiceOperationsMap"`
}

// firmwareRepository defines struct for firmwareRepository
type FirmwareRepository struct {
	Id                      string   `json:"id"`
	Name                    *string  `json:"name"`
	SourceLocation          *string  `json:"sourceLocation"`
	SourceType              *string  `json:"sourceType"`
	DiskLocation            *string  `json:"diskLocation"`
	Filename                *string  `json:"filename"`
	Md5Hash                 *string  `json:"md5Hash"`
	Username                *string  `json:"username"`
	Password                *string  `json:"password"`
	DownloadStatus          *string  `json:"downloadStatus"`
	CreatedDate             *string  `json:"createdDate"`
	CreatedBy               *string  `json:"createdBy"`
	UpdatedDate             *string  `json:"updatedDate"`
	UpdatedBy               *string  `json:"updatedBy"`
	DefaultCatalog          bool     `json:"defaultCatalog"`
	Embedded                bool     `json:"embedded"`
	State                   *string  `json:"state"`
	SoftwareComponents      []string `json:"softwareComponents"`
	SoftwareBundles         []string `json:"softwareBundles"`
	Deployments             []string `json:"deployments"`
	BundleCount             int      `json:"bundleCount"`
	ComponentCount          int      `json:"componentCount"`
	UserBundleCount         int      `json:"userBundleCount"`
	Minimal                 bool     `json:"minimal"`
	DownloadProgress        int      `json:"downloadProgress"`
	ExtractProgress         int      `json:"extractProgress"`
	FileSizeInGigabytes     *string  `json:"fileSizeInGigabytes"`
	SignedKeySourceLocation *string  `json:"signedKeySourceLocation"`
	Signature               *string  `json:"signature"`
	Rcmapproved             bool     `json:"rcmapproved"`
}

// BlockServiceOperationsMap defines struct for BlockServiceOperationsMap
type BlockServiceOperationsMap struct {
}

// "guid": null,
// "id": "asm::scaleio::cloudlink",
// "displayName": "Cloud Link Center Settings",

type ComponentResources struct {
	Id          *string              `json:"id"`
	Guid        string               `json:"guid"`
	DisplayName string               `json:"displayName"`
	Parameters  []ResourceParameters `json:"parameters"`
}

type ResourceParameters struct {
	Guid                     *string     `json:"guid"`
	Id                       string      `json:"id"`
	Value                    string      `json:"value"`
	Type                     string      `json:"type"`
	DisplayName              string      `json:"displayName"`
	Required                 bool        `json:"required"`
	RequiredAtDeployment     bool        `json:"requiredAtDeployment"`
	HideFromTemplate         bool        `json:"hideFromTemplate"`
	Min                      *string     `json:"min"`
	Max                      *string     `json:"max"`
	DependencyTarget         *string     `json:"dependencyTarget"`
	DependencyValue          *string     `json:"dependencyValue"`
	Dependencies             []string    `json:"dependencies"`
	Networks                 *string     `json:"networks"`
	NetworkIpAddressList     *string     `json:"networkIpAddressList"`
	NetworkConfiguration     *string     `json:"networkConfiguration"`
	RaidConfiguration        *string     `json:"raidConfiguration"`
	Options                  []string    `json:"options"`
	ToolTip                  string      `json:"toolTip"`
	ReadOnly                 bool        `json:"readOnly"`
	Generated                bool        `json:"generated"`
	Group                    *string     `json:"group"`
	InfoIcon                 bool        `json:"infoIcon"`
	MaxLength                int         `json:"maxLength"`
	Step                     int         `json:"step"`
	OptionsSortable          bool        `json:"optionsSortable"`
	PreservedForDeployment   bool        `json:"preservedForDeployment"`
	ScaleIODiskConfiguration *string     `json:"scaleIODiskConfiguration"`
	ProtectionDomainSettings *string     `json:"protectionDomainSettings"`
	FaultSetSettings         *string     `json:"faultSetSettings"`
	Attributes               interface{} `json:"attributes"`
	VdsConfiguration         *string     `json:"vdsConfiguration"`
	NodeSelection            *string     `json:"nodeSelection"`
}

// TemplateComponent defines struct for TemplateComponent
type TemplateComponent struct {
	Id                string               `json:"id"`
	ComponentID       string               `json:"componentID"`
	Name              string               `json:"name"`
	Type              string               `json:"type"`
	SubType           string               `json:"subType"`
	Resources         []ComponentResources `json:"resources"`
	RefId             *string              `json:"refId"`
	Cloned            bool                 `json:"cloned"`
	ClonedFromId      *string              `json:"clonedFromId"`
	ManageFirmware    bool                 `json:"manageFirmware"`
	Brownfield        bool                 `json:"brownfield"`
	Instances         int                  `json:"instances"`
	ClonedFromAsmGuid *string              `json:"clonedFromAsmGuid"`
	Ip                *string              `json:"ip"`
}

type TemplateComponentGen struct {
	ID          string `json:"id"`
	ComponentID string `json:"componentID"`
	Name        string `json:"name"`
	Type        string `json:"type"`
	SubType     string `json:"subType"`
	Resources   []struct {
		GUID        interface{} `json:"guid"`
		ID          string      `json:"id"`
		DisplayName string      `json:"displayName"`
		Parameters  []struct {
			GUID                 interface{}   `json:"guid"`
			ID                   string        `json:"id"`
			Value                string        `json:"value"`
			Type                 string        `json:"type"`
			DisplayName          string        `json:"displayName"`
			Required             bool          `json:"required"`
			RequiredAtDeployment bool          `json:"requiredAtDeployment"`
			HideFromTemplate     bool          `json:"hideFromTemplate"`
			Min                  interface{}   `json:"min"`
			Max                  interface{}   `json:"max"`
			DependencyTarget     interface{}   `json:"dependencyTarget"`
			DependencyValue      interface{}   `json:"dependencyValue"`
			Dependencies         []interface{} `json:"dependencies"`
			Networks             interface{}   `json:"networks"`
			NetworkIPAddressList interface{}   `json:"networkIpAddressList"`
			NetworkConfiguration interface{}   `json:"networkConfiguration"`
			RaidConfiguration    interface{}   `json:"raidConfiguration"`
			Options              []struct {
				ID               interface{}   `json:"id"`
				Value            string        `json:"value"`
				Name             string        `json:"name"`
				DependencyTarget interface{}   `json:"dependencyTarget"`
				DependencyValue  interface{}   `json:"dependencyValue"`
				Dependencies     []interface{} `json:"dependencies"`
				Attributes       struct {
				} `json:"attributes"`
			} `json:"options"`
			ToolTip                  string      `json:"toolTip"`
			ReadOnly                 bool        `json:"readOnly"`
			Generated                bool        `json:"generated"`
			Group                    interface{} `json:"group"`
			InfoIcon                 bool        `json:"infoIcon"`
			MaxLength                int         `json:"maxLength"`
			Step                     int         `json:"step"`
			OptionsSortable          bool        `json:"optionsSortable"`
			PreservedForDeployment   bool        `json:"preservedForDeployment"`
			ScaleIODiskConfiguration interface{} `json:"scaleIODiskConfiguration"`
			ProtectionDomainSettings interface{} `json:"protectionDomainSettings"`
			FaultSetSettings         interface{} `json:"faultSetSettings"`
			Attributes               struct {
			} `json:"attributes"`
			VdsConfiguration interface{} `json:"vdsConfiguration"`
			NodeSelection    interface{} `json:"nodeSelection"`
		} `json:"parameters"`
	} `json:"resources"`
	RefID             interface{} `json:"refId"`
	Cloned            bool        `json:"cloned"`
	ClonedFromID      interface{} `json:"clonedFromId"`
	ManageFirmware    bool        `json:"manageFirmware"`
	Brownfield        bool        `json:"brownfield"`
	Instances         int         `json:"instances"`
	ClonedFromAsmGUID interface{} `json:"clonedFromAsmGuid"`
	IP                interface{} `json:"ip"`
}
