package utils

import (
	"fmt"
	"io"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
	listPaddingLeft   = 2
	paginationPadding = 4
	listWidth         = 20
	listHeight        = 14
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
type Item struct {
	ID   int    // ID is the unique identifier for the item.
	Name string // Name is the display name of the item.
	Data any
}

// FilterValue returns the value to be used for filtering list items.
// In this case, it's the name of the item.
func (i *Item) FilterValue() string { return i.Name }

// itemDelegate is a struct responsible for rendering and interacting with list items.
type itemDelegate struct{}

// Height returns the height of the item (always 1).
func (itemDelegate) Height() int { return 1 }

// Spacing returns the spacing between items (always 0).
func (itemDelegate) Spacing() int { return 0 }

// Update returns the command to update the list. (always nil).
func (itemDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }

// Render renders the list items with the selected item highlighted.
//
//nolint:gocritic //required for rendering list items and implementing ItemDelegate interface
func (itemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(*Item)
	if !ok {
		return
	}

	str := fmt.Sprintf("%3d. %s", index+1, i.Name)

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
	choice   *Item      // choice is the selected item.
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
			m.quitting = true
			return m, tea.Quit

		case "enter":
			i, ok := m.list.SelectedItem().(*Item)
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

func RenderList(items []*Item) (*Item, error) {
	listItems := make([]list.Item, 0)

	for i := range items {
		listItems = append(listItems, items[i])
	}

	l := list.New(listItems, itemDelegate{}, listWidth, listHeight)
	l.Title = "Select the application where you want to add the environment!"
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(true)
	l.Styles.PaginationStyle = paginationStyle
	l.Styles.HelpStyle = helpStyle
	l.SetShowStatusBar(false)

	m := model{list: l}

	if _, er := tea.NewProgram(&m, tea.WithAltScreen()).Run(); er != nil {
		return nil, er
	}

	return m.choice, nil
}
