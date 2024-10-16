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
	"fmt"
	"testing"

	types "github.com/dell/goscaleio/types/v1"
	"github.com/stretchr/testify/assert"
)

func TestSdt(t *testing.T) {
	system := getSystem()
	newName := fmt.Sprintf("example-sdt_%v", randString(10))
	var sdtID string

	t.Run("Create sdt", func(t *testing.T) {
		assert.NotNil(t, system)
		pd := getProtectionDomain(t)
		assert.NotNil(t, pd)

		sdtParam := &types.SdtParam{
			Name:               "example-sdt",
			IPList:             []*types.SdtIPList{{SdtIP: types.SdtIP{IP: "192.168.0.1", Role: "StorageAndHost"}}},
			StoragePort:        12200,
			NvmePort:           4420,
			DiscoveryPort:      8009,
			ProtectionDomainID: pd.ProtectionDomain.ID,
		}
		resp, err := pd.CreateSdt(sdtParam)
		assert.Nil(t, err)
		assert.NotNil(t, resp.ID)
		sdtID = resp.ID
	})

	t.Run("Get All Sdt", func(t *testing.T) {
		hosts, err := system.GetAllSdts()
		assert.Nil(t, err)
		assert.NotNil(t, hosts)
	})

	t.Run("Get Sdt By ID", func(t *testing.T) {
		host, err := system.GetSdtByID(sdtID)
		assert.Nil(t, err)
		assert.NotNil(t, host)
	})

	t.Run("Rename Sdt", func(t *testing.T) {
		err := system.RenameSdt(sdtID, newName)
		assert.Nil(t, err)
	})

	t.Run("Set Sdt NvmePort", func(t *testing.T) {
		err := system.SetSdtNvmePort(sdtID, 4422)
		assert.Nil(t, err)
	})

	t.Run("Set Sdt StoragePort", func(t *testing.T) {
		err := system.SetSdtStoragePort(sdtID, 12300)
		assert.Nil(t, err)
	})

	t.Run("Set Sdt DiscoveryPort", func(t *testing.T) {
		err := system.SetSdtDiscoveryPort(sdtID, 8010)
		assert.Nil(t, err)
	})

	t.Run("Add Sdt Target IP", func(t *testing.T) {
		err := system.AddSdtTargetIP(sdtID, "192.168.0.2", "StorageAndHost")
		assert.Nil(t, err)
	})

	t.Run("Modify Sdt IP and Role", func(t *testing.T) {
		err := system.ModifySdtIPRole(sdtID, "192.168.0.2", "StorageOnly")
		assert.Nil(t, err)
	})

	t.Run("Add Sdt Target IP", func(t *testing.T) {
		err := system.RemoveSdtTargetIP(sdtID, "192.168.0.2")
		assert.Nil(t, err)
	})



	t.Cleanup(func() {
		err := system.DeleteNvmeHost(sdtID)
		assert.Nil(t, err)
	})
}
