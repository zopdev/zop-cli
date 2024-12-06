package gcp

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"
	"zop.dev/cli/zop/models"

	"gofr.dev/pkg/gofr"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/cloudresourcemanager/v1"
	"google.golang.org/api/iam/v1"
	"google.golang.org/api/option"
)

const (
	tokenURL = "https://oauth2.googleapis.com/token"
)

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
		config := newServiceAccountConfig(projectID, serviceAccountName)

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

type serviceAccountConfig struct {
	ProjectID          string
	ServiceAccountName string
	Roles              []string
}

func newServiceAccountConfig(projectID, serviceAccountName string) *serviceAccountConfig {
	return &serviceAccountConfig{
		ProjectID:          projectID,
		ServiceAccountName: serviceAccountName,
		Roles: []string{
			"roles/editor",
			"roles/container.admin",
			"roles/resourcemanager.projectIamAdmin",
			"roles/iam.roleAdmin",
			"roles/secretmanager.admin",
			"roles/servicenetworking.networksAdmin",
			"roles/storage.admin",
			"roles/dns.admin",
			"roles/artifactregistry.admin",
			"roles/pubsub.admin",
		},
	}
}

func checkProjectAccess(ctx context.Context, projectID string, accessToken *oauth2.Token) error {
	crmService, err := cloudresourcemanager.NewService(ctx, option.WithTokenSource(oauth2.StaticTokenSource(accessToken)))
	if err != nil {
		return fmt.Errorf("failed to create Cloud Resource Manager client: %v", err)
	}

	_, err = crmService.Projects.Get(projectID).Do()
	if err != nil {
		return fmt.Errorf("project %s is not accessible: %v", projectID, err)
	}

	return nil
}

func createServiceAccount(ctx context.Context, config *serviceAccountConfig) (*iam.ServiceAccount, error) {
	iamService, err := iam.NewService(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to create IAM client: %v", err)
	}

	serviceAccountEmail := fmt.Sprintf("%s@%s.iam.gserviceaccount.com",
		config.ServiceAccountName, config.ProjectID)

	// Check if service account exists
	_, err = iamService.Projects.ServiceAccounts.Get(
		fmt.Sprintf("projects/%s/serviceAccounts/%s",
			config.ProjectID, serviceAccountEmail)).Do()

	if err == nil {
		return nil, fmt.Errorf("service account %s already exists", serviceAccountEmail)
	}

	request := &iam.CreateServiceAccountRequest{
		AccountId: config.ServiceAccountName,
		ServiceAccount: &iam.ServiceAccount{
			DisplayName: config.ServiceAccountName,
			Description: "Service account for ZOP",
		},
	}

	serviceAccount, err := iamService.Projects.ServiceAccounts.Create(
		fmt.Sprintf("projects/%s", config.ProjectID), request).Do()
	if err != nil {
		return nil, fmt.Errorf("failed to create service account: %v", err)
	}

	return serviceAccount, nil
}

func createServiceAccountKey(ctx context.Context, serviceAccount *iam.ServiceAccount) ([]byte, error) {
	iamService, err := iam.NewService(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to create IAM client: %v", err)
	}

	key, err := iamService.Projects.ServiceAccounts.Keys.Create(
		serviceAccount.Name,
		&iam.CreateServiceAccountKeyRequest{
			PrivateKeyType: "TYPE_GOOGLE_CREDENTIALS_FILE",
		}).Do()
	if err != nil {
		return nil, fmt.Errorf("failed to create service account key: %v", err)
	}

	return []byte(key.PrivateKeyData), nil
}

func assignRoles(ctx context.Context, config *serviceAccountConfig, serviceAccount *iam.ServiceAccount) error {
	crmService, err := cloudresourcemanager.NewService(ctx)
	if err != nil {
		return fmt.Errorf("failed to create Cloud Resource Manager client: %v", err)
	}

	policy, err := crmService.Projects.GetIamPolicy(config.ProjectID,
		&cloudresourcemanager.GetIamPolicyRequest{}).Do()
	if err != nil {
		return fmt.Errorf("failed to get IAM policy: %v", err)
	}

	member := fmt.Sprintf("serviceAccount:%s", serviceAccount.Email)

	for _, role := range config.Roles {
		found := false
		for _, binding := range policy.Bindings {
			if binding.Role == role {
				binding.Members = append(binding.Members, member)
				found = true
				break
			}
		}
		if !found {
			policy.Bindings = append(policy.Bindings, &cloudresourcemanager.Binding{
				Role:    role,
				Members: []string{member},
			})
		}
	}

	_, err = crmService.Projects.SetIamPolicy(config.ProjectID,
		&cloudresourcemanager.SetIamPolicyRequest{Policy: policy}).Do()
	if err != nil {
		return fmt.Errorf("failed to set IAM policy: %v", err)
	}

	return nil
}

func refreshAccessToken(ctx *gofr.Context, clientID, clientSecret, refreshToken string) (*oauth2.Token, error) {
	// Prepare the token URL and data
	data := url.Values{
		"client_id":     {clientID},
		"client_secret": {clientSecret},
		"refresh_token": {refreshToken},
		"grant_type":    {"refresh_token"},
	}

	// Make the request to get the new access token
	resp, err := http.PostForm(tokenURL, data)
	if err != nil {
		return nil, fmt.Errorf("failed to refresh token: %v", err)
	}
	defer resp.Body.Close()

	// Parse the response and extract the new token
	var newToken oauth2.Token
	err = json.NewDecoder(resp.Body).Decode(&newToken)
	if err != nil {
		return nil, fmt.Errorf("failed to decode token response: %v", err)
	}

	return &newToken, nil
}

func fetchProjects(ctx *gofr.Context, clientID, clientSecret string, token *oauth2.Token) ([]*cloudresourcemanager.Project, error) {
	config := &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		Scopes:       []string{"https://www.googleapis.com/auth/cloud-platform.read-only"},
		Endpoint:     google.Endpoint,
	}

	client := config.Client(ctx, token)

	// Use the Google Cloud Resource Manager API to fetch projects
	cloudService, err := cloudresourcemanager.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		ctx.Errorf("Failed to create Cloud Resource Manager client: %v", err)

		return nil, err
	}

	resp, err := cloudService.Projects.List().Do()
	if err != nil {
		ctx.Errorf("Failed to list projects: %v", err)

		return nil, err
	}

	return resp.Projects, nil
}
