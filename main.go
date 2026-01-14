package main

import (
	"fmt"
	"os"

	"github.com/viveksb007/bpftui/internal/tui"
	"github.com/viveksb007/gobpftool/pkg/maps"
	"github.com/viveksb007/gobpftool/pkg/prog"
)

func main() {
	// Create the real BPF services
	progSvc := prog.NewService()
	mapsSvc := maps.NewService()

	// Create adapters for the TUI
	progAdapter := tui.NewProgServiceAdapter(progSvc)
	mapsAdapter := tui.NewMapsServiceAdapter(mapsSvc)

	// Run the TUI with real services
	if err := tui.RunWithServices(progAdapter, mapsAdapter); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
