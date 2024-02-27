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
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	path "path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/dell/goscaleio/api"
	types "github.com/dell/goscaleio/types/v1"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
)

var (
	errNewClient = errors.New("missing endpoint")
	errSysCerts  = errors.New("Unable to initialize cert pool from system")
)

// GatewayClient is client for Gateway server
type GatewayClient struct {
	http     *http.Client
	api      api.Client
	host     string
	username string
	password string
	token    string
	version  string
}

// NewGateway returns a new gateway client.
func NewGateway(host string, username, password string, insecure, useCerts bool) (*GatewayClient, error) {

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
				InsecureSkipVerify: true, // #nosec G402
			},
		}
	}

	if !insecure || useCerts {
		pool, err := x509.SystemCertPool()
		if err != nil {
			return nil, errSysCerts
		}

		gc.http.Transport = &http.Transport{  // #nosec G402
			TLSClientConfig: &tls.Config{
				RootCAs:            pool,
				InsecureSkipVerify: insecure,
			},
		}
	}

	version, err := gc.GetVersion()
	if err != nil {
		return nil, err
	}

	if version == "3.5" {
		gc.version = version
		//No need to create token
	} else {
		bodyData := map[string]interface{}{
			"username": username,
			"password": password,
		}

		body, _ := json.Marshal(bodyData)

		req, err := http.NewRequest("POST", host+"/rest/auth/login", bytes.NewBuffer(body))
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

		version, err = gc.GetVersion()
		if err != nil {
			return nil, err
		}
		gc.version = version
	}

	return gc, nil
}

// GetVersion returns version
func (gc *GatewayClient) GetVersion() (string, error) {

	req, httpError := http.NewRequest("GET", gc.host+"/api/version", nil)
	if httpError != nil {
		return "", httpError
	}

	if gc.token != "" {
		req.Header.Set("Authorization", "Bearer "+gc.token)
	}

	req.Header.Set("Content-Type", "application/json")

	client := gc.http
	resp, httpRespError := client.Do(req)
	if httpRespError != nil {
		return "", httpRespError
	}

	// parse the response
	switch {
	case resp == nil:
		return "", errNilReponse
	case !(resp.StatusCode >= 200 && resp.StatusCode <= 299):
		return "", nil
	}

	version, err := extractString(resp)
	if err != nil {
		return "", err
	}

	versionRX := regexp.MustCompile(`^(\d+?\.\d+?).*$`)
	if m := versionRX.FindStringSubmatch(version); len(m) > 0 {
		return m[1], nil
	}
	return version, nil
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
	if gc.version == "4.0" {
		req.Header.Set("Authorization", "Bearer "+gc.token)

		err := setCookie(req.Header, gc.host)
		if err != nil {
			return nil, fmt.Errorf("Error While Handling Cookie: %s", err)
		}
	} else {
		req.Header.Set("Authorization", "Basic "+base64.StdEncoding.EncodeToString([]byte(gc.username+":"+gc.password)))
	}
	client := gc.http
	response, httpRespError := client.Do(req)

	if httpRespError != nil {
		return &gatewayResponse, httpRespError
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
	if gc.version == "4.0" {
		req.Header.Set("Authorization", "Bearer "+gc.token)

		err := setCookie(req.Header, gc.host)
		if err != nil {
			return nil, fmt.Errorf("Error While Handling Cookie: %s", err)
		}
	} else {
		req.Header.Set("Authorization", "Basic "+base64.StdEncoding.EncodeToString([]byte(gc.username+":"+gc.password)))
	}
	client := gc.http
	response, httpRespError := client.Do(req)

	if httpRespError != nil {
		return &gatewayResponse, httpRespError
	}

	responseString, _ := extractString(response)

	if response.StatusCode == 200 {

		var parseCSVData map[string]interface{}

		err := json.Unmarshal([]byte(responseString), &parseCSVData)

		if err != nil {
			return &gatewayResponse, fmt.Errorf("Error While Parsing Response Data For CSV: %s", err)
		}

		if parseCSVData["masterMdm"] != nil {
			gatewayResponse.Data = responseString

			gatewayResponse.StatusCode = response.StatusCode

			return &gatewayResponse, nil
		}

		gatewayResponse.StatusCode = 500

		return &gatewayResponse, fmt.Errorf("Error For Parse CSV: Unable to detect a Primary MDM in the CSV file. All the details about the Primary MDM are needed for extending your PowerFlex system. The Primary MDM will not be reinstalled")

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

	if gc.version == "4.0" {
		req.Header.Set("Authorization", "Bearer "+gc.token)

		err := setCookie(req.Header, gc.host)
		if err != nil {
			return nil, fmt.Errorf("Error While Handling Cookie: %s", err)
		}

	} else {
		req.Header.Set("Authorization", "Basic "+base64.StdEncoding.EncodeToString([]byte(gc.username+":"+gc.password)))
	}

	req.Header.Set("Content-Type", "application/json")

	client := gc.http
	httpResp, httpRespError := client.Do(req)
	if httpRespError != nil {
		return packageParam, httpRespError
	}

	responseString, error := extractString(httpResp)
	if error != nil {
		return packageParam, fmt.Errorf("Error Extracting Response: %s", error)
	}

	if httpResp.StatusCode == 200 {

		if gc.version == "4.0" {
			err := storeCookie(httpResp.Header, gc.host)
			if err != nil {
				return packageParam, fmt.Errorf("Error While Storing cookie: %s", err)
			}
		}

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
	if gc.version == "4.0" {
		req.Header.Set("Authorization", "Bearer "+gc.token)

		err := setCookie(req.Header, gc.host)
		if err != nil {
			return nil, fmt.Errorf("Error While Handling Cookie: %s", err)
		}
	} else {
		req.Header.Set("Authorization", "Basic "+base64.StdEncoding.EncodeToString([]byte(gc.username+":"+gc.password)))
	}
	req.Header.Set("Content-Type", "application/json")

	client := gc.http
	httpResp, httpRespError := client.Do(req)
	if httpRespError != nil {
		return &gatewayResponse, httpRespError
	}

	responseString, error := extractString(httpResp)
	if error != nil {
		return &gatewayResponse, fmt.Errorf("Error Extracting Response: %s", error)
	}

	if httpResp.StatusCode != 200 {

		err := json.Unmarshal([]byte(responseString), &gatewayResponse)

		if err != nil {
			return &gatewayResponse, fmt.Errorf("Error Validating MDM Details: %s", err)
		}

		return &gatewayResponse, nil
	} else if httpResp.StatusCode == 200 && responseString == "" {
		gatewayResponse.Message = "Wrong Primary MDM IP, Please provide valid Primary MDM IP"

		return &gatewayResponse, fmt.Errorf("Wrong Primary MDM IP, Please provide valid Primary MDM IP")
	}

	if gc.version == "4.0" {
		err := storeCookie(httpResp.Header, gc.host)
		if err != nil {
			return &gatewayResponse, fmt.Errorf("Error While Storing cookie: %s", err)
		}
	}

	var mdmTopologyDetails types.MDMTopologyDetails

	err := json.Unmarshal([]byte(responseString), &mdmTopologyDetails)

	if err != nil {
		return &gatewayResponse, fmt.Errorf("Error Validating MDM Details: %s", err)
	}

	gatewayResponse.StatusCode = 200

	gatewayResponse.Data = strings.Join(mdmTopologyDetails.SdcIps, ",")

	return &gatewayResponse, nil
}

// GetClusterDetails used for get MDM cluster details
func (gc *GatewayClient) GetClusterDetails(mdmTopologyParam []byte, requireJSONOutput bool) (*types.GatewayResponse, error) {
	var gatewayResponse types.GatewayResponse

	req, httpError := http.NewRequest("POST", gc.host+"/im/types/Configuration/instances", bytes.NewBuffer(mdmTopologyParam))
	if httpError != nil {
		return &gatewayResponse, httpError
	}
	if gc.version == "4.0" {
		req.Header.Set("Authorization", "Bearer "+gc.token)

		err := setCookie(req.Header, gc.host)
		if err != nil {
			return nil, fmt.Errorf("Error While Handling Cookie: %s", err)
		}
	} else {
		req.Header.Set("Authorization", "Basic "+base64.StdEncoding.EncodeToString([]byte(gc.username+":"+gc.password)))
	}
	req.Header.Set("Content-Type", "application/json")

	client := gc.http
	httpResp, httpRespError := client.Do(req)
	if httpRespError != nil {
		return &gatewayResponse, httpRespError
	}

	responseString, error := extractString(httpResp)
	if error != nil {
		return &gatewayResponse, fmt.Errorf("Error Extracting Response: %s", error)
	}

	if httpResp.StatusCode != 200 {

		err := json.Unmarshal([]byte(responseString), &gatewayResponse)

		if err != nil {
			return &gatewayResponse, fmt.Errorf("Error Validating MDM Details: %s", err)
		}

		return &gatewayResponse, nil
	}

	if responseString == "" {
		return &gatewayResponse, fmt.Errorf("Error Getting Cluster Details")
	}

	if gc.version == "4.0" {
		err := storeCookie(httpResp.Header, gc.host)
		if err != nil {
			return &gatewayResponse, fmt.Errorf("Error While Storing cookie: %s", err)
		}
	}

	if requireJSONOutput {
		gatewayResponse.StatusCode = 200

		gatewayResponse.Data = responseString

		return &gatewayResponse, nil
	}

	var mdmTopologyDetails types.MDMTopologyDetails

	err := json.Unmarshal([]byte(responseString), &mdmTopologyDetails)

	if err != nil {
		return &gatewayResponse, fmt.Errorf("Error For Get Cluster Details: %s", err)
	}

	gatewayResponse.StatusCode = 200

	gatewayResponse.ClusterDetails = mdmTopologyDetails

	return &gatewayResponse, nil
}

// DeletePackage used for delete packages from gateway server
func (gc *GatewayClient) DeletePackage(packageName string) (*types.GatewayResponse, error) {

	var gatewayResponse types.GatewayResponse

	req, httpError := http.NewRequest("DELETE", gc.host+"/im/types/installationPackages/instances/actions/delete::"+packageName, nil)
	if httpError != nil {
		return &gatewayResponse, httpError
	}
	if gc.version == "4.0" {
		req.Header.Set("Authorization", "Bearer "+gc.token)

		err := setCookie(req.Header, gc.host)
		if err != nil {
			return nil, fmt.Errorf("Error While Handling Cookie: %s", err)
		}
	} else {
		req.Header.Set("Authorization", "Basic "+base64.StdEncoding.EncodeToString([]byte(gc.username+":"+gc.password)))
	}
	req.Header.Set("Content-Type", "application/json")

	client := gc.http
	httpResp, httpRespError := client.Do(req)
	if httpRespError != nil {
		return &gatewayResponse, httpRespError
	}

	responseString, error := extractString(httpResp)
	if error != nil {
		return &gatewayResponse, fmt.Errorf("Error Extracting Response: %s", error)
	}

	if httpResp.StatusCode != 200 {

		err := json.Unmarshal([]byte(responseString), &gatewayResponse)

		if err != nil {
			return &gatewayResponse, fmt.Errorf("Error For Delete Package: %s", err)
		}

		return &gatewayResponse, nil
	}

	if gc.version == "4.0" {
		err := storeCookie(httpResp.Header, gc.host)
		if err != nil {
			return &gatewayResponse, fmt.Errorf("Error While Storing cookie: %s", err)
		}
	}

	gatewayResponse.StatusCode = 200

	return &gatewayResponse, nil
}

// BeginInstallation used for start installation
func (gc *GatewayClient) BeginInstallation(jsonStr, mdmUsername, mdmPassword, liaPassword string, allowNonSecureCommunicationWithMdm, allowNonSecureCommunicationWithLia, disableNonMgmtComponentsAuth, expansion bool) (*types.GatewayResponse, error) {

	var gatewayResponse types.GatewayResponse

	mapData, jsonParseError := jsonToMap(jsonStr)
	if jsonParseError != nil {
		return &gatewayResponse, jsonParseError
	}

	mapData["mdmPassword"] = mdmPassword
	mapData["mdmUser"] = mdmUsername
	mapData["liaPassword"] = liaPassword
	mapData["liaLdapInitialMode"] = "NATIVE_AUTHENTICATION"

	secureData := map[string]interface{}{
		"allowNonSecureCommunicationWithMdm": allowNonSecureCommunicationWithMdm,
		"allowNonSecureCommunicationWithLia": allowNonSecureCommunicationWithLia,
		"disableNonMgmtComponentsAuth":       disableNonMgmtComponentsAuth,
	}
	mapData["securityConfiguration"] = secureData

	finalJSON, _ := json.Marshal(mapData)

	u, _ := url.Parse(gc.host + "/im/types/Configuration/actions/install")
	q := u.Query()

	if gc.version == "4.0" && !expansion {
		q.Set("noSecurityBootstrap", "false")
	} else {
		q.Set("noUpload", "false")
		q.Set("noInstall", "false")
		q.Set("noConfigure", "false")
		q.Set("noLinuxDevValidation", "false")
		q.Set("globalZeroPadPolicy", "false")
	}

	if expansion {
		q.Set("extend", strconv.FormatBool(expansion))
	}

	u.RawQuery = q.Encode()

	req, httpError := http.NewRequest("POST", u.String(), bytes.NewBuffer(finalJSON))
	if httpError != nil {
		return &gatewayResponse, httpError
	}
	if gc.version == "4.0" {
		req.Header.Set("Authorization", "Bearer "+gc.token)

		err := setCookie(req.Header, gc.host)
		if err != nil {
			return nil, fmt.Errorf("Error While Handling Cookie: %s", err)
		}
	} else {
		req.Header.Set("Authorization", "Basic "+base64.StdEncoding.EncodeToString([]byte(gc.username+":"+gc.password)))
	}
	req.Header.Set("Content-Type", "application/json")

	client := gc.http

	httpResp, httpRespError := client.Do(req)
	if httpRespError != nil {
		return &gatewayResponse, httpRespError
	}

	if httpResp.StatusCode != 202 {

		responseString, error := extractString(httpResp)
		if error != nil {
			return &gatewayResponse, fmt.Errorf("Error Extracting Response: %s", error)
		}

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
	if gc.version == "4.0" {
		req.Header.Set("Authorization", "Bearer "+gc.token)

		err := setCookie(req.Header, gc.host)
		if err != nil {
			return nil, fmt.Errorf("Error While Handling Cookie: %s", err)
		}
	} else {
		req.Header.Set("Authorization", "Basic "+base64.StdEncoding.EncodeToString([]byte(gc.username+":"+gc.password)))
	}
	req.Header.Set("Content-Type", "application/json")

	client := gc.http

	httpResp, httpRespError := client.Do(req)
	if httpRespError != nil {
		return &gatewayResponse, httpRespError
	}

	responseString, error := extractString(httpResp)
	if error != nil {
		return &gatewayResponse, fmt.Errorf("Error Extracting Response: %s", error)
	}

	if httpResp.StatusCode != 200 {

		err := json.Unmarshal([]byte(responseString), &gatewayResponse)

		if err != nil {
			return &gatewayResponse, fmt.Errorf("Error For Move To Next Phase: %s", err)
		}

		return &gatewayResponse, nil
	}

	if gc.version == "4.0" {
		err := storeCookie(httpResp.Header, gc.host)
		if err != nil {
			return &gatewayResponse, fmt.Errorf("Error While Storing cookie: %s", err)
		}
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
	if gc.version == "4.0" {
		req.Header.Set("Authorization", "Bearer "+gc.token)

		err := setCookie(req.Header, gc.host)
		if err != nil {
			return nil, fmt.Errorf("Error While Handling Cookie: %s", err)
		}
	} else {
		req.Header.Set("Authorization", "Basic "+base64.StdEncoding.EncodeToString([]byte(gc.username+":"+gc.password)))
	}
	req.Header.Set("Content-Type", "application/json")

	client := gc.http

	httpResp, httpRespError := client.Do(req)
	if httpRespError != nil {
		return &gatewayResponse, httpRespError
	}

	responseString, error := extractString(httpResp)
	if error != nil {
		return &gatewayResponse, fmt.Errorf("Error Extracting Response: %s", error)
	}

	if httpResp.StatusCode != 200 {

		err := json.Unmarshal([]byte(responseString), &gatewayResponse)

		if err != nil {
			return &gatewayResponse, fmt.Errorf("Error For Retry Phase: %s", err)
		}

		return &gatewayResponse, nil
	}

	if gc.version == "4.0" {
		err := storeCookie(httpResp.Header, gc.host)
		if err != nil {
			return &gatewayResponse, fmt.Errorf("Error While Storing cookie: %s", err)
		}
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
	if gc.version == "4.0" {
		req.Header.Set("Authorization", "Bearer "+gc.token)

		err := setCookie(req.Header, gc.host)
		if err != nil {
			return nil, fmt.Errorf("Error While Handling Cookie: %s", err)
		}
	} else {
		req.Header.Set("Authorization", "Basic "+base64.StdEncoding.EncodeToString([]byte(gc.username+":"+gc.password)))
	}
	req.Header.Set("Content-Type", "application/json")

	client := gc.http

	httpResp, httpRespError := client.Do(req)
	if httpRespError != nil {
		return &gatewayResponse, httpRespError
	}

	responseString, error := extractString(httpResp)
	if error != nil {
		return &gatewayResponse, fmt.Errorf("Error Extracting Response: %s", error)
	}

	if httpResp.StatusCode != 200 {

		err := json.Unmarshal([]byte(responseString), &gatewayResponse)

		if err != nil {
			return &gatewayResponse, fmt.Errorf("Error For Abort Operation: %s", err)
		}

		return &gatewayResponse, nil
	}

	if gc.version == "4.0" {
		err := storeCookie(httpResp.Header, gc.host)
		if err != nil {
			return &gatewayResponse, fmt.Errorf("Error While Storing cookie: %s", err)
		}
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
	if gc.version == "4.0" {
		req.Header.Set("Authorization", "Bearer "+gc.token)

		err := setCookie(req.Header, gc.host)
		if err != nil {
			return nil, fmt.Errorf("Error While Handling Cookie: %s", err)
		}
	} else {
		req.Header.Set("Authorization", "Basic "+base64.StdEncoding.EncodeToString([]byte(gc.username+":"+gc.password)))
	}
	req.Header.Set("Content-Type", "application/json")

	client := gc.http

	httpResp, httpRespError := client.Do(req)
	if httpRespError != nil {
		return &gatewayResponse, httpRespError
	}

	responseString, error := extractString(httpResp)
	if error != nil {
		return &gatewayResponse, fmt.Errorf("Error Extracting Response: %s", error)
	}

	if httpResp.StatusCode != 200 {

		err := json.Unmarshal([]byte(responseString), &gatewayResponse)

		if err != nil {
			return &gatewayResponse, fmt.Errorf("Error For Clear Queue Commands: %s", err)
		}

		return &gatewayResponse, nil
	}

	if gc.version == "4.0" {
		err := storeCookie(httpResp.Header, gc.host)
		if err != nil {
			return &gatewayResponse, fmt.Errorf("Error While Storing cookie: %s", err)
		}
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
	if gc.version == "4.0" {
		req.Header.Set("Authorization", "Bearer "+gc.token)

		err := setCookie(req.Header, gc.host)
		if err != nil {
			return nil, fmt.Errorf("Error While Handling Cookie: %s", err)
		}
	} else {
		req.Header.Set("Authorization", "Basic "+base64.StdEncoding.EncodeToString([]byte(gc.username+":"+gc.password)))
	}
	req.Header.Set("Content-Type", "application/json")

	client := gc.http

	httpResp, httpRespError := client.Do(req)
	if httpRespError != nil {
		return &gatewayResponse, httpRespError
	}

	responseString, error := extractString(httpResp)
	if error != nil {
		return &gatewayResponse, fmt.Errorf("Error Extracting Response: %s", error)
	}

	if httpResp.StatusCode != 200 {

		err := json.Unmarshal([]byte(responseString), &gatewayResponse)

		if err != nil {
			return &gatewayResponse, fmt.Errorf("Error For Move To Ideal Phase: %s", err)
		}

		return &gatewayResponse, nil
	}

	if gc.version == "4.0" {
		err := storeCookie(httpResp.Header, gc.host)
		if err != nil {
			return &gatewayResponse, fmt.Errorf("Error While Storing cookie: %s", err)
		}
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
	if gc.version == "4.0" {
		req.Header.Set("Authorization", "Bearer "+gc.token)

		err := setCookie(req.Header, gc.host)
		if err != nil {
			return nil, fmt.Errorf("Error While Handling Cookie: %s", err)
		}
	} else {
		req.Header.Set("Authorization", "Basic "+base64.StdEncoding.EncodeToString([]byte(gc.username+":"+gc.password)))
	}
	req.Header.Set("Content-Type", "application/json")

	client := gc.http
	httpResp, httpRespError := client.Do(req)
	if httpRespError != nil {
		return mdmQueueCommandDetails, httpRespError
	}

	responseString, error := extractString(httpResp)
	if error != nil {
		return mdmQueueCommandDetails, fmt.Errorf("Error Extracting Response: %s", error)
	}

	if httpResp.StatusCode == 200 {

		if gc.version == "4.0" {
			err := storeCookie(httpResp.Header, gc.host)
			if err != nil {
				return mdmQueueCommandDetails, fmt.Errorf("Error While Storing cookie: %s", err)
			}
		}

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

		if currentPhase == mdmQueueCommandDetail.AllowedPhase && (mdmQueueCommandDetail.CommandState == "pending" || mdmQueueCommandDetail.CommandState == "running") {
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

// UninstallCluster used for uninstallation of cluster
func (gc *GatewayClient) UninstallCluster(jsonStr, mdmUsername, mdmPassword, liaPassword string, allowNonSecureCommunicationWithMdm, allowNonSecureCommunicationWithLia, disableNonMgmtComponentsAuth, expansion bool) (*types.GatewayResponse, error) {

	var gatewayResponse types.GatewayResponse

	clusterData, jsonParseError := jsonToMap(jsonStr)
	if jsonParseError != nil {
		return &gatewayResponse, jsonParseError
	}

	clusterData["mdmPassword"] = mdmPassword
	clusterData["mdmUser"] = mdmUsername
	clusterData["liaPassword"] = liaPassword
	clusterData["liaLdapInitialMode"] = "NATIVE_AUTHENTICATION"

	secureData := map[string]interface{}{
		"allowNonSecureCommunicationWithMdm": allowNonSecureCommunicationWithMdm,
		"allowNonSecureCommunicationWithLia": allowNonSecureCommunicationWithLia,
		"disableNonMgmtComponentsAuth":       disableNonMgmtComponentsAuth,
	}
	clusterData["securityConfiguration"] = secureData

	finalJSON, _ := json.Marshal(clusterData)

	u, _ := url.Parse(gc.host + "/im/types/Configuration/actions/uninstall")

	req, httpError := http.NewRequest("POST", u.String(), bytes.NewBuffer(finalJSON))
	if httpError != nil {
		return &gatewayResponse, httpError
	}
	if gc.version == "4.0" {
		req.Header.Set("Authorization", "Bearer "+gc.token)

		err := setCookie(req.Header, gc.host)
		if err != nil {
			return nil, fmt.Errorf("Error While Handling Cookie: %s", err)
		}
	} else {
		req.Header.Set("Authorization", "Basic "+base64.StdEncoding.EncodeToString([]byte(gc.username+":"+gc.password)))
	}
	req.Header.Set("Content-Type", "application/json")

	client := gc.http

	httpResp, httpRespError := client.Do(req)
	if httpRespError != nil {
		return &gatewayResponse, httpRespError
	}

	if httpResp.StatusCode != 202 {

		responseString, error := extractString(httpResp)
		if error != nil {
			return &gatewayResponse, fmt.Errorf("Error Extracting Response: %s", error)
		}

		err := json.Unmarshal([]byte(responseString), &gatewayResponse)

		if err != nil {
			return &gatewayResponse, fmt.Errorf("Error For Uninstall Cluster: %s", err)
		}

		return &gatewayResponse, nil
	}

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

const configFile = "/home/.cookie_config.yaml"

var globalCookie string

// CookieConfig represents the YAML structure
type CookieConfig struct {
	Hosts []Host `yaml:"hosts"`
}

// Host represents individual hosts in the YAML structure
type Host struct {
	Name           string `yaml:"name"`
	LegacyGWCookie string `yaml:"cookie"`
}

func storeCookie(header http.Header, host string) error {
	if header != nil && header["Set-Cookie"] != nil {

		newCookie := strings.Split(header["Set-Cookie"][0], ";")[0]
		sanitizedCookie := strings.ReplaceAll(strings.Split(newCookie, "=")[1], "|", "_")

		// Load existing configuration
		config, err := loadConfig()
		if err != nil {
			return err
		}

		// Check if the host already exists, and update or add accordingly
		found := false
		for i, h := range config.Hosts {
			if h.Name == host {
				config.Hosts[i].LegacyGWCookie = sanitizedCookie
				found = true
				break
			}
		}

		// If the host is not found, add a new host
		if !found {
			config.Hosts = append(config.Hosts, Host{Name: host, LegacyGWCookie: sanitizedCookie})
		}

		// Update the global variable directly
		globalCookie = sanitizedCookie

		err = writeConfig(config)
		if err != nil {
			return err
		}
	}

	return nil
}

func setCookie(header http.Header, host string) error {

	if globalCookie != "" {
		header.Set("Cookie", "LEGACYGWCOOKIE="+strings.ReplaceAll(globalCookie, "_", "|"))
	} else {
		config, err := loadConfig()
		if err != nil {
			return err
		}

		// Check if the host already exists and set the globalCookie
		for _, h := range config.Hosts {
			if h.Name == host {
				globalCookie = h.LegacyGWCookie
				header.Set("Cookie", "LEGACYGWCOOKIE="+strings.ReplaceAll(globalCookie, "_", "|"))
				break
			}
		}
	}

	return nil
}

func loadConfig() (*CookieConfig, error) {
	if _, err := os.Stat(configFile); err == nil {
		data, err := ioutil.ReadFile(configFile)
		if err != nil {
			return nil, err
		}

		var config CookieConfig
		err = yaml.Unmarshal(data, &config)
		if err != nil {
			return nil, err
		}

		return &config, nil
	}

	return &CookieConfig{}, nil
}

func writeConfig(config *CookieConfig) error {
	data, err := yaml.Marshal(&config)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(configFile, data, 0644)
	if err != nil {
		return err
	}

	return nil
}
