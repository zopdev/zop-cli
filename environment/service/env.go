package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"gofr.dev/pkg/gofr"

	"zop.dev/cli/zop/utils"
)

const listTitle = "Select the application where you want to add the environment!"

var (
	// ErrUnableToRenderApps is returned when the application list cannot be rendered.
	ErrUnableToRenderApps = errors.New("unable to render the list of applications")

	// ErrConnectingZopAPI is returned when there is an error connecting to the Zop API.
	ErrConnectingZopAPI = errors.New("unable to connect to Zop API")

	// ErrorAddingEnv is returned when there is an error adding an environment.
	ErrorAddingEnv = errors.New("unable to add environment")

	// ErrNoApplicationSelected is returned when no application is selected.
	ErrNoApplicationSelected = errors.New("no application selected")

	// ErrorFetchingEnvironments is returned when there is an error fetching environments for a given application.
	ErrorFetchingEnvironments = errors.New("unable to fetch environments")
)

// Service represents the application service that handles application and environment operations.
type Service struct {
	appGet ApplicationGetter
}

// New creates a new Service instance with the provided ApplicationGetter.
func New(appGet ApplicationGetter) *Service {
	return &Service{appGet: appGet}
}

// Add prompts the user to add environments to a selected application.
// It returns the number of environments added and an error, if any.
func (s *Service) Add(ctx *gofr.Context) (int, error) {
	app, err := s.getSelectedApplication(ctx)
	if err != nil {
		return 0, err
	}

	ctx.Out.Println("Selected application: ", app.Name)
	ctx.Out.Println("Please provide names of environment to be added...")

	var (
		input string
		level = 1
	)

	// Loop to gather environment names from the user and add them to the application.
	for {
		ctx.Out.Print("Enter environment name: ")

		_, _ = fmt.Scanf("%s", &input)

		err = postEnvironment(ctx, &Environment{Name: input, Level: level, ApplicationID: app.ID})
		if err != nil {
			return level, err
		}

		level++

		ctx.Out.Print("Do you wish to add more? (y/n) ")

		_, _ = fmt.Scanf("%s", &input)

		if input == "n" {
			break
		}
	}

	return level, nil
}

func (s *Service) List(ctx *gofr.Context) ([]Environment, error) {
	app, err := s.getSelectedApplication(ctx)
	if err != nil {
		return nil, err
	}

	resp, err := ctx.GetHTTPService("api-service").
		Get(ctx, fmt.Sprintf("applications/%d/environments", app.ID), nil)
	if err != nil {
		ctx.Logger.Errorf("unable to connect to Zop API! %v", err)

		return nil, ErrConnectingZopAPI
	}

	var data struct {
		Envs []Environment `json:"data"`
	}

	err = getResponse(resp, &data)
	if err != nil {
		ctx.Logger.Errorf("unable to fetch environments, could not unmarshall response %v", err)

		return nil, ErrorFetchingEnvironments
	}

	return data.Envs, nil
}

// getSelectedApplication renders a list of applications for the user to select from.
// It returns the selected application or an error if no selection is made.
func (s *Service) getSelectedApplication(ctx *gofr.Context) (*utils.Item, error) {
	apps, err := s.appGet.List(ctx)
	if err != nil {
		return nil, err
	}

	items := make([]*utils.Item, 0)

	for _, app := range apps {
		items = append(items, &utils.Item{ID: app.ID, Name: app.Name})
	}

	choice, err := utils.RenderList(listTitle, items)
	if err != nil {
		ctx.Logger.Errorf("unable to render the list of applications! %v", err)

		return nil, ErrUnableToRenderApps
	}

	if choice == nil {
		return nil, ErrNoApplicationSelected
	}

	return choice, nil
}

// postEnvironment sends a POST request to the API to add the provided environment to the application.
// It returns an error if the request fails or the response status code is not created (201).
func postEnvironment(ctx *gofr.Context, env *Environment) error {
	body, _ := json.Marshal(env)

	resp, err := ctx.GetHTTPService("api-service").
		PostWithHeaders(ctx, "environments", nil, body, map[string]string{
			"Content-Type": "application/json",
		})
	if err != nil {
		ctx.Logger.Errorf("unable to connect to Zop API! %v", err)

		return ErrConnectingZopAPI
	}

	if resp.StatusCode != http.StatusCreated {
		var errResp struct {
			Errors any `json:"errors,omitempty"`
		}

		err = getResponse(resp, &errResp)
		if err != nil {
			ctx.Logger.Errorf("unable to add environment!, could not decode error message %v", err)
		}

		ctx.Logger.Errorf("unable to add environment! %v", resp)

		return ErrorAddingEnv
	}

	return nil
}

// getResponse reads the HTTP response body and unmarshals it into the provided interface.
func getResponse(resp *http.Response, i any) error {
	defer resp.Body.Close()

	b, _ := io.ReadAll(resp.Body)

	err := json.Unmarshal(b, i)
	if err != nil {
		return err
	}

	return nil
}
