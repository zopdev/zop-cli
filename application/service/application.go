// Package service provides functionalities for managing applications and their environments.
// It includes methods for adding a new application and listing existing applications.
package service

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"gofr.dev/pkg/gofr"
)

// Service provides methods for managing applications.
type Service struct{}

// New creates a new instance of Service.
//
// Returns:
//
//	A pointer to a Service instance.
func New() *Service {
	return &Service{}
}

// Add adds a new application and optionally its environments.
//
// Parameters:
//   - ctx: The application context containing dependencies and utilities.
//   - name: The name of the application to be added.
//
// Returns:
//
//	An error if the application or environments could not be added.
func (*Service) Add(ctx *gofr.Context, name string) error {
	var (
		envs  []Environment
		input string
	)

	app := &Application{Name: name}
	api := ctx.GetHTTPService("api-service")
	level := 1

	ctx.Out.Print("Do you wish to add environments to the application? (y/n) ")

	_, _ = fmt.Scanf("%s", &input)

	for {
		if input != "y" {
			break
		}

		ctx.Out.Print("Enter environment name: ")

		_, _ = fmt.Scanf("%s", &input)
		envs = append(envs, Environment{Name: input, Level: level})
		level++

		ctx.Out.Print("Do you wish to add more? (y/n) ")

		_, _ = fmt.Scanf("%s", &input)

		if input == "n" {
			break
		}
	}

	app.Envs = envs
	body, _ := json.Marshal(app)

	resp, err := api.PostWithHeaders(ctx, "applications", nil, body, map[string]string{
		"Content-Type": "application/json",
	})
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return getAPIError(resp)
	}

	return nil
}

// List retrieves all applications and their environments.
//
// Parameters:
//   - ctx: The application context containing dependencies and utilities.
//
// Returns:
//
//	A slice of applications and an error, if any.
func (*Service) List(ctx *gofr.Context) ([]Application, error) {
	api := ctx.GetHTTPService("api-service")

	reps, err := api.Get(ctx, "applications", nil)
	if err != nil {
		return nil, err
	}
	defer reps.Body.Close()

	var apps struct {
		Data []Application `json:"data"`
	}

	body, _ := io.ReadAll(reps.Body)

	err = json.Unmarshal(body, &apps)
	if err != nil {
		return nil, &ErrAPIService{
			StatusCode: http.StatusInternalServerError,
			Message:    "Internal Server Error",
		}
	}

	return apps.Data, nil
}
