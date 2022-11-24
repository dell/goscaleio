package goscaleio

import (
	"fmt"
	"net/http"
	"software/template"
	"time"
)

// CreateVolume creates a blank template
func (c *Client) CreateTemplate(
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
func (c *Client) DiscardTemplate(
	templateId string) (interface{}, error) {
	defer TimeSpent("GetTemplate", time.Now())

	path := "/Api/V1/ui/templates/discardtemplate"

	type deletebody struct {
		RequestObj []string `json:"requestObj"`
	}
	delete := deletebody{RequestObj: []string{templateId}}
	fmt.Println("\n\ndelete", delete, "\n\n")
	response, err := c.authorizedJSONWithRetry(http.MethodPost, path, delete)
	if err != nil {
		return nil, err
	}

	return response, nil
}
