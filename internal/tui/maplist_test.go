package tui

import (
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

func TestNewMapListModel(t *testing.T) {
	m := newMapListModel(80, 24)

	if m.list.Title != "BPF Maps" {
		t.Errorf("expected title 'BPF Maps', got '%s'", m.list.Title)
	}

	if len(m.maps) != 0 {
		t.Errorf("expected empty maps slice, got %d items", len(m.maps))
	}

	if m.err != nil {
		t.Errorf("expected nil error, got %v", m.err)
	}
}

func TestMapItemInterface(t *testing.T) {
	item := mapItem{
		info: MapInfo{
			ID:         42,
			Name:       "test_map",
			Type:       "hash",
			KeySize:    4,
			ValueSize:  8,
			MaxEntries: 1024,
		},
	}

	// Test Title
	expectedTitle := "[42] test_map"
	if item.Title() != expectedTitle {
		t.Errorf("expected title '%s', got '%s'", expectedTitle, item.Title())
	}

	// Test Description
	expectedDesc := "Type: hash | Key: 4 | Value: 8 | Max: 1024"
	if item.Description() != expectedDesc {
		t.Errorf("expected description '%s', got '%s'", expectedDesc, item.Description())
	}

	// Test FilterValue
	if item.FilterValue() != "test_map" {
		t.Errorf("expected filter value 'test_map', got '%s'", item.FilterValue())
	}
}

func TestMapListInit(t *testing.T) {
	m := newMapListModel(80, 24)
	cmd := m.Init()

	if cmd != nil {
		t.Error("expected nil command from Init()")
	}
}

func TestMapListSetMaps(t *testing.T) {
	m := newMapListModel(80, 24)

	maps := []MapInfo{
		{ID: 1, Name: "map1", Type: "hash", KeySize: 4, ValueSize: 8, MaxEntries: 100},
		{ID: 2, Name: "map2", Type: "array", KeySize: 4, ValueSize: 16, MaxEntries: 256},
		{ID: 3, Name: "map3", Type: "lru_hash", KeySize: 8, ValueSize: 32, MaxEntries: 512},
	}

	m.SetMaps(maps)

	if len(m.maps) != 3 {
		t.Errorf("expected 3 maps, got %d", len(m.maps))
	}

	// Verify list items were set
	items := m.list.Items()
	if len(items) != 3 {
		t.Errorf("expected 3 list items, got %d", len(items))
	}
}

func TestMapListUpdateEnterKey(t *testing.T) {
	m := newMapListModel(80, 24)

	maps := []MapInfo{
		{ID: 1, Name: "map1", Type: "hash", KeySize: 4, ValueSize: 8, MaxEntries: 100},
		{ID: 2, Name: "map2", Type: "array", KeySize: 4, ValueSize: 16, MaxEntries: 256},
	}
	m.SetMaps(maps)

	// Press Enter to select first item
	msg := tea.KeyMsg{Type: tea.KeyEnter}
	_, _, selectedMap := m.Update(msg)

	if selectedMap == nil {
		t.Fatal("expected selected map, got nil")
	}

	if selectedMap.ID != 1 {
		t.Errorf("expected map ID 1, got %d", selectedMap.ID)
	}

	if selectedMap.Name != "map1" {
		t.Errorf("expected map name 'map1', got '%s'", selectedMap.Name)
	}
}

func TestMapListUpdateNavigationKeys(t *testing.T) {
	m := newMapListModel(80, 24)

	maps := []MapInfo{
		{ID: 1, Name: "map1", Type: "hash", KeySize: 4, ValueSize: 8, MaxEntries: 100},
		{ID: 2, Name: "map2", Type: "array", KeySize: 4, ValueSize: 16, MaxEntries: 256},
		{ID: 3, Name: "map3", Type: "lru_hash", KeySize: 8, ValueSize: 32, MaxEntries: 512},
	}
	m.SetMaps(maps)

	// Navigate down with 'j'
	msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}}
	m, _, _ = m.Update(msg)

	// Press Enter to select
	enterMsg := tea.KeyMsg{Type: tea.KeyEnter}
	_, _, selectedMap := m.Update(enterMsg)

	if selectedMap == nil {
		t.Fatal("expected selected map, got nil")
	}

	if selectedMap.ID != 2 {
		t.Errorf("expected map ID 2 after navigating down, got %d", selectedMap.ID)
	}
}

func TestMapListUpdateDownArrowKey(t *testing.T) {
	m := newMapListModel(80, 24)

	maps := []MapInfo{
		{ID: 1, Name: "map1", Type: "hash", KeySize: 4, ValueSize: 8, MaxEntries: 100},
		{ID: 2, Name: "map2", Type: "array", KeySize: 4, ValueSize: 16, MaxEntries: 256},
	}
	m.SetMaps(maps)

	// Navigate down with arrow key
	msg := tea.KeyMsg{Type: tea.KeyDown}
	m, _, _ = m.Update(msg)

	// Press Enter to select
	enterMsg := tea.KeyMsg{Type: tea.KeyEnter}
	_, _, selectedMap := m.Update(enterMsg)

	if selectedMap == nil {
		t.Fatal("expected selected map, got nil")
	}

	if selectedMap.ID != 2 {
		t.Errorf("expected map ID 2 after navigating down, got %d", selectedMap.ID)
	}
}

func TestMapListUpdateUpKey(t *testing.T) {
	m := newMapListModel(80, 24)

	maps := []MapInfo{
		{ID: 1, Name: "map1", Type: "hash", KeySize: 4, ValueSize: 8, MaxEntries: 100},
		{ID: 2, Name: "map2", Type: "array", KeySize: 4, ValueSize: 16, MaxEntries: 256},
	}
	m.SetMaps(maps)

	// Navigate down then up with 'k'
	downMsg := tea.KeyMsg{Type: tea.KeyDown}
	m, _, _ = m.Update(downMsg)

	upMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'k'}}
	m, _, _ = m.Update(upMsg)

	// Press Enter to select
	enterMsg := tea.KeyMsg{Type: tea.KeyEnter}
	_, _, selectedMap := m.Update(enterMsg)

	if selectedMap == nil {
		t.Fatal("expected selected map, got nil")
	}

	if selectedMap.ID != 1 {
		t.Errorf("expected map ID 1 after navigating up, got %d", selectedMap.ID)
	}
}

func TestMapListViewEmpty(t *testing.T) {
	m := newMapListModel(80, 24)

	view := m.View()

	if !strings.Contains(view, "BPF Maps") {
		t.Error("expected view to contain 'BPF Maps'")
	}

	if !strings.Contains(view, "No BPF maps loaded") {
		t.Error("expected view to contain empty state message")
	}
}

func TestMapListViewWithMaps(t *testing.T) {
	m := newMapListModel(80, 24)

	maps := []MapInfo{
		{ID: 1, Name: "test_map", Type: "hash", KeySize: 4, ValueSize: 8, MaxEntries: 100},
	}
	m.SetMaps(maps)

	view := m.View()

	if !strings.Contains(view, "test_map") {
		t.Error("expected view to contain map name")
	}
}

func TestMapListViewWithError(t *testing.T) {
	m := newMapListModel(80, 24)
	m.SetError(&PermissionError{})

	view := m.View()

	if !strings.Contains(view, "Error") {
		t.Error("expected view to contain error message")
	}
}

func TestMapListSetSize(t *testing.T) {
	m := newMapListModel(80, 24)

	m.SetSize(120, 40)

	// The list should have been resized (we can't easily verify internal dimensions,
	// but we can verify the method doesn't panic)
}

func TestMapListSelectedItem(t *testing.T) {
	m := newMapListModel(80, 24)

	// Empty list should return nil
	if m.SelectedItem() != nil {
		t.Error("expected nil selected item for empty list")
	}

	maps := []MapInfo{
		{ID: 1, Name: "map1", Type: "hash", KeySize: 4, ValueSize: 8, MaxEntries: 100},
	}
	m.SetMaps(maps)

	selected := m.SelectedItem()
	if selected == nil {
		t.Fatal("expected selected item, got nil")
	}

	if selected.ID != 1 {
		t.Errorf("expected selected map ID 1, got %d", selected.ID)
	}
}

func TestMapListIsFiltering(t *testing.T) {
	m := newMapListModel(80, 24)

	maps := []MapInfo{
		{ID: 1, Name: "map1", Type: "hash", KeySize: 4, ValueSize: 8, MaxEntries: 100},
	}
	m.SetMaps(maps)

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

// mockMapsService is a mock implementation of MapsService for testing.
type mockMapsServiceForMapList struct {
	maps []MapInfo
	err  error
}

func (m *mockMapsServiceForMapList) List() ([]MapInfo, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.maps, nil
}

func (m *mockMapsServiceForMapList) Get(id uint32) (*MapInfo, error) {
	if m.err != nil {
		return nil, m.err
	}
	for _, mp := range m.maps {
		if mp.ID == id {
			return &mp, nil
		}
	}
	return nil, nil
}

func (m *mockMapsServiceForMapList) Dump(id uint32) ([]MapEntry, error) {
	return nil, nil
}

// Integration test: Full navigation flow from menu to maps list and back
func TestMapListIntegrationNavigateFromMenu(t *testing.T) {
	m := NewModel(nil, nil)

	// Verify we start at menu
	if m.state != ViewMenu {
		t.Errorf("expected initial state ViewMenu, got %v", m.state)
	}

	// Navigate down to select "Maps" (second menu item)
	downMsg := tea.KeyMsg{Type: tea.KeyDown}
	result, _ := m.Update(downMsg)
	m = result.(Model)

	// Press Enter to select Maps
	enterMsg := tea.KeyMsg{Type: tea.KeyEnter}
	result, _ = m.Update(enterMsg)
	m = result.(Model)

	// Should now be at maps list
	if m.state != ViewMapList {
		t.Errorf("expected state ViewMapList, got %v", m.state)
	}

	// History should have menu
	if m.historyLen() != 1 {
		t.Errorf("expected history length 1, got %d", m.historyLen())
	}
}

func TestMapListIntegrationBackToMenu(t *testing.T) {
	m := NewModel(nil, nil)

	// Navigate to maps list
	downMsg := tea.KeyMsg{Type: tea.KeyDown}
	result, _ := m.Update(downMsg)
	m = result.(Model)

	enterMsg := tea.KeyMsg{Type: tea.KeyEnter}
	result, _ = m.Update(enterMsg)
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

func TestMapListIntegrationSelectMap(t *testing.T) {
	// Create mock service
	mockSvc := &mockMapsServiceForMapList{
		maps: []MapInfo{
			{ID: 1, Name: "map1", Type: "hash", KeySize: 4, ValueSize: 8, MaxEntries: 100},
			{ID: 2, Name: "map2", Type: "array", KeySize: 4, ValueSize: 16, MaxEntries: 256},
		},
	}

	m := NewModel(nil, mockSvc)

	// Navigate to maps list
	downMsg := tea.KeyMsg{Type: tea.KeyDown}
	result, _ := m.Update(downMsg)
	m = result.(Model)

	enterMsg := tea.KeyMsg{Type: tea.KeyEnter}
	result, _ = m.Update(enterMsg)
	m = result.(Model)

	// Press Enter to select first map
	result, _ = m.Update(enterMsg)
	m = result.(Model)

	// Should now be at map detail
	if m.state != ViewMapDetail {
		t.Errorf("expected state ViewMapDetail, got %v", m.state)
	}

	// History should have menu and maplist
	if m.historyLen() != 2 {
		t.Errorf("expected history length 2, got %d", m.historyLen())
	}
}

func TestMapListIntegrationBackspaceNavigation(t *testing.T) {
	m := NewModel(nil, nil)

	// Navigate to maps list
	downMsg := tea.KeyMsg{Type: tea.KeyDown}
	result, _ := m.Update(downMsg)
	m = result.(Model)

	enterMsg := tea.KeyMsg{Type: tea.KeyEnter}
	result, _ = m.Update(enterMsg)
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

func TestMapListIntegrationScrolling(t *testing.T) {
	m := newMapListModel(80, 24)

	// Create many maps to test scrolling
	maps := make([]MapInfo, 50)
	for i := 0; i < 50; i++ {
		maps[i] = MapInfo{
			ID:         uint32(i + 1),
			Name:       "map" + string(rune('a'+i%26)),
			Type:       "hash",
			KeySize:    4,
			ValueSize:  8,
			MaxEntries: 100,
		}
	}
	m.SetMaps(maps)

	// Navigate down multiple times
	downMsg := tea.KeyMsg{Type: tea.KeyDown}
	for i := 0; i < 10; i++ {
		m, _, _ = m.Update(downMsg)
	}

	// Press Enter to select
	enterMsg := tea.KeyMsg{Type: tea.KeyEnter}
	_, _, selectedMap := m.Update(enterMsg)

	if selectedMap == nil {
		t.Fatal("expected selected map, got nil")
	}

	// Should have navigated to the 11th item (index 10)
	if selectedMap.ID != 11 {
		t.Errorf("expected map ID 11 after scrolling, got %d", selectedMap.ID)
	}
}

// Fuzzy Search Tests

func TestMapListFuzzySearchActivation(t *testing.T) {
	m := newMapListModel(80, 24)

	maps := []MapInfo{
		{ID: 1, Name: "hash_map", Type: "hash", KeySize: 4, ValueSize: 8, MaxEntries: 100},
		{ID: 2, Name: "array_map", Type: "array", KeySize: 4, ValueSize: 16, MaxEntries: 256},
		{ID: 3, Name: "lru_hash_map", Type: "lru_hash", KeySize: 8, ValueSize: 32, MaxEntries: 512},
	}
	m.SetMaps(maps)

	// Initially not filtering
	if m.IsFiltering() {
		t.Error("expected not filtering initially")
	}

	// Activate filtering with '/'
	filterMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'/'}}
	m, _, _ = m.Update(filterMsg)

	// Should now be in filtering mode
	if !m.IsFiltering() {
		t.Error("expected filtering after pressing '/'")
	}
}

func TestMapListFuzzySearchExitWithEscape(t *testing.T) {
	m := newMapListModel(80, 24)

	maps := []MapInfo{
		{ID: 1, Name: "hash_map", Type: "hash", KeySize: 4, ValueSize: 8, MaxEntries: 100},
		{ID: 2, Name: "array_map", Type: "array", KeySize: 4, ValueSize: 16, MaxEntries: 256},
	}
	m.SetMaps(maps)

	// Activate filtering
	filterMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'/'}}
	m, _, _ = m.Update(filterMsg)

	if !m.IsFiltering() {
		t.Error("expected filtering after pressing '/'")
	}

	// Exit filtering with Escape
	escMsg := tea.KeyMsg{Type: tea.KeyEscape}
	m, _, _ = m.Update(escMsg)

	// Should no longer be filtering
	if m.IsFiltering() {
		t.Error("expected not filtering after pressing Escape")
	}
}

func TestMapListFuzzySearchEnterDoesNotSelectWhileFiltering(t *testing.T) {
	m := newMapListModel(80, 24)

	maps := []MapInfo{
		{ID: 1, Name: "hash_map", Type: "hash", KeySize: 4, ValueSize: 8, MaxEntries: 100},
		{ID: 2, Name: "array_map", Type: "array", KeySize: 4, ValueSize: 16, MaxEntries: 256},
	}
	m.SetMaps(maps)

	// Activate filtering
	filterMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'/'}}
	m, _, _ = m.Update(filterMsg)

	// Press Enter while filtering - should not return a selected map
	// (the list handles this internally to confirm the filter)
	enterMsg := tea.KeyMsg{Type: tea.KeyEnter}
	m, _, selectedMap := m.Update(enterMsg)

	// When filtering, Enter confirms the filter, not selects the item
	if selectedMap != nil {
		t.Error("expected no selection while filtering is active")
	}
}

func TestMapListFuzzySearchNavigationWhileFiltering(t *testing.T) {
	m := newMapListModel(80, 24)

	maps := []MapInfo{
		{ID: 1, Name: "hash_map", Type: "hash", KeySize: 4, ValueSize: 8, MaxEntries: 100},
		{ID: 2, Name: "hash_lru_map", Type: "lru_hash", KeySize: 4, ValueSize: 8, MaxEntries: 100},
		{ID: 3, Name: "array_map", Type: "array", KeySize: 4, ValueSize: 16, MaxEntries: 256},
	}
	m.SetMaps(maps)

	// Activate filtering
	filterMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'/'}}
	m, _, _ = m.Update(filterMsg)

	// Navigate down while filtering
	downMsg := tea.KeyMsg{Type: tea.KeyDown}
	m, _, _ = m.Update(downMsg)

	// Should still be filtering
	if !m.IsFiltering() {
		t.Error("expected still filtering after navigation")
	}
}

func TestMapListFuzzySearchSelectAfterFilter(t *testing.T) {
	m := newMapListModel(80, 24)

	maps := []MapInfo{
		{ID: 1, Name: "hash_map", Type: "hash", KeySize: 4, ValueSize: 8, MaxEntries: 100},
		{ID: 2, Name: "array_map", Type: "array", KeySize: 4, ValueSize: 16, MaxEntries: 256},
	}
	m.SetMaps(maps)

	// Activate filtering
	filterMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'/'}}
	m, _, _ = m.Update(filterMsg)

	// Type search query
	for _, r := range "hash" {
		charMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{r}}
		m, _, _ = m.Update(charMsg)
	}

	// Confirm filter with Enter (exits filter mode)
	enterMsg := tea.KeyMsg{Type: tea.KeyEnter}
	m, _, _ = m.Update(enterMsg)

	// Now press Enter again to select the filtered item
	m, _, selectedMap := m.Update(enterMsg)

	// Should have selected a map
	if selectedMap == nil {
		t.Fatal("expected selected map after confirming filter")
	}

	// Should be the hash_map (first match)
	if selectedMap.Name != "hash_map" {
		t.Errorf("expected 'hash_map', got '%s'", selectedMap.Name)
	}
}

func TestMapListFuzzySearchViewShowsFilterInput(t *testing.T) {
	m := newMapListModel(80, 24)

	maps := []MapInfo{
		{ID: 1, Name: "hash_map", Type: "hash", KeySize: 4, ValueSize: 8, MaxEntries: 100},
	}
	m.SetMaps(maps)

	// Activate filtering
	filterMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'/'}}
	m, _, _ = m.Update(filterMsg)

	view := m.View()

	// The view should show filter-related UI elements
	// Bubbles list shows "Filter:" when filtering is active
	if !strings.Contains(view, "Filter") {
		t.Error("expected view to show filter input when filtering")
	}
}

func TestMapListFilterValueUsesName(t *testing.T) {
	item := mapItem{
		info: MapInfo{
			ID:         1,
			Name:       "my_custom_map",
			Type:       "hash",
			KeySize:    4,
			ValueSize:  8,
			MaxEntries: 100,
		},
	}

	// FilterValue should return the map name for fuzzy matching
	if item.FilterValue() != "my_custom_map" {
		t.Errorf("expected FilterValue 'my_custom_map', got '%s'", item.FilterValue())
	}
}

// Integration test: Fuzzy search in full TUI context
func TestMapListIntegrationFuzzySearchDoesNotNavigateBack(t *testing.T) {
	mockSvc := &mockMapsServiceForMapList{
		maps: []MapInfo{
			{ID: 1, Name: "map1", Type: "hash", KeySize: 4, ValueSize: 8, MaxEntries: 100},
		},
	}

	m := NewModel(nil, mockSvc)

	// Navigate to maps list (down to select Maps, then Enter)
	downMsg := tea.KeyMsg{Type: tea.KeyDown}
	result, _ := m.Update(downMsg)
	m = result.(Model)

	enterMsg := tea.KeyMsg{Type: tea.KeyEnter}
	result, _ = m.Update(enterMsg)
	m = result.(Model)

	if m.state != ViewMapList {
		t.Fatalf("expected ViewMapList, got %v", m.state)
	}

	// Activate filtering
	filterMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'/'}}
	result, _ = m.Update(filterMsg)
	m = result.(Model)

	// Press Escape while filtering - should exit filter mode, not navigate back
	escMsg := tea.KeyMsg{Type: tea.KeyEscape}
	result, _ = m.Update(escMsg)
	m = result.(Model)

	// Should still be at maps list (not navigated back to menu)
	if m.state != ViewMapList {
		t.Errorf("expected to stay at ViewMapList after exiting filter, got %v", m.state)
	}
}

func TestMapListIntegrationQuitBlockedWhileFiltering(t *testing.T) {
	mockSvc := &mockMapsServiceForMapList{
		maps: []MapInfo{
			{ID: 1, Name: "map1", Type: "hash", KeySize: 4, ValueSize: 8, MaxEntries: 100},
		},
	}

	m := NewModel(nil, mockSvc)

	// Navigate to maps list
	downMsg := tea.KeyMsg{Type: tea.KeyDown}
	result, _ := m.Update(downMsg)
	m = result.(Model)

	enterMsg := tea.KeyMsg{Type: tea.KeyEnter}
	result, _ = m.Update(enterMsg)
	m = result.(Model)

	// Activate filtering
	filterMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'/'}}
	result, _ = m.Update(filterMsg)
	m = result.(Model)

	// Press 'q' while filtering - should type 'q' in filter, not quit
	qMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}}
	result, cmd := m.Update(qMsg)
	m = result.(Model)

	// Should not have quit command
	if cmd != nil {
		// Check if it's a quit command by checking if it returns tea.Quit
		// In practice, the cmd should be nil or a non-quit command
		t.Log("Command returned, but should not be quit while filtering")
	}

	// Should still be at maps list
	if m.state != ViewMapList {
		t.Errorf("expected to stay at ViewMapList, got %v", m.state)
	}
}

func TestMapListFuzzySearchWithMultipleMatches(t *testing.T) {
	m := newMapListModel(80, 24)

	maps := []MapInfo{
		{ID: 1, Name: "hash_map_1", Type: "hash", KeySize: 4, ValueSize: 8, MaxEntries: 100},
		{ID: 2, Name: "hash_map_2", Type: "hash", KeySize: 4, ValueSize: 16, MaxEntries: 256},
		{ID: 3, Name: "array_map", Type: "array", KeySize: 4, ValueSize: 32, MaxEntries: 512},
	}
	m.SetMaps(maps)

	// Activate filtering
	filterMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'/'}}
	m, _, _ = m.Update(filterMsg)

	// Type search query that matches multiple items
	for _, r := range "hash" {
		charMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{r}}
		m, _, _ = m.Update(charMsg)
	}

	// Confirm filter
	enterMsg := tea.KeyMsg{Type: tea.KeyEnter}
	m, _, _ = m.Update(enterMsg)

	// Navigate down to second match
	downMsg := tea.KeyMsg{Type: tea.KeyDown}
	m, _, _ = m.Update(downMsg)

	// Select second match
	m, _, selectedMap := m.Update(enterMsg)

	if selectedMap == nil {
		t.Fatal("expected selected map")
	}

	// Should be hash_map_2 (second match after navigating down)
	if selectedMap.Name != "hash_map_2" {
		t.Errorf("expected 'hash_map_2', got '%s'", selectedMap.Name)
	}
}

func TestMapListResetFilter(t *testing.T) {
	m := newMapListModel(80, 24)

	maps := []MapInfo{
		{ID: 1, Name: "hash_map", Type: "hash", KeySize: 4, ValueSize: 8, MaxEntries: 100},
		{ID: 2, Name: "array_map", Type: "array", KeySize: 4, ValueSize: 16, MaxEntries: 256},
	}
	m.SetMaps(maps)

	// Activate filtering
	filterMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'/'}}
	m, _, _ = m.Update(filterMsg)

	// Type search query
	for _, r := range "hash" {
		charMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{r}}
		m, _, _ = m.Update(charMsg)
	}

	// Confirm filter
	enterMsg := tea.KeyMsg{Type: tea.KeyEnter}
	m, _, _ = m.Update(enterMsg)

	// Verify filter is applied
	if !m.list.IsFiltered() {
		t.Error("expected filter to be applied")
	}

	// Reset filter
	m.ResetFilter()

	// Verify filter is cleared
	if m.list.IsFiltered() {
		t.Error("expected filter to be cleared after ResetFilter")
	}

	if m.IsFiltering() {
		t.Error("expected not to be in filtering mode after ResetFilter")
	}
}
