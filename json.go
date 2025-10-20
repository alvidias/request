package request

import "encoding/json"

type jsonEncoder struct{}

func (j jsonEncoder) Marshal(v any) ([]byte, error) {
	return json.Marshal(v)
}

type jsonDecoder struct{}

func (j jsonDecoder) Unmarshal(data []byte, v any) error {
	return json.Unmarshal(data, v)
}
