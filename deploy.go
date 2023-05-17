package goscaleio

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	path "path/filepath"
	"strconv"
	"strings"

	types "github.com/dell/goscaleio/types/v1"
)

var (
	errNewClient = errors.New("missing endpoint")
	errSysCerts  = errors.New("Unable to initialize cert pool from system")
)

// GatewayClient is client for Gateway server
type GatewayClient struct {
	http     *http.Client
	host     string
	username string
	password string
}

// NewGateway returns a new gateway client.
func NewGateway(
	host string, username, password string, insecure, useCerts bool) (*GatewayClient, error) {

	if host == "" {
		return nil, errNewClient
	}

	gc := &GatewayClient{
		http:     &http.Client{},
		host:     host,
		username: username,
		password: password,
	}

	if insecure {
		gc.http.Transport = &http.Transport{
			/* #nosec G402 */
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		}
	}

	if !insecure || useCerts {
		pool, err := x509.SystemCertPool()
		if err != nil {
			return nil, errSysCerts
		}

		gc.http.Transport = &http.Transport{
			/* #nosec G402 */
			TLSClientConfig: &tls.Config{
				RootCAs:            pool,
				InsecureSkipVerify: insecure,
			},
		}
	}

	return gc, nil
}

// UploadPackages used for upload packge to gateway server
func (gc *GatewayClient) UploadPackages(filePaths []string) (*types.GatewayResponse, error) {
	var gatewayResponse types.GatewayResponse

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	for _, filePath := range filePaths {

		info, err := os.Stat(filePath)

		if err != nil {
			return &gatewayResponse, err
		}

		if !info.IsDir() && (strings.HasSuffix(filePath, ".tar") || strings.HasSuffix(filePath, ".rpm")) {

			file, filePathError := os.Open(path.Clean(filePath))
			if filePathError != nil {
				return &gatewayResponse, filePathError
			}

			part, fileReaderError := writer.CreateFormFile("files", path.Base(filePath))
			if fileReaderError != nil {
				return &gatewayResponse, fileReaderError
			}
			_, fileContentError := io.Copy(part, file)
			if fileContentError != nil {
				return &gatewayResponse, fileContentError
			}
		} else {
			return &gatewayResponse, fmt.Errorf("invalid file type, please provide valid file type")
		}
	}

	fileWriterError := writer.Close()
	if fileWriterError != nil {
		return &gatewayResponse, fileWriterError
	}

	req, httpError := http.NewRequest("POST", gc.host+"/im/types/installationPackages/instances/actions/uploadPackages", body)
	if httpError != nil {
		return &gatewayResponse, httpError
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("Authorization", "Basic "+base64.StdEncoding.EncodeToString([]byte(gc.username+":"+gc.password)))
	client := gc.http
	response, httpReqError := client.Do(req)

	if httpReqError != nil {
		return &gatewayResponse, httpReqError
	}

	if response.StatusCode != 200 {
		responseString, _ := extractString(response)

		err := json.Unmarshal([]byte(responseString), &gatewayResponse)

		if err != nil {
			return &gatewayResponse, fmt.Errorf("Error For Uploading Package: %s", err)
		}

		return &gatewayResponse, fmt.Errorf("Error For Uploading Package: %s", gatewayResponse.Message)
	}

	gatewayResponse.StatusCode = 200

	return &gatewayResponse, nil
}

// ParseCSV used for upload csv to gateway server and parse it
func (gc *GatewayClient) ParseCSV(filePath string) (*types.GatewayResponse, error) {
	var gatewayResponse types.GatewayResponse

	file, filePathError := os.Open(path.Clean(filePath))
	if filePathError != nil {
		return &gatewayResponse, filePathError
	}

	defer func() error {
		if err := file.Close(); err != nil {
			return err
		}
		return nil
	}()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	part, fileReaderError := writer.CreateFormFile("file", path.Base(filePath))
	if fileReaderError != nil {
		return &gatewayResponse, fileReaderError
	}
	_, fileContentError := io.Copy(part, file)
	if fileContentError != nil {
		return &gatewayResponse, fileContentError
	}
	fileWriterError := writer.Close()
	if fileWriterError != nil {
		return &gatewayResponse, fileWriterError
	}

	req, httpError := http.NewRequest("POST", gc.host+"/im/types/Configuration/instances/actions/parseFromCSV", body)
	if httpError != nil {
		return &gatewayResponse, httpError
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("Authorization", "Basic "+base64.StdEncoding.EncodeToString([]byte(gc.username+":"+gc.password)))
	client := gc.http
	response, httpReqError := client.Do(req)

	if httpReqError != nil {
		return &gatewayResponse, httpReqError
	}

	responseString, _ := extractString(response)

	if response.StatusCode == 200 {

		gatewayResponse.Data = responseString

		gatewayResponse.StatusCode = response.StatusCode

		return &gatewayResponse, nil
	}

	err := json.Unmarshal([]byte(responseString), &gatewayResponse)

	if err != nil {
		return &gatewayResponse, fmt.Errorf("Error While Parsing Response Data For CSV: %s", err)
	}

	return &gatewayResponse, fmt.Errorf("Error For Parse CSV: %s", gatewayResponse.Message)
}

// GetPackageDetails used for get package details
func (gc *GatewayClient) GetPackageDetails() ([]*types.PackageDetails, error) {

	var packageParam []*types.PackageDetails

	req, httpError := http.NewRequest("GET", gc.host+"/im/types/installationPackages/instances?onlyLatest=false&_search=false", nil)
	if httpError != nil {
		return packageParam, httpError
	}
	req.Header.Set("Authorization", "Basic "+base64.StdEncoding.EncodeToString([]byte(gc.username+":"+gc.password)))
	req.Header.Set("Content-Type", "application/json")

	client := gc.http
	httpRes, httpReqError := client.Do(req)
	if httpReqError != nil {
		return packageParam, httpReqError
	}

	responseString, _ := extractString(httpRes)

	if httpRes.StatusCode == 200 {

		err := json.Unmarshal([]byte(responseString), &packageParam)

		if err != nil {
			return packageParam, fmt.Errorf("Error For Get Package Details: %s", err)
		}

		return packageParam, nil
	}

	return packageParam, nil
}

// ValidateMDMDetails used for validate mdm details
func (gc *GatewayClient) ValidateMDMDetails(mdmTopologyParam []byte) (*types.GatewayResponse, error) {
	var gatewayResponse types.GatewayResponse

	req, httpError := http.NewRequest("POST", gc.host+"/im/types/Configuration/instances", bytes.NewBuffer(mdmTopologyParam))
	if httpError != nil {
		return &gatewayResponse, httpError
	}
	req.Header.Set("Authorization", "Basic "+base64.StdEncoding.EncodeToString([]byte(gc.username+":"+gc.password)))
	req.Header.Set("Content-Type", "application/json")

	client := gc.http
	httpRes, httpReqError := client.Do(req)
	if httpReqError != nil {
		return &gatewayResponse, httpReqError
	}

	responseString, _ := extractString(httpRes)

	if httpRes.StatusCode != 200 {

		err := json.Unmarshal([]byte(responseString), &gatewayResponse)

		if err != nil {
			return &gatewayResponse, fmt.Errorf("Error For Validate MDM Details: %s", err)
		}

		return &gatewayResponse, nil
	}

	var mdmTopologyDetails types.MDMTopologyDetails

	err := json.Unmarshal([]byte(responseString), &mdmTopologyDetails)

	if err != nil {
		return &gatewayResponse, fmt.Errorf("Error For Validate MDM Details: %s", err)
	}

	gatewayResponse.StatusCode = 200

	gatewayResponse.Data = strings.Join(mdmTopologyDetails.SdcIps, ",")

	return &gatewayResponse, nil
}

// DeletePackage used for delete packages from gateway server
func (gc *GatewayClient) DeletePackage(packageName string) (*types.GatewayResponse, error) {

	var gatewayResponse types.GatewayResponse

	req, httpError := http.NewRequest("DELETE", gc.host+"/im/types/installationPackages/instances/actions/delete::"+packageName, nil)
	if httpError != nil {
		return &gatewayResponse, httpError
	}
	req.Header.Set("Authorization", "Basic "+base64.StdEncoding.EncodeToString([]byte(gc.username+":"+gc.password)))
	req.Header.Set("Content-Type", "application/json")

	client := gc.http
	httpRes, httpReqError := client.Do(req)
	if httpReqError != nil {
		return &gatewayResponse, httpReqError
	}

	responseString, _ := extractString(httpRes)

	if httpRes.StatusCode != 200 {

		err := json.Unmarshal([]byte(responseString), &gatewayResponse)

		if err != nil {
			return &gatewayResponse, fmt.Errorf("Error For Delete Package: %s", err)
		}

		return &gatewayResponse, nil
	}

	gatewayResponse.StatusCode = 200

	return &gatewayResponse, nil
}

// BeginInstallation used for start installation
func (gc *GatewayClient) BeginInstallation(jsonStr, mdmUsername, mdmPassword, liaPassword string, expansion bool) (*types.GatewayResponse, error) {

	var gatewayResponse types.GatewayResponse

	mapData, jsonParseError := jsonToMap(jsonStr)
	if jsonParseError != nil {
		return &gatewayResponse, jsonParseError
	}

	mapData["mdmPassword"] = mdmPassword
	mapData["mdmUser"] = mdmUsername
	mapData["liaPassword"] = liaPassword

	secureData := map[string]interface{}{
		"allowNonSecureCommunicationWithMdm": true,
		"allowNonSecureCommunicationWithLia": true,
		"disableNonMgmtComponentsAuth":       false,
	}
	mapData["securityConfiguration"] = secureData

	finalJSON, _ := json.Marshal(mapData)

	u, _ := url.Parse(gc.host + "/im/types/Configuration/actions/install")
	q := u.Query()
	q.Set("noUpload", "false")
	q.Set("noInstall", "false")
	q.Set("noConfigure", "false")
	q.Set("noLinuxDevValidation", "false")
	q.Set("globalZeroPadPolicy", "false")
	q.Set("extend", strconv.FormatBool(expansion))
	u.RawQuery = q.Encode()

	req, httpError := http.NewRequest("POST", u.String(), bytes.NewBuffer(finalJSON))
	if httpError != nil {
		return &gatewayResponse, httpError
	}
	req.Header.Set("Authorization", "Basic "+base64.StdEncoding.EncodeToString([]byte(gc.username+":"+gc.password)))
	req.Header.Set("Content-Type", "application/json")

	client := gc.http

	httpRes, httpReqError := client.Do(req)
	if httpReqError != nil {
		return &gatewayResponse, httpReqError
	}

	if httpRes.StatusCode != 202 {

		responseString, _ := extractString(httpRes)

		err := json.Unmarshal([]byte(responseString), &gatewayResponse)

		if err != nil {
			return &gatewayResponse, fmt.Errorf("Error For Begin Installation: %s", err)
		}

		return &gatewayResponse, nil
	}

	gatewayResponse.StatusCode = 200

	return &gatewayResponse, nil
}

// MoveToNextPhase used for move to next phases in installation
func (gc *GatewayClient) MoveToNextPhase() (*types.GatewayResponse, error) {

	var gatewayResponse types.GatewayResponse

	req, httpError := http.NewRequest("POST", gc.host+"/im/types/ProcessPhase/actions/moveToNextPhase", nil)
	if httpError != nil {
		return &gatewayResponse, httpError
	}
	req.Header.Set("Authorization", "Basic "+base64.StdEncoding.EncodeToString([]byte(gc.username+":"+gc.password)))
	req.Header.Set("Content-Type", "application/json")

	client := gc.http

	httpRes, httpReqError := client.Do(req)
	if httpReqError != nil {
		return &gatewayResponse, httpReqError
	}

	responseString, _ := extractString(httpRes)

	if httpRes.StatusCode != 200 {

		err := json.Unmarshal([]byte(responseString), &gatewayResponse)

		if err != nil {
			return &gatewayResponse, fmt.Errorf("Error For Move To Next Phase: %s", err)
		}

		return &gatewayResponse, nil
	}

	gatewayResponse.StatusCode = 200

	return &gatewayResponse, nil
}

// RetryPhase used for re run to failed phases in installation
func (gc *GatewayClient) RetryPhase() (*types.GatewayResponse, error) {

	var gatewayResponse types.GatewayResponse

	req, httpError := http.NewRequest("POST", gc.host+"/im/types/Command/instances/actions/retry/", nil)
	if httpError != nil {
		return &gatewayResponse, httpError
	}
	req.Header.Set("Authorization", "Basic "+base64.StdEncoding.EncodeToString([]byte(gc.username+":"+gc.password)))
	req.Header.Set("Content-Type", "application/json")

	client := gc.http

	httpRes, httpReqError := client.Do(req)
	if httpReqError != nil {
		return &gatewayResponse, httpReqError
	}

	responseString, _ := extractString(httpRes)

	if httpRes.StatusCode != 200 {

		err := json.Unmarshal([]byte(responseString), &gatewayResponse)

		if err != nil {
			return &gatewayResponse, fmt.Errorf("Error For Retry Phase: %s", err)
		}

		return &gatewayResponse, nil
	}

	gatewayResponse.StatusCode = 200

	return &gatewayResponse, nil
}

// AbortOperation used for abort installation operation
func (gc *GatewayClient) AbortOperation() (*types.GatewayResponse, error) {

	var gatewayResponse types.GatewayResponse

	req, httpError := http.NewRequest("POST", gc.host+"/im/types/Command/instances/actions/abort", nil)
	if httpError != nil {
		return &gatewayResponse, httpError
	}
	req.Header.Set("Authorization", "Basic "+base64.StdEncoding.EncodeToString([]byte(gc.username+":"+gc.password)))
	req.Header.Set("Content-Type", "application/json")

	client := gc.http

	httpRes, httpReqError := client.Do(req)
	if httpReqError != nil {
		return &gatewayResponse, httpReqError
	}

	responseString, _ := extractString(httpRes)

	if httpRes.StatusCode != 200 {

		err := json.Unmarshal([]byte(responseString), &gatewayResponse)

		if err != nil {
			return &gatewayResponse, fmt.Errorf("Error For Abort Operation: %s", err)
		}

		return &gatewayResponse, nil
	}

	gatewayResponse.StatusCode = 200

	return &gatewayResponse, nil
}

// ClearQueueCommand used for clear all commands in queue
func (gc *GatewayClient) ClearQueueCommand() (*types.GatewayResponse, error) {

	var gatewayResponse types.GatewayResponse

	req, httpError := http.NewRequest("POST", gc.host+"/im/types/Command/instances/actions/clear", nil)
	if httpError != nil {
		return &gatewayResponse, httpError
	}
	req.Header.Set("Authorization", "Basic "+base64.StdEncoding.EncodeToString([]byte(gc.username+":"+gc.password)))
	req.Header.Set("Content-Type", "application/json")

	client := gc.http

	httpRes, httpReqError := client.Do(req)
	if httpReqError != nil {
		return &gatewayResponse, httpReqError
	}

	responseString, _ := extractString(httpRes)

	if httpRes.StatusCode != 200 {

		err := json.Unmarshal([]byte(responseString), &gatewayResponse)

		if err != nil {
			return &gatewayResponse, fmt.Errorf("Error For Clear Queue Commands: %s", err)
		}

		return &gatewayResponse, nil
	}

	gatewayResponse.StatusCode = 200

	return &gatewayResponse, nil
}

// MoveToIdlePhase used for move gateway installer to idle state
func (gc *GatewayClient) MoveToIdlePhase() (*types.GatewayResponse, error) {

	var gatewayResponse types.GatewayResponse

	req, httpError := http.NewRequest("POST", gc.host+"/im/types/ProcessPhase/actions/moveToIdlePhase", nil)
	if httpError != nil {
		return &gatewayResponse, httpError
	}
	req.Header.Set("Authorization", "Basic "+base64.StdEncoding.EncodeToString([]byte(gc.username+":"+gc.password)))
	req.Header.Set("Content-Type", "application/json")

	client := gc.http

	httpRes, httpReqError := client.Do(req)
	if httpReqError != nil {
		return &gatewayResponse, httpReqError
	}

	responseString, _ := extractString(httpRes)

	if httpRes.StatusCode != 200 {

		err := json.Unmarshal([]byte(responseString), &gatewayResponse)

		if err != nil {
			return &gatewayResponse, fmt.Errorf("Error For Move To Ideal Phase: %s", err)
		}

		return &gatewayResponse, nil
	}

	gatewayResponse.StatusCode = 200

	return &gatewayResponse, nil
}

// GetInQueueCommand used for get in queue commands
func (gc *GatewayClient) GetInQueueCommand() ([]types.MDMQueueCommandDetails, error) {

	var mdmQueueCommandDetails []types.MDMQueueCommandDetails

	req, httpError := http.NewRequest("GET", gc.host+"/im/types/Command/instances", nil)
	if httpError != nil {
		return mdmQueueCommandDetails, httpError
	}
	req.Header.Set("Authorization", "Basic "+base64.StdEncoding.EncodeToString([]byte(gc.username+":"+gc.password)))
	req.Header.Set("Content-Type", "application/json")

	client := gc.http
	httpRes, httpReqError := client.Do(req)
	if httpReqError != nil {
		return mdmQueueCommandDetails, httpReqError
	}

	responseString, _ := extractString(httpRes)

	if httpRes.StatusCode == 200 {

		var queueCommandDetails map[string][]interface{}

		err := json.Unmarshal([]byte(responseString), &queueCommandDetails)

		if err != nil {
			return mdmQueueCommandDetails, fmt.Errorf("Error For Get In Queue Commands: %s", err)
		}

		var commandList []interface{}

		for _, value := range queueCommandDetails {
			commandList = append(commandList, value...)
		}

		mdmCommands, _ := json.Marshal(commandList)

		err = json.Unmarshal([]byte(mdmCommands), &mdmQueueCommandDetails)

		if err != nil {
			return mdmQueueCommandDetails, fmt.Errorf("Error For Get In Queue Commands: %s", err)
		}

		return mdmQueueCommandDetails, nil
	}

	return mdmQueueCommandDetails, nil
}

// CheckForCompletionQueueCommands used for check queue commands completed or not
func (gc *GatewayClient) CheckForCompletionQueueCommands(currentPhase string) (*types.GatewayResponse, error) {
	var gatewayResponse types.GatewayResponse

	mdmQueueCommandDetails, err := gc.GetInQueueCommand()

	if err != nil {
		return &gatewayResponse, err
	}

	checkCompleted := "Completed"

	var errMsg bytes.Buffer

	for _, mdmQueueCommandDetail := range mdmQueueCommandDetails {

		if currentPhase == mdmQueueCommandDetail.AllowedPhase && mdmQueueCommandDetail.CommandState == "pending" {
			checkCompleted = "Running"
			break
		} else if currentPhase == mdmQueueCommandDetail.AllowedPhase && mdmQueueCommandDetail.CommandState == "failed" {
			checkCompleted = "Failed"
			errMsg.WriteString(mdmQueueCommandDetail.TargetEntityIdentifier + ": " + mdmQueueCommandDetail.Message + ", ")
		}
	}

	if len(errMsg.String()) > 0 {
		gatewayResponse.Message = errMsg.String()[:len(errMsg.String())-2]
	}

	gatewayResponse.Data = checkCompleted

	gatewayResponse.StatusCode = 200

	return &gatewayResponse, nil
}

// jsonToMap used for convert json to map
func jsonToMap(jsonStr string) (map[string]interface{}, error) {
	result := make(map[string]interface{})
	err := json.Unmarshal([]byte(jsonStr), &result)
	if err != nil {
		return nil, err
	}
	return result, nil
}
