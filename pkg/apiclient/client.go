package apiclient

import (
	"bytes"
	"encoding/json"
	"net/http"
)

type Client struct {
	BaseURL string
	Headers map[string]string
	client  *http.Client
}

func NewClient(baseURL string, headers map[string]string) *Client {
	return &Client{
		BaseURL: baseURL,
		Headers: headers,
		client:  &http.Client{},
	}
}

type JSONRequest struct {
	Method  string
	Url     string
	Headers map[string]string
	Data    map[string]interface{}
}

func (c *Client) NewJSONRequest(method string, subUrl string, headers map[string]string, data map[string]interface{}) *JSONRequest {
	allHeaders := map[string]string{}
	for k, v := range c.Headers {
		allHeaders[k] = v
	}
	for k, v := range headers {
		allHeaders[k] = v
	}
	return &JSONRequest{
		Method:  method,
		Url:     c.BaseURL + subUrl,
		Headers: allHeaders,
		Data:    data,
	}
}

func (c *Client) Do(req *JSONRequest) (*http.Response, error) {
	jsonData, err := json.Marshal(req.Data)
	if err != nil {
		return nil, err
	}

	request, err := http.NewRequest(req.Method, req.Url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}

	for k, v := range req.Headers {
		request.Header.Set(k, v)
	}
	return c.client.Do(request)
}
