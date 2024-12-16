package handler

import (
	"gofr.dev/pkg/gofr"

	"zop.dev/cli/zop/application/service"
)

type ApplicationService interface {
	Add(ctx *gofr.Context, name string) error
	List(ctx *gofr.Context) ([]service.Application, error)
}
