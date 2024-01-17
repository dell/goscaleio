package goscaleio

import (
	"errors"
	"math"
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

	client, err := NewClientWithArgs(svr.URL, "", math.MaxInt64, true, false)
	client.configConnect.Version = "4.5"
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
			client, err := NewClientWithArgs(svr.URL, "", math.MaxInt64, true, false)
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
			client, err := NewClientWithArgs(svr.URL, "", math.MaxInt64, true, false)
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
