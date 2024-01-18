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

	"github.com/stretchr/testify/assert"
)

func TestGetNodes(t *testing.T) {
	allNodes, err := C.GetAllNodes()
	assert.Nil(t, err)
	assert.NotNil(t, allNodes)
}

func TestGetNodeByID(t *testing.T) {
	allNodes, err := C.GetAllNodes()
	assert.Nil(t, err)
	assert.NotNil(t, allNodes)

	if len(allNodes) > 0 {
		node, err := C.GetNodeByID(allNodes[0].RefID)
		assert.Nil(t, err)
		assert.NotNil(t, node)
	}

	node, err := C.GetNodeByID(invalidIdentifier)
	assert.NotNil(t, err)
	assert.Nil(t, node)
}

func TestGetNodeByFilters(t *testing.T) {
	allNodes, err := C.GetNodeByFilters("invalid", "invalid")
	assert.NotNil(t, err)
	assert.Nil(t, allNodes)
}

func TestGetNodePoolByID(t *testing.T) {
	nodePool, err := C.GetNodePoolByID(-2)
	assert.Nil(t, err)
	assert.Equal(t, nodePool.DeviceGroup.GroupName, "")
}
