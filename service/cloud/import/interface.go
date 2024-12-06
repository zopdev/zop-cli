package _import

import (
	"gofr.dev/pkg/gofr"

	"zop/models"
)

type AccountStore interface {
	GetAccounts(ctx *gofr.Context) ([]models.AccountStore, error)
}
