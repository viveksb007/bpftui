package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
)

// mapDumpModel manages the map dump view state.
type mapDumpModel struct {
	mapID    uint32
	mapName  string
	entries  []MapEntry
	viewport viewport.Model
	width    int
	height   int
	ready    bool
	loading  bool
	err      error
}

// newMapDumpModel creates a new map dump model.
func newMapDumpModel(width, height int) mapDumpModel {
	return mapDumpModel{
		width:  width,
		height: height,
	}
}

// SetMapDump sets the map entries to display.
func (m *mapDumpModel) SetMapDump(mapID uint32, mapName string, entries []MapEntry) {
	m.mapID = mapID
	m.mapName = mapName
	m.entries = entries
	m.loading = false
	m.err = nil
	m.updateViewport()
}

// SetError sets an error state for the dump view.
func (m *mapDumpModel) SetError(err error) {
	m.err = err
	m.loading = false
	m.updateViewport()
}

// SetLoading sets the loading state.
func (m *mapDumpModel) SetLoading(loading bool) {
	m.loading = loading
	m.updateViewport()
}

// SetSize updates the viewport dimensions.
func (m *mapDumpModel) SetSize(width, height int) {
	m.width = width
	m.height = height
	if m.ready {
		m.viewport.Width = width
		m.viewport.Height = height - 4 // Leave room for title and help bar
	}
	m.updateViewport()
}

// updateViewport refreshes the viewport content.
func (m *mapDumpModel) updateViewport() {
	content := m.renderContent()

	if !m.ready {
		m.viewport = viewport.New(m.width, m.height-4)
		m.viewport.SetContent(content)
		m.ready = true
	} else {
		m.viewport.SetContent(content)
	}
}

// renderContent generates the dump view content.
func (m *mapDumpModel) renderContent() string {
	if m.err != nil {
		return errorStyle.Render(fmt.Sprintf("Error: %v", m.err))
	}

	if m.loading {
		return dimStyle.Render("Loading map contents...")
	}

	if len(m.entries) == 0 {
		return dimStyle.Render("Map contains no entries")
	}

	var b strings.Builder

	for i, entry := range m.entries {
		// Key
		b.WriteString(labelStyle.Render("Key:   "))
		b.WriteString(valueStyle.Render(formatHex(entry.Key)))
		b.WriteString("\n")

		// Value
		b.WriteString(labelStyle.Render("Value: "))
		b.WriteString(valueStyle.Render(formatHex(entry.Value)))
		b.WriteString("\n")

		// Separator between entries (except for last entry)
		if i < len(m.entries)-1 {
			b.WriteString(dimStyle.Render("---"))
			b.WriteString("\n")
		}
	}

	return b.String()
}

// formatHex converts a byte slice to a space-separated hex string.
func formatHex(data []byte) string {
	if len(data) == 0 {
		return "(empty)"
	}

	parts := make([]string, len(data))
	for i, b := range data {
		parts[i] = fmt.Sprintf("%02x", b)
	}
	return strings.Join(parts, " ")
}

// Init implements tea.Model for mapDumpModel.
func (m mapDumpModel) Init() tea.Cmd {
	return nil
}

// Update handles messages for the map dump view.
func (m mapDumpModel) Update(msg tea.Msg) (mapDumpModel, tea.Cmd) {
	// Handle viewport scrolling
	var cmd tea.Cmd
	m.viewport, cmd = m.viewport.Update(msg)
	return m, cmd
}

// View renders the map dump view.
func (m mapDumpModel) View() string {
	var title string
	if m.mapName != "" {
		title = titleStyle.Render(fmt.Sprintf("Map Dump: %s (ID: %d)", m.mapName, m.mapID))
	} else if m.mapID != 0 {
		title = titleStyle.Render(fmt.Sprintf("Map Dump: ID %d", m.mapID))
	} else {
		title = titleStyle.Render("Map Dump")
	}

	if !m.ready {
		return title + "\n\nLoading..."
	}

	return title + "\n\n" + m.viewport.View()
}

// GetMapID returns the ID of the map being dumped.
func (m mapDumpModel) GetMapID() uint32 {
	return m.mapID
}

// GetEntryCount returns the number of entries in the dump.
func (m mapDumpModel) GetEntryCount() int {
	return len(m.entries)
}

// HasError returns true if there's an error state.
func (m mapDumpModel) HasError() bool {
	return m.err != nil
}

// GetError returns the current error, if any.
func (m mapDumpModel) GetError() error {
	return m.err
}

// IsLoading returns true if the dump is loading.
func (m mapDumpModel) IsLoading() bool {
	return m.loading
}
