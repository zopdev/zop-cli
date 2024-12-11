package service

import (
	"encoding/json"
	"io"
	"net/http"
)

type ErrAPIService struct {
	StatusCode int
	Message    string
}

func (e *ErrAPIService) Error() string {
	return e.Message
}

var errInternal = &ErrAPIService{
	StatusCode: http.StatusInternalServerError,
	Message:    "error in POST /application zop-api, invalid response",
}

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

type Environment struct {
	Name            string `json:"name"`
	Order           int    `json:"order"`
	DeploymentSpace any    `json:"deploymentSpace,omitempty"`
}

type Application struct {
	Name string        `json:"name"`
	Envs []Environment `json:"environments,omitempty"`
}
