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
