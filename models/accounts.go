package models

type AccountStore struct {
	AccountID string `json:"account_id"`
	Value     []byte `json:"value"`
}

type UserAccount struct {
	RefreshToken string `json:"refresh_token"`
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
}

type ServiceAccount struct {
	Type                    string `json:"type"`
	ProjectID               string `json:"project_id"`
	PrivateKeyID            string `json:"private_key_id"`
	PrivateKey              string `json:"private_key"`
	ClientEmail             string `json:"client_email"`
	ClientID                string `json:"client_id"`
	AuthURI                 string `json:"auth_uri"`
	TokenURI                string `json:"token_uri"`
	AuthProviderX509CertURL string `json:"auth_provider_x509_cert_url"`
	ClientX509CertURL       string `json:"client_x509_cert_url"`
	UniverseDomain          string `json:"universe_domain"`
}

type CloudAccountResponse struct {
	Name            string `json:"name"`
	Provider        string `json:"provider"`
	ID              int64  `json:"id,omitempty"`
	ProviderID      string `json:"providerId"`
	ProviderDetails any    `json:"providerDetails"`
	Credentials     any    `json:"credentials,omitempty"`
	CreatedAt       string `json:"createdAt"`
	UpdatedAt       string `json:"updatedAt"`
	DeletedAt       string `json:"deletedAt,omitempty"`
}

type PostCloudAccountRequest struct {
	Name        string `json:"name"`
	Provider    string `json:"provider"`
	Credentials any    `json:"credentials"`
}
