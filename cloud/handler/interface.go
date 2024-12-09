package handler

import (
	"gofr.dev/pkg/gofr"
	"zop.dev/cli/zop/cloud/service/list"
)

// AccountImporter is an interface for importing cloud accounts to zop api.
// It has a PostAccounts method that is used to import all local cloud accounts to the zop api to store and validate those cloud accounts.
type AccountImporter interface {
	PostAccounts(ctx *gofr.Context) error
}

// AccountGetter is an interface for getting cloud accounts from the zop api.
type AccountGetter interface {
	GetAccounts(ctx *gofr.Context) ([]*list.CloudAccountResponse, error)
}
