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
	indepthDetails, err := GC.GetUploadComplianceDetails(details.ID, false)
	assert.Nil(t, err)
	assert.NotEmpty(t, indepthDetails.ID)
	assert.NotEmpty(t, indepthDetails.State)

	details2, err2 := GC.GetUploadComplianceDetailsUsingID(details.ID)
	assert.Nil(t, err2)
	assert.NotNil(t, details2.ID)
	assert.NotNil(t, details2.State)
	details3, err3 := GC.GetFirmwareRepositoryDetailsUsingName(details.Name)
	assert.Nil(t, err3)
	assert.NotNil(t, details3.ID)
	assert.NotNil(t, details3.State)
}

func TestApproveUnsignedFile(t *testing.T) {
	var unsigned string
	if os.Getenv("GOSCALEIO_UNSIGNED_COMPLIANCE_FILE_ID") != "" {
		unsigned = os.Getenv("GOSCALEIO_UNSIGNED_COMPLIANCE_FILE_ID")
	}
	err := GC.ApproveUnsignedFile(unsigned)
	assert.Nil(t, err)
}
