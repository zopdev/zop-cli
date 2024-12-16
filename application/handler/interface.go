package handler

import (
	"gofr.dev/pkg/gofr"

	"zop.dev/cli/zop/application/service"
)

// ApplicationService defines the methods required for application management.
type ApplicationService interface {
	// Add adds a new application with the specified name.
	//
	// Parameters:
	//   - ctx: The application context containing dependencies and utilities.
	//   - name: The name of the application to be added.
	//
	// Returns:
	//   An error if the application could not be added.
	Add(ctx *gofr.Context, name string) error

	// List retrieves the list of applications along with their environments.
	//
	// Parameters:
	//   - ctx: The application context containing dependencies and utilities.
	//
	// Returns:
	//   A slice of applications and an error, if any.
	List(ctx *gofr.Context) ([]service.Application, error)
}
