package tui

import (
	"fmt"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

// progItem represents a BPF program in the list.
type progItem struct {
	info ProgramInfo
}

// FilterValue implements list.Item interface for fuzzy filtering.
func (i progItem) FilterValue() string { return i.info.Name }

// Title returns the program title for display (ID and Name).
func (i progItem) Title() string {
	return fmt.Sprintf("[%d] %s", i.info.ID, i.info.Name)
}

// Description returns the program description for display (Type and Tag).
func (i progItem) Description() string {
	return fmt.Sprintf("Type: %s | Tag: %s", i.info.Type, i.info.Tag)
}

// progListModel manages the programs list state.
type progListModel struct {
	list     list.Model
	programs []ProgramInfo
	err      error
}

// newProgListModel creates a new programs list model.
func newProgListModel(width, height int) progListModel {
	// Create delegate for custom item rendering
	delegate := list.NewDefaultDelegate()
	delegate.Styles.SelectedTitle = selectedStyle
	delegate.Styles.SelectedDesc = selectedStyle.Foreground(dimStyle.GetForeground())

	// Calculate list dimensions (leave room for title and help bar)
	listHeight := height - 6
	if listHeight < 3 {
		listHeight = 3
	}

	l := list.New([]list.Item{}, delegate, width, listHeight)
	l.Title = "BPF Programs"
	l.SetShowStatusBar(true)
	l.SetFilteringEnabled(true)
	l.SetShowHelp(false)
	l.Styles.Title = titleStyle

	return progListModel{
		list:     l,
		programs: []ProgramInfo{},
	}
}

// Init implements tea.Model for progListModel.
func (m progListModel) Init() tea.Cmd {
	return nil
}

// SetPrograms updates the list with new program data.
func (m *progListModel) SetPrograms(programs []ProgramInfo) {
	m.programs = programs
	items := make([]list.Item, len(programs))
	for i, prog := range programs {
		items[i] = progItem{info: prog}
	}
	m.list.SetItems(items)
}

// Update handles messages for the programs list.
// Returns the updated model, an optional command, and the selected program if Enter was pressed.
func (m progListModel) Update(msg tea.Msg) (progListModel, tea.Cmd, *ProgramInfo) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		// Don't handle enter if we're filtering
		if m.list.FilterState() == list.Filtering {
			break
		}

		switch msg.String() {
		case "enter":
			// Get selected item and return its program info
			if item, ok := m.list.SelectedItem().(progItem); ok {
				return m, nil, &item.info
			}
		}
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd, nil
}

// View renders the programs list.
func (m progListModel) View() string {
	if len(m.programs) == 0 && m.err == nil {
		return titleStyle.Render("BPF Programs") + "\n\n" +
			dimStyle.Render("No BPF programs loaded")
	}

	if m.err != nil {
		return titleStyle.Render("BPF Programs") + "\n\n" +
			errorStyle.Render(fmt.Sprintf("Error: %v", m.err))
	}

	return m.list.View()
}

// SetSize updates the list dimensions.
func (m *progListModel) SetSize(width, height int) {
	listHeight := height - 6
	if listHeight < 3 {
		listHeight = 3
	}
	m.list.SetSize(width, listHeight)
}

// SetError sets an error state for the list.
func (m *progListModel) SetError(err error) {
	m.err = err
}

// SelectedItem returns the currently selected program, if any.
func (m progListModel) SelectedItem() *ProgramInfo {
	if item, ok := m.list.SelectedItem().(progItem); ok {
		return &item.info
	}
	return nil
}

// IsFiltering returns true if the list is currently in filtering mode.
func (m progListModel) IsFiltering() bool {
	return m.list.FilterState() == list.Filtering
}

// ResetFilter clears any active filter and restores the full list.
func (m *progListModel) ResetFilter() {
	m.list.ResetFilter()
}
