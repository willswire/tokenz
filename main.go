package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Define styles
var (
	focusedStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	blurredStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
)

// Model representing the application state
type model struct {
	key   textinput.Model
	value textinput.Model
}

// Initialize the input model with two text inputs
func createModel() model {
	key := textinput.New()
	key.Placeholder = "key"
	key.CharLimit = 32
	key.Focus()
	key.Cursor.Style = focusedStyle
	key.PromptStyle = focusedStyle
	key.TextStyle = focusedStyle

	value := textinput.New()
	value.Placeholder = "value"
	value.CharLimit = 128
	value.EchoMode = textinput.EchoPassword
	value.EchoCharacter = '•'
	value.PromptStyle = blurredStyle
	value.TextStyle = blurredStyle
	value.Cursor.Style = blurredStyle

	return model{key: key, value: value}
}

// Initialize the application
func (m model) Init() tea.Cmd {
	return textinput.Blink
}

// Update the application state based on user input
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			return m, tea.Quit
		case "tab", "shift+tab", "enter", "up", "down":
			if handleNavigation(msg.String(), &m) {
				return m, tea.Quit
			}
		}
	}

	cmd := m.updateInputs(msg)
	return m, cmd
}

// Handle navigation keys and focus changes
func handleNavigation(key string, m *model) bool {
	if key == "enter" && m.value.Focused() {
		return saveInputs(m)
	}

	if key == "up" || key == "shift+tab" {
		m.key.Focus()
		m.value.Blur()
	} else {
		m.value.Focus()
		m.key.Blur()
	}

	updateStyles(m)
	return false
}

// Update styles based on focus
func updateStyles(m *model) {
	if m.key.Focused() {
		m.key.PromptStyle = focusedStyle
		m.key.TextStyle = focusedStyle
		m.key.Cursor.Style = focusedStyle
		m.value.PromptStyle = blurredStyle
		m.value.TextStyle = blurredStyle
		m.value.Cursor.Style = blurredStyle
	} else {
		m.value.PromptStyle = focusedStyle
		m.value.TextStyle = focusedStyle
		m.value.Cursor.Style = focusedStyle
		m.key.PromptStyle = blurredStyle
		m.key.TextStyle = blurredStyle
		m.key.Cursor.Style = blurredStyle
	}
}

// Update input fields based on user input
func (m *model) updateInputs(msg tea.Msg) tea.Cmd {
	cmds := make([]tea.Cmd, 2)
	m.key, cmds[0] = m.key.Update(msg)
	m.value, cmds[1] = m.value.Update(msg)
	return tea.Batch(cmds...)
}

// Save the value to the keychain
func saveInputs(m *model) bool {
	fmt.Printf("%s=%s\n", m.key.Value(), m.value.Value())
	return true
}

// Render the view
func (m model) View() string {
	var b strings.Builder
	b.WriteString(m.key.View())
	b.WriteRune('\n')
	b.WriteString(m.value.View())
	b.WriteRune('\n')
	return b.String()
}

// Main function to run the application
func main() {
	if _, err := tea.NewProgram(createModel()).Run(); err != nil {
		fmt.Printf("could not start program: %s\n", err)
		os.Exit(1)
	}
}
