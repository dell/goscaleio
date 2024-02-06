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
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	types "github.com/dell/goscaleio/types/v1"
	log "github.com/sirupsen/logrus"
)

// DeployService used to deploy service
func (gc *GatewayClient) DeployService(deploymentName, deploymentDesc, serviceTemplateID, firmwareRepositoryId string) (*types.ServiceResponse, error) {
	defer TimeSpent("DeployService", time.Now())

	path := fmt.Sprintf("/Api/V1/ServiceTemplate/%v?forDeployment=true", serviceTemplateID)

	req, httpError := http.NewRequest("GET", gc.host+path, nil)
	if httpError != nil {
		return nil, httpError
	}

	if gc.version == "4.0" {
		req.Header.Set("Authorization", "Bearer "+gc.token)

		setCookie(req.Header, gc.host)

	} else {
		req.Header.Set("Authorization", "Basic "+base64.StdEncoding.EncodeToString([]byte(gc.username+":"+gc.password)))
	}

	req.Header.Set("Content-Type", "application/json")

	client := gc.http
	httpResp, httpRespError := client.Do(req)
	if httpRespError != nil {
		return nil, httpRespError
	}

	responseString, _ := extractString(httpResp)

	if httpResp.StatusCode == 200 {
		var templateData map[string]interface{}

		parseError := json.Unmarshal([]byte(responseString), &templateData)

		if parseError != nil {
			return nil, fmt.Errorf("Error While Parsing Response Data For Template: %s", parseError)
		}

		deploymentPayload := map[string]interface{}{
			"deploymentName":        deploymentName,
			"deploymentDescription": deploymentDesc,
			"serviceTemplate":       templateData,
			"updateServerFirmware":  true,
			"firmwareRepositoryId":  firmwareRepositoryId, //TODO
		}

		deploymentPayloadJson, _ := json.Marshal(deploymentPayload)

		req, httpError := http.NewRequest("POST", gc.host+"/Api/V1/Deployment", bytes.NewBuffer(deploymentPayloadJson))
		if httpError != nil {
			return nil, httpError
		}
		if gc.version == "4.0" {
			req.Header.Set("Authorization", "Bearer "+gc.token)

			setCookie(req.Header, gc.host)
		} else {
			req.Header.Set("Authorization", "Basic "+base64.StdEncoding.EncodeToString([]byte(gc.username+":"+gc.password)))
		}
		req.Header.Set("Content-Type", "application/json")

		client := gc.http
		httpResp, httpRespError := client.Do(req)
		if httpRespError != nil {
			return nil, httpRespError
		}

		responseString, error := extractString(httpResp)
		if error != nil {
			return nil, fmt.Errorf("Error Extracting Response: %s", error)
		}

		if httpResp.StatusCode == 200 {

			var deploymentResponse types.ServiceResponse

			parseError := json.Unmarshal([]byte(responseString), &deploymentResponse)

			deploymentResponse.StatusCode = 200

			if parseError != nil {
				return nil, fmt.Errorf("Error While Parsing Response Data For Deployment: %s", parseError)
			}

			return &deploymentResponse, nil

		} else {
			var deploymentResponse types.ServiceFailedResponse

			parseError := json.Unmarshal([]byte(responseString), &deploymentResponse)

			deploymentResponse.StatusCode = 400

			if parseError != nil {
				return nil, fmt.Errorf("Error While Parsing Response Data For Deployment: %s", parseError)
			}

			return nil, fmt.Errorf("Error While Parsing Response Data For Deployment: %s", deploymentResponse.Messages[0].DisplayMessage)
		}

	}

	return nil, nil
}

func (gc *GatewayClient) GetServiceDetailsByID(deploymentID string, newToken bool) (*types.ServiceResponse, error) {

	defer TimeSpent("GetServiceDetailsByID", time.Now())

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

	path := fmt.Sprintf("/Api/V1/Deployment/%v", deploymentID)

	req, httpError := http.NewRequest("GET", gc.host+path, nil)
	if httpError != nil {
		return nil, httpError
	}

	if gc.version == "4.0" {
		req.Header.Set("Authorization", "Bearer "+gc.token)

		setCookie(req.Header, gc.host)

	} else {
		req.Header.Set("Authorization", "Basic "+base64.StdEncoding.EncodeToString([]byte(gc.username+":"+gc.password)))
	}

	req.Header.Set("Content-Type", "application/json")

	client := gc.http
	httpResp, httpRespError := client.Do(req)
	if httpRespError != nil {
		return nil, httpRespError
	}

	responseString, _ := extractString(httpResp)

	if httpResp.StatusCode == 200 {

		var deploymentResponse types.ServiceResponse

		parseError := json.Unmarshal([]byte(responseString), &deploymentResponse)

		if parseError != nil {
			return nil, fmt.Errorf("Error While Parsing Response Data For Deployment: %s", parseError)
		}

		return &deploymentResponse, nil

	} else {
		var deploymentResponse types.ServiceFailedResponse

		parseError := json.Unmarshal([]byte(responseString), &deploymentResponse)

		if parseError != nil {
			return nil, fmt.Errorf("Error While Parsing Response Data For Deployment: %s", parseError)
		}

		return nil, fmt.Errorf("Error While Parsing Response Data For Deployment: %s", deploymentResponse.Messages[0].DisplayMessage)
	}
}

func (gc *GatewayClient) GetServiceDetailsByFilter(filter, value string) ([]types.ServiceResponse, error) {

	defer TimeSpent("GetServiceDetailsByFilter", time.Now())

	encodedValue := url.QueryEscape(value)

	path := fmt.Sprintf("/Api/V1/Deployment?filter=eq,%v,%v", filter, encodedValue)

	req, httpError := http.NewRequest("GET", gc.host+path, nil)
	if httpError != nil {
		return nil, httpError
	}

	if gc.version == "4.0" {
		req.Header.Set("Authorization", "Bearer "+gc.token)

		setCookie(req.Header, gc.host)

	} else {
		req.Header.Set("Authorization", "Basic "+base64.StdEncoding.EncodeToString([]byte(gc.username+":"+gc.password)))
	}

	req.Header.Set("Content-Type", "application/json")

	client := gc.http
	httpResp, httpRespError := client.Do(req)
	if httpRespError != nil {
		return nil, httpRespError
	}

	responseString, _ := extractString(httpResp)

	if httpResp.StatusCode == 200 {

		var deploymentResponse []types.ServiceResponse

		parseError := json.Unmarshal([]byte(responseString), &deploymentResponse)

		if parseError != nil {
			return nil, fmt.Errorf("Error While Parsing Response Data For Deployment: %s", parseError)
		}

		return deploymentResponse, nil

	} else {
		var deploymentResponse types.ServiceFailedResponse

		parseError := json.Unmarshal([]byte(responseString), &deploymentResponse)

		if parseError != nil {
			return nil, fmt.Errorf("Error While Parsing Response Data For Deployment: %s", parseError)
		}

		return nil, fmt.Errorf("Error While Parsing Response Data For Deployment: %s", deploymentResponse.Messages[0].DisplayMessage)
	}
}

func (gc *GatewayClient) GetAllServiceDetails() ([]types.ServiceResponse, error) {

	defer TimeSpent("DeploGetServiceDetailsByIDyService", time.Now())

	req, httpError := http.NewRequest("GET", gc.host+"/Api/V1/Deployment/", nil)
	if httpError != nil {
		return nil, httpError
	}

	if gc.version == "4.0" {
		req.Header.Set("Authorization", "Bearer "+gc.token)

		setCookie(req.Header, gc.host)

	} else {
		req.Header.Set("Authorization", "Basic "+base64.StdEncoding.EncodeToString([]byte(gc.username+":"+gc.password)))
	}

	req.Header.Set("Content-Type", "application/json")

	client := gc.http
	httpResp, httpRespError := client.Do(req)
	if httpRespError != nil {
		return nil, httpRespError
	}

	responseString, _ := extractString(httpResp)

	if httpResp.StatusCode == 200 {

		var deploymentResponse []types.ServiceResponse

		parseError := json.Unmarshal([]byte(responseString), &deploymentResponse)

		if parseError != nil {
			return nil, fmt.Errorf("Error While Parsing Response Data For Deployment: %s", parseError)
		}

		return deploymentResponse, nil

	} else {
		var deploymentResponse types.ServiceFailedResponse

		parseError := json.Unmarshal([]byte(responseString), &deploymentResponse)

		if parseError != nil {
			return nil, fmt.Errorf("Error While Parsing Response Data For Deployment: %s", parseError)
		}

		return nil, fmt.Errorf("Error While Parsing Response Data For Deployment: %s", deploymentResponse.Messages[0].DisplayMessage)
	}
}

func (gc *GatewayClient) DeleteService(serviceId string) (*types.ServiceResponse, error) {

	var deploymentResponse types.ServiceResponse

	deploymentResponse.StatusCode = 400

	defer TimeSpent("DeleteService", time.Now())

	req, httpError := http.NewRequest("DELETE", gc.host+"/Api/V1/Deployment/"+serviceId+"?serversInInventory=remove&resourceState=managed", nil)
	if httpError != nil {
		return nil, httpError
	}

	if gc.version == "4.0" {
		req.Header.Set("Authorization", "Bearer "+gc.token)

		setCookie(req.Header, gc.host)

	} else {
		req.Header.Set("Authorization", "Basic "+base64.StdEncoding.EncodeToString([]byte(gc.username+":"+gc.password)))
	}

	req.Header.Set("Content-Type", "application/json")

	client := gc.http
	httpResp, httpRespError := client.Do(req)
	if httpRespError != nil {
		return nil, httpRespError
	}

	if httpResp.StatusCode == 204 {

		deploymentResponse.StatusCode = 200

		return &deploymentResponse, nil
	}

	return &deploymentResponse, nil
}
