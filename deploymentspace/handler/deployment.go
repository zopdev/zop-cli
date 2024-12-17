package handler

import "gofr.dev/pkg/gofr"

type Handler struct {
	deployService DeploymentService
}

func New(depSvc DeploymentService) *Handler {
	return &Handler{
		deployService: depSvc,
	}
}

func (h *Handler) Add(ctx *gofr.Context) (any, error) {
	err := h.deployService.Add(ctx)
	if err != nil {
		return nil, err
	}

	return "Deployment Created", nil
}
