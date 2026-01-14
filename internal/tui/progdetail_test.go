package tui

import (
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

func TestNewProgDetailModel(t *testing.T) {
	m := newProgDetailModel(80, 24)

	if m.program != nil {
		t.Error("expected program to be nil initially")
	}
	if m.cursor != -1 {
		t.Errorf("expected cursor to be -1, got %d", m.cursor)
	}
	if m.width != 80 {
		t.Errorf("expected width to be 80, got %d", m.width)
	}
	if m.height != 24 {
		t.Errorf("expected height to be 24, got %d", m.height)
	}
}

func TestProgDetailModel_SetProgram(t *testing.T) {
	m := newProgDetailModel(80, 24)

	// Test with program that has maps
	prog := &ProgramInfo{
		ID:          1,
		Name:        "test_prog",
		Type:        "kprobe",
		Tag:         "abc123",
		GPL:         true,
		LoadedAt:    "2024-01-01 12:00:00",
		UID:         0,
		BytesXlated: 100,
		BytesJIT:    200,
		MemLock:     4096,
		MapIDs:      []uint32{10, 20, 30},
	}

	m.SetProgram(prog)

	if m.program != prog {
		t.Error("expected program to be set")
	}
	if len(m.mapIDs) != 3 {
		t.Errorf("expected 3 map IDs, got %d", len(m.mapIDs))
	}
	if m.cursor != 0 {
		t.Errorf("expected cursor to be 0 when maps exist, got %d", m.cursor)
	}
}

func TestProgDetailModel_SetProgramNoMaps(t *testing.T) {
	m := newProgDetailModel(80, 24)

	prog := &ProgramInfo{
		ID:     1,
		Name:   "test_prog",
		MapIDs: []uint32{},
	}

	m.SetProgram(prog)

	if m.cursor != -1 {
		t.Errorf("expected cursor to be -1 when no maps, got %d", m.cursor)
	}
	if len(m.mapIDs) != 0 {
		t.Errorf("expected 0 map IDs, got %d", len(m.mapIDs))
	}
}

func TestProgDetailModel_CursorNavigation(t *testing.T) {
	m := newProgDetailModel(80, 24)

	prog := &ProgramInfo{
		ID:     1,
		Name:   "test_prog",
		MapIDs: []uint32{10, 20, 30},
	}
	m.SetProgram(prog)

	// Initial cursor should be at 0
	if m.cursor != 0 {
		t.Errorf("expected cursor at 0, got %d", m.cursor)
	}

	// Move down
	m, _, _ = m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}})
	if m.cursor != 1 {
		t.Errorf("expected cursor at 1 after down, got %d", m.cursor)
	}

	// Move down again
	m, _, _ = m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}})
	if m.cursor != 2 {
		t.Errorf("expected cursor at 2 after second down, got %d", m.cursor)
	}

	// Try to move down past end - should stay at 2
	m, _, _ = m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}})
	if m.cursor != 2 {
		t.Errorf("expected cursor to stay at 2, got %d", m.cursor)
	}

	// Move up
	m, _, _ = m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'k'}})
	if m.cursor != 1 {
		t.Errorf("expected cursor at 1 after up, got %d", m.cursor)
	}

	// Move up to beginning
	m, _, _ = m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'k'}})
	if m.cursor != 0 {
		t.Errorf("expected cursor at 0, got %d", m.cursor)
	}

	// Try to move up past beginning - should stay at 0
	m, _, _ = m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'k'}})
	if m.cursor != 0 {
		t.Errorf("expected cursor to stay at 0, got %d", m.cursor)
	}
}

func TestProgDetailModel_ArrowKeyNavigation(t *testing.T) {
	m := newProgDetailModel(80, 24)

	prog := &ProgramInfo{
		ID:     1,
		Name:   "test_prog",
		MapIDs: []uint32{10, 20},
	}
	m.SetProgram(prog)

	// Move down with arrow key
	m, _, _ = m.Update(tea.KeyMsg{Type: tea.KeyDown})
	if m.cursor != 1 {
		t.Errorf("expected cursor at 1 after down arrow, got %d", m.cursor)
	}

	// Move up with arrow key
	m, _, _ = m.Update(tea.KeyMsg{Type: tea.KeyUp})
	if m.cursor != 0 {
		t.Errorf("expected cursor at 0 after up arrow, got %d", m.cursor)
	}
}

func TestProgDetailModel_EnterSelectsMap(t *testing.T) {
	m := newProgDetailModel(80, 24)

	prog := &ProgramInfo{
		ID:     1,
		Name:   "test_prog",
		MapIDs: []uint32{10, 20, 30},
	}
	m.SetProgram(prog)

	// Move to second map
	m, _, _ = m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}})

	// Press enter
	var selectedMapID *uint32
	m, _, selectedMapID = m.Update(tea.KeyMsg{Type: tea.KeyEnter})

	if selectedMapID == nil {
		t.Fatal("expected map ID to be selected")
	}
	if *selectedMapID != 20 {
		t.Errorf("expected selected map ID to be 20, got %d", *selectedMapID)
	}
}

func TestProgDetailModel_EnterNoMaps(t *testing.T) {
	m := newProgDetailModel(80, 24)

	prog := &ProgramInfo{
		ID:     1,
		Name:   "test_prog",
		MapIDs: []uint32{},
	}
	m.SetProgram(prog)

	// Press enter when no maps
	var selectedMapID *uint32
	m, _, selectedMapID = m.Update(tea.KeyMsg{Type: tea.KeyEnter})

	if selectedMapID != nil {
		t.Error("expected no map ID to be selected when no maps exist")
	}
}

func TestProgDetailModel_ViewRendersAllFields(t *testing.T) {
	m := newProgDetailModel(80, 24)

	prog := &ProgramInfo{
		ID:          42,
		Name:        "my_bpf_prog",
		Type:        "tracepoint",
		Tag:         "deadbeef",
		GPL:         true,
		LoadedAt:    "2024-06-15 10:30:00",
		UID:         1000,
		BytesXlated: 512,
		BytesJIT:    1024,
		MemLock:     8192,
		MapIDs:      []uint32{100, 200},
	}
	m.SetProgram(prog)

	view := m.View()

	// Check that all fields are rendered
	expectedFields := []string{
		"42",          // ID
		"my_bpf_prog", // Name
		"tracepoint",  // Type
		"deadbeef",    // Tag
		"Yes",         // GPL
		"2024-06-15",  // LoadedAt (partial)
		"1000",        // UID
		"512",         // BytesXlated
		"1024",        // BytesJIT
		"8192",        // MemLock
		"Map ID: 100", // First map
		"Map ID: 200", // Second map
	}

	for _, field := range expectedFields {
		if !strings.Contains(view, field) {
			t.Errorf("expected view to contain %q", field)
		}
	}
}

func TestProgDetailModel_ViewNoProgram(t *testing.T) {
	m := newProgDetailModel(80, 24)

	view := m.View()

	if !strings.Contains(view, "No program selected") {
		t.Error("expected view to show 'No program selected' when no program is set")
	}
}

func TestProgDetailModel_ViewNoMaps(t *testing.T) {
	m := newProgDetailModel(80, 24)

	prog := &ProgramInfo{
		ID:     1,
		Name:   "test_prog",
		MapIDs: []uint32{},
	}
	m.SetProgram(prog)

	view := m.View()

	if !strings.Contains(view, "No associated maps") {
		t.Error("expected view to show 'No associated maps' when program has no maps")
	}
}

func TestProgDetailModel_ViewGPLNo(t *testing.T) {
	m := newProgDetailModel(80, 24)

	prog := &ProgramInfo{
		ID:   1,
		Name: "test_prog",
		GPL:  false,
	}
	m.SetProgram(prog)

	content := m.renderContent()

	if !strings.Contains(content, "No") {
		t.Error("expected GPL to show 'No' when false")
	}
}

func TestProgDetailModel_HasMaps(t *testing.T) {
	m := newProgDetailModel(80, 24)

	// No program set
	if m.HasMaps() {
		t.Error("expected HasMaps to return false when no program set")
	}

	// Program with no maps
	m.SetProgram(&ProgramInfo{ID: 1, MapIDs: []uint32{}})
	if m.HasMaps() {
		t.Error("expected HasMaps to return false when program has no maps")
	}

	// Program with maps
	m.SetProgram(&ProgramInfo{ID: 1, MapIDs: []uint32{10}})
	if !m.HasMaps() {
		t.Error("expected HasMaps to return true when program has maps")
	}
}

func TestProgDetailModel_SelectedMapID(t *testing.T) {
	m := newProgDetailModel(80, 24)

	// No program set
	if m.SelectedMapID() != nil {
		t.Error("expected SelectedMapID to return nil when no program set")
	}

	// Program with maps
	m.SetProgram(&ProgramInfo{ID: 1, MapIDs: []uint32{10, 20, 30}})

	selected := m.SelectedMapID()
	if selected == nil {
		t.Fatal("expected SelectedMapID to return a value")
	}
	if *selected != 10 {
		t.Errorf("expected selected map ID to be 10, got %d", *selected)
	}

	// Move cursor and check again
	m, _, _ = m.Update(tea.KeyMsg{Type: tea.KeyDown})
	selected = m.SelectedMapID()
	if selected == nil || *selected != 20 {
		t.Errorf("expected selected map ID to be 20 after moving cursor")
	}
}

func TestProgDetailModel_SetSize(t *testing.T) {
	m := newProgDetailModel(80, 24)

	m.SetSize(120, 40)

	if m.width != 120 {
		t.Errorf("expected width to be 120, got %d", m.width)
	}
	if m.height != 40 {
		t.Errorf("expected height to be 40, got %d", m.height)
	}
}

func TestProgDetailModel_SelectedMapHighlighted(t *testing.T) {
	m := newProgDetailModel(80, 24)

	prog := &ProgramInfo{
		ID:     1,
		Name:   "test_prog",
		MapIDs: []uint32{10, 20},
	}
	m.SetProgram(prog)

	content := m.renderContent()

	// First map should be highlighted (has ▶ prefix)
	if !strings.Contains(content, "▶ Map ID: 10") {
		t.Error("expected first map to be highlighted with ▶")
	}

	// Second map should not be highlighted
	if strings.Contains(content, "▶ Map ID: 20") {
		t.Error("expected second map to not be highlighted")
	}
}
