// Copyright © 2021 - 2022 Dell Inc. or its subsidiaries. All Rights Reserved.
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
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"testing"

	siotypes "github.com/AnshumanPradipPatil1506/goscaleio/types/v1/template"
	"github.com/stretchr/testify/assert"
)

var (
	empty = ""
	tab   = "\t"
	// TemplateName   string = ""
	// TemplateID     string = ""
	// TemplateString string = ""
	// TemplateGet    siotypes.DefaultTemplate
)

func PrettyJson(data interface{}) (string, error) {
	buffer := new(bytes.Buffer)
	encoder := json.NewEncoder(buffer)
	encoder.SetIndent(empty, tab)

	err := encoder.Encode(data)
	if err != nil {
		return empty, err
	}
	return buffer.String(), nil
}

func createTemplateFunction(tplstr string) (*siotypes.DefaultTemplate, error) {
	resultTemplate := siotypes.DefaultTemplate{}
	templateString := siotypes.TemplateData
	mytemplate := siotypes.DefaultTemplate{}
	json.Unmarshal([]byte(templateString), &mytemplate)
	response, err := C.CreateTemplate().FromString(tplstr)

	if err != nil {
		return &resultTemplate, err
	}

	b, err := io.ReadAll(response.(*http.Response).Body)
	if err != nil {
		return &resultTemplate, err
	}
	json.Unmarshal(b, &resultTemplate)
	fmt.Println("\nTemplate Created- ID ", resultTemplate.ID)
	fmt.Println("\nTemplate Created- TemplateName ", resultTemplate.TemplateName)
	return &resultTemplate, err
}

func getTemplateFunction(nm string) (*siotypes.DefaultTemplate, error) {
	resultTemplate := siotypes.DefaultTemplate{}

	if nm == "" {
		return &resultTemplate, fmt.Errorf("No template created named %v", nm)
	}
	responseGet, err := C.GetTemplate(nm)
	if err != nil {
		return &resultTemplate, err
	}

	resultTemplatesGet := []siotypes.DefaultTemplate{}
	b, err := io.ReadAll(responseGet.(*http.Response).Body)
	if err != nil {
		return &resultTemplate, err
	}
	json.Unmarshal(b, &resultTemplatesGet)
	resultTemplate = resultTemplatesGet[0]
	fmt.Println("\nGET Template -ID ", resultTemplate.ID)
	fmt.Println("\nGET Template -Name ", resultTemplate.TemplateName)
	return &resultTemplate, nil
}

func deleteTemplateFunction(tplid string) (int, error) {
	statuscode := 0
	if tplid == "" {
		return statuscode, fmt.Errorf("No template created of id %v", tplid)
	}
	responseDelete, err := C.DeleteTemplate(tplid)

	if err != nil {
		return statuscode, err
	}
	statuscode = responseDelete.(*http.Response).StatusCode
	fmt.Println("\nDeleted Template - response code ", statuscode)

	return statuscode, err
}

func updateTemplateFunction(template siotypes.DefaultTemplate) (int, error) {
	update, err := PrettyJson(template)

	statuscode := 0

	response, err := C.UpdateTemplate(update, template.ID)

	if err != nil {
		return statuscode, err
	}
	statuscode = response.(*http.Response).StatusCode

	fmt.Println("\nUpdated Template - Description ", template.TemplateDescription)
	fmt.Println("\nUpdated Template - response code ", statuscode)
	return statuscode, err
}

func TestTemplateFeature(t *testing.T) {
	TemplateName := randString(10)

	tplstr := `{
		"templateName": "` + TemplateName + `",
		"templateDescription": "` + TemplateName + ` template",
		"templateType": "Software Only",
		"draft": true
		}`
	Templt, err := createTemplateFunction(tplstr)
	assert.Nil(t, err)
	assert.NotNil(t, Templt)
	assert.NotNil(t, Templt.ID)

	getTemplt, err := getTemplateFunction(Templt.TemplateName)
	assert.Nil(t, err)
	assert.NotNil(t, getTemplt)

	getTemplt.TemplateDescription = getTemplt.TemplateDescription + " updated"
	updatestatuscode, err := updateTemplateFunction(*getTemplt)
	assert.Nil(t, err)
	assert.NotNil(t, updatestatuscode)
	assert.EqualValues(t, 204, updatestatuscode)

	deletestatuscode, err := deleteTemplateFunction(Templt.ID)
	assert.Nil(t, err)
	assert.NotNil(t, deletestatuscode)
	assert.EqualValues(t, 204, deletestatuscode)
}
