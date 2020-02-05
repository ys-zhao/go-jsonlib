package jsonlib

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/clbanning/mxj"
	"github.com/pkg/errors"
)

// J2XOptions ...
type J2XOptions struct {
	WithIndent bool   // well formatted json
	Prefix     string // the json marshal prefix
	Indent     string // the json marshal indent
	RootTag    string // the json root tag
}

// J2XOption is a function on the options for a json2xml.
type J2XOption func(*J2XOptions) error

// DefaultJ2XOptions ...
func DefaultJ2XOptions() *J2XOptions {
	return &J2XOptions{
		WithIndent: false,
		Prefix:     "",
		Indent:     "\t",
		RootTag:    "doc",
	}
}

// J2XWithIndent ...
func J2XWithIndent(withIndent bool, prefix, indent string) J2XOption {
	return func(opts *J2XOptions) error {
		opts.WithIndent = withIndent
		opts.Prefix = prefix
		opts.Indent = indent
		return nil
	}
}

// J2XWithRootTag ...
func J2XWithRootTag(rootTag string) J2XOption {
	return func(opts *J2XOptions) error {
		opts.RootTag = rootTag
		return nil
	}
}

// X2JOptions ...
type X2JOptions struct {
	WithIndent bool   // well formatted json
	Prefix     string // the json marshal prefix
	Indent     string // the json marshal indent
	OmitRoot   bool   // omit the root element
}

// X2JOption is a function on the options for a xml2json.
type X2JOption func(*X2JOptions) error

// DefaultX2JOptions ...
func DefaultX2JOptions() *X2JOptions {
	return &X2JOptions{
		WithIndent: false,
		Prefix:     "",
		Indent:     "\t",
		OmitRoot:   false,
	}
}

// X2JWithIndent ...
func X2JWithIndent(withIndent bool, prefix, indent string) X2JOption {
	return func(opts *X2JOptions) error {
		opts.WithIndent = withIndent
		opts.Prefix = prefix
		opts.Indent = indent
		return nil
	}
}

// X2JWithOmitRoot ...
func X2JWithOmitRoot(omitRoot bool) X2JOption {
	return func(opts *X2JOptions) error {
		opts.OmitRoot = omitRoot
		return nil
	}
}

// JSONLibrary json library interface
type JSONLibrary interface {
	RequestJSON(method, url string, headers map[string]string, req interface{}, res interface{}) error
	GetJSON(url string, headers map[string]string, res interface{}) error
	PostJSON(url string, headers map[string]string, req interface{}, res interface{}) error
	PutJSON(url string, headers map[string]string, req interface{}, res interface{}) error
	DeleteJSON(url string, headers map[string]string, req interface{}, res interface{}) error
	ParseJSONRequest(r *http.Request, res interface{}) error
	JSON2XML(jsonStr string, opts ...J2XOption) (string, error)
	XML2JSON(xmlStr string, opts ...X2JOption) (string, error)
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

// JSON2XML ...
func JSON2XML(jsonStr string, opts ...J2XOption) (string, error) {
	return Default.JSON2XML(jsonStr, opts...)
}

// XML2JSON ...
func XML2JSON(xmlStr string, opts ...X2JOption) (string, error) {
	return Default.XML2JSON(xmlStr, opts...)
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

func (m *jsonLibrary) JSON2XML(jsonStr string, newOpts ...J2XOption) (string, error) {
	// get options
	opts := DefaultJ2XOptions()
	for _, opt := range newOpts {
		opt(opts)
	}
	// unmarshal from xml
	mv, err := mxj.NewMapJson([]byte(jsonStr))
	if err != nil {
		return "", errors.Wrap(err, "jsonlib: failed to unmarsh json data")
	}
	var xmlData []byte
	if opts.WithIndent {
		xmlData, err = mv.XmlIndent(opts.Prefix, opts.Indent, opts.RootTag)
	} else {
		xmlData, err = mv.Xml(opts.RootTag)
	}
	if err != nil {
		return "", errors.Wrap(err, "jsonlib: failed to marshal xml data")
	}
	return string(xmlData), nil
}

func (m *jsonLibrary) XML2JSON(xmlStr string, newOpts ...X2JOption) (string, error) {
	// get options
	opts := DefaultX2JOptions()
	for _, opt := range newOpts {
		opt(opts)
	}
	// unmarshal from xml
	mv, err := mxj.NewMapXml([]byte(xmlStr))
	if err != nil {
		return "", errors.Wrap(err, "jsonlib: failed to unmarsh xml data")
	}
	// omit the root element
	if opts.OmitRoot {
		for _, value := range mv {
			mv = value.(map[string]interface{})
			break
		}
	}
	// marsh to json
	var jsonData []byte
	if opts.WithIndent {
		jsonData, err = mv.JsonIndent(opts.Prefix, opts.Indent)
	} else {
		jsonData, err = mv.Json()
	}
	if err != nil {
		return "", errors.Wrap(err, "jsonlib: failed to marshal json data")
	}
	return string(jsonData), nil
}
