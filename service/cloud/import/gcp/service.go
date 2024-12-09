// Package gcp provides a service for importing GCP service accounts into zop api service.
package gcp

import (
	"encoding/json"
	"fmt"
	"net/http"

	"gofr.dev/pkg/gofr"
	"gofr.dev/pkg/gofr/cmd/terminal"
	"zop.dev/cli/zop/models"
	cloudImporter "zop.dev/cli/zop/service/cloud/import"
)

var (
	// ErrInvalidOrExpiredToken is returned when the token is invalid or expired for the gcloud user account
	// and a new token cannot be generated. User is advised to run gcloud auth login to refresh the token.
	ErrInvalidOrExpiredToken = fmt.Errorf("invalid or expired token, please login again")
)

// Service is a service for importing GCP service accounts into zop api service.
type Service struct {
	store cloudImporter.AccountStore
}

func New(store cloudImporter.AccountStore) *Service {
	return &Service{
		store: store,
	}
}

type ErrAPIService struct {
	StatusCode int
	Message    string
}

func (e *ErrAPIService) Error() string {
	return fmt.Sprintf("error from api service: %s, status code: %d", e.Message, e.StatusCode)
}

// PostAccounts posts the GCP service accounts to the api service.
// It fetches the accounts from the store layer and posts them to the api service.
// If an account is of type user account, it generates a token and then creates a service account.
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
			ctx.Logger.Errorf("error getting service accounts: %v", er)

			continue
		}

		for _, svcAcc := range svAccs {
			body, er := json.Marshal(&models.PostCloudAccountRequest{
				Name:        acc.AccountID,
				Provider:    "gcp",
				Credentials: svcAcc,
			})
			if er != nil {
				ctx.Logger.Errorf("error marshaling account creds: %v", er)
				continue
			}

			resp, er := api.PostWithHeaders(ctx, "cloud-accounts", nil, body, map[string]string{
				"Content-Type": "application/json",
			})
			if er != nil {
				ctx.Logger.Errorf("error posting account: %v", er)
				continue
			}

			if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusConflict {
				ctx.Logger.Errorf("error posting account: %v", resp.Body)

				return &ErrAPIService{StatusCode: resp.StatusCode, Message: "could not connect to the zop-api service"}
			}

			resp.Body.Close()
		}
	}

	return nil
}
