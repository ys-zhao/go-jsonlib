package jsonlib

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/pkg/errors"
)

// JSONLibrary json library interface
type JSONLibrary interface {
	GetJSON(url string, headers map[string]string, res interface{}) error
	PostJSON(url string, headers map[string]string, req interface{}, res interface{}) error
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

// GetJSON ...
func GetJSON(url string, headers map[string]string, res interface{}) error {
	return Default.GetJSON(url, headers, res)
}

// PostJSON ...
func PostJSON(url string, headers map[string]string, req interface{}, res interface{}) error {
	return Default.PostJSON(url, headers, req, res)
}

func (m *jsonLibrary) getClient() *http.Client {
	if m._client != nil {
		return m._client
	}
	return &http.Client{Timeout: 30 * time.Second}
}

// GetJSON ...
func (m *jsonLibrary) GetJSON(url string, headers map[string]string, res interface{}) error {
	client := m.getClient()
	// new request
	reqt, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return errors.Wrapf(err, "jsonlib: failed to create request. url:'%s'", url)
	}
	// set header
	if headers != nil {
		for key, value := range headers {
			reqt.Header.Set(key, value)
		}
	}
	resp, err := client.Do(reqt)
	if err != nil {
		return errors.Wrap(err, "jsonlib: failed to get json")
	}
	defer resp.Body.Close()
	// check status
	if (resp.StatusCode/100)*100 != http.StatusOK {
		return fmt.Errorf("jsonlib: failed to get json. Status='%s', Url='%s'", resp.Status, url)
	}
	// json parse
	return json.NewDecoder(resp.Body).Decode(res)
}

// PostJSON ...
func (m *jsonLibrary) PostJSON(url string, headers map[string]string, req interface{}, res interface{}) error {
	// get client
	client := m.getClient()
	// get post body
	data, err := json.Marshal(req)
	if err != nil {
		return errors.Wrap(err, "jsonlib: failed to marshal request")
	}
	// new request
	reqt, err := http.NewRequest("POST", url, bytes.NewBuffer(data))
	if err != nil {
		return errors.Wrapf(err, "jsonlib: failed to create request. url:'%s'", url)
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
			resp.Status, url, string(data))
	}
	// decode
	if err := json.NewDecoder(resp.Body).Decode(res); err != nil {
		return errors.Wrapf(err, "jsonlib: failed to decode json, url:%s", url)
	}
	// json parse
	return nil
}
