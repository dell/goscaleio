package goscaleio

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDrvCfgQuerySystems(t *testing.T) {
	defaultExecFunc := executeFunc
	afterEach := func() {
		executeFunc = defaultExecFunc
		SCINIMockMode = false
	}
	tests := []struct {
		name      string
		setup     func()
		expectErr bool
		expectOut *[]ConfiguredCluster
	}{ //todo: test cases
		{
			name: "scini in mock mode",
			setup: func() {
				SCINIMockMode = true

			},
			expectErr: false,
			expectOut: &[]ConfiguredCluster{{
				SystemID: mockSystem,
				SdcID:    mockGUID,
			}},
		},
		{
			name: "passing test",
			setup: func() {
				executeFunc = func(_ string, _ ...string) ([]byte, error) {
					return []byte("MDM-ID aaaa SDC ID bbbb"), nil
				}
			},
			expectErr: false,
			expectOut: &[]ConfiguredCluster{{
				SystemID: "aaaa",
				SdcID:    "bbbb",
			}},
		},
		{
			name: "execute cmd returns failure",
			setup: func() {
				executeFunc = func(_ string, _ ...string) ([]byte, error) {
					return nil, assert.AnError
				}
			},
			expectErr: true,
		},
		{
			name: "bad output from exec",
			setup: func() {
				executeFunc = func(_ string, _ ...string) ([]byte, error) {
					return []byte("MDMAAAAA-ID aaaa SDC ID bbbb"), nil
				}
			},
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setup != nil {
				tt.setup()
			}
			defer afterEach()
			out, err := DrvCfgQuerySystems()
			if tt.expectErr {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err)
			}

			if tt.expectOut != nil {
				assert.Equal(t, tt.expectOut, out)
			}
		})
	}
}

func TestDrvCfgIsSDCInstalled(t *testing.T) {
	defaultStatFileFunc := statFileFunc
	afterEach := func() {
		statFileFunc = defaultStatFileFunc
		SCINIMockMode = false
	}
	tests := []struct {
		name        string
		setup       func()
		expectedOut bool
	}{
		{
			name: "scini in mock mode",
			setup: func() {
				SCINIMockMode = true

			},
			expectedOut: true,
		},
		{
			name:        "failing test - no SDC",
			expectedOut: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setup != nil {
				tt.setup()
			}
			defer afterEach()
			out := DrvCfgIsSDCInstalled()
			assert.Equal(t, out, tt.expectedOut)
		})
	}
}

func TestDrvCfgQueryGUIDd(t *testing.T) {
	defaultStatFileFunc := statFileFunc
	afterEach := func() {
		statFileFunc = defaultStatFileFunc
		SCINIMockMode = false
	}
	tests := []struct {
		name        string
		setup       func()
		expectedOut string
		expectError bool
	}{
		{
			name: "scini in mock mode",
			setup: func() {
				SCINIMockMode = true

			},
			expectedOut: mockGUID,
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setup != nil {
				tt.setup()
			}
			defer afterEach()
			out, err := DrvCfgQueryGUID()
			if tt.expectError {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err)
			}
			assert.Equal(t, out, tt.expectedOut)
		})
	}
}
