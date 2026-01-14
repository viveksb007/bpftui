package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
)

// progDetailModel manages the program detail view state.
type progDetailModel struct {
	program  *ProgramInfo
	viewport viewport.Model
	mapIDs   []uint32 // For navigation to maps
	cursor   int      // Selected map ID index (-1 means no map selected)
	width    int
	height   int
	ready    bool
}

// newProgDetailModel creates a new program detail model.
func newProgDetailModel(width, height int) progDetailModel {
	return progDetailModel{
		width:  width,
		height: height,
		cursor: -1,
	}
}

// SetProgram sets the program to display.
func (m *progDetailModel) SetProgram(prog *ProgramInfo) {
	m.program = prog
	m.cursor = -1
	if prog != nil && len(prog.MapIDs) > 0 {
		m.mapIDs = prog.MapIDs
		m.cursor = 0 // Start with first map selected if maps exist
	} else {
		m.mapIDs = nil
	}
	m.updateViewport()
}

// SetSize updates the viewport dimensions.
func (m *progDetailModel) SetSize(width, height int) {
	m.width = width
	m.height = height
	if m.ready {
		m.viewport.Width = width
		m.viewport.Height = height - 4 // Leave room for title and help bar
	}
	m.updateViewport()
}

// updateViewport refreshes the viewport content.
func (m *progDetailModel) updateViewport() {
	if m.program == nil {
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
func (m *progDetailModel) renderContent() string {
	if m.program == nil {
		return dimStyle.Render("No program selected")
	}

	var b strings.Builder
	p := m.program

	// Program info section
	b.WriteString(labelStyle.Render("ID:          "))
	b.WriteString(valueStyle.Render(fmt.Sprintf("%d", p.ID)))
	b.WriteString("\n")

	b.WriteString(labelStyle.Render("Name:        "))
	b.WriteString(valueStyle.Render(p.Name))
	b.WriteString("\n")

	b.WriteString(labelStyle.Render("Type:        "))
	b.WriteString(valueStyle.Render(p.Type))
	b.WriteString("\n")

	b.WriteString(labelStyle.Render("Tag:         "))
	b.WriteString(valueStyle.Render(p.Tag))
	b.WriteString("\n")

	b.WriteString(labelStyle.Render("GPL:         "))
	gplStr := "No"
	if p.GPL {
		gplStr = "Yes"
	}
	b.WriteString(valueStyle.Render(gplStr))
	b.WriteString("\n")

	b.WriteString(labelStyle.Render("Loaded At:   "))
	b.WriteString(valueStyle.Render(p.LoadedAt))
	b.WriteString("\n")

	b.WriteString(labelStyle.Render("UID:         "))
	b.WriteString(valueStyle.Render(fmt.Sprintf("%d", p.UID)))
	b.WriteString("\n")

	b.WriteString(labelStyle.Render("Bytes Xlated:"))
	b.WriteString(valueStyle.Render(fmt.Sprintf("%d", p.BytesXlated)))
	b.WriteString("\n")

	b.WriteString(labelStyle.Render("Bytes JIT:   "))
	b.WriteString(valueStyle.Render(fmt.Sprintf("%d", p.BytesJIT)))
	b.WriteString("\n")

	b.WriteString(labelStyle.Render("MemLock:     "))
	b.WriteString(valueStyle.Render(fmt.Sprintf("%d", p.MemLock)))
	b.WriteString("\n")

	// Associated Maps section
	b.WriteString("\n")
	b.WriteString(titleStyle.Render("Associated Maps"))
	b.WriteString("\n")

	if len(m.mapIDs) == 0 {
		b.WriteString(dimStyle.Render("No associated maps"))
		b.WriteString("\n")
	} else {
		for i, mapID := range m.mapIDs {
			if i == m.cursor {
				b.WriteString(selectedStyle.Render(fmt.Sprintf("▶ Map ID: %d", mapID)))
			} else {
				b.WriteString(normalStyle.Render(fmt.Sprintf("  Map ID: %d", mapID)))
			}
			b.WriteString("\n")
		}
		b.WriteString("\n")
		b.WriteString(dimStyle.Render("Press Enter to view map details"))
	}

	return b.String()
}

// Init implements tea.Model for progDetailModel.
func (m progDetailModel) Init() tea.Cmd {
	return nil
}

// Update handles messages for the program detail view.
// Returns the updated model, an optional command, and the selected map ID if Enter was pressed on a map.
func (m progDetailModel) Update(msg tea.Msg) (progDetailModel, tea.Cmd, *uint32) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			if len(m.mapIDs) > 0 && m.cursor > 0 {
				m.cursor--
				m.updateViewport()
			}
			return m, nil, nil

		case "down", "j":
			if len(m.mapIDs) > 0 && m.cursor < len(m.mapIDs)-1 {
				m.cursor++
				m.updateViewport()
			}
			return m, nil, nil

		case "enter":
			// If a map is selected, return its ID for navigation
			if len(m.mapIDs) > 0 && m.cursor >= 0 && m.cursor < len(m.mapIDs) {
				selectedMapID := m.mapIDs[m.cursor]
				return m, nil, &selectedMapID
			}
			return m, nil, nil
		}
	}

	// Handle viewport scrolling for other keys
	var cmd tea.Cmd
	m.viewport, cmd = m.viewport.Update(msg)
	return m, cmd, nil
}

// View renders the program detail view.
func (m progDetailModel) View() string {
	if m.program == nil {
		return titleStyle.Render("Program Details") + "\n\n" +
			dimStyle.Render("No program selected")
	}

	title := titleStyle.Render(fmt.Sprintf("Program: %s", m.program.Name))

	if !m.ready {
		return title + "\n\nLoading..."
	}

	return title + "\n\n" + m.viewport.View()
}

// HasMaps returns true if the program has associated maps.
func (m progDetailModel) HasMaps() bool {
	return len(m.mapIDs) > 0
}

// SelectedMapID returns the currently selected map ID, or nil if none selected.
func (m progDetailModel) SelectedMapID() *uint32 {
	if len(m.mapIDs) > 0 && m.cursor >= 0 && m.cursor < len(m.mapIDs) {
		return &m.mapIDs[m.cursor]
	}
	return nil
}
