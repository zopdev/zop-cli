package utils

import (
	"encoding/json"
	"io"
	"net/http"
)

// GetResponse reads the HTTP response body and unmarshals it into the provided interface.
func GetResponse(resp *http.Response, i any) error {
	defer resp.Body.Close()

	b, _ := io.ReadAll(resp.Body)

	err := json.Unmarshal(b, i)
	if err != nil {
		return err
	}

	return nil
}
