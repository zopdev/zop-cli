package export

import (
	"gofr.dev/pkg/gofr"

	"zop.dev/cli/zop/models"
)

// AccountStore is an interface for getting accounts from the store layer.
type AccountStore interface {
	GetAccounts(ctx *gofr.Context) ([]models.AccountStore, error)
}
