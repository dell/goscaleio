package inttests

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestDeployUploadPackage function to test upload packge with dummy path of packages
func TestDeployUploadPackage(t *testing.T) {
	_, err := GC.UploadPackages("/home/krunal/Work/Software/abc.txt")
	assert.NotNil(t, err)
}

// TestDeployParseCSV function to test parse csv function with dummy path of CSV file
func TestDeployParseCSV(t *testing.T) {
	err := GC.ParseCSV("/test/test.csv")
	assert.NotNil(t, err)
}

// TestDeployGetPackage function to test Get Packge Details function
func TestDeployGetPackgeDetails(t *testing.T) {
	res, err := GC.GetPackgeDetails()
	assert.NotNil(t, res)
	assert.Nil(t, err)
}
