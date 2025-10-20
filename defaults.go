package request

import "net/http"

type ClientDefaults func(*Client)

func BaseURL(url URL) ClientDefaults {
	return func(c *Client) {
		c.Defaults.BaseURL = url
	}
}

func DefaultHeaders(headers http.Header) ClientDefaults {
	return func(c *Client) {
		c.Defaults.Headers = headers
	}
}
