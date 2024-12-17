package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

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

	// ErrConnectingZopAPI is returned when there is an error connecting to the Zop API.
	ErrConnectingZopAPI = errors.New("unable to connect to Zop API")

	// ErrGettingDeploymentOptions is returned when there is an error adding an environment.
	ErrGettingDeploymentOptions = errors.New("unable to get deployment options")

	// ErrorFetchingEnvironments is returned when there is an error fetching environments for a given application.
	ErrorFetchingEnvironments = errors.New("unable to fetch environments")

	// ErrNoOptionsFound is returned when there are no options available for selection.
	ErrNoOptionsFound = errors.New("no options available for selection")
)

type Service struct {
	cloudGet CloudAccountService
	envGet   EnvironmentService
}

func New(cloudGet CloudAccountService, envGet EnvironmentService) *Service {
	return &Service{
		cloudGet: cloudGet,
		envGet:   envGet,
	}
}

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
		opt, er := getSelectedOption(ctx, option.Data.Option)
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
		return ErrGettingDeploymentOptions
	}

	return resp.Body.Close()
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

func getParameters(opt map[string]any, options *apiResponse) string {
	params := "?"

	for _, v := range options.Data.Next.Params {
		params += fmt.Sprintf("&%s=%s", v, opt[v])
	}

	return params
}
