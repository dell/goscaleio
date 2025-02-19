package goscaleio

import (
	"errors"
	"os"
	"syscall"
	"testing"

	"github.com/stretchr/testify/assert"
)

// MockSyscall is a mock implementation of Syscaller
type MockSyscall struct {
	ReturnErrno syscall.Errno
}

func (m MockSyscall) Syscall(trap, a1, a2, a3 uintptr) (uintptr, uintptr, syscall.Errno) {
	return 0, 0, m.ReturnErrno
}

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
	}{
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
	defaultOsOpen := openFileFunc

	afterEach := func() {
		statFileFunc = defaultStatFileFunc
		SCINIMockMode = false
		openFileFunc = defaultOsOpen
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
		{
			name: "error opening SDC device",
			setup: func() {
				defaultOsOpen = func(_ string) (*os.File, error) {
					return nil, assert.AnError
				}
			},
			expectedOut: "",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setup != nil {
				tt.setup()
			}
			defer afterEach()
			syscaller := MockSyscall{ReturnErrno: 0}
			out, err := DrvCfgQueryGUID(syscaller)
			if tt.expectError {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err)
			}
			assert.Equal(t, out, tt.expectedOut)
		})
	}
}

func TestDrvCfgQueryRescan(t *testing.T) {
	defaultOsOpen := openFileFunc
	afterEach := func() {
		openFileFunc = defaultOsOpen
	}
	tests := []struct {
		name        string
		setup       func()
		expectedOut string
		expectError bool
	}{
		{
			name: "error opening SDC device",
			setup: func() {
				defaultOsOpen = func(_ string) (*os.File, error) {
					return nil, errors.New("open error")
				}
			},
			expectedOut: "",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setup != nil {
				tt.setup()
			}
			defer afterEach()
			syscaller := MockSyscall{ReturnErrno: 0}
			out, err := DrvCfgQueryRescan(syscaller)
			if tt.expectError {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err)
			}
			assert.Equal(t, tt.expectedOut, out)
		})
	}
}

func TestIoctl(t *testing.T) {
	tests := []struct {
		fd, op, arg uintptr
		mockErrno   syscall.Errno
		expectedErr error
	}{
		{1, 2, 3, 0, nil},
		{4, 5, 6, 1, syscall.Errno(1)},
	}

	for _, tt := range tests {
		mockSyscall := MockSyscall{ReturnErrno: tt.mockErrno}
		err := ioctl(mockSyscall, tt.fd, tt.op, tt.arg)
		if !errors.Is(err, tt.expectedErr) {
			t.Errorf("expected %v, got %v", tt.expectedErr, err)
		}
	}
}

func Test_IO(t *testing.T) {
	tests := []struct {
		t, nr, expected uintptr
	}{
		{0x1, 0x2, _IOC(0x0, 0x1, 0x2, 0)},
		{0x3, 0x4, _IOC(0x0, 0x3, 0x4, 0)},
		{0x5, 0x6, _IOC(0x0, 0x5, 0x6, 0)},
	}

	for _, tt := range tests {
		result := _IO(tt.t, tt.nr)
		if result != tt.expected {
			t.Errorf("expected %v, got %v", tt.expected, result)
		}
	}
}

func TestIOC(t *testing.T) {
	dir := uintptr(0)
	tuint := uintptr(1)
	nr := uintptr(2)
	size := uintptr(3)

	expected := (dir << 30) | (tuint << 8) | nr | (size << 16)

	result := _IOC(dir, tuint, nr, size)

	if result != expected {
		t.Errorf("Expected: %v, but got: %v", expected, result)
	}
}
