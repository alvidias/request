package request

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"
)

var errPOSTWrongDataType = errors.New("unrecognized data type for POST request")

type Requester interface {
	Do(req *http.Request) (*http.Response, error)
}

type JSONEncoder interface {
	Marshal(data any) ([]byte, error)
}

type JSONDecoder interface {
	Unmarshal(data []byte, v any) error
}

type Response *http.Response

type URL string

type Client struct {
	client      Requester
	jsonEncoder JSONEncoder
	Defaults    struct {
		BaseURL URL
		Headers http.Header
	}
}

func NewClient(opts ...ClientDefaults) *Client {
	c := &Client{
		client:      new(http.Client),
		jsonEncoder: new(jsonEncoder),
	}

	for _, opt := range opts {
		opt(c)
	}

	return c
}

func (c *Client) Post(url URL, contentType string, data any) (Response, error) {
	d, err := validateToByte(data)

	if err != nil {
		return nil, err
	}

	req, err := c.newRequest("POST", string(url), bytes.NewBuffer(d))

	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", contentType)

	return c.client.Do(req)
}

func (c *Client) Get(url URL) (Response, error) {
	req, err := c.newRequest("GET", string(url), nil)

	if err != nil {
		return nil, err
	}

	return c.client.Do(req)
}

func (c *Client) GetJSON(url URL) *JSONResponse {
	r := &JSONResponse{}

	r.Response, r.err = c.Get(url)

	return r
}

func (c *Client) PostJSON(url URL, data any) *JSONResponse {
	r := &JSONResponse{}
	enc := data

	if IsStructMapOrSlice(data) {
		js, err := c.jsonEncoder.Marshal(data)

		if err != nil {
			r.err = err

			return r
		}

		enc = js
	}

	r.Response, r.err = c.Post(url, "application/json", enc)

	return r
}

type JSONResponse struct {
	*http.Response
	body []byte
	err  error
}

func (j *JSONResponse) Parse(a any) error {
	if j.err != nil {
		return j.err
	}

	var c []byte
	var err error

	if j.body != nil {
		c = j.body
	} else {
		c, err = io.ReadAll(j.Body)

		if err != nil {
			return err
		}

		defer j.Body.Close()

		j.body = c
	}

	return json.Unmarshal(c, a)
}

func validateToByte(d any) ([]byte, error) {
	switch data := d.(type) {
	case []byte:
		return d.([]byte), nil
	case string:
		return []byte(data), nil
	default:
		return nil, errPOSTWrongDataType
	}
}

func (c *Client) newRequest(method string, URL string, data io.Reader) (*http.Request, error) {
	var req *http.Request
	var err error

	if strings.HasPrefix(URL, "http") {
		req, err = http.NewRequest(method, URL, data)
	} else {
		base := string(c.Defaults.BaseURL)
		req, err = http.NewRequest(method, base+URL, data)
	}

	if err != nil {
		return nil, err
	}

	if c.Defaults.Headers != nil {
		req.Header = c.Defaults.Headers
	}

	return req, nil
}
