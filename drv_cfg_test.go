package goscaleio

import (
	"errors"
	"fmt"
	"os"
	"syscall"
	"testing"
	"unsafe"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

// MockSyscall is a mock implementation of Syscaller
type MockSyscall struct {
	ReturnErrno syscall.Errno
	SyscallFunc func(trap, a1, a2, a3 uintptr) (uintptr, uintptr, syscall.Errno)
}

func (m MockSyscall) Syscall(trap, a1, a2, a3 uintptr) (uintptr, uintptr, syscall.Errno) {
	if m.SyscallFunc != nil {
		return m.SyscallFunc(trap, a1, a2, a3)
	}
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

func TestDrvCfgQueryGUID(t *testing.T) {
	tests := []struct {
		name        string
		mockMode    bool
		mockSyscall func(trap, a1, a2, a3 uintptr) (uintptr, uintptr, syscall.Errno) // simulate system call behavior
		mockOpen    func(name string) (*os.File, error)
		expected    string
		expectErr   bool
	}{
		{
			name:     "Mock mode",
			mockMode: true,
			expected: mockGUID,
		},
		{
			name: "Open device error",
			mockOpen: func(name string) (*os.File, error) {
				return nil, fmt.Errorf("open device error")
			},
			expectErr: true,
		},
		{
			name: "Ioctl_error",
			mockSyscall: func(trap, a1, a2, a3 uintptr) (uintptr, uintptr, syscall.Errno) {
				return 0, 0, syscall.EIO // Simulate an I/O error
			},
			mockOpen: func(name string) (*os.File, error) {
				return os.NewFile(0, name), nil
			},
			expectErr: true,
		},
		{
			name: "Invalid_RC",
			mockSyscall: func(trap, a1, a2, a3 uintptr) (uintptr, uintptr, syscall.Errno) {
				var buf ioctlGUID
				buf.rc[0] = 0x00                        // Set an invalid return code
				return 0, 0, 0
			},
			mockOpen: func(name string) (*os.File, error) {
				return os.NewFile(0, name), nil
			},
			expectErr: true,
		},
		{
			name: "Successful query",
			mockSyscall: func(trap, a1, a2, a3 uintptr) (uintptr, uintptr, syscall.Errno) {
				var buf ioctlGUID
				buf.rc[0] = 0x41
				uuidBytes, _ := uuid.Parse("D7C07724-A481-42D6-B1A7-0739A3F28BB0")
				copy(buf.uuid[:], uuidBytes[:])

				return 0, 0, 0
			},
			mockOpen: func(name string) (*os.File, error) {
				return os.NewFile(0, name), nil
			},
			expected: "D7C07724-A481-42D6-B1A7-0739A3F28BB0",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			SCINIMockMode = tt.mockMode
			openFileFunc = tt.mockOpen

			// Mock the ioctlWrapper function
			ioctlWrapper = func(_ Syscaller, _, _ uintptr, arg *ioctlGUID) error {
				if tt.name == "Ioctl_error" {
					return syscall.EIO // Simulate an I/O error
				}
				buf := (*ioctlGUID)(unsafe.Pointer(arg))
				if tt.name == "Invalid_RC" {
					buf.rc[0] = 0x00 // Set an invalid return code
				} else {
					buf.rc[0] = 0x41
					uuidBytes, _ := uuid.Parse("D7C07724-A481-42D6-B1A7-0739A3F28BB0")
					copy(buf.uuid[:], uuidBytes[:])
				}

				return nil
			}

			syscaller := &MockSyscall{
				SyscallFunc: tt.mockSyscall,
			}

			guid, err := DrvCfgQueryGUID(syscaller)
			if tt.expectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, guid)
			}
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

// func TestIoctlWrapper(t *testing.T) {
// 	tests := []struct {
// 		fd, op, arg uintptr
// 		mockErrno   syscall.Errno
// 		expectedErr error
// 	}{
// 		{
// 			1,
// 			2,
// 			uintptr(unsafe.Pointer(&ioctlGUID{})),
// 			0,
// 			nil,
// 		},
// 		{
// 			4,
// 			5,
// 			uintptr(unsafe.Pointer(&ioctlGUID{})),
// 			syscall.EINVAL,
// 			syscall.EINVAL,
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(fmt.Sprintf("fd=%d,op=%d,arg=%d", tt.fd, tt.op, tt.arg), func(t *testing.T) {
// 			fmt.Printf("Testing with fd=%d, op=%d, arg=%d, mockErrno=%d\n", tt.fd, tt.op, tt.arg, tt.mockErrno)
// 			mockSyscall := MockSyscall{
// 				SyscallFunc: func(trap, a1, a2, a3 uintptr) (uintptr, uintptr, syscall.Errno) {
// 					fmt.Printf("MockSyscall called with trap=%d, a1=%d, a2=%d, a3=%d\n", trap, a1, a2, a3)
// 					if tt.mockErrno != 0 {
// 						fmt.Printf("Returning mockErrno=%d\n", tt.mockErrno)
// 						return 0, 0, tt.mockErrno
// 					}
// 					return 0, 0, 0
// 				},
// 			}
// 			err := ioctlWrapper(mockSyscall, tt.fd, tt.op, tt.arg)
// 			if !errors.Is(err, tt.expectedErr) {
// 				t.Errorf("expected %v, got %v", tt.expectedErr, err)
// 			}
// 		})
// 	}
// }

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
