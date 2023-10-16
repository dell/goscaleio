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

type emptyBody struct {
}

func (c *Client) GetSystemLimits() (syslimit *types.Limit, err error) {
	defer TimeSpent("GetSystemLimits", time.Now())
	var e emptyBody
	path := "/api/instances/System/action/querySystemLimits"
	err = c.getJSONWithRetry(
		http.MethodPost, path, e, &syslimit)
	if err != nil {
		return nil, err
	}

	return syslimit, nil
}

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
