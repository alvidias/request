package request

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_IsStructMapOrSlice(t *testing.T) {
	cases := []struct {
		name   string
		data   any
		expect bool
	}{
		{
			name:   "With map data",
			data:   map[string]string{"test": "test"},
			expect: true,
		},
		{
			name: "With struct data",
			data: struct {
				Test string `json:"test"`
			}{Test: "test"},
			expect: true,
		},
		{
			name:   "With slice data",
			data:   []string{"test", "test2"},
			expect: true,
		},
		{
			name:   "With ptr to map data",
			data:   &map[string]string{"test": "test"},
			expect: true,
		},
		{
			name: "With ptr to struct data",
			data: &struct {
				Test string `json:"test"`
			}{Test: "test"},
			expect: true,
		},
		{
			name:   "With ptr to slice data",
			data:   &[]string{"test", "test2"},
			expect: true,
		},
		{
			name:   "With ptr to slice data",
			data:   &[]string{"test", "test2"},
			expect: true,
		},
		{
			name:   "With nil data",
			data:   nil,
			expect: false,
		},
		{
			name:   "With string data",
			data:   "test",
			expect: false,
		},
		{
			name:   "With bytes data",
			data:   []byte("test"),
			expect: false,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			assert.Equal(t, c.expect, IsStructMapOrSlice(c.data))
		})
	}
}
