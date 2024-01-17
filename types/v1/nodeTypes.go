package goscaleio

type NodeDetails struct {
	RefID               string          `json:"refId"`
	IPAddress           string          `json:"ipAddress"`
	CurrentIPAddress    string          `json:"currentIpAddress"`
	ServiceTag          string          `json:"serviceTag"`
	Model               string          `json:"model"`
	DeviceType          string          `json:"deviceType"`
	DiscoverDeviceType  string          `json:"discoverDeviceType"`
	DisplayName         string          `json:"displayName"`
	ManagedState        string          `json:"managedState"`
	State               string          `json:"state"`
	InUse               bool            `json:"inUse"`
	CustomFirmware      bool            `json:"customFirmware"`
	NeedsAttention      bool            `json:"needsAttention"`
	Manufacturer        string          `json:"manufacturer"`
	SystemID            string          `json:"systemId"`
	Health              string          `json:"health"`
	HealthMessage       string          `json:"healthMessage"`
	OperatingSystem     string          `json:"operatingSystem"`
	NumberOfCPUs        int             `json:"numberOfCPUs"`
	Nics                int             `json:"nics"`
	MemoryInGB          int             `json:"memoryInGB"`
	ComplianceCheckDate string          `json:"complianceCheckDate"`
	DiscoveredDate      string          `json:"discoveredDate"`
	DeviceGroupList     DeviceGroupList `json:"deviceGroupList"`
	DetailLink          DetailLink      `json:"detailLink"`
	CredID              string          `json:"credId"`
	Compliance          string          `json:"compliance"`
	FailuresCount       int             `json:"failuresCount"`
	Facts               string          `json:"facts"`
	PuppetCertName      string          `json:"puppetCertName"`
	FlexosMaintMode     int             `json:"flexosMaintMode"`
	EsxiMaintMode       int             `json:"esxiMaintMode"`
}

type DeviceGroupList struct {
	DeviceGroup    []DeviceGroup `json:"deviceGroup"`
	ManagedDevices []NodeDetails `json:"managedDevices"`
}

type DeviceGroup struct {
	GroupSeqID       int           `json:"groupSeqId"`
	GroupName        string        `json:"groupName"`
	GroupDescription string        `json:"groupDescription"`
	CreatedDate      string        `json:"createdDate"`
	CreatedBy        string        `json:"createdBy"`
	UpdatedDate      string        `json:"updatedDate"`
	UpdatedBy        string        `json:"updatedBy"`
	GroupUserList    GroupUserList `json:"groupUserList"`
}

type GroupUserList struct {
	TotalRecords int          `json:"totalRecords"`
	GroupUsers   []GroupUsers `json:"groupUsers"`
}

type GroupUsers struct {
	UserSeqID string `json:"userSeqId"`
	UserName  string `json:"userName"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Role      string `json:"role"`
	Enabled   bool   `json:"enabled"`
}

type DetailLink struct {
	Title string `json:"title"`
	Href  string `json:"href"`
	Rel   string `json:"rel"`
}

type ManagedDeviceList struct {
	ManagedDevices []NodeDetails `json:"managedDevices"`
}

type NodePoolDetails struct {
	DeviceGroup       DeviceGroup       `json:"deviceGroup"`
	ManagedDeviceList ManagedDeviceList `json:"managedDeviceList"`
	GroupUserList     GroupUserList     `json:"groupUserList"`
}
