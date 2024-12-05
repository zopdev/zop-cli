package _import

import "gofr.dev/pkg/gofr"

type AccountGetter interface {
	GetAccounts(ctx *gofr.Context) ([]string, error)
}
