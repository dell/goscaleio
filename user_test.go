package goscaleio

import (
	"encoding/json"
	"errors"
	"fmt"
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
		server   *httptest.Server
	}
	cases := map[string]testCase{
		"success": {
			user: types.UserParam{
				Name:     "testUser",
				UserRole: "Monitor",
				Password: "default",
			},
			server: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				w.WriteHeader(http.StatusCreated)
			})),
		},
		"invalid user role": {
			user: types.UserParam{
				Name:     "newUser",
				UserRole: "Role",
				Password: "password",
			},
			expected: errors.New("userRole should get on Monitor, Configure, Administrator, Security, FrontendConfig, BackendConfig, but its value is Role"),
			server: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte(`{"message":"userRole should get on Monitor, Configure, Administrator, Security, FrontendConfig, BackendConfig, but its value is Role","httpStatusCode":400,"errorCode":0}`))
			})),
		},
	}

	for name, tc := range cases {
		tc := tc
		t.Run(name, func(_ *testing.T) {
			defer tc.server.Close()

			client, err := NewClientWithArgs(tc.server.URL, "", math.MaxInt64, true, false)
			if err != nil {
				t.Fatal(err)
			}

			s := System{
				client: client,
			}
			_, err2 := s.CreateUser(&tc.user)
			errorCheck(t, tc.expected, err2, "CreateUser")
		})
	}
}

func TestGetUser(t *testing.T) {
	type testCase struct {
		id       string
		name     string
		expected error
		server   *httptest.Server
	}
	cases := map[string]testCase{
		"success": {
			id:       "eeb2dec800000001",
			name:     "",
			expected: nil,
			server: httptest.NewServer(http.HandlerFunc(func(_ http.ResponseWriter, _ *http.Request) {
			})),
		},
		"error with API call": {
			id:       "eeb2dec800000001",
			name:     "",
			expected: fmt.Errorf("unable to get user"),
			server: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte(`{"message":"unable to get user","httpStatusCode":400,"errorCode":0}`))
			})),
		},
	}

	for name, tc := range cases {
		tc := tc
		t.Run(name, func(_ *testing.T) {
			defer tc.server.Close()

			client, err := NewClientWithArgs(tc.server.URL, "", math.MaxInt64, true, false)
			if err != nil {
				t.Fatal(err)
			}

			s := System{
				System: &types.System{ID: "system-id"},
				client: client,
			}
			_, err2 := s.GetUser()
			errorCheck(t, tc.expected, err2, "GetUser")
		})
	}
}

func TestGetUserByIDName(t *testing.T) {
	type testCase struct {
		id       string
		name     string
		expected error
		server   *httptest.Server
	}
	cases := map[string]testCase{
		"success with first call to get by ID": {
			id:       "eeb2dec800000001",
			name:     "",
			expected: nil,
			server: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				switch r.RequestURI {
				case "/api/instances/User::eeb2dec800000001":
					data, err := json.Marshal(&types.User{})
					if err != nil {
						t.Fatal(err)
					}
					w.Write(data)
					return
				default:
					http.NotFound(w, r)
				}
			})),
		},
		"error with first API call to get by id": {
			id:       "eeb2dec800000001",
			name:     "",
			expected: fmt.Errorf("unable to get user by name"),
			server: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				switch r.RequestURI {
				case "/api/instances/User::eeb2dec800000001":
					w.WriteHeader(http.StatusBadRequest)
					w.Write([]byte(`{"message":"unable to get user by name","httpStatusCode":400,"errorCode":0}`))
				default:
					http.NotFound(w, r)
				}
			})),
		},
		"success with second API call to get by name": {
			id:       "",
			name:     "known-user-name-2",
			expected: nil,
			server: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				switch r.RequestURI {
				case "/api/instances/User::eeb2dec800000001":
					data, err := json.Marshal(&types.User{})
					if err != nil {
						t.Fatal(err)
					}
					w.Write(data)
				case "/api/instances/System::system-id/relationships/User":
					users := make([]types.User, 0)
					users = append(users, types.User{Name: "known-user-name-1"})
					users = append(users, types.User{Name: "known-user-name-2"})
					users = append(users, types.User{Name: "known-user-name-3"})
					data, err := json.Marshal(users)
					if err != nil {
						t.Fatal(err)
					}
					w.Write(data)
				default:
					http.NotFound(w, r)
				}
			})),
		},
		"error from second API call to get by user name": {
			id:       "",
			name:     "known-user-name-1",
			expected: fmt.Errorf("unable to get user by name"),
			server: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				switch r.RequestURI {
				case "/api/instances/User::eeb2dec800000001":
					data, err := json.Marshal(&types.User{})
					if err != nil {
						t.Fatal(err)
					}
					w.Write(data)
				case "/api/instances/System::system-id/relationships/User":
					w.WriteHeader(http.StatusBadRequest)
					w.Write([]byte(`{"message":"unable to get user by name","httpStatusCode":400,"errorCode":0}`))
				default:
					http.NotFound(w, r)
				}
			})),
		},
		"unable to find user by name": {
			id:       "",
			name:     "unknown-user-name",
			expected: errors.New("couldn't find user by name"),
			server: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				switch r.RequestURI {
				case "/api/types/User::eeb2dec800000001":
					data, err := json.Marshal(&types.User{})
					if err != nil {
						t.Fatal(err)
					}
					w.Write(data)
				case "/api/instances/System::system-id/relationships/User":
					users := make([]types.User, 0)
					users = append(users, types.User{Name: "known-user-name-1"})
					users = append(users, types.User{Name: "known-user-name-2"})
					data, err := json.Marshal(users)
					if err != nil {
						t.Fatal(err)
					}
					w.Write(data)
				default:
					http.NotFound(w, r)
				}
			})),
		},
		"invalid request with empty parameters": {
			id:       "",
			name:     "",
			expected: errors.New("user name or ID is mandatory, please enter a valid value"),
			server:   httptest.NewServer(http.HandlerFunc(func(_ http.ResponseWriter, _ *http.Request) {})),
		},
	}

	for name, tc := range cases {
		tc := tc
		t.Run(name, func(_ *testing.T) {
			defer tc.server.Close()
			client, err := NewClientWithArgs(tc.server.URL, "", math.MaxInt64, true, false)
			if err != nil {
				t.Fatal(err)
			}

			s := System{
				client: client,
				System: &types.System{ID: "system-id"},
			}
			_, err2 := s.GetUserByIDName(tc.id, tc.name)
			errorCheck(t, tc.expected, err2, "GetUserByIDName")
		})
	}
}

func TestRemoveUser(t *testing.T) {
	type testCase struct {
		id       string
		expected error
		server   *httptest.Server
	}
	cases := map[string]testCase{
		"success": {
			id:       "eeb2dec800000001",
			expected: nil,
			server: httptest.NewServer(http.HandlerFunc(func(_ http.ResponseWriter, _ *http.Request) {
			})),
		},
		"user not found": {
			id:       "eeb2dec800000005",
			expected: errors.New("User not found. Please check that you have the correct user name"),
			server: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte(`{"message":"User not found. Please check that you have the correct user name","httpStatusCode":400,"errorCode":0}`))
			})),
		},
	}

	for name, tc := range cases {
		defer tc.server.Close()

		tc := tc
		t.Run(name, func(_ *testing.T) {
			client, err := NewClientWithArgs(tc.server.URL, "", math.MaxInt64, true, false)
			if err != nil {
				t.Fatal(err)
			}

			s := System{
				client: client,
			}
			err2 := s.RemoveUser(tc.id)
			errorCheck(t, tc.expected, err2, "RemoveUser")
		})
	}
}

func TestSetRole(t *testing.T) {
	type testCase struct {
		role     types.UserRoleParam
		id       string
		expected error
		server   *httptest.Server
	}
	cases := map[string]testCase{
		"success": {
			role: types.UserRoleParam{
				UserRole: "Monitor",
			},
			id:       "eeb2dec800000001",
			expected: nil,
			server: httptest.NewServer(http.HandlerFunc(func(_ http.ResponseWriter, _ *http.Request) {
			})),
		},
		"error from API": {
			role: types.UserRoleParam{
				UserRole: "any",
			},
			id:       "eeb2dec800000001",
			expected: errors.New("userRole should get one of the following values: Monitor, Configure, Administrator, Security, FrontendConfig, BackendConfig, but its value is any"),
			server: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte(`{"message":"userRole should get one of the following values: Monitor, Configure, Administrator, Security, FrontendConfig, BackendConfig, but its value is any","httpStatusCode":400,"errorCode":0}`))
			})),
		},
	}

	for name, tc := range cases {
		defer tc.server.Close()

		tc := tc
		t.Run(name, func(_ *testing.T) {
			client, err := NewClientWithArgs(tc.server.URL, "", math.MaxInt64, true, false)
			if err != nil {
				t.Fatal(err)
			}

			s := System{
				client: client,
			}
			err2 := s.SetUserRole(&tc.role, tc.id)
			errorCheck(t, tc.expected, err2, "SetUserRole")
		})
	}
}

func errorCheck(t *testing.T, expected error, err error, name string) {
	if expected != nil && err == nil {
		t.Errorf("%s did not work as expected, \n\tgot: %s \n\twant: %v", name, err, expected)
	}
	if err != nil {
		if expected == nil || err.Error() != expected.Error() {
			t.Errorf("%s did not work as expected, \n\tgot: %s \n\twant: %v", name, err, expected)
		}
	}
}
