package request

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestNewClient(t *testing.T) {
	t.Run("New client", func(t *testing.T) {
		client := NewClient()

		if client == nil {
			t.Fatal("client not created")
		}
	})

	t.Run("New client with defaults", func(t *testing.T) {
		baseUrl := URL("http://localhost:11434")
		headers := http.Header{"Content-Type": {"application/json"}}

		client := NewClient(
			BaseURL(baseUrl),
			DefaultHeaders(http.Header{"Content-Type": {"application/json"}}),
		)

		assert.Equal(t, baseUrl, client.Defaults.BaseURL)
		assert.EqualValues(t, headers, client.Defaults.Headers)
	})
}

func TestClient_Post(t *testing.T) {
	client := NewClient()

	mockClient := &ClientMock{}
	client.client = mockClient
	bytesData := []byte(`{"test": "test"}`)
	stringData := "This is a test"

	cases := []struct {
		name     string
		ct       string
		data     any
		mockData *bytes.Buffer
		err      error
	}{
		{
			name:     "With bytes data",
			ct:       "application/json",
			data:     bytesData,
			mockData: bytes.NewBuffer(bytesData),
		},
		{
			name:     "With string data",
			ct:       "application/json",
			data:     stringData,
			mockData: bytes.NewBuffer([]byte(stringData)),
		},
		{
			name:     "With string data",
			ct:       "application/json",
			data:     map[string]string{"test": "test"},
			mockData: nil,
			err:      errPOSTWrongDataType,
		},
	}

	url := "http://localhost:11434/api/generate"

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if c.err == nil {
				mockClient.On(
					"Do",
					mock.MatchedBy(func(r *http.Request) bool {
						body, err := r.GetBody()
						assert.Nil(t, err)

						bd, err := io.ReadAll(body)
						assert.Nil(t, err)

						return r.Method == http.MethodPost &&
							r.URL.String() == url &&
							r.Header.Get("Content-Type") == c.ct &&
							string(bd) == c.mockData.String()
					}),
				).Return(nil, nil)
			}

			_, err := client.Post(URL(url), c.ct, c.data)

			if c.err != nil {
				assert.ErrorIs(t, err, c.err)
			} else {
				assert.Nil(t, err)
			}

			mockClient.AssertExpectations(t)
		})
	}
}

func TestClient_PostJSON(t *testing.T) {
	client := NewClient()
	mockClient := &ClientMock{}
	client.client = mockClient
	dataString := "This is a test"
	url := "http://localhost:11434/api/generate"

	mockClient.On("Do",
		mock.MatchedBy(func(r *http.Request) bool {
			body, err := r.GetBody()
			assert.Nil(t, err)

			bd, err := io.ReadAll(body)
			assert.Nil(t, err)

			return r.Method == http.MethodPost &&
				r.URL.String() == url &&
				r.Header.Get("Content-Type") == "application/json" &&
				string(bd) == dataString
		}),
	).Return(nil, nil)

	r := client.PostJSON(URL(url), dataString)

	assert.Nil(t, r.err)
	mockClient.AssertExpectations(t)
}

func TestClient_PostJSON_JsonData(t *testing.T) {
	client := NewClient()
	mockClient := &ClientMock{}
	client.client = mockClient

	type TestStruct struct {
		Test string `json:"test"`
	}

	structData := TestStruct{Test: "test"}
	ptrSructData := &TestStruct{Test: "test"}

	cases := []struct {
		name     string
		data     any
		mockData string
	}{
		{
			name:     "With map data",
			data:     map[string]string{"test": "test"},
			mockData: `{"test":"test"}`,
		},
		{
			name:     "With struct data",
			data:     structData,
			mockData: `{"test":"test"}`,
		},
		{
			name:     "With struct ptr data",
			data:     ptrSructData,
			mockData: `{"test":"test"}`,
		},
	}

	url := "http://localhost:11434/api/generate"

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			mockClient.On("Do",
				mock.MatchedBy(func(r *http.Request) bool {
					body, err := r.GetBody()
					assert.Nil(t, err)

					bd, err := io.ReadAll(body)
					assert.Nil(t, err)

					return r.Method == http.MethodPost &&
						r.URL.String() == url &&
						r.Header.Get("Content-Type") == "application/json" &&
						string(bd) == c.mockData
				}),
			).Return(nil, nil)

			r := client.PostJSON(URL(url), c.data)

			assert.Nil(t, r.err)
			mockClient.AssertExpectations(t)
		})
	}
}

func TestClient_PostJSON_WithEncodingError(t *testing.T) {
	client := NewClient()
	mockClient := &ClientMock{}
	mockEncoder := &EncoderMock{}
	client.client = mockClient
	client.jsonEncoder = mockEncoder
	data := map[string]string{"test": "test"}

	mockEncoder.On("Marshal", data).Return([]byte{}, fmt.Errorf("error marshaling data"))
	mockClient.AssertNotCalled(t, "Post")

	r := client.PostJSON("http://localhost:11434/api/generate", data)

	assert.NotNil(t, r.err)
}

func TestClient_JSONResponse_Parse(t *testing.T) {
	r := io.NopCloser(strings.NewReader(`{"test":"test"}`))
	resp := &http.Response{}
	resp.Body = r

	type Data struct {
		Test string `json:"test"`
	}

	d := &Data{}

	jsonResp := &JSONResponse{Response: resp}

	err := jsonResp.Parse(d)

	assert.Nil(t, err)
}

func TestClient_Get(t *testing.T) {
	client := NewClient()
	mockClient := &ClientMock{}
	client.client = mockClient
	url := "/test-url"

	mockClient.On("Do", mock.MatchedBy(func(r *http.Request) bool {
		return r.Method == http.MethodGet &&
			r.URL.String() == url &&
			r.Body == nil
	})).Return(nil, nil)

	_, err := client.Get(URL(url))
	assert.Nil(t, err)
}

func TestClient_GetJSON(t *testing.T) {
	client := NewClient()
	mockClient := &ClientMock{}
	client.client = mockClient
	url := "/test-url"

	mockClient.On("Do", mock.MatchedBy(func(r *http.Request) bool {
		return r.Method == http.MethodGet &&
			r.URL.String() == url &&
			r.Body == nil
	})).Return(nil, nil)

	resp := client.GetJSON(URL(url))
	assert.Nil(t, resp.err)
}

func TestClient_newRequest(t *testing.T) {
	base := "http://localhost:8080"
	headers := http.Header{"Authorization": {"Bearer test_token"}}
	client := NewClient(BaseURL("http://localhost:8080"), DefaultHeaders(headers))
	url := base + "/test-full-url"
	path := "/test-url"
	method := http.MethodGet

	req, err := client.newRequest(method, url, nil)

	assert.Nil(t, err)
	assert.Equal(t, method, req.Method)
	assert.Equal(t, url, req.URL.String())
	assert.Equal(t, headers, req.Header)

	req, err = client.newRequest(method, path, nil)
	assert.Nil(t, err)
	assert.Equal(t, method, req.Method)
	assert.Equal(t, base+path, req.URL.String())
	assert.Equal(t, headers, req.Header)
}
