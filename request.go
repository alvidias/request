package request

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func ReadJson(r *http.Request, data any) error {
	c, err := io.ReadAll(r.Body)

	if err != nil {
		return fmt.Errorf("error reading request's body: %w", err)
	}

	defer r.Body.Close()

	err = json.Unmarshal(c, &data)

	if err != nil {
		return fmt.Errorf("error unmarshalling JSON: %w", err)
	}

	return nil
}
