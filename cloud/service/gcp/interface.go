package gcp

import (
	"gofr.dev/pkg/gofr"

	"zop.dev/cli/zop/cloud/store/gcp"
)

// AccountStore is an interface for getting accounts from the store layer.
type AccountStore interface {
	GetAccounts(ctx *gofr.Context) ([]gcp.AccountStore, error)
}
