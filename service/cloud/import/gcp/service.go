package gcp

import (
	"encoding/json"
	"gofr.dev/pkg/gofr"
	"gofr.dev/pkg/gofr/cmd/terminal"

	"zop/models"
	cloudImporter "zop/service/cloud/import"
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
		svAcc := getServiceAccount(acc.Value)
		if svAcc == nil {
			continue
		}

		body, er := json.Marshal(&models.PostCloudAccountRequest{
			Name:        acc.AccountID,
			Provider:    "gcp",
			Credentials: svAcc,
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

	return nil
}

func getServiceAccount(value []byte) *models.ServiceAccount {
	var acc models.ServiceAccount
	err := json.Unmarshal(value, &acc)
	if err != nil {
		return nil
	}

	if acc.PrivateKey == "" {
		return nil
	}

	return &acc
}
