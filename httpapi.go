package httpapi

import (
	"encoding/json"
	"fmt"
	"github.com/ActiveState/log"
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
	// Create a logger with a prefix unique to this request. Assuming
	// that a new request object is created per request, its memory
	// address (%p) should give us an unique identifier.
	l := log.New()
	l.SetPrefix(fmt.Sprintf("[HTTP:%p] ", r))
	
	l.Infof("%+v", r)
	request := reflect.New(reflect.TypeOf(h.RequestStruct)).Interface().(RequestParams)

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		l.Error(err)
		return
	}
	if err := json.Unmarshal(body, request); err != nil {
		l.Errorf("Unable to decode JSON body in POST request (%s). Original body was: %s", err, string(body))
		return
	}

	response, err := request.HandleRequest()
	if err != nil {
		l.Errorf("Failed to handle this request -- %s", err)
		http.Error(w, err.Error(), 500)
	} else {
		data, err := json.Marshal(response)
		if err != nil {
			err = fmt.Errorf("Unable to response into JSON: %s", err)
			http.Error(w, err.Error(), 500)
		} else {
			_, err = w.Write(data)
			if err != nil {
				l.Errorf("Unable to write http response: %s", err)
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
