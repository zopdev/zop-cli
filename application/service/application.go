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

func (s *Service) AddApplication(ctx *gofr.Context, name string) error {
	var (
		envs  []Environment
		input string
	)

	app := &Application{Name: name}
	api := ctx.GetHTTPService("api-service")
	order := 1

	ctx.Out.Print("Do you wish to add environments to the application? (y/n) ")

	fmt.Scanf("%s", &input)

	for {
		if input == "y" {
			ctx.Out.Println("Enter environment name:")
			fmt.Scanf("%s", &input)

			envs = append(envs, Environment{Name: input, Order: order})
			order++
		} else {
			break
		}

		ctx.Out.Print("Do you wish to add more? (y/n) ")
		fmt.Scanf("%s", &input)
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
