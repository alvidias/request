package request

import (
	"net/http"

	"github.com/stretchr/testify/mock"
)

type ClientMock struct {
	mock.Mock
}

func (m *ClientMock) Do(req *http.Request) (*http.Response, error) {
	args := m.Called(req)

	respIface := args.Get(0)
	var resp *http.Response

	if respIface != nil {
		resp = respIface.(*http.Response)
	}

	return resp, args.Error(1)
}

func (m *ClientMock) PostJSON(url string, data any) (*http.Response, error) {
	args := m.Called(url, data)

	respIface := args.Get(0)
	var resp *http.Response

	if respIface != nil {
		resp = respIface.(*http.Response)
	}

	return resp, args.Error(1)
}

type EncoderMock struct {
	mock.Mock
}

func (m *EncoderMock) Marshal(data any) ([]byte, error) {
	args := m.Called(data)

	return args.Get(0).([]byte), args.Error(1)
}

func (m *EncoderMock) Unmarshal(data []byte, v any) error {
	args := m.Called(data, v)

	return args.Error(0)
}
