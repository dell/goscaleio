package inttests

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestDeployUploadPackage function to test upload packge with dummy path of packages
func TestDeployUploadPackage(t *testing.T) {
	var filePaths []string
	filePaths = append(filePaths, "/home/abc.txt")
	_, err := GC.UploadPackages(filePaths)
	assert.NotNil(t, err)
}

// TestDeployParseCSV function to test parse csv function with dummy path of CSV file
func TestDeployParseCSV(t *testing.T) {
	_, err := GC.ParseCSV("/test/test.csv")
	assert.NotNil(t, err)
}

// TestDeployGetPackage function to test Get Packge Details
func TestDeployGetPackgeDetails(t *testing.T) {
	res, err := GC.GetPackageDetails()
	assert.NotNil(t, res)
	assert.Nil(t, err)
}

// TestDeployValidateMDMDetails function to test Retrival of MDM Topology Function
func TestDeployValidateMDMDetails(t *testing.T) {
	mapData := map[string]interface{}{
		"mdmUser":     "admin",
		"mdmPassword": "Password123",
	}
	mapData["mdmIps"] = []string{"10.247.103.161"}

	secureData := map[string]interface{}{
		"allowNonSecureCommunicationWithMdm": true,
		"allowNonSecureCommunicationWithLia": true,
		"disableNonMgmtComponentsAuth":       false,
	}
	mapData["securityConfiguration"] = secureData

	jsonres, _ := json.Marshal(mapData)

	res, err := GC.ValidateMDMDetails(jsonres)

	assert.NotNil(t, res)

	assert.EqualValues(t, res.StatusCode, 200)

	assert.Nil(t, err)
}

func TestDeployGetClusterDetails(t *testing.T) {
	mapData := map[string]interface{}{
		"mdmUser":     "admin",
		"mdmPassword": "Password123",
	}
	mapData["mdmIps"] = []string{"10.247.103.164"}

	secureData := map[string]interface{}{
		"allowNonSecureCommunicationWithMdm": true,
		"allowNonSecureCommunicationWithLia": true,
		"disableNonMgmtComponentsAuth":       false,
	}
	mapData["securityConfiguration"] = secureData

	jsonres, _ := json.Marshal(mapData)

	res, err := GC.GetClusterDetails(jsonres, false)

	assert.NotNil(t, res)

	assert.EqualValues(t, res.StatusCode, 200)

	assert.Nil(t, err)
}

// TestDeployDeletePackge function to test Delete Functionality
func TestDeployDeletePackge(t *testing.T) {
	res, err := GC.DeletePackage("ABC")
	assert.EqualValues(t, res.StatusCode, 500)
	assert.NotNil(t, res)
	assert.Nil(t, err)
}

// TestDeployBeginInstallation function to test Begin Installation with Parsed CSV Data
func TestDeployBeginInstallation(t *testing.T) {
	_, err := GC.BeginInstallation("", "admin", "Password", "Password", true, true, false, true)
	assert.NotNil(t, err)
}

// TestDeployUninstallCluster function to test Begin Installation with Parsed CSV Data
func TestDeployUninstallCluster(t *testing.T) {
	jsonString := `{"snmpIp":null,"installationId":"4e7618620d6948f0","mdmIPs":["10.247.103.161","10.247.103.162"],"safeIPsForLia":null,"mdmUser":"admin","mdmPassword":null,"liaLdapUsername":null,"liaLdapPassword":null,"liaLdapInitialMode":null,"liaLdapInitialUri":null,"liaLdapInitialBaseDn":null,"liaLdapInitialGroup":null,"liaLdapInitialUsernameDnFormat":null,"liaPassword":null,"liaLdapInitialSearchFilterFormat":null,"licenseKey":null,"licenseType":null,"masterMdm":{"node":{"ostype":"unknown","nodeName":null,"nodeIPs":["10.247.103.161"],"domain":null,"userName":null,"password":null,"liaLdapUsername":null,"liaLdapPassword":null,"liaPassword":null},"nodeInfo":null,"tunables":null,"rollbackVersion":null,"mdmIPs":["10.247.103.161"],"name":null,"id":"8298098732753089792","ipForActor":null,"managementIPs":["10.247.103.161"],"virtIpIntfsList":null},"isClusterOptimized":null,"slaveMdmSet":[{"node":{"ostype":"unknown","nodeName":null,"nodeIPs":["10.247.103.162"],"domain":null,"userName":null,"password":null,"liaLdapUsername":null,"liaLdapPassword":null,"liaPassword":null},"nodeInfo":null,"tunables":null,"rollbackVersion":null,"mdmIPs":["10.247.103.162"],"name":null,"id":"5124358629956574209","ipForActor":null,"managementIPs":["10.247.103.162"],"virtIpIntfsList":null}],"tbSet":[{"node":{"ostype":"unknown","nodeName":null,"nodeIPs":["10.247.103.160"],"domain":null,"userName":null,"password":null,"liaLdapUsername":null,"liaLdapPassword":null,"liaPassword":null},"nodeInfo":null,"tunables":null,"rollbackVersion":null,"mdmIPs":["10.247.103.160"],"name":null,"id":"3749187462515833602","tbIPs":["10.247.103.160"]}],"standbyMdmSet":[],"standbyTbSet":[],"sdsList":[{"node":{"ostype":"unknown","nodeName":null,"nodeIPs":["10.247.103.162"],"domain":null,"userName":null,"password":null,"liaLdapUsername":null,"liaLdapPassword":null,"liaPassword":null},"nodeInfo":null,"tunables":null,"rollbackVersion":null,"sdsName":"sds2","protectionDomain":"default","protectionDomainId":"4373788827873968128","faultSet":"","faultSetId":"0","allIPs":["10.247.103.162"],"sdsOnlyIPs":[],"sdcOnlyIPs":[],"devices":[],"rfCached":false,"rfCachedPools":[],"rfCachedDevices":[],"rfCacheType":null,"flashAccDevices":[],"nvdimmAccDevices":[],"useRmCache":false,"optimized":false,"packageNumber":0,"optimizedNumOfIOBufs":3,"port":7072,"id":"-1832116130873868287"},{"node":{"ostype":"unknown","nodeName":null,"nodeIPs":["10.247.103.161"],"domain":null,"userName":null,"password":null,"liaLdapUsername":null,"liaLdapPassword":null,"liaPassword":null},"nodeInfo":null,"tunables":null,"rollbackVersion":null,"sdsName":"sds1","protectionDomain":"default","protectionDomainId":"4373788827873968128","faultSet":"","faultSetId":"0","allIPs":["10.247.103.161"],"sdsOnlyIPs":[],"sdcOnlyIPs":[],"devices":[],"rfCached":false,"rfCachedPools":[],"rfCachedDevices":[],"rfCacheType":null,"flashAccDevices":[],"nvdimmAccDevices":[],"useRmCache":false,"optimized":false,"packageNumber":0,"optimizedNumOfIOBufs":3,"port":7072,"id":"-1832116135168835584"},{"node":{"ostype":"unknown","nodeName":null,"nodeIPs":["10.247.103.160"],"domain":null,"userName":null,"password":null,"liaLdapUsername":null,"liaLdapPassword":null,"liaPassword":null},"nodeInfo":null,"tunables":null,"rollbackVersion":null,"sdsName":"sds3","protectionDomain":"default","protectionDomainId":"4373788827873968128","faultSet":"","faultSetId":"0","allIPs":["10.247.103.160"],"sdsOnlyIPs":[],"sdcOnlyIPs":[],"devices":[],"rfCached":false,"rfCachedPools":[],"rfCachedDevices":[],"rfCacheType":null,"flashAccDevices":[],"nvdimmAccDevices":[],"useRmCache":false,"optimized":false,"packageNumber":0,"optimizedNumOfIOBufs":3,"port":7072,"id":"-1832116126578900990"}],"sdcList":[{"node":{"ostype":"unknown","nodeName":null,"nodeIPs":["10.247.103.161"],"domain":null,"userName":null,"password":null,"liaLdapUsername":null,"liaLdapPassword":null,"liaPassword":null},"nodeInfo":null,"tunables":null,"rollbackVersion":null,"guid":"44E0832B-4D83-4557-954A-C8C3FDB1E325","splitterRpaIp":null,"sdcName":null,"isOnESX":"NO","optimized":false,"id":"4784734193563205634"},{"node":{"ostype":"unknown","nodeName":null,"nodeIPs":["10.247.103.162"],"domain":null,"userName":null,"password":null,"liaLdapUsername":null,"liaLdapPassword":null,"liaPassword":null},"nodeInfo":null,"tunables":null,"rollbackVersion":null,"guid":"54C2C77B-953D-415D-B088-30AC11610565","splitterRpaIp":null,"sdcName":null,"isOnESX":"NO","optimized":false,"id":"4784734189268238337"},{"node":{"ostype":"unknown","nodeName":null,"nodeIPs":["10.247.103.160"],"domain":null,"userName":null,"password":null,"liaLdapUsername":null,"liaLdapPassword":null,"liaPassword":null},"nodeInfo":null,"tunables":null,"rollbackVersion":null,"guid":"AA378BC3-AA19-4C50-8746-824B012D8C23","splitterRpaIp":null,"sdcName":null,"isOnESX":"NO","optimized":false,"id":"4784734184973271040"}],"vasaProviderList":[],"numOfVolumes":0,"sdrList":[],"callHomeConfiguration":null,"remoteSyslogConfiguration":null,"systemVersionName":"DellEMC PowerFlex Version: R3_6.700.103","securityConfiguration":null,"virtualIps":null,"securityCommunicationEnabled":true,"skipCollectLogsOnEsx":false,"protectionDomains":[{"name":"default","storagePools":[{"name":"default","mediaType":"MEDIA_TYPE_SSD","externalAccelerationType":"EXTERNAL_ACCELERATION_TYPE_READ_AND_WRITE","dataLayout":"PERFORMANCE_OPTIMIZED","compressionMethod":"COMPRESSION_METHOD_INVALID","spefAccPoolName":null,"shouldApplyZeroPadding":false,"writeAtomicitySize":null,"overProvisioningFactor":null,"maxCompressionRatio":null,"perfProfile":null,"rplJournalCapacity":null}],"accelerationPools":[]}],"upgradeableByCurrentUser":true,"ignoredSdcs":[],"clustered":true,"sdsAndMdmIps":["10.247.103.161","10.247.103.162","10.247.103.160"],"sdcIps":["10.247.103.161","10.247.103.162","10.247.103.160"],"upgradeRunning":false,"pd2Sps":{"default":["default"]},"primaryMdm":null,"secondaryMdm":null,"tb":null}`
	_, err := GC.UninstallCluster(jsonString, "admin", "Password", "Password", true, true, false, true)
	assert.NotNil(t, err)
}

// TestDeployMoveToNextPhase function to test Move to Next Phase Functionality
func TestDeployMoveToNextPhase(t *testing.T) {
	res, err := GC.MoveToNextPhase()
	assert.NotNil(t, res)
	assert.EqualValues(t, res.StatusCode, 500)
	assert.Nil(t, err)
}

// TestDeployRetryPhase function to test Retry Failed Phase  Functionality
func TestDeployRetryPhase(t *testing.T) {
	res, err := GC.RetryPhase()
	assert.NotNil(t, res)
	assert.EqualValues(t, res.StatusCode, 200)
	assert.Nil(t, err)
}

// TestDeployAbortOperation function to test Abort Operation Functionality
func TestDeployAbortOperation(t *testing.T) {
	res, err := GC.AbortOperation()
	assert.NotNil(t, res)
	assert.EqualValues(t, res.StatusCode, 200)
	assert.Nil(t, err)
}

// TestDeployClearQueueCommand function to test Clear Queue Commands Functionality
func TestDeployClearQueueCommand(t *testing.T) {
	res, err := GC.ClearQueueCommand()
	assert.NotNil(t, res)
	assert.Nil(t, err)
}

// TestDeployMoveToIdlePhase function to test Move to Ideal Phase Functionality
func TestDeployMoveToIdlePhase(t *testing.T) {
	res, err := GC.MoveToIdlePhase()
	assert.NotNil(t, res)
	assert.EqualValues(t, res.StatusCode, 200)
	assert.Nil(t, err)
}

// TestDeployGetInQueueCommand function to test Queue Command Detail API
func TestDeployGetInQueueCommand(t *testing.T) {
	_, err := GC.GetInQueueCommand()
	assert.Nil(t, err)
}

// TestDeployCheckForCompletionQueueCommands function to test Queue Command Completed or not
func TestDeployCheckForCompletionQueueCommands(t *testing.T) {
	res, err := GC.CheckForCompletionQueueCommands("query")
	assert.EqualValues(t, res.StatusCode, 200)
	assert.NotNil(t, res)
	assert.Nil(t, err)
}
