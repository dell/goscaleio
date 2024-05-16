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
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"

	types "github.com/dell/goscaleio/types/v1"
	log "github.com/sirupsen/logrus"
)

// UploadCompliance function is used for uploading the compliance file.
func (gc *GatewayClient) UploadCompliance(uploadComplianceParam *types.UploadComplianceParam) (*types.UploadComplianceTopologyDetails, error) {
	var uploadResponse types.UploadComplianceTopologyDetails
	jsonData, err := json.Marshal(uploadComplianceParam)
	if err != nil {
		return &uploadResponse, err
	}

	req, httpError := http.NewRequest("POST", gc.host+"/Api/V1/FirmwareRepository", bytes.NewBuffer(jsonData))
	if httpError != nil {
		return &uploadResponse, httpError
	}

	req.Header.Set("Authorization", "Bearer "+gc.token)
	setCookieError := setCookie(req.Header, gc.host)
	if setCookieError != nil {
		return nil, fmt.Errorf("Error While Handling Cookie: %s", setCookieError)
	}
	req.Header.Set("Content-Type", "application/json")

	client := gc.http
	httpResp, httpRespError := client.Do(req)
	if httpRespError != nil {
		return &uploadResponse, httpRespError
	}

	responseString, err := extractString(httpResp)
	if err != nil {
		return &uploadResponse, fmt.Errorf("Error Extracting Response: %s", err)
	}

	if httpResp.StatusCode != http.StatusCreated {
		return &uploadResponse, fmt.Errorf("Error while uploading Compliance File")
	}

	if responseString == "" {
		return &uploadResponse, fmt.Errorf("Error while uploading Compliance File")
	}

	err = storeCookie(httpResp.Header, gc.host)
	if err != nil {
		return &uploadResponse, fmt.Errorf("Error While Storing cookie: %s", err)
	}

	err = json.Unmarshal([]byte(responseString), &uploadResponse)
	if err != nil {
		return &uploadResponse, fmt.Errorf("Error getting upload compliance details: %s", err)
	}

	return &uploadResponse, nil
}

// GetUploadComplianceDetails function is used for getting the details of the compliance upload
func (gc *GatewayClient) GetUploadComplianceDetails(id string, newToken bool) (*types.UploadComplianceTopologyDetails, error) {
	var getUploadCompResponse types.UploadComplianceTopologyDetails

	if newToken {
		bodyData := map[string]interface{}{
			"username": gc.username,
			"password": gc.password,
		}

		body, _ := json.Marshal(bodyData)

		req, err := http.NewRequest("POST", gc.host+"/rest/auth/login", bytes.NewBuffer(body))
		if err != nil {
			return nil, err
		}

		req.Header.Add("Content-Type", "application/json")

		resp, err := gc.http.Do(req)
		if err != nil {
			return nil, err
		}

		defer func() {
			if err := resp.Body.Close(); err != nil {
				doLog(log.WithError(err).Error, "")
			}
		}()

		// parse the response
		switch {
		case resp == nil:
			return nil, errNilReponse
		case !(resp.StatusCode >= 200 && resp.StatusCode <= 299):
			return nil, gc.api.ParseJSONError(resp)
		}

		bs, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}

		responseBody := string(bs)
		result := make(map[string]interface{})
		jsonErr := json.Unmarshal([]byte(responseBody), &result)
		if err != nil {
			return nil, fmt.Errorf("Error For Uploading Package: %s", jsonErr)
		}

		token := result["access_token"].(string)
		gc.token = token
	}

	req, httpError := http.NewRequest("GET", gc.host+"/Api/V1/FirmwareRepository/"+id, nil)
	if httpError != nil {
		return &getUploadCompResponse, httpError
	}

	req.Header.Set("Authorization", "Bearer "+gc.token)
	setCookieError := setCookie(req.Header, gc.host)
	if setCookieError != nil {
		return nil, fmt.Errorf("Error While Handling Cookie: %s", setCookieError)
	}
	req.Header.Set("Content-Type", "application/json")

	client := gc.http
	httpResp, httpRespError := client.Do(req)
	if httpRespError != nil {
		return &getUploadCompResponse, httpRespError
	}

	responseString, err := extractString(httpResp)
	if err != nil {
		return &getUploadCompResponse, fmt.Errorf("Error Extracting Response: %s", err)
	}

	if httpResp.StatusCode != http.StatusOK {
		return &getUploadCompResponse, fmt.Errorf("Error while getting Compliance details")
	}

	if responseString == "" {
		return &getUploadCompResponse, fmt.Errorf("Error Getting Compliance Details")
	}

	err3 := storeCookie(httpResp.Header, gc.host)
	if err3 != nil {
		return &getUploadCompResponse, fmt.Errorf("Error While Storing cookie: %s", err3)
	}

	err = json.Unmarshal([]byte(responseString), &getUploadCompResponse)
	if err != nil {
		return &getUploadCompResponse, fmt.Errorf("Error getting upload compliance details: %s", err)
	}

	return &getUploadCompResponse, nil
}

// ApproveUnsignedFile is used for approving the unsigned file to upload
func (gc *GatewayClient) ApproveUnsignedFile(id string) error {
	jsonData := []byte(`{}`)

	req, httpError := http.NewRequest("PUT", gc.host+"/Api/V1/FirmwareRepository/"+id+"/allowunsignedfile", bytes.NewBuffer(jsonData))
	if httpError != nil {
		return httpError
	}

	req.Header.Set("Authorization", "Bearer "+gc.token)
	setCookieError := setCookie(req.Header, gc.host)
	if setCookieError != nil {
		return fmt.Errorf("Error While Handling Cookie: %s", setCookieError)
	}
	req.Header.Set("Content-Type", "application/json")

	client := gc.http
	httpResp, httpRespError := client.Do(req)
	if httpRespError != nil {
		return httpRespError
	}

	if httpResp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("Error while approving the unsigned Compliance file")
	}

	err3 := storeCookie(httpResp.Header, gc.host)
	if err3 != nil {
		return fmt.Errorf("Error While Storing cookie: %s", err3)
	}

	return nil
}

// GetAllUploadComplianceDetails returns all the firmware repository
func (gc *GatewayClient) GetAllUploadComplianceDetails() (*[]types.UploadComplianceTopologyDetails, error) {
	var getUploadCompResponse []types.UploadComplianceTopologyDetails
	req, httpError := http.NewRequest("GET", gc.host+"/Api/V1/FirmwareRepository/", nil)
	if httpError != nil {
		return &getUploadCompResponse, httpError
	}

	req.Header.Set("Authorization", "Bearer "+gc.token)
	setCookieError := setCookie(req.Header, gc.host)
	if setCookieError != nil {
		return nil, fmt.Errorf("Error While Handling Cookie: %s", setCookieError)
	}
	req.Header.Set("Content-Type", "application/json")

	client := gc.http
	httpResp, httpRespError := client.Do(req)
	if httpRespError != nil {
		return &getUploadCompResponse, httpRespError
	}

	responseString, err := extractString(httpResp)
	if err != nil {
		return &getUploadCompResponse, fmt.Errorf("Error Extracting Response: %s", err)
	}

	if httpResp.StatusCode != http.StatusOK {
		return &getUploadCompResponse, fmt.Errorf("Error while getting Compliance details")
	}

	if responseString == "" {
		return &getUploadCompResponse, fmt.Errorf("Error Getting Compliance Details")
	}

	err3 := storeCookie(httpResp.Header, gc.host)
	if err3 != nil {
		return &getUploadCompResponse, fmt.Errorf("Error While Storing cookie: %s", err3)
	}

	err = json.Unmarshal([]byte(responseString), &getUploadCompResponse)
	if err != nil {
		return &getUploadCompResponse, fmt.Errorf("Error getting upload compliance details: %s", err)
	}

	return &getUploadCompResponse, nil
}

// GetUploadComplianceDetailsUsingFilter filters the firmware repository based on name
func (gc *GatewayClient) GetUploadComplianceDetailsUsingFilter(name string) (*types.UploadComplianceTopologyDetails, error) {
	frDetails, err := gc.GetAllUploadComplianceDetails()
	if err != nil {
		return nil, err
	}

	for _, fr := range *frDetails {
		if fr.Name == name {
			return &fr, nil
		}
	}
	return nil, errors.New("couldn't find the firmware repository")
}

// GetUploadComplianceDetailsUsingID returns all the details of the firmware repository using ID
func (gc *GatewayClient) GetUploadComplianceDetailsUsingID(id string) (*types.FirmwareRepositoryDetails, error) {
	var frResponse types.FirmwareRepositoryDetails

	u, err := url.Parse(gc.host + "/Api/V1/FirmwareRepository/" + id)
	if err != nil {
		return &frResponse, err
	}
	q := u.Query()
	q.Set("components", "true")
	u.RawQuery = q.Encode()

	req, httpError := http.NewRequest("GET", u.String(), nil)
	if httpError != nil {
		return &frResponse, httpError
	}

	req.Header.Set("Authorization", "Bearer "+gc.token)
	setCookieError := setCookie(req.Header, gc.host)
	if setCookieError != nil {
		return nil, fmt.Errorf("Error While Handling Cookie: %s", setCookieError)
	}
	req.Header.Set("Content-Type", "application/json")

	client := gc.http
	httpResp, httpRespError := client.Do(req)
	if httpRespError != nil {
		return &frResponse, httpRespError
	}

	responseString, err := extractString(httpResp)
	if err != nil {
		return &frResponse, fmt.Errorf("Error Extracting Response: %s", err)
	}

	if httpResp.StatusCode != http.StatusOK {
		return &frResponse, fmt.Errorf("Error while getting Compliance details")
	}

	if responseString == "" {
		return &frResponse, fmt.Errorf("Error Getting Compliance Details")
	}

	err3 := storeCookie(httpResp.Header, gc.host)
	if err3 != nil {
		return &frResponse, fmt.Errorf("Error While Storing cookie: %s", err3)
	}

	err = json.Unmarshal([]byte(responseString), &frResponse)
	if err != nil {
		return &frResponse, fmt.Errorf("Error getting upload compliance details: %s", err)
	}

	return &frResponse, nil
}

// GetFirmwareRepositoryDetailsUsingName returns all the details of the firmware repository using name
func (gc *GatewayClient) GetFirmwareRepositoryDetailsUsingName(name string) (*types.FirmwareRepositoryDetails, error) {
	var fr *types.UploadComplianceTopologyDetails
	var frDetails *types.FirmwareRepositoryDetails
	var err error
	fr, err = gc.GetUploadComplianceDetailsUsingFilter(name)
	if err != nil {
		return frDetails, err
	}
	frDetails, err = gc.GetUploadComplianceDetailsUsingID(fr.ID)
	if err != nil {
		return frDetails, err
	}
	return frDetails, err
}
