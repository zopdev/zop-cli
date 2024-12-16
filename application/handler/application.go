package handler

import (
	"errors"
	"fmt"
	"sort"
	"strings"

	"gofr.dev/pkg/gofr"
	"gofr.dev/pkg/gofr/cmd/terminal"
)

var (
	ErrorApplicationNameNotProvided = errors.New("please enter application name, -name=<application_name>")
)

type Handler struct {
	appAdd ApplicationService
}

func New(appAdd ApplicationService) *Handler {
	return &Handler{
		appAdd: appAdd,
	}
}

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
