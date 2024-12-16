package handler

import (
	"gofr.dev/pkg/gofr"

	"zop.dev/cli/zop/application/service"
)

type ApplicationService interface {
	AddApplication(ctx *gofr.Context, name string) error
	GetApplications(ctx *gofr.Context) ([]service.Application, error)
}
