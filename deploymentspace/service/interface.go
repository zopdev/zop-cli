package service

import (
	"gofr.dev/pkg/gofr"

	cloudSvc "zop.dev/cli/zop/cloud/service/list"
	envSvc "zop.dev/cli/zop/environment/service"
)

// CloudAccountService defines the interface for managing cloud accounts.
// It provides methods to retrieve cloud accounts available to the user.
type CloudAccountService interface {
	// GetAccounts retrieves a list of cloud accounts.
	//
	// Parameters:
	//  - ctx: The context object containing request and session details.
	//
	// Returns:
	//  - A slice of pointers to CloudAccountResponse objects.
	//  - An error if the retrieval fails.
	GetAccounts(ctx *gofr.Context) ([]*cloudSvc.CloudAccountResponse, error)
}

// EnvironmentService defines the interface for managing environments.
// It provides methods to list the environments associated with a user.
type EnvironmentService interface {
	// List retrieves a list of environments.
	//
	// Parameters:
	//  - ctx: The context object containing request and session details.
	//
	// Returns:
	//  - A slice of Environment objects.
	//  - An error if the retrieval fails.
	List(ctx *gofr.Context) ([]envSvc.Environment, error)
}
