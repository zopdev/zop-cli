package service

import (
	"errors"
	"fmt"

	"gofr.dev/pkg/gofr"

	cloudSvc "zop.dev/cli/zop/cloud/service/list"
	"zop.dev/cli/zop/utils"
)

const (
	accListTitle         = "Select the cloud account where you want to add the deployment!"
	appListTitle         = "Select the application where you want to add the deployment!"
	deploymentSpaceTitle = "Select the deployment space where you want to add the deployment!"
)

var (
	// ErrUnableToRenderList is returned when the application list cannot be rendered.
	ErrUnableToRenderList = errors.New("unable to render the list")

	// ErrConnectingZopAPI is returned when there is an error connecting to the Zop API.
	ErrConnectingZopAPI = errors.New("unable to connect to Zop API")

	// ErrGettingDeploymentOptions is returned when there is an error adding an environment.
	ErrGettingDeploymentOptions = errors.New("unable to get deployment options")

	// ErrorFetchingEnvironments is returned when there is an error fetching environments for a given application.
	ErrorFetchingEnvironments = errors.New("unable to fetch environments")
)

type Service struct {
	cloudGet CloudAccountService
	appGet   ApplicationService
	envGet   EnvironmentService
}

func New(cloudGet CloudAccountService, appGet ApplicationService, envGet EnvironmentService) *Service {
	return &Service{
		cloudGet: cloudGet,
		appGet:   appGet,
		envGet:   envGet,
	}
}

func (s *Service) Add(ctx *gofr.Context) error {
	//var request = make(map[string]any)

	// apiSvc := ctx.GetHTTPService("api-service")
	cloudAcc, err := s.getSelectedCloudAccount(ctx)
	if err != nil {
		return err
	}

	//request["cloud"]

	ctx.Out.Println("Selected cloud account: ", cloudAcc.Name)

	app, err := s.getSelectedApplication(ctx)
	if err != nil {
		return err
	}

	ctx.Out.Println("Selected application: ", app.Name)

	return nil
}

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

func (s *Service) getSelectedApplication(ctx *gofr.Context) (*utils.Item, error) {
	apps, err := s.appGet.List(ctx)
	if err != nil {
		return nil, err
	}

	items := make([]*utils.Item, 0)

	for _, app := range apps {
		items = append(items, &utils.Item{ID: app.ID, Name: app.Name})
	}

	choice, err := utils.RenderList(appListTitle, items)
	if err != nil {
		ctx.Logger.Errorf("unable to render the list of applications! %v", err)

		return nil, ErrUnableToRenderList
	}

	if choice == nil {
		return nil, &ErrNoItemSelected{"application"}
	}

	return choice, nil
}

func (s *Service) getSelectedOptions(ctx *gofr.Context, id int64) (*DeploymentSpaceOptions, error) {
	resp, err := ctx.GetHTTPService("api-service").
		Get(ctx, fmt.Sprintf("cloud-accounts/%d/deployment-space/options", id), nil)

	if err != nil {
		ctx.Logger.Errorf("error connecting to zop api! %v", err)

		return nil, ErrConnectingZopAPI
	}

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
