package service

import (
	"errors"
	"fmt"

	"gofr.dev/pkg/gofr"

	cloudSvc "zop.dev/cli/zop/cloud/service/list"
	envSvc "zop.dev/cli/zop/environment/service"
	"zop.dev/cli/zop/utils"
)

const (
	accListTitle         = "Select the cloud account where you want to add the deployment!"
	deploymentSpaceTitle = "Select the deployment space where you want to add the deployment!"
)

var (
	// ErrUnableToRenderList is returned when the application list cannot be rendered.
	ErrUnableToRenderList = errors.New("unable to render the list")

	// ErrNoOptionsFound is returned when there are no options available for selection.
	ErrNoOptionsFound = errors.New("no options available for selection")
)

func (s *Service) getSelectedCloudAccount(ctx *gofr.Context) (*cloudSvc.CloudAccountResponse, error) {
	accounts, err := s.cloudGet.GetAccounts(ctx)
	if err != nil {
		ctx.Logger.Errorf("unable to fetch cloud accounts! %v", err)
	}

	items := make([]*utils.Item, 0)
	for _, acc := range accounts {
		items = append(items, &utils.Item{ID: acc.ID, Name: acc.Name, Data: acc})
	}

	choice, err := utils.RenderList(accListTitle, items)
	if err != nil {
		ctx.Logger.Errorf("unable to render the list of cloud accounts! %v", err)

		return nil, ErrUnableToRenderList
	}

	if choice == nil || choice.Data == nil {
		return nil, &ErrNoItemSelected{"cloud account"}
	}

	return choice.Data.(*cloudSvc.CloudAccountResponse), nil
}

func (s *Service) getSelectedEnvironment(ctx *gofr.Context) (*envSvc.Environment, error) {
	envs, err := s.envGet.List(ctx)
	if err != nil {
		ctx.Logger.Errorf("unable to fetch environments! %v", err)

		return nil, ErrorFetchingEnvironments
	}

	items := make([]*utils.Item, 0)

	for _, env := range envs {
		items = append(items, &utils.Item{ID: env.ID, Name: env.Name, Data: &env})
	}

	choice, err := utils.RenderList("Select the environment where you want to add the deployment!", items)
	if err != nil {
		ctx.Logger.Errorf("unable to render the list of environments! %v", err)

		return nil, ErrUnableToRenderList
	}

	if choice == nil {
		return nil, &ErrNoItemSelected{"environment"}
	}

	return choice.Data.(*envSvc.Environment), nil
}

func getDeploymentSpaceOptions(ctx *gofr.Context, id int64) (*DeploymentSpaceOptions, error) {
	resp, err := ctx.GetHTTPService("api-service").
		Get(ctx, fmt.Sprintf("cloud-accounts/%d/deployment-space/options", id), nil)
	if err != nil {
		ctx.Logger.Errorf("error connecting to zop api! %v", err)

		return nil, ErrConnectingZopAPI
	}

	defer resp.Body.Close()

	var opts struct {
		Options []*DeploymentSpaceOptions `json:"data"`
	}

	err = utils.GetResponse(resp, &opts)
	if err != nil {
		ctx.Logger.Errorf("error fetching deployment space options! %v", err)

		return nil, ErrGettingDeploymentOptions
	}

	items := make([]*utils.Item, 0)

	for _, opt := range opts.Options {
		items = append(items, &utils.Item{Name: opt.Name, Data: opt})
	}

	choice, err := utils.RenderList(deploymentSpaceTitle, items)
	if err != nil {
		ctx.Logger.Errorf("unable to render the list of deployment spaces! %v", err)

		return nil, ErrUnableToRenderList
	}

	if choice == nil || choice.Data == nil {
		return nil, &ErrNoItemSelected{"deployment space"}
	}

	return choice.Data.(*DeploymentSpaceOptions), nil
}

func getSelectedOption(ctx *gofr.Context, items []map[string]any) (map[string]any, error) {
	listI := make([]*utils.Item, 0)

	if len(items) == 0 {
		return nil, ErrNoOptionsFound
	}

	for _, item := range items {
		listI = append(listI, &utils.Item{Name: item["name"].(string), Data: item})
	}

	choice, err := utils.RenderList("Select the option", listI)
	if err != nil {
		ctx.Logger.Errorf("unable to render the list of environments! %v", err)

		return nil, ErrUnableToRenderList
	}

	if choice == nil || choice.Data == nil {
		return nil, &ErrNoItemSelected{items[0]["type"].(string)}
	}

	return choice.Data.(map[string]any), nil
}
