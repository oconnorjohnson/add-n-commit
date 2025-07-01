package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/oconnorjohnson/add-n-commit/internal/app"
	"github.com/oconnorjohnson/add-n-commit/internal/config"
)

func main() {
	// Define command-line flags
	var (
		setKey    = flag.String("set-key", "", "Set the OpenAI API key")
		showKey   = flag.Bool("show-key", false, "Show the current OpenAI API key")
		deleteKey = flag.Bool("delete-key", false, "Delete the stored OpenAI API key")
		configure = flag.Bool("config", false, "Open configuration editor")
		showHelp  = flag.Bool("help", false, "Show help")
		version   = flag.Bool("version", false, "Show version")
	)

	flag.Parse()

	// Handle version
	if *version {
		fmt.Println("anc (add-n-commit) v1.0.0")
		fmt.Println("AI-powered git commit message generator")
		return
	}

	// Handle help
	if *showHelp {
		printHelp()
		return
	}

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		// If config doesn't exist, create a default one
		cfg = config.Default()
	}

	// Handle key management commands
	if *setKey != "" {
		if err := handleSetKey(cfg, *setKey); err != nil {
			log.Fatal(err)
		}
		return
	}

	if *showKey {
		handleShowKey(cfg)
		return
	}

	if *deleteKey {
		if err := handleDeleteKey(cfg); err != nil {
			log.Fatal(err)
		}
		return
	}

	if *configure {
		// Run configuration editor
		p := tea.NewProgram(
			config.NewConfigEditor(cfg),
			tea.WithAltScreen(),
		)
		if _, err := p.Run(); err != nil {
			log.Fatal(err)
		}
		return
	}

	// Check if we're in a git repository
	if err := checkGitRepo(); err != nil {
		log.Fatal("Error: Not in a git repository")
	}

	// Create and run the app
	p := tea.NewProgram(
		app.New(cfg),
		tea.WithAltScreen(),
		tea.WithMouseCellMotion(),
	)

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

func printHelp() {
	fmt.Println(`anc - AI-powered git commit message generator

USAGE:
    anc [options]

OPTIONS:
    --set-key <key>    Set the OpenAI API key
    --show-key         Show the current OpenAI API key (masked)
    --delete-key       Delete the stored OpenAI API key
    --config           Open interactive configuration editor
    --version          Show version information
    --help             Show this help message

INTERACTIVE MODE:
    Run 'anc' without options to enter interactive mode where you can:
    - Select files to stage
    - Choose commit message generation mode
    - Review and edit generated messages
    - Commit your changes

CONFIGURATION:
    API keys are stored in ~/.config/anc/config.json
    You can also set the OPENAI_API_KEY environment variable

EXAMPLES:
    anc                          # Enter interactive mode
    anc --set-key sk-...        # Set your OpenAI API key
    anc --show-key              # View your current API key (masked)
    anc --delete-key            # Remove stored API key
    anc --config                # Open configuration editor`)
}

func handleSetKey(cfg *config.Config, key string) error {
	// Validate key format (basic check)
	if !strings.HasPrefix(key, "sk-") {
		return fmt.Errorf("invalid API key format: OpenAI API keys should start with 'sk-'")
	}

	cfg.OpenAIKey = key
	if err := cfg.Save(); err != nil {
		return fmt.Errorf("failed to save configuration: %w", err)
	}

	fmt.Println("✓ OpenAI API key saved successfully")
	fmt.Println("Configuration stored in ~/.config/anc/config.json")
	return nil
}

func handleShowKey(cfg *config.Config) {
	if cfg.OpenAIKey == "" {
		fmt.Println("No OpenAI API key configured")
		fmt.Println("Use 'anc --set-key <key>' to set one")
		return
	}

	// Mask the key for security
	maskedKey := maskAPIKey(cfg.OpenAIKey)
	fmt.Printf("Current OpenAI API key: %s\n", maskedKey)
}

func handleDeleteKey(cfg *config.Config) error {
	if cfg.OpenAIKey == "" {
		fmt.Println("No OpenAI API key to delete")
		return nil
	}

	cfg.OpenAIKey = ""
	if err := cfg.Save(); err != nil {
		return fmt.Errorf("failed to save configuration: %w", err)
	}

	fmt.Println("✓ OpenAI API key deleted successfully")
	return nil
}

func maskAPIKey(key string) string {
	if len(key) <= 8 {
		return "********"
	}
	// Show first 3 and last 4 characters
	return key[:3] + "..." + key[len(key)-4:]
} 