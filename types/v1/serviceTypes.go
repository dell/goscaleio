package goscaleio

type ServiceFailedResponse struct {
	DetailMessage string     `json:"detailMessage,omitempty"`
	Status        int        `json:"status,omitempty"`
	Timestamp     string     `json:"timestamp,omitempty"`
	Error         string     `json:"error,omitempty"`
	Path          string     `json:"path,omitempty"`
	Messages      []Messages `json:"messages,omitempty"`
}

type DeploymentPayload struct {
	DeploymentName        string          `json:"deploymentName,omitempty"`
	DeploymentDescription string          `json:"deploymentDescription,omitempty"`
	ServiceTemplate       TemplateDetails `json:"serviceTemplate,omitempty"`
	UpdateServerFirmware  bool            `json:"updateServerFirmware,omitempty"`
	FirmwareRepositoryID  string          `json:"firmwareRepositoryId,omitempty"`
	Status                string          `json:"status,omitempty"`
}

type ServiceResponse struct {
	ID                           string                       `json:"id,omitempty"`
	DeploymentName               string                       `json:"deploymentName,omitempty"`
	DeploymentDescription        string                       `json:"deploymentDescription,omitempty"`
	DeploymentValid              DeploymentValid              `json:"deploymentValid,omitempty"`
	Retry                        bool                         `json:"retry,omitempty"`
	Teardown                     bool                         `json:"teardown,omitempty"`
	TeardownAfterCancel          bool                         `json:"teardownAfterCancel,omitempty"`
	RemoveService                bool                         `json:"removeService,omitempty"`
	CreatedDate                  string                       `json:"createdDate,omitempty"`
	CreatedBy                    string                       `json:"createdBy,omitempty"`
	UpdatedDate                  string                       `json:"updatedDate,omitempty"`
	UpdatedBy                    string                       `json:"updatedBy,omitempty"`
	DeploymentScheduledDate      string                       `json:"deploymentScheduledDate,omitempty"`
	DeploymentStartedDate        string                       `json:"deploymentStartedDate,omitempty"`
	DeploymentFinishedDate       string                       `json:"deploymentFinishedDate,omitempty"`
	ScheduleDate                 string                       `json:"scheduleDate,omitempty"`
	Status                       string                       `json:"status,omitempty"`
	Compliant                    bool                         `json:"compliant,omitempty"`
	DeploymentDevice             []DeploymentDevice           `json:"deploymentDevice,omitempty"`
	Vms                          []Vms                        `json:"vms,omitempty"`
	UpdateServerFirmware         bool                         `json:"updateServerFirmware,omitempty"`
	UseDefaultCatalog            bool                         `json:"useDefaultCatalog,omitempty"`
	FirmwareRepository           FirmwareRepository           `json:"firmwareRepository,omitempty"`
	FirmwareRepositoryID         string                       `json:"firmwareRepositoryId,omitempty"`
	LicenseRepository            LicenseRepository            `json:"licenseRepository,omitempty"`
	LicenseRepositoryID          string                       `json:"licenseRepositoryId,omitempty"`
	IndividualTeardown           bool                         `json:"individualTeardown,omitempty"`
	DeploymentHealthStatusType   string                       `json:"deploymentHealthStatusType,omitempty"`
	AssignedUsers                []AssignedUsers              `json:"assignedUsers,omitempty"`
	AllUsersAllowed              bool                         `json:"allUsersAllowed,omitempty"`
	Owner                        string                       `json:"owner,omitempty"`
	NoOp                         bool                         `json:"noOp,omitempty"`
	FirmwareInit                 bool                         `json:"firmwareInit,omitempty"`
	DisruptiveFirmware           bool                         `json:"disruptiveFirmware,omitempty"`
	PreconfigureSVM              bool                         `json:"preconfigureSVM,omitempty"`
	PreconfigureSVMAndUpdate     bool                         `json:"preconfigureSVMAndUpdate,omitempty"`
	ServicesDeployed             string                       `json:"servicesDeployed,omitempty"`
	PrecalculatedDeviceHealth    string                       `json:"precalculatedDeviceHealth,omitempty"`
	LifecycleModeReasons         []string                     `json:"lifecycleModeReasons,omitempty"`
	JobDetails                   []JobDetails                 `json:"jobDetails,omitempty"`
	NumberOfDeployments          int                          `json:"numberOfDeployments,omitempty"`
	OperationType                string                       `json:"operationType,omitempty"`
	OperationStatus              string                       `json:"operationStatus,omitempty"`
	OperationData                string                       `json:"operationData,omitempty"`
	DeploymentValidationResponse DeploymentValidationResponse `json:"deploymentValidationResponse,omitempty"`
	CurrentStepCount             string                       `json:"currentStepCount,omitempty"`
	TotalNumOfSteps              string                       `json:"totalNumOfSteps,omitempty"`
	CurrentStepMessage           string                       `json:"currentStepMessage,omitempty"`
	CustomImage                  string                       `json:"customImage,omitempty"`
	OriginalDeploymentID         string                       `json:"originalDeploymentId,omitempty"`
	CurrentBatchCount            string                       `json:"currentBatchCount,omitempty"`
	TotalBatchCount              string                       `json:"totalBatchCount,omitempty"`
	Brownfield                   bool                         `json:"brownfield,omitempty"`
	OverallDeviceHealth          string                       `json:"overallDeviceHealth,omitempty"`
	Vds                          bool                         `json:"vds,omitempty"`
	ScaleUp                      bool                         `json:"scaleUp,omitempty"`
	LifecycleMode                bool                         `json:"lifecycleMode,omitempty"`
	CanMigratevCLSVMs            bool                         `json:"canMigratevCLSVMs,omitempty"`
	TemplateValid                bool                         `json:"templateValid,omitempty"`
	ConfigurationChange          bool                         `json:"configurationChange,omitempty"`
	DetailMessage                string                       `json:"detailMessage,omitempty"`
	Timestamp                    string                       `json:"timestamp,omitempty"`
	Error                        string                       `json:"error,omitempty"`
	Path                         string                       `json:"path,omitempty"`
	Messages                     []Messages                   `json:"messages,omitempty"`
}
