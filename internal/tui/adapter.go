package tui

import (
	"github.com/viveksb007/gobpftool/pkg/maps"
	"github.com/viveksb007/gobpftool/pkg/prog"
)

// ProgServiceAdapter adapts gobpftool's prog.Service to our ProgService interface.
type ProgServiceAdapter struct {
	svc prog.Service
}

// NewProgServiceAdapter creates a new adapter for the prog service.
func NewProgServiceAdapter(svc prog.Service) *ProgServiceAdapter {
	return &ProgServiceAdapter{svc: svc}
}

// List returns all loaded BPF programs.
func (a *ProgServiceAdapter) List() ([]ProgramInfo, error) {
	progs, err := a.svc.List()
	if err != nil {
		return nil, err
	}

	result := make([]ProgramInfo, len(progs))
	for i, p := range progs {
		result[i] = ProgramInfo{
			ID:          p.ID,
			Type:        p.Type,
			Name:        p.Name,
			Tag:         p.Tag,
			GPL:         p.GPL,
			LoadedAt:    p.LoadedAt.Format("2006-01-02 15:04:05"),
			UID:         p.UID,
			BytesXlated: p.BytesXlated,
			BytesJIT:    p.BytesJIT,
			MemLock:     p.MemLock,
			MapIDs:      p.MapIDs,
		}
	}
	return result, nil
}

// Get returns program info by ID.
func (a *ProgServiceAdapter) Get(id uint32) (*ProgramInfo, error) {
	p, err := a.svc.GetByID(id)
	if err != nil {
		return nil, err
	}

	return &ProgramInfo{
		ID:          p.ID,
		Type:        p.Type,
		Name:        p.Name,
		Tag:         p.Tag,
		GPL:         p.GPL,
		LoadedAt:    p.LoadedAt.Format("2006-01-02 15:04:05"),
		UID:         p.UID,
		BytesXlated: p.BytesXlated,
		BytesJIT:    p.BytesJIT,
		MemLock:     p.MemLock,
		MapIDs:      p.MapIDs,
	}, nil
}

// MapsServiceAdapter adapts gobpftool's maps.Service to our MapsService interface.
type MapsServiceAdapter struct {
	svc maps.Service
}

// NewMapsServiceAdapter creates a new adapter for the maps service.
func NewMapsServiceAdapter(svc maps.Service) *MapsServiceAdapter {
	return &MapsServiceAdapter{svc: svc}
}

// List returns all loaded BPF maps.
func (a *MapsServiceAdapter) List() ([]MapInfo, error) {
	mapList, err := a.svc.List()
	if err != nil {
		return nil, err
	}

	result := make([]MapInfo, len(mapList))
	for i, m := range mapList {
		result[i] = MapInfo{
			ID:         m.ID,
			Type:       m.Type,
			Name:       m.Name,
			KeySize:    m.KeySize,
			ValueSize:  m.ValueSize,
			MaxEntries: m.MaxEntries,
			Flags:      m.Flags,
			MemLock:    m.MemLock,
			LoadedAt:   m.LoadedAt.Format("2006-01-02 15:04:05"),
			UID:        m.UID,
		}
	}
	return result, nil
}

// Get returns map info by ID.
func (a *MapsServiceAdapter) Get(id uint32) (*MapInfo, error) {
	m, err := a.svc.GetByID(id)
	if err != nil {
		return nil, err
	}

	return &MapInfo{
		ID:         m.ID,
		Type:       m.Type,
		Name:       m.Name,
		KeySize:    m.KeySize,
		ValueSize:  m.ValueSize,
		MaxEntries: m.MaxEntries,
		Flags:      m.Flags,
		MemLock:    m.MemLock,
		LoadedAt:   m.LoadedAt.Format("2006-01-02 15:04:05"),
		UID:        m.UID,
	}, nil
}

// Dump returns all entries in the map.
func (a *MapsServiceAdapter) Dump(id uint32) ([]MapEntry, error) {
	entries, err := a.svc.Dump(id)
	if err != nil {
		return nil, err
	}

	result := make([]MapEntry, len(entries))
	for i, e := range entries {
		result[i] = MapEntry{
			Key:   e.Key,
			Value: e.Value,
		}
	}
	return result, nil
}
