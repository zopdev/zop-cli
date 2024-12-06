package _import

import (
	"gofr.dev/pkg/gofr"

	"zop.dev/cli/zop/models"
)

type AccountStore interface {
	GetAccounts(ctx *gofr.Context) ([]models.AccountStore, error)
}
