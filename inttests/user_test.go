/*
 *
 * Copyright © 2020 Dell Inc. or its subsidiaries. All Rights Reserved.
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

package inttests

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestGetAllUsers will return all user instances
func TestGetAllUsers(t *testing.T) {
	system := getSystem()
	assert.NotNil(t, system)

	allUsers, err := system.GetUser()
	assert.Nil(t, err)
	assert.NotZero(t, len(allUsers))
}
