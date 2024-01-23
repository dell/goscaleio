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

package goscaleio

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	types "github.com/dell/goscaleio/types/v1"
)

// GetTemplateByID gets the node details based on ID
func (c *Client) GetTemplateByID(id string) (*types.TemplateDetails, error) {
	defer TimeSpent("GetTemplateByID", time.Now())

	path := fmt.Sprintf("/Api/V1/template/%v", id)

	var template types.TemplateDetails
	err := c.getJSONWithRetry(http.MethodGet, path, nil, &template)
	if err != nil {
		return nil, err
	}
	return &template, nil
}

// GetAllTemplates gets all the Template details
func (c *Client) GetAllTemplates() ([]types.TemplateDetails, error) {
	defer TimeSpent("GetAllTemplates", time.Now())

	path := fmt.Sprintf("/Api/V1/template")

	var templates types.TemplateDetailsFilter
	err := c.getJSONWithRetry(http.MethodGet, path, nil, &templates)
	if err != nil {
		return nil, err
	}
	return templates.TemplateDetails, nil
}

// GetTemplateByFilters gets the Template details based on the provided filter
func (c *Client) GetTemplateByFilters(key string, value string) ([]types.TemplateDetails, error) {
	defer TimeSpent("GetTemplateByFilters", time.Now())

	path := fmt.Sprintf("/Api/V1/template?filter=eq,%v,%v", key, value)

	var templates types.TemplateDetailsFilter
	err := c.getJSONWithRetry(http.MethodGet, path, nil, &templates)
	if err != nil {
		return nil, err
	}

	if len(templates.TemplateDetails) == 0 {
		return nil, errors.New("Couldn't find templates with the given filter")
	}
	return templates.TemplateDetails, nil
}
