/*
 *
 * Copyright © 2021-2024 Dell Inc. or its subsidiaries. All Rights Reserved.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */

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

	types "github.com/dell/goscaleio/types/v1"
	"github.com/stretchr/testify/assert"
)

func TestCreateAndDeleteSSOUser(t *testing.T) {
	details, err := C.CreateSSOUser(&types.SSOUserCreateParam{
		UserName: "IntegrationTestSSOUser",
		Role:     "Monitor",
		Password: "Ssouser123!",
		Type:     "Local",
	})
	assert.Nil(t, err)
	assert.NotEmpty(t, details)

	details, err = C.GetSSOUser(details.ID)
	assert.Nil(t, err)
	assert.NotEmpty(t, details)

	details1, err := C.GetSSOUserByFilters("username", "IntegrationTestSSOUser")
	assert.Nil(t, err)
	assert.NotEmpty(t, details1)

	details, err = C.ModifySSOUser(details.ID, &types.SSOUserModifyParam{
		Role: "Technician",
	})
	assert.Nil(t, err)
	assert.NotEmpty(t, details)

	err = C.ResetSSOUserPassword(details.ID, &types.SSOUserModifyParam{Password: "Ssouser1234#"})
	assert.Nil(t, err)

	err = C.DeleteSSOUser(details.ID)
	assert.Nil(t, err)
}
