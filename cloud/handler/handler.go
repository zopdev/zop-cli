// Package handler is used to import data from external sources
// this package has an Import(ctx *gofr.Context) method that is used to import all
// local cloud accounts to the zop api to store and validate those cloud accounts.
package handler

import (
	"fmt"

	"gofr.dev/pkg/gofr"
)

const (
	successMessage = "Successfully Imported!"
	maxNameLength  = 20
)

type Handler struct {
	accountService AccountImporter
	accountGetter  AccountGetter
}

func New(accountService AccountImporter, accountGetter AccountGetter) *Handler {
	return &Handler{
		accountService: accountService,
		accountGetter:  accountGetter,
	}
}

// Import is a handler for importing cloud accounts to zop api.
func (h *Handler) Import(ctx *gofr.Context) (any, error) {
	err := h.accountService.PostAccounts(ctx)
	if err != nil {
		return nil, err
	}

	return successMessage, nil
}

// List is a handler for listing all cloud accounts.
func (h *Handler) List(ctx *gofr.Context) (any, error) {
	accounts, err := h.accountGetter.GetAccounts(ctx)
	if err != nil {
		return nil, err
	}

	if len(accounts) == 0 {
		return "No accounts found", nil
	}

	header := fmt.Sprintf("%-20s %-20s %-20s %-20s %-20s\n",
		"Name", "Provider", "ProviderID", "UpdateAt", "CreatedAt")
	rows := ""

	for _, account := range accounts {
		if len(account.Name) > maxNameLength {
			account.Name = account.Name[:17] + "..."
		}

		rows += fmt.Sprintf("%-20s %-20s %-20s %-20s %-20s\n",
			account.Name, account.Provider, account.ProviderID, account.UpdatedAt, account.CreatedAt)
	}

	return header + rows, nil
}
