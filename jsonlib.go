package jsonlib

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/pkg/errors"
)

// JSONLibrary json library interface
type JSONLibrary interface {
	RequestJSON(method, url string, headers map[string]string, req interface{}, res interface{}) error
	GetJSON(url string, headers map[string]string, res interface{}) error
	PostJSON(url string, headers map[string]string, req interface{}, res interface{}) error
	PutJSON(url string, headers map[string]string, req interface{}, res interface{}) error
	DeleteJSON(url string, headers map[string]string, req interface{}, res interface{}) error
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

func (m *jsonLibrary) getClient() *http.Client {
	if m._client != nil {
		return m._client
	}
	return &http.Client{Timeout: 30 * time.Second}
}

// RequestJSON ...
func (m *jsonLibrary) RequestJSON(method, url string, headers map[string]string, req interface{}, res interface{}) error {
	// get client
	client := m.getClient()
	// get body
	var body io.Reader
	var sData string
	if req != nil {
		data, err := json.Marshal(req)
		if err != nil {
			return errors.Wrap(err, "jsonlib: failed to marshal request")
		}
		sData = string(data)
		body = bytes.NewBuffer(data)
	}
	// new request
	reqt, err := http.NewRequest(method, url, body)
	if err != nil {
		return errors.Wrapf(err, "jsonlib: failed to new request. url:'%s'", url)
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
		return errors.Wrapf(err, "jsonlib: failed to post json, url:%s", url)
	}
	defer resp.Body.Close()
	// check status
	if (resp.StatusCode/100)*100 != http.StatusOK {
		return fmt.Errorf("jsonlib: failed to postJSON. status:'%s', url:'%s', data:'%s'",
			resp.Status, url, sData)
	}
	// decode
	if res != nil {
		if err := json.NewDecoder(resp.Body).Decode(res); err != nil {
			return errors.Wrapf(err, "jsonlib: failed to decode response, url:%s", url)
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
