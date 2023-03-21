package goscaleio

import (
	"bytes"
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"crypto/x509"
)

var (
	errNewClient = errors.New("missing endpoint")
	errSysCerts  = errors.New("Unable to initialize cert pool from system")
)

type gatewayclient struct {
	http     *http.Client
	host     string
	username string
	password string
}

// NewGateway returns a new gateway client.
func NewGateway(
	host string, username, password string, insecure,useCerts bool) (GatewayClient, error) {

	if host == "" {
		return nil, errNewClient
	}

	gc := &gatewayclient{
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

// GatewayClient is an  inetface which has all the functionalities for the gateway.
type GatewayClient interface {
	UploadPackages(fliePath string) error
	ParseCSV(filePath string) error
	BeginInstallation(jsonStr, mdmUsername, mdmPassword, liaPassword string) error
}

func (gc *gatewayclient) UploadPackages(filePath string) error {
	file, err1 := os.Open(filePath)
	if err1 != nil {
		return err1
	}
	defer file.Close()
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	part, err2 := writer.CreateFormFile("files", filepath.Base(filePath))
	if err2 != nil {
		return err2
	}
	_, err3 := io.Copy(part, file)
	if err3 != nil {
		return err3
	}
	err4 := writer.Close()
	if err4 != nil {
		return err4
	}

	req, err5 := http.NewRequest("POST", gc.host+"/im/types/installationPackages/instances/actions/uploadPackages", body)
	if err5 != nil {
		return err5
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("Authorization", "Basic "+base64.StdEncoding.EncodeToString([]byte(gc.username+":"+gc.password)))
	client := gc.http
	_, err6 := client.Do(req)
	if err6 != nil {
		return err6
	}
	return nil
}

func (gc *gatewayclient) ParseCSV(filePath string) error {
	file, err1 := os.Open(filePath)
	if err1 != nil {
		return err1
	}
	defer file.Close()
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	part, err2 := writer.CreateFormFile("file", filepath.Base(filePath))
	if err2 != nil {
		return err2
	}
	_, err3 := io.Copy(part, file)
	if err3 != nil {
		return err3
	}
	err4 := writer.Close()
	if err4 != nil {
		return err4
	}

	req, err5 := http.NewRequest("POST", gc.host+"/im/types/Configuration/instances/actions/parseFromCSV", body)
	if err5 != nil {
		return err5
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("Authorization", "Basic "+base64.StdEncoding.EncodeToString([]byte(gc.username+":"+gc.password)))
	fmt.Println(body)
	client := gc.http
	_, err6 := client.Do(req)
	if err6 != nil {
		return err6
	}

	return nil

}

func (gc *gatewayclient) BeginInstallation(jsonStr, mdmUsername, mdmPassword, liaPassword string) error {
	mapData := jsonToMap(jsonStr)
	mapData["mdmPassword"] = mdmPassword
	mapData["mdmUser"] = mdmUsername
	mapData["liaPassword"] = liaPassword
	finalJSON, _ := json.Marshal(mapData)
	req, err1 := http.NewRequest("POST", gc.host+"/im/types/Configuration/actions/install", bytes.NewBuffer(finalJSON))
	if err1 != nil {
		return err1
	}
	req.Header.Set("Authorization", "Basic "+base64.StdEncoding.EncodeToString([]byte(gc.username+":"+gc.password)))
	req.Header.Set("Content-Type", "application/json")
	client := gc.http
	_, err2 := client.Do(req)
	if err2 != nil {
		return err2
	}
	return nil
}

func jsonToMap(jsonStr string) map[string]interface{} {
	result := make(map[string]interface{})
	json.Unmarshal([]byte(jsonStr), &result)
	return result
}
