package goscaleio

type NodeDetails struct {
	RefID string `json:"refId"`
	// RefType             interface{}   `json:"refType"`
	IPAddress          string `json:"ipAddress"`
	CurrentIPAddress   string `json:"currentIpAddress"`
	ServiceTag         string `json:"serviceTag"`
	Model              string `json:"model"`
	DeviceType         string `json:"deviceType"`
	DiscoverDeviceType string `json:"discoverDeviceType"`
	DisplayName        string `json:"displayName"`
	ManagedState       string `json:"managedState"`
	State              string `json:"state"`
	InUse              bool   `json:"inUse"`
	// ServiceReferences   []interface{} `json:"serviceReferences"`
	// StatusMessage       interface{} `json:"statusMessage"`
	// FirmwareName        interface{} `json:"firmwareName"`
	CustomFirmware  bool   `json:"customFirmware"`
	NeedsAttention  bool   `json:"needsAttention"`
	Manufacturer    string `json:"manufacturer"`
	SystemID        string `json:"systemId"`
	Health          string `json:"health"`
	HealthMessage   string `json:"healthMessage"`
	OperatingSystem string `json:"operatingSystem"`
	NumberOfCPUs    int    `json:"numberOfCPUs"`
	// CPUType             interface{} `json:"cpuType"`
	Nics       int `json:"nics"`
	MemoryInGB int `json:"memoryInGB"`
	// InfraTemplateDate   interface{} `json:"infraTemplateDate"`
	// InfraTemplateID     interface{} `json:"infraTemplateId"`
	// ServerTemplateDate  interface{} `json:"serverTemplateDate"`
	// ServerTemplateID    interface{} `json:"serverTemplateId"`
	// InventoryDate       interface{} `json:"inventoryDate"`
	ComplianceCheckDate string          `json:"complianceCheckDate"`
	DiscoveredDate      string          `json:"discoveredDate"`
	DeviceGroupList     DeviceGroupList `json:"deviceGroupList"`
	DetailLink          DetailLink      `json:"detailLink"`
	CredID              string          `json:"credId"`
	Compliance          string          `json:"compliance"`
	FailuresCount       int             `json:"failuresCount"`
	// ChassisID           interface{}     `json:"chassisId"`
	Facts string `json:"facts"`
	// ParsedFacts        interface{}   `json:"parsedFacts"`
	// Config             interface{}   `json:"config"`
	// Hostname           interface{}   `json:"hostname"`
	// OsIPAddress        interface{}   `json:"osIpAddress"`
	// OsAdminCredential  interface{}   `json:"osAdminCredential"`
	// OsImageType        interface{}   `json:"osImageType"`
	// LastJobs           interface{}   `json:"lastJobs"`
	PuppetCertName string `json:"puppetCertName"`
	// SvmAdminCredential interface{}   `json:"svmAdminCredential"`
	// SvmName            interface{}   `json:"svmName"`
	// SvmIPAddress       interface{}   `json:"svmIpAddress"`
	// SvmImageType       interface{}   `json:"svmImageType"`
	FlexosMaintMode int `json:"flexosMaintMode"`
	EsxiMaintMode   int `json:"esxiMaintMode"`
	// VMList             []interface{} `json:"vmList"`
}

type DeviceGroupList struct {
	// Paging      interface{} `json:"paging"`
	DeviceGroup []DeviceGroup `json:"deviceGroup"`
}

type DeviceGroup struct {
	// Link              interface{} `json:"link"`
	GroupSeqID       int    `json:"groupSeqId"`
	GroupName        string `json:"groupName"`
	GroupDescription string `json:"groupDescription"`
	CreatedDate      string `json:"createdDate"`
	CreatedBy        string `json:"createdBy"`
	UpdatedDate      string `json:"updatedDate"`
	UpdatedBy        string `json:"updatedBy"`
	// ManagedDeviceList interface{} `json:"managedDeviceList"`
	// GroupUserList     interface{} `json:"groupUserList"`
}

type DetailLink struct {
	Title string `json:"title"`
	Href  string `json:"href"`
	Rel   string `json:"rel"`
	// Type  interface{} `json:"type"`
}

type NodePoolDetails struct {
	DeviceGroupList
}
