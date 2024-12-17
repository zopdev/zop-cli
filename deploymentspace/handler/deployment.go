package handler

import "gofr.dev/pkg/gofr"

type Handler struct{}

func New() *Handler {
	return &Handler{}
}

func (h *Handler) Add(ctx *gofr.Context) (any, error) {
	// Add your code here
	return nil, nil
}
