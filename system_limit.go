// Copyright © 2019 - 2023 Dell Inc. or its subsidiaries. All Rights Reserved.
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
	"net/http"
	"time"

	types "github.com/dell/goscaleio/types/v1"
)

type QuerySystemLimitsParam struct {
}

// GetSystemLimits gets list of sytem limits

func (c *Client) GetSystemLimits() (systemLimits *types.Limit, err error) {
	defer TimeSpent("GetSystemLimits", time.Now())
	var body QuerySystemLimitsParam
	path := "/api/instances/System/action/querySystemLimits"
	err = c.getJSONWithRetry(
		http.MethodPost, path, body, &systemLimits)
	if err != nil {
		return nil, err
	}

	return systemLimits, nil
}

// GetMaxVol returns max volume size in GB
func (c *Client) GetMaxVol() (sys string, err error) {
	defer TimeSpent("GetMaxVol", time.Now())
	maxlimitType, err := c.GetSystemLimits()

	if err != nil {
		return "", err
	}

	for _, systype := range maxlimitType.SystemLimitEntryList {

		if systype.Type == "volumeSizeGb" {
			return systype.MaxVal, nil
		}

	}
	return "", errors.New("couldn't get max vol size")
}
