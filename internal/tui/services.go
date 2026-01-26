package tui

import "errors"

// ProgramInfo represents information about a BPF program.
// This mirrors the structure from gobpftool's prog.ProgramInfo.
type ProgramInfo struct {
	ID          uint32
	Type        string
	Name        string
	Tag         string
	GPL         bool
	LoadedAt    string
	UID         uint32
	BytesXlated uint32
	BytesJIT    uint32
	MemLock     uint32
	MapIDs      []uint32
	Pinned      bool     // Whether the program is pinned
	PinnedPaths []string // Paths where the program is pinned
}

// MapInfo represents information about a BPF map.
// This mirrors the structure from gobpftool's maps.MapInfo.
type MapInfo struct {
	ID          uint32
	Type        string
	Name        string
	KeySize     uint32
	ValueSize   uint32
	MaxEntries  uint32
	Flags       uint32
	MemLock     uint32
	LoadedAt    string
	UID         uint32
	Pinned      bool     // Whether the map is pinned
	PinnedPaths []string // Paths where the map is pinned
}

// MapEntry represents a key-value entry in a BPF map.
type MapEntry struct {
	Key   []byte
	Value []byte
}

// ProgService defines the interface for BPF program operations.
type ProgService interface {
	List() ([]ProgramInfo, error)
	Get(id uint32) (*ProgramInfo, error)
}

// MapsService defines the interface for BPF map operations.
type MapsService interface {
	List() ([]MapInfo, error)
	Get(id uint32) (*MapInfo, error)
	Dump(id uint32) ([]MapEntry, error)
}

// PermissionError indicates insufficient permissions for BPF operations.
type PermissionError struct {
	Err error
}

func (e *PermissionError) Error() string {
	if e.Err != nil {
		return "insufficient permissions: " + e.Err.Error()
	}
	return "insufficient permissions"
}

func (e *PermissionError) Unwrap() error {
	return e.Err
}

// IsPermissionError checks if an error is a permission error.
func IsPermissionError(err error) bool {
	var permErr *PermissionError
	return errors.As(err, &permErr)
}
