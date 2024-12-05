package _import

import (
	"gofr.dev/pkg/gofr"
)

type Handler struct {
	accountStore AccountGetter
}

func New(getter AccountGetter) *Handler {
	return &Handler{
		accountStore: getter,
	}
}

func (h *Handler) Import(ctx *gofr.Context) (any, error) {
	acc, err := h.accountStore.GetAccounts(ctx)
	if err != nil {
		return nil, err
	}

	return acc, nil
}
