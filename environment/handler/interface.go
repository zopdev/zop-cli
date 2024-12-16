package handler

import (
	"gofr.dev/pkg/gofr"

	"zop.dev/cli/zop/environment/service"
)

type EnvironmentService interface {
	Add(ctx *gofr.Context) (int, error)
	List(ctx *gofr.Context) ([]service.Environment, error)
}
