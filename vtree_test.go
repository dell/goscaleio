package goscaleio

import (
	"encoding/json"
	"errors"
	"math"
	"net/http"
	"net/http/httptest"
	"testing"

	types "github.com/dell/goscaleio/types/v1"
	"github.com/stretchr/testify/assert"
)

func TestGetVTrees(t *testing.T) {
	t.Run("successful retrieval", func(t *testing.T) {
		svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusNoContent)
		}))
		defer svr.Close()

		client, err := NewClientWithArgs(svr.URL, "", math.MaxInt64, true, false)
		client.configConnect.Version = "3.6"
		if err != nil {
			t.Fatal(err)
		}

		vTreeDetails, err := client.GetVTrees()
		assert.Equal(t, len(vTreeDetails), 0)
		assert.Nil(t, err)
	})

	t.Run("error condition", func(t *testing.T) {
		svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			http.Error(w, `{"error": "500 Internal Server Error"}`, http.StatusInternalServerError)
		}))
		defer svr.Close()

		client, err := NewClientWithArgs(svr.URL, "", math.MaxInt64, true, false)
		client.configConnect.Version = "3.6"
		if err != nil {
			t.Fatal(err)
		}

		vTreeDetails, err := client.GetVTrees()
		assert.Nil(t, vTreeDetails)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "500 Internal Server Error")
	})
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
			expected: errors.New("404 Not Found"),
		},
	}

	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/instances/VTree::b21581e400000001" {
			w.WriteHeader(http.StatusOK)
			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte(`{"id": "b21581e400000001", "name": "test-vtree"}`))
		} else if r.URL.Path == "/api/instances/VTree::b21581e400000002" {
			w.WriteHeader(http.StatusNotFound)
			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte(`{"error": {"message": "The VTree was not found"}}`))
		}
	}))
	defer svr.Close()

	for _, tc := range cases {
		tc := tc
		t.Run("", func(_ *testing.T) {
			client, err := NewClientWithArgs(svr.URL, "", math.MaxInt64, true, false)
			if err != nil {
				t.Fatal(err)
			}

			_, err = client.GetVTreeByID(tc.id)
			if err != nil {
				if tc.expected == nil {
					t.Errorf("Getting VTree by ID did not work as expected, \n\tgot: %s \n\twant: %v", err, tc.expected)
				} else {
					if err.Error() != tc.expected.Error() {
						t.Errorf("Getting VTree by ID did not work as expected, \n\tgot: %s \n\twant: %s", err, tc.expected)
					}
				}
			} else {
				if tc.expected != nil {
					t.Errorf("Expected error but got none, \n\twant: %s", tc.expected)
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
			expected: errors.New("404 Not Found"),
		},
	}

	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/types/VTree/instances/action/queryBySelectedIds" {
			var payload types.VTreeQueryBySelectedIDsParam
			if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			if payload.IDs[0] == "b21581e400000001" {
				w.WriteHeader(http.StatusOK)
				w.Header().Set("Content-Type", "application/json")
				w.Write([]byte(`[{"id": "b21581e400000001", "name": "VTree1"}]`))
			} else if payload.IDs[0] == "b21581e400000002" {
				w.WriteHeader(http.StatusNotFound)
				w.Header().Set("Content-Type", "application/json")
				w.Write([]byte(`{"error": {"message": "Query selected Instances type: VTree - Got no statistics for id b21581e300000002. It doesn't exist"}}`))
			}
		}
	}))
	defer svr.Close()

	for _, tc := range cases {
		tc := tc
		t.Run("", func(_ *testing.T) {
			client, err := NewClientWithArgs(svr.URL, "", math.MaxInt64, true, false)
			if err != nil {
				t.Fatal(err)
			}

			_, err = client.GetVTreeInstances(tc.ids)
			if err != nil {
				if tc.expected == nil {
					t.Errorf("Getting VTree by ID did not work as expected, \n\tgot: %s \n\twant: %v", err, tc.expected)
				} else {
					if err.Error() != tc.expected.Error() {
						t.Errorf("Getting VTree by ID did not work as expected, \n\tgot: %s \n\twant: %s", err, tc.expected)
					}
				}
			} else {
				if tc.expected != nil {
					t.Errorf("Expected error but got none, \n\twant: %s", tc.expected)
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
			expected: errors.New("404 Not Found"),
		},
	}

	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/instances/Volume::b21581e400000002" {
			w.WriteHeader(http.StatusNotFound)
			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte(`{"error": {"message": "Invalid volume"}}`))
		}
	}))
	defer svr.Close()

	for _, tc := range cases {
		tc := tc
		t.Run("", func(t *testing.T) {
			client, err := NewClientWithArgs(svr.URL, "", math.MaxInt64, true, false)
			if err != nil {
				t.Fatal(err)
			}

			_, err = client.GetVTreeByVolumeID(tc.id)
			if err != nil {
				if tc.expected == nil {
					t.Errorf("Getting VTree by Volume ID did not work as expected, \n\tgot: %s \n\twant: %v", err, tc.expected)
				} else {
					if err.Error() != tc.expected.Error() {
						t.Errorf("Getting VTree by Volume ID did not work as expected, \n\tgot: %s \n\twant: %s", err, tc.expected)
					}
				}
			} else {
				if tc.expected != nil {
					t.Errorf("Expected error but got none, \n\twant: %s", tc.expected)
				}
			}
		})
	}
}
