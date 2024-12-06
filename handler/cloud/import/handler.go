// Package _import is used to import data from external sources
// this package has an Import(ctx *gofr.Context) method that is used to import all
// local cloud accounts to the zop api to store and validate those cloud accounts.
package _import

import (
	"gofr.dev/pkg/gofr"
)

type Handler struct {
	accountService AccountImporter
}

func New(getter AccountImporter) *Handler {
	return &Handler{
		accountService: getter,
	}
}

func (h *Handler) Import(ctx *gofr.Context) (any, error) {
	err := h.accountService.PostAccounts(ctx)
	if err != nil {
		return nil, err
	}

	return "Successfully Imported!", nil
}
