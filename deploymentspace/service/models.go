package service

// DeploymentSpaceOptions represents the deployment space options in the system.
//
// It includes information such as the name, path, and type of the deployment space.
type DeploymentSpaceOptions struct {
	Name string `json:"name"` // The name of the deployment space.
	Path string `json:"path"` // The API path to access the deployment space.
	Type string `json:"type"` // The type of the deployment space.
}

// DeploymentOption represents a list of deployment options and information about the next page, if available.
type DeploymentOption struct {
	Option []map[string]any `json:"options"` // A slice of options for deployment.
	Next   *Next            `json:"next"`    // Information about the next page of options.
}

// Next provides details about the subsequent page of deployment options.
//
// It includes the name, path, and query parameters for accessing the next page.
type Next struct {
	Name   string            `json:"name"`   // The name of the next page.
	Path   string            `json:"path"`   // The API path to access the next page.
	Params map[string]string `json:"params"` // Query parameters required for the next page.
}

// apiResponse represents the structure of the API response for deployment options.
//
// It contains a data field with deployment options.
type apiResponse struct {
	Data *DeploymentOption `json:"data"` // The deployment option data.
}
