package goscaleio

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"errors"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	path "path/filepath"
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

// GatewayFunction is an  interface which has all the functionalities for the gateway.
type GatewayFunction interface {
	UploadPackages(fliePath string) error
	ParseCSV(filePath string) error
	BeginInstallation(jsonStr, mdmUsername, mdmPassword, liaPassword string) error
}

// UploadPackages used for upload packge to gateway server
func (gc *GatewayClient) UploadPackages(filePath string) error {

	file, filePathError := os.Open(path.Clean(filePath))
	if filePathError != nil {
		return filePathError
	}
	defer func() error {
		if err := file.Close(); err != nil {
			return err
		}
		return nil
	}()
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	part, fileReaderError := writer.CreateFormFile("files", path.Base(filePath))
	if fileReaderError != nil {
		return fileReaderError
	}
	_, fileContentError := io.Copy(part, file)
	if fileContentError != nil {
		return fileContentError
	}
	fileWriterError := writer.Close()
	if fileWriterError != nil {
		return fileWriterError
	}

	req, httpError := http.NewRequest("POST", gc.host+"/im/types/installationPackages/instances/actions/uploadPackages", body)
	if httpError != nil {
		return httpError
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("Authorization", "Basic "+base64.StdEncoding.EncodeToString([]byte(gc.username+":"+gc.password)))
	client := gc.http
	_, httpReqError := client.Do(req)
	if httpReqError != nil {
		return httpReqError
	}
	return nil
}

// ParseCSV used for upload CSV to gateway server and parse it
func (gc *GatewayClient) ParseCSV(filePath string) error {

	file, filePathError := os.Open(path.Clean(filePath))
	if filePathError != nil {
		return filePathError
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
		return fileReaderError
	}
	_, fileContentError := io.Copy(part, file)
	if fileContentError != nil {
		return fileContentError
	}
	fileWriterError := writer.Close()
	if fileWriterError != nil {
		return fileWriterError
	}

	req, httpError := http.NewRequest("POST", gc.host+"/im/types/Configuration/instances/actions/parseFromCSV", body)
	if httpError != nil {
		return httpError
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("Authorization", "Basic "+base64.StdEncoding.EncodeToString([]byte(gc.username+":"+gc.password)))
	client := gc.http
	_, httpReqError := client.Do(req)
	if httpReqError != nil {
		return httpReqError
	}

	return nil

}

// BeginInstallation used for start installation
func (gc *GatewayClient) BeginInstallation(jsonStr, mdmUsername, mdmPassword, liaPassword string) error {

	mapData, jsonParseError := jsonToMap(jsonStr)
	if jsonParseError != nil {
		return jsonParseError
	}
	mapData["mdmPassword"] = mdmPassword
	mapData["mdmUser"] = mdmUsername
	mapData["liaPassword"] = liaPassword
	finalJSON, _ := json.Marshal(mapData)

	req, httpError := http.NewRequest("POST", gc.host+"/im/types/Configuration/actions/install", bytes.NewBuffer(finalJSON))
	if httpError != nil {
		return httpError
	}
	req.Header.Set("Authorization", "Basic "+base64.StdEncoding.EncodeToString([]byte(gc.username+":"+gc.password)))
	req.Header.Set("Content-Type", "application/json")
	client := gc.http
	_, httpReqError := client.Do(req)
	if httpReqError != nil {
		return httpReqError
	}
	return nil
}

func jsonToMap(jsonStr string) (map[string]interface{}, error) {
	result := make(map[string]interface{})
	err := json.Unmarshal([]byte(jsonStr), &result)
	if err != nil {
		return nil, err
	}
	return result, nil
}
