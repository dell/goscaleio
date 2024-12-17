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
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	types "github.com/dell/goscaleio/types/v1"
	"github.com/google/uuid"
)

func TestUploadCompliance(t *testing.T) {
	type testCase struct {
		server      *httptest.Server
		expectedErr error
	}

	cases := map[string]testCase{
		"success": {
			server: httptest.NewServer(http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
				switch req.RequestURI {
				case "/api/version":
					resp.WriteHeader(http.StatusOK)
					fmt.Fprintln(resp, "4.0")
				case "/rest/auth/login":
					resp.WriteHeader(http.StatusOK)
					fmt.Fprintln(resp, `{"access_token":"mock_access_token"}`)
				case "/Api/V1/FirmwareRepository":
					resp.WriteHeader(http.StatusCreated)
					content, err := json.Marshal(types.UploadComplianceTopologyDetails{
						ID: uuid.NewString(),
					})
					if err != nil {
						t.Fatal(err)
					}

					resp.Write(content)
				default:
					resp.WriteHeader(http.StatusBadRequest)
					resp.Write([]byte(`{"message":"no route handled","httpStatusCode":400,"errorCode":0}`))
				}
			})),
			expectedErr: nil,
		},
		"error: bad request": {
			server: httptest.NewServer(http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
				switch req.RequestURI {
				case "/api/version":
					resp.WriteHeader(http.StatusOK)
					fmt.Fprintln(resp, "4.0")
				case "/rest/auth/login":
					resp.WriteHeader(http.StatusOK)
					fmt.Fprintln(resp, `{"access_token":"mock_access_token"}`)
				case "/Api/V1/FirmwareRepository":
					resp.WriteHeader(http.StatusBadRequest)
				default:
					resp.WriteHeader(http.StatusBadRequest)
					resp.Write([]byte(`{"message":"no route handled","httpStatusCode":400,"errorCode":0}`))
				}
			})),
			expectedErr: fmt.Errorf("Error while uploading Compliance File"),
		},
		"error: empty response": {
			server: httptest.NewServer(http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
				switch req.RequestURI {
				case "/api/version":
					resp.WriteHeader(http.StatusOK)
					fmt.Fprintln(resp, "4.0")
				case "/rest/auth/login":
					resp.WriteHeader(http.StatusOK)
					fmt.Fprintln(resp, `{"access_token":"mock_access_token"}`)
				case "/Api/V1/FirmwareRepository":
					resp.WriteHeader(http.StatusCreated)
				default:
					resp.WriteHeader(http.StatusBadRequest)
					resp.Write([]byte(`{"message":"no route handled","httpStatusCode":400,"errorCode":0}`))
				}
			})),
			expectedErr: fmt.Errorf("Error while uploading Compliance File"),
		},
		"error: unable to unmarshal": {
			server: httptest.NewServer(http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
				switch req.RequestURI {
				case "/api/version":
					resp.WriteHeader(http.StatusOK)
					fmt.Fprintln(resp, "4.0")
				case "/rest/auth/login":
					resp.WriteHeader(http.StatusOK)
					fmt.Fprintln(resp, `{"access_token":"mock_access_token"}`)
				case "/Api/V1/FirmwareRepository":
					resp.WriteHeader(http.StatusCreated)
					resp.Write([]byte(`{abc}`))
				default:
					resp.WriteHeader(http.StatusBadRequest)
					resp.Write([]byte(`{"message":"no route handled","httpStatusCode":400,"errorCode":0}`))
				}
			})),
			expectedErr: fmt.Errorf("Error getting upload compliance details: invalid character 'a' looking for beginning of object key string"),
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			defer tc.server.Close()

			gc, err := NewGateway(tc.server.URL, "test_username", "test_password", false, false)
			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}

			uploadParams := types.UploadComplianceParam{
				Username: "user",
				Password: "password",
			}

			_, err = gc.UploadCompliance(&uploadParams)
			if err != nil {
				if tc.expectedErr.Error() != err.Error() {
					t.Fatal(err)
				}
			}
		})
	}
}

func TestGetUploadComplianceDetails(t *testing.T) {
	uploadID := uuid.NewString()
	type testCase struct {
		server      *httptest.Server
		expectedErr error
	}

	cases := map[string]testCase{
		"success": {
			server: httptest.NewServer(http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
				switch req.RequestURI {
				case "/api/version":
					resp.WriteHeader(http.StatusOK)
					fmt.Fprintln(resp, "4.0")
				case "/rest/auth/login":
					resp.WriteHeader(http.StatusOK)
					fmt.Fprintln(resp, `{"access_token":"mock_access_token"}`)
				case fmt.Sprintf("/Api/V1/FirmwareRepository/%s", uploadID):
					resp.WriteHeader(http.StatusOK)
					content, err := json.Marshal(types.UploadComplianceTopologyDetails{
						ID:             uuid.NewString(),
						Name:           "myUploadCompliance",
						DefaultCatalog: true,
					})
					if err != nil {
						t.Fatal(err)
					}

					resp.Write(content)
				default:
					resp.WriteHeader(http.StatusBadRequest)
					resp.Write([]byte(`{"message":"no route handled","httpStatusCode":400,"errorCode":0}`))
				}
			})),
			expectedErr: nil,
		},
		"error: bad request": {
			server: httptest.NewServer(http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
				switch req.RequestURI {
				case "/api/version":
					resp.WriteHeader(http.StatusOK)
					fmt.Fprintln(resp, "4.0")
				case "/rest/auth/login":
					resp.WriteHeader(http.StatusOK)
					fmt.Fprintln(resp, `{"access_token":"mock_access_token"}`)
				case fmt.Sprintf("/Api/V1/FirmwareRepository/%s", uploadID):
					resp.WriteHeader(http.StatusBadRequest)
				default:
					resp.WriteHeader(http.StatusBadRequest)
					resp.Write([]byte(`{"message":"no route handled","httpStatusCode":400,"errorCode":0}`))
				}
			})),
			expectedErr: fmt.Errorf("Error while getting Compliance details"),
		},
		"error: empty response": {
			server: httptest.NewServer(http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
				switch req.RequestURI {
				case "/api/version":
					resp.WriteHeader(http.StatusOK)
					fmt.Fprintln(resp, "4.0")
				case "/rest/auth/login":
					resp.WriteHeader(http.StatusOK)
					fmt.Fprintln(resp, `{"access_token":"mock_access_token"}`)
				case fmt.Sprintf("/Api/V1/FirmwareRepository/%s", uploadID):
					resp.WriteHeader(http.StatusOK)
				default:
					resp.WriteHeader(http.StatusBadRequest)
					resp.Write([]byte(`{"message":"no route handled","httpStatusCode":400,"errorCode":0}`))
				}
			})),
			expectedErr: fmt.Errorf("Error Getting Compliance Details"),
		},
		"error: unable to unmarshal": {
			server: httptest.NewServer(http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
				switch req.RequestURI {
				case "/api/version":
					resp.WriteHeader(http.StatusOK)
					fmt.Fprintln(resp, "4.0")
				case "/rest/auth/login":
					resp.WriteHeader(http.StatusOK)
					fmt.Fprintln(resp, `{"access_token":"mock_access_token"}`)
				case fmt.Sprintf("/Api/V1/FirmwareRepository/%s", uploadID):
					resp.WriteHeader(http.StatusOK)
					resp.Write([]byte(`{abc}`))
				default:
					resp.WriteHeader(http.StatusBadRequest)
					resp.Write([]byte(`{"message":"no route handled","httpStatusCode":400,"errorCode":0}`))
				}
			})),
			expectedErr: fmt.Errorf("Error getting upload compliance details: invalid character 'a' looking for beginning of object key string"),
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			defer tc.server.Close()

			gc, err := NewGateway(tc.server.URL, "test_username", "test_password", false, false)
			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}

			_, err = gc.GetUploadComplianceDetails(uploadID, true)
			if err != nil {
				if tc.expectedErr.Error() != err.Error() {
					t.Fatal(err)
				}
			}
		})
	}
}

func TestApproveUnsignedFile(t *testing.T) {
	uploadID := uuid.NewString()
	type testCase struct {
		server      *httptest.Server
		expectedErr error
	}

	cases := map[string]testCase{
		"success": {
			server: httptest.NewServer(http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
				switch req.RequestURI {
				case "/api/version":
					resp.WriteHeader(http.StatusOK)
					fmt.Fprintln(resp, "4.0")
				case "/rest/auth/login":
					resp.WriteHeader(http.StatusOK)
					fmt.Fprintln(resp, `{"access_token":"mock_access_token"}`)
				case fmt.Sprintf("/Api/V1/FirmwareRepository/%s/allowunsignedfile", uploadID):
					resp.WriteHeader(http.StatusNoContent)
				default:
					resp.WriteHeader(http.StatusBadRequest)
					resp.Write([]byte(`{"message":"no route handled","httpStatusCode":400,"errorCode":0}`))
				}
			})),
			expectedErr: nil,
		},
		"error: bad response": {
			server: httptest.NewServer(http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
				switch req.RequestURI {
				case "/api/version":
					resp.WriteHeader(http.StatusOK)
					fmt.Fprintln(resp, "4.0")
				case "/rest/auth/login":
					resp.WriteHeader(http.StatusOK)
					fmt.Fprintln(resp, `{"access_token":"mock_access_token"}`)
				case fmt.Sprintf("/Api/V1/FirmwareRepository/%s/allowunsignedfile", uploadID):
					resp.WriteHeader(http.StatusBadRequest)
				default:
					resp.WriteHeader(http.StatusBadRequest)
					resp.Write([]byte(`{"message":"no route handled","httpStatusCode":400,"errorCode":0}`))
				}
			})),
			expectedErr: fmt.Errorf("Error while approving the unsigned Compliance file"),
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			defer tc.server.Close()

			gc, err := NewGateway(tc.server.URL, "test_username", "test_password", false, false)
			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}

			err = gc.ApproveUnsignedFile(uploadID)
			if err != nil {
				if tc.expectedErr.Error() != err.Error() {
					t.Fatal(err)
				}
			}
		})
	}
}

func TestGetAllUploadComplianceDetails(t *testing.T) {
	type testCase struct {
		server      *httptest.Server
		expectedErr error
	}

	cases := map[string]testCase{
		"success": {
			server: httptest.NewServer(http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
				switch req.RequestURI {
				case "/api/version":
					resp.WriteHeader(http.StatusOK)
					fmt.Fprintln(resp, "4.0")
				case "/rest/auth/login":
					resp.WriteHeader(http.StatusOK)
					fmt.Fprintln(resp, `{"access_token":"mock_access_token"}`)
				case "/Api/V1/FirmwareRepository/":
					resp.WriteHeader(http.StatusOK)
					content, err := json.Marshal([]types.UploadComplianceTopologyDetails{
						{
							ID:   uuid.NewString(),
							Name: "myUploadA",
						},
						{
							ID:   uuid.NewString(),
							Name: "myUploadB",
						},
					})
					if err != nil {
						t.Fatal(err)
					}

					resp.Write(content)
				default:
					resp.WriteHeader(http.StatusBadRequest)
					resp.Write([]byte(`{"message":"no route handled","httpStatusCode":400,"errorCode":0}`))
				}
			})),
			expectedErr: nil,
		},
		"error: bad request": {
			server: httptest.NewServer(http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
				switch req.RequestURI {
				case "/api/version":
					resp.WriteHeader(http.StatusOK)
					fmt.Fprintln(resp, "4.0")
				case "/rest/auth/login":
					resp.WriteHeader(http.StatusOK)
					fmt.Fprintln(resp, `{"access_token":"mock_access_token"}`)
				case "/Api/V1/FirmwareRepository/":
					resp.WriteHeader(http.StatusBadRequest)
				default:
					resp.WriteHeader(http.StatusBadRequest)
					resp.Write([]byte(`{"message":"no route handled","httpStatusCode":400,"errorCode":0}`))
				}
			})),
			expectedErr: fmt.Errorf("Error while getting Compliance details"),
		},
		"error: empty response": {
			server: httptest.NewServer(http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
				switch req.RequestURI {
				case "/api/version":
					resp.WriteHeader(http.StatusOK)
					fmt.Fprintln(resp, "4.0")
				case "/rest/auth/login":
					resp.WriteHeader(http.StatusOK)
					fmt.Fprintln(resp, `{"access_token":"mock_access_token"}`)
				case "/Api/V1/FirmwareRepository/":
					resp.WriteHeader(http.StatusOK)
				default:
					resp.WriteHeader(http.StatusBadRequest)
					resp.Write([]byte(`{"message":"no route handled","httpStatusCode":400,"errorCode":0}`))
				}
			})),
			expectedErr: fmt.Errorf("Error Getting Compliance Details"),
		},
		"error: unable to unmarshal": {
			server: httptest.NewServer(http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
				switch req.RequestURI {
				case "/api/version":
					resp.WriteHeader(http.StatusOK)
					fmt.Fprintln(resp, "4.0")
				case "/rest/auth/login":
					resp.WriteHeader(http.StatusOK)
					fmt.Fprintln(resp, `{"access_token":"mock_access_token"}`)
				case "/Api/V1/FirmwareRepository/":
					resp.WriteHeader(http.StatusOK)
					resp.Write([]byte(`{abc}`))
				default:
					resp.WriteHeader(http.StatusBadRequest)
					resp.Write([]byte(`{"message":"no route handled","httpStatusCode":400,"errorCode":0}`))
				}
			})),
			expectedErr: fmt.Errorf("Error getting upload compliance details: invalid character 'a' looking for beginning of object key string"),
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			defer tc.server.Close()

			gc, err := NewGateway(tc.server.URL, "test_username", "test_password", false, false)
			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}

			_, err = gc.GetAllUploadComplianceDetails()
			if err != nil {
				if tc.expectedErr.Error() != err.Error() {
					t.Fatal(err)
				}
			}
		})
	}
}

func TestGetUploadComplianceDetailsUsingFilter(t *testing.T) {
	searchName := "myUploadA"
	type testCase struct {
		server      *httptest.Server
		expectedErr error
	}

	cases := map[string]testCase{
		"success": {
			server: httptest.NewServer(http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
				switch req.RequestURI {
				case "/api/version":
					resp.WriteHeader(http.StatusOK)
					fmt.Fprintln(resp, "4.0")
				case "/rest/auth/login":
					resp.WriteHeader(http.StatusOK)
					fmt.Fprintln(resp, `{"access_token":"mock_access_token"}`)
				case "/Api/V1/FirmwareRepository/":
					resp.WriteHeader(http.StatusOK)
					content, err := json.Marshal([]types.UploadComplianceTopologyDetails{
						{
							ID:   uuid.NewString(),
							Name: searchName,
						},
						{
							ID:   uuid.NewString(),
							Name: "myUploadB",
						},
					})
					if err != nil {
						t.Fatal(err)
					}

					resp.Write(content)
				default:
					resp.WriteHeader(http.StatusBadRequest)
					resp.Write([]byte(`{"message":"no route handled","httpStatusCode":400,"errorCode":0}`))
				}
			})),
			expectedErr: nil,
		},
		"error: bad request": {
			server: httptest.NewServer(http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
				switch req.RequestURI {
				case "/api/version":
					resp.WriteHeader(http.StatusOK)
					fmt.Fprintln(resp, "4.0")
				case "/rest/auth/login":
					resp.WriteHeader(http.StatusOK)
					fmt.Fprintln(resp, `{"access_token":"mock_access_token"}`)
				case "/Api/V1/FirmwareRepository/":
					resp.WriteHeader(http.StatusBadRequest)
				default:
					resp.WriteHeader(http.StatusBadRequest)
					resp.Write([]byte(`{"message":"no route handled","httpStatusCode":400,"errorCode":0}`))
				}
			})),
			expectedErr: fmt.Errorf("Error while getting Compliance details"),
		},
		"error: unable to find firmware": {
			server: httptest.NewServer(http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
				switch req.RequestURI {
				case "/api/version":
					resp.WriteHeader(http.StatusOK)
					fmt.Fprintln(resp, "4.0")
				case "/rest/auth/login":
					resp.WriteHeader(http.StatusOK)
					fmt.Fprintln(resp, `{"access_token":"mock_access_token"}`)
				case "/Api/V1/FirmwareRepository/":
					resp.WriteHeader(http.StatusOK)
					content, err := json.Marshal([]types.UploadComplianceTopologyDetails{
						{
							ID:   uuid.NewString(),
							Name: "otherFirmwareName",
						},
					})
					if err != nil {
						t.Fatal(err)
					}

					resp.Write(content)
				default:
					resp.WriteHeader(http.StatusBadRequest)
					resp.Write([]byte(`{"message":"no route handled","httpStatusCode":400,"errorCode":0}`))
				}
			})),
			expectedErr: fmt.Errorf("couldn't find the firmware repository"),
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			defer tc.server.Close()

			gc, err := NewGateway(tc.server.URL, "test_username", "test_password", false, false)
			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}

			_, err = gc.GetUploadComplianceDetailsUsingFilter(searchName)
			if err != nil {
				if tc.expectedErr.Error() != err.Error() {
					t.Fatal(err)
				}
			}
		})
	}
}

func TestGetUploadComplianceDetailsUsingID(t *testing.T) {
	searchID := uuid.NewString()
	type testCase struct {
		server      *httptest.Server
		expectedErr error
	}

	cases := map[string]testCase{
		"success": {
			server: httptest.NewServer(http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
				switch req.RequestURI {
				case "/api/version":
					resp.WriteHeader(http.StatusOK)
					fmt.Fprintln(resp, "4.0")
				case "/rest/auth/login":
					resp.WriteHeader(http.StatusOK)
					fmt.Fprintln(resp, `{"access_token":"mock_access_token"}`)
				case fmt.Sprintf("/Api/V1/FirmwareRepository/%s?components=true", searchID):
					resp.WriteHeader(http.StatusOK)
					content, err := json.Marshal(types.FirmwareRepositoryDetails{
						ID:   searchID,
						Name: "myFirmwareRepository",
					})
					if err != nil {
						t.Fatal(err)
					}

					resp.Write(content)
				default:
					resp.WriteHeader(http.StatusBadRequest)
					resp.Write([]byte(`{"message":"no route handled","httpStatusCode":400,"errorCode":0}`))
				}
			})),
			expectedErr: nil,
		},
		"error: bad request": {
			server: httptest.NewServer(http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
				switch req.RequestURI {
				case "/api/version":
					resp.WriteHeader(http.StatusOK)
					fmt.Fprintln(resp, "4.0")
				case "/rest/auth/login":
					resp.WriteHeader(http.StatusOK)
					fmt.Fprintln(resp, `{"access_token":"mock_access_token"}`)
				case fmt.Sprintf("/Api/V1/FirmwareRepository/%s?components=true", searchID):
					resp.WriteHeader(http.StatusBadRequest)
				default:
					resp.WriteHeader(http.StatusBadRequest)
					resp.Write([]byte(`{"message":"no route handled","httpStatusCode":400,"errorCode":0}`))
				}
			})),
			expectedErr: fmt.Errorf("Error while getting Compliance details"),
		},
		"error: empty response": {
			server: httptest.NewServer(http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
				switch req.RequestURI {
				case "/api/version":
					resp.WriteHeader(http.StatusOK)
					fmt.Fprintln(resp, "4.0")
				case "/rest/auth/login":
					resp.WriteHeader(http.StatusOK)
					fmt.Fprintln(resp, `{"access_token":"mock_access_token"}`)
				case fmt.Sprintf("/Api/V1/FirmwareRepository/%s?components=true", searchID):
					resp.WriteHeader(http.StatusOK)
				default:
					resp.WriteHeader(http.StatusBadRequest)
					resp.Write([]byte(`{"message":"no route handled","httpStatusCode":400,"errorCode":0}`))
				}
			})),
			expectedErr: fmt.Errorf("Error Getting Compliance Details"),
		},
		"error: unable to unmarshal": {
			server: httptest.NewServer(http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
				switch req.RequestURI {
				case "/api/version":
					resp.WriteHeader(http.StatusOK)
					fmt.Fprintln(resp, "4.0")
				case "/rest/auth/login":
					resp.WriteHeader(http.StatusOK)
					fmt.Fprintln(resp, `{"access_token":"mock_access_token"}`)
				case fmt.Sprintf("/Api/V1/FirmwareRepository/%s?components=true", searchID):
					resp.WriteHeader(http.StatusOK)
					resp.Write([]byte(`{abc}`))
				default:
					resp.WriteHeader(http.StatusBadRequest)
					resp.Write([]byte(`{"message":"no route handled","httpStatusCode":400,"errorCode":0}`))
				}
			})),
			expectedErr: fmt.Errorf("Error getting upload compliance details: invalid character 'a' looking for beginning of object key string"),
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			defer tc.server.Close()

			gc, err := NewGateway(tc.server.URL, "test_username", "test_password", false, false)
			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}

			_, err = gc.GetUploadComplianceDetailsUsingID(searchID)
			if err != nil {
				if tc.expectedErr.Error() != err.Error() {
					t.Fatal(err)
				}
			}
		})
	}
}

func TestGetFirmwareRepositoryDetailsUsingName(t *testing.T) {
	searchName := "myFirmwareRepository"
	searchID := uuid.NewString()
	type testCase struct {
		server      *httptest.Server
		expectedErr error
	}

	cases := map[string]testCase{
		"success": {
			server: httptest.NewServer(http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
				switch req.RequestURI {
				case "/api/version":
					resp.WriteHeader(http.StatusOK)
					fmt.Fprintln(resp, "4.0")
				case "/rest/auth/login":
					resp.WriteHeader(http.StatusOK)
					fmt.Fprintln(resp, `{"access_token":"mock_access_token"}`)
				case "/Api/V1/FirmwareRepository/":
					resp.WriteHeader(http.StatusOK)
					content, err := json.Marshal([]types.UploadComplianceTopologyDetails{
						{
							ID:   searchID,
							Name: searchName,
						},
						{
							ID:   uuid.NewString(),
							Name: "myUploadB",
						},
					})
					if err != nil {
						t.Fatal(err)
					}

					resp.Write(content)
				case fmt.Sprintf("/Api/V1/FirmwareRepository/%s?components=true", searchID):
					resp.WriteHeader(http.StatusOK)
					content, err := json.Marshal(types.FirmwareRepositoryDetails{
						ID:   searchID,
						Name: "myFirmwareRepository",
					})
					if err != nil {
						t.Fatal(err)
					}

					resp.Write(content)
				default:
					resp.WriteHeader(http.StatusBadRequest)
					resp.Write([]byte(`{"message":"no route handled","httpStatusCode":400,"errorCode":0}`))
				}
			})),
			expectedErr: nil,
		},
		"error: bad request using filter": {
			server: httptest.NewServer(http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
				switch req.RequestURI {
				case "/api/version":
					resp.WriteHeader(http.StatusOK)
					fmt.Fprintln(resp, "4.0")
				case "/rest/auth/login":
					resp.WriteHeader(http.StatusOK)
					fmt.Fprintln(resp, `{"access_token":"mock_access_token"}`)
				case "/Api/V1/FirmwareRepository/":
					resp.WriteHeader(http.StatusBadRequest)
				default:
					resp.WriteHeader(http.StatusBadRequest)
					resp.Write([]byte(`{"message":"no route handled","httpStatusCode":400,"errorCode":0}`))
				}
			})),
			expectedErr: fmt.Errorf("Error while getting Compliance details"),
		},
		"error: bad request using ID": {
			server: httptest.NewServer(http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
				switch req.RequestURI {
				case "/api/version":
					resp.WriteHeader(http.StatusOK)
					fmt.Fprintln(resp, "4.0")
				case "/rest/auth/login":
					resp.WriteHeader(http.StatusOK)
					fmt.Fprintln(resp, `{"access_token":"mock_access_token"}`)
				case "/Api/V1/FirmwareRepository/":
					resp.WriteHeader(http.StatusOK)
					content, err := json.Marshal([]types.UploadComplianceTopologyDetails{
						{
							ID:   searchID,
							Name: searchName,
						},
						{
							ID:   uuid.NewString(),
							Name: "myUploadB",
						},
					})
					if err != nil {
						t.Fatal(err)
					}

					resp.Write(content)
				case fmt.Sprintf("/Api/V1/FirmwareRepository/%s?components=true", searchID):
					resp.WriteHeader(http.StatusBadRequest)
				default:
					resp.WriteHeader(http.StatusBadRequest)
					resp.Write([]byte(`{"message":"no route handled","httpStatusCode":400,"errorCode":0}`))
				}
			})),
			expectedErr: fmt.Errorf("Error while getting Compliance details"),
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			defer tc.server.Close()

			gc, err := NewGateway(tc.server.URL, "test_username", "test_password", false, false)
			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}

			_, err = gc.GetFirmwareRepositoryDetailsUsingName(searchName)
			if err != nil {
				if tc.expectedErr.Error() != err.Error() {
					t.Fatal(err)
				}
			}
		})
	}
}

func TestDeleteFirmwareRepository(t *testing.T) {
	searchID := uuid.NewString()
	type testCase struct {
		server      *httptest.Server
		expectedErr error
	}

	cases := map[string]testCase{
		"success": {
			server: httptest.NewServer(http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
				switch req.RequestURI {
				case "/api/version":
					resp.WriteHeader(http.StatusOK)
					fmt.Fprintln(resp, "4.0")
				case "/rest/auth/login":
					resp.WriteHeader(http.StatusOK)
					fmt.Fprintln(resp, `{"access_token":"mock_access_token"}`)
				case fmt.Sprintf("/Api/V1/FirmwareRepository/%s", searchID):
					resp.WriteHeader(http.StatusNoContent)
				default:
					resp.WriteHeader(http.StatusBadRequest)
					resp.Write([]byte(`{"message":"no route handled","httpStatusCode":400,"errorCode":0}`))
				}
			})),
			expectedErr: nil,
		},
		"error: bad request": {
			server: httptest.NewServer(http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
				switch req.RequestURI {
				case "/api/version":
					resp.WriteHeader(http.StatusOK)
					fmt.Fprintln(resp, "4.0")
				case "/rest/auth/login":
					resp.WriteHeader(http.StatusOK)
					fmt.Fprintln(resp, `{"access_token":"mock_access_token"}`)
				case fmt.Sprintf("/Api/V1/FirmwareRepository/%s", searchID):
					resp.WriteHeader(http.StatusBadRequest)
				default:
					resp.WriteHeader(http.StatusBadRequest)
					resp.Write([]byte(`{"message":"no route handled","httpStatusCode":400,"errorCode":0}`))
				}
			})),
			expectedErr: fmt.Errorf("Error while deleting firmware repository"),
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			defer tc.server.Close()

			gc, err := NewGateway(tc.server.URL, "test_username", "test_password", false, false)
			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}

			err = gc.DeleteFirmwareRepository(searchID)
			if err != nil {
				if tc.expectedErr.Error() != err.Error() {
					t.Fatal(err)
				}
			}
		})
	}
}

func TestTestConnection(t *testing.T) {
	type testCase struct {
		server      *httptest.Server
		expectedErr error
	}

	cases := map[string]testCase{
		"success": {
			server: httptest.NewServer(http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
				switch req.RequestURI {
				case "/api/version":
					resp.WriteHeader(http.StatusOK)
					fmt.Fprintln(resp, "4.0")
				case "/rest/auth/login":
					resp.WriteHeader(http.StatusOK)
					fmt.Fprintln(resp, `{"access_token":"mock_access_token"}`)
				case "/Api/V1/FirmwareRepository/connection":
					resp.WriteHeader(http.StatusNoContent)
				default:
					resp.WriteHeader(http.StatusBadRequest)
					resp.Write([]byte(`{"message":"no route handled","httpStatusCode":400,"errorCode":0}`))
				}
			})),
			expectedErr: nil,
		},
		"error: bad request": {
			server: httptest.NewServer(http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
				switch req.RequestURI {
				case "/api/version":
					resp.WriteHeader(http.StatusOK)
					fmt.Fprintln(resp, "4.0")
				case "/rest/auth/login":
					resp.WriteHeader(http.StatusOK)
					fmt.Fprintln(resp, `{"access_token":"mock_access_token"}`)
				case "/Api/V1/FirmwareRepository/connection":
					resp.WriteHeader(http.StatusBadRequest)
				default:
					resp.WriteHeader(http.StatusBadRequest)
					resp.Write([]byte(`{"message":"no route handled","httpStatusCode":400,"errorCode":0}`))
				}
			})),
			expectedErr: fmt.Errorf("Error while connecting to the source location. Please chack the credentials"),
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			defer tc.server.Close()

			gc, err := NewGateway(tc.server.URL, "test_username", "test_password", false, false)
			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}

			params := types.UploadComplianceParam{
				Username: "username",
				Password: "password",
			}

			err = gc.TestConnection(&params)
			if err != nil {
				if tc.expectedErr.Error() != err.Error() {
					t.Fatal(err)
				}
			}
		})
	}
}
