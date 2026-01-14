package tui

import (
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

// menuItem represents an item in the main menu.
type menuItem struct {
	title       string
	description string
	target      ViewState
}

// FilterValue implements list.Item interface for fuzzy filtering.
func (i menuItem) FilterValue() string { return i.title }

// Title returns the menu item title for display.
func (i menuItem) Title() string { return i.title }

// Description returns the menu item description for display.
func (i menuItem) Description() string { return i.description }

// menuModel manages the main menu state.
type menuModel struct {
	list list.Model
}

// newMenuModel creates a new menu model with default items.
func newMenuModel(width, height int) menuModel {
	items := []list.Item{
		menuItem{
			title:       "Programs",
			description: "Browse loaded BPF programs",
			target:      ViewProgList,
		},
		menuItem{
			title:       "Maps",
			description: "Browse loaded BPF maps",
			target:      ViewMapList,
		},
	}

	// Create delegate for custom item rendering
	delegate := list.NewDefaultDelegate()
	delegate.Styles.SelectedTitle = selectedStyle
	delegate.Styles.SelectedDesc = selectedStyle.Foreground(dimStyle.GetForeground())

	// Calculate list dimensions (leave room for title and help bar)
	listHeight := height - 6
	if listHeight < 3 {
		listHeight = 3
	}

	l := list.New(items, delegate, width, listHeight)
	l.Title = "BPF TUI Explorer"
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.SetShowHelp(false)
	l.Styles.Title = titleStyle

	return menuModel{list: l}
}

// Init implements tea.Model for menuModel.
func (m menuModel) Init() tea.Cmd {
	return nil
}

// Update handles messages for the menu.
// Returns the updated model, an optional command, and the selected ViewState if Enter was pressed.
func (m menuModel) Update(msg tea.Msg) (menuModel, tea.Cmd, *ViewState) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			// Get selected item and return its target view
			if item, ok := m.list.SelectedItem().(menuItem); ok {
				return m, nil, &item.target
			}
		}
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd, nil
}

// View renders the menu.
func (m menuModel) View() string {
	return m.list.View()
}

// SetSize updates the menu dimensions.
func (m *menuModel) SetSize(width, height int) {
	listHeight := height - 6
	if listHeight < 3 {
		listHeight = 3
	}
	m.list.SetSize(width, listHeight)
}

// SelectedItem returns the currently selected menu item, if any.
func (m menuModel) SelectedItem() *menuItem {
	if item, ok := m.list.SelectedItem().(menuItem); ok {
		return &item
	}
	return nil
}
