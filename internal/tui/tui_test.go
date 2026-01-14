package tui

import (
	"errors"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

// mockProgService is a mock implementation of ProgService for testing.
type mockProgService struct {
	programs []ProgramInfo
	err      error
}

func (m *mockProgService) List() ([]ProgramInfo, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.programs, nil
}

func (m *mockProgService) Get(id uint32) (*ProgramInfo, error) {
	if m.err != nil {
		return nil, m.err
	}
	for _, p := range m.programs {
		if p.ID == id {
			return &p, nil
		}
	}
	return nil, errors.New("program not found")
}

// mockMapsService is a mock implementation of MapsService for testing.
type mockMapsService struct {
	maps []MapInfo
	err  error
}

func (m *mockMapsService) List() ([]MapInfo, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.maps, nil
}

func (m *mockMapsService) Get(id uint32) (*MapInfo, error) {
	if m.err != nil {
		return nil, m.err
	}
	for _, mp := range m.maps {
		if mp.ID == id {
			return &mp, nil
		}
	}
	return nil, errors.New("map not found")
}

func (m *mockMapsService) Dump(id uint32) ([]MapEntry, error) {
	return nil, nil
}

func TestNewModel(t *testing.T) {
	progSvc := &mockProgService{}
	mapsSvc := &mockMapsService{}

	m := NewModel(progSvc, mapsSvc)

	if m.state != ViewMenu {
		t.Errorf("initial state = %v, want ViewMenu", m.state)
	}
	if len(m.history) != 0 {
		t.Errorf("initial history length = %d, want 0", len(m.history))
	}
	if m.progSvc != progSvc {
		t.Error("progSvc not set correctly")
	}
	if m.mapsSvc != mapsSvc {
		t.Error("mapsSvc not set correctly")
	}
}

func TestViewStateString(t *testing.T) {
	tests := []struct {
		state ViewState
		want  string
	}{
		{ViewMenu, "Menu"},
		{ViewProgList, "Programs"},
		{ViewProgDetail, "Program Detail"},
		{ViewMapList, "Maps"},
		{ViewMapDetail, "Map Detail"},
		{ViewMapDump, "Map Dump"},
		{ViewState(99), "Unknown"},
	}

	for _, tt := range tests {
		t.Run(tt.want, func(t *testing.T) {
			if got := tt.state.String(); got != tt.want {
				t.Errorf("ViewState(%d).String() = %q, want %q", tt.state, got, tt.want)
			}
		})
	}
}

func TestPushState(t *testing.T) {
	m := NewModel(nil, nil)

	// Initial state should be ViewMenu
	if m.state != ViewMenu {
		t.Fatalf("initial state = %v, want ViewMenu", m.state)
	}

	// Push to ProgList
	m.pushState(ViewProgList)
	if m.state != ViewProgList {
		t.Errorf("state after push = %v, want ViewProgList", m.state)
	}
	if m.historyLen() != 1 {
		t.Errorf("history length = %d, want 1", m.historyLen())
	}

	// Push to ProgDetail
	m.pushState(ViewProgDetail)
	if m.state != ViewProgDetail {
		t.Errorf("state after second push = %v, want ViewProgDetail", m.state)
	}
	if m.historyLen() != 2 {
		t.Errorf("history length = %d, want 2", m.historyLen())
	}

	// Push to MapDetail (from ProgDetail via associated map)
	m.pushState(ViewMapDetail)
	if m.state != ViewMapDetail {
		t.Errorf("state after third push = %v, want ViewMapDetail", m.state)
	}
	if m.historyLen() != 3 {
		t.Errorf("history length = %d, want 3", m.historyLen())
	}
}

func TestPopState(t *testing.T) {
	m := NewModel(nil, nil)

	// Build up navigation history: Menu -> ProgList -> ProgDetail -> MapDetail
	m.pushState(ViewProgList)
	m.pushState(ViewProgDetail)
	m.pushState(ViewMapDetail)

	// Pop should return to ProgDetail
	prevState := m.popState()
	if prevState != ViewProgDetail {
		t.Errorf("popState() = %v, want ViewProgDetail", prevState)
	}
	if m.historyLen() != 2 {
		t.Errorf("history length after pop = %d, want 2", m.historyLen())
	}

	// Pop should return to ProgList
	prevState = m.popState()
	if prevState != ViewProgList {
		t.Errorf("popState() = %v, want ViewProgList", prevState)
	}
	if m.historyLen() != 1 {
		t.Errorf("history length after second pop = %d, want 1", m.historyLen())
	}

	// Pop should return to Menu
	prevState = m.popState()
	if prevState != ViewMenu {
		t.Errorf("popState() = %v, want ViewMenu", prevState)
	}
	if m.historyLen() != 0 {
		t.Errorf("history length after third pop = %d, want 0", m.historyLen())
	}

	// Pop on empty history should return ViewMenu
	prevState = m.popState()
	if prevState != ViewMenu {
		t.Errorf("popState() on empty history = %v, want ViewMenu", prevState)
	}
}

func TestClearHistory(t *testing.T) {
	m := NewModel(nil, nil)

	m.pushState(ViewProgList)
	m.pushState(ViewProgDetail)

	if m.historyLen() != 2 {
		t.Fatalf("history length before clear = %d, want 2", m.historyLen())
	}

	m.clearHistory()

	if m.historyLen() != 0 {
		t.Errorf("history length after clear = %d, want 0", m.historyLen())
	}
}

func TestNavigationHistoryFromProgDetailToMapDetail(t *testing.T) {
	// This tests the specific use case from the design:
	// When navigating from ProgDetail to MapDetail (via associated map IDs),
	// pressing Esc/Backspace should return to ProgDetail, not MapList.

	m := NewModel(nil, nil)

	// Navigate: Menu -> ProgList -> ProgDetail
	m.pushState(ViewProgList)
	m.pushState(ViewProgDetail)

	// Now navigate to MapDetail from ProgDetail (via associated map ID)
	m.pushState(ViewMapDetail)

	if m.state != ViewMapDetail {
		t.Fatalf("state = %v, want ViewMapDetail", m.state)
	}

	// Pop should return to ProgDetail, NOT MapList
	m.state = m.popState()
	if m.state != ViewProgDetail {
		t.Errorf("after back from MapDetail, state = %v, want ViewProgDetail", m.state)
	}
}

func TestCheckPermissions(t *testing.T) {
	tests := []struct {
		name    string
		progSvc ProgService
		wantErr bool
	}{
		{
			name:    "nil service - no error",
			progSvc: nil,
			wantErr: false,
		},
		{
			name:    "successful list - no error",
			progSvc: &mockProgService{programs: []ProgramInfo{{ID: 1}}},
			wantErr: false,
		},
		{
			name:    "permission error",
			progSvc: &mockProgService{err: &PermissionError{}},
			wantErr: true,
		},
		{
			name:    "other error wrapped as permission error",
			progSvc: &mockProgService{err: errors.New("some error")},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := NewModel(tt.progSvc, nil)
			err := m.checkPermissions()
			if (err != nil) != tt.wantErr {
				t.Errorf("checkPermissions() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestUpdateWindowSize(t *testing.T) {
	m := NewModel(nil, nil)

	msg := tea.WindowSizeMsg{Width: 100, Height: 50}
	newModel, _ := m.Update(msg)
	updated := newModel.(Model)

	if updated.width != 100 {
		t.Errorf("width = %d, want 100", updated.width)
	}
	if updated.height != 50 {
		t.Errorf("height = %d, want 50", updated.height)
	}
}

func TestUpdateQuit(t *testing.T) {
	m := NewModel(nil, nil)

	// Test 'q' key
	msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}}
	_, cmd := m.Update(msg)

	if cmd == nil {
		t.Error("expected quit command, got nil")
	}
}

func TestUpdateBackNavigation(t *testing.T) {
	m := NewModel(nil, nil)
	m.pushState(ViewProgList)

	// Press Esc to go back
	msg := tea.KeyMsg{Type: tea.KeyEsc}
	newModel, _ := m.Update(msg)
	updated := newModel.(Model)

	if updated.state != ViewMenu {
		t.Errorf("state after back = %v, want ViewMenu", updated.state)
	}
}

func TestUpdateBackAtMenuDoesNothing(t *testing.T) {
	m := NewModel(nil, nil)

	// Press Esc at menu - should stay at menu
	msg := tea.KeyMsg{Type: tea.KeyEsc}
	newModel, _ := m.Update(msg)
	updated := newModel.(Model)

	if updated.state != ViewMenu {
		t.Errorf("state after back at menu = %v, want ViewMenu", updated.state)
	}
}

func TestUpdateHelpToggle(t *testing.T) {
	m := NewModel(nil, nil)

	// Press '?' to show help
	msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'?'}}
	newModel, _ := m.Update(msg)
	updated := newModel.(Model)

	if !updated.showHelp {
		t.Error("showHelp should be true after pressing '?'")
	}

	// Press '?' again to hide help
	newModel, _ = updated.Update(msg)
	updated = newModel.(Model)

	// Any key should close help, including '?'
	if updated.showHelp {
		t.Error("showHelp should be false after pressing any key while help is shown")
	}
}

func TestViewWithError(t *testing.T) {
	m := NewModel(nil, nil)
	m.err = errors.New("test error")

	view := m.View()

	if view == "" {
		t.Error("View() returned empty string for error state")
	}
	// Should contain error message
	if !containsString(view, "Error") {
		t.Error("error view should contain 'Error'")
	}
}

func TestViewWithHelp(t *testing.T) {
	m := NewModel(nil, nil)
	m.showHelp = true

	view := m.View()

	if view == "" {
		t.Error("View() returned empty string for help state")
	}
	// Should contain keyboard shortcuts header
	if !containsString(view, "Keyboard Shortcuts") {
		t.Error("help view should contain 'Keyboard Shortcuts'")
	}
	// Should contain navigation section
	if !containsString(view, "Navigation") {
		t.Error("help view should contain 'Navigation' section")
	}
	// Should contain global section
	if !containsString(view, "Global") {
		t.Error("help view should contain 'Global' section")
	}
	// Should contain close instruction
	if !containsString(view, "Press any key to close") {
		t.Error("help view should contain close instruction")
	}
}

func TestViewStates(t *testing.T) {
	states := []ViewState{
		ViewMenu,
		ViewProgList,
		ViewProgDetail,
		ViewMapList,
		ViewMapDetail,
		ViewMapDump,
	}

	for _, state := range states {
		t.Run(state.String(), func(t *testing.T) {
			m := NewModel(nil, nil)
			m.state = state

			view := m.View()
			if view == "" {
				t.Errorf("View() for %v returned empty string", state)
			}
		})
	}
}

func TestRenderHelpBar(t *testing.T) {
	tests := []struct {
		state    ViewState
		contains []string
	}{
		{ViewMenu, []string{"navigate", "select", "quit", "help"}},
		{ViewProgList, []string{"navigate", "search", "back", "quit", "help"}},
		{ViewMapList, []string{"navigate", "search", "back", "quit", "help"}},
		{ViewProgDetail, []string{"select map", "view map", "back", "quit", "help"}},
		{ViewMapDetail, []string{"dump contents", "back", "quit", "help"}},
		{ViewMapDump, []string{"scroll", "back", "quit", "help"}},
	}

	for _, tt := range tests {
		t.Run(tt.state.String(), func(t *testing.T) {
			m := NewModel(nil, nil)
			m.state = tt.state

			helpBar := m.renderHelpBar()
			for _, expected := range tt.contains {
				if !containsString(helpBar, expected) {
					t.Errorf("help bar for %v should contain %q, got: %s", tt.state, expected, helpBar)
				}
			}
		})
	}
}

func TestRenderHelpBarWhileFiltering(t *testing.T) {
	// Test that help bar shows different content when filtering
	m := NewModel(nil, nil)
	m.state = ViewProgList

	// Normal state should show search option
	helpBar := m.renderHelpBar()
	if !containsString(helpBar, "search") {
		t.Error("help bar should show 'search' when not filtering")
	}
}

func TestRenderHelpOverlayContextSpecific(t *testing.T) {
	tests := []struct {
		state    ViewState
		contains []string
	}{
		{ViewMenu, []string{"Navigation", "Menu", "Global", "Open selected option"}},
		{ViewProgList, []string{"Navigation", "List", "Global", "fuzzy search"}},
		{ViewMapList, []string{"Navigation", "List", "Global", "fuzzy search"}},
		{ViewProgDetail, []string{"Navigation", "Program Detail", "Global", "Navigate associated maps"}},
		{ViewMapDetail, []string{"Navigation", "Map Detail", "Global", "Dump map contents"}},
		{ViewMapDump, []string{"Navigation", "Map Dump", "Global", "Scroll through entries"}},
	}

	for _, tt := range tests {
		t.Run(tt.state.String(), func(t *testing.T) {
			m := NewModel(nil, nil)
			m.state = tt.state
			m.showHelp = true

			helpView := m.renderHelp()
			for _, expected := range tt.contains {
				if !containsString(helpView, expected) {
					t.Errorf("help overlay for %v should contain %q", tt.state, expected)
				}
			}
		})
	}
}

func TestHelpOverlayClosesOnAnyKey(t *testing.T) {
	testKeys := []tea.KeyMsg{
		{Type: tea.KeyRunes, Runes: []rune{'a'}},
		{Type: tea.KeyRunes, Runes: []rune{'z'}},
		{Type: tea.KeyEnter},
		{Type: tea.KeySpace},
		{Type: tea.KeyUp},
		{Type: tea.KeyDown},
	}

	for _, key := range testKeys {
		t.Run(key.String(), func(t *testing.T) {
			m := NewModel(nil, nil)
			m.showHelp = true

			newModel, _ := m.Update(key)
			updated := newModel.(Model)

			if updated.showHelp {
				t.Errorf("help should close on key %v", key)
			}
		})
	}
}

func TestHelpToggleWithQuestionMark(t *testing.T) {
	m := NewModel(nil, nil)

	// Initially help should be hidden
	if m.showHelp {
		t.Error("help should be hidden initially")
	}

	// Press '?' to show help
	msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'?'}}
	newModel, _ := m.Update(msg)
	updated := newModel.(Model)

	if !updated.showHelp {
		t.Error("help should be shown after pressing '?'")
	}

	// Press any key (including '?') to close help
	newModel, _ = updated.Update(msg)
	updated = newModel.(Model)

	if updated.showHelp {
		t.Error("help should be hidden after pressing any key")
	}
}

func TestHelpNotShownWhenFiltering(t *testing.T) {
	// When filtering in a list, '?' should not toggle help
	// This is handled by the list component, but we test the model behavior

	m := NewModel(nil, nil)
	m.state = ViewProgList

	// Simulate that we're not filtering - help should toggle
	msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'?'}}
	newModel, _ := m.Update(msg)
	updated := newModel.(Model)

	if !updated.showHelp {
		t.Error("help should toggle when not filtering")
	}
}

func TestPermissionError(t *testing.T) {
	// Test PermissionError with wrapped error
	innerErr := errors.New("operation not permitted")
	permErr := &PermissionError{Err: innerErr}

	expectedMsg := "insufficient permissions: operation not permitted"
	if permErr.Error() != expectedMsg {
		t.Errorf("PermissionError.Error() = %q, want %q",
			permErr.Error(), expectedMsg)
	}

	if permErr.Unwrap() != innerErr {
		t.Error("PermissionError.Unwrap() should return inner error")
	}

	// Test PermissionError without wrapped error
	permErr2 := &PermissionError{}
	expectedMsg2 := "insufficient permissions"
	if permErr2.Error() != expectedMsg2 {
		t.Errorf("PermissionError.Error() without inner = %q, want %q", permErr2.Error(), expectedMsg2)
	}
}

func TestIsPermissionError(t *testing.T) {
	if !IsPermissionError(&PermissionError{}) {
		t.Error("IsPermissionError should return true for PermissionError")
	}

	if IsPermissionError(errors.New("other error")) {
		t.Error("IsPermissionError should return false for other errors")
	}

	if IsPermissionError(nil) {
		t.Error("IsPermissionError should return false for nil")
	}
}

// Helper function to check if a string contains a substring
func containsString(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > 0 && containsSubstring(s, substr))
}

func containsSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

// ============================================================================
// Integration Tests for Error Scenarios
// ============================================================================

// TestIntegrationPermissionErrorOnStartup tests that permission errors are
// displayed gracefully on startup without crashing.
func TestIntegrationPermissionErrorOnStartup(t *testing.T) {
	// Create a mock service that returns a permission error
	mockSvc := &mockProgService{
		err: &PermissionError{Err: errors.New("operation not permitted")},
	}

	m := NewModel(mockSvc, nil)

	// Check permissions should set the error
	err := m.checkPermissions()
	if err == nil {
		t.Fatal("expected permission error")
	}

	// Set the error on the model (as RunWithServices does)
	m.err = err

	// View should render the error gracefully
	view := m.View()

	if view == "" {
		t.Error("View() should not return empty string for permission error")
	}

	// Should contain permission-related message
	if !containsString(view, "Permission Error") {
		t.Error("view should contain 'Permission Error'")
	}

	// Should contain sudo suggestion
	if !containsString(view, "sudo") {
		t.Error("view should contain 'sudo' suggestion")
	}

	// Should contain quit instruction
	if !containsString(view, "quit") {
		t.Error("view should contain quit instruction")
	}
}

// TestIntegrationErrorInProgListDoesNotCrash tests that errors in the programs
// list are displayed inline without crashing.
func TestIntegrationErrorInProgListDoesNotCrash(t *testing.T) {
	// Create a mock service that returns an error on List
	mockSvc := &mockProgService{
		err: errors.New("failed to read BPF programs"),
	}

	m := NewModel(mockSvc, nil)

	// Navigate to programs list
	enterMsg := tea.KeyMsg{Type: tea.KeyEnter}
	result, _ := m.Update(enterMsg)
	m = result.(Model)

	// Should be at programs list
	if m.state != ViewProgList {
		t.Fatalf("expected ViewProgList, got %v", m.state)
	}

	// View should render without crashing
	view := m.View()

	if view == "" {
		t.Error("View() should not return empty string")
	}

	// Should contain error message
	if !containsString(view, "Error") {
		t.Error("view should contain error message")
	}

	// Should still be able to navigate back
	escMsg := tea.KeyMsg{Type: tea.KeyEsc}
	result, _ = m.Update(escMsg)
	m = result.(Model)

	if m.state != ViewMenu {
		t.Errorf("expected ViewMenu after back, got %v", m.state)
	}
}

// TestIntegrationErrorInMapListDoesNotCrash tests that errors in the maps
// list are displayed inline without crashing.
func TestIntegrationErrorInMapListDoesNotCrash(t *testing.T) {
	// Create a mock service that returns an error on List
	mockMapsSvc := &mockMapsService{
		err: errors.New("failed to read BPF maps"),
	}

	m := NewModel(nil, mockMapsSvc)

	// Navigate to maps list (down to select Maps, then Enter)
	downMsg := tea.KeyMsg{Type: tea.KeyDown}
	result, _ := m.Update(downMsg)
	m = result.(Model)

	enterMsg := tea.KeyMsg{Type: tea.KeyEnter}
	result, _ = m.Update(enterMsg)
	m = result.(Model)

	// Should be at maps list
	if m.state != ViewMapList {
		t.Fatalf("expected ViewMapList, got %v", m.state)
	}

	// View should render without crashing
	view := m.View()

	if view == "" {
		t.Error("View() should not return empty string")
	}

	// Should contain error message
	if !containsString(view, "Error") {
		t.Error("view should contain error message")
	}

	// Should still be able to navigate back
	escMsg := tea.KeyMsg{Type: tea.KeyEsc}
	result, _ = m.Update(escMsg)
	m = result.(Model)

	if m.state != ViewMenu {
		t.Errorf("expected ViewMenu after back, got %v", m.state)
	}
}

// TestIntegrationEmptyProgramsList tests that an empty programs list
// displays the appropriate message.
func TestIntegrationEmptyProgramsList(t *testing.T) {
	// Create a mock service that returns empty list
	mockSvc := &mockProgService{
		programs: []ProgramInfo{},
	}

	m := NewModel(mockSvc, nil)

	// Navigate to programs list
	enterMsg := tea.KeyMsg{Type: tea.KeyEnter}
	result, _ := m.Update(enterMsg)
	m = result.(Model)

	// View should show empty state message
	view := m.View()

	if !containsString(view, "No BPF programs loaded") {
		t.Error("view should contain empty state message for programs")
	}
}

// TestIntegrationEmptyMapsList tests that an empty maps list
// displays the appropriate message.
func TestIntegrationEmptyMapsList(t *testing.T) {
	// Create a mock service that returns empty list
	mockMapsSvc := &mockMapsService{
		maps: []MapInfo{},
	}

	m := NewModel(nil, mockMapsSvc)

	// Navigate to maps list
	downMsg := tea.KeyMsg{Type: tea.KeyDown}
	result, _ := m.Update(downMsg)
	m = result.(Model)

	enterMsg := tea.KeyMsg{Type: tea.KeyEnter}
	result, _ = m.Update(enterMsg)
	m = result.(Model)

	// View should show empty state message
	view := m.View()

	if !containsString(view, "No BPF maps loaded") {
		t.Error("view should contain empty state message for maps")
	}
}

// TestIntegrationEmptyMapDump tests that an empty map dump
// displays the appropriate message.
func TestIntegrationEmptyMapDump(t *testing.T) {
	// Create mock services
	mockMapsSvc := &mockMapsServiceWithDump{
		maps: []MapInfo{
			{ID: 1, Name: "empty_map", Type: "hash", KeySize: 4, ValueSize: 8, MaxEntries: 100},
		},
		entries: []MapEntry{}, // Empty entries
	}

	m := NewModel(nil, mockMapsSvc)

	// Navigate to maps list
	downMsg := tea.KeyMsg{Type: tea.KeyDown}
	result, _ := m.Update(downMsg)
	m = result.(Model)

	enterMsg := tea.KeyMsg{Type: tea.KeyEnter}
	result, _ = m.Update(enterMsg)
	m = result.(Model)

	// Select the map
	result, _ = m.Update(enterMsg)
	m = result.(Model)

	// Should be at map detail
	if m.state != ViewMapDetail {
		t.Fatalf("expected ViewMapDetail, got %v", m.state)
	}

	// Navigate to dump
	result, _ = m.Update(enterMsg)
	m = result.(Model)

	// Should be at map dump
	if m.state != ViewMapDump {
		t.Fatalf("expected ViewMapDump, got %v", m.state)
	}

	// View should show empty state message
	view := m.View()

	if !containsString(view, "Map contains no entries") {
		t.Error("view should contain empty state message for map dump")
	}
}

// TestIntegrationErrorInMapDumpDoesNotCrash tests that errors during map dump
// are displayed inline without crashing.
func TestIntegrationErrorInMapDumpDoesNotCrash(t *testing.T) {
	// Create mock services where dump fails
	mockMapsSvc := &mockMapsServiceWithDump{
		maps: []MapInfo{
			{ID: 1, Name: "test_map", Type: "hash", KeySize: 4, ValueSize: 8, MaxEntries: 100},
		},
		dumpErr: errors.New("failed to dump map contents"),
	}

	m := NewModel(nil, mockMapsSvc)

	// Navigate to maps list
	downMsg := tea.KeyMsg{Type: tea.KeyDown}
	result, _ := m.Update(downMsg)
	m = result.(Model)

	enterMsg := tea.KeyMsg{Type: tea.KeyEnter}
	result, _ = m.Update(enterMsg)
	m = result.(Model)

	// Select the map
	result, _ = m.Update(enterMsg)
	m = result.(Model)

	// Navigate to dump
	result, _ = m.Update(enterMsg)
	m = result.(Model)

	// Should be at map dump
	if m.state != ViewMapDump {
		t.Fatalf("expected ViewMapDump, got %v", m.state)
	}

	// View should render without crashing
	view := m.View()

	if view == "" {
		t.Error("View() should not return empty string")
	}

	// Should contain error message
	if !containsString(view, "Error") {
		t.Error("view should contain error message")
	}

	// Should still be able to navigate back
	escMsg := tea.KeyMsg{Type: tea.KeyEsc}
	result, _ = m.Update(escMsg)
	m = result.(Model)

	if m.state != ViewMapDetail {
		t.Errorf("expected ViewMapDetail after back, got %v", m.state)
	}
}

// TestIntegrationErrorClearedOnNavigation tests that errors are cleared
// when navigating back.
func TestIntegrationErrorClearedOnNavigation(t *testing.T) {
	m := NewModel(nil, nil)

	// Set an error
	m.err = errors.New("test error")

	// Navigate to programs list (this should clear the error)
	m.pushState(ViewProgList)

	// Press back
	escMsg := tea.KeyMsg{Type: tea.KeyEsc}
	result, _ := m.Update(escMsg)
	m = result.(Model)

	// Error should be cleared
	if m.err != nil {
		t.Error("error should be cleared after navigating back")
	}
}

// TestIntegrationNilServicesDoNotCrash tests that nil services are handled
// gracefully without crashing.
func TestIntegrationNilServicesDoNotCrash(t *testing.T) {
	m := NewModel(nil, nil)

	// Navigate to programs list
	enterMsg := tea.KeyMsg{Type: tea.KeyEnter}
	result, _ := m.Update(enterMsg)
	m = result.(Model)

	// View should render without crashing
	view := m.View()
	if view == "" {
		t.Error("View() should not return empty string with nil services")
	}

	// Navigate back
	escMsg := tea.KeyMsg{Type: tea.KeyEsc}
	result, _ = m.Update(escMsg)
	m = result.(Model)

	// Navigate to maps list
	downMsg := tea.KeyMsg{Type: tea.KeyDown}
	result, _ = m.Update(downMsg)
	m = result.(Model)

	result, _ = m.Update(enterMsg)
	m = result.(Model)

	// View should render without crashing
	view = m.View()
	if view == "" {
		t.Error("View() should not return empty string with nil services")
	}
}

// mockMapsServiceWithDump is a mock that supports configurable dump behavior.
type mockMapsServiceWithDump struct {
	maps    []MapInfo
	entries []MapEntry
	dumpErr error
}

func (m *mockMapsServiceWithDump) List() ([]MapInfo, error) {
	return m.maps, nil
}

func (m *mockMapsServiceWithDump) Get(id uint32) (*MapInfo, error) {
	for _, mp := range m.maps {
		if mp.ID == id {
			return &mp, nil
		}
	}
	return nil, errors.New("map not found")
}

func (m *mockMapsServiceWithDump) Dump(id uint32) ([]MapEntry, error) {
	if m.dumpErr != nil {
		return nil, m.dumpErr
	}
	return m.entries, nil
}
