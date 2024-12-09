package list

import (
	"encoding/json"
	"io"

	"gofr.dev/pkg/gofr"
)

type Service struct {
}

func New() *Service {
	return &Service{}
}

func (*Service) GetAccounts(ctx *gofr.Context) ([]*CloudAccountResponse, error) {
	api := ctx.GetHTTPService("api-service")

	reps, err := api.Get(ctx, "/cloud-accounts", nil)
	if err != nil {
		return nil, err
	}
	defer reps.Body.Close()

	var accounts struct {
		Data []*CloudAccountResponse `json:"data"`
	}

	body, _ := io.ReadAll(reps.Body)

	err = json.Unmarshal(body, &accounts)
	if err != nil {
		return nil, err
	}

	return accounts.Data, nil
}
