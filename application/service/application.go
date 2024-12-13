package service

import (
	"encoding/json"
	"fmt"
	"net/http"

	"gofr.dev/pkg/gofr"
)

type Service struct {
}

func New() *Service {
	return &Service{}
}

func (*Service) AddApplication(ctx *gofr.Context, name string) error {
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

	resp, err := api.PostWithHeaders(ctx, "application", nil, body, nil)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return getAPIError(resp)
	}

	return nil
}
