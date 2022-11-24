package goscaleio

import (
	"net/http"
	"software/template"
	"time"
)

func (c *Client) CreateTemplate() *Client {
	// func (c *Client) CreateTemplate(value interface{}) *Client {
	// c.FringeObject = value

	// switch v := c.FringeObject.(type) {
	// case string:
	// 	fmt.Printf("String:")
	// case template.DefaultTemplate:
	// 	fmt.Printf("template.DefaultTemplate:")
	// default:
	// 	fmt.Println("I don't know, ask stackoverflow.", v)
	// }

	return c
}

// CreateTemplate creates a blank template
func (c *Client) FromModel(
	templateParameters template.DefaultTemplate) (interface{}, error) {
	defer TimeSpent("CreateTemplate", time.Now())

	path := "api/v1/ServiceTemplate"

	backResponse, err := c.authorizedJSONWithRetry(
		http.MethodPost, path, templateParameters)
	if err != nil {
		return nil, err
	}

	return backResponse, nil
}
func (c *Client) FromString(
	templateString string) (interface{}, error) {
	defer TimeSpent("CreateTemplate", time.Now())

	path := "api/v1/ServiceTemplate"
	backResponse, err := c.authorizedJSONWithRetry(
		http.MethodPost, path, templateString)
	if err != nil {
		return nil, err
	}

	return backResponse, nil
}
func (c *Client) UpdateTemplate(
	templateString, templateID string) (interface{}, error) {
	defer TimeSpent("CreateTemplate", time.Now())

	path := "/api/v1/ServiceTemplate/" + templateID
	backResponse, err := c.authorizedJSONWithRetry(
		http.MethodPut, path, templateString)
	if err != nil {
		return nil, err
	}

	return backResponse, nil
}
func (c *Client) GetTemplate(
	templateName string) (interface{}, error) {
	defer TimeSpent("GetTemplate", time.Now())

	path := "/api/v1/ServiceTemplate?filter=eq,name," + templateName

	var body interface{}
	response, err := c.authorizedJSONWithRetry(http.MethodGet, path, body)
	if err != nil {
		return nil, err
	}

	return response, nil
}
func (c *Client) DeleteTemplate(
	templateId string) (interface{}, error) {

	defer TimeSpent("DeleteTemplate", time.Now())

	path := "/api/v1/ServiceTemplate/" + templateId

	var body interface{}
	response, err := c.authorizedJSONWithRetry(http.MethodDelete, path, body)
	if err != nil {
		return nil, err
	}

	return response, nil
}
