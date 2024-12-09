package list

type CloudAccountResponse struct {
	// Name is the name of the cloud account.
	Name string `json:"name"`

	// ID is a unique identifier for the cloud account.
	ID int64 `json:"id,omitempty"`

	// Provider is the name of the cloud service provider.
	Provider string `json:"provider"`

	// ProviderID is the identifier for the provider account.
	ProviderID string `json:"providerId"`

	// ProviderDetails contains additional details specific to the provider.
	ProviderDetails any `json:"providerDetails"`

	// CreatedAt is the timestamp of when the cloud account was created.
	CreatedAt string `json:"createdAt"`

	// UpdatedAt is the timestamp of the last update to the cloud account.
	UpdatedAt string `json:"updatedAt"`

	// DeletedAt is the timestamp of when the cloud account was deleted, if applicable.
	DeletedAt string `json:"deletedAt,omitempty"`
}
