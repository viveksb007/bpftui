package tui

import (
	"fmt"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

// mapItem represents a BPF map in the list.
type mapItem struct {
	info MapInfo
}

// FilterValue implements list.Item interface for fuzzy filtering.
func (i mapItem) FilterValue() string { return i.info.Name }

// Title returns the map title for display (ID and Name).
func (i mapItem) Title() string {
	return fmt.Sprintf("[%d] %s", i.info.ID, i.info.Name)
}

// Description returns the map description for display (Type, KeySize, ValueSize, MaxEntries).
func (i mapItem) Description() string {
	return fmt.Sprintf("Type: %s | Key: %d | Value: %d | Max: %d",
		i.info.Type, i.info.KeySize, i.info.ValueSize, i.info.MaxEntries)
}

// mapListModel manages the maps list state.
type mapListModel struct {
	list list.Model
	maps []MapInfo
	err  error
}

// newMapListModel creates a new maps list model.
func newMapListModel(width, height int) mapListModel {
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
	l.Title = "BPF Maps"
	l.SetShowStatusBar(true)
	l.SetFilteringEnabled(true)
	l.SetShowHelp(false)
	l.Styles.Title = titleStyle

	return mapListModel{
		list: l,
		maps: []MapInfo{},
	}
}

// Init implements tea.Model for mapListModel.
func (m mapListModel) Init() tea.Cmd {
	return nil
}

// SetMaps updates the list with new map data.
func (m *mapListModel) SetMaps(maps []MapInfo) {
	m.maps = maps
	items := make([]list.Item, len(maps))
	for i, mapInfo := range maps {
		items[i] = mapItem{info: mapInfo}
	}
	m.list.SetItems(items)
}

// Update handles messages for the maps list.
// Returns the updated model, an optional command, and the selected map if Enter was pressed.
func (m mapListModel) Update(msg tea.Msg) (mapListModel, tea.Cmd, *MapInfo) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		// Don't handle enter if we're filtering
		if m.list.FilterState() == list.Filtering {
			break
		}

		switch msg.String() {
		case "enter":
			// Get selected item and return its map info
			if item, ok := m.list.SelectedItem().(mapItem); ok {
				return m, nil, &item.info
			}
		}
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd, nil
}

// View renders the maps list.
func (m mapListModel) View() string {
	if len(m.maps) == 0 && m.err == nil {
		return titleStyle.Render("BPF Maps") + "\n\n" +
			dimStyle.Render("No BPF maps loaded")
	}

	if m.err != nil {
		return titleStyle.Render("BPF Maps") + "\n\n" +
			errorStyle.Render(fmt.Sprintf("Error: %v", m.err))
	}

	return m.list.View()
}

// SetSize updates the list dimensions.
func (m *mapListModel) SetSize(width, height int) {
	listHeight := height - 6
	if listHeight < 3 {
		listHeight = 3
	}
	m.list.SetSize(width, listHeight)
}

// SetError sets an error state for the list.
func (m *mapListModel) SetError(err error) {
	m.err = err
}

// SelectedItem returns the currently selected map, if any.
func (m mapListModel) SelectedItem() *MapInfo {
	if item, ok := m.list.SelectedItem().(mapItem); ok {
		return &item.info
	}
	return nil
}

// IsFiltering returns true if the list is currently in filtering mode.
func (m mapListModel) IsFiltering() bool {
	return m.list.FilterState() == list.Filtering
}

// ResetFilter clears any active filter and restores the full list.
func (m *mapListModel) ResetFilter() {
	m.list.ResetFilter()
}
