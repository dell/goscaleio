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
	templates, err := GC.DeployService("Test-Create", "Test", "8150d563-639d-464e-80c4-a435ed10f132", "8aaaee208c8c467e018cd37813250614")
	assert.Nil(t, err)
	assert.NotNil(t, templates)
}

func TestGetAllDeployeService(t *testing.T) {
	templates, err := GC.GetAllServiceDetails()
	assert.Nil(t, err)
	assert.NotNil(t, templates)

	if len(templates) > 0 {
		template, err := GC.GetServiceDetailsByID(templates[0].ID)
		assert.Nil(t, err)
		assert.NotNil(t, template)
	}

	template, err := GC.GetServiceDetailsByID(invalidIdentifier)
	assert.NotNil(t, err)
	assert.Nil(t, template)
}

func TestGetDeployeServiceByName(t *testing.T) {
	templates, err := GC.GetAllServiceDetails()
	assert.Nil(t, err)
	assert.NotNil(t, templates)

	if len(templates) > 0 {
		template, err := GC.GetServiceDetailsByFilter(templates[0].DeploymentName, "name")
		assert.Nil(t, err)
		assert.NotNil(t, template)
	}

	template, err := GC.GetServiceDetailsByFilter("invalid", "invalid")
	assert.NotNil(t, err)
	assert.Nil(t, template)
}
