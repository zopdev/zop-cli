package handler

import "gofr.dev/pkg/gofr"

type DeploymentService interface {
	Add(ctx *gofr.Context) error
}
