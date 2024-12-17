// Package handler provides the CMD handler logic for managing environments.
package handler

import (
	"bytes"
	"fmt"
	"sort"
	"text/tabwriter"

	"gofr.dev/pkg/gofr"
)

const padding = 2

// Handler is responsible for managing environment-related operations.
type Handler struct {
	envSvc EnvironmentService
}

// New creates a new Handler with the given EnvAdder service.
func New(envSvc EnvironmentService) *Handler {
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

func (h *Handler) List(ctx *gofr.Context) (any, error) {
	envs, err := h.envSvc.List(ctx)
	if err != nil {
		return nil, err
	}

	sort.Slice(envs, func(i, j int) bool { return envs[i].ID < envs[j].ID })

	b := bytes.NewBuffer([]byte{})

	// Print a table of all the environments in the application
	writer := tabwriter.NewWriter(b, 0, 0, padding, ' ', tabwriter.Debug)

	// Print table headers
	fmt.Fprintln(writer, "Name\tLevel\tCreatedAt\tUpdatedAt")

	// Print rows for each environment
	for _, env := range envs {
		fmt.Fprintf(writer, "%s\t%d\t%s\t%s\n",
			env.Name,
			env.Level,
			env.CreatedAt,
			env.UpdatedAt,
		)
	}

	// Flush the writer to output the table
	writer.Flush()

	return b.String(), nil
}
