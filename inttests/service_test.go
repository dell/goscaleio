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
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDeployeService(t *testing.T) {
	_, err := GC.DeployService("Test-Create", "Test", "453c41eb-d72a-4ed1-ad16-bacdffbdd766", "8aaaee208c8c467e018cd37813250614", "3")
	assert.NotNil(t, err)
}

func TestGetAllDeployeService(t *testing.T) {
	deployments, err := GC.GetAllServiceDetails()
	assert.Nil(t, err)
	assert.NotNil(t, deployments)

	if len(deployments) > 0 {
		template, err := GC.GetServiceDetailsByID(deployments[0].ID, false)
		assert.Nil(t, err)
		assert.NotNil(t, template)

		_, err = GC.UpdateService(template.ID, "Test-Update-K", "Test-Update-K", "5", "pfmc-k8s-20230809-160")
		assert.NotNil(t, err)
	}

	template, err := GC.GetServiceDetailsByID(invalidIdentifier, false)
	assert.NotNil(t, err)
	assert.Nil(t, template)
}

func TestGetDeployeServiceByName(t *testing.T) {
	deployments, err := GC.GetAllServiceDetails()
	assert.Nil(t, err)
	assert.NotNil(t, deployments)

	if len(deployments) > 0 {
		template, err := GC.GetServiceDetailsByFilter("name", deployments[0].DeploymentName)
		assert.Nil(t, err)
		assert.NotNil(t, template)
	}

	template, err := GC.GetServiceDetailsByFilter("invalid", "invalid")
	assert.NotNil(t, err)
	assert.Nil(t, template)
}

func TestGetAllComplianceReports(t *testing.T) {
	complianceReport, err := GC.GetServiceComplianceDetails("8aaa806c9035ea9101903ab07ff623b2")
	assert.Nil(t, err)
	assert.NotNil(t, complianceReport)
}

func TestGetAllComplianceReportsByFilter(t *testing.T) {
	complianceReport, err := GC.GetServiceComplianceDetailsByFilter("8aaa806c9035ea9101903ab07ff623b2", "ServiceTag", "VMware-42 0c 5e 1c 03 b5 bf a1-a7 cc 69 60 f9 52 80 fa-SW")
	assert.Nil(t, err)
	assert.NotNil(t, complianceReport)
	// invalid filter
	complianceReport, err = GC.GetServiceComplianceDetailsByFilter("8aaa806c9035ea9101903ab07ff623b2", "NotValid", "VMware-42 0c 5e 1c 03 b5 bf a1-a7 cc 69 60 f9 52 80 fa-SW")
	assert.Error(t, err)
}
