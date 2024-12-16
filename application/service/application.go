package service

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"gofr.dev/pkg/gofr"
)

type Service struct {
}

func New() *Service {
	return &Service{}
}

func (*Service) Add(ctx *gofr.Context, name string) error {
	var (
		envs  []Environment
		input string
	)

	app := &Application{Name: name}
	api := ctx.GetHTTPService("api-service")
	order := 1

	ctx.Out.Print("Do you wish to add environments to the application? (y/n) ")

	_, _ = fmt.Scanf("%s", &input)

	for {
		if input != "y" {
			break
		}

		ctx.Out.Print("Enter environment name: ")

		_, _ = fmt.Scanf("%s", &input)
		envs = append(envs, Environment{Name: input, Order: order})
		order++

		ctx.Out.Print("Do you wish to add more? (y/n) ")

		_, _ = fmt.Scanf("%s", &input)

		if input == "n" {
			break
		}
	}

	app.Envs = envs
	body, _ := json.Marshal(app)

	resp, err := api.PostWithHeaders(ctx, "applications", nil, body, map[string]string{
		"Content-Type": "application/json",
	})
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return getAPIError(resp)
	}

	return nil
}

func (*Service) List(ctx *gofr.Context) ([]Application, error) {
	api := ctx.GetHTTPService("api-service")

	reps, err := api.Get(ctx, "applications", nil)
	if err != nil {
		return nil, err
	}
	defer reps.Body.Close()

	var apps struct {
		Data []Application `json:"data"`
	}

	body, _ := io.ReadAll(reps.Body)

	err = json.Unmarshal(body, &apps)
	if err != nil {
		return nil, &ErrAPIService{
			StatusCode: http.StatusInternalServerError,
			Message:    "Internal Server Error",
		}
	}

	return apps.Data, nil
}
