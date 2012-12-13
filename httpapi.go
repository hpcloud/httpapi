package httpapi

import (
	"encoding/json"
	"fmt"
	"github.com/srid/log"
	"io/ioutil"
	"net/http"
	"reflect"
)

// Handler handles a single API endpoint
type Handler struct {
	// RequestStruct is a struct to store the fields of request
	// paramemters, passed as JSON from the client.
	RequestStruct interface{}
}

// FIXME: best way to report errors in ServerHTTP?
func (h Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Infof("httpapi -- %+v", r)
	request := reflect.New(reflect.TypeOf(h.RequestStruct)).Interface().(RequestParams)

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Error(err)
		return
	}
	if err := json.Unmarshal(body, request); err != nil {
		log.Errorf("Error decoding JSON body in POST request (%s). Original body was: %s", err, string(body))
		return
	}

	response, err := request.HandleRequest()
	if err != nil {
		log.Errorf("Request error -- %s", err)
		http.Error(w, err.Error(), 500)
	} else {
		data, err := json.Marshal(response)
		if err != nil {
			err = fmt.Errorf("Error encoding response into JSON: %s", err)
			http.Error(w, err.Error(), 500)
		} else {
			_, err = w.Write(data)
			if err != nil {
				log.Errorf("Failed to write http response: %s", err)
			}
		}
	}
}

type RequestParams interface {
	// HandleRequest is called when a request comes in. POST body will
	// be decoded into the receiver; returned value will be encoded to
	// JSON before responding to the client.
	// FIXME: pre-define errors for appropriate HTTP codes (404, 500) ...
	HandleRequest() (interface{}, error)
}
