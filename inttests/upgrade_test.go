package inttests

import (
	"os"
	"testing"
	"time"

	types "github.com/dell/goscaleio/types/v1"
	"github.com/stretchr/testify/assert"
)

// TestDeployUploadPackage function to test upload packge with dummy path of packages
func TestUploadCompliance(t *testing.T) {
	var sourceLocation string
	if os.Getenv("GOSCALEIO_COMPLIANCE_ENDPOINT") != "" {
		sourceLocation = os.Getenv("GOSCALEIO_COMPLIANCE_ENDPOINT")
	}
	ucParam := &types.UploadComplianceParam{
		SourceLocation: sourceLocation,
	}
	details, err := GC.UploadCompliance(ucParam)
	assert.Nil(t, err)
	assert.NotNil(t, details.ID)
	assert.NotNil(t, details.State)
	time.Sleep(5 * time.Second)
	indepthDetails, err := GC.GetUploadComplianceDetails(details.ID)
	assert.Nil(t, err)
	assert.NotEmpty(t, indepthDetails.ID)
	assert.NotEmpty(t, indepthDetails.State)
}

func TestApproveUnsignedFile(t *testing.T) {
	var unsigned string
	if os.Getenv("GOSCALEIO_UNSIGNED_COMPLIANCE_FILE_ID") != "" {
		unsigned = os.Getenv("GOSCALEIO_UNSIGNED_COMPLIANCE_FILE_ID")
	}
	err := GC.ApproveUnsignedFile(unsigned)
	assert.Nil(t, err)
}
