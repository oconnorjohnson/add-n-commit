package config

import (
	"fmt"
	"strconv"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// ConfigEditor is a TUI for editing configuration
type ConfigEditor struct {
	config      *Config
	inputs      []textinput.Model
	focusIndex  int
	saved       bool
	err         error
}

// NewConfigEditor creates a new configuration editor
func NewConfigEditor(cfg *Config) *ConfigEditor {
	inputs := make([]textinput.Model, 6)
	
	// API Key
	inputs[0] = textinput.New()
	inputs[0].Placeholder = "sk-..."
	inputs[0].SetValue(cfg.OpenAIKey)
	inputs[0].EchoMode = textinput.EchoPassword
	inputs[0].CharLimit = 100
	
	// Model
	inputs[1] = textinput.New()
	inputs[1].Placeholder = "o4-mini"
	inputs[1].SetValue(cfg.Model)
	inputs[1].CharLimit = 50
	
	// Default Mode
	inputs[2] = textinput.New()
	inputs[2].Placeholder = "interactive/all/by-file"
	inputs[2].SetValue(cfg.DefaultMode)
	inputs[2].CharLimit = 20
	
	// Temperature
	inputs[3] = textinput.New()
	inputs[3].Placeholder = "1.0"
	inputs[3].SetValue(fmt.Sprintf("%.1f", cfg.Temperature))
	inputs[3].CharLimit = 5
	
	// System Prompt All
	inputs[4] = textinput.New()
	inputs[4].Placeholder = "System prompt for all-in-one mode..."
	inputs[4].SetValue(cfg.SystemPromptAll)
	inputs[4].CharLimit = 500
	
	// System Prompt File
	inputs[5] = textinput.New()
	inputs[5].Placeholder = "System prompt for file-by-file mode..."
	inputs[5].SetValue(cfg.SystemPromptFile)
	inputs[5].CharLimit = 500
	
	// Focus on first input
	inputs[0].Focus()
	
	return &ConfigEditor{
		config: cfg,
		inputs: inputs,
	}
}

func (e *ConfigEditor) Init() tea.Cmd {
	return textinput.Blink
}

func (e *ConfigEditor) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			return e, tea.Quit
			
		case "tab", "down":
			e.focusIndex++
			if e.focusIndex >= len(e.inputs) {
				e.focusIndex = 0
			}
			e.updateFocus()
			return e, textinput.Blink
			
		case "shift+tab", "up":
			e.focusIndex--
			if e.focusIndex < 0 {
				e.focusIndex = len(e.inputs) - 1
			}
			e.updateFocus()
			return e, textinput.Blink
			
		case "ctrl+s", "enter":
			if e.focusIndex == len(e.inputs)-1 || msg.String() == "ctrl+s" {
				// Save configuration
				if err := e.saveConfig(); err != nil {
					e.err = err
				} else {
					e.saved = true
					return e, tea.Quit
				}
			} else {
				// Move to next field
				e.focusIndex++
				e.updateFocus()
				return e, textinput.Blink
			}
		}
	}
	
	// Update the focused input
	cmd := e.updateInputs(msg)
	
	return e, cmd
}

func (e *ConfigEditor) View() string {
	if e.saved {
		return lipgloss.NewStyle().
			Foreground(lipgloss.Color("42")).
			Bold(true).
			Render("✓ Configuration saved successfully!")
	}
	
	if e.err != nil {
		return lipgloss.NewStyle().
			Foreground(lipgloss.Color("196")).
			Bold(true).
			Render(fmt.Sprintf("Error: %v", e.err))
	}
	
	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("205")).
		MarginBottom(1)
	
	labelStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("86"))
	
	helpStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("241"))
	
	s := titleStyle.Render("Configure add-n-commit") + "\n\n"
	
	labels := []string{
		"OpenAI API Key:",
		"Model:",
		"Default Mode:",
		"Temperature:",
		"System Prompt (All):",
		"System Prompt (File):",
	}
	
	for i, input := range e.inputs {
		s += labelStyle.Render(labels[i]) + "\n"
		s += input.View() + "\n\n"
	}
	
	s += helpStyle.Render("Tab/↓: next field • Shift+Tab/↑: previous field • Ctrl+S: save • Esc: cancel")
	
	return s
}

func (e *ConfigEditor) updateFocus() {
	for i := range e.inputs {
		if i == e.focusIndex {
			e.inputs[i].Focus()
		} else {
			e.inputs[i].Blur()
		}
	}
}

func (e *ConfigEditor) updateInputs(msg tea.Msg) tea.Cmd {
	cmds := make([]tea.Cmd, len(e.inputs))
	
	for i := range e.inputs {
		e.inputs[i], cmds[i] = e.inputs[i].Update(msg)
	}
	
	return tea.Batch(cmds...)
}

func (e *ConfigEditor) saveConfig() error {
	// Update config from inputs
	e.config.OpenAIKey = e.inputs[0].Value()
	e.config.Model = e.inputs[1].Value()
	e.config.DefaultMode = e.inputs[2].Value()
	
	// Parse temperature
	temp, err := strconv.ParseFloat(e.inputs[3].Value(), 32)
	if err != nil {
		return fmt.Errorf("invalid temperature value: %w", err)
	}
	e.config.Temperature = float32(temp)
	
	e.config.SystemPromptAll = e.inputs[4].Value()
	e.config.SystemPromptFile = e.inputs[5].Value()
	
	// Validate default mode
	if e.config.DefaultMode != "interactive" && 
	   e.config.DefaultMode != "all" && 
	   e.config.DefaultMode != "by-file" {
		return fmt.Errorf("invalid default mode: must be 'interactive', 'all', or 'by-file'")
	}
	
	// Save to file
	return e.config.Save()
} 