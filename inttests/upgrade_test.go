package inttests

import (
	"testing"
	"time"

	types "github.com/dell/goscaleio/types/v1"
	"github.com/stretchr/testify/assert"
)

// TestDeployUploadPackage function to test upload packge with dummy path of packages
func TestUploadCompliance(t *testing.T) {
	ucParam := &types.UploadComplianceParam{
		SourceLocation:   "https://100.65.27.72/artifactory/vxfm-yum-release/pfmp20/RCM/Denver/RCMs/SoftwareOnly/PowerFlex_Software_4.5.0.0_287_r1.zip"            ,
	}
	details, err := GC.UploadCompliance(ucParam)
	assert.Nil(t, err)
	assert.NotNil(t, details.ID)
	assert.NotNil(t, details.State)
	time.Sleep(5*time.Second)
	indepthDetails,err :=GC.GetUploadComplianceDetails(details.ID)
	assert.Nil(t, err)
	assert.NotEmpty(t, indepthDetails.ID)
	assert.NotEmpty(t, indepthDetails.State)
}

func TestApproveUnsignedFile(t *testing.T) {
	err :=GC.ApproveUnsignedFile("8aaa3fd38f4c78eb018f4dad3781001d")
	assert.Nil(t, err)
}