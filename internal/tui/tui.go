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
	menu       menuModel
	progList   progListModel
	progDetail progDetailModel
	mapList    mapListModel
	mapDetail  mapDetailModel
	mapDump    mapDumpModel

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
		state:      ViewMenu,
		history:    make([]ViewState, 0),
		progSvc:    progSvc,
		mapsSvc:    mapsSvc,
		menu:       newMenuModel(80, 24),       // Default size, will be updated on WindowSizeMsg
		progList:   newProgListModel(80, 24),   // Default size, will be updated on WindowSizeMsg
		progDetail: newProgDetailModel(80, 24), // Default size, will be updated on WindowSizeMsg
		mapList:    newMapListModel(80, 24),    // Default size, will be updated on WindowSizeMsg
		mapDetail:  newMapDetailModel(80, 24),  // Default size, will be updated on WindowSizeMsg
		mapDump:    newMapDumpModel(80, 24),    // Default size, will be updated on WindowSizeMsg
		keys:       defaultKeyMap,
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
		m.progDetail.SetSize(msg.Width, msg.Height)
		m.mapList.SetSize(msg.Width, msg.Height)
		m.mapDetail.SetSize(msg.Width, msg.Height)
		m.mapDump.SetSize(msg.Width, msg.Height)
		return m, nil

	default:
		// Pass other messages (like FilterMatchesMsg) to the active sub-model
		// This is necessary for async operations like list filtering to work
		return m.handleOtherMsg(msg)
	}
}

// handleOtherMsg passes non-key messages to the appropriate sub-model.
// This is necessary for async operations like list filtering.
func (m Model) handleOtherMsg(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch m.state {
	case ViewProgList:
		m.progList, cmd, _ = m.progList.Update(msg)
	case ViewMapList:
		m.mapList, cmd, _ = m.mapList.Update(msg)
	case ViewMapDump:
		m.mapDump, cmd = m.mapDump.Update(msg)
	}

	return m, cmd
}

// handleKeyMsg processes keyboard input.
func (m Model) handleKeyMsg(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	// Handle quit from any view (but not when filtering)
	if key.Matches(msg, m.keys.Quit) {
		// Don't quit if we're in the middle of filtering
		if m.state == ViewProgList && m.progList.IsFiltering() {
			// Let the list handle it
		} else if m.state == ViewMapList && m.mapList.IsFiltering() {
			// Let the list handle it
		} else {
			return m, tea.Quit
		}
	}

	// Handle help toggle (but not when filtering)
	if key.Matches(msg, m.keys.Help) {
		if m.state == ViewProgList && m.progList.IsFiltering() {
			// Let the list handle it
		} else if m.state == ViewMapList && m.mapList.IsFiltering() {
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
		if m.state == ViewMapList && m.mapList.IsFiltering() {
			return m.handleMapListKeys(msg)
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
	case ViewProgDetail:
		return m.handleProgDetailKeys(msg)
	case ViewMapList:
		return m.handleMapListKeys(msg)
	case ViewMapDetail:
		return m.handleMapDetailKeys(msg)
	case ViewMapDump:
		return m.handleMapDumpKeys(msg)
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
		switch *targetView {
		case ViewProgList:
			m.loadPrograms()
		case ViewMapList:
			m.loadMaps()
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
		m.progDetail.SetProgram(selectedProg)
	}

	return m, cmd
}

// handleProgDetailKeys handles keyboard input in the program detail view.
func (m Model) handleProgDetailKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var selectedMapID *uint32
	m.progDetail, cmd, selectedMapID = m.progDetail.Update(msg)

	// If a map was selected, navigate to map detail view
	if selectedMapID != nil {
		m.pushState(ViewMapDetail)
		m.loadMapByID(*selectedMapID)
	}

	return m, cmd
}

// handleMapListKeys handles keyboard input in the maps list view.
func (m Model) handleMapListKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var selectedMap *MapInfo
	m.mapList, cmd, selectedMap = m.mapList.Update(msg)

	// If a map was selected, navigate to detail view
	if selectedMap != nil {
		m.pushState(ViewMapDetail)
		m.mapDetail.SetMap(selectedMap)
	}

	return m, cmd
}

// handleMapDetailKeys handles keyboard input in the map detail view.
func (m Model) handleMapDetailKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var dumpSelected bool
	m.mapDetail, cmd, dumpSelected = m.mapDetail.Update(msg)

	// If Dump was selected, navigate to map dump view
	if dumpSelected {
		m.pushState(ViewMapDump)
		m.loadMapDump(m.mapDetail.GetMapID())
	}

	return m, cmd
}

// handleMapDumpKeys handles keyboard input in the map dump view.
func (m Model) handleMapDumpKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	m.mapDump, cmd = m.mapDump.Update(msg)
	return m, cmd
}

// loadPrograms fetches programs from the service and updates the list.
func (m *Model) loadPrograms() {
	// Reset any existing filter when entering the list
	m.progList.ResetFilter()

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

// loadMaps fetches maps from the service and updates the list.
func (m *Model) loadMaps() {
	// Reset any existing filter when entering the list
	m.mapList.ResetFilter()

	if m.mapsSvc == nil {
		m.mapList.SetMaps([]MapInfo{})
		return
	}

	maps, err := m.mapsSvc.List()
	if err != nil {
		m.mapList.SetError(err)
		return
	}

	m.mapList.SetMaps(maps)
}

// loadMapByID fetches a specific map by ID and sets it in the map detail view.
func (m *Model) loadMapByID(id uint32) {
	if m.mapsSvc == nil {
		return
	}

	mapInfo, err := m.mapsSvc.Get(id)
	if err != nil {
		m.err = err
		return
	}

	m.mapDetail.SetMap(mapInfo)
}

// loadMapDump fetches map entries and sets them in the map dump view.
func (m *Model) loadMapDump(id uint32) {
	if m.mapsSvc == nil {
		m.mapDump.SetMapDump(id, "", []MapEntry{})
		return
	}

	// Get map name from current map detail if available
	mapName := ""
	if mapInfo := m.mapDetail.GetMapInfo(); mapInfo != nil && mapInfo.ID == id {
		mapName = mapInfo.Name
	}

	entries, err := m.mapsSvc.Dump(id)
	if err != nil {
		m.mapDump.SetError(err)
		return
	}

	m.mapDump.SetMapDump(id, mapName, entries)
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

// renderHelp displays the help overlay with context-appropriate shortcuts.
func (m Model) renderHelp() string {
	var content string

	content += "Navigation:\n"
	content += "  ↑/k      Move up\n"
	content += "  ↓/j      Move down\n"
	content += "  Enter    Select\n"

	// Context-specific shortcuts
	switch m.state {
	case ViewMenu:
		content += "\nMenu:\n"
		content += "  Enter    Open selected option\n"

	case ViewProgList, ViewMapList:
		content += "\nList:\n"
		content += "  /        Start fuzzy search\n"
		content += "  Esc      Exit search / Go back\n"
		content += "  Enter    View details\n"

	case ViewProgDetail:
		content += "\nProgram Detail:\n"
		content += "  ↑/↓      Navigate associated maps\n"
		content += "  Enter    View selected map\n"
		content += "  Esc      Go back to list\n"

	case ViewMapDetail:
		content += "\nMap Detail:\n"
		content += "  Enter    Dump map contents\n"
		content += "  Esc      Go back\n"

	case ViewMapDump:
		content += "\nMap Dump:\n"
		content += "  ↑/↓      Scroll through entries\n"
		content += "  Esc      Go back to map detail\n"
	}

	// Global shortcuts
	content += "\nGlobal:\n"
	content += "  ?        Toggle this help\n"
	content += "  q        Quit application\n"

	return titleStyle.Render("Keyboard Shortcuts") + "\n\n" +
		helpStyle.Render(content) + "\n" +
		dimStyle.Render("Press any key to close")
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
func (m Model) renderProgDetail() string {
	return m.progDetail.View() + "\n" + m.renderHelpBar()
}

// renderMapList displays the maps list.
func (m Model) renderMapList() string {
	return m.mapList.View() + "\n" + m.renderHelpBar()
}

// renderMapDetail displays map details.
func (m Model) renderMapDetail() string {
	return m.mapDetail.View() + "\n" + m.renderHelpBar()
}

// renderMapDump displays map contents.
func (m Model) renderMapDump() string {
	return m.mapDump.View() + "\n" + m.renderHelpBar()
}

// renderHelpBar displays context-appropriate shortcuts at the bottom.
func (m Model) renderHelpBar() string {
	var shortcuts string
	switch m.state {
	case ViewMenu:
		shortcuts = "↑/↓: navigate • enter: select • q: quit • ?: help"
	case ViewProgList, ViewMapList:
		if m.state == ViewProgList && m.progList.IsFiltering() {
			shortcuts = "↑/↓: navigate • enter: select • esc: cancel search"
		} else if m.state == ViewMapList && m.mapList.IsFiltering() {
			shortcuts = "↑/↓: navigate • enter: select • esc: cancel search"
		} else {
			shortcuts = "↑/↓: navigate • enter: select • /: search • esc: back • q: quit • ?: help"
		}
	case ViewProgDetail:
		shortcuts = "↑/↓: select map • enter: view map • esc: back • q: quit • ?: help"
	case ViewMapDetail:
		shortcuts = "enter: dump contents • esc: back • q: quit • ?: help"
	case ViewMapDump:
		shortcuts = "↑/↓: scroll • esc: back • q: quit • ?: help"
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
