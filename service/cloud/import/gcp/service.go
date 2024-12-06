package gcp

import (
	"encoding/json"
	"fmt"
	"gofr.dev/pkg/gofr"
	"gofr.dev/pkg/gofr/cmd/terminal"

	"zop.dev/cli/zop/models"
	cloudImporter "zop.dev/cli/zop/service/cloud/import"
)

var (
	ErrInvalidOrExpiredToken = fmt.Errorf("invalid or expired token, please login again")
)

type Service struct {
	store cloudImporter.AccountStore
}

func New(store cloudImporter.AccountStore) *Service {
	return &Service{
		store: store,
	}
}

func (s *Service) PostAccounts(ctx *gofr.Context) error {
	accounts, err := s.store.GetAccounts(ctx)
	if err != nil {
		return err
	}

	api := ctx.GetHTTPService("api-service")

	defer terminal.NewDotSpinner(ctx.Out).Spin(ctx).Stop()

	for _, acc := range accounts {
		svAccs, er := getServiceAccounts(ctx, acc.Value)
		if er != nil {
			continue
		}

		for _, svcAcc := range svAccs {
			body, er := json.Marshal(&models.PostCloudAccountRequest{
				Name:        acc.AccountID,
				Provider:    "gcp",
				Credentials: svcAcc,
			})
			if er != nil {
				ctx.Logger.Errorf("error marshalling account creds: %v", er)
				continue
			}

			resp, er := api.PostWithHeaders(ctx, "cloud-accounts", nil, body, map[string]string{
				"Content-Type": "application/json",
			})
			if er != nil {
				ctx.Logger.Errorf("error posting account: %v", er)
				continue
			}

			if resp.StatusCode != 201 && resp.StatusCode != 409 {
				ctx.Logger.Errorf("error posting account: %v", resp.Body)
				continue
			}
		}
	}

	return nil
}
