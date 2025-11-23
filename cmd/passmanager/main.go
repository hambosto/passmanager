package main

import (
	"fmt"
	"os"
	"path/filepath"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/hambosto/passmanager/config"
	"github.com/hambosto/passmanager/internal/presentation/tui"
)

func main() {
	// Load or create config
	configPath := config.GetConfigPath()
	cfg, err := config.LoadConfig(configPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading config: %v\n", err)
		cfg = config.DefaultConfig()
	}

	// Ensure config directory exists
	configDir := filepath.Dir(configPath)
	if err := os.MkdirAll(configDir, 0o700); err != nil {
		fmt.Fprintf(os.Stderr, "Error creating config directory: %v\n", err)
		os.Exit(1)
	}

	// Save default config if it doesn't exist
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		if err := cfg.Save(configPath); err != nil {
			fmt.Fprintf(os.Stderr, "Warning: Could not save default config: %v\n", err)
		}
	}

	// Get vault path from config
	vaultPath := cfg.Storage.VaultPath

	// Initialize TUI application with config
	app := tui.NewAppWithConfig(vaultPath, cfg)

	// Run Bubble Tea program
	p := tea.NewProgram(app, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
