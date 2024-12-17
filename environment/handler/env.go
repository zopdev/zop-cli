package handler

import (
	"bytes"
	"fmt"
	"sort"
	"text/tabwriter"

	"gofr.dev/pkg/gofr"
)

const padding = 2

type Handler struct {
	envSvc EnvironmentService
}

func New(envSvc EnvironmentService) *Handler {
	return &Handler{envSvc: envSvc}
}

func (h *Handler) Add(ctx *gofr.Context) (any, error) {
	n, err := h.envSvc.Add(ctx)
	if err != nil {
		return nil, err
	}

	return fmt.Sprintf("%d enviromnets added", n), nil
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
