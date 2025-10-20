package request

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMarshal(t *testing.T) {
	data := map[string]string{"test": "test"}
	encoder := jsonEncoder{}

	e, err := encoder.Marshal(data)
	assert.Nil(t, err)
	assert.Equal(t, `{"test":"test"}`, string(e))
}

func TestUnmarshal(t *testing.T) {
	data := `{"test":"test"}`
	var result map[string]string
	encoder := jsonDecoder{}

	err := encoder.Unmarshal([]byte(data), &result)
	assert.Nil(t, err)
	assert.Equal(t, map[string]string{"test": "test"}, result)
}
