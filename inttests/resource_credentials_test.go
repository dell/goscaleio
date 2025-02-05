/*
 *
 * Copyright © 2025 Dell Inc. or its subsidiaries. All Rights Reserved.
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

	siotypes "github.com/dell/goscaleio/types/v1"
	"github.com/stretchr/testify/assert"
)

var exampleSSHKey = `-----BEGIN OPENSSH PRIVATE KEY-----
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
-----END OPENSSH PRIVATE KEY-----`

// TestGetResourceCredentials Get All Resource Credentials
func TestGetResourceCredentials(t *testing.T) {
	system := getSystem()
	creds, err := system.GetResourceCredentials()
	assert.Nil(t, err)
	assert.NotNil(t, creds)
}

// TestGetResourceCredential Gets a specific Resource Credential
func TestGetResourceCredential(t *testing.T) {
	system := getSystem()
	creds, err := system.GetResourceCredentials()
	assert.Nil(t, err)
	assert.NotNil(t, creds)
	id := creds.Credentials[0].Credential.ID
	cred, err1 := system.GetResourceCredential(id)
	assert.Nil(t, err1)
	assert.NotNil(t, cred)
}

// TestCreateModifyResourceCredentialForNode Creates a new Resource Credential of node type
func TestCreateModifyResourceCredentialForNode(t *testing.T) {
	system := getSystem()
	// Create
	cred, err1 := system.CreateNodeResourceCredential(siotypes.ServerCredential{
		Username: "testFT3",
		Password: "testFT4",
		Label:    "testFT3",
	})
	assert.Nil(t, err1)
	assert.NotNil(t, cred)
	id := cred.Credential.ID
	// Modify
	credMod, err2 := system.ModifyNodeResourceCredential(siotypes.ServerCredential{
		Username:            "testFT6",
		Password:            "testFT1",
		Label:               "testFT7",
		SNMPv3SecurityName:  "ConTest",
		SNMPv3SecurityLevel: "1",
	}, id)
	assert.Nil(t, err2)
	assert.NotNil(t, credMod)
	// Modify SNMPv3 Level 2 should error bc no MD5
	_, err3 := system.ModifyNodeResourceCredential(siotypes.ServerCredential{
		Username:            "testFT6",
		Password:            "testFT1",
		Label:               "testFT7",
		SNMPv3SecurityName:  "ConTest",
		SNMPv3SecurityLevel: "2",
	}, id)
	assert.NotNil(t, err3)
	// Modify SNMPv3 Level 3 should error bc no DES
	_, err4 := system.ModifyNodeResourceCredential(siotypes.ServerCredential{
		Username:                        "testFT6",
		Password:                        "testFT1",
		Label:                           "testFT7",
		SNMPv3SecurityName:              "ConTestLevel3",
		SNMPv3MD5AuthenticationPassword: "ConMD5",
		SNMPv3SecurityLevel:             "3",
	}, id)
	assert.NotNil(t, err4)
	// Modify SNMPv3 Level 3 should Pass with all Values Set
	type3Cred, err5 := system.ModifyNodeResourceCredential(siotypes.ServerCredential{
		Username:                        "testFT6",
		Password:                        "testFT1",
		Label:                           "testFT7",
		SNMPv3SecurityName:              "ConTestLevel3",
		SNMPv3MD5AuthenticationPassword: "ConTestMD5",
		SNMPv3DesPrivatePassword:        "ConTestDES",
		SNMPv3SecurityLevel:             "3",
	}, id)
	assert.Nil(t, err5)
	assert.NotNil(t, type3Cred)

	// Modify Set PrivateKey without Name should fail
	_, err6 := system.ModifyNodeResourceCredential(siotypes.ServerCredential{
		Username:                        "testFT6",
		Password:                        "testFT1",
		Label:                           "testFT7",
		SNMPv3SecurityName:              "ConTestLevel3",
		SNMPv3MD5AuthenticationPassword: "ConTestMD5",
		SNMPv3DesPrivatePassword:        "ConTestDES",
		SNMPv3SecurityLevel:             "3",
		SSHPrivateKey:                   exampleSSHKey,
	}, id)
	assert.NotNil(t, err6)

	// Modify Set PrivateKey without Name should fail
	privateKeySet, err7 := system.ModifyNodeResourceCredential(siotypes.ServerCredential{
		Username:                        "testFT6",
		Password:                        "testFT1",
		Label:                           "testFT7",
		SNMPv3SecurityName:              "ConTestLevel3",
		SNMPv3MD5AuthenticationPassword: "ConTestMD5",
		SNMPv3DesPrivatePassword:        "ConTestDES",
		SNMPv3SecurityLevel:             "3",
		SSHPrivateKey:                   exampleSSHKey,
		KeyPairName:                     "ConTestKeyPair",
	}, id)
	assert.Nil(t, err7)
	assert.NotNil(t, privateKeySet)
	// Delete
	errDel := system.DeleteResourceCredential(id)
	assert.Nil(t, errDel)
}

// TestCreateModifyResourceCredentialForSwitch creates, modifies and deletes new Resource Credential of switch type
func TestCreateModifyResourceCredentialForSwitch(t *testing.T) {
	system := getSystem()
	// Create
	cred, err1 := system.CreateSwitchResourceCredential(siotypes.IomCredential{
		Username:              "testFT3",
		Password:              "testFT4",
		Label:                 "testFT3",
		SNMPv2CommunityString: "public",
	})
	assert.Nil(t, err1)
	assert.NotNil(t, cred)
	id := cred.Credential.ID
	// Modify
	credMod, err2 := system.ModifySwitchResourceCredential(siotypes.IomCredential{
		Username: "testFT1",
		Password: "testFT2",
		Label:    "testFT3",
	}, id)
	assert.Nil(t, err2)
	assert.NotNil(t, credMod)

	// Modify Partial Private Key Should Fail
	_, err3 := system.ModifySwitchResourceCredential(siotypes.IomCredential{
		Username:      "testFT1",
		Password:      "testFT2",
		Label:         "testFT3",
		SSHPrivateKey: exampleSSHKey,
	}, id)
	assert.NotNil(t, err3)

	// Modify Add Private Key Should succeed
	privateKey, err4 := system.ModifySwitchResourceCredential(siotypes.IomCredential{
		Username:      "testFT1",
		Password:      "testFT2",
		Label:         "testFT3",
		SSHPrivateKey: exampleSSHKey,
		KeyPairName:   "ConTestKeyPair",
	}, id)
	assert.Nil(t, err4)
	assert.NotNil(t, privateKey)
	// Delete
	errDel := system.DeleteResourceCredential(id)
	assert.Nil(t, errDel)
}

// TestCreateModifyResourceCredentialForvCenter create, modify and delete vCenter Credential
func TestCreateModifyResourceCredentialForvCenter(t *testing.T) {
	system := getSystem()
	// Create
	cred, err := system.CreateVCenterResourceCredential(siotypes.VCenterCredential{
		Username: "testFT3",
		Password: "testFT4",
		Label:    "testFT3",
		Domain:   "1.1.1.1",
	})
	assert.Nil(t, err)
	assert.NotNil(t, cred)
	id := cred.Credential.ID
	// Modify
	credMod, err2 := system.ModifyVCenterResourceCredential(siotypes.VCenterCredential{
		Username: "testFT4",
		Password: "testFT5",
		Label:    "testFT1",
		Domain:   "2.2.2.2",
	}, id)
	assert.Nil(t, err2)
	assert.NotNil(t, credMod)
	// Delete
	errDel := system.DeleteResourceCredential(id)
	assert.Nil(t, errDel)
}

// TestCreateModifyResourceCredentialForElementManager create, modify and delete ElementManager Credential
func TestCreateModifyResourceCredentialForElementManager(t *testing.T) {
	system := getSystem()
	// Create
	cred, err := system.CreateElementManagerResourceCredential(siotypes.EMCredential{
		Username: "testFT3",
		Password: "testFT4",
		Label:    "testFT3",
		Domain:   "1.1.1.1",
	})
	assert.Nil(t, err)
	assert.NotNil(t, cred)
	id := cred.Credential.ID
	// Modify
	credMod, err2 := system.ModifyElementManagerResourceCredential(siotypes.EMCredential{
		Username: "testFT4",
		Password: "testFT5",
		Label:    "testFT1",
		Domain:   "2.2.2.2",
	}, id)
	assert.Nil(t, err2)
	assert.NotNil(t, credMod)
	// Delete
	errDel := system.DeleteResourceCredential(id)
	assert.Nil(t, errDel)
}

// TestCreateModifyResourceCredentialForScaleIO create, modify and delete ScaleIO Credential
func TestCreateModifyResourceCredentialForScaleIO(t *testing.T) {
	system := getSystem()
	// Create
	cred, err := system.CreateScaleIOResourceCredential(siotypes.ScaleIOCredential{
		AdminUsername: "testFT3",
		AdminPassword: "testFT4",
		Label:         "testFT3",
		OSUsername:    "testOSFT3",
		OSPassword:    "testOSFT4",
	})
	assert.Nil(t, err)
	assert.NotNil(t, cred)
	id := cred.Credential.ID
	// Modify
	credMod, err2 := system.ModifyScaleIOResourceCredential(siotypes.ScaleIOCredential{
		AdminUsername: "testFT4",
		AdminPassword: "testFT5",
		Label:         "testFT6",
		OSUsername:    "testOSFT1",
		OSPassword:    "testOSFT2",
	}, id)
	assert.Nil(t, err2)
	assert.NotNil(t, credMod)
	// Delete
	errDel := system.DeleteResourceCredential(id)
	assert.Nil(t, errDel)
}

// TestCreateModifyResourceCredentialForPresentationServer create, modify and delete PresentationServer Credential
func TestCreateModifyResourceCredentialForPresentationServer(t *testing.T) {
	system := getSystem()
	// Create
	cred, err := system.CreatePresentationServerResourceCredential(siotypes.PSCredential{
		Label:    "testFT3",
		Username: "testuser1",
		Password: "testFT4",
	})
	assert.Nil(t, err)
	assert.NotNil(t, cred)
	id := cred.Credential.ID
	// Modify
	credMod, err2 := system.ModifyPresentationServerResourceCredential(siotypes.PSCredential{
		Label:    "testFT5",
		Username: "testuser2",
		Password: "testFT3",
	}, id)
	assert.Nil(t, err2)
	assert.NotNil(t, credMod)
	// Delete
	errDel := system.DeleteResourceCredential(id)
	assert.Nil(t, errDel)
}

// TestCreateModifyResourceCredentialForOsAdmin create, modify and delete OsAdmin Credential
func TestCreateModifyResourceCredentialForOsAdmin(t *testing.T) {
	system := getSystem()
	// Create
	cred, err := system.CreateOsAdminResourceCredential(siotypes.OSAdminCredential{
		Password: "testFT4",
		Label:    "testFT3",
	})
	assert.Nil(t, err)
	assert.NotNil(t, cred)
	id := cred.Credential.ID
	// Modify
	credMod, err2 := system.ModifyOsAdminResourceCredential(siotypes.OSAdminCredential{
		Password:      "testFT4",
		Label:         "testFT3",
		KeyPairName:   "ConTestKeyPair",
		SSHPrivateKey: exampleSSHKey,
	}, id)
	assert.Nil(t, err2)
	assert.NotNil(t, credMod)
	// Delete
	errDel := system.DeleteResourceCredential(id)
	assert.Nil(t, errDel)
}

// TestCreateModifyResourceCredentialForOsUser create, modify and delete OsUser Credential
func TestCreateModifyResourceCredentialForOsUser(t *testing.T) {
	system := getSystem()
	// Create
	cred, err := system.CreateOsUserResourceCredential(siotypes.OSUserCredential{
		Password: "testFT4",
		Label:    "testFT3",
		Username: "test",
	})
	assert.Nil(t, err)
	assert.NotNil(t, cred)
	id := cred.Credential.ID
	// Modify
	credMod, err2 := system.ModifyOsUserResourceCredential(siotypes.OSUserCredential{
		Password:      "testFT4",
		Label:         "testFT3",
		Username:      "tes3t",
		KeyPairName:   "ConTestKeyPair",
		SSHPrivateKey: exampleSSHKey,
	}, id)
	assert.Nil(t, err2)
	assert.NotNil(t, credMod)
	// Delete
	errDel := system.DeleteResourceCredential(id)
	assert.Nil(t, errDel)
}
