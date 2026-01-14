package tui

import (
	"errors"
	"strings"
	"testing"
)

func TestNewMapDumpModel(t *testing.T) {
	m := newMapDumpModel(80, 24)

	if m.mapID != 0 {
		t.Errorf("expected mapID to be 0, got %d", m.mapID)
	}
	if m.mapName != "" {
		t.Errorf("expected mapName to be empty, got %q", m.mapName)
	}
	if len(m.entries) != 0 {
		t.Errorf("expected entries to be empty, got %d entries", len(m.entries))
	}
	if m.width != 80 {
		t.Errorf("expected width to be 80, got %d", m.width)
	}
	if m.height != 24 {
		t.Errorf("expected height to be 24, got %d", m.height)
	}
	if m.loading {
		t.Error("expected loading to be false initially")
	}
	if m.err != nil {
		t.Error("expected err to be nil initially")
	}
}

func TestMapDumpModel_SetMapDump(t *testing.T) {
	m := newMapDumpModel(80, 24)

	entries := []MapEntry{
		{Key: []byte{0x01, 0x02}, Value: []byte{0x0a, 0x0b}},
		{Key: []byte{0x03, 0x04}, Value: []byte{0x0c, 0x0d}},
	}

	m.SetMapDump(42, "test_map", entries)

	if m.mapID != 42 {
		t.Errorf("expected mapID to be 42, got %d", m.mapID)
	}
	if m.mapName != "test_map" {
		t.Errorf("expected mapName to be 'test_map', got %q", m.mapName)
	}
	if len(m.entries) != 2 {
		t.Errorf("expected 2 entries, got %d", len(m.entries))
	}
	if m.loading {
		t.Error("expected loading to be false after SetMapDump")
	}
	if m.err != nil {
		t.Error("expected err to be nil after SetMapDump")
	}
}

func TestMapDumpModel_SetMapDumpClearsError(t *testing.T) {
	m := newMapDumpModel(80, 24)

	// Set an error first
	m.SetError(errors.New("test error"))
	if m.err == nil {
		t.Error("expected error to be set")
	}

	// SetMapDump should clear the error
	m.SetMapDump(1, "map", []MapEntry{})
	if m.err != nil {
		t.Error("expected error to be cleared after SetMapDump")
	}
}

func TestMapDumpModel_SetError(t *testing.T) {
	m := newMapDumpModel(80, 24)

	testErr := errors.New("dump failed")
	m.SetError(testErr)

	if m.err != testErr {
		t.Error("expected error to be set")
	}
	if m.loading {
		t.Error("expected loading to be false after SetError")
	}
}

func TestMapDumpModel_SetLoading(t *testing.T) {
	m := newMapDumpModel(80, 24)

	m.SetLoading(true)
	if !m.loading {
		t.Error("expected loading to be true")
	}

	m.SetLoading(false)
	if m.loading {
		t.Error("expected loading to be false")
	}
}

func TestMapDumpModel_SetSize(t *testing.T) {
	m := newMapDumpModel(80, 24)

	m.SetSize(120, 40)

	if m.width != 120 {
		t.Errorf("expected width to be 120, got %d", m.width)
	}
	if m.height != 40 {
		t.Errorf("expected height to be 40, got %d", m.height)
	}
}

func TestFormatHex(t *testing.T) {
	tests := []struct {
		name     string
		input    []byte
		expected string
	}{
		{
			name:     "empty slice",
			input:    []byte{},
			expected: "(empty)",
		},
		{
			name:     "nil slice",
			input:    nil,
			expected: "(empty)",
		},
		{
			name:     "single byte",
			input:    []byte{0x0a},
			expected: "0a",
		},
		{
			name:     "multiple bytes",
			input:    []byte{0x01, 0x02, 0x03},
			expected: "01 02 03",
		},
		{
			name:     "all zeros",
			input:    []byte{0x00, 0x00, 0x00, 0x00},
			expected: "00 00 00 00",
		},
		{
			name:     "all ff",
			input:    []byte{0xff, 0xff},
			expected: "ff ff",
		},
		{
			name:     "mixed values",
			input:    []byte{0x00, 0x0a, 0x10, 0xff},
			expected: "00 0a 10 ff",
		},
		{
			name:     "lowercase hex",
			input:    []byte{0xAB, 0xCD, 0xEF},
			expected: "ab cd ef",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := formatHex(tt.input)
			if result != tt.expected {
				t.Errorf("formatHex(%v) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestMapDumpModel_RenderContentEmpty(t *testing.T) {
	m := newMapDumpModel(80, 24)
	m.SetMapDump(1, "empty_map", []MapEntry{})

	content := m.renderContent()

	if !strings.Contains(content, "Map contains no entries") {
		t.Error("expected content to show 'Map contains no entries' for empty map")
	}
}

func TestMapDumpModel_RenderContentLoading(t *testing.T) {
	m := newMapDumpModel(80, 24)
	m.SetLoading(true)

	content := m.renderContent()

	if !strings.Contains(content, "Loading map contents") {
		t.Error("expected content to show loading message")
	}
}

func TestMapDumpModel_RenderContentError(t *testing.T) {
	m := newMapDumpModel(80, 24)
	m.SetError(errors.New("permission denied"))

	content := m.renderContent()

	if !strings.Contains(content, "Error") {
		t.Error("expected content to show error prefix")
	}
	if !strings.Contains(content, "permission denied") {
		t.Error("expected content to show error message")
	}
}

func TestMapDumpModel_RenderContentWithEntries(t *testing.T) {
	m := newMapDumpModel(80, 24)

	entries := []MapEntry{
		{Key: []byte{0x01, 0x02, 0x03, 0x04}, Value: []byte{0x0a, 0x0b, 0x0c, 0x0d}},
		{Key: []byte{0x05, 0x06}, Value: []byte{0x0e, 0x0f, 0x10}},
	}
	m.SetMapDump(42, "test_map", entries)

	content := m.renderContent()

	// Check first entry
	if !strings.Contains(content, "01 02 03 04") {
		t.Error("expected content to contain first key hex")
	}
	if !strings.Contains(content, "0a 0b 0c 0d") {
		t.Error("expected content to contain first value hex")
	}

	// Check second entry
	if !strings.Contains(content, "05 06") {
		t.Error("expected content to contain second key hex")
	}
	if !strings.Contains(content, "0e 0f 10") {
		t.Error("expected content to contain second value hex")
	}

	// Check separator
	if !strings.Contains(content, "---") {
		t.Error("expected content to contain separator between entries")
	}

	// Check labels
	if !strings.Contains(content, "Key:") {
		t.Error("expected content to contain 'Key:' label")
	}
	if !strings.Contains(content, "Value:") {
		t.Error("expected content to contain 'Value:' label")
	}
}

func TestMapDumpModel_RenderContentSingleEntry(t *testing.T) {
	m := newMapDumpModel(80, 24)

	entries := []MapEntry{
		{Key: []byte{0xaa}, Value: []byte{0xbb}},
	}
	m.SetMapDump(1, "single", entries)

	content := m.renderContent()

	// Should have key and value
	if !strings.Contains(content, "aa") {
		t.Error("expected content to contain key hex")
	}
	if !strings.Contains(content, "bb") {
		t.Error("expected content to contain value hex")
	}

	// Should NOT have separator for single entry
	lines := strings.Split(content, "\n")
	separatorCount := 0
	for _, line := range lines {
		if strings.Contains(line, "---") {
			separatorCount++
		}
	}
	if separatorCount > 0 {
		t.Error("expected no separator for single entry")
	}
}

func TestMapDumpModel_ViewWithMapName(t *testing.T) {
	m := newMapDumpModel(80, 24)
	m.SetMapDump(42, "my_map", []MapEntry{})

	view := m.View()

	if !strings.Contains(view, "Map Dump: my_map") {
		t.Error("expected view title to contain map name")
	}
	if !strings.Contains(view, "ID: 42") {
		t.Error("expected view title to contain map ID")
	}
}

func TestMapDumpModel_ViewWithoutMapName(t *testing.T) {
	m := newMapDumpModel(80, 24)
	m.SetMapDump(42, "", []MapEntry{})

	view := m.View()

	if !strings.Contains(view, "Map Dump: ID 42") {
		t.Error("expected view title to show ID when no name")
	}
}

func TestMapDumpModel_ViewNoMapSet(t *testing.T) {
	m := newMapDumpModel(80, 24)

	view := m.View()

	if !strings.Contains(view, "Map Dump") {
		t.Error("expected view to contain 'Map Dump' title")
	}
}

func TestMapDumpModel_GetMapID(t *testing.T) {
	m := newMapDumpModel(80, 24)

	if m.GetMapID() != 0 {
		t.Errorf("expected GetMapID to return 0 initially, got %d", m.GetMapID())
	}

	m.SetMapDump(42, "test", []MapEntry{})

	if m.GetMapID() != 42 {
		t.Errorf("expected GetMapID to return 42, got %d", m.GetMapID())
	}
}

func TestMapDumpModel_GetEntryCount(t *testing.T) {
	m := newMapDumpModel(80, 24)

	if m.GetEntryCount() != 0 {
		t.Errorf("expected GetEntryCount to return 0 initially, got %d", m.GetEntryCount())
	}

	entries := []MapEntry{
		{Key: []byte{1}, Value: []byte{2}},
		{Key: []byte{3}, Value: []byte{4}},
		{Key: []byte{5}, Value: []byte{6}},
	}
	m.SetMapDump(1, "test", entries)

	if m.GetEntryCount() != 3 {
		t.Errorf("expected GetEntryCount to return 3, got %d", m.GetEntryCount())
	}
}

func TestMapDumpModel_HasError(t *testing.T) {
	m := newMapDumpModel(80, 24)

	if m.HasError() {
		t.Error("expected HasError to return false initially")
	}

	m.SetError(errors.New("test"))

	if !m.HasError() {
		t.Error("expected HasError to return true after SetError")
	}

	m.SetMapDump(1, "test", []MapEntry{})

	if m.HasError() {
		t.Error("expected HasError to return false after SetMapDump")
	}
}

func TestMapDumpModel_GetError(t *testing.T) {
	m := newMapDumpModel(80, 24)

	if m.GetError() != nil {
		t.Error("expected GetError to return nil initially")
	}

	testErr := errors.New("test error")
	m.SetError(testErr)

	if m.GetError() != testErr {
		t.Error("expected GetError to return the set error")
	}
}

func TestMapDumpModel_IsLoading(t *testing.T) {
	m := newMapDumpModel(80, 24)

	if m.IsLoading() {
		t.Error("expected IsLoading to return false initially")
	}

	m.SetLoading(true)

	if !m.IsLoading() {
		t.Error("expected IsLoading to return true after SetLoading(true)")
	}

	m.SetLoading(false)

	if m.IsLoading() {
		t.Error("expected IsLoading to return false after SetLoading(false)")
	}
}

func TestMapDumpModel_EmptyKeyValue(t *testing.T) {
	m := newMapDumpModel(80, 24)

	entries := []MapEntry{
		{Key: []byte{}, Value: []byte{0x01}},
		{Key: []byte{0x02}, Value: []byte{}},
	}
	m.SetMapDump(1, "test", entries)

	content := m.renderContent()

	// Empty key/value should show "(empty)"
	if !strings.Contains(content, "(empty)") {
		t.Error("expected content to show '(empty)' for empty key or value")
	}
}

func TestMapDumpModel_ErrorPriorityOverLoading(t *testing.T) {
	m := newMapDumpModel(80, 24)

	// Set both loading and error
	m.loading = true
	m.err = errors.New("test error")

	content := m.renderContent()

	// Error should take priority - loading is checked first in renderContent
	// but SetError sets loading to false, so let's test the actual priority
	if strings.Contains(content, "Loading") {
		t.Error("expected error to take priority over loading state")
	}
}
