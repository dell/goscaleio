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
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	types "github.com/dell/goscaleio/types/v1"
)

func TestGetAllDevices(t *testing.T) {
	type testCase struct {
		server      *httptest.Server
		expectedErr error
	}

	cases := map[string]testCase{
		"get all devices success": {
			server: httptest.NewServer(http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
				switch req.RequestURI {
				case "/api/types/Device/instances":
					resp.WriteHeader(http.StatusOK)
					response := []types.Device{
						{ID: "mock-device-id-1"},
						{ID: "mock-device-id-2"},
					}

					content, err := json.Marshal(response)
					if err != nil {
						t.Fatal(err)
					}

					resp.Write(content)
				default:
					resp.WriteHeader(http.StatusNoContent)
				}
			})),
			expectedErr: nil,
		},
		"get one device success": {
			server: httptest.NewServer(http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
				deviceID := "mock-device-id-1"
				switch req.RequestURI {
				case fmt.Sprintf("/api/instances/Sds::%v/relationships/Device", deviceID):
					resp.WriteHeader(http.StatusOK)
					response := types.Device{
						ID: "1",
					}

					content, err := json.Marshal(response)
					if err != nil {
						t.Fatal(err)
					}

					resp.Write(content)
				default:
					resp.WriteHeader(http.StatusNoContent)
				}
			})),
			expectedErr: nil,
		},
		"get device by field success": {
			server: httptest.NewServer(http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
				switch req.RequestURI {
				case "/api/types/Device/instances":
					resp.WriteHeader(http.StatusOK)
					response := []types.Device{
						{ID: "mock-device-id-1"},
						{ID: "mock-device-id-2"},
					}

					content, err := json.Marshal(response)
					if err != nil {
						t.Fatal(err)
					}

					resp.Write(content)
				default:
					resp.WriteHeader(http.StatusNoContent)
				}
			})),
			expectedErr: nil,
		},
		"bad request": {
			server: httptest.NewServer(http.HandlerFunc(func(resp http.ResponseWriter, _ *http.Request) {
				resp.WriteHeader(http.StatusBadRequest)
				resp.Write([]byte(`{"message":"bad request","httpStatusCode":400,"errorCode":0}`))
			})),
			expectedErr: errors.New("bad request"),
		},
	}

	for _, tc := range cases {
		client, err := NewClientWithArgs(tc.server.URL, "3.6", math.MaxInt64, true, false)
		if err != nil {
			t.Fatal(err)
		}

		s := System{
			client: client,
		}

		_, err = s.GetAllDevice()
		if err != nil {
			if tc.expectedErr.Error() != err.Error() {
				t.Fatal(err)
			}
		}

		tc.server.Close()
	}
}

func TestGetDevice(t *testing.T) {
	type testCase struct {
		server      *httptest.Server
		expectedErr error
	}

	cases := map[string]testCase{
		"get one device success": {
			server: httptest.NewServer(http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
				deviceID := "mock-device-id-1"
				switch req.RequestURI {
				case fmt.Sprintf("/api/instances/Sds::%v/relationships/Device", deviceID):
					resp.WriteHeader(http.StatusOK)
					response := types.Device{
						ID: "mock-device-id-1",
					}

					content, err := json.Marshal(response)
					if err != nil {
						t.Fatal(err)
					}

					resp.Write(content)
				default:
					resp.WriteHeader(http.StatusNoContent)
				}
			})),
			expectedErr: nil,
		},
		"bad request": {
			server: httptest.NewServer(http.HandlerFunc(func(resp http.ResponseWriter, _ *http.Request) {
				resp.WriteHeader(http.StatusBadRequest)
				resp.Write([]byte(`{"message":"bad request","httpStatusCode":400,"errorCode":0}`))
			})),
			expectedErr: errors.New("bad request"),
		},
	}

	for _, tc := range cases {
		client, err := NewClientWithArgs(tc.server.URL, "3.6", math.MaxInt64, true, false)
		if err != nil {
			t.Fatal(err)
		}

		s := System{
			client: client,
		}

		_, err = s.GetDevice("1")
		if err != nil {
			if tc.expectedErr.Error() != err.Error() {
				t.Fatal(err)
			}
		}

		tc.server.Close()
	}
}

func TestGetDeviceByField(t *testing.T) {
	type testCase struct {
		server      *httptest.Server
		expectedErr error
	}

	cases := map[string]testCase{
		"get device by field success": {
			server: httptest.NewServer(http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
				switch req.RequestURI {
				case "/api/types/Device/instances":
					resp.WriteHeader(http.StatusOK)
					response := []types.Device{
						{ID: "mock-device-id-1", Name: "mock-device-name-1"},
						{ID: "mock-device-id-2", Name: "mock-device-name-2"},
					}

					content, err := json.Marshal(response)
					if err != nil {
						t.Fatal(err)
					}

					resp.Write(content)
				default:
					resp.WriteHeader(http.StatusNoContent)
				}
			})),
			expectedErr: nil,
		},
		"get device by field failure": {
			server: httptest.NewServer(http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
				switch req.RequestURI {
				case "/api/types/Device/instances":
					resp.WriteHeader(http.StatusOK)
					response := []types.Device{
						{ID: "mock-device-id-1", Name: "mock-device-name-1"},
						{ID: "mock-device-id-2", Name: "mock-device-name-2"},
					}

					content, err := json.Marshal(response)
					if err != nil {
						t.Fatal(err)
					}

					resp.Write(content)
				default:
					resp.WriteHeader(http.StatusNoContent)
				}
			})),
			expectedErr: errors.New("couldn't find device"),
		},
		"bad request": {
			server: httptest.NewServer(http.HandlerFunc(func(resp http.ResponseWriter, _ *http.Request) {
				resp.WriteHeader(http.StatusBadRequest)
				resp.Write([]byte(`{"message":"bad request","httpStatusCode":400,"errorCode":0}`))
			})),
			expectedErr: errors.New("bad request"),
		},
	}

	deviceFields := map[string]map[string]string{
		"get device by field success": {"ID": "mock-device-id-1", "Name": "mock-device-name-1"},
		"get device by field failure": {"Name": "mock-device-name-invalid"},
		"bad request":                 {},
	}

	for id, tc := range cases {
		t.Run(id, func(t *testing.T) {
			client, err := NewClientWithArgs(tc.server.URL, "", math.MaxInt64, true, false)
			client.configConnect.Version = "4.0"
			if err != nil {
				t.Fatal(err)
			}

			s := System{
				client: client,
			}

			temp := deviceFields[id]
			for fieldKey, fieldValue := range temp {
				_, err = s.GetDeviceByField(fieldKey, fieldValue)
				if err != nil {
					if tc.expectedErr.Error() != err.Error() {
						t.Fatal(err)
					}
				}
			}

		})
	}
}

func TestSDSFindDevice(t *testing.T) {
	type testCase struct {
		server      *httptest.Server
		expectedErr error
	}

	cases := map[string]testCase{
		"get device by field success": {
			server: httptest.NewServer(http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
				switch req.RequestURI {
				case fmt.Sprintf("/api/instances/Sds::%v/relationships/Device", "mock-sds-id"):
					resp.WriteHeader(http.StatusOK)
					response := []types.Device{
						{ID: "mock-device-id-1", Name: "mock-device-name-1"},
						{ID: "mock-device-id-2", Name: "mock-device-name-2"},
					}

					content, err := json.Marshal(response)
					if err != nil {
						t.Fatal(err)
					}

					resp.Write(content)
				default:
					resp.WriteHeader(http.StatusNoContent)
				}
			})),
			expectedErr: nil,
		},
		"get device by field failure": {
			server: httptest.NewServer(http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
				switch req.RequestURI {
				case fmt.Sprintf("/api/instances/Sds::%v/relationships/Device", "mock-sds-name-invalid"):
					resp.WriteHeader(http.StatusOK)
					response := []types.Device{
						{ID: "mock-device-id-1", Name: "mock-device-name-1"},
						{ID: "mock-device-id-2", Name: "mock-device-name-2"},
					}

					content, err := json.Marshal(response)
					if err != nil {
						t.Fatal(err)
					}

					resp.Write(content)
				default:
					resp.WriteHeader(http.StatusNoContent)
				}
			})),
			expectedErr: errors.New("couldn't find device"),
		},

		"bad request": {
			server: httptest.NewServer(http.HandlerFunc(func(resp http.ResponseWriter, _ *http.Request) {
				resp.WriteHeader(http.StatusBadRequest)
				resp.Write([]byte(`{"message":"bad request","httpStatusCode":400,"errorCode":0}`))
			})),
			expectedErr: errors.New("bad request"),
		},
	}

	deviceFields := map[string]map[string]string{
		"get device by field success": {"ID": "mock-device-id-1", "Name": "mock-device-name-1"},
		"get device by field error":   {"ID": "mock-device-id-invalid", "Name": "mock-device-name-invalid"},
		"bad request":                 {},
	}

	for id, tc := range cases {
		t.Run(id, func(t *testing.T) {
			client, err := NewClientWithArgs(tc.server.URL, "", math.MaxInt64, true, false)
			client.configConnect.Version = "4.0"
			if err != nil {
				t.Fatal(err)
			}

			sds := NewSds(client)
			sds.Sds.Name = "mock-sds-name"
			sds.Sds.ID = "mock-sds-id"
			for fieldKey, fieldValue := range deviceFields[id] {
				_, err = sds.FindDevice(fieldKey, fieldValue)
				if err != nil {
					if tc.expectedErr.Error() != err.Error() {
						t.Fatal(err)
					}
				}
			}

		})
	}
}

func TestSDSGetDevice(t *testing.T) {
	type testCase struct {
		server      *httptest.Server
		expectedErr error
	}

	cases := map[string]testCase{
		"get device by field success": {
			server: httptest.NewServer(http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
				switch req.RequestURI {
				case fmt.Sprintf("/api/instances/Sds::%v/relationships/Device", "mock-device-id-1"):
					resp.WriteHeader(http.StatusOK)
					response := []types.Device{
						{ID: "mock-device-id-1", Name: "mock-device-name-1"},
						{ID: "mock-device-id-2", Name: "mock-device-name-2"},
					}

					content, err := json.Marshal(response)
					if err != nil {
						t.Fatal(err)
					}

					resp.Write(content)
				default:
					resp.WriteHeader(http.StatusNoContent)
				}
			})),
			expectedErr: nil,
		},
		"bad request": {
			server: httptest.NewServer(http.HandlerFunc(func(resp http.ResponseWriter, _ *http.Request) {
				resp.WriteHeader(http.StatusBadRequest)
				resp.Write([]byte(`{"message":"bad request","httpStatusCode":400,"errorCode":0}`))
			})),
			expectedErr: errors.New("bad request"),
		},
	}

	for id, tc := range cases {
		t.Run(id, func(t *testing.T) {
			client, err := NewClientWithArgs(tc.server.URL, "", math.MaxInt64, true, false)
			client.configConnect.Version = "4.0"
			if err != nil {
				t.Fatal(err)
			}

			sds := NewSds(client)
			sds.Sds.Name = "mock-sds-name"
			sds.Sds.ID = "mock-sds-id"
			_, err = sds.GetDevice()
			if err != nil {
				if tc.expectedErr.Error() != err.Error() {
					t.Fatal(err)
				}
			}
		})
	}
}

func TestStoragePoolGetDevice(t *testing.T) {
	type testCase struct {
		server      *httptest.Server
		expectedErr error
	}

	cases := map[string]testCase{
		"get device by field success": {
			server: httptest.NewServer(http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
				switch req.RequestURI {
				case fmt.Sprintf("/api/instances/StoragePool::%v/relationships/Device", "mock-storage-pool-name"):
					resp.WriteHeader(http.StatusOK)
					response := []types.Device{
						{ID: "mock-device-id-1", Name: "mock-device-name-1"},
						{ID: "mock-device-id-2", Name: "mock-device-name-2"},
					}

					content, err := json.Marshal(response)
					if err != nil {
						t.Fatal(err)
					}

					resp.Write(content)
				default:
					resp.WriteHeader(http.StatusNoContent)
				}
			})),
			expectedErr: nil,
		},
		"bad request": {
			server: httptest.NewServer(http.HandlerFunc(func(resp http.ResponseWriter, _ *http.Request) {
				resp.WriteHeader(http.StatusBadRequest)
				resp.Write([]byte(`{"message":"bad request","httpStatusCode":400,"errorCode":0}`))
			})),
			expectedErr: errors.New("bad request"),
		},
	}

	for id, tc := range cases {
		t.Run(id, func(t *testing.T) {
			client, err := NewClientWithArgs(tc.server.URL, "", math.MaxInt64, true, false)
			client.configConnect.Version = "4.0"
			if err != nil {
				t.Fatal(err)
			}

			pool := NewStoragePool(client)
			pool.StoragePool.Name = "mock-storage-pool-name"

			_, err = pool.GetDevice()
			if err != nil {
				if tc.expectedErr.Error() != err.Error() {
					t.Fatal(err)
				}
			}
		})
	}
}

func TestStoragePoolFindDevice(t *testing.T) {
	type testCase struct {
		server      *httptest.Server
		expectedErr error
	}

	cases := map[string]testCase{
		"get device by field success": {
			server: httptest.NewServer(http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
				switch req.RequestURI {
				case fmt.Sprintf("/api/instances/StoragePool::%v/relationships/Device", "mock-storage-pool-id"):
					resp.WriteHeader(http.StatusOK)
					response := []types.Device{
						{ID: "mock-device-id-1", Name: "mock-device-name-1"},
						{ID: "mock-device-id-2", Name: "mock-device-name-2"},
					}

					content, err := json.Marshal(response)
					if err != nil {
						t.Fatal(err)
					}

					resp.Write(content)
				default:
					resp.WriteHeader(http.StatusNoContent)
				}
			})),
			expectedErr: nil,
		},
		"get device by field failure": {
			server: httptest.NewServer(http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
				switch req.RequestURI {
				case fmt.Sprintf("/api/instances/StoragePool::%v/relationships/Device", "InvalidStoragePoolID"):
					resp.WriteHeader(http.StatusOK)
					response := []types.Device{
						{ID: "mock-device-id-1", Name: "mock-device-name-1"},
						{ID: "mock-device-id-2", Name: "mock-device-name-2"},
					}

					content, err := json.Marshal(response)
					if err != nil {
						t.Fatal(err)
					}

					resp.Write(content)
				default:
					resp.WriteHeader(http.StatusNoContent)
				}
			})),
			expectedErr: errors.New("couldn't find device"),
		},

		"bad request": {
			server: httptest.NewServer(http.HandlerFunc(func(resp http.ResponseWriter, _ *http.Request) {
				resp.WriteHeader(http.StatusBadRequest)
				resp.Write([]byte(`{"message":"bad request","httpStatusCode":400,"errorCode":0}`))
			})),
			expectedErr: errors.New("bad request"),
		},
	}

	deviceFields := map[string]map[string]string{
		"get device by field success": {"ID": "mock-device-id-1", "Name": "mock-device-name-1"},
		"get device by field error":   {"ID": "mock-device-id-invalid", "Name": "mock-device-name-invalid"},
		"bad request":                 {},
	}

	for id, tc := range cases {
		t.Run(id, func(t *testing.T) {
			client, err := NewClientWithArgs(tc.server.URL, "", math.MaxInt64, true, false)
			client.configConnect.Version = "4.0"
			if err != nil {
				t.Fatal(err)
			}

			pool := NewStoragePool(client)
			pool.StoragePool.Name = "mock-storage-pool-name"
			pool.StoragePool.ID = "mock-storage-pool-id"
			for fieldKey, fieldValue := range deviceFields[id] {
				_, err = pool.FindDevice(fieldKey, fieldValue)
				if err != nil {
					if tc.expectedErr.Error() != err.Error() {
						t.Fatal(err)
					}
				}
			}

		})
	}
}

func TestStoragePoolSetDeviceName(t *testing.T) {
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

	tests := map[string]func(t *testing.T) (*httptest.Server, types.System, []checkFn){
		"success": func(_ *testing.T) (*httptest.Server, types.System, []checkFn) {
			systemID := "mock-system-id"
			mockDeviceID := "mock-device-id"
			url := fmt.Sprintf("/api/instances/Device::%v/action/setDeviceName", mockDeviceID)
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method == http.MethodPost && strings.EqualFold(r.URL.Path, url) {
					w.WriteHeader(http.StatusOK)
					return
				}
				http.NotFound(w, r)
			}))
			system := types.System{
				ID: systemID,
			}
			return server, system, check(hasNoError)
		},
		"error response": func(_ *testing.T) (*httptest.Server, types.System, []checkFn) {
			systemID := "mock-system-id"
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				http.Error(w, "Server Error", http.StatusInternalServerError)
			}))
			system := types.System{
				ID: systemID,
			}
			return server, system, check(hasError)
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			server, _, checkFns := tc(t)
			defer server.Close()

			client, _ := NewClientWithArgs(server.URL, "", math.MaxInt64, true, false)
			storagePool := NewStoragePool(client)
			err := storagePool.SetDeviceName("mock-device-id", "mock-device-name")

			for _, checkFn := range checkFns {
				checkFn(t, err)
			}
		})
	}
}

func TestStoragePoolSetMediaType(t *testing.T) {
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

	tests := map[string]func(t *testing.T) (*httptest.Server, types.System, []checkFn){
		"success": func(_ *testing.T) (*httptest.Server, types.System, []checkFn) {
			systemID := "mock-system-id"
			deviceID := "mock-device-id"
			url := fmt.Sprintf("/api/instances/Device::%v/action/setMediaType", deviceID)
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method == http.MethodPost && strings.EqualFold(r.URL.Path, url) {
					w.WriteHeader(http.StatusOK)
					return
				}
				http.NotFound(w, r)
			}))
			system := types.System{
				ID: systemID,
			}
			return server, system, check(hasNoError)
		},
		"error response": func(_ *testing.T) (*httptest.Server, types.System, []checkFn) {
			systemID := "mock-system-id"
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				http.Error(w, "Server Error", http.StatusInternalServerError)
			}))
			system := types.System{
				ID: systemID,
			}
			return server, system, check(hasError)
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			server, _, checkFns := tc(t)
			defer server.Close()

			client, _ := NewClientWithArgs(server.URL, "", math.MaxInt64, true, false)
			storagePool := NewStoragePool(client)
			err := storagePool.SetDeviceMediaType("mock-device-id", "mock-media-type")

			for _, checkFn := range checkFns {
				checkFn(t, err)
			}
		})
	}
}

func TestStoragePoolSetDeviceExternalAccelerationType(t *testing.T) {
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

	tests := map[string]func(t *testing.T) (*httptest.Server, types.System, []checkFn){
		"success": func(_ *testing.T) (*httptest.Server, types.System, []checkFn) {
			systemID := "mock-system-id"
			deviceID := "mock-device-id"
			url := fmt.Sprintf("/api/instances/Device::%v/action/setExternalAccelerationType", deviceID)
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method == http.MethodPost && strings.EqualFold(r.URL.Path, url) {
					w.WriteHeader(http.StatusOK)
					return
				}
				http.NotFound(w, r)
			}))
			system := types.System{
				ID: systemID,
			}
			return server, system, check(hasNoError)
		},
		"error response": func(_ *testing.T) (*httptest.Server, types.System, []checkFn) {
			systemID := "mock-system-id"
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				http.Error(w, "Server Error", http.StatusInternalServerError)
			}))
			system := types.System{
				ID: systemID,
			}
			return server, system, check(hasError)
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			server, _, checkFns := tc(t)
			defer server.Close()

			client, _ := NewClientWithArgs(server.URL, "", math.MaxInt64, true, false)
			storagePool := NewStoragePool(client)
			err := storagePool.SetDeviceExternalAccelerationType("mock-device-id", "mock-acceleration-type")

			for _, checkFn := range checkFns {
				checkFn(t, err)
			}
		})
	}
}

func TestStoragePoolSetDeviceCapacityLimit(t *testing.T) {
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

	tests := map[string]func(t *testing.T) (*httptest.Server, types.System, []checkFn){
		"success": func(_ *testing.T) (*httptest.Server, types.System, []checkFn) {
			systemID := "mock-system-id"
			deviceID := "mock-device-id"
			url := fmt.Sprintf("/api/instances/Device::%v/action/setDeviceCapacityLimit", deviceID)
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method == http.MethodPost && strings.EqualFold(r.URL.Path, url) {
					w.WriteHeader(http.StatusOK)
					return
				}
				http.NotFound(w, r)
			}))
			system := types.System{
				ID: systemID,
			}
			return server, system, check(hasNoError)
		},
		"error response": func(_ *testing.T) (*httptest.Server, types.System, []checkFn) {
			systemID := "mock-system-id"
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				http.Error(w, "Server Error", http.StatusInternalServerError)
			}))
			system := types.System{
				ID: systemID,
			}
			return server, system, check(hasError)
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			server, _, checkFns := tc(t)
			defer server.Close()

			client, _ := NewClientWithArgs(server.URL, "", math.MaxInt64, true, false)
			storagePool := NewStoragePool(client)
			err := storagePool.SetDeviceCapacityLimit("mock-device-id", "100G")

			for _, checkFn := range checkFns {
				checkFn(t, err)
			}
		})
	}
}

func TestStoragePoolUpdateDeviceOriginalPathways(t *testing.T) {
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

	tests := map[string]func(t *testing.T) (*httptest.Server, types.System, []checkFn){
		"success": func(_ *testing.T) (*httptest.Server, types.System, []checkFn) {
			systemID := "mock-system-id"
			deviceID := "mock-device-id"
			url := fmt.Sprintf("/api/instances/Device::%v/action/updateDeviceOriginalPathname", deviceID)
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method == http.MethodPost && strings.EqualFold(r.URL.Path, url) {
					w.WriteHeader(http.StatusOK)
					return
				}
				http.NotFound(w, r)
			}))
			system := types.System{
				ID: systemID,
			}
			return server, system, check(hasNoError)
		},
		"error response": func(_ *testing.T) (*httptest.Server, types.System, []checkFn) {
			systemID := "mock-system-id"
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				http.Error(w, "Server Error", http.StatusInternalServerError)
			}))
			system := types.System{
				ID: systemID,
			}
			return server, system, check(hasError)
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			server, _, checkFns := tc(t)
			defer server.Close()

			client, _ := NewClientWithArgs(server.URL, "", math.MaxInt64, true, false)
			storagePool := NewStoragePool(client)
			err := storagePool.UpdateDeviceOriginalPathways("mock-device-id")

			for _, checkFn := range checkFns {
				checkFn(t, err)
			}
		})
	}
}

func TestStoragePoolRemoveDevice(t *testing.T) {
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

	tests := map[string]func(t *testing.T) (*httptest.Server, types.System, []checkFn){
		"success": func(_ *testing.T) (*httptest.Server, types.System, []checkFn) {
			systemID := "mock-system-id"
			deviceID := "mock-device-id"
			url := fmt.Sprintf("/api/instances/Device::%v/action/removeDevice", deviceID)
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method == http.MethodPost && strings.EqualFold(r.URL.Path, url) {
					w.WriteHeader(http.StatusOK)
					return
				}
				http.NotFound(w, r)
			}))
			system := types.System{
				ID: systemID,
			}
			return server, system, check(hasNoError)
		},
		"error response": func(_ *testing.T) (*httptest.Server, types.System, []checkFn) {
			systemID := "mock-system-id"
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				http.Error(w, "Server Error", http.StatusInternalServerError)
			}))
			system := types.System{
				ID: systemID,
			}
			return server, system, check(hasError)
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			server, _, checkFns := tc(t)
			defer server.Close()

			client, _ := NewClientWithArgs(server.URL, "", math.MaxInt64, true, false)
			storagePool := NewStoragePool(client)
			err := storagePool.RemoveDevice("mock-device-id")

			for _, checkFn := range checkFns {
				checkFn(t, err)
			}
		})
	}
}

func TestStoragePoolAttachDevice(t *testing.T) {
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

	tests := map[string]func(t *testing.T) (*httptest.Server, types.System, []checkFn){
		"success": func(_ *testing.T) (*httptest.Server, types.System, []checkFn) {
			systemID := "mock-system-id"
			url := "/api/types/Device/instances"
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method == http.MethodPost && strings.EqualFold(r.URL.Path, url) {
					w.WriteHeader(http.StatusOK)
					return
				}
				http.NotFound(w, r)
			}))
			system := types.System{
				ID: systemID,
			}
			return server, system, check(hasNoError)
		},
		"error response": func(_ *testing.T) (*httptest.Server, types.System, []checkFn) {
			systemID := "mock-system-id"
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				http.Error(w, "Server Error", http.StatusInternalServerError)
			}))
			system := types.System{
				ID: systemID,
			}
			return server, system, check(hasError)
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			server, _, checkFns := tc(t)
			defer server.Close()

			client, _ := NewClientWithArgs(server.URL, "", math.MaxInt64, true, false)
			storagePool := NewStoragePool(client)
			deviceParam := &types.DeviceParam{
				Name: "mock-device-name",
			}
			_, err := storagePool.AttachDevice(deviceParam)

			for _, checkFn := range checkFns {
				checkFn(t, err)
			}
		})
	}
}
