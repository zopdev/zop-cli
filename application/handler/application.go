package handler

import (
	"errors"

	"gofr.dev/pkg/gofr"
)

var (
	ErrorApplicationNameNotProvided = errors.New("please enter application name, -name=<application_name>")
)

type Handler struct {
	appAdd ApplicationAdder
}

func New(appAdd ApplicationAdder) *Handler {
	return &Handler{
		appAdd: appAdd,
	}
}

func (h *Handler) Add(ctx *gofr.Context) (any, error) {
	name := ctx.Param("name")
	if name == "" {
		return nil, ErrorApplicationNameNotProvided
	}

	err := h.appAdd.AddApplication(ctx, name)
	if err != nil {
		return nil, err
	}

	return "Application " + name + " added successfully!", nil
}
