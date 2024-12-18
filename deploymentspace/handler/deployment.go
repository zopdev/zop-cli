// Package handler is used to import data from external sources
// this package has an Add(ctx *gofr.Context) method that is used to configure deployment space
// for the environments of applications.
package handler

import "gofr.dev/pkg/gofr"

// Handler is responsible for handling requests related to deployment operations.
type Handler struct {
	deployService DeploymentService
}

// New initializes a new Handler instance.
//
// Parameters:
//   - depSvc: An instance of DeploymentService used to manage deployment-related operations.
//
// Returns:
//   - A pointer to the Handler instance.
func New(depSvc DeploymentService) *Handler {
	return &Handler{
		deployService: depSvc,
	}
}

// Add processes a deployment creation request.
//
// Parameters:
//   - ctx: The context object containing request and session details.
//
// Returns:
//   - A success message if the deployment is created successfully.
//   - An error if the deployment creation fails.
func (h *Handler) Add(ctx *gofr.Context) (any, error) {
	err := h.deployService.Add(ctx)
	if err != nil {
		return nil, err
	}

	return "Deployment Created", nil
}
