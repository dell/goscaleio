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
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetAllDeployeService(t *testing.T) {
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	}))
	defer svr.Close()

	GC, err := NewGateway(svr.URL, "", "", true, true)
	if err != nil {
		t.Fatal(err)
	}

	templateDetails, err := GC.GetAllServiceDetails()
	assert.Equal(t, len(templateDetails), 0)
	assert.Nil(t, err)
}

func TestGetDeployeServiceByID(t *testing.T) {
	type testCase struct {
		id       string
		expected error
	}

	cases := []testCase{
		{
			id:       "sdnasgw",
			expected: nil,
		},
		{
			id:       "sdnasgw1",
			expected: errors.New("The template cannot be found"),
		},
	}

	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	}))
	defer svr.Close()

	for _, tc := range cases {
		tc := tc
		t.Run("", func(ts *testing.T) {
			GC, err := NewGateway(svr.URL, "", "", true, true)
			if err != nil {
				t.Fatal(err)
			}

			_, err = GC.GetServiceDetailsByID(tc.id, false)
			if err != nil {
				if tc.expected == nil {
					t.Errorf("Getting template by ID did not work as expected, \n\tgot: %s \n\twant: %v", err, tc.expected)
				} else {
					if err.Error() != tc.expected.Error() {
						t.Errorf("Getting template by ID did not work as expected, \n\tgot: %s \n\twant: %s", err, tc.expected)
					}
				}
			}
		})
	}
}

func TestGetDeployeServiceByFilters(t *testing.T) {
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	}))
	defer svr.Close()

	client, err := NewGateway(svr.URL, "", "", true, true)
	if err != nil {
		t.Fatal(err)
	}

	templates, err := client.GetServiceDetailsByFilter("Name", "Test")
	assert.Equal(t, len(templates), 0)
	assert.NotNil(t, err)
}

func TestGetDeployeServiceByIDNegative(t *testing.T) {
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w, `{"error":"Internal Server Error"}`)
	}))
	defer svr.Close()

	client, err := NewGateway(svr.URL, "", "", true, true)
	if err != nil {
		t.Fatal(err)
	}

	templates, err := client.GetServiceDetailsByID("Test", false)
	assert.Nil(t, templates)
	assert.NotNil(t, err)
}

func TestGetDeployeServiceByFiltersNegative(t *testing.T) {
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w, `{"error":"Internal Server Error"}`)
	}))
	defer svr.Close()

	client, err := NewGateway(svr.URL, "", "", true, true)
	if err != nil {
		t.Fatal(err)
	}

	templates, err := client.GetServiceDetailsByFilter("Name", "Test")
	assert.Nil(t, templates)
	assert.NotNil(t, err)
}

func TestGetAllDeployeServiceNegative(t *testing.T) {
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w, `{"error":"Internal Server Error"}`)
	}))
	defer svr.Close()

	client, err := NewGateway(svr.URL, "", "", true, true)
	if err != nil {
		t.Fatal(err)
	}

	templates, err := client.GetAllServiceDetails()
	assert.Nil(t, templates)
	assert.NotNil(t, err)
}
