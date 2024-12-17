// Package handler provides the CMD handler logic for managing environments.
package handler

import (
	"fmt"

	"gofr.dev/pkg/gofr"
)

// Handler is responsible for managing environment-related operations.
type Handler struct {
	envSvc EnvAdder
}

// New creates a new Handler with the given EnvAdder service.
func New(envSvc EnvAdder) *Handler {
	return &Handler{envSvc: envSvc}
}

// Add handles the HTTP request to add environments. It delegates the task
// to the EnvAdder service and returns a success message or an error.
//
// Parameters:
//   - ctx: The GoFR context containing request data.
//
// Returns:
//   - A success message indicating how many environments were added, or an error
//     if the operation failed.
func (h *Handler) Add(ctx *gofr.Context) (any, error) {
	n, err := h.envSvc.Add(ctx)
	if err != nil {
		return nil, err
	}

	return fmt.Sprintf("%d environments added", n), nil
}
