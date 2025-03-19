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

func TestGetTemplates(t *testing.T) {
	templates, err := GC.GetAllTemplates()
	assert.Nil(t, err)
	assert.NotNil(t, templates)
}

func TestGetTemplateByID(t *testing.T) {
	templates, err := GC.GetAllTemplates()
	assert.Nil(t, err)
	assert.NotNil(t, templates)

	if len(templates) > 0 {
		template, err := GC.GetTemplateByID(templates[0].ID)
		assert.Nil(t, err)
		assert.NotNil(t, template)
	}

	template, err := GC.GetTemplateByID(invalidIdentifier)
	assert.NotNil(t, err)
	assert.Nil(t, template)
}

func TestGetTemplateByFilters(t *testing.T) {
	templates, err := GC.GetTemplateByFilters("name", "invalid")
	assert.NotNil(t, err)
	assert.Nil(t, templates)
}

func TestCloneTemplate(t *testing.T) {
	templates, err := GC.GetTemplateByFilters("category", "Sample Templates")
	assert.Nil(t, err)
	assert.NotNil(t, templates)

	system := getSystem()
	assert.NotNil(t, system)

	errCT := GC.CloneTemplate(system, templates[0].TemplateName+"- Copy", templates[0].OriginalTemplateID)
	assert.Nil(t, errCT)
}
