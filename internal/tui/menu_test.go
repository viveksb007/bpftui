package tui

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

func TestNewMenuModel(t *testing.T) {
	m := newMenuModel(80, 24)

	// Verify menu has correct number of items
	if m.list.Items() == nil {
		t.Fatal("expected menu items to be initialized")
	}

	items := m.list.Items()
	if len(items) != 2 {
		t.Errorf("expected 2 menu items, got %d", len(items))
	}

	// Verify first item is Programs
	if item, ok := items[0].(menuItem); ok {
		if item.title != "Programs" {
			t.Errorf("expected first item title 'Programs', got '%s'", item.title)
		}
		if item.target != ViewProgList {
			t.Errorf("expected first item target ViewProgList, got %v", item.target)
		}
	} else {
		t.Error("first item is not a menuItem")
	}

	// Verify second item is Maps
	if item, ok := items[1].(menuItem); ok {
		if item.title != "Maps" {
			t.Errorf("expected second item title 'Maps', got '%s'", item.title)
		}
		if item.target != ViewMapList {
			t.Errorf("expected second item target ViewMapList, got %v", item.target)
		}
	} else {
		t.Error("second item is not a menuItem")
	}
}

func TestMenuItemInterface(t *testing.T) {
	item := menuItem{
		title:       "Test",
		description: "Test description",
		target:      ViewProgList,
	}

	if item.Title() != "Test" {
		t.Errorf("expected Title() 'Test', got '%s'", item.Title())
	}

	if item.Description() != "Test description" {
		t.Errorf("expected Description() 'Test description', got '%s'", item.Description())
	}

	if item.FilterValue() != "Test" {
		t.Errorf("expected FilterValue() 'Test', got '%s'", item.FilterValue())
	}
}

func TestMenuInit(t *testing.T) {
	m := newMenuModel(80, 24)
	cmd := m.Init()

	if cmd != nil {
		t.Error("expected Init() to return nil command")
	}
}

func TestMenuUpdateEnterKey(t *testing.T) {
	m := newMenuModel(80, 24)

	// First item (Programs) should be selected by default
	msg := tea.KeyMsg{Type: tea.KeyEnter}
	updatedMenu, _, targetView := m.Update(msg)

	if targetView == nil {
		t.Fatal("expected targetView to be set on Enter key")
	}

	if *targetView != ViewProgList {
		t.Errorf("expected target ViewProgList, got %v", *targetView)
	}

	// Verify menu state is preserved
	if updatedMenu.list.Items() == nil {
		t.Error("menu items should be preserved after update")
	}
}

func TestMenuUpdateDownThenEnter(t *testing.T) {
	m := newMenuModel(80, 24)

	// Move down to Maps
	downMsg := tea.KeyMsg{Type: tea.KeyDown}
	m, _, _ = m.Update(downMsg)

	// Press Enter to select Maps
	enterMsg := tea.KeyMsg{Type: tea.KeyEnter}
	_, _, targetView := m.Update(enterMsg)

	if targetView == nil {
		t.Fatal("expected targetView to be set on Enter key")
	}

	if *targetView != ViewMapList {
		t.Errorf("expected target ViewMapList after moving down, got %v", *targetView)
	}
}

func TestMenuUpdateNavigationKeys(t *testing.T) {
	m := newMenuModel(80, 24)

	// Test down arrow
	downMsg := tea.KeyMsg{Type: tea.KeyDown}
	m, _, targetView := m.Update(downMsg)
	if targetView != nil {
		t.Error("down key should not trigger navigation")
	}

	// Test up arrow
	upMsg := tea.KeyMsg{Type: tea.KeyUp}
	m, _, targetView = m.Update(upMsg)
	if targetView != nil {
		t.Error("up key should not trigger navigation")
	}

	// Test j key (vim down)
	jMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}}
	m, _, targetView = m.Update(jMsg)
	if targetView != nil {
		t.Error("j key should not trigger navigation")
	}

	// Test k key (vim up)
	kMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'k'}}
	_, _, targetView = m.Update(kMsg)
	if targetView != nil {
		t.Error("k key should not trigger navigation")
	}
}

func TestMenuView(t *testing.T) {
	m := newMenuModel(80, 24)
	view := m.View()

	if view == "" {
		t.Error("expected non-empty view output")
	}

	// View should contain the title
	if !containsString(view, "BPF TUI Explorer") {
		t.Error("view should contain title 'BPF TUI Explorer'")
	}
}

func TestMenuSetSize(t *testing.T) {
	m := newMenuModel(80, 24)

	// Change size
	m.SetSize(120, 40)

	// Verify the list was resized (we can't directly check dimensions,
	// but we can verify it doesn't panic and still renders)
	view := m.View()
	if view == "" {
		t.Error("expected non-empty view after resize")
	}
}

func TestMenuSelectedItem(t *testing.T) {
	m := newMenuModel(80, 24)

	// First item should be selected by default
	item := m.SelectedItem()
	if item == nil {
		t.Fatal("expected selected item to be non-nil")
	}

	if item.title != "Programs" {
		t.Errorf("expected selected item 'Programs', got '%s'", item.title)
	}

	// Move down and check selection
	downMsg := tea.KeyMsg{Type: tea.KeyDown}
	m, _, _ = m.Update(downMsg)

	item = m.SelectedItem()
	if item == nil {
		t.Fatal("expected selected item to be non-nil after moving down")
	}

	if item.title != "Maps" {
		t.Errorf("expected selected item 'Maps' after moving down, got '%s'", item.title)
	}
}

// Integration tests for menu with root model

func TestMenuIntegrationSelectPrograms(t *testing.T) {
	m := NewModel(nil, nil)

	// Verify we start at menu
	if m.state != ViewMenu {
		t.Fatalf("initial state = %v, want ViewMenu", m.state)
	}

	// Press Enter to select Programs (first item)
	enterMsg := tea.KeyMsg{Type: tea.KeyEnter}
	newModel, _ := m.Update(enterMsg)
	updated := newModel.(Model)

	if updated.state != ViewProgList {
		t.Errorf("state after selecting Programs = %v, want ViewProgList", updated.state)
	}

	// History should contain ViewMenu
	if updated.historyLen() != 1 {
		t.Errorf("history length = %d, want 1", updated.historyLen())
	}
}

func TestMenuIntegrationSelectMaps(t *testing.T) {
	m := NewModel(nil, nil)

	// Move down to Maps
	downMsg := tea.KeyMsg{Type: tea.KeyDown}
	newModel, _ := m.Update(downMsg)
	m = newModel.(Model)

	// Press Enter to select Maps
	enterMsg := tea.KeyMsg{Type: tea.KeyEnter}
	newModel, _ = m.Update(enterMsg)
	updated := newModel.(Model)

	if updated.state != ViewMapList {
		t.Errorf("state after selecting Maps = %v, want ViewMapList", updated.state)
	}

	// History should contain ViewMenu
	if updated.historyLen() != 1 {
		t.Errorf("history length = %d, want 1", updated.historyLen())
	}
}

func TestMenuIntegrationBackFromProgList(t *testing.T) {
	m := NewModel(nil, nil)

	// Navigate to Programs
	enterMsg := tea.KeyMsg{Type: tea.KeyEnter}
	newModel, _ := m.Update(enterMsg)
	m = newModel.(Model)

	if m.state != ViewProgList {
		t.Fatalf("state = %v, want ViewProgList", m.state)
	}

	// Press Esc to go back
	escMsg := tea.KeyMsg{Type: tea.KeyEsc}
	newModel, _ = m.Update(escMsg)
	updated := newModel.(Model)

	if updated.state != ViewMenu {
		t.Errorf("state after back = %v, want ViewMenu", updated.state)
	}

	if updated.historyLen() != 0 {
		t.Errorf("history length after back = %d, want 0", updated.historyLen())
	}
}

func TestMenuIntegrationQuitFromMenu(t *testing.T) {
	m := NewModel(nil, nil)

	// Press q to quit
	quitMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}}
	_, cmd := m.Update(quitMsg)

	if cmd == nil {
		t.Error("expected quit command, got nil")
	}
}
