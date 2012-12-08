package httpapi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/srid/log"
	"io/ioutil"
	"net/http"
	"reflect"
)

// APIHandler handles a single API endpoint
type APIHandler struct {
	// RequestStruct is a struct to store the fields of request
	// paramemters, passed as JSON from the client.
	RequestStruct interface{}
}

func (h APIHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Infof("httpapi %s -- %s", h.RequestStruct, r)
	request := reflect.New(reflect.TypeOf(h.RequestStruct)).Interface().(RequestParams)

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		// FIXME: best way to report errors?
		log.Error(err)
		return
	}
	if err := json.Unmarshal(body, request); err != nil {
		log.Error(err)
		return
	}
	log.Info(request)
	data, err := request.HandleRequest()
	if err != nil {
		log.Error(err)
		fmt.Fprintf(w, "FAIL %s", err)
		return
	}
	fmt.Fprintf(w, data)
}

type RequestParams interface {
	// HandleRequest is called when a request comes in. It must return
	// the response string or the error object.
	// FIXME: pre-define errors for appropriate HTTP codes (404, 500) ...
	HandleRequest() (string, error)
}

// RequestPost initiates a POST request from the client side.
// FIXME: somehow separate client-side functions from the server-side
func RequestPost(url string, r RequestParams) ([]byte, error) {
	data, err := json.Marshal(r)
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
