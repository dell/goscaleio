package inttests

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestDeployUploadPackage function to test upload packge with dummy path of packages
func TestDeployUploadPackage(t *testing.T) {
	err := GC.UploadPackages("/test")
	assert.NotNil(t, err)
}

// TestDeployParseCSV function to test parse csv function with dummy path of CSV file
func TestDeployParseCSV(t *testing.T) {
	err := GC.ParseCSV("/test/test.csv")
	assert.NotNil(t, err)
}
