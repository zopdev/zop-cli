package gcp

// AccountStore stores the accountId and creds value for gcp account.
type AccountStore struct {
	AccountID string `json:"account_id"`
	Value     []byte `json:"value"`
}
