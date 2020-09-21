package goscaleio

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	types "github.com/dell/goscaleio/types/v1"
	"github.com/stretchr/testify/assert"
)

func Test_FindVolumes(t *testing.T) {
	type checkFn func(*testing.T, []*Volume, error)
	check := func(fns ...checkFn) []checkFn { return fns }

	hasNoError := func(t *testing.T, vols []*Volume, err error) {
		if err != nil {
			t.Fatalf("expected no error")
		}
	}

	checkLength := func(length int) func(t *testing.T, vols []*Volume, err error) {
		return func(t *testing.T, vols []*Volume, err error) {
			assert.Equal(t, length, len(vols))
		}
	}

	hasError := func(t *testing.T, vols []*Volume, err error) {
		if err == nil {
			t.Fatalf("expected error")
		}
	}

	tests := map[string]func(t *testing.T) (*httptest.Server, types.Sdc, []checkFn){
		"success": func(t *testing.T) (*httptest.Server, types.Sdc, []checkFn) {
			sdcID := "000001111a2222b"
			href := fmt.Sprintf("/api/instances/Sdc::%s/relationships/Volume", sdcID)
			sdc := types.Sdc{
				ID: sdcID,
				Links: []*types.Link{
					{
						Rel:  "/api/Sdc/relationship/Volume",
						HREF: href,
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
				vols := []types.Volume{{}, {}, {}}
				respData, err := json.Marshal(vols)
				if err != nil {
					t.Fatal(err)
				}
				fmt.Fprintln(w, string(respData))
			}))
			return ts, sdc, check(hasNoError, checkLength(3))
		},
		"error from GetVolume": func(t *testing.T) (*httptest.Server, types.Sdc, []checkFn) {
			sdcID := "someID"
			href := fmt.Sprintf("/api/instances/Sdc::%s/relationships/Volume", sdcID)
			sdc := types.Sdc{
				ID: sdcID,
				Links: []*types.Link{
					{
						Rel:  "/api/Sdc/relationship/Volume",
						HREF: href,
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
			return ts, sdc, check(hasError)
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			ts, sdc, checkFns := tc(t)
			defer ts.Close()

			client, err := NewClientWithArgs(ts.URL, "", true, false)
			if err != nil {
				t.Fatal(err)
			}

			sdcClient := NewSdc(client, &sdc)
			vols, err := sdcClient.FindVolumes()
			for _, checkFn := range checkFns {
				checkFn(t, vols, err)
			}
		})
	}
}
