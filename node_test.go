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

func TestGetNodes(t *testing.T) {
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	}))
	defer svr.Close()

	client, err := NewGateway(svr.URL, "", "", true, false)
	client.version = "4.5"
	if err != nil {
		t.Fatal(err)
	}

	nodeDetails, err := client.GetAllNodes()
	assert.Equal(t, len(nodeDetails), 0)
	assert.Nil(t, err)
}

func TestGetNodeByID(t *testing.T) {
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
			expected: errors.New("The node cannot be found"),
		},
	}

	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	}))
	defer svr.Close()

	for _, tc := range cases {
		tc := tc
		t.Run("", func(ts *testing.T) {
			client, err := NewGateway(svr.URL, "", "", true, false)
			client.version = "4.5"
			if err != nil {
				t.Fatal(err)
			}

			_, err = client.GetNodeByID(tc.id)
			if err != nil {
				if tc.expected == nil {
					t.Errorf("Getting node by ID did not work as expected, \n\tgot: %s \n\twant: %v", err, tc.expected)
				} else {
					if err.Error() != tc.expected.Error() {
						t.Errorf("Getting node by ID did not work as expected, \n\tgot: %s \n\twant: %s", err, tc.expected)
					}
				}
			}
		})
	}
}

func TestGetNodePoolByID(t *testing.T) {
	type testCase struct {
		id       int
		expected error
	}

	cases := []testCase{
		{
			id:       1,
			expected: nil,
		},
		{
			id:       -100,
			expected: errors.New("The nodepool cannot be found"),
		},
	}

	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	}))
	defer svr.Close()

	for _, tc := range cases {
		tc := tc
		t.Run("", func(ts *testing.T) {
			client, err := NewGateway(svr.URL, "", "", true, false)
			client.version = "4.5"
			if err != nil {
				t.Fatal(err)
			}

			_, err = client.GetNodePoolByID(tc.id)
			if err != nil {
				if tc.expected == nil {
					t.Errorf("Getting nodepool by ID did not work as expected, \n\tgot: %s \n\twant: %v", err, tc.expected)
				} else {
					if err.Error() != tc.expected.Error() {
						t.Errorf("Getting nodepool by ID did not work as expected, \n\tgot: %s \n\twant: %s", err, tc.expected)
					}
				}
			}
		})
	}
}

func TestGetNodeByFilters(t *testing.T) {
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	}))
	defer svr.Close()

	client, err := NewGateway(svr.URL, "", "", true, false)
	client.version = "4.5"
	if err != nil {
		t.Fatal(err)
	}

	nodeDetails, err := client.GetNodeByFilters("ipAddress", "1.1.1.1")
	assert.Equal(t, len(nodeDetails), 0)
	assert.NotNil(t, err)
}

func TestGetNodePoolByName(t *testing.T) {
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	}))
	defer svr.Close()

	client, err := NewGateway(svr.URL, "", "", true, false)
	client.version = "4.5"
	if err != nil {
		t.Fatal(err)
	}

	NodePoolDetails, err := client.GetNodePoolByName("nodepool")
	assert.Nil(t, NodePoolDetails)
	assert.NotNil(t, err)
}

func TestGetNodePoolByNameError(t *testing.T) {
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintln(w, `{"error":"Resource not found"}`)
	}))
	defer svr.Close()

	client, err := NewGateway(svr.URL, "", "", true, false)
	client.version = "4.5"
	if err != nil {
		t.Fatal(err)
	}

	NodePoolDetails, err := client.GetNodePoolByName("nodepool")
	assert.Nil(t, NodePoolDetails)
	assert.NotNil(t, err)
}

func TestGetNodePoolByIDNegative(t *testing.T) {
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintln(w, `{"error":"Resource not found"}`)
	}))
	defer svr.Close()

	client, err := NewGateway(svr.URL, "", "", true, false)
	client.version = "4.5"
	if err != nil {
		t.Fatal(err)
	}

	NodePoolDetails, err := client.GetNodePoolByID(-100)
	assert.Nil(t, NodePoolDetails)
	assert.NotNil(t, err)
}

func TestGetNodeByIDNegative(t *testing.T) {
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintln(w, `{"error":"Resource not found"}`)
	}))
	defer svr.Close()

	client, err := NewGateway(svr.URL, "", "", true, false)
	client.version = "4.5"
	if err != nil {
		t.Fatal(err)
	}

	node, err := client.GetNodeByID("-100")
	assert.Nil(t, node)
	assert.NotNil(t, err)
}

func TestGetAllNodesNegative(t *testing.T) {
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w, `{"error":"Internal Server Error"}`)
	}))
	defer svr.Close()

	client, err := NewGateway(svr.URL, "", "", true, false)
	client.version = "4.5"
	if err != nil {
		t.Fatal(err)
	}

	nodes, err := client.GetAllNodes()
	assert.Nil(t, nodes)
	assert.NotNil(t, err)
}
