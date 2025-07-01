package ui

import (
	"fmt"
	"io"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/oconnorjohnson/add-n-commit/internal/git"
)

// Styles
var (
	TitleStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("205")).
		MarginBottom(1)

	SubtleStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("241"))

	SelectedStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("170")).
		Bold(true)

	NormalStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("252"))

	StatusStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("86"))

	ErrorStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("196")).
		Bold(true)

	SuccessStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("42")).
		Bold(true)
)

// Key bindings
type KeyMap struct {
	Up       key.Binding
	Down     key.Binding
	Left     key.Binding
	Right    key.Binding
	Enter    key.Binding
	Space    key.Binding
	Quit     key.Binding
	Help     key.Binding
	ToggleAll key.Binding
}

var Keys = KeyMap{
	Up: key.NewBinding(
		key.WithKeys("up", "k"),
		key.WithHelp("↑/k", "up"),
	),
	Down: key.NewBinding(
		key.WithKeys("down", "j"),
		key.WithHelp("↓/j", "down"),
	),
	Left: key.NewBinding(
		key.WithKeys("left", "h"),
		key.WithHelp("←/h", "left"),
	),
	Right: key.NewBinding(
		key.WithKeys("right", "l"),
		key.WithHelp("→/l", "right"),
	),
	Enter: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "select"),
	),
	Space: key.NewBinding(
		key.WithKeys(" "),
		key.WithHelp("space", "toggle"),
	),
	Quit: key.NewBinding(
		key.WithKeys("q", "ctrl+c"),
		key.WithHelp("q", "quit"),
	),
	Help: key.NewBinding(
		key.WithKeys("?"),
		key.WithHelp("?", "help"),
	),
	ToggleAll: key.NewBinding(
		key.WithKeys("a"),
		key.WithHelp("a", "toggle all"),
	),
}

// Helper functions
func Title(s string) string {
	return TitleStyle.Render(s)
}

func Subtle(s string) string {
	return SubtleStyle.Render(s)
}

func CenterInWindow(s string, width, height int) string {
	lines := strings.Split(s, "\n")
	maxLineWidth := 0
	for _, line := range lines {
		if len(line) > maxLineWidth {
			maxLineWidth = len(line)
		}
	}

	// Center horizontally
	leftPadding := (width - maxLineWidth) / 2
	if leftPadding < 0 {
		leftPadding = 0
	}

	// Center vertically
	topPadding := (height - len(lines)) / 2
	if topPadding < 0 {
		topPadding = 0
	}

	// Apply padding
	var result strings.Builder
	for i := 0; i < topPadding; i++ {
		result.WriteString("\n")
	}

	for _, line := range lines {
		result.WriteString(strings.Repeat(" ", leftPadding))
		result.WriteString(line)
		result.WriteString("\n")
	}

	return result.String()
}

// FileItem represents a file in the list
type FileItem struct {
	File     git.File
	Selected bool
}

func (i FileItem) FilterValue() string { return i.File.Path }

// FileDelegate handles rendering of file items
type fileDelegate struct{}

func NewFileDelegate() list.DefaultDelegate {
	d := list.NewDefaultDelegate()
	d.ShowDescription = false
	return d
}

func (d fileDelegate) Height() int                             { return 1 }
func (d fileDelegate) Spacing() int                            { return 0 }
func (d fileDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }

func (d fileDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(FileItem)
	if !ok {
		return
	}

	str := fmt.Sprintf("%s %s", i.File.Status, i.File.Path)

	fn := NormalStyle.Render
	if index == m.Index() {
		fn = func(s ...string) string {
			return SelectedStyle.Render("> " + strings.Join(s, " "))
		}
	}

	if i.Selected {
		str = "[✓] " + str
	} else {
		str = "[ ] " + str
	}

	fmt.Fprint(w, fn(str))
}

// ModeItem represents a commit mode in the list
type ModeItem struct {
	Name string
	Mode int
}

func (i ModeItem) FilterValue() string { return i.Name }

// ModeDelegate handles rendering of mode items
type modeDelegate struct{}

func NewModeDelegate() list.DefaultDelegate {
	d := list.NewDefaultDelegate()
	d.ShowDescription = false
	return d
}

func (d modeDelegate) Height() int                             { return 1 }
func (d modeDelegate) Spacing() int                            { return 0 }
func (d modeDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }

func (d modeDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(ModeItem)
	if !ok {
		return
	}

	fn := NormalStyle.Render
	if index == m.Index() {
		fn = func(s ...string) string {
			return SelectedStyle.Render("> " + strings.Join(s, " "))
		}
	}

	fmt.Fprint(w, fn(i.Name))
} 