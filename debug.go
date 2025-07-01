//go:build debug
// +build debug

package main

import (
	"fmt"
	"log"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/oconnorjohnson/add-n-commit/internal/app"
	"github.com/oconnorjohnson/add-n-commit/internal/config"
)

func main() {
	// Enable debug logging
	if f, err := tea.LogToFile("debug.log", "debug"); err != nil {
		fmt.Println("fatal:", err)
		os.Exit(1)
	} else {
		defer f.Close()
	}

	log.Println("Starting anc in debug mode...")

	// Check if we're in a git repository
	if err := checkGitRepo(); err != nil {
		log.Fatal("Error: Not in a git repository")
	}

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Printf("Warning: Could not load config: %v", err)
		cfg = config.Default()
	}

	log.Printf("Config loaded: API Key present: %v", cfg.OpenAIKey != "")

	// Create and run the app
	p := tea.NewProgram(
		app.New(cfg),
		tea.WithAltScreen(),
		tea.WithMouseCellMotion(),
	)

	log.Println("Starting TUI...")
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}

func checkGitRepo() error {
	_, err := os.Stat(".git")
	if os.IsNotExist(err) {
		return fmt.Errorf("not in a git repository")
	}
	return err
} 