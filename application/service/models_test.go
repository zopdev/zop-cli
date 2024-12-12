package service

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetAPIError(t *testing.T) {
	tests := []struct {
		name         string
		responseBody interface{}
		responseCode int
		expectError  *ErrAPIService
	}{
		{
			name: "Valid error response",
			responseBody: MockErrorResponse{
				Error: "Something went wrong",
			},
			responseCode: http.StatusBadRequest,
			expectError: &ErrAPIService{
				StatusCode: http.StatusBadRequest,
				Message:    "Something went wrong",
			},
		},
		{
			name:         "Malformed JSON response",
			responseBody: "{error: unquoted string}",
			responseCode: http.StatusBadRequest,
			expectError:  errInternal,
		},
		{
			name:         "Empty response body",
			responseBody: "",
			responseCode: http.StatusBadRequest,
			expectError:  errInternal,
		},
		{
			name:         "Error reading body",
			responseBody: nil,
			responseCode: http.StatusInternalServerError,
			expectError:  errInternal,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var bodyReader io.ReadCloser

			if tt.responseBody == nil {
				// Simulate error reading body by providing a reader that always errors
				bodyReader = io.NopCloser(&errorReader{})
			} else {
				b, err := json.Marshal(tt.responseBody)
				if err != nil {
					t.Fatalf("Failed to marshal test response body: %v", err)
				}

				bodyReader = io.NopCloser(bytes.NewBuffer(b))
			}

			// Create a mock response
			resp := &http.Response{
				StatusCode: tt.responseCode,
				Body:       bodyReader,
			}

			// Call the function and check results
			actualError := getAPIError(resp)

			require.Equal(t, tt.expectError.StatusCode, actualError.StatusCode, "Unexpected status code")
			require.Equal(t, tt.expectError.Message, actualError.Message, "Unexpected error message")
		})
	}
}

type MockErrorResponse struct {
	Error string `json:"error"`
}

type errorReader struct{}

func (*errorReader) Read(_ []byte) (n int, err error) {
	return 0, io.ErrUnexpectedEOF
}

func (*errorReader) Close() error {
	return nil
}

func Test_ErrAPIError(t *testing.T) {
	tests := []struct {
		name string
		err  *ErrAPIService
	}{
		{
			name: "Valid error",
			err: &ErrAPIService{
				StatusCode: http.StatusBadRequest,
				Message:    "Something went wrong",
			},
		},
		{
			name: "Empty error message",
			err: &ErrAPIService{
				StatusCode: http.StatusBadRequest,
				Message:    "",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := tt.err.Error()
			require.Equal(t, tt.err.Message, actual, "Unexpected error message")
		})
	}
}
