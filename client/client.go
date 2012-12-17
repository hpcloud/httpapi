package client

import (
	"fmt"
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
)

type Client http.Client

var DefaultClient = &Client{}

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
func (c *Client) DoRequest(req *http.Request, response interface{}) error{
	resp, err := (*http.Client)(c).Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	err = json.Unmarshal(data, response)
	if err != nil {
		return err
	}

	// XXX: accept other codes
	if !(resp.StatusCode == 200 || resp.StatusCode == 302) {
		return fmt.Errorf("HTTP request with failure code (%d); body -- %v",
			resp.StatusCode, response)
	}
	
	return nil
}

// Post is a version of http.Post accepting JSON params and returning
// the same.
func Post(url string, params interface{}) (Hash, error) {
	req, err := NewRequest("POST", url, params)
	if err != nil {
		return nil, err
	}

	var h Hash
	return h, DefaultClient.DoRequest(req, &h)
}
