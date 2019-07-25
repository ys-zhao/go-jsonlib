package jsonlib

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"time"
)

// JSONLibrary json library interface
type JSONLibrary interface {
	RequestJSON(method, url string, headers map[string]string, req interface{}, res interface{}) error
	GetJSON(url string, headers map[string]string, res interface{}) error
	PostJSON(url string, headers map[string]string, req interface{}, res interface{}) error
	PutJSON(url string, headers map[string]string, req interface{}, res interface{}) error
	DeleteJSON(url string, headers map[string]string, req interface{}, res interface{}) error
	ParseJSONRequest(r *http.Request, res interface{}) error
}

// Error struct
type Error struct {
	InternalError error
	Message       string

	URL     string
	Method  string
	Request []byte

	Status   int
	Response []byte
}

func (m *Error) Error() string {
	return fmt.Sprintf("error=%s; message=%s; url=%s; method=%s; request=%s; status=%d; response=%s",
		m.InternalError, m.Message, m.URL, m.Method, m.Request, m.Status, m.Response)
}

type jsonLibrary struct {
	_client *http.Client
}

// Default the default json library
var Default = New(nil)

// New return a new json library
func New(client *http.Client) JSONLibrary {
	if client == nil {
		client = &http.Client{Timeout: 30 * time.Second}
	}
	return &jsonLibrary{
		_client: client,
	}
}

// RequestJSON ...
func RequestJSON(method, url string, headers map[string]string, req interface{}, res interface{}) error {
	return Default.RequestJSON(method, url, headers, req, res)
}

// GetJSON ...
func GetJSON(url string, headers map[string]string, res interface{}) error {
	return Default.GetJSON(url, headers, res)
}

// PostJSON ...
func PostJSON(url string, headers map[string]string, req interface{}, res interface{}) error {
	return Default.PostJSON(url, headers, req, res)
}

// PutJSON ...
func PutJSON(url string, headers map[string]string, req interface{}, res interface{}) error {
	return Default.PutJSON(url, headers, req, res)
}

// DeleteJSON ...
func DeleteJSON(url string, headers map[string]string, req interface{}, res interface{}) error {
	return Default.DeleteJSON(url, headers, req, res)
}

// ParseJSONRequest ...
func ParseJSONRequest(r *http.Request, res interface{}) error {
	return Default.ParseJSONRequest(r, res)
}

func (m *jsonLibrary) getClient() *http.Client {
	if m._client != nil {
		return m._client
	}
	return &http.Client{Timeout: 30 * time.Second}
}

// RequestJSON ...
func (m *jsonLibrary) RequestJSON(method, url string, headers map[string]string, req interface{}, res interface{}) (err error) {
	// get client
	client := m.getClient()
	// get body
	var body io.Reader
	var reqData []byte
	if req != nil {
		if reqData, err = json.Marshal(req); err != nil {
			return &Error{
				InternalError: err,
				URL:           url,
				Method:        method,
				Request:       reqData,
				Message:       "jsonlib: failed to marshal request",
			}
		}
		body = bytes.NewBuffer(reqData)
	}
	// new request
	reqt, err := http.NewRequest(method, url, body)
	if err != nil {
		return &Error{
			InternalError: err,
			URL:           url,
			Method:        method,
			Request:       reqData,
			Message:       "jsonlib: failed to new request",
		}
	}
	// set header
	reqt.Header.Set("Content-Type", "application/json")
	if headers != nil {
		for key, value := range headers {
			reqt.Header.Set(key, value)
		}
	}
	// send the request
	resp, err := client.Do(reqt)
	if err != nil {
		return &Error{
			InternalError: err,
			URL:           url,
			Method:        method,
			Request:       reqData,
			Message:       "jsonlib: failed to send request",
		}
	}
	defer resp.Body.Close()
	resData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return &Error{
			InternalError: err,
			URL:           url,
			Method:        method,
			Request:       reqData,
			Message:       "jsonlib: failed to read response",
		}
	}
	// check status
	if (resp.StatusCode/100)*100 != http.StatusOK {
		return &Error{
			Method:   method,
			URL:      url,
			Request:  reqData,
			Status:   resp.StatusCode,
			Response: resData,
			Message:  "jsonlib: failed to receive response",
		}
	}
	// decode
	if res != nil {
		if err := json.Unmarshal(resData, res); err != nil {
			return &Error{
				Method:   method,
				URL:      url,
				Request:  reqData,
				Status:   resp.StatusCode,
				Response: resData,
				Message:  "jsonlib: failed to unmarshal response",
			}
		}
	}
	// return
	return nil
}

// GetJSON ...
func (m *jsonLibrary) GetJSON(url string, headers map[string]string, res interface{}) error {
	return m.RequestJSON("GET", url, headers, nil, res)
}

// PostJSON ...
func (m *jsonLibrary) PostJSON(url string, headers map[string]string, req interface{}, res interface{}) error {
	return m.RequestJSON("POST", url, headers, req, res)
}

// PutJSON ...
func (m *jsonLibrary) PutJSON(url string, headers map[string]string, req interface{}, res interface{}) error {
	return m.RequestJSON("PUT", url, headers, req, res)
}

// DeleteJSON ...
func (m *jsonLibrary) DeleteJSON(url string, headers map[string]string, req interface{}, res interface{}) error {
	return m.RequestJSON("DELETE", url, headers, req, res)
}

// ParseJSONRequest ...
func (m *jsonLibrary) ParseJSONRequest(r *http.Request, res interface{}) error {
	defer r.Body.Close()
	return json.NewDecoder(r.Body).Decode(res)
}
