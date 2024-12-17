// Package service provides functionalities for interacting with applications and environments.
// It supports selecting an application and adding environments to it by communicating with an external API.
// it gives users a text-based user interface (TUI) for displaying and selecting items
// using the Charmbracelet bubbletea and list packages.
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
	// listWidth defines the width of the list.
	listWidth = 20
	// listHeight defines the height of the list.
	listHeight = 14
	// listPaddingLeft defines the left padding of the list items.
	listPaddingLeft = 2
	// paginationPadding defines the padding for pagination controls.
	paginationPadding = 4
)

//nolint:gochecknoglobals //required TUI styles for displaying the list
var (
	// itemStyle defines the default style for list items.
	itemStyle = lipgloss.NewStyle().PaddingLeft(listPaddingLeft)
	// selectedItemStyle defines the style for the selected list item.
	selectedItemStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("170"))
	// paginationStyle defines the style for pagination controls.
	paginationStyle = list.DefaultStyles().PaginationStyle.PaddingLeft(paginationPadding)
	// helpStyle defines the style for the help text.
	helpStyle = list.DefaultStyles().HelpStyle
)

// item represents a single item in the list.
type item struct {
	id   int    // ID is the unique identifier for the item.
	name string // Name is the display name of the item.
}

// FilterValue returns the value to be used for filtering list items.
// In this case, it's the name of the item.
func (i *item) FilterValue() string {
	return i.name
}

// itemDelegate is a struct responsible for rendering and interacting with list items.
type itemDelegate struct{}

// Height returns the height of the item (always 1).
func (itemDelegate) Height() int { return 1 }

// Spacing returns the spacing between items (always 0).
func (itemDelegate) Spacing() int { return 0 }

// Update is used to handle updates to the item model. It doesn't do anything in this case.
func (itemDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd {
	return nil
}

// Render renders a single list item, applying the selected item style if it's the currently selected item.
//
//nolint:gocritic //required for rendering list items and implementing ItemDelegate interface
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

// model represents the state of the TUI interface, including the list and selected item.
type model struct {
	choice   *item      // choice is the selected item.
	quitting bool       // quitting indicates if the application is quitting.
	list     list.Model // list holds the list of items displayed in the TUI.
}

// Init initializes the model, returning nil for no commands.
func (*model) Init() tea.Cmd {
	return nil
}

// Update handles updates from messages, such as key presses or window resizing.
// It updates the list and handles quitting or selecting an item.
func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.list.SetWidth(msg.Width)
		return m, nil

	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "q", "ctrl+c":
			m.quitting = true // Set quitting to true when 'q' or 'ctrl+c' is pressed.
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

// View renders the view of the current model, displaying the list to the user.
func (m *model) View() string {
	return "\n" + m.list.View()
}
