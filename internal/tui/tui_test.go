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
	// Should contain keyboard shortcuts
	if !containsString(view, "Keyboard Shortcuts") {
		t.Error("help view should contain 'Keyboard Shortcuts'")
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
		contains string
	}{
		{ViewMenu, "navigate"},
		{ViewProgList, "search"},
		{ViewMapList, "search"},
		{ViewProgDetail, "back"},
		{ViewMapDetail, "back"},
		{ViewMapDump, "back"},
	}

	for _, tt := range tests {
		t.Run(tt.state.String(), func(t *testing.T) {
			m := NewModel(nil, nil)
			m.state = tt.state

			helpBar := m.renderHelpBar()
			if !containsString(helpBar, tt.contains) {
				t.Errorf("help bar for %v should contain %q", tt.state, tt.contains)
			}
		})
	}
}

func TestPermissionError(t *testing.T) {
	// Test PermissionError with wrapped error
	innerErr := errors.New("operation not permitted")
	permErr := &PermissionError{Err: innerErr}

	if permErr.Error() != "insufficient permissions: operation not permitted" {
		t.Errorf("PermissionError.Error() = %q, want %q",
			permErr.Error(), "insufficient permissions: operation not permitted")
	}

	if permErr.Unwrap() != innerErr {
		t.Error("PermissionError.Unwrap() should return inner error")
	}

	// Test PermissionError without wrapped error
	permErr2 := &PermissionError{}
	if permErr2.Error() != "insufficient permissions - try running with sudo" {
		t.Errorf("PermissionError.Error() without inner = %q", permErr2.Error())
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
