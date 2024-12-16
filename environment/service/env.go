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
	ErrUnableToRenderApps     = errors.New("unable to render the list of applications")
	ErrConnectingZopAPI       = errors.New("unable to connect to Zop API")
	ErrorAddingEnv            = errors.New("unable to add environment")
	ErrNoApplicationSelected  = errors.New("no application selected")
	ErrorFetchingEnvironments = errors.New("unable to fetch environments")
)

type Service struct {
	appGet ApplicationGetter
}

func New(appGet ApplicationGetter) *Service { return &Service{appGet: appGet} }

func (s *Service) Add(ctx *gofr.Context) (int, error) {
	app, err := s.getSelectedApplication(ctx)
	if err != nil {
		return 0, err
	}

	ctx.Out.Println("Selected application: ", app.name)
	ctx.Out.Println("Please provide names of environment to be added...")

	var (
		input string
		level = 1
	)

	for {
		ctx.Out.Print("Enter environment name: ")

		_, _ = fmt.Scanf("%s", &input)

		err = postEnvironment(ctx, &Environment{Name: input, Level: level, ApplicationID: int64(app.id)})
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
		Get(ctx, fmt.Sprintf("applications/%d/environments", app.id), nil)
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

func (s *Service) getSelectedApplication(ctx *gofr.Context) (*item, error) {
	apps, err := s.appGet.List(ctx)
	if err != nil {
		return nil, err
	}

	items := make([]list.Item, 0)

	for _, app := range apps {
		items = append(items, &item{app.ID, app.Name})
	}

	l := list.New(items, itemDelegate{}, listWidth, listHeight)
	l.Title = "Select the application where you want to add the environment!"
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(true)
	l.Styles.PaginationStyle = paginationStyle
	l.Styles.HelpStyle = helpStyle
	l.SetShowStatusBar(false)

	m := model{list: l}

	if _, er := tea.NewProgram(&m).Run(); er != nil {
		ctx.Logger.Errorf("unable to render the list of applications! %v", er)

		return nil, ErrUnableToRenderApps
	}

	if m.choice == nil {
		return nil, ErrNoApplicationSelected
	}

	return m.choice, nil
}

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

func getResponse(resp *http.Response, i any) error {
	defer resp.Body.Close()

	b, _ := io.ReadAll(resp.Body)

	err := json.Unmarshal(b, i)
	if err != nil {
		return err
	}

	return nil
}
