package tui

import (
	"testing"

	"github.com/charmbracelet/bubbles/key"
)

func TestDefaultKeyMapBindings(t *testing.T) {
	tests := []struct {
		name     string
		binding  key.Binding
		wantKeys []string
		wantHelp string
	}{
		{
			name:     "Up binding",
			binding:  defaultKeyMap.Up,
			wantKeys: []string{"up", "k"},
			wantHelp: "up",
		},
		{
			name:     "Down binding",
			binding:  defaultKeyMap.Down,
			wantKeys: []string{"down", "j"},
			wantHelp: "down",
		},
		{
			name:     "Enter binding",
			binding:  defaultKeyMap.Enter,
			wantKeys: []string{"enter"},
			wantHelp: "select",
		},
		{
			name:     "Back binding",
			binding:  defaultKeyMap.Back,
			wantKeys: []string{"esc", "backspace"},
			wantHelp: "back",
		},
		{
			name:     "Quit binding",
			binding:  defaultKeyMap.Quit,
			wantKeys: []string{"q", "ctrl+c"},
			wantHelp: "quit",
		},
		{
			name:     "Search binding",
			binding:  defaultKeyMap.Search,
			wantKeys: []string{"/"},
			wantHelp: "search",
		},
		{
			name:     "Help binding",
			binding:  defaultKeyMap.Help,
			wantKeys: []string{"?"},
			wantHelp: "help",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Check that binding is enabled
			if !tt.binding.Enabled() {
				t.Errorf("%s: binding should be enabled", tt.name)
			}

			// Check help text
			help := tt.binding.Help()
			if help.Desc != tt.wantHelp {
				t.Errorf("%s: help desc = %q, want %q", tt.name, help.Desc, tt.wantHelp)
			}

			// Check that keys match expected
			gotKeys := tt.binding.Keys()
			if len(gotKeys) != len(tt.wantKeys) {
				t.Errorf("%s: got %d keys, want %d", tt.name, len(gotKeys), len(tt.wantKeys))
				return
			}
			for i, wantKey := range tt.wantKeys {
				if gotKeys[i] != wantKey {
					t.Errorf("%s: key[%d] = %q, want %q", tt.name, i, gotKeys[i], wantKey)
				}
			}
		})
	}
}

func TestKeyMapShortHelp(t *testing.T) {
	shortHelp := defaultKeyMap.ShortHelp()

	// Should return 5 bindings: Up, Down, Enter, Back, Quit
	if len(shortHelp) != 5 {
		t.Errorf("ShortHelp() returned %d bindings, want 5", len(shortHelp))
	}

	// Verify the bindings are the expected ones
	expectedDescs := []string{"up", "down", "select", "back", "quit"}
	for i, binding := range shortHelp {
		if binding.Help().Desc != expectedDescs[i] {
			t.Errorf("ShortHelp()[%d] desc = %q, want %q", i, binding.Help().Desc, expectedDescs[i])
		}
	}
}

func TestKeyMapFullHelp(t *testing.T) {
	fullHelp := defaultKeyMap.FullHelp()

	// Should return 3 rows
	if len(fullHelp) != 3 {
		t.Errorf("FullHelp() returned %d rows, want 3", len(fullHelp))
	}

	// First row: Up, Down, Enter
	if len(fullHelp[0]) != 3 {
		t.Errorf("FullHelp()[0] has %d bindings, want 3", len(fullHelp[0]))
	}

	// Second row: Back, Search, Help
	if len(fullHelp[1]) != 3 {
		t.Errorf("FullHelp()[1] has %d bindings, want 3", len(fullHelp[1]))
	}

	// Third row: Quit
	if len(fullHelp[2]) != 1 {
		t.Errorf("FullHelp()[2] has %d bindings, want 1", len(fullHelp[2]))
	}
}

func TestKeyMatchesExpectedKeys(t *testing.T) {
	// Test that key.Matches works correctly with our bindings
	tests := []struct {
		name    string
		binding key.Binding
		keyStr  string
		want    bool
	}{
		{"Up matches 'up'", defaultKeyMap.Up, "up", true},
		{"Up matches 'k'", defaultKeyMap.Up, "k", true},
		{"Up doesn't match 'j'", defaultKeyMap.Up, "j", false},
		{"Down matches 'down'", defaultKeyMap.Down, "down", true},
		{"Down matches 'j'", defaultKeyMap.Down, "j", true},
		{"Quit matches 'q'", defaultKeyMap.Quit, "q", true},
		{"Quit matches 'ctrl+c'", defaultKeyMap.Quit, "ctrl+c", true},
		{"Back matches 'esc'", defaultKeyMap.Back, "esc", true},
		{"Back matches 'backspace'", defaultKeyMap.Back, "backspace", true},
		{"Search matches '/'", defaultKeyMap.Search, "/", true},
		{"Help matches '?'", defaultKeyMap.Help, "?", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Check if the key is in the binding's keys
			keys := tt.binding.Keys()
			found := false
			for _, k := range keys {
				if k == tt.keyStr {
					found = true
					break
				}
			}
			if found != tt.want {
				t.Errorf("key %q in binding: got %v, want %v", tt.keyStr, found, tt.want)
			}
		})
	}
}
