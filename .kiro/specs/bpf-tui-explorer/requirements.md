# Requirements Document

## Introduction

This feature adds an interactive Terminal User Interface (TUI) to gobpftool that allows users to explore BPF programs and maps more intuitively. Instead of running individual CLI commands, users can navigate through lists of programs and maps, view details, and dump map contents all within a single interactive session. The TUI will use the Bubbletea framework for a modern, responsive terminal experience and include fuzzy search capabilities for quickly finding specific programs or maps.

## Requirements

### Requirement 1: TUI Application Entry Point

**User Story:** As a user, I want to launch an interactive TUI mode so that I can explore BPF programs and maps without running multiple CLI commands.

#### Acceptance Criteria

1. WHEN the user runs `gobpftool tui` THEN the system SHALL launch the interactive TUI application
2. WHEN the TUI launches THEN the system SHALL display a main menu with options to explore Programs or Maps
3. IF the user does not have sufficient permissions THEN the system SHALL display a clear error message and exit gracefully
4. WHEN the user presses `q` or `Ctrl+C` at any screen THEN the system SHALL exit the TUI cleanly

### Requirement 2: BPF Programs List View

**User Story:** As a user, I want to see a scrollable list of all loaded BPF programs so that I can browse and select programs to inspect.

#### Acceptance Criteria

1. WHEN the user selects "Programs" from the main menu THEN the system SHALL display a list of all loaded BPF programs
2. WHEN displaying the programs list THEN the system SHALL show ID, Name, Type, and Tag for each program
3. WHEN the programs list is displayed THEN the system SHALL support keyboard navigation using arrow keys or j/k
4. WHEN the list exceeds the terminal height THEN the system SHALL enable scrolling through the list
5. WHEN the user presses Enter on a selected program THEN the system SHALL navigate to the program detail view
6. WHEN the user presses Escape or Backspace THEN the system SHALL return to the main menu

### Requirement 3: BPF Program Detail View

**User Story:** As a user, I want to view detailed information about a specific BPF program so that I can understand its configuration and associated maps.

#### Acceptance Criteria

1. WHEN viewing a program's details THEN the system SHALL display all ProgramInfo fields (ID, Type, Name, Tag, GPL, LoadedAt, UID, BytesXlated, BytesJIT, MemLock, MapIDs)
2. WHEN the program has associated maps THEN the system SHALL display the list of associated map IDs
3. IF the program has associated maps THEN the system SHALL allow the user to navigate to a map's detail view by selecting its ID
4. WHEN the user presses Escape or Backspace THEN the system SHALL return to the programs list

### Requirement 4: BPF Maps List View

**User Story:** As a user, I want to see a scrollable list of all loaded BPF maps so that I can browse and select maps to inspect or dump.

#### Acceptance Criteria

1. WHEN the user selects "Maps" from the main menu THEN the system SHALL display a list of all loaded BPF maps
2. WHEN displaying the maps list THEN the system SHALL show ID, Name, Type, KeySize, ValueSize, and MaxEntries for each map
3. WHEN the maps list is displayed THEN the system SHALL support keyboard navigation using arrow keys or j/k
4. WHEN the list exceeds the terminal height THEN the system SHALL enable scrolling through the list
5. WHEN the user presses Enter on a selected map THEN the system SHALL navigate to the map detail view
6. WHEN the user presses Escape or Backspace THEN the system SHALL return to the main menu

### Requirement 5: BPF Map Detail View

**User Story:** As a user, I want to view detailed information about a specific BPF map so that I can understand its configuration.

#### Acceptance Criteria

1. WHEN viewing a map's details THEN the system SHALL display all MapInfo fields (ID, Type, Name, KeySize, ValueSize, MaxEntries, Flags, MemLock, LoadedAt, UID)
2. WHEN viewing a map's details THEN the system SHALL provide an option to dump the map contents
3. WHEN the user presses Escape or Backspace THEN the system SHALL return to the maps list

### Requirement 6: Map Dump View

**User Story:** As a user, I want to dump and view the contents of a BPF map so that I can inspect the key-value pairs stored in it.

#### Acceptance Criteria

1. WHEN the user selects "Dump" on a map THEN the system SHALL retrieve and display all key-value entries
2. WHEN displaying map entries THEN the system SHALL show keys and values in hexadecimal format
3. WHEN the dump results exceed the terminal height THEN the system SHALL enable scrolling through the entries
4. IF the map is empty THEN the system SHALL display a message indicating the map has no entries
5. IF an error occurs during dump THEN the system SHALL display the error message without crashing
6. WHEN the user presses Escape or Backspace THEN the system SHALL return to the map detail view

### Requirement 7: Fuzzy Search

**User Story:** As a user, I want to fuzzy search through programs and maps so that I can quickly find specific items without scrolling through long lists.

#### Acceptance Criteria

1. WHEN the user presses `/` in the programs or maps list view THEN the system SHALL activate fuzzy search mode
2. WHEN fuzzy search is active THEN the system SHALL display a search input field
3. WHEN the user types in the search field THEN the system SHALL filter the list in real-time using fuzzy matching
4. WHEN fuzzy search matches items THEN the system SHALL highlight the matching portions of text
5. WHEN the filtered list is displayed THEN the system SHALL allow the user to navigate through results using arrow keys or j/k while keeping the search input active
6. WHEN the user presses Enter in search mode THEN the system SHALL navigate to the detail view of the currently selected item
7. WHEN the user presses Escape in search mode THEN the system SHALL exit search mode and restore the full list
8. WHEN no items match the search query THEN the system SHALL display a "no results" message

### Requirement 8: Help and Keyboard Shortcuts

**User Story:** As a user, I want to see available keyboard shortcuts so that I can navigate the TUI efficiently.

#### Acceptance Criteria

1. WHEN the TUI is running THEN the system SHALL display a help bar at the bottom showing common shortcuts
2. WHEN the user presses `?` THEN the system SHALL display a help overlay with all available keyboard shortcuts
3. WHEN the help overlay is displayed THEN the system SHALL close it when the user presses any key
4. WHEN displaying shortcuts THEN the system SHALL show context-appropriate shortcuts for the current view
