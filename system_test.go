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
	"github.com/stretchr/testify/assert"
)

func TestModifyPerformanceProfile(t *testing.T) {
	type testCase struct {
		perfProfile string
		expected    error
	}

	cases := []testCase{
		{
			"HighPerformance",
			nil,
		},

		{
			"Compact1",
			errors.New("500 Internal Server Error"),
		},
	}

	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/instances/System/action/setMdmPerformanceParameters" {
			if r.Method == http.MethodPost {
				var param types.ChangeMdmPerfProfile
				_ = json.NewDecoder(r.Body).Decode(&param)
				switch param.PerfProfile {
				case "Compact1":
					w.WriteHeader(http.StatusInternalServerError)
					w.Write([]byte(`{"error": "perfProfile should get one of the following values: Compact, HighPerformance, but its value is Compact1"}`))
				default:
					w.WriteHeader(http.StatusOK)
				}
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

			s := System{
				client: client,
			}

			err = s.ModifyPerformanceProfileMdmCluster(tc.perfProfile)
			if err != nil {
				if tc.expected == nil {
					t.Errorf("Modifying performance profile of MDM cluster did not work as expected, \n\tgot: %s \n\twant: %v", err, tc.expected)
				} else {
					if err.Error() != tc.expected.Error() {
						t.Errorf("Modifying performance profile of MDM cluster did not work as expected, \n\tgot: %s \n\twant: %s", err, tc.expected)
					}
				}
			}
		})
	}
}

func TestAddStandByMDM(t *testing.T) {
	type testCase struct {
		ips           []string
		role          string
		server        *httptest.Server
		expected      string
		expectedError error
	}

	cases := map[string]testCase{
		"success": {
			ips:  []string{"10.xx.xx.xxx"},
			role: "Manager",
			server: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				w.WriteHeader(http.StatusOK)
				data, err := json.Marshal(types.Mdm{
					ID: "mdm-id-1",
				})
				if err != nil {
					t.Fatal(err)
				}
				w.Write(data)
			})),
			expected:      "mdm-id-1",
			expectedError: nil,
		},
		"error with API call": {
			ips:  []string{"10.xx.xx.xxx"},
			role: "Manager",
			server: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				w.WriteHeader(http.StatusInternalServerError)
			})),
			expected:      "",
			expectedError: fmt.Errorf("EOF"),
		},
	}

	for name, tc := range cases {
		tc := tc
		defer tc.server.Close()

		t.Run(name, func(_ *testing.T) {
			client, err := NewClientWithArgs(tc.server.URL, "", math.MaxInt64, true, false)
			if err != nil {
				t.Fatal(err)
			}

			s := System{
				client: client,
			}

			payload := types.StandByMdm{
				IPs:  tc.ips,
				Role: tc.role,
			}
			result, err := s.AddStandByMdm(&payload)
			assert.Equal(t, tc.expected, result)
			assert.Equal(t, tc.expectedError, err)
		})
	}
}

func TestRemoveStandByMDM(t *testing.T) {
	type testCase struct {
		mdmID    string
		expected error
	}

	cases := []testCase{
		{
			"1d9004d91b4ba503",
			nil,
		},

		{
			"1d9004d91b4ba504",
			errors.New("404 Not Found"),
		},
		{
			"1d9004d91b4ba505",
			errors.New("500 Internal Server Error"),
		},
	}

	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/instances/System/action/removeStandbyMdm" {
			if r.Method == http.MethodPost {
				var param types.RemoveStandByMdmParam
				_ = json.NewDecoder(r.Body).Decode(&param)
				switch param.ID {
				case "1d9004d91b4ba504":
					w.WriteHeader(http.StatusNotFound)
					w.Write([]byte(`{"error": "404 Not Found"}`))
				case "1d9004d91b4ba505":
					w.WriteHeader(http.StatusInternalServerError)
					w.Write([]byte(`{"error": "The MDM could not be found"}`))
				default:
					w.WriteHeader(http.StatusOK)
				}
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

			s := System{
				client: client,
			}

			err = s.RemoveStandByMdm(tc.mdmID)
			if err != nil {
				if tc.expected == nil {
					t.Errorf("Adding standby mdm did not work as expected, \n\tgot: %s \n\twant: %v", err, tc.expected)
				} else {
					if err.Error() != tc.expected.Error() {
						t.Errorf("Adding standby mdm did not work as expected, \n\tgot: %s \n\twant: %s", err, tc.expected)
					}
				}
			} else if tc.expected != nil {
				t.Errorf("Expected error but got nil, \n\twant: %s", tc.expected)
			}
		})
	}
}

func TestChangeMDMOwnership(t *testing.T) {
	type testCase struct {
		mdmID    string
		expected error
	}

	cases := []testCase{
		{
			"7f328d0b71711802",
			nil,
		},

		{
			"1d9004d91b4ba504",
			errors.New("500 Internal Server Error"),
		},
	}

	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/instances/System/action/changeMdmOwnership" {
			if r.Method == http.MethodPost {
				var param types.ChangeMdmOwnerShip
				_ = json.NewDecoder(r.Body).Decode(&param)
				switch param.ID {
				case "1d9004d91b4ba504":
					w.WriteHeader(http.StatusInternalServerError)
					w.Write([]byte(`{"error": "The MDM could not be found"}`))
				default:
					w.WriteHeader(http.StatusOK)
				}
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

			s := System{
				client: client,
			}

			err = s.ChangeMdmOwnerShip(tc.mdmID)
			if err != nil {
				if tc.expected == nil {
					t.Errorf("Changing MDM ownership did not work as expected, \n\tgot: %s \n\twant: %v", err, tc.expected)
				} else {
					if err.Error() != tc.expected.Error() {
						t.Errorf("Changing MDM ownership did not work as expected, \n\tgot: %s \n\twant: %s", err, tc.expected)
					}
				}
			}
		})
	}
}

func TestSwitchClusterMode(t *testing.T) {
	type testCase struct {
		mode                string
		addSecondaryMdmList []string
		addTBMdmList        []string
		expected            error
	}

	cases := []testCase{
		{
			"Success - FiveNodes",
			[]string{"1728fe1657674303"},
			[]string{"463a6129033de104"},
			nil,
		},
		{
			"Failure - FiveNodes",
			[]string{"1728fe1657674311"},
			[]string{"463a6129033de112"},
			errors.New("500 Internal Server Error"),
		},
	}

	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/instances/System/action/switchClusterMode" {
			if r.Method == http.MethodPost {
				var param types.SwitchClusterMode
				_ = json.NewDecoder(r.Body).Decode(&param)
				switch param.Mode {
				case "Failure - FiveNodes":
					w.WriteHeader(http.StatusInternalServerError)
					w.Write([]byte(`{"error": "The MDM could not be found"}`))
				default:
					w.WriteHeader(http.StatusOK)
				}
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

			s := System{
				client: client,
			}

			payload := types.SwitchClusterMode{
				Mode:             tc.mode,
				AddSecondaryMdms: tc.addSecondaryMdmList,
				AddTBMdms:        tc.addTBMdmList,
			}

			err = s.SwitchClusterMode(&payload)
			if err != nil {
				if tc.expected == nil {
					t.Errorf("Switching MDM cluster did not work as expected, \n\tgot: %s \n\twant: %v", err, tc.expected)
				} else {
					if err.Error() != tc.expected.Error() {
						t.Errorf("Switching MDM cluster did not work as expected, \n\tgot: %s \n\twant: %s", err, tc.expected)
					}
				}
			}
		})
	}
}

func TestGetMDMClusterDetails(t *testing.T) {
	testCases := map[string]struct {
		expected *types.MdmCluster
		server   *httptest.Server
		err      error
	}{
		"success": {
			expected: &types.MdmCluster{
				ID: "mdm-cluster-id",
			},
			server: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				w.WriteHeader(http.StatusOK)
				data, err := json.Marshal(types.MdmCluster{
					ID: "mdm-cluster-id",
				})
				if err != nil {
					t.Fatal(err)
				}
				w.Write(data)
			})),
			err: nil,
		},
		"error with API call": {
			expected: nil,
			server: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				w.WriteHeader(http.StatusInternalServerError)
			})),
			err: errors.New("EOF"),
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			defer tc.server.Close()

			client, err := NewClientWithArgs(tc.server.URL, "", math.MaxInt64, true, false)
			client.configConnect.Version = "3.6"
			if err != nil {
				t.Fatal(err)
			}

			s := System{
				client: client,
			}

			mdmDetails, err := s.GetMDMClusterDetails()

			assert.Equal(t, tc.expected, mdmDetails)
			assert.Equal(t, tc.err, err)
		})
	}
}

func TestRenameMdm(t *testing.T) {
	type testCase struct {
		id       string
		newName  string
		expected error
	}

	cases := []testCase{
		{
			"0e4f0a2f5978ae02",
			"mdm_renamed",
			nil,
		},
		{
			"FiveNodes",
			"mdm_renamed",
			errors.New("500 Internal Server Error"),
		},
	}

	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/instances/System/action/renameMdm" {
			if r.Method == http.MethodPost {
				var param types.RenameMdm
				_ = json.NewDecoder(r.Body).Decode(&param)
				switch param.ID {
				case "FiveNodes":
					w.WriteHeader(http.StatusInternalServerError)
					w.Write([]byte(`{"error": "An MDM with the same name already exists"}`))
				default:
					w.WriteHeader(http.StatusOK)
				}
			}
		}
	}))
	defer svr.Close()

	for _, tc := range cases {
		tc := tc
		t.Run("", func(_ *testing.T) {
			client, err := NewClientWithArgs(svr.URL, "", math.MaxInt64, true, false)
			client.configConnect.Version = "3.6"
			if err != nil {
				t.Fatal(err)
			}

			s := System{
				client: client,
			}

			payload := types.RenameMdm{
				ID:      tc.id,
				NewName: tc.newName,
			}

			err = s.RenameMdm(&payload)
			if err != nil {
				if tc.expected == nil {
					t.Errorf("Renaming MDM did not work as expected, \n\tgot: %s \n\twant: %v", err, tc.expected)
				} else {
					if err.Error() != tc.expected.Error() {
						t.Errorf("Renaming MDM did not work as expected, \n\tgot: %s \n\twant: %s", err, tc.expected)
					}
				}
			}
		})
	}
}

func TestGetSystems(t *testing.T) {
	// Define test cases
	testCases := []struct {
		name     string
		input    *Client
		expected []*types.System
		server   *httptest.Server
		err      error
	}{
		{
			name:  "Test case 1",
			input: &Client{},
			expected: []*types.System{
				{
					ID:   "system1",
					Name: "System 1",
				},
				{
					ID:   "system2",
					Name: "System 2",
				},
			},
			server: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				w.WriteHeader(http.StatusOK)
				data, err := json.Marshal([]*types.System{
					{
						ID:   "system1",
						Name: "System 1",
					},
					{
						ID:   "system2",
						Name: "System 2",
					},
				})
				if err != nil {
					t.Fatal(err)
				}
				w.Write(data)
			})),
			err: nil,
		},
		{
			name:     "Test case 2",
			input:    &Client{},
			expected: nil,
			server: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				w.WriteHeader(http.StatusInternalServerError)
			})),
			err: errors.New("err: problem getting instances: EOF"),
		},
	}

	// Run test cases
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			defer tc.server.Close()
			client, err := NewClientWithArgs(tc.server.URL, "", math.MaxInt64, true, false)

			// Call the GetSystems function
			result, err := client.GetSystems()

			// Assert the results
			if tc.err != nil {
				assert.EqualError(t, err, tc.err.Error())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expected, result)
			}
		})
	}
}

func TestFindSystem(t *testing.T) {
	// Define test cases
	testCases := map[string]struct {
		instanceID string
		name       string
		href       string
		expected   *types.System
		server     *httptest.Server
		err        error
	}{
		"success find by instanceID": {
			instanceID: "system1",
			name:       "",
			href:       "",
			expected: &types.System{
				ID:   "system1",
				Name: "System 1",
			},
			server: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				switch r.RequestURI {
				case "/api/types/System/instances":
					w.WriteHeader(http.StatusOK)
					data, err := json.Marshal([]*types.System{
						{
							ID:   "system1",
							Name: "System 1",
						},
						{
							ID:   "system2",
							Name: "System 2",
						},
					})
					if err != nil {
						t.Fatal(err)
					}
					w.Write(data)
				default:
					w.WriteHeader(http.StatusInternalServerError)
				}
			})),
			err: nil,
		},
		"success find by name": {
			instanceID: "",
			name:       "System 2",
			href:       "",
			expected: &types.System{
				ID:   "system2",
				Name: "System 2",
			},
			server: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				switch r.RequestURI {
				case "/api/types/System/instances":
					w.WriteHeader(http.StatusOK)
					data, err := json.Marshal([]*types.System{
						{
							ID:   "system1",
							Name: "System 1",
						},
						{
							ID:   "system2",
							Name: "System 2",
						},
					})
					if err != nil {
						t.Fatal(err)
					}
					w.Write(data)
				default:
					w.WriteHeader(http.StatusInternalServerError)
				}
			})),
			err: nil,
		},
		"error unable to find matching system by instanceID or name": {
			instanceID: "",
			name:       "System 3",
			href:       "",
			expected:   &types.System{},
			server: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				switch r.RequestURI {
				case "/api/types/System/instances":
					w.WriteHeader(http.StatusOK)
					data, err := json.Marshal([]*types.System{
						{
							ID:   "system1",
							Name: "System 1",
						},
						{
							ID:   "system2",
							Name: "System 2",
						},
					})
					if err != nil {
						t.Fatal(err)
					}
					w.Write(data)
				default:
					w.WriteHeader(http.StatusInternalServerError)
				}
			})),
			err: fmt.Errorf("err: systemid or systemname not found"),
		},
		"error from API call": {
			instanceID: "",
			name:       "",
			href:       "",
			expected:   &types.System{},
			server: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				switch r.RequestURI {
				case "/api/types/System/instances":
					w.WriteHeader(http.StatusInternalServerError)
				default:
					w.WriteHeader(http.StatusInternalServerError)
				}
			})),
			err: fmt.Errorf("err: problem getting instances: EOF"),
		},
	}

	// Run test cases
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			defer tc.server.Close()
			client, err := NewClientWithArgs(tc.server.URL, "", math.MaxInt64, true, false)

			// Call the GetSystems function
			result, err := client.FindSystem(tc.instanceID, tc.name, tc.href)

			// Assert the results
			if tc.err != nil {
				assert.EqualError(t, err, tc.err.Error())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expected, result.System)
			}
		})
	}
}

func TestSystemGetStatistics(t *testing.T) {
	// Define test cases
	testCases := map[string]struct {
		system   *System
		expected *types.Statistics
		server   *httptest.Server
		err      error
	}{
		"success": {
			system: &System{
				System: &types.System{
					Links: []*types.Link{
						{
							Rel:  "/api/System/relationship/Statistics",
							HREF: "/api/System/relationship/Statistics/system-1",
						},
					},
				},
			},
			expected: &types.Statistics{
				NumOfStoragePools: 3,
			},
			server: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				switch r.RequestURI {
				case "/api/System/relationship/Statistics/system-1":
					data, err := json.Marshal(types.Statistics{
						NumOfStoragePools: 3,
					})
					if err != nil {
						t.Fatal(err)
					}
					w.Write(data)
				default:
					w.WriteHeader(http.StatusInternalServerError)
				}
			})),
			err: nil,
		},
		"error due to no links in system": {
			system: &System{
				System: &types.System{
					Links: []*types.Link{
						{},
					},
				},
			},
			expected: nil,
			server: httptest.NewServer(http.HandlerFunc(func(_ http.ResponseWriter, _ *http.Request) {
			})),
			err: fmt.Errorf("Error: problem finding link"),
		},
		"error with API call": {
			system: &System{
				System: &types.System{
					Links: []*types.Link{
						{
							Rel:  "/api/System/relationship/Statistics",
							HREF: "/api/System/relationship/Statistics/system-1",
						},
					},
				},
			},
			expected: nil,
			server: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				w.WriteHeader(http.StatusInternalServerError)
			})),
			err: fmt.Errorf("EOF"),
		},
	}

	// Run test cases
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			defer tc.server.Close()
			client, err := NewClientWithArgs(tc.server.URL, "", math.MaxInt64, true, false)

			tc.system.client = client

			// Call the GetStatistics function
			result, err := tc.system.GetStatistics()

			// Assert the results
			if tc.err != nil {
				assert.EqualError(t, err, tc.err.Error())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expected, result)
			}
		})
	}
}

func TestCreateSnapshotConsistencyGroup(t *testing.T) {
	// Define test cases
	testCases := map[string]struct {
		system   *System
		expected *types.SnapshotVolumesResp
		server   *httptest.Server
		err      error
	}{
		"success": {
			system: &System{
				System: &types.System{
					Links: []*types.Link{
						{
							Rel:  "self",
							HREF: "/api/System/instances/system-1",
						},
					},
				},
			},
			expected: &types.SnapshotVolumesResp{
				VolumeIDList:    []string{"volume-1", "volume-2"},
				SnapshotGroupID: "group-1",
			},
			server: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				switch r.RequestURI {
				case "/api/System/instances/system-1/action/snapshotVolumes":
					data, err := json.Marshal(&types.SnapshotVolumesResp{
						VolumeIDList:    []string{"volume-1", "volume-2"},
						SnapshotGroupID: "group-1",
					})
					if err != nil {
						t.Fatal(err)
					}
					w.Write(data)
				default:
					w.WriteHeader(http.StatusInternalServerError)
				}
			})),
			err: nil,
		},
		"error due to no links in system": {
			system: &System{
				System: &types.System{
					Links: []*types.Link{
						{},
					},
				},
			},
			expected: nil,
			server: httptest.NewServer(http.HandlerFunc(func(_ http.ResponseWriter, _ *http.Request) {
			})),
			err: fmt.Errorf("Error: problem finding link"),
		},
		"error with API call": {
			system: &System{
				System: &types.System{
					Links: []*types.Link{
						{
							Rel:  "self",
							HREF: "/api/System/instances/system-1",
						},
					},
				},
			},
			expected: nil,
			server: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				w.WriteHeader(http.StatusInternalServerError)
			})),
			err: fmt.Errorf("EOF"),
		},
	}

	// Run test cases
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			defer tc.server.Close()
			client, err := NewClientWithArgs(tc.server.URL, "", math.MaxInt64, true, false)

			tc.system.client = client

			// Call the CreateSnapshotConsistencyGroup function
			result, err := tc.system.CreateSnapshotConsistencyGroup(&types.SnapshotVolumesParam{})

			// Assert the results
			if tc.err != nil {
				assert.EqualError(t, err, tc.err.Error())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expected, result)
			}
		})
	}
}
