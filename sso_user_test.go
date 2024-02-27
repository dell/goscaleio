// Copyright © 2024 Dell Inc. or its subsidiaries. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//      http://www.apache.org/licenses/LICENSE-2.0
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package goscaleio

import (
	"errors"
	"fmt"
	"math"
	"net/http"
	"net/http/httptest"
	"testing"

	types "github.com/dell/goscaleio/types/v1"
	"github.com/stretchr/testify/assert"
)

func TestCreateSSOUser(t *testing.T) {
	type testCase struct {
		user     types.SSOUserCreateParam
		expected error
	}

	cases := []testCase{
		{
			user: types.SSOUserCreateParam{
				UserName: "testUser",
				Role:     "Monitor",
				Password: "default",
				Type:     "Local",
			},
			expected: nil,
		},
		{
			user: types.SSOUserCreateParam{
				UserName: "admin",
				Role:     "Monitor",
				Password: "default",
				Type:     "Local",
			},
			expected: errors.New("Invalid enum value"),
		},
	}
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	}))
	defer svr.Close()

	for _, tc := range cases {
		tc := tc
		t.Run("", func(ts *testing.T) {
			client, err := NewClientWithArgs(svr.URL, "", math.MaxInt64, true, false)
			if err != nil {
				t.Fatal(err)
			}

			_, err2 := client.CreateSSOUser(&tc.user)
			if err2 != nil {
				if tc.expected == nil {
					t.Errorf("Creating user did not work as expected, \n\tgot: %s \n\twant: %v", err2, tc.expected)
				} else {
					if err2.Error() != tc.expected.Error() {
						t.Errorf("Creating user did not work as expected, \n\tgot: %s \n\twant: %s", err2, tc.expected)
					}
				}
			}
		})
	}
}

func TestGetSSOUser(t *testing.T) {
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
			id:       "123",
			expected: errors.New("error getting user details"),
		},
	}
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	}))
	defer svr.Close()

	for _, tc := range cases {
		tc := tc
		t.Run("", func(ts *testing.T) {
			client, err := NewClientWithArgs(svr.URL, "", math.MaxInt64, true, false)
			if err != nil {
				t.Fatal(err)
			}

			_, err2 := client.GetSSOUser(tc.id)
			if err2 != nil {
				if tc.expected == nil {
					t.Errorf("Getting user details did not work as expected, \n\tgot: %s \n\twant: %v", err2, tc.expected)
				} else {
					if err2.Error() != tc.expected.Error() {
						t.Errorf("Getting user details did not work as expected, \n\tgot: %s \n\twant: %s", err2, tc.expected)
					}
				}
			}
		})
	}
}

func TestGetSSOUserByFilters(t *testing.T) {
	type testCase struct {
		username string
		expected error
	}
	cases := []testCase{
		{
			username: "admin",
			expected: nil,
		},
	}
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	}))
	defer svr.Close()

	for _, tc := range cases {
		tc := tc
		t.Run("", func(ts *testing.T) {
			client, err := NewClientWithArgs(svr.URL, "", math.MaxInt64, true, false)
			if err != nil {
				t.Fatal(err)
			}

			_, err2 := client.GetSSOUserByFilters(tc.username, "admin")
			if err2 != nil {
				if tc.expected == nil {
					t.Errorf("Getting user details did not work as expected, \n\tgot: %s \n\twant: %v", err2, tc.expected)
				} else {
					if err2.Error() != tc.expected.Error() {
						t.Errorf("Getting user details did not work as expected, \n\tgot: %s \n\twant: %s", err2, tc.expected)
					}
				}
			}
		})
	}
}

func TestModifySSOUser(t *testing.T) {
	type testCase struct {
		user     types.SSOUserModifyParam
		id       string
		expected error
	}
	cases := []testCase{
		{
			user: types.SSOUserModifyParam{
				Role: "Monitor",
			},
			id:       "eeb2dec800000001",
			expected: nil,
		},
		{
			user: types.SSOUserModifyParam{
				Role: "Monitor1",
			},
			id:       "eeb2dec800000001",
			expected: errors.New("Invalid enum value"),
		},
	}
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	}))
	defer svr.Close()

	for _, tc := range cases {
		tc := tc
		t.Run("", func(ts *testing.T) {
			client, err := NewClientWithArgs(svr.URL, "", math.MaxInt64, true, false)
			if err != nil {
				t.Fatal(err)
			}

			_, err2 := client.ModifySSOUser(tc.id, &tc.user)
			if err2 != nil {
				if tc.expected == nil {
					t.Errorf("Modifying user did not work as expected, \n\tgot: %s \n\twant: %v", err2, tc.expected)
				} else {
					if err2.Error() != tc.expected.Error() {
						t.Errorf("Modifying user did not work as expected, \n\tgot: %s \n\twant: %s", err2, tc.expected)
					}
				}
			}
		})
	}
}

func TestResetSSOUserPassword(t *testing.T) {
	type testCase struct {
		user     types.SSOUserModifyParam
		id       string
		expected error
	}
	cases := []testCase{
		{
			user: types.SSOUserModifyParam{
				Password: "default",
			},
			id:       "eeb2dec800000001",
			expected: nil,
		},
	}
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	}))
	defer svr.Close()

	for _, tc := range cases {
		tc := tc
		t.Run("", func(ts *testing.T) {
			client, err := NewClientWithArgs(svr.URL, "", math.MaxInt64, true, false)
			if err != nil {
				t.Fatal(err)
			}

			err2 := client.ResetSSOUserPassword(tc.id, &tc.user)
			if err2 != nil {
				if tc.expected == nil {
					t.Errorf("Resetting user password did not work as expected, \n\tgot: %s \n\twant: %v", err2, tc.expected)
				} else {
					if err2.Error() != tc.expected.Error() {
						t.Errorf("Resetting user password did not work as expected, \n\tgot: %s \n\twant: %s", err2, tc.expected)
					}
				}
			}
		})
	}
}

func TestDeleteSSOUser(t *testing.T) {
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
			expected: errors.New("HTTP 404 Not Found"),
		},
	}
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	}))
	defer svr.Close()

	for _, tc := range cases {
		tc := tc
		t.Run("", func(ts *testing.T) {
			client, err := NewClientWithArgs(svr.URL, "", math.MaxInt64, true, false)
			if err != nil {
				t.Fatal(err)
			}

			err2 := client.DeleteSSOUser(tc.id)
			if err2 != nil {
				if tc.expected == nil {
					t.Errorf("Deleting user did not work as expected, \n\tgot: %s \n\twant: %v", err2, tc.expected)
				} else {
					if err2.Error() != tc.expected.Error() {
						t.Errorf("Deleting user did not work as expected, \n\tgot: %s \n\twant: %s", err2, tc.expected)
					}
				}
			}
		})
	}
}

func TestGetSSOUserNegative(t *testing.T) {
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintln(w, `{"error":"User 93634330-6ffd-4d17-a22a-d3ec701e73d4 not found"}`)
	}))
	defer svr.Close()

	client, err := NewClientWithArgs(svr.URL, "", math.MaxInt64, true, false)
	if err != nil {
		t.Fatal(err)
	}

	details, err := client.GetSSOUser("93634330-6ffd-4d17-a22a-d3ec701e73d4")
	assert.Nil(t, details)
	assert.NotNil(t, err)
}

func TestDeleteSSOUserNegative(t *testing.T) {
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintln(w, `{"error":"User 93634330-6ffd-4d17-a22a-d3ec701e73d4 not found"}`)
	}))
	defer svr.Close()

	client, err := NewClientWithArgs(svr.URL, "", math.MaxInt64, true, false)
	if err != nil {
		t.Fatal(err)
	}

	err = client.DeleteSSOUser("93634330-6ffd-4d17-a22a-d3ec701e73d4")
	assert.NotNil(t, err)
}

func TestCreateSSOUserNegative(t *testing.T) {
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(w, `{"error":"Invalid enum value"}`)
	}))
	defer svr.Close()

	user := &types.SSOUserCreateParam{
		UserName: "testUser",
		Role:     "Monitor1",
		Password: "default",
		Type:     "Local",
	}

	client, err := NewClientWithArgs(svr.URL, "", math.MaxInt64, true, false)
	if err != nil {
		t.Fatal(err)
	}
	_, err = client.CreateSSOUser(user)
	assert.NotNil(t, err)
}
