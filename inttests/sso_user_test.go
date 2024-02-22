package inttests

import (
	"testing"

	types "github.com/dell/goscaleio/types/v1"
	"github.com/stretchr/testify/assert"
)

func TestCreateAndDeleteSSOUser(t *testing.T) {
	details, err := C.CreateSSOUser(&types.SSOUserCreateParam{
		UserName: "IntegrationTestSSOUser",
		Role:     "Monitor",
		Password: "Ssouser123!",
		Type:     "Local",
	})
	assert.Nil(t, err)
	assert.NotEmpty(t, details)

	err = C.DeleteSSOUser(details.ID)
	assert.Nil(t, err)
}
