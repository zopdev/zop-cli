// Package service provides structures and interfaces for managing deployment options and related operations.
package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"gofr.dev/pkg/gofr"

	"zop.dev/cli/zop/utils"
)

var (

	// ErrConnectingZopAPI is returned when there is an error connecting to the Zop API.
	ErrConnectingZopAPI = errors.New("unable to connect to Zop API")

	// ErrGettingDeploymentOptions is returned when there is an error adding an environment.
	ErrGettingDeploymentOptions = errors.New("unable to get deployment options")

	// ErrorFetchingEnvironments is returned when there is an error fetching environments for a given application.
	ErrorFetchingEnvironments = errors.New("unable to fetch environments")

	// ErrUnknown is returned when an unknown error occurs while processing the request.
	ErrUnknown = errors.New("unknown error occurred while processing the request")
)

// Service represents the core service that handles cloud account and environment-related operations.
type Service struct {
	cloudGet CloudAccountService
	envGet   EnvironmentService
}

// New initializes a new Service instance.
//
// Parameters:
//   - cloudGet: A CloudAccountService instance for retrieving cloud accounts.
//   - envGet: An EnvironmentService instance for retrieving environments.
//
// Returns:
//   - A pointer to the Service instance.
func New(cloudGet CloudAccountService, envGet EnvironmentService) *Service {
	return &Service{
		cloudGet: cloudGet,
		envGet:   envGet,
	}
}

// Add handles the addition of a deployment configuration.
//
// This function selects a cloud account and environment, retrieves deployment options,
// processes the options, and submits the deployment request.
//
// Parameters:
//   - ctx: The context object containing request and session details.
//
// Returns:
//   - An error if any step in the process fails.
func (s *Service) Add(ctx *gofr.Context) error {
	var request = make(map[string]any)

	cloudAcc, err := s.getSelectedCloudAccount(ctx)
	if err != nil {
		return err
	}

	request["cloudAccount"] = cloudAcc

	ctx.Out.Println("Selected cloud account: ", cloudAcc.Name)

	env, err := s.getSelectedEnvironment(ctx)
	if err != nil {
		return err
	}

	ctx.Out.Println("Selected environment"+
		":", env.Name)

	options, err := getDeploymentSpaceOptions(ctx, cloudAcc.ID)
	if err != nil {
		return err
	}

	request[options.Type] = options

	if er := processOptions(ctx, request, options.Path); er != nil {
		return er
	}

	return submitDeployment(ctx, env.ID, request)
}

func processOptions(ctx *gofr.Context, request map[string]any, path string) error {
	var optionName string

	api := ctx.GetHTTPService("api-service")

	resp, err := api.Get(ctx, path[1:], nil)
	if err != nil {
		ctx.Logger.Errorf("error connecting to zop api! %v", err)
		return ErrConnectingZopAPI
	}

	var option apiResponse
	if er := utils.GetResponse(resp, &option); er != nil {
		ctx.Logger.Errorf("error fetching deployment options! %v", er)

		return ErrGettingDeploymentOptions
	}

	resp.Body.Close()

	for {
		optionName = "option"

		if option.Data.Metadata != nil {
			optionName = option.Data.Metadata.Name
		}

		opt, er := getSelectedOption(ctx, option.Data.Option, optionName)
		if er != nil {
			return er
		}

		updateRequestWithOption(request, opt)

		if option.Data.Next == nil {
			break
		}

		params := getParameters(opt, &option)

		resp, er = api.Get(ctx, option.Data.Next.Path[1:]+params, nil)
		if er != nil {
			ctx.Logger.Errorf("error connecting to zop api! %v", er)
			return ErrConnectingZopAPI
		}

		option.Data = nil

		er = utils.GetResponse(resp, &option)
		if er != nil {
			ctx.Logger.Errorf("error fetching deployment options! %v", er)
			return ErrGettingDeploymentOptions
		}

		resp.Body.Close()
	}

	return nil
}

func updateRequestWithOption(request, opt map[string]any) {
	keys := strings.Split(opt["type"].(string), ".")
	current := request

	for i, key := range keys {
		if i == len(keys)-1 {
			current[key] = opt
			break
		}

		if _, exists := current[key]; !exists {
			current[key] = make(map[string]any)
		}

		current = current[key].(map[string]any)
	}
}

func submitDeployment(ctx *gofr.Context, envID int64, request map[string]any) error {
	b, err := json.Marshal(request)
	if err != nil {
		return err
	}

	resp, err := ctx.GetHTTPService("api-service").
		PostWithHeaders(ctx, fmt.Sprintf("environments/%d/deploymentspace", envID), nil, b, map[string]string{
			"Content-Type": "application/json",
		})
	if err != nil {
		return ErrConnectingZopAPI
	}

	if resp.StatusCode != http.StatusCreated {
		var er ErrorResponse

		err = utils.GetResponse(resp, &er)
		if err != nil {
			return ErrUnknown
		}

		return &er
	}

	return resp.Body.Close()
}

func getParameters(opt map[string]any, options *apiResponse) string {
	params := "?"

	for _, v := range options.Data.Next.Params {
		params += fmt.Sprintf("&%s=%s", v, opt[v])
	}

	return params
}
