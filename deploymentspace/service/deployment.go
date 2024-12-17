package service

import "gofr.dev/pkg/gofr"

type Service struct {
	cloudGet CloudAccountService
	appGet   ApplicationService
}

func New() *Service {
	return &Service{}
}

func (*Service) Add(ctx *gofr.Context) error {
	apiSvc := ctx.GetHTTPService("api-service")

	return nil
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
