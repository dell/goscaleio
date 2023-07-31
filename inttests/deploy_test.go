package inttests

import (
	"encoding/json"
	"testing"
	"os"

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
	clusterData := map[string]interface{}{
		"mdmUser":     "admin",
		"mdmPassword": string(os.Getenv("GOSCALEIO_MDMPASSWORD")),
	}
	clusterData["mdmIps"] = []string{string(os.Getenv("GOSCALEIO_MDMIP"))}

	secureData := map[string]interface{}{
		"allowNonSecureCommunicationWithMdm": true,
		"allowNonSecureCommunicationWithLia": true,
		"disableNonMgmtComponentsAuth":       false,
	}
	clusterData["securityConfiguration"] = secureData

	jsonres, _ := json.Marshal(clusterData)

	res, err := GC.ValidateMDMDetails(jsonres)

	assert.NotNil(t, res)

	assert.EqualValues(t, res.StatusCode, 200)

	assert.Nil(t, err)
}

func TestDeployGetClusterDetails(t *testing.T) {
	clusterData := map[string]interface{}{
		"mdmUser":     "admin",
		"mdmPassword": string(os.Getenv("GOSCALEIO_MDMPASSWORD")),
	}
	clusterData["mdmIps"] = []string{string(os.Getenv("GOSCALEIO_MDMIP"))}

	secureData := map[string]interface{}{
		"allowNonSecureCommunicationWithMdm": true,
		"allowNonSecureCommunicationWithLia": true,
		"disableNonMgmtComponentsAuth":       false,
	}
	clusterData["securityConfiguration"] = secureData

	jsonres, _ := json.Marshal(clusterData)

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
	_, err := GC.UninstallCluster("", "admin", "Password", "Password", true, true, false, true)
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
