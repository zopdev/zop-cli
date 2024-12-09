package export

import "gofr.dev/pkg/gofr"

// AccountImporter is an interface for importing cloud accounts to zop api.
// It has a PostAccounts method that is used to import all local cloud accounts to the zop api to store and validate those cloud accounts.
type AccountImporter interface {
	PostAccounts(ctx *gofr.Context) error
}
