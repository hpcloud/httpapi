package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type Client http.Client

var DefaultClient *Client

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
func (c *Client) DoRequest(req *http.Request, response interface{}) error {
	resp, err := (*http.Client)(c).Transport.RoundTrip(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if !IsHttpSuccess(resp.StatusCode) {
		return fmt.Errorf("HTTP request failed with code %d and body << %v >>.",
			resp.StatusCode, string(data))
	}

	err = json.Unmarshal(data, response)
	if err != nil {
		// fmt.Printf("==> %v <==\n", string(data))
		return fmt.Errorf("Response not in JSON format; %v", err)
	}

	return nil
}

// Post is a version of http.Post accepting JSON params and returning
// the same.
func Post(url string, params interface{}, response interface{}) error {
	req, err := NewRequest("POST", url, params)
	if err != nil {
		return err
	}

	return DefaultClient.DoRequest(req, response)
}

func IsHttpSuccess(code int) bool {
	// XXX: accept other codes
	return (code == 200 || code == 201 || code == 302)
}

func init() {
	DefaultClient = &Client{Transport: http.DefaultTransport}
}
