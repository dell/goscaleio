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

	details, err = C.GetSSOUser(details.ID)
	assert.Nil(t, err)
	assert.NotEmpty(t, details)

	details1, err := C.GetSSOUserByFilters("username", "IntegrationTestSSOUser")
	assert.Nil(t, err)
	assert.NotEmpty(t, details1)

	details, err = C.ModifySSOUser(details.ID, &types.SSOUserModifyParam{
		Role: "Technician",
	})
	assert.Nil(t, err)
	assert.NotEmpty(t, details)

	err = C.ResetSSOUserPassword(details.ID, &types.SSOUserModifyParam{Password: "Ssouser1234#"})
	assert.Nil(t, err)

	err = C.DeleteSSOUser(details.ID)
	assert.Nil(t, err)
}
