package handler

import "gofr.dev/pkg/gofr"

// DeploymentService defines the interface for deployment-related operations.
//
// It contains methods that allow adding and managing deployments.
type DeploymentService interface {
	// Add creates a new deployment.
	//
	// Parameters:
	//  - ctx: The context object containing request and session details.
	//
	// Returns:
	//  - An error if the deployment creation fails.
	Add(ctx *gofr.Context) error
}
