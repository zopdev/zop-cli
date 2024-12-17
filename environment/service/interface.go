package service

import (
	"gofr.dev/pkg/gofr"
	appSvc "zop.dev/cli/zop/application/service"
)

// ApplicationGetter interface is used to abstract the process of fetching application data,
// which can be implemented by any service that has access to application-related data.
type ApplicationGetter interface {
	// List fetches a list of applications from the service.
	// It returns a slice of Application objects and an error if the request fails.
	List(ctx *gofr.Context) ([]appSvc.Application, error)
}
