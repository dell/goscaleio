package inttests

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetVTrees(t *testing.T) {
	allVTrees, err := C.GetVTrees(context.Background())
	assert.Nil(t, err)
	assert.NotNil(t, allVTrees)
}

func TestGetVTreeByID(t *testing.T) {
	ctx := context.Background()

	allVTrees, err := C.GetVTrees(ctx)
	assert.Nil(t, err)
	assert.NotNil(t, allVTrees)

	vTree, err := C.GetVTreeByID(ctx, allVTrees[0].ID)
	assert.Nil(t, err)
	assert.NotNil(t, vTree)

	vTree, err = C.GetVTreeByID(ctx, invalidIdentifier)
	assert.NotNil(t, err)
	assert.Nil(t, vTree)
}

func TestGetVTreeInstances(t *testing.T) {
	ctx := context.Background()

	allVTrees, err := C.GetVTrees(ctx)
	assert.Nil(t, err)
	assert.NotNil(t, allVTrees)

	allVTrees, err = C.GetVTreeInstances(ctx, []string{allVTrees[0].ID})
	assert.Nil(t, err)
	assert.NotNil(t, allVTrees)

	allVTrees, err = C.GetVTreeInstances(ctx, []string{invalidIdentifier})
	assert.NotNil(t, err)
	assert.Nil(t, allVTrees)
}

func TestGetVTreeByVolumeID(t *testing.T) {
	ctx := context.Background()

	allVTrees, err := C.GetVTrees(ctx)
	assert.Nil(t, err)
	assert.NotNil(t, allVTrees)

	vTree, err := C.GetVTreeByVolumeID(ctx, allVTrees[0].RootVolumes[0])
	assert.Nil(t, err)
	assert.NotNil(t, vTree)

	vTree, err = C.GetVTreeByVolumeID(ctx, invalidIdentifier)
	assert.NotNil(t, err)
	assert.Nil(t, vTree)
}
