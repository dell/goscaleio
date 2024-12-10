package goscaleio

import (
	"context"
	"errors"
	"math"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetVTrees(t *testing.T) {
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	}))
	defer svr.Close()

	client, err := NewClientWithArgs(svr.URL, "", math.MaxInt64, true, false)
	client.configConnect.Version = "3.6"
	if err != nil {
		t.Fatal(err)
	}

	vTreeDetails, err := client.GetVTrees(context.Background())
	assert.Equal(t, len(vTreeDetails), 0)
	assert.Nil(t, err)
}

func TestGetVTreeByID(t *testing.T) {
	type testCase struct {
		id       string
		expected error
	}

	cases := []testCase{
		{
			id:       "b21581e400000001",
			expected: nil,
		},
		{
			id:       "b21581e400000002",
			expected: errors.New("The VTree was not found"),
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

			_, err = client.GetVTreeByID(context.Background(), tc.id)
			if err != nil {
				if tc.expected == nil {
					t.Errorf("Getting VTree by ID did not work as expected, \n\tgot: %s \n\twant: %v", err, tc.expected)
				} else {
					if err.Error() != tc.expected.Error() {
						t.Errorf("Getting VTree by ID did not work as expected, \n\tgot: %s \n\twant: %s", err, tc.expected)
					}
				}
			}
		})
	}
}

func TestGetVTreeInstances(t *testing.T) {
	type testCase struct {
		ids      []string
		expected error
	}

	cases := []testCase{
		{
			ids:      []string{"b21581e400000001"},
			expected: nil,
		},
		{
			ids:      []string{"b21581e400000002"},
			expected: errors.New("Query selected Instances type: VTree - Got no statistics for id b21581e300000002. It doesn't exist"),
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

			_, err = client.GetVTreeInstances(context.Background(), tc.ids)
			if err != nil {
				if tc.expected == nil {
					t.Errorf("Getting VTree by ID did not work as expected, \n\tgot: %s \n\twant: %v", err, tc.expected)
				} else {
					if err.Error() != tc.expected.Error() {
						t.Errorf("Getting VTree by ID did not work as expected, \n\tgot: %s \n\twant: %s", err, tc.expected)
					}
				}
			}
		})
	}
}

func TestGetVTreeByVolumeID(t *testing.T) {
	type testCase struct {
		id       string
		expected error
	}

	cases := []testCase{
		{
			id:       "3c855e2900000001",
			expected: nil,
		},
		{
			id:       "b21581e400000002",
			expected: errors.New("Invalid volume"),
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

			_, err = client.GetVTreeByVolumeID(context.Background(), tc.id)
			if err != nil {
				if tc.expected == nil {
					t.Errorf("Getting VTree by Volume ID did not work as expected, \n\tgot: %s \n\twant: %v", err, tc.expected)
				} else {
					if err.Error() != tc.expected.Error() {
						t.Errorf("Getting VTree by Volume ID did not work as expected, \n\tgot: %s \n\twant: %s", err, tc.expected)
					}
				}
			}
		})
	}
}
