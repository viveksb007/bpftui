# Implementation Plan

- [x] 1. Set up project structure and dependencies
  - Create new `bpftui` repository with `go mod init github.com/viveksb007/bpftui`
  - Add dependencies: bubbletea, bubbles, lipgloss, and github.com/viveksb007/gobpftool
  - Create `main.go` entry point that calls `tui.Run()`
  - Create `internal/tui/` directory structure
  - _Requirements: 1.1, 1.4_

- [x] 2. Implement core TUI framework
  - [x] 2.1 Create key bindings and styles
    - Implement `keys.go` with keyMap struct and defaultKeyMap bindings (up/down/enter/back/quit/search/help)
    - Implement `styles.go` with Lipgloss styles (titleStyle, selectedStyle, helpStyle, errorStyle)
    - Write unit tests for key binding definitions
    - _Requirements: 8.1, 8.4_

  - [x] 2.2 Create root model and TUI entry point
    - Implement `tui.go` with Model struct containing ViewState, navigation history, services, and sub-models
    - Implement `Run()` function that initializes services and starts Bubbletea program
    - Implement `Init()`, `Update()`, and `View()` methods for root model
    - Implement navigation helpers `pushState()` and `popState()` for history management
    - Add permission check on startup with graceful error handling
    - Write unit tests for state transitions and navigation history
    - _Requirements: 1.1, 1.3, 1.4_

- [x] 3. Implement main menu component
  - Implement `menu.go` with menuModel using Bubbles list.Model
  - Create menuItem struct with "Programs" and "Maps" options
  - Implement `Init()`, `Update()`, `View()` methods for menu
  - Handle Enter key to navigate to ProgList or MapList
  - Handle q/Ctrl+C to quit application
  - Write unit tests for menu navigation
  - _Requirements: 1.2, 1.4_

- [x] 4. Implement programs list component
  - [x] 4.1 Create programs list model
    - Implement `proglist.go` with progListModel using Bubbles list.Model
    - Create progItem struct implementing list.Item interface (Title, Description, FilterValue)
    - Display ID, Name, Type, Tag for each program
    - _Requirements: 2.1, 2.2_

  - [x] 4.2 Add navigation and scrolling
    - Implement keyboard navigation (arrow keys, j/k)
    - Enable scrolling when list exceeds terminal height
    - Handle Enter to navigate to ProgDetail
    - Handle Esc/Backspace to return to main menu
    - Write unit tests for list navigation
    - _Requirements: 2.3, 2.4, 2.5, 2.6_

- [ ] 5. Implement program detail component
  - Implement `progdetail.go` with progDetailModel using Bubbles viewport.Model
  - Display all ProgramInfo fields (ID, Type, Name, Tag, GPL, LoadedAt, UID, BytesXlated, BytesJIT, MemLock)
  - Display associated MapIDs as selectable list
  - Handle cursor navigation through MapIDs
  - Handle Enter on MapID to navigate to MapDetail (push current state to history)
  - Handle Esc/Backspace to return to programs list
  - Write unit tests for detail rendering and map navigation
  - _Requirements: 3.1, 3.2, 3.3, 3.4_

- [ ] 6. Implement maps list component
  - [ ] 6.1 Create maps list model
    - Implement `maplist.go` with mapListModel using Bubbles list.Model
    - Create mapItem struct implementing list.Item interface
    - Display ID, Name, Type, KeySize, ValueSize, MaxEntries for each map
    - _Requirements: 4.1, 4.2_

  - [ ] 6.2 Add navigation and scrolling
    - Implement keyboard navigation (arrow keys, j/k)
    - Enable scrolling when list exceeds terminal height
    - Handle Enter to navigate to MapDetail
    - Handle Esc/Backspace to return to main menu
    - Write unit tests for list navigation
    - _Requirements: 4.3, 4.4, 4.5, 4.6_

- [ ] 7. Implement map detail component
  - Implement `mapdetail.go` with mapDetailModel using Bubbles viewport.Model
  - Display all MapInfo fields (ID, Type, Name, KeySize, ValueSize, MaxEntries, Flags, MemLock, LoadedAt, UID)
  - Add "Dump Contents" action option
  - Handle Enter on Dump to navigate to MapDump
  - Handle Esc/Backspace to return to previous view (using navigation history)
  - Write unit tests for detail rendering
  - _Requirements: 5.1, 5.2, 5.3_

- [ ] 8. Implement map dump component
  - Implement `mapdump.go` with mapDumpModel using Bubbles viewport.Model
  - Fetch map entries using maps.Service
  - Display key-value pairs in hexadecimal format
  - Enable scrolling through entries
  - Handle empty map state with "Map contains no entries" message
  - Handle errors gracefully without crashing
  - Handle Esc/Backspace to return to map detail
  - Write unit tests for hex formatting and error handling
  - _Requirements: 6.1, 6.2, 6.3, 6.4, 6.5, 6.6_

- [ ] 9. Implement fuzzy search
  - [ ] 9.1 Add fuzzy search to programs list
    - Enable Bubbles list.Model built-in filtering with `/` key
    - Display search input field when active
    - Filter list in real-time using fuzzy matching
    - Allow navigation through filtered results while search is active
    - Handle Enter to navigate to selected item's detail view
    - Handle Esc to exit search mode and restore full list
    - Display "No matches found" when no results
    - _Requirements: 7.1, 7.2, 7.3, 7.5, 7.6, 7.7, 7.8_

  - [ ] 9.2 Add fuzzy search to maps list
    - Enable same fuzzy search functionality for maps list
    - Write unit tests for search filtering behavior
    - _Requirements: 7.1, 7.2, 7.3, 7.5, 7.6, 7.7, 7.8_

- [ ] 10. Implement help system
  - Add help bar at bottom of each view showing context-appropriate shortcuts
  - Implement help overlay triggered by `?` key
  - Close help overlay on any key press
  - Show different shortcuts based on current view context
  - Write unit tests for help rendering
  - _Requirements: 8.1, 8.2, 8.3, 8.4_

- [ ] 11. Add error handling and empty states
  - Implement permission error handling with clear message on startup
  - Add empty state messages for programs list, maps list, and map dump
  - Display inline errors using errorStyle without crashing
  - Write integration tests for error scenarios
  - _Requirements: 1.3, 6.4, 6.5_

- [ ] 12. Final integration and polish
  - Wire all components together in root model Update() switch
  - Test full navigation flow: Menu → List → Detail → Dump → Back
  - Test navigation from ProgDetail → MapDetail → Back to ProgDetail
  - Verify terminal resize handling
  - Add README.md with installation and usage instructions
  - _Requirements: 1.1, 1.2, 1.4, 2.5, 2.6, 3.4, 4.5, 4.6, 5.3, 6.6_
