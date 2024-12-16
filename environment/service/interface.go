package service

import (
	"gofr.dev/pkg/gofr"
	appSvc "zop.dev/cli/zop/application/service"
)

type ApplicationGetter interface {
	List(ctx *gofr.Context) ([]appSvc.Application, error)
}
