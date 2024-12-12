package service

import (
	"fmt"
	"io"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
	listWidth         = 20
	listHeight        = 14
	listPaddingLeft   = 2
	paginationPadding = 4
)

//nolint:gochecknoglobals //required TUI styles for displaying the list
var (
	itemStyle         = lipgloss.NewStyle().PaddingLeft(listPaddingLeft)
	selectedItemStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("170"))
	paginationStyle   = list.DefaultStyles().PaginationStyle.PaddingLeft(paginationPadding)
	helpStyle         = list.DefaultStyles().HelpStyle
)

type item struct {
	id   int
	name string
}

func (i *item) FilterValue() string { return i.name }

type itemDelegate struct{}

func (itemDelegate) Height() int                             { return 1 }
func (itemDelegate) Spacing() int                            { return 0 }
func (itemDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }

// Render renders the list items with the selected item highlighted.
//
//nolint:gocritic //required to render the list items.
func (itemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(*item)
	if !ok {
		return
	}

	str := fmt.Sprintf("%3d. %s", index+1, i.name)

	fn := itemStyle.Render
	if index == m.Index() {
		fn = func(s ...string) string {
			return selectedItemStyle.Render("> " + strings.Join(s, " "))
		}
	}

	fmt.Fprint(w, fn(str))
}

type model struct {
	choice   *item
	quitting bool
	list     list.Model
}

func (*model) Init() tea.Cmd {
	return nil
}

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.list.SetWidth(msg.Width)
		return m, nil

	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "q", "ctrl+c":
			m.quitting = true
			return m, tea.Quit

		case "enter":
			i, ok := m.list.SelectedItem().(*item)
			if ok {
				m.choice = i
			}

			return m, tea.Quit
		}
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)

	return m, cmd
}

func (m *model) View() string {
	return "\n" + m.list.View()
}
