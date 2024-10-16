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
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestNvmeHost(t *testing.T) {
	system := getSystem()
	name := fmt.Sprintf("nvme_%v", randString(10))
	newName := fmt.Sprintf("nvme_new_%v", randString(10))
	var hostID string

	t.Run("Create NVMe Host", func(t *testing.T) {
		assert.NotNil(t, system)
		nvmeHostParam := types.NvmeHostParam{
			Name:        name,
			Nqn:         fmt.Sprintf("nqn.2014-08.org.nvmexpress:uuid:%v", uuid.New()),
			MaxNumPaths: 4,
		}
		resp, err := system.CreateNvmeHost(nvmeHostParam)
		assert.Nil(t, err)
		assert.NotNil(t, resp.ID)
		hostID = resp.ID
	})

	t.Run("Get All NVMe Hosts", func(t *testing.T) {
		hosts, err := system.GetAllNvmeHosts()
		assert.Nil(t, err)
		assert.NotNil(t, hosts)
	})

	t.Run("Get NVMe Host By ID", func(t *testing.T) {
		host, err := system.GetNvmeHostByID(hostID)
		assert.Nil(t, err)
		assert.NotNil(t, host)
	})

	t.Run("Change NVMe Host Name", func(t *testing.T) {
		err := system.ChangeNvmeHostName(hostID, newName)
		assert.Nil(t, err)
	})

	t.Run("Change NVMe Host MaxNumPaths", func(t *testing.T) {
		err := system.ChangeNvmeHostMaxNumPaths(hostID, 6)
		assert.Nil(t, err)
	})

	t.Run("Change NVMe Host MaxNumSysPorts", func(t *testing.T) {
		err := system.ChangeNvmeHostMaxNumSysPorts(hostID, 8)
		assert.Nil(t, err)
	})

	t.Run("Get Host NVMe Controllers", func(t *testing.T) {
		controllers, err := system.GetHostNvmeControllers(types.NvmeHost{ID: hostID})
		assert.Nil(t, err)
		assert.NotNil(t, controllers)
	})

	t.Cleanup(func() {
		err := system.DeleteNvmeHost(hostID)
		assert.Nil(t, err)
	})
}
