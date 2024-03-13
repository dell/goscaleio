package goscaleio

import (
	"errors"
	"math"
	"net/http"
	"net/http/httptest"
	"testing"

	types "github.com/dell/goscaleio/types/v1"
)

func TestCreateUser(t *testing.T) {
	type testCase struct {
		user     types.UserParam
		expected error
	}
	cases := []testCase{
		{
			user: types.UserParam{
				Name:     "testUser",
				UserRole: "Monitor",
				Password: "default",
			},
			expected: nil,
		},
		{
			user: types.UserParam{
				Name:     "newUser",
				UserRole: "Role",
				Password: "password",
			},
			expected: errors.New("userRole should get on Monitor, Configure, Administrator, Security, FrontendConfig, BackendConfig, but its value is Role"),
		},
	}
	svr := httptest.NewServer(http.HandlerFunc(func(_ http.ResponseWriter, _ *http.Request) {
	}))
	defer svr.Close()

	for _, tc := range cases {
		tc := tc
		t.Run("", func(_ *testing.T) {
			client, err := NewClientWithArgs(svr.URL, "", math.MaxInt64, true, false)
			if err != nil {
				t.Fatal(err)
			}

			s := System{
				client: client,
			}
			_, err2 := s.CreateUser(&tc.user)
			if err2 != nil {
				if tc.expected == nil {
					t.Errorf("Creating User did not work as expected, \n\tgot: %s \n\twant: %v", err2, tc.expected)
				} else {
					if err2.Error() != tc.expected.Error() {
						t.Errorf("Creating User did not work as expected, \n\tgot: %s \n\twant: %s", err2, tc.expected)
					}
				}
			}
		})
	}
}

func TestGetUserByIDName(t *testing.T) {
	type testCase struct {
		id       string
		name     string
		expected error
	}
	cases := []testCase{
		{
			id:       "eeb2dec800000001",
			name:     "",
			expected: nil,
		},
		{
			id:       "",
			name:     "",
			expected: errors.New("user name or ID is mandatory, please enter a valid value"),
		},
	}
	svr := httptest.NewServer(http.HandlerFunc(func(_ http.ResponseWriter, _ *http.Request) {
	}))
	defer svr.Close()

	for _, tc := range cases {
		tc := tc
		t.Run("", func(_ *testing.T) {
			client, err := NewClientWithArgs(svr.URL, "", math.MaxInt64, true, false)
			if err != nil {
				t.Fatal(err)
			}

			s := System{
				client: client,
			}
			_, err2 := s.GetUserByIDName(tc.id, tc.name)
			if err2 != nil {
				if tc.expected == nil {
					t.Errorf("Creating User did not work as expected, \n\tgot: %s \n\twant: %v", err2, tc.expected)
				} else {
					if err2.Error() != tc.expected.Error() {
						t.Errorf("Creating User did not work as expected, \n\tgot: %s \n\twant: %s", err2, tc.expected)
					}
				}
			}
		})
	}
}

func TestRemoveUser(t *testing.T) {
	type testCase struct {
		id       string
		expected error
	}
	cases := []testCase{
		{
			id:       "eeb2dec800000001",
			expected: nil,
		},
		{
			id:       "eeb2dec800000005",
			expected: errors.New("User not found. Please check that you have the correct user name"),
		},
	}
	svr := httptest.NewServer(http.HandlerFunc(func(_ http.ResponseWriter, _ *http.Request) {
	}))
	defer svr.Close()

	for _, tc := range cases {
		tc := tc
		t.Run("", func(_ *testing.T) {
			client, err := NewClientWithArgs(svr.URL, "", math.MaxInt64, true, false)
			if err != nil {
				t.Fatal(err)
			}

			s := System{
				client: client,
			}
			err2 := s.RemoveUser(tc.id)
			if err2 != nil {
				if tc.expected == nil {
					t.Errorf("Removing User did not work as expected, \n\tgot: %s \n\twant: %v", err2, tc.expected)
				} else {
					if err2.Error() != tc.expected.Error() {
						t.Errorf("Removing User did not work as expected, \n\tgot: %s \n\twant: %s", err2, tc.expected)
					}
				}
			}
		})
	}
}

func TestSetRole(t *testing.T) {
	type testCase struct {
		role     types.UserRoleParam
		id       string
		expected error
	}
	cases := []testCase{
		{
			role: types.UserRoleParam{
				UserRole: "Monitor",
			},
			id:       "eeb2dec800000001",
			expected: nil,
		},
		{
			role: types.UserRoleParam{
				UserRole: "any",
			},
			id:       "eeb2dec800000001",
			expected: errors.New("userRole should get one of the following values: Monitor, Configure, Administrator, Security, FrontendConfig, BackendConfig, but its value is any"),
		},
	}
	svr := httptest.NewServer(http.HandlerFunc(func(_ http.ResponseWriter, _ *http.Request) {
	}))
	defer svr.Close()

	for _, tc := range cases {
		tc := tc
		t.Run("", func(_ *testing.T) {
			client, err := NewClientWithArgs(svr.URL, "", math.MaxInt64, true, false)
			if err != nil {
				t.Fatal(err)
			}

			s := System{
				client: client,
			}
			err2 := s.SetUserRole(&tc.role, tc.id)
			if err2 != nil {
				if tc.expected == nil {
					t.Errorf("Removing User did not work as expected, \n\tgot: %s \n\twant: %v", err2, tc.expected)
				} else {
					if err2.Error() != tc.expected.Error() {
						t.Errorf("Removing User did not work as expected, \n\tgot: %s \n\twant: %s", err2, tc.expected)
					}
				}
			}
		})
	}
}
