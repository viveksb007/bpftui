package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
)

// mapDetailModel manages the map detail view state.
type mapDetailModel struct {
	mapInfo  *MapInfo
	viewport viewport.Model
	cursor   int // 0 = Dump option selected
	width    int
	height   int
	ready    bool
}

// newMapDetailModel creates a new map detail model.
func newMapDetailModel(width, height int) mapDetailModel {
	return mapDetailModel{
		width:  width,
		height: height,
		cursor: 0,
	}
}

// SetMap sets the map to display.
func (m *mapDetailModel) SetMap(mapInfo *MapInfo) {
	m.mapInfo = mapInfo
	m.cursor = 0 // Reset cursor to Dump option
	m.updateViewport()
}

// SetSize updates the viewport dimensions.
func (m *mapDetailModel) SetSize(width, height int) {
	m.width = width
	m.height = height
	if m.ready {
		m.viewport.Width = width
		m.viewport.Height = height - 4 // Leave room for title and help bar
	}
	m.updateViewport()
}

// updateViewport refreshes the viewport content.
func (m *mapDetailModel) updateViewport() {
	if m.mapInfo == nil {
		return
	}

	content := m.renderContent()

	if !m.ready {
		m.viewport = viewport.New(m.width, m.height-4)
		m.viewport.SetContent(content)
		m.ready = true
	} else {
		m.viewport.SetContent(content)
	}
}

// renderContent generates the detail view content.
func (m *mapDetailModel) renderContent() string {
	if m.mapInfo == nil {
		return dimStyle.Render("No map selected")
	}

	var b strings.Builder
	mi := m.mapInfo

	// Map info section
	b.WriteString(labelStyle.Render("ID:          "))
	b.WriteString(valueStyle.Render(fmt.Sprintf("%d", mi.ID)))
	b.WriteString("\n")

	b.WriteString(labelStyle.Render("Name:        "))
	b.WriteString(valueStyle.Render(mi.Name))
	b.WriteString("\n")

	b.WriteString(labelStyle.Render("Type:        "))
	b.WriteString(valueStyle.Render(mi.Type))
	b.WriteString("\n")

	b.WriteString(labelStyle.Render("Key Size:    "))
	b.WriteString(valueStyle.Render(fmt.Sprintf("%d", mi.KeySize)))
	b.WriteString("\n")

	b.WriteString(labelStyle.Render("Value Size:  "))
	b.WriteString(valueStyle.Render(fmt.Sprintf("%d", mi.ValueSize)))
	b.WriteString("\n")

	b.WriteString(labelStyle.Render("Max Entries: "))
	b.WriteString(valueStyle.Render(fmt.Sprintf("%d", mi.MaxEntries)))
	b.WriteString("\n")

	b.WriteString(labelStyle.Render("Flags:       "))
	b.WriteString(valueStyle.Render(fmt.Sprintf("%d", mi.Flags)))
	b.WriteString("\n")

	b.WriteString(labelStyle.Render("MemLock:     "))
	b.WriteString(valueStyle.Render(fmt.Sprintf("%d", mi.MemLock)))
	b.WriteString("\n")

	b.WriteString(labelStyle.Render("Loaded At:   "))
	b.WriteString(valueStyle.Render(mi.LoadedAt))
	b.WriteString("\n")

	b.WriteString(labelStyle.Render("UID:         "))
	b.WriteString(valueStyle.Render(fmt.Sprintf("%d", mi.UID)))
	b.WriteString("\n")

	// Actions section
	b.WriteString("\n")
	b.WriteString(titleStyle.Render("Actions"))
	b.WriteString("\n")

	// Dump Contents action
	if m.cursor == 0 {
		b.WriteString(selectedStyle.Render("▶ Dump Contents"))
	} else {
		b.WriteString(normalStyle.Render("  Dump Contents"))
	}
	b.WriteString("\n")
	b.WriteString("\n")
	b.WriteString(dimStyle.Render("Press Enter to dump map contents"))

	return b.String()
}

// Init implements tea.Model for mapDetailModel.
func (m mapDetailModel) Init() tea.Cmd {
	return nil
}

// Update handles messages for the map detail view.
// Returns the updated model, an optional command, and whether Dump was selected.
func (m mapDetailModel) Update(msg tea.Msg) (mapDetailModel, tea.Cmd, bool) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			// Currently only one action, so no movement needed
			// But keep the handler for future extensibility
			return m, nil, false

		case "down", "j":
			// Currently only one action, so no movement needed
			return m, nil, false

		case "enter":
			// If Dump is selected and we have a map, signal to navigate to MapDump
			if m.mapInfo != nil && m.cursor == 0 {
				return m, nil, true
			}
			return m, nil, false
		}
	}

	// Handle viewport scrolling for other keys
	var cmd tea.Cmd
	m.viewport, cmd = m.viewport.Update(msg)
	return m, cmd, false
}

// View renders the map detail view.
func (m mapDetailModel) View() string {
	if m.mapInfo == nil {
		return titleStyle.Render("Map Details") + "\n\n" +
			dimStyle.Render("No map selected")
	}

	title := titleStyle.Render(fmt.Sprintf("Map: %s", m.mapInfo.Name))

	if !m.ready {
		return title + "\n\nLoading..."
	}

	return title + "\n\n" + m.viewport.View()
}

// GetMapInfo returns the currently displayed map info.
func (m mapDetailModel) GetMapInfo() *MapInfo {
	return m.mapInfo
}

// GetMapID returns the ID of the currently displayed map, or 0 if none.
func (m mapDetailModel) GetMapID() uint32 {
	if m.mapInfo != nil {
		return m.mapInfo.ID
	}
	return 0
}
