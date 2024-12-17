package service

import (
	"gofr.dev/pkg/gofr"

	appSvc "zop.dev/cli/zop/application/service"
	cloudSvc "zop.dev/cli/zop/cloud/service/list"
	envSvc "zop.dev/cli/zop/environment/service"
)

type CloudAccountService interface {
	GetAccounts(ctx *gofr.Context) ([]*cloudSvc.CloudAccountResponse, error)
}

type ApplicationService interface {
	List(ctx *gofr.Context) ([]appSvc.Application, error)
}

type EnvironmentService interface {
	List(ctx *gofr.Context) ([]envSvc.Environment, error)
}
