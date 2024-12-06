package gcp

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"time"

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

func getServiceAccounts(ctx *gofr.Context, value []byte) ([]*models.ServiceAccount, error) {
	var acc models.ServiceAccount

	err := json.Unmarshal(value, &acc)
	if err != nil {
		return nil, err
	}

	if acc.PrivateKey == "" {
		return generateNewServiceAccount(ctx, value)
	}

	return []*models.ServiceAccount{&acc}, nil
}

func generateNewServiceAccount(ctx *gofr.Context, value []byte) ([]*models.ServiceAccount, error) {
	var acc models.UserAccount

	err := json.Unmarshal(value, &acc)
	if err != nil {
		return nil, err
	}

	token, err := refreshAccessToken(ctx, acc.ClientID, acc.ClientSecret, acc.RefreshToken)
	if err != nil {
		return nil, ErrInvalidOrExpiredToken
	}

	projects, err := fetchProjects(ctx, acc.ClientID, acc.ClientSecret, token)
	if err != nil {
		return nil, err
	}

	var serviceAccounts []*models.ServiceAccount

	for _, project := range projects {
		projectID := project.ProjectId
		serviceAccountName := fmt.Sprintf("zop-dev-%v", time.Now().Unix())
		config := NewServiceAccountConfig(projectID, serviceAccountName)

		if err = checkProjectAccess(ctx, config.ProjectID, token); err != nil {
			ctx.Logger.Errorf("Project access check failed: %v", err)
			continue
		}

		serviceAccount, err := createServiceAccount(ctx, config)
		if err != nil {
			ctx.Logger.Errorf("Failed to create service account: %v", err)
			continue
		}

		key, err := createServiceAccountKey(ctx, serviceAccount)
		if err != nil {
			ctx.Logger.Errorf("Failed to create service account key: %v", err)
			continue
		}

		decodedKey, err := base64.StdEncoding.DecodeString(string(key))
		if err != nil {
			ctx.Errorf("Failed to decode Base64 string: %v", err)
			continue
		}

		if err = assignRoles(ctx, config, serviceAccount); err != nil {
			ctx.Logger.Errorf("Failed to assign roles: %v", err)
			continue
		}

		var svAcc models.ServiceAccount

		err = json.Unmarshal(decodedKey, &svAcc)
		if err != nil {
			ctx.Logger.Errorf("Failed to unmarshal service account key: %v", err)
			continue
		}

		serviceAccounts = append(serviceAccounts, &svAcc)
	}

	return serviceAccounts, nil
}
