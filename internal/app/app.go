package app

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/oconnorjohnson/add-n-commit/internal/config"
	"github.com/oconnorjohnson/add-n-commit/internal/git"
	"github.com/oconnorjohnson/add-n-commit/internal/openai"
	"github.com/oconnorjohnson/add-n-commit/internal/ui"
)

type state int

const (
	stateFileSelection state = iota
	stateModeSelection
	stateGenerating
	stateReviewing
	stateEditing
	stateCommitting
	stateSuccess
	stateError
	stateConfig
)

type commitMode int

const (
	modeAllInOne commitMode = iota
	modeByFile
	modeCustomPrompt
)

type Model struct {
	state       state
	config      *config.Config
	files       []git.File
	fileList    list.Model
	modeList    list.Model
	spinner     spinner.Model
	textarea    textarea.Model
	textinput   textinput.Model
	apiKeyInput textinput.Model
	
	selectedFiles   []string
	selectedMode    commitMode
	generatedMsg    string
	customPrompt    string
	errorMsg        string
	successMsg      string
	
	width  int
	height int
	
	openaiClient *openai.Client
}

// New creates a new app model
func New(cfg *config.Config) *Model {
	m := &Model{
		state:    stateFileSelection,
		config:   cfg,
		spinner:  spinner.New(),
		textarea: textarea.New(),
	}
	
	// Initialize text input for custom prompt
	ti := textinput.New()
	ti.Placeholder = "Enter additional context for commit message generation..."
	ti.CharLimit = 200
	m.textinput = ti
	
	// Initialize API key input
	apiKeyInput := textinput.New()
	apiKeyInput.Placeholder = "Enter your OpenAI API key..."
	apiKeyInput.CharLimit = 100
	apiKeyInput.EchoMode = textinput.EchoPassword
	m.apiKeyInput = apiKeyInput
	
	// Set up spinner
	m.spinner.Spinner = spinner.Dot
	
	// Initialize OpenAI client if API key is available
	if cfg.OpenAIKey != "" {
		m.openaiClient = openai.NewClient(cfg.OpenAIKey, cfg.Model, cfg.Temperature)
	}
	
	return m
}

func (m *Model) Init() tea.Cmd {
	// Check if we need to configure API key first
	if m.config.OpenAIKey == "" {
		m.state = stateConfig
		m.apiKeyInput.Focus()
		return textinput.Blink
	}
	
	return tea.Batch(
		m.loadFiles,
		m.spinner.Tick,
	)
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil
		
	case tea.KeyMsg:
		// Handle Ctrl+C globally
		if key.Matches(msg, key.NewBinding(key.WithKeys("ctrl+c"))) {
			return m, tea.Quit
		}
		
		switch m.state {
		case stateConfig:
			return m.updateConfig(msg)
		case stateFileSelection:
			return m.updateFileSelection(msg)
		case stateModeSelection:
			return m.updateModeSelection(msg)
		case stateReviewing:
			return m.updateReviewing(msg)
		case stateEditing:
			return m.updateEditing(msg)
		case stateSuccess, stateError:
			return m, tea.Quit
		}
		
	case filesLoadedMsg:
		m.files = msg.files
		m.setupFileList()
		return m, nil
		
	case commitModeSelectedMsg:
		m.selectedMode = msg.mode
		if m.selectedMode == modeCustomPrompt {
			m.state = stateEditing
			m.textinput.Focus()
			return m, textinput.Blink
		}
		m.state = stateGenerating
		return m, m.generateCommitMessage()
		
	case commitMessageGeneratedMsg:
		m.generatedMsg = msg.message
		m.state = stateReviewing
		m.textarea.SetValue(m.generatedMsg)
		return m, nil
		
	case commitSuccessMsg:
		m.successMsg = "✓ Changes committed successfully!"
		m.state = stateSuccess
		return m, nil
		
	case errorMsg:
		m.errorMsg = msg.err.Error()
		m.state = stateError
		return m, nil
	}
	
	// Update sub-components
	var cmd tea.Cmd
	switch m.state {
	case stateFileSelection:
		m.fileList, cmd = m.fileList.Update(msg)
	case stateModeSelection:
		m.modeList, cmd = m.modeList.Update(msg)
	case stateGenerating:
		m.spinner, cmd = m.spinner.Update(msg)
	case stateEditing:
		if m.selectedMode == modeCustomPrompt {
			m.textinput, cmd = m.textinput.Update(msg)
		} else {
			m.textarea, cmd = m.textarea.Update(msg)
		}
	case stateConfig:
		m.apiKeyInput, cmd = m.apiKeyInput.Update(msg)
	}
	
	return m, cmd
}

func (m *Model) View() string {
	if m.width == 0 {
		return "Loading..."
	}
	
	var content string
	
	switch m.state {
	case stateConfig:
		content = m.viewConfig()
	case stateFileSelection:
		content = m.viewFileSelection()
	case stateModeSelection:
		content = m.viewModeSelection()
	case stateGenerating:
		content = m.viewGenerating()
	case stateReviewing:
		content = m.viewReviewing()
	case stateEditing:
		content = m.viewEditing()
	case stateSuccess:
		content = m.viewSuccess()
	case stateError:
		content = m.viewError()
	}
	
	return ui.CenterInWindow(content, m.width, m.height)
}

// View helpers
func (m *Model) viewConfig() string {
	return fmt.Sprintf(
		"%s\n\n%s\n\n%s",
		ui.Title("Configure OpenAI API Key"),
		m.apiKeyInput.View(),
		ui.Subtle("Press Enter to save, Esc to quit"),
	)
}

func (m *Model) viewFileSelection() string {
	if len(m.files) == 0 {
		return ui.Title("No changes detected") + "\n\n" + ui.Subtle("Make some changes and run again!")
	}
	
	return fmt.Sprintf(
		"%s\n\n%s\n\n%s",
		ui.Title("Select files to stage"),
		m.fileList.View(),
		ui.Subtle("Space: toggle, a: all/none, Enter: continue, q: quit"),
	)
}

func (m *Model) viewModeSelection() string {
	return fmt.Sprintf(
		"%s\n\n%s\n\n%s",
		ui.Title("Select commit message mode"),
		m.modeList.View(),
		ui.Subtle("Enter: select, q: quit"),
	)
}

func (m *Model) viewGenerating() string {
	return fmt.Sprintf(
		"%s %s",
		m.spinner.View(),
		ui.Title("Generating commit message..."),
	)
}

func (m *Model) viewReviewing() string {
	return fmt.Sprintf(
		"%s\n\n%s\n\n%s",
		ui.Title("Review commit message"),
		m.textarea.View(),
		ui.Subtle("Enter: commit, e: edit, r: regenerate, q: quit"),
	)
}

func (m *Model) viewEditing() string {
	if m.selectedMode == modeCustomPrompt {
		return fmt.Sprintf(
			"%s\n\n%s\n\n%s",
			ui.Title("Enter additional context"),
			m.textinput.View(),
			ui.Subtle("Enter: generate message, Esc: back"),
		)
	}
	
	return fmt.Sprintf(
		"%s\n\n%s\n\n%s",
		ui.Title("Edit commit message"),
		m.textarea.View(),
		ui.Subtle("Ctrl+Enter: commit, Esc: cancel"),
	)
}

func (m *Model) viewSuccess() string {
	style := lipgloss.NewStyle().
		Foreground(lipgloss.Color("42")).
		Bold(true)
		
	return fmt.Sprintf(
		"%s\n\n%s",
		style.Render(m.successMsg),
		ui.Subtle("Press any key to exit"),
	)
}

func (m *Model) viewError() string {
	style := lipgloss.NewStyle().
		Foreground(lipgloss.Color("196")).
		Bold(true)
		
	return fmt.Sprintf(
		"%s\n\n%s\n\n%s",
		ui.Title("Error"),
		style.Render(m.errorMsg),
		ui.Subtle("Press any key to exit"),
	)
}

// Helper methods
func (m *Model) setupFileList() {
	items := make([]list.Item, len(m.files))
	for i, f := range m.files {
		items[i] = ui.FileItem{
			File:     f,
			Selected: false,
		}
	}
	
	delegate := ui.NewFileDelegate()
	m.fileList = list.New(items, delegate, m.width-4, m.height-10)
	m.fileList.Title = "Files"
	m.fileList.SetShowStatusBar(false)
	m.fileList.SetFilteringEnabled(false)
	m.fileList.Styles.Title = ui.TitleStyle
}

func (m *Model) setupModeList() {
	items := []list.Item{
		ui.ModeItem{Name: "All-in-one summary", Mode: int(modeAllInOne)},
		ui.ModeItem{Name: "File-by-file summary", Mode: int(modeByFile)},
		ui.ModeItem{Name: "Custom prompt", Mode: int(modeCustomPrompt)},
	}
	
	delegate := ui.NewModeDelegate()
	m.modeList = list.New(items, delegate, m.width-4, 10)
	m.modeList.Title = "Modes"
	m.modeList.SetShowStatusBar(false)
	m.modeList.SetFilteringEnabled(false)
	m.modeList.Styles.Title = ui.TitleStyle
}

// Update helpers
func (m *Model) updateConfig(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.Type {
	case tea.KeyEnter:
		apiKey := m.apiKeyInput.Value()
		if apiKey == "" {
			m.errorMsg = "API key cannot be empty"
			m.state = stateError
			return m, nil
		}
		
		m.config.OpenAIKey = apiKey
		if err := m.config.Save(); err != nil {
			m.errorMsg = fmt.Sprintf("Failed to save config: %v", err)
			m.state = stateError
			return m, nil
		}
		
		m.openaiClient = openai.NewClient(apiKey, m.config.Model, m.config.Temperature)
		m.state = stateFileSelection
		return m, m.loadFiles
		
	case tea.KeyEsc:
		return m, tea.Quit
	}
	
	var cmd tea.Cmd
	m.apiKeyInput, cmd = m.apiKeyInput.Update(msg)
	return m, cmd
}

func (m *Model) updateFileSelection(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "q":
		return m, tea.Quit
		
	case " ":
		if i, ok := m.fileList.SelectedItem().(ui.FileItem); ok {
			idx := m.fileList.Index()
			items := m.fileList.Items()
			i.Selected = !i.Selected
			items[idx] = i
			m.fileList.SetItems(items)
		}
		
	case "a":
		// Toggle all
		items := m.fileList.Items()
		allSelected := true
		for _, item := range items {
			if fi, ok := item.(ui.FileItem); ok && !fi.Selected {
				allSelected = false
				break
			}
		}
		
		for i, item := range items {
			if fi, ok := item.(ui.FileItem); ok {
				fi.Selected = !allSelected
				items[i] = fi
			}
		}
		m.fileList.SetItems(items)
		
	case "enter":
		// Collect selected files
		m.selectedFiles = []string{}
		for _, item := range m.fileList.Items() {
			if fi, ok := item.(ui.FileItem); ok && fi.Selected {
				m.selectedFiles = append(m.selectedFiles, fi.File.Path)
			}
		}
		
		if len(m.selectedFiles) == 0 {
			m.errorMsg = "No files selected"
			m.state = stateError
			return m, nil
		}
		
		// Stage selected files
		if err := git.StageFiles(m.selectedFiles); err != nil {
			m.errorMsg = fmt.Sprintf("Failed to stage files: %v", err)
			m.state = stateError
			return m, nil
		}
		
		m.setupModeList()
		m.state = stateModeSelection
		return m, nil
	}
	
	return m, nil
}

func (m *Model) updateModeSelection(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "q":
		return m, tea.Quit
		
	case "enter":
		if i, ok := m.modeList.SelectedItem().(ui.ModeItem); ok {
			return m, func() tea.Msg {
				return commitModeSelectedMsg{mode: commitMode(i.Mode)}
			}
		}
	}
	
	return m, nil
}

func (m *Model) updateReviewing(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "q":
		return m, tea.Quit
		
	case "enter":
		return m, m.commitChanges()
		
	case "e":
		m.state = stateEditing
		m.textarea.Focus()
		return m, textarea.Blink
		
	case "r":
		m.state = stateGenerating
		return m, m.generateCommitMessage()
	}
	
	return m, nil
}

func (m *Model) updateEditing(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	if m.selectedMode == modeCustomPrompt {
		switch msg.Type {
		case tea.KeyEnter:
			m.customPrompt = m.textinput.Value()
			m.state = stateGenerating
			return m, m.generateCommitMessage()
			
		case tea.KeyEsc:
			m.state = stateModeSelection
			return m, nil
		}
		return m, nil
	}
	
	switch msg.Type {
	case tea.KeyCtrlC:
		return m, tea.Quit
		
	case tea.KeyEsc:
		m.state = stateReviewing
		m.textarea.SetValue(m.generatedMsg)
		return m, nil
		
	case tea.KeyCtrlS, tea.KeyCtrlD:
		m.generatedMsg = m.textarea.Value()
		return m, m.commitChanges()
	}
	
	return m, nil
}

// Commands
func (m *Model) loadFiles() tea.Msg {
	files, err := git.GetStatus()
	if err != nil {
		return errorMsg{err: err}
	}
	
	// Filter out already staged files if not auto-staging all
	if m.config != nil && !m.config.AutoStageAll {
		var unstaged []git.File
		for _, f := range files {
			if !f.IsStaged {
				unstaged = append(unstaged, f)
			}
		}
		files = unstaged
	}
	
	return filesLoadedMsg{files: files}
}

func (m *Model) generateCommitMessage() tea.Cmd {
	return func() tea.Msg {
		if m.openaiClient == nil {
			return errorMsg{err: fmt.Errorf("OpenAI client not initialized. Please set your API key.")}
		}
		
		var message string
		var err error
		
		switch m.selectedMode {
		case modeAllInOne:
			diff, err := git.GetStagedDiff()
			if err != nil {
				return errorMsg{err: err}
			}
			
			message, err = m.openaiClient.GenerateCommitMessage(
				m.config.SystemPromptAll,
				diff,
			)
			
		case modeByFile:
			files, err := git.GetStagedFiles()
			if err != nil {
				return errorMsg{err: err}
			}
			
			var messages []string
			for _, file := range files {
				diff, err := git.GetStagedDiffForFile(file)
				if err != nil {
					continue
				}
				
				msg, err := m.openaiClient.GenerateCommitMessage(
					m.config.SystemPromptFile,
					diff,
				)
				if err != nil {
					continue
				}
				
				messages = append(messages, fmt.Sprintf("%s: %s", file, msg))
			}
			
			message = strings.Join(messages, "\n")
			
		case modeCustomPrompt:
			diff, err := git.GetStagedDiff()
			if err != nil {
				return errorMsg{err: err}
			}
			
			message, err = m.openaiClient.GenerateCommitMessageWithContext(
				m.config.SystemPromptAll,
				diff,
				m.customPrompt,
			)
		}
		
		if err != nil {
			return errorMsg{err: err}
		}
		
		return commitMessageGeneratedMsg{message: message}
	}
}

func (m *Model) commitChanges() tea.Cmd {
	return func() tea.Msg {
		message := m.textarea.Value()
		if message == "" {
			message = m.generatedMsg
		}
		
		if err := git.Commit(message); err != nil {
			return errorMsg{err: err}
		}
		
		return commitSuccessMsg{}
	}
}

// Message types
type filesLoadedMsg struct {
	files []git.File
}

type commitModeSelectedMsg struct {
	mode commitMode
}

type commitMessageGeneratedMsg struct {
	message string
}

type commitSuccessMsg struct{}

type errorMsg struct {
	err error
} 