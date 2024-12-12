package handler

import "gofr.dev/pkg/gofr"

type ApplicationAdder interface {
	AddApplication(ctx *gofr.Context, name string) error
}
