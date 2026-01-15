# bpftui

An interactive Terminal User Interface (TUI) for exploring BPF programs and maps. Built with [Bubbletea](https://github.com/charmbracelet/bubbletea) and powered by [gobpftool](https://github.com/viveksb007/gobpftool).

[![Build](https://github.com/viveksb007/bpftui/actions/workflows/build.yml/badge.svg)](https://github.com/viveksb007/bpftui/actions/workflows/build.yml)
![License](https://img.shields.io/badge/License-MIT-blue.svg)

## Features

- Browse loaded BPF programs and maps
- Fuzzy search to quickly find what you're looking for
- Dump map contents in hex
- Jump from a program directly to its associated maps
- Vim-style keyboard navigation
- Press `?` for help

## Installation

### Homebrew (Linux)

```bash
brew tap viveksb007/tap
brew install bpftui
```

### From Source

```bash
# Clone the repository
git clone https://github.com/viveksb007/bpftui.git
cd bpftui

# Build
go build -o bpftui .
```

### Requirements

- Linux with BPF support
- Root privileges (or CAP_BPF capability) to access BPF information

## Usage

```bash
# Run with sudo (required for BPF access)
sudo ./bpftui
```

### Navigation

| Key | Action |
|-----|--------|
| `↑` / `k` | Move up |
| `↓` / `j` | Move down |
| `Enter` | Select / Confirm |
| `Esc` / `Backspace` | Go back |
| `/` | Start fuzzy search (in list views) |
| `?` | Toggle help overlay |
| `q` / `Ctrl+C` | Quit |

### Views

#### Main Menu
The starting point with two options:
- **Programs** - Browse loaded BPF programs
- **Maps** - Browse loaded BPF maps

#### Programs List
Displays all loaded BPF programs with:
- Program ID
- Name
- Type (kprobe, tracepoint, xdp, etc.)
- Tag

Use `/` to fuzzy search by program name.

#### Program Detail
Shows detailed information about a selected program:
- ID, Name, Type, Tag
- GPL license status
- Load time and UID
- Bytes translated and JIT compiled
- Memory lock size
- Associated map IDs (selectable - press Enter to view map details)

#### Maps List
Displays all loaded BPF maps with:
- Map ID
- Name
- Type (hash, array, etc.)
- Key size, Value size, Max entries

Use `/` to fuzzy search by map name.

#### Map Detail
Shows detailed information about a selected map:
- ID, Name, Type
- Key size, Value size, Max entries
- Flags, Memory lock
- Load time and UID
- **Dump Contents** action - view map entries

#### Map Dump
Displays map contents in hexadecimal format:
```
Key:   01 02 03 04
Value: 0a 0b 0c 0d 0e 0f 10 11
---
Key:   01 02 03 05
Value: 1a 1b 1c 1d 1e 1f 20 21
```

## Troubleshooting

### Permission Denied

If you see a permission error, make sure you're running with root privileges:

```bash
sudo ./bpftui
```

Alternatively, you can grant the `CAP_BPF` capability:

```bash
sudo setcap cap_bpf+ep ./bpftui
./bpftui
```

### No Programs/Maps Found

If the lists are empty, it means no BPF programs or maps are currently loaded on your system. You can load some test programs using tools like `bpftrace` or by running BPF-based applications.

## Development

### Running Tests

```bash
go test ./... -v
```

### Project Structure

```
bpftui/
├── main.go              # Application entry point
├── go.mod               # Module definition
├── go.sum
├── internal/
│   └── tui/
│       ├── tui.go       # Main TUI model and entry point
│       ├── keys.go      # Key bindings
│       ├── styles.go    # Lipgloss styles
│       ├── services.go  # Service interfaces and types
│       ├── adapter.go   # Adapters for gobpftool services
│       ├── menu.go      # Main menu component
│       ├── proglist.go  # Programs list component
│       ├── progdetail.go # Program detail component
│       ├── maplist.go   # Maps list component
│       ├── mapdetail.go # Map detail component
│       └── mapdump.go   # Map dump component
└── README.md
```

## Dependencies

- [bubbletea](https://github.com/charmbracelet/bubbletea) - TUI framework
- [bubbles](https://github.com/charmbracelet/bubbles) - TUI components
- [lipgloss](https://github.com/charmbracelet/lipgloss) - Style definitions
- [gobpftool](https://github.com/viveksb007/gobpftool) - BPF program/map access

## License

MIT License - see [LICENSE](LICENSE) for details.

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.
