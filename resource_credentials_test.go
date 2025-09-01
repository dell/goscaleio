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

// Copyright © 2025 Dell Inc. or its subsidiaries. All Rights Reserved.
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
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"net/http/httptest"
	"testing"

	types "github.com/dell/goscaleio/types/v1"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestResourceCredentialsGet(t *testing.T) {
	searchID := uuid.NewString()
	rcs := []types.ResourceCredential{
		{
			Credential: types.CredObj{
				Type: "exampleType",
				ID:   searchID,
			},
		},
		{
			Credential: types.CredObj{
				Type: "exampleType2",
				ID:   uuid.NewString(),
			},
		},
	}
	rc := types.ResourceCredentials{
		Credentials: rcs,
	}

	server := httptest.NewServer(http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
		switch req.RequestURI {
		case "/api/v1/Credential":
			resp.WriteHeader(http.StatusOK)
			content, err := json.Marshal(rc)
			if err != nil {
				t.Fatal(err)
			}

			resp.Write(content)
		case fmt.Sprintf("/api/v1/Credential/%s", searchID):
			resp.WriteHeader(http.StatusOK)
			var cred types.CredObj
			for _, val := range rcs {
				if val.Credential.ID == searchID {
					cred = val.Credential
				}
			}

			content, err := json.Marshal(cred)
			if err != nil {
				t.Fatal(err)
			}

			resp.Write(content)
		default:
			resp.WriteHeader(http.StatusNoContent)
		}
	}))
	defer server.Close()

	client, err := NewClientWithArgs(server.URL, "3.6", math.MaxInt64, true, false)
	if err != nil {
		t.Fatal(err)
	}

	system := NewSystem(client)

	response, err := system.GetResourceCredentials()
	if err != nil {
		t.Fatal(err)
	}

	if len(response.Credentials) != 2 {
		t.Errorf("expected %d, got %d", 2, len(response.Credentials))
	}

	res, err := system.GetResourceCredential(searchID)
	if err != nil || res == nil {
		t.Fatal(err)
	}
}

func TestResourceCredentialsCreateModifyDelete(t *testing.T) {
	searchID := uuid.NewString()
	cred := types.CredObj{
		Type: "exampleType",
		ID:   searchID,
	}

	server := httptest.NewServer(http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
		switch req.RequestURI {
		case "/api/v1/Credential":
			resp.WriteHeader(http.StatusOK)
			content, err := json.Marshal(cred)
			if err != nil {
				t.Fatal(err)
			}

			resp.Write(content)
		case fmt.Sprintf("/api/v1/Credential/%s", searchID):
			resp.WriteHeader(http.StatusOK)
			content, err := json.Marshal(cred)
			if err != nil {
				t.Fatal(err)
			}

			resp.Write(content)
		default:
			resp.WriteHeader(http.StatusNoContent)
		}
	}))
	defer server.Close()

	client, err := NewClientWithArgs(server.URL, "3.6", math.MaxInt64, true, false)
	if err != nil {
		t.Fatal(err)
	}

	system := NewSystem(client)

	// Snmpv2 defaults Test
	_, err = system.CreateNodeResourceCredential(
		types.ServerCredential{
			Username: "test",
		},
	)
	if err != nil {
		t.Fatal(err)
	}

	// Snmpv2 name set Test
	_, err = system.CreateNodeResourceCredential(
		types.ServerCredential{
			Username:              "test",
			SNMPv2CommunityString: "some-test-string",
		},
	)
	if err != nil {
		t.Fatal(err)
	}

	// Snmpv3 Assert error for level 2 no MD5 pass set
	_, err = system.CreateNodeResourceCredential(
		types.ServerCredential{
			Username:            "test",
			SNMPv3SecurityName:  "security-test",
			SNMPv3SecurityLevel: "2",
		},
	)
	assert.NotNil(t, err)

	// Snmpv3 Assert error for level 3 no Des pass set
	_, err = system.CreateNodeResourceCredential(
		types.ServerCredential{
			Username:                        "test",
			SNMPv3SecurityName:              "security-test",
			SNMPv3SecurityLevel:             "3",
			SNMPv3MD5AuthenticationPassword: "some-md5-pass",
		},
	)
	assert.NotNil(t, err)

	// Snmpv3 with level 1 should pass with no passwords set
	_, err = system.CreateNodeResourceCredential(
		types.ServerCredential{
			Username:            "test",
			SNMPv3SecurityName:  "security-test",
			SNMPv3SecurityLevel: "1",
		},
	)
	if err != nil {
		t.Fatal(err)
	}

	// Modify Successfully
	_, err = system.ModifyNodeResourceCredential(
		types.ServerCredential{
			Username:            "test2",
			SNMPv3SecurityName:  "security-test",
			SNMPv3SecurityLevel: "1",
		}, searchID,
	)
	if err != nil {
		t.Fatal(err)
	}

	// Create Switch RC Successfully
	_, err = system.CreateSwitchResourceCredential(
		types.IomCredential{
			Username: "test",
		},
	)
	if err != nil {
		t.Fatal(err)
	}

	// Expect failure: inconsistent credentials config
	_, err = system.CreateSwitchResourceCredential(
		types.IomCredential{
			SSHPrivateKey: "test",
			KeyPairName:   "",
		},
	)
	assert.Error(t, err)

	// Modify Switch RC Successfully
	_, err = system.ModifySwitchResourceCredential(
		types.IomCredential{
			Username: "test",
		},
		searchID,
	)
	if err != nil {
		t.Fatal(err)
	}

	// Expect failure: inconsistent credentials config
	_, err = system.ModifySwitchResourceCredential(
		types.IomCredential{
			SSHPrivateKey: "test",
			KeyPairName:   "",
		}, "123",
	)
	assert.Error(t, err)

	// Create VCenter RC Successfully
	_, err = system.CreateVCenterResourceCredential(
		types.VCenterCredential{
			Username: "test",
		},
	)
	if err != nil {
		t.Fatal(err)
	}

	// Modify VCenter RC Successfully
	_, err = system.ModifyVCenterResourceCredential(
		types.VCenterCredential{
			Username: "test",
		},
		searchID,
	)
	if err != nil {
		t.Fatal(err)
	}

	// Create ElementManager RC Successfully
	_, err = system.CreateElementManagerResourceCredential(
		types.EMCredential{
			Username: "test",
		},
	)
	if err != nil {
		t.Fatal(err)
	}

	// Modify ElementManager RC Successfully
	_, err = system.ModifyElementManagerResourceCredential(
		types.EMCredential{
			Username: "test",
		},
		searchID,
	)
	if err != nil {
		t.Fatal(err)
	}

	// Create ScaleIOCredential RC Successfully
	_, err = system.CreateScaleIOResourceCredential(
		types.ScaleIOCredential{
			AdminUsername: "test",
		},
	)
	if err != nil {
		t.Fatal(err)
	}

	// Modify ElementManager RC Successfully
	_, err = system.ModifyScaleIOResourceCredential(
		types.ScaleIOCredential{
			AdminUsername: "test",
		},
		searchID,
	)
	if err != nil {
		t.Fatal(err)
	}

	// Create PresentationServer RC Successfully
	_, err = system.CreatePresentationServerResourceCredential(
		types.PSCredential{
			Username: "test",
		},
	)
	if err != nil {
		t.Fatal(err)
	}

	// Modify PresentationServer RC Successfully
	_, err = system.ModifyPresentationServerResourceCredential(
		types.PSCredential{
			Username: "test",
		},
		searchID,
	)
	if err != nil {
		t.Fatal(err)
	}

	// Create OsAdmin RC Successfully
	_, err = system.CreateOsAdminResourceCredential(
		types.OSAdminCredential{
			Label: "test",
		},
	)
	if err != nil {
		t.Fatal(err)
	}

	// Expect failure: inconsistent credentials config
	_, err = system.CreateOsAdminResourceCredential(
		types.OSAdminCredential{
			SSHPrivateKey: "test",
			KeyPairName:   "",
		},
	)
	assert.Error(t, err)

	// Modify OsAdmin RC Successfully
	_, err = system.ModifyOsAdminResourceCredential(
		types.OSAdminCredential{
			Label: "test",
		},
		searchID,
	)
	if err != nil {
		t.Fatal(err)
	}

	// Expect failure: inconsistent credentials config
	_, err = system.ModifyOsAdminResourceCredential(
		types.OSAdminCredential{
			SSHPrivateKey: "test",
			KeyPairName:   "",
		}, "123",
	)
	assert.Error(t, err)

	// Create OsUser RC Successfully
	_, err = system.CreateOsUserResourceCredential(
		types.OSUserCredential{
			Username: "test",
		},
	)
	if err != nil {
		t.Fatal(err)
	}

	// Expect failure: inconsistent credentials config
	_, err = system.CreateOsUserResourceCredential(
		types.OSUserCredential{
			SSHPrivateKey: "test",
			KeyPairName:   "",
		},
	)
	assert.Error(t, err)

	// Modify OsUser RC Successfully
	_, err = system.ModifyOsUserResourceCredential(
		types.OSUserCredential{
			Username: "test",
		},
		searchID,
	)
	if err != nil {
		t.Fatal(err)
	}

	// Expect failure: inconsistent credentials config
	_, err = system.ModifyOsUserResourceCredential(
		types.OSUserCredential{
			SSHPrivateKey: "test",
			KeyPairName:   "",
		}, "123",
	)
	assert.Error(t, err)

	// Delete Should work successfully

	err = system.DeleteResourceCredential(searchID)
	if err != nil {
		t.Fatal(err)
	}
}
