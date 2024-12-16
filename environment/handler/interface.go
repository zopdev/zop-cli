package handler

import "gofr.dev/pkg/gofr"

type EnvAdder interface {
	Add(ctx *gofr.Context) (int, error)
}
