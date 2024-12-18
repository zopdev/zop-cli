// Package service defines models and utility functions for error handling
// and data structures representing applications and their environments.
package service

import (
	"encoding/json"
	"io"
	"net/http"
)

// ErrAPIService represents an error returned by the API service.
type ErrAPIService struct {
	StatusCode int    // HTTP status code of the error
	Message    string // Message describing the error
}

// Error returns the error message for ErrAPIService.
//
// Returns:
//
//	A string describing the error.
func (e *ErrAPIService) Error() string {
	return e.Message
}

// Predefined internal error for API response issues.
var errInternal = &ErrAPIService{
	StatusCode: http.StatusInternalServerError,
	Message:    "error in /applications zop-api, invalid response",
}

// getAPIError extracts and constructs an ErrAPIService from an HTTP response.
//
// Parameters:
//   - resp: The HTTP response containing the error details.
//
// Returns:
//
//	A pointer to an ErrAPIService.
func getAPIError(resp *http.Response) *ErrAPIService {
	var errResp struct {
		Error string `json:"error"`
	}

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return errInternal
	}

	err = json.Unmarshal(b, &errResp)
	if err != nil {
		return errInternal
	}

	return &ErrAPIService{
		StatusCode: resp.StatusCode,
		Message:    errResp.Error,
	}
}

// Environment represents an environment associated with an application.
type Environment struct {
	Name            string `json:"name"`                      // Name of the environment
	Level           int    `json:"level"`                     // Priority Level of the environment
	DeploymentSpace any    `json:"deploymentSpace,omitempty"` // DeploymentSpace information for the environment
}

// Application represents an application with its associated environments.
type Application struct {
	ID   int64         `json:"id"`                     // Unique identifier of the application
	Name string        `json:"name"`                   // Name of the application
	Envs []Environment `json:"environments,omitempty"` // List of associated environments
}
