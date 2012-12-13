package client

import (
	"fmt"
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
)

type Client http.Client

// Value type for json decoded values for request params and response
// body of a HTTP API
type Hash map[string]interface{}

// NewRequest is wrapper over http.NewRequest handling json encoding
// for params.
func NewRequest(method string, url string, params interface{}) (*http.Request, error) {
	body, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(method, url, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	return req, nil
}

// DoRequest a wrapper over Do handling json encoding
func (c *Client) DoRequest(req *http.Request) (Hash, error){
	resp, err := (*http.Client)(c).Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	
	var v Hash
	err = json.Unmarshal(data, &v)
	if err != nil {
		return nil, err
	}

	// XXX: accept other codes
	if !(resp.StatusCode == 200 || resp.StatusCode == 302) {
		return nil, fmt.Errorf("HTTP request with failure code (%d); body -- %v",
			resp.StatusCode, v)
	}
	
	return v, nil
}

// RequestPost initiates a POST request from the client side. Accepts
// params as JSON, and returns response as decoded JSON.
func RequestPost(url string, params interface{}) (Hash, error) {
	data, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}
	response, err := http.Post(
		url, "application/json", bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	data, err = ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	var reply Hash
	err = json.Unmarshal(data, &reply)
	if err != nil {
		return nil, err
	}

	return reply, nil
}
