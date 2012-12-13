package client

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
)

// RequestPost initiates a POST request from the client side.
func RequestPost(url string, params interface{}) ([]byte, error) {
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

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}
