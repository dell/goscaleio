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
	"fmt"
	"net/http"

	types "github.com/dell/goscaleio/types/v1"
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

	if httpResp.StatusCode != 201 {
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
func (gc *GatewayClient) GetUploadComplianceDetails(id string) (*types.UploadComplianceTopologyDetails, error) {
	var getUploadCompResponse types.UploadComplianceTopologyDetails

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

	if httpResp.StatusCode != 200 {
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

	if httpResp.StatusCode != 204 {
		return fmt.Errorf("Error while approving the unsigned Compliance file")
	}

	err3 := storeCookie(httpResp.Header, gc.host)
	if err3 != nil {
		return fmt.Errorf("Error While Storing cookie: %s", err3)
	}

	return nil
}
