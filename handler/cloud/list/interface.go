package list

import (
	"gofr.dev/pkg/gofr"

	"zop.dev/cli/zop/models"
)

// AccountGetter is an interface for getting cloud accounts from the zop api.
type AccountGetter interface {
	GetAccounts(ctx *gofr.Context) ([]models.CloudAccountResponse, error)
}
