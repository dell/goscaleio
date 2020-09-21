package goscaleio

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	types "github.com/dell/goscaleio/types/v1"
)

func Test_GetVolumeStatistics(t *testing.T) {
	type checkFn func(*testing.T, error)
	check := func(fns ...checkFn) []checkFn { return fns }

	hasNoError := func(t *testing.T, err error) {
		if err != nil {
			t.Fatalf("expected no error")
		}
	}

	hasError := func(t *testing.T, err error) {
		if err == nil {
			t.Fatalf("expected error")
		}
	}

	tests := map[string]func(t *testing.T) (*httptest.Server, *types.Volume, []checkFn){
		"success": func(t *testing.T) (*httptest.Server, *types.Volume, []checkFn) {
			VolumeID := "000001111a2222b"
			href := fmt.Sprintf("/api/instances/Volume::%s/relationships/Statistics", VolumeID)
			vol := types.Volume{
				ID: VolumeID,
				Links: []*types.Link{
					{
						Rel:  "/api/Volume/relationship/Statistics",
						HREF: fmt.Sprintf("/api/instances/Volume::%s/relationships/Statistics", VolumeID),
					},
				},
			}
			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method != http.MethodGet {
					t.Fatal(fmt.Errorf("wrong method. Expected %s; but got %s", http.MethodGet, r.Method))
				}
				if r.URL.Path != href {
					t.Fatal(fmt.Errorf("wrong path. Expected %s; but got %s", href, r.URL.Path))
				}
				volumeStats := types.VolumeStatistics{}
				respData, err := json.Marshal(volumeStats)
				if err != nil {
					t.Fatal(err)
				}
				fmt.Fprintln(w, string(respData))
			}))
			return ts, &vol, check(hasNoError)
		},
		"error from getJSONWithRetry": func(t *testing.T) (*httptest.Server, *types.Volume, []checkFn) {
			VolumeID := "000001111a2222b"
			href := fmt.Sprintf("/api/instances/Volume::%s/relationships/Statistics", VolumeID)
			vol := types.Volume{
				ID: VolumeID,
				Links: []*types.Link{
					{
						Rel:  "/api/Volume/relationship/Statistics",
						HREF: fmt.Sprintf("/api/instances/Volume::%s/relationships/Statistics", VolumeID),
					},
				},
			}
			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method != http.MethodGet {
					t.Fatal(fmt.Errorf("wrong method. Expected %s; but got %s", http.MethodGet, r.Method))
				}
				if r.URL.Path != href {
					t.Fatal(fmt.Errorf("wrong path. Expected %s; but got %s", href, r.URL.Path))
				}
				http.NotFound(w, r)
			}))
			return ts, &vol, check(hasError)
		},
		"error from GetLink": func(t *testing.T) (*httptest.Server, *types.Volume, []checkFn) {
			VolumeID := "000001111a2222b"
			vol := types.Volume{
				ID: VolumeID,
				Links: []*types.Link{
					{
						Rel:  "noLink error",
						HREF: fmt.Sprintf("/api/instances/Volume::%s/relationships/Statistics", VolumeID),
					},
				},
			}
			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				volumeStats := types.VolumeStatistics{}
				respData, err := json.Marshal(volumeStats)
				if err != nil {
					t.Fatal(err)
				}
				fmt.Fprintln(w, string(respData))
			}))
			return ts, &vol, check(hasError)
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			ts, vol, checkFns := tc(t)
			defer ts.Close()

			client, err := NewClientWithArgs(ts.URL, "", true, false)
			if err != nil {
				t.Fatal(err)
			}

			volClient := NewVolume(client)
			volClient.Volume = vol
			_, err = volClient.GetVolumeStatistics()
			for _, checkFn := range checkFns {
				checkFn(t, err)
			}
		})
	}
}
