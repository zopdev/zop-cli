package service

import (
	"gofr.dev/pkg/gofr"
	appSvc "zop.dev/cli/zop/application/service"
)

type ApplicationGetter interface {
	GetApplications(ctx *gofr.Context) ([]appSvc.Application, error)
}
