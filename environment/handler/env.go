package handler

import (
	"fmt"
	"sort"

	"gofr.dev/pkg/gofr"
)

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

	// Print a table of all the environments in the application
	ctx.Out.Println("ID\tName")

	for _, env := range envs {
		ctx.Out.Printf("%d\t%s\n", env.ID, env.Name)
	}

	return nil, nil
}
