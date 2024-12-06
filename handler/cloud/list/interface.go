package list

import (
	"gofr.dev/pkg/gofr"
	"zop/models"
)

type AccountGetter interface {
	GetAccounts(ctx *gofr.Context) ([]models.CloudAccountResponse, error)
}
