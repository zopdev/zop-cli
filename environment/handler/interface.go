package handler

import "gofr.dev/pkg/gofr"

type EnvAdder interface {
	AddEnvironments(ctx *gofr.Context) (int, error)
}
