// Package handler provides functionalities for managing applications, including adding new applications and
// listing existing ones with their environments.
// It acts as the CMD handler layer, connecting the application logic to the user interface.
package handler

import (
	"errors"
	"fmt"
	"sort"
	"strings"

	"gofr.dev/pkg/gofr"
	"gofr.dev/pkg/gofr/cmd/terminal"
)

// Errors returned by the handler package.
var (
	// ErrorApplicationNameNotProvided indicates that the application name was not provided as a parameter.
	ErrorApplicationNameNotProvided = errors.New("please enter application name, -name=<application_name>")
)

// Handler represents the HTTP handler responsible for managing applications.
type Handler struct {
	appAdd ApplicationService
}

// New creates a new instance of Handler.
//
// Parameters:
//   - appAdd: An implementation of the ApplicationService interface used to manage applications.
//
// Returns:
//
//	A pointer to a Handler instance.
func New(appAdd ApplicationService) *Handler {
	return &Handler{
		appAdd: appAdd,
	}
}

// Add handles the addition of a new application.
//
// Parameters:
//   - ctx: The application context containing dependencies and utilities.
//
// Returns:
//
//	A success message and an error, if any.
func (h *Handler) Add(ctx *gofr.Context) (any, error) {
	name := ctx.Param("name")
	if name == "" {
		return nil, ErrorApplicationNameNotProvided
	}

	err := h.appAdd.Add(ctx, name)
	if err != nil {
		return nil, err
	}

	return "Application " + name + " added successfully!", nil
}

// List retrieves and displays all applications along with their environments.
//
// Parameters:
//   - ctx: The application context containing dependencies and utilities.
//
// Returns:
//
//	A newline-separated string and an error, if any.
func (h *Handler) List(ctx *gofr.Context) (any, error) {
	apps, err := h.appAdd.List(ctx)
	if err != nil {
		return nil, err
	}

	ctx.Out.Println("Applications and their environments:\n")

	s := strings.Builder{}

	for i, app := range apps {
		ctx.Out.Printf("%d.", i+1)
		ctx.Out.SetColor(terminal.Cyan)
		ctx.Out.Printf(" %s \n\t", app.Name)
		ctx.Out.ResetColor()

		sort.Slice(app.Envs, func(i, j int) bool { return app.Envs[i].Order < app.Envs[j].Order })

		for _, env := range app.Envs {
			s.WriteString(fmt.Sprintf("%s > ", env.Name))
		}

		ctx.Out.SetColor(terminal.Green)
		ctx.Out.Println(s.String()[:s.Len()-2])
		ctx.Out.ResetColor()
		s.Reset()
	}

	return "\n", nil
}
