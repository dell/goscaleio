package inttests

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestUploadPackage function to test upload packge with dummy path of packages
func TestUploadPackage(t *testing.T) {
	err := GC.UploadPackages("/test")
	assert.NotNil(t, err)
}

// TestParseCSV function to test parse csv function with dummy path of CSV file
func TestParseCSV(t *testing.T) {
	err := GC.ParseCSV("/test/test.csv")
	assert.NotNil(t, err)
}
