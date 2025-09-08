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

package inttests

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetVTrees(t *testing.T) {
	allVTrees, err := C.GetVTrees()
	assert.Nil(t, err)
	assert.NotNil(t, allVTrees)
}

func TestGetVTreeByID(t *testing.T) {
	allVTrees, err := C.GetVTrees()
	assert.Nil(t, err)
	assert.NotNil(t, allVTrees)

	vTree, err := C.GetVTreeByID(allVTrees[0].ID)
	assert.Nil(t, err)
	assert.NotNil(t, vTree)

	vTree, err = C.GetVTreeByID(invalidIdentifier)
	assert.NotNil(t, err)
	assert.Nil(t, vTree)
}

func TestGetVTreeInstances(t *testing.T) {
	allVTrees, err := C.GetVTrees()
	assert.Nil(t, err)
	assert.NotNil(t, allVTrees)

	allVTrees, err = C.GetVTreeInstances([]string{allVTrees[0].ID})
	assert.Nil(t, err)
	assert.NotNil(t, allVTrees)

	allVTrees, err = C.GetVTreeInstances([]string{invalidIdentifier})
	assert.NotNil(t, err)
	assert.Nil(t, allVTrees)
}

func TestGetVTreeByVolumeID(t *testing.T) {
	allVTrees, err := C.GetVTrees()
	assert.Nil(t, err)
	assert.NotNil(t, allVTrees)

	vTree, err := C.GetVTreeByVolumeID(allVTrees[0].RootVolumes[0])
	assert.Nil(t, err)
	assert.NotNil(t, vTree)

	vTree, err = C.GetVTreeByVolumeID(invalidIdentifier)
	assert.NotNil(t, err)
	assert.Nil(t, vTree)
}
