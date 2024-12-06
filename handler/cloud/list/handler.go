// Package list provides the handler for listing cloud accounts.
package list

import (
	"fmt"

	"gofr.dev/pkg/gofr"
)

type Handler struct {
	accountGetter AccountGetter
}

func New(getter AccountGetter) *Handler {
	return &Handler{
		accountGetter: getter,
	}
}

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
		rows += fmt.Sprintf("%-20s %-20s %-20s %-20s %-20s\n",
			account.Name[:17]+"...", account.Provider, account.ProviderID, account.UpdatedAt, account.CreatedAt)
	}

	return header + rows, nil
}
