package client

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
)

// RequestPost initiates a POST request from the client side. Accepts
// params as JSON, and returns response as decoded JSON.
func RequestPost(url string, params interface{}) (map[string]interface{}, error) {
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

	var reply map[string]interface{}
	err = json.Unmarshal(data, &reply)
	if err != nil {
		return nil, err
	}

	return reply, nil
}
