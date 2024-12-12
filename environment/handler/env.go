package handler

import (
	"fmt"

	"gofr.dev/pkg/gofr"
)

type Handler struct {
	envSvc EnvAdder
}

func New(envSvc EnvAdder) *Handler {
	return &Handler{envSvc: envSvc}
}

func (h *Handler) AddEnvironment(ctx *gofr.Context) (any, error) {
	n, err := h.envSvc.AddEnvironments(ctx)
	if err != nil {
		return nil, err
	}

	return fmt.Sprintf("%d enviromnets added", n), nil
}
