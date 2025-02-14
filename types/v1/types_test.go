package goscaleio

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMetaData(t *testing.T) {
	vp := &VolumeParam{}
	assert.NotNil(t, vp.MetaData())
}

func TestIntString_MarshalJSON(t *testing.T) {
	is := IntString(123)
	obj, err := is.MarshalJSON()
	assert.NoError(t, err)
	assert.Equal(t, []byte(`"123"`), obj)
}

func TestGetBoolType(t *testing.T) {
	assert.Equal(t, "TRUE", GetBoolType(true))
	assert.Equal(t, "FALSE", GetBoolType(false))
}

func TestError(t *testing.T) {
	// Test case: error with message in details
	e := Error{
		Message:      errorWithDetails,
		ErrorDetails: []ErrorMessageDetails{{Error: "error1", ErrorMessage: "message1"}},
	}
	assert.EqualError(t, e, "message1")

	// Test case: error with untranslatable error in details
	e = Error{
		Message:      errorWithDetails,
		ErrorDetails: []ErrorMessageDetails{{Error: "error1"}},
	}
	assert.EqualError(t, e, errorWithDetails)

	// Test case: error with translatable error in details
	e = Error{
		Message:      errorWithDetails,
		ErrorDetails: []ErrorMessageDetails{{Error: "ALREADY_EXISTS"}},
	}
	assert.EqualError(t, e, "Already exists")

	// Test case: error without details
	e = Error{
		Message: "message2",
	}
	assert.EqualError(t, e, "message2")
}
