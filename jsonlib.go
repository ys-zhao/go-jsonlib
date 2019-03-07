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
	GetJSON(url string, res interface{}) error
	PostJSON(url string, req interface{}, res interface{}) error
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
func GetJSON(url string, res interface{}) error {
	return Default.GetJSON(url, res)
}

// PostJSON ...
func PostJSON(url string, req interface{}, res interface{}) error {
	return Default.PostJSON(url, req, res)
}

func (m *jsonLibrary) getClient() *http.Client {
	if m._client != nil {
		return m._client
	}
	return &http.Client{Timeout: 30 * time.Second}
}

// GetJSON ...
func (m *jsonLibrary) GetJSON(url string, res interface{}) error {
	client := m.getClient()
	resp, err := client.Get(url)
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
func (m *jsonLibrary) PostJSON(url string, req interface{}, res interface{}) error {
	client := m.getClient()
	data, _ := json.Marshal(req)
	resp, err := client.Post(url, "application/json", bytes.NewBuffer(data))
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
