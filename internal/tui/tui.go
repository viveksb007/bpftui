package tui

import (
	"fmt"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

// ViewState represents the current view in the TUI.
type ViewState int

const (
	ViewMenu ViewState = iota
	ViewProgList
	ViewProgDetail
	ViewMapList
	ViewMapDetail
	ViewMapDump
)

// String returns a human-readable name for the view state.
func (v ViewState) String() string {
	switch v {
	case ViewMenu:
		return "Menu"
	case ViewProgList:
		return "Programs"
	case ViewProgDetail:
		return "Program Detail"
	case ViewMapList:
		return "Maps"
	case ViewMapDetail:
		return "Map Detail"
	case ViewMapDump:
		return "Map Dump"
	default:
		return "Unknown"
	}
}

// Model is the root model for the TUI application.
type Model struct {
	// Current view state
	state ViewState

	// Navigation history for back navigation
	// When navigating from ProgDetail -> MapDetail, we push ProgDetail to history
	// On Esc/Backspace, we pop from history to return to the correct view
	history []ViewState

	// Services for data access
	progSvc ProgService
	mapsSvc MapsService

	// Sub-models for each view
	menu     menuModel
	progList progListModel

	// Terminal dimensions
	width  int
	height int

	// Error state
	err error

	// Key bindings
	keys keyMap

	// Help visibility
	showHelp bool
}

// NewModel creates a new root model with the given services.
func NewModel(progSvc ProgService, mapsSvc MapsService) Model {
	return Model{
		state:    ViewMenu,
		history:  make([]ViewState, 0),
		progSvc:  progSvc,
		mapsSvc:  mapsSvc,
		menu:     newMenuModel(80, 24),     // Default size, will be updated on WindowSizeMsg
		progList: newProgListModel(80, 24), // Default size, will be updated on WindowSizeMsg
		keys:     defaultKeyMap,
	}
}

// pushState saves the current state to history and transitions to a new state.
func (m *Model) pushState(newState ViewState) {
	m.history = append(m.history, m.state)
	m.state = newState
}

// popState returns to the previous state from history.
// If history is empty, returns to ViewMenu.
func (m *Model) popState() ViewState {
	if len(m.history) == 0 {
		return ViewMenu
	}
	lastIdx := len(m.history) - 1
	prevState := m.history[lastIdx]
	m.history = m.history[:lastIdx]
	return prevState
}

// historyLen returns the current navigation history length.
func (m *Model) historyLen() int {
	return len(m.history)
}

// clearHistory clears the navigation history.
func (m *Model) clearHistory() {
	m.history = m.history[:0]
}

// checkPermissions verifies the user has sufficient permissions for BPF operations.
func (m *Model) checkPermissions() error {
	if m.progSvc == nil {
		return nil // No service configured, skip check
	}
	_, err := m.progSvc.List()
	if err != nil {
		if IsPermissionError(err) {
			return err
		}
		// Wrap other errors as potential permission issues
		return &PermissionError{Err: err}
	}
	return nil
}

// Init implements tea.Model.
func (m Model) Init() tea.Cmd {
	// Check permissions on startup
	if err := m.checkPermissions(); err != nil {
		m.err = err
	}
	return nil
}

// Update implements tea.Model.
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		return m.handleKeyMsg(msg)

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.menu.SetSize(msg.Width, msg.Height)
		m.progList.SetSize(msg.Width, msg.Height)
		return m, nil
	}

	return m, nil
}

// handleKeyMsg processes keyboard input.
func (m Model) handleKeyMsg(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	// Handle quit from any view (but not when filtering)
	if key.Matches(msg, m.keys.Quit) {
		// Don't quit if we're in the middle of filtering
		if m.state == ViewProgList && m.progList.IsFiltering() {
			// Let the list handle it
		} else {
			return m, tea.Quit
		}
	}

	// Handle help toggle (but not when filtering)
	if key.Matches(msg, m.keys.Help) {
		if m.state == ViewProgList && m.progList.IsFiltering() {
			// Let the list handle it
		} else {
			m.showHelp = !m.showHelp
			return m, nil
		}
	}

	// If help is showing, any key closes it
	if m.showHelp {
		m.showHelp = false
		return m, nil
	}

	// Handle back navigation (but not when filtering)
	if key.Matches(msg, m.keys.Back) {
		// Don't navigate back if we're filtering - let the list handle escape
		if m.state == ViewProgList && m.progList.IsFiltering() {
			return m.handleProgListKeys(msg)
		}

		if m.state != ViewMenu {
			m.state = m.popState()
			m.err = nil // Clear any errors when navigating back
		}
		return m, nil
	}

	// View-specific key handling
	switch m.state {
	case ViewMenu:
		return m.handleMenuKeys(msg)
	case ViewProgList:
		return m.handleProgListKeys(msg)
	}

	return m, nil
}

// handleMenuKeys handles keyboard input in the menu view.
func (m Model) handleMenuKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var targetView *ViewState
	m.menu, cmd, targetView = m.menu.Update(msg)

	// If a menu item was selected, navigate to that view
	if targetView != nil {
		m.pushState(*targetView)

		// Load data for the target view
		if *targetView == ViewProgList {
			m.loadPrograms()
		}
	}

	return m, cmd
}

// handleProgListKeys handles keyboard input in the programs list view.
func (m Model) handleProgListKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var selectedProg *ProgramInfo
	m.progList, cmd, selectedProg = m.progList.Update(msg)

	// If a program was selected, navigate to detail view
	if selectedProg != nil {
		m.pushState(ViewProgDetail)
		// TODO: Set selected program in progDetail model when implemented
	}

	return m, cmd
}

// loadPrograms fetches programs from the service and updates the list.
func (m *Model) loadPrograms() {
	if m.progSvc == nil {
		m.progList.SetPrograms([]ProgramInfo{})
		return
	}

	programs, err := m.progSvc.List()
	if err != nil {
		m.progList.SetError(err)
		return
	}

	m.progList.SetPrograms(programs)
}

// View implements tea.Model.
func (m Model) View() string {
	// Show error if present
	if m.err != nil {
		return m.renderError()
	}

	// Show help overlay if active
	if m.showHelp {
		return m.renderHelp()
	}

	// Render current view
	switch m.state {
	case ViewMenu:
		return m.renderMenu()
	case ViewProgList:
		return m.renderProgList()
	case ViewProgDetail:
		return m.renderProgDetail()
	case ViewMapList:
		return m.renderMapList()
	case ViewMapDetail:
		return m.renderMapDetail()
	case ViewMapDump:
		return m.renderMapDump()
	default:
		return "Unknown view"
	}
}

// renderError displays an error message.
func (m Model) renderError() string {
	return errorStyle.Render(fmt.Sprintf("Error: %v\n\nPress 'q' to quit.", m.err))
}

// renderHelp displays the help overlay.
func (m Model) renderHelp() string {
	help := titleStyle.Render("Keyboard Shortcuts") + "\n\n"
	help += helpStyle.Render("↑/k      - Move up\n")
	help += helpStyle.Render("↓/j      - Move down\n")
	help += helpStyle.Render("Enter    - Select\n")
	help += helpStyle.Render("Esc      - Go back\n")
	help += helpStyle.Render("/        - Search\n")
	help += helpStyle.Render("?        - Toggle help\n")
	help += helpStyle.Render("q        - Quit\n")
	help += "\n" + dimStyle.Render("Press any key to close")
	return help
}

// renderMenu displays the main menu.
func (m Model) renderMenu() string {
	return m.menu.View() + "\n" + m.renderHelpBar()
}

// renderProgList displays the programs list.
func (m Model) renderProgList() string {
	return m.progList.View() + "\n" + m.renderHelpBar()
}

// renderProgDetail displays program details.
// Placeholder - will be implemented in task 5.
func (m Model) renderProgDetail() string {
	return titleStyle.Render("Program Details") + "\n\n" +
		"Program details will be displayed here.\n\n" +
		m.renderHelpBar()
}

// renderMapList displays the maps list.
// Placeholder - will be implemented in task 6.
func (m Model) renderMapList() string {
	return titleStyle.Render("BPF Maps") + "\n\n" +
		"Maps list will be displayed here.\n\n" +
		m.renderHelpBar()
}

// renderMapDetail displays map details.
// Placeholder - will be implemented in task 7.
func (m Model) renderMapDetail() string {
	return titleStyle.Render("Map Details") + "\n\n" +
		"Map details will be displayed here.\n\n" +
		m.renderHelpBar()
}

// renderMapDump displays map contents.
// Placeholder - will be implemented in task 8.
func (m Model) renderMapDump() string {
	return titleStyle.Render("Map Dump") + "\n\n" +
		"Map contents will be displayed here.\n\n" +
		m.renderHelpBar()
}

// renderHelpBar displays context-appropriate shortcuts at the bottom.
func (m Model) renderHelpBar() string {
	var shortcuts string
	switch m.state {
	case ViewMenu:
		shortcuts = "↑/↓: navigate • enter: select • q: quit • ?: help"
	case ViewProgList, ViewMapList:
		shortcuts = "↑/↓: navigate • enter: select • /: search • esc: back • q: quit • ?: help"
	default:
		shortcuts = "↑/↓: navigate • enter: select • esc: back • q: quit • ?: help"
	}
	return helpStyle.Render(shortcuts)
}

// Run starts the TUI application.
// This is the main entry point called from main.go.
func Run() error {
	return RunWithServices(nil, nil)
}

// RunWithServices starts the TUI application with the provided services.
// If services are nil, the TUI will run in a limited mode.
func RunWithServices(progSvc ProgService, mapsSvc MapsService) error {
	m := NewModel(progSvc, mapsSvc)

	// Check permissions before starting
	if err := m.checkPermissions(); err != nil {
		// Set error but still start TUI to show error gracefully
		m.err = err
	}

	p := tea.NewProgram(m, tea.WithAltScreen())
	_, err := p.Run()
	return err
}
