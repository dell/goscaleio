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
	siotypes "github.com/dell/goscaleio/types/v1"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

// TestGetAllUsers will return all user instances
func TestGetAllUsers(t *testing.T) {
	system := getSystem()
	assert.NotNil(t, system)

	allUsers, err := system.GetUser()
	assert.Nil(t, err)
	assert.NotZero(t, len(allUsers))
}

func TestCreateAndDeleteUser(t *testing.T) {
	// Get the System
	system := getSystem()
	assert.NotNil(t, system)

	// Create a new User
	userParams := siotypes.UserParam{
		Name:     "testUser",
		Password: os.Getenv("USER_PASSWORD"),
		UserRole: "Security",
	}
	resp, err1 := system.CreateUser(&userParams)
	assert.Nil(t, err1)
	assert.NotEmpty(t, resp)

	// Fetch the User which you just now created
	user, err2 := system.GetUserByID(resp.ID)
	assert.Nil(t, err2)
	assert.Equal(t, "testUser", user.Name)
	assert.Equal(t, "Security", user.UserRole)

	// Change the user role
	userRoleParams := siotypes.UserRoleParam{
		UserRole: "Configure",
	}
	err3 := system.SetUserRole(&userRoleParams, resp.ID)
	assert.Nil(t, err3)

	// Remove the user
	err4 := system.RemoveUser(resp.ID)
	assert.Nil(t, err4)

}
