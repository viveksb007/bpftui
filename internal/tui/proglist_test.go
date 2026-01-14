package tui

import (
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

func TestNewProgListModel(t *testing.T) {
	m := newProgListModel(80, 24)

	if m.list.Title != "BPF Programs" {
		t.Errorf("expected title 'BPF Programs', got '%s'", m.list.Title)
	}

	if len(m.programs) != 0 {
		t.Errorf("expected empty programs slice, got %d items", len(m.programs))
	}

	if m.err != nil {
		t.Errorf("expected nil error, got %v", m.err)
	}
}

func TestProgItemInterface(t *testing.T) {
	item := progItem{
		info: ProgramInfo{
			ID:   42,
			Name: "test_prog",
			Type: "kprobe",
			Tag:  "abc123",
		},
	}

	// Test Title
	expectedTitle := "[42] test_prog"
	if item.Title() != expectedTitle {
		t.Errorf("expected title '%s', got '%s'", expectedTitle, item.Title())
	}

	// Test Description
	expectedDesc := "Type: kprobe | Tag: abc123"
	if item.Description() != expectedDesc {
		t.Errorf("expected description '%s', got '%s'", expectedDesc, item.Description())
	}

	// Test FilterValue
	if item.FilterValue() != "test_prog" {
		t.Errorf("expected filter value 'test_prog', got '%s'", item.FilterValue())
	}
}

func TestProgListInit(t *testing.T) {
	m := newProgListModel(80, 24)
	cmd := m.Init()

	if cmd != nil {
		t.Error("expected nil command from Init()")
	}
}

func TestProgListSetPrograms(t *testing.T) {
	m := newProgListModel(80, 24)

	programs := []ProgramInfo{
		{ID: 1, Name: "prog1", Type: "kprobe", Tag: "tag1"},
		{ID: 2, Name: "prog2", Type: "tracepoint", Tag: "tag2"},
		{ID: 3, Name: "prog3", Type: "xdp", Tag: "tag3"},
	}

	m.SetPrograms(programs)

	if len(m.programs) != 3 {
		t.Errorf("expected 3 programs, got %d", len(m.programs))
	}

	// Verify list items were set
	items := m.list.Items()
	if len(items) != 3 {
		t.Errorf("expected 3 list items, got %d", len(items))
	}
}

func TestProgListUpdateEnterKey(t *testing.T) {
	m := newProgListModel(80, 24)

	programs := []ProgramInfo{
		{ID: 1, Name: "prog1", Type: "kprobe", Tag: "tag1"},
		{ID: 2, Name: "prog2", Type: "tracepoint", Tag: "tag2"},
	}
	m.SetPrograms(programs)

	// Press Enter to select first item
	msg := tea.KeyMsg{Type: tea.KeyEnter}
	_, _, selectedProg := m.Update(msg)

	if selectedProg == nil {
		t.Fatal("expected selected program, got nil")
	}

	if selectedProg.ID != 1 {
		t.Errorf("expected program ID 1, got %d", selectedProg.ID)
	}

	if selectedProg.Name != "prog1" {
		t.Errorf("expected program name 'prog1', got '%s'", selectedProg.Name)
	}
}

func TestProgListUpdateNavigationKeys(t *testing.T) {
	m := newProgListModel(80, 24)

	programs := []ProgramInfo{
		{ID: 1, Name: "prog1", Type: "kprobe", Tag: "tag1"},
		{ID: 2, Name: "prog2", Type: "tracepoint", Tag: "tag2"},
		{ID: 3, Name: "prog3", Type: "xdp", Tag: "tag3"},
	}
	m.SetPrograms(programs)

	// Navigate down with 'j'
	msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}}
	m, _, _ = m.Update(msg)

	// Press Enter to select
	enterMsg := tea.KeyMsg{Type: tea.KeyEnter}
	_, _, selectedProg := m.Update(enterMsg)

	if selectedProg == nil {
		t.Fatal("expected selected program, got nil")
	}

	if selectedProg.ID != 2 {
		t.Errorf("expected program ID 2 after navigating down, got %d", selectedProg.ID)
	}
}

func TestProgListUpdateDownArrowKey(t *testing.T) {
	m := newProgListModel(80, 24)

	programs := []ProgramInfo{
		{ID: 1, Name: "prog1", Type: "kprobe", Tag: "tag1"},
		{ID: 2, Name: "prog2", Type: "tracepoint", Tag: "tag2"},
	}
	m.SetPrograms(programs)

	// Navigate down with arrow key
	msg := tea.KeyMsg{Type: tea.KeyDown}
	m, _, _ = m.Update(msg)

	// Press Enter to select
	enterMsg := tea.KeyMsg{Type: tea.KeyEnter}
	_, _, selectedProg := m.Update(enterMsg)

	if selectedProg == nil {
		t.Fatal("expected selected program, got nil")
	}

	if selectedProg.ID != 2 {
		t.Errorf("expected program ID 2 after navigating down, got %d", selectedProg.ID)
	}
}

func TestProgListUpdateUpKey(t *testing.T) {
	m := newProgListModel(80, 24)

	programs := []ProgramInfo{
		{ID: 1, Name: "prog1", Type: "kprobe", Tag: "tag1"},
		{ID: 2, Name: "prog2", Type: "tracepoint", Tag: "tag2"},
	}
	m.SetPrograms(programs)

	// Navigate down then up with 'k'
	downMsg := tea.KeyMsg{Type: tea.KeyDown}
	m, _, _ = m.Update(downMsg)

	upMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'k'}}
	m, _, _ = m.Update(upMsg)

	// Press Enter to select
	enterMsg := tea.KeyMsg{Type: tea.KeyEnter}
	_, _, selectedProg := m.Update(enterMsg)

	if selectedProg == nil {
		t.Fatal("expected selected program, got nil")
	}

	if selectedProg.ID != 1 {
		t.Errorf("expected program ID 1 after navigating up, got %d", selectedProg.ID)
	}
}

func TestProgListViewEmpty(t *testing.T) {
	m := newProgListModel(80, 24)

	view := m.View()

	if !strings.Contains(view, "BPF Programs") {
		t.Error("expected view to contain 'BPF Programs'")
	}

	if !strings.Contains(view, "No BPF programs loaded") {
		t.Error("expected view to contain empty state message")
	}
}

func TestProgListViewWithPrograms(t *testing.T) {
	m := newProgListModel(80, 24)

	programs := []ProgramInfo{
		{ID: 1, Name: "test_prog", Type: "kprobe", Tag: "abc123"},
	}
	m.SetPrograms(programs)

	view := m.View()

	if !strings.Contains(view, "test_prog") {
		t.Error("expected view to contain program name")
	}
}

func TestProgListViewWithError(t *testing.T) {
	m := newProgListModel(80, 24)
	m.SetError(&PermissionError{})

	view := m.View()

	if !strings.Contains(view, "Error") {
		t.Error("expected view to contain error message")
	}
}

func TestProgListSetSize(t *testing.T) {
	m := newProgListModel(80, 24)

	m.SetSize(120, 40)

	// The list should have been resized (we can't easily verify internal dimensions,
	// but we can verify the method doesn't panic)
}

func TestProgListSelectedItem(t *testing.T) {
	m := newProgListModel(80, 24)

	// Empty list should return nil
	if m.SelectedItem() != nil {
		t.Error("expected nil selected item for empty list")
	}

	programs := []ProgramInfo{
		{ID: 1, Name: "prog1", Type: "kprobe", Tag: "tag1"},
	}
	m.SetPrograms(programs)

	selected := m.SelectedItem()
	if selected == nil {
		t.Fatal("expected selected item, got nil")
	}

	if selected.ID != 1 {
		t.Errorf("expected selected program ID 1, got %d", selected.ID)
	}
}

func TestProgListIsFiltering(t *testing.T) {
	m := newProgListModel(80, 24)

	programs := []ProgramInfo{
		{ID: 1, Name: "prog1", Type: "kprobe", Tag: "tag1"},
	}
	m.SetPrograms(programs)

	// Initially not filtering
	if m.IsFiltering() {
		t.Error("expected not filtering initially")
	}

	// Activate filtering with '/'
	filterMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'/'}}
	m, _, _ = m.Update(filterMsg)

	if !m.IsFiltering() {
		t.Error("expected filtering after pressing '/'")
	}
}

// Integration test: Full navigation flow from menu to programs list and back
func TestProgListIntegrationNavigateFromMenu(t *testing.T) {
	m := NewModel(nil, nil)

	// Verify we start at menu
	if m.state != ViewMenu {
		t.Errorf("expected initial state ViewMenu, got %v", m.state)
	}

	// Press Enter to select Programs (first menu item)
	enterMsg := tea.KeyMsg{Type: tea.KeyEnter}
	result, _ := m.Update(enterMsg)
	m = result.(Model)

	// Should now be at programs list
	if m.state != ViewProgList {
		t.Errorf("expected state ViewProgList, got %v", m.state)
	}

	// History should have menu
	if m.historyLen() != 1 {
		t.Errorf("expected history length 1, got %d", m.historyLen())
	}
}

func TestProgListIntegrationBackToMenu(t *testing.T) {
	m := NewModel(nil, nil)

	// Navigate to programs list
	enterMsg := tea.KeyMsg{Type: tea.KeyEnter}
	result, _ := m.Update(enterMsg)
	m = result.(Model)

	// Press Escape to go back
	escMsg := tea.KeyMsg{Type: tea.KeyEscape}
	result, _ = m.Update(escMsg)
	m = result.(Model)

	// Should be back at menu
	if m.state != ViewMenu {
		t.Errorf("expected state ViewMenu after back, got %v", m.state)
	}

	// History should be empty
	if m.historyLen() != 0 {
		t.Errorf("expected empty history, got %d", m.historyLen())
	}
}

func TestProgListIntegrationSelectProgram(t *testing.T) {
	// Create mock service
	mockSvc := &mockProgService{
		programs: []ProgramInfo{
			{ID: 1, Name: "prog1", Type: "kprobe", Tag: "tag1"},
			{ID: 2, Name: "prog2", Type: "tracepoint", Tag: "tag2"},
		},
	}

	m := NewModel(mockSvc, nil)

	// Navigate to programs list
	enterMsg := tea.KeyMsg{Type: tea.KeyEnter}
	result, _ := m.Update(enterMsg)
	m = result.(Model)

	// Press Enter to select first program
	result, _ = m.Update(enterMsg)
	m = result.(Model)

	// Should now be at program detail
	if m.state != ViewProgDetail {
		t.Errorf("expected state ViewProgDetail, got %v", m.state)
	}

	// History should have menu and proglist
	if m.historyLen() != 2 {
		t.Errorf("expected history length 2, got %d", m.historyLen())
	}
}

func TestProgListIntegrationBackspaceNavigation(t *testing.T) {
	m := NewModel(nil, nil)

	// Navigate to programs list
	enterMsg := tea.KeyMsg{Type: tea.KeyEnter}
	result, _ := m.Update(enterMsg)
	m = result.(Model)

	// Press Backspace to go back
	backspaceMsg := tea.KeyMsg{Type: tea.KeyBackspace}
	result, _ = m.Update(backspaceMsg)
	m = result.(Model)

	// Should be back at menu
	if m.state != ViewMenu {
		t.Errorf("expected state ViewMenu after backspace, got %v", m.state)
	}
}
