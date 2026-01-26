package tui

import (
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

func TestNewMapDetailModel(t *testing.T) {
	m := newMapDetailModel(80, 24)

	if m.mapInfo != nil {
		t.Error("expected mapInfo to be nil initially")
	}
	if m.cursor != 0 {
		t.Errorf("expected cursor to be 0, got %d", m.cursor)
	}
	if m.width != 80 {
		t.Errorf("expected width to be 80, got %d", m.width)
	}
	if m.height != 24 {
		t.Errorf("expected height to be 24, got %d", m.height)
	}
}

func TestMapDetailModel_SetMap(t *testing.T) {
	m := newMapDetailModel(80, 24)

	mapInfo := &MapInfo{
		ID:         1,
		Name:       "test_map",
		Type:       "hash",
		KeySize:    4,
		ValueSize:  8,
		MaxEntries: 1024,
		Flags:      0,
		MemLock:    4096,
		LoadedAt:   "2024-01-01 12:00:00",
		UID:        0,
	}

	m.SetMap(mapInfo)

	if m.mapInfo != mapInfo {
		t.Error("expected mapInfo to be set")
	}
	if m.cursor != 0 {
		t.Errorf("expected cursor to be 0, got %d", m.cursor)
	}
}

func TestMapDetailModel_SetMapResetsState(t *testing.T) {
	m := newMapDetailModel(80, 24)

	// Set first map
	m.SetMap(&MapInfo{ID: 1, Name: "map1"})

	// Set second map - cursor should reset
	m.SetMap(&MapInfo{ID: 2, Name: "map2"})

	if m.cursor != 0 {
		t.Errorf("expected cursor to reset to 0, got %d", m.cursor)
	}
	if m.mapInfo.ID != 2 {
		t.Errorf("expected map ID to be 2, got %d", m.mapInfo.ID)
	}
}

func TestMapDetailModel_EnterSelectsDump(t *testing.T) {
	m := newMapDetailModel(80, 24)

	mapInfo := &MapInfo{
		ID:   1,
		Name: "test_map",
	}
	m.SetMap(mapInfo)

	// Press enter - should signal dump selection
	var dumpSelected bool
	m, _, dumpSelected = m.Update(tea.KeyMsg{Type: tea.KeyEnter})

	if !dumpSelected {
		t.Error("expected dump to be selected on Enter")
	}
}

func TestMapDetailModel_EnterNoMap(t *testing.T) {
	m := newMapDetailModel(80, 24)

	// Press enter when no map is set
	var dumpSelected bool
	m, _, dumpSelected = m.Update(tea.KeyMsg{Type: tea.KeyEnter})

	if dumpSelected {
		t.Error("expected dump not to be selected when no map is set")
	}
}

func TestMapDetailModel_ViewRendersAllFields(t *testing.T) {
	m := newMapDetailModel(80, 24)

	mapInfo := &MapInfo{
		ID:          42,
		Name:        "my_bpf_map",
		Type:        "hash",
		KeySize:     8,
		ValueSize:   16,
		MaxEntries:  2048,
		Flags:       1,
		MemLock:     8192,
		LoadedAt:    "2024-06-15 10:30:00",
		UID:         1000,
		Pinned:      true,
		PinnedPaths: []string{"/sys/fs/bpf/my_map"},
	}
	m.SetMap(mapInfo)

	view := m.View()

	// Check that all fields are rendered
	expectedFields := []string{
		"42",         // ID
		"my_bpf_map", // Name
		"hash",       // Type
		"8",          // KeySize
		"16",         // ValueSize
		"2048",       // MaxEntries
		"1",          // Flags (need to be careful, this appears in multiple places)
		"8192",       // MemLock
		"2024-06-15", // LoadedAt (partial)
		"1000",       // UID
		"Yes",        // Pinned
		"/sys/fs/bpf/my_map", // PinnedPath
		"Dump Contents",
	}

	for _, field := range expectedFields {
		if !strings.Contains(view, field) {
			t.Errorf("expected view to contain %q", field)
		}
	}
}

func TestMapDetailModel_ViewNoMap(t *testing.T) {
	m := newMapDetailModel(80, 24)

	view := m.View()

	if !strings.Contains(view, "No map selected") {
		t.Error("expected view to show 'No map selected' when no map is set")
	}
}

func TestMapDetailModel_DumpOptionHighlighted(t *testing.T) {
	m := newMapDetailModel(80, 24)

	mapInfo := &MapInfo{
		ID:   1,
		Name: "test_map",
	}
	m.SetMap(mapInfo)

	content := m.renderContent()

	// Dump option should be highlighted (has ▶ prefix)
	if !strings.Contains(content, "▶ Dump Contents") {
		t.Error("expected Dump Contents to be highlighted with ▶")
	}
}

func TestMapDetailModel_SetSize(t *testing.T) {
	m := newMapDetailModel(80, 24)

	m.SetSize(120, 40)

	if m.width != 120 {
		t.Errorf("expected width to be 120, got %d", m.width)
	}
	if m.height != 40 {
		t.Errorf("expected height to be 40, got %d", m.height)
	}
}

func TestMapDetailModel_GetMapInfo(t *testing.T) {
	m := newMapDetailModel(80, 24)

	// No map set
	if m.GetMapInfo() != nil {
		t.Error("expected GetMapInfo to return nil when no map set")
	}

	// Map set
	mapInfo := &MapInfo{ID: 42, Name: "test"}
	m.SetMap(mapInfo)

	if m.GetMapInfo() != mapInfo {
		t.Error("expected GetMapInfo to return the set map")
	}
}

func TestMapDetailModel_GetMapID(t *testing.T) {
	m := newMapDetailModel(80, 24)

	// No map set
	if m.GetMapID() != 0 {
		t.Errorf("expected GetMapID to return 0 when no map set, got %d", m.GetMapID())
	}

	// Map set
	m.SetMap(&MapInfo{ID: 42, Name: "test"})

	if m.GetMapID() != 42 {
		t.Errorf("expected GetMapID to return 42, got %d", m.GetMapID())
	}
}

func TestMapDetailModel_KeyNavigation(t *testing.T) {
	m := newMapDetailModel(80, 24)
	m.SetMap(&MapInfo{ID: 1, Name: "test"})

	// Test up key - should not change cursor (only one action)
	m, _, _ = m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'k'}})
	if m.cursor != 0 {
		t.Errorf("expected cursor to stay at 0, got %d", m.cursor)
	}

	// Test down key - should not change cursor (only one action)
	m, _, _ = m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}})
	if m.cursor != 0 {
		t.Errorf("expected cursor to stay at 0, got %d", m.cursor)
	}

	// Test arrow keys
	m, _, _ = m.Update(tea.KeyMsg{Type: tea.KeyUp})
	if m.cursor != 0 {
		t.Errorf("expected cursor to stay at 0, got %d", m.cursor)
	}

	m, _, _ = m.Update(tea.KeyMsg{Type: tea.KeyDown})
	if m.cursor != 0 {
		t.Errorf("expected cursor to stay at 0, got %d", m.cursor)
	}
}

func TestMapDetailModel_ViewTitle(t *testing.T) {
	m := newMapDetailModel(80, 24)

	mapInfo := &MapInfo{
		ID:   1,
		Name: "my_special_map",
	}
	m.SetMap(mapInfo)

	view := m.View()

	if !strings.Contains(view, "Map: my_special_map") {
		t.Error("expected view title to contain map name")
	}
}

func TestMapDetailModel_RenderContentNoMap(t *testing.T) {
	m := newMapDetailModel(80, 24)

	content := m.renderContent()

	if !strings.Contains(content, "No map selected") {
		t.Error("expected renderContent to show 'No map selected' when no map is set")
	}
}

func TestMapDetailModel_ActionsSection(t *testing.T) {
	m := newMapDetailModel(80, 24)
	m.SetMap(&MapInfo{ID: 1, Name: "test"})

	content := m.renderContent()

	if !strings.Contains(content, "Actions") {
		t.Error("expected content to contain 'Actions' section header")
	}

	if !strings.Contains(content, "Press Enter to dump map contents") {
		t.Error("expected content to contain help text for dump action")
	}
}
