package _import

import "gofr.dev/pkg/gofr"

type AccountImporter interface {
	PostAccounts(ctx *gofr.Context) error
}
