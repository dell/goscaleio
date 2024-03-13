package goscaleio

import (
	"errors"
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
			errors.New("perfProfile should get one of the following values: Compact, HighPerformance, but its value is Compact1"),
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
		ips      []string
		role     string
		expected error
	}

	cases := []testCase{
		{
			[]string{"10.xx.xx.xxx"},
			"Manager",
			errors.New("An invalid IP or host-name specified"),
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

			payload := types.StandByMdm{
				IPs:  tc.ips,
				Role: tc.role,
			}
			_, err = s.AddStandByMdm(&payload)
			if err != nil {
				if tc.expected == nil {
					t.Errorf("Adding standby mdm did not work as expected, \n\tgot: %s \n\twant: %v", err, tc.expected)
				} else {
					if err.Error() != tc.expected.Error() {
						t.Errorf("Adding standby mdm did not work as expected, \n\tgot: %s \n\twant: %s", err, tc.expected)
					}
				}
			}
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
			errors.New("The MDM could not be found"),
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

			err = s.RemoveStandByMdm(tc.mdmID)
			if err != nil {
				if tc.expected == nil {
					t.Errorf("Adding standby mdm did not work as expected, \n\tgot: %s \n\twant: %v", err, tc.expected)
				} else {
					if err.Error() != tc.expected.Error() {
						t.Errorf("Adding standby mdm did not work as expected, \n\tgot: %s \n\twant: %s", err, tc.expected)
					}
				}
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
			errors.New("The MDM could not be found"),
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
			"FiveNodes",
			[]string{"1728fe1657674303"},
			[]string{"463a6129033de104"},
			nil,
		},
		{
			"FiveNodes",
			[]string{"1728fe1657674311"},
			[]string{"463a6129033de112"},
			errors.New("The MDM could not be found"),
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
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	}))
	defer svr.Close()

	client, err := NewClientWithArgs(svr.URL, "", math.MaxInt64, true, false)
	client.configConnect.Version = "3.6"
	if err != nil {
		t.Fatal(err)
	}

	s := System{
		client: client,
	}

	mdmDetails, err1 := s.GetMDMClusterDetails()
	assert.NotNil(t, mdmDetails)
	assert.Nil(t, err1)
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
			errors.New("An MDM with the same name already exists"),
		},
	}

	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNoContent)
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
