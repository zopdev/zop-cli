package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"gofr.dev/pkg/gofr"
)

var (
	// ErrUnableToRenderApps is returned when the application list cannot be rendered.
	ErrUnableToRenderApps = errors.New("unable to render the list of applications")
	// ErrConnectingZopAPI is returned when there is an error connecting to the Zop API.
	ErrConnectingZopAPI = errors.New("unable to connect to Zop API")
	// ErrorAddingEnv is returned when there is an error adding an environment.
	ErrorAddingEnv = errors.New("unable to add environment")
	// ErrNoApplicationSelected is returned when no application is selected.
	ErrNoApplicationSelected = errors.New("no application selected")
)

// Service represents the application service that handles application and environment operations.
type Service struct {
	appGet ApplicationGetter // appGet is responsible for fetching the list of applications.
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

	ctx.Out.Println("Selected application: ", app.name)
	ctx.Out.Println("Please provide names of environments to be added...")

	var (
		input string
		level = 1
	)

	// Loop to gather environment names from the user and add them to the application.
	for {
		ctx.Out.Print("Enter environment name: ")

		_, _ = fmt.Scanf("%s", &input)

		err = postEnvironment(ctx, &Environment{Name: input, Level: level, ApplicationID: int64(app.id)})
		if err != nil {
			return level, err
		}

		level++

		// Ask the user if they want to add more environments.
		ctx.Out.Print("Do you wish to add more? (y/n) ")

		_, _ = fmt.Scanf("%s", &input)

		if input == "n" {
			break
		}
	}

	return level, nil
}

// getSelectedApplication renders a list of applications for the user to select from.
// It returns the selected application or an error if no selection is made.
func (s *Service) getSelectedApplication(ctx *gofr.Context) (*item, error) {
	apps, err := s.appGet.List(ctx)
	if err != nil {
		return nil, err
	}

	// Prepare a list of items for the user to select from.
	items := make([]list.Item, 0)
	for _, app := range apps {
		items = append(items, &item{app.ID, app.Name})
	}

	// Initialize the list component for application selection.
	l := list.New(items, itemDelegate{}, listWidth, listHeight)
	l.Title = "Select the application where you want to add the environment!"
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(true)
	l.Styles.PaginationStyle = paginationStyle
	l.Styles.HelpStyle = helpStyle
	l.SetShowStatusBar(false)

	m := model{list: l}

	// Render the list using the bubbletea program.
	if _, er := tea.NewProgram(&m).Run(); er != nil {
		ctx.Logger.Errorf("unable to render the list of applications! %v", er)
		return nil, ErrUnableToRenderApps
	}

	if m.choice == nil {
		return nil, ErrNoApplicationSelected
	}

	return m.choice, nil
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
