package request

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestClient_ClientDefaults(t *testing.T) {
	url := URL("http://localhost:11434/api/generate")
	client := NewClient(BaseURL(url))

	assert.Equal(t, url, client.Defaults.BaseURL)
}

func TestClient_DefaultHeaders(t *testing.T) {
	headers := http.Header(map[string][]string{
		"Content-Type": {"application/json"},
	})

	client := NewClient(DefaultHeaders(headers))
	assert.Equal(t, headers, client.Defaults.Headers)
}
