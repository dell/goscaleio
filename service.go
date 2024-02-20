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
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
)

// DeployService used to deploy service
func (gc *GatewayClient) DeployService(deploymentName, deploymentDesc, serviceTemplateID, firmwareRepositoryID string) (*types.ServiceResponse, error) {
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
			"firmwareRepositoryId":  firmwareRepositoryID, //TODO
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

func (gc *GatewayClient) UpdateService(deploymentID, deploymentName, deploymentDesc string, nodes int) (*types.ServiceResponse, error) {
	defer TimeSpent("UpdateService", time.Now())

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

		var deploymentPayloadJson []byte

		if nodes > 0 {

			var deploymentData map[string]interface{}

			uuid := uuid.New().String()

			parseError := json.Unmarshal([]byte(responseString), &deploymentData)
			if parseError != nil {
				return nil, fmt.Errorf("Error While Parsing Response Data For Deployment: %s", parseError)
			}

			deploymentData["deploymentName"] = deploymentName

			deploymentData["deploymentDescription"] = deploymentDesc

			// Access the "components" field
			serviceTemplate, ok := deploymentData["serviceTemplate"].(map[string]interface{})
			if !ok {
				fmt.Println("Error: serviceTemplate field not found or not a map[string]interface{}")
				return nil, fmt.Errorf("Error While Parsing Response Data For Deployment: %s", ok)
			}

			components, ok := serviceTemplate["components"].([]interface{})
			if !ok {
				fmt.Println("Error: components field not found or not a []interface{}")
				return nil, fmt.Errorf("Error While Parsing Response Data For Deployment: %s", ok)
			}

			// Find the component with type "SERVER"
			var serverComponent map[string]interface{}

			for _, comp := range components {
				comp := comp.(map[string]interface{})
				if comp["type"].(string) == "SERVER" {
					serverComponent = comp
					break
				}
			}

			// Deep copy the component
			clonedComponent := make(map[string]interface{})
			for key, value := range serverComponent {
				clonedComponent[key] = value
			}

			// Modify ID and GUID of the cloned component
			clonedComponent["id"] = uuid
			clonedComponent["name"] = uuid
			clonedComponent["brownfield"] = false

			clonedComponent["identifier"] = nil
			clonedComponent["asmGUID"] = nil
			clonedComponent["puppetCertName"] = nil
			clonedComponent["osPuppetCertName"] = nil
			clonedComponent["managementIpAddress"] = nil

			// Deep copy resources
			resources, ok := clonedComponent["resources"].([]interface{})
			if !ok {
				fmt.Println("Error: resources field not found or not a []interface{}")
				return nil, fmt.Errorf("Error While Parsing Response Data For Deployment: %s", ok)
			}

			clonedResources := make([]interface{}, len(resources))
			for i, res := range resources {
				resCopy := make(map[string]interface{})
				for k, v := range res.(map[string]interface{}) {
					resCopy[k] = v
				}
				clonedResources[i] = resCopy
			}
			clonedComponent["resources"] = clonedResources

			// Exclude list of parameters to skip
			excludeList := map[string]bool{
				"razor_image":         true,
				"scaleio_enabled":     true,
				"scaleio_role":        true,
				"compression_enabled": true,
				"replication_enabled": true,
			}

			// Iterate over resources to modify parameters
			for _, comp := range clonedResources {
				comp := comp.(map[string]interface{})
				if comp["id"].(string) == "asm::server" {

					comp["guid"] = nil

					parameters, ok := comp["parameters"].([]interface{})
					if !ok {
						fmt.Println("Error: components field not found or not a []interface{}")
						return nil, fmt.Errorf("Error While Parsing Response Data For Deployment: %s", ok)
					}

					clonedParams := make([]interface{}, len(parameters))
					for i, param := range parameters {
						paramCopy := make(map[string]interface{})
						for k, v := range param.(map[string]interface{}) {
							paramCopy[k] = v
						}
						clonedParams[i] = paramCopy
					}

					for _, parameter := range clonedParams {
						parameter := parameter.(map[string]interface{})
						if !excludeList[parameter["id"].(string)] {

							if parameter["id"].(string) == "scaleio_mdm_role" {
								parameter["guid"] = nil
								parameter["value"] = "standby_mdm"
							} else {
								parameter["guid"] = nil
								parameter["value"] = nil
							}

						}
					}

					// Update parameters in the component
					comp["parameters"] = clonedParams
				}
			}

			// Append the cloned component back to the components array
			components = append(components, clonedComponent)

			// Update serviceTemplate with modified components
			serviceTemplate["components"] = components

			// Update deploymentData with modified serviceTemplate
			deploymentData["serviceTemplate"] = serviceTemplate

			// Update other fields as needed
			deploymentData["scaleUp"] = true
			deploymentData["retry"] = true

			// Marshal deploymentData to JSON
			deploymentPayloadJson, _ = json.Marshal(deploymentData)

		} else {

			deploymentResponse, jsonParseError := jsonToMap(responseString)
			if jsonParseError != nil {
				return nil, jsonParseError
			}

			deploymentResponse["deploymentName"] = deploymentName

			deploymentResponse["deploymentDescription"] = deploymentDesc

			deploymentPayloadJson, _ = json.Marshal(deploymentResponse)
		}

		fmt.Println("==================================")

		fmt.Println(string(deploymentPayloadJson))

		req, httpError := http.NewRequest("PUT", gc.host+"/Api/V1/Deployment/"+deploymentID, bytes.NewBuffer(deploymentPayloadJson))
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

		//return nil, fmt.Errorf("Error While Parsing Response Data For Deployment:")

	} else {
		var deploymentResponse types.ServiceFailedResponse

		parseError := json.Unmarshal([]byte(responseString), &deploymentResponse)

		if parseError != nil {
			return nil, fmt.Errorf("Error While Parsing Response Data For Deployment: %s", parseError)
		}

		return nil, fmt.Errorf("Error While Parsing Response Data For Deployment: %s", deploymentResponse.Messages[0].DisplayMessage)
	}
}

// Function to check if string is not present in list
func contains(list []string, str string) bool {
	for _, s := range list {
		if s == str {
			return true
		}
	}
	return false
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
