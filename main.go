package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/user"
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
type appModel struct {
	tokenNameInput  textinput.Model
	tokenValueInput textinput.Model
}

// Initialize the input model with two text inputs
func newModel() appModel {
	tokenNameInput := textinput.New()
	tokenNameInput.Placeholder = "Token Key"
	tokenNameInput.CharLimit = 32
	tokenNameInput.Focus()
	tokenNameInput.Cursor.Style = focusedStyle
	tokenNameInput.PromptStyle = focusedStyle
	tokenNameInput.TextStyle = focusedStyle

	tokenValueInput := textinput.New()
	tokenValueInput.Placeholder = "Token Value"
	tokenValueInput.CharLimit = 128
	tokenValueInput.EchoMode = textinput.EchoPassword
	tokenValueInput.EchoCharacter = '•'
	tokenValueInput.PromptStyle = blurredStyle
	tokenValueInput.TextStyle = blurredStyle
	tokenValueInput.Cursor.Style = blurredStyle

	return appModel{tokenNameInput: tokenNameInput, tokenValueInput: tokenValueInput}
}

// Initialize the application
func (m appModel) Init() tea.Cmd {
	return textinput.Blink
}

// Update the application state based on user input
func (m appModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
func handleNavigation(key string, m *appModel) bool {
	if key == "enter" && m.tokenValueInput.Focused() {
		return saveToken(m)
	}

	if key == "up" || key == "shift+tab" {
		m.tokenNameInput.Focus()
		m.tokenValueInput.Blur()
	} else {
		m.tokenValueInput.Focus()
		m.tokenNameInput.Blur()
	}

	updateStyles(m)
	return false
}

// Update styles based on focus
func updateStyles(m *appModel) {
	if m.tokenNameInput.Focused() {
		m.tokenNameInput.PromptStyle = focusedStyle
		m.tokenNameInput.TextStyle = focusedStyle
		m.tokenNameInput.Cursor.Style = focusedStyle
		m.tokenValueInput.PromptStyle = blurredStyle
		m.tokenValueInput.TextStyle = blurredStyle
		m.tokenValueInput.Cursor.Style = blurredStyle
	} else {
		m.tokenValueInput.PromptStyle = focusedStyle
		m.tokenValueInput.TextStyle = focusedStyle
		m.tokenValueInput.Cursor.Style = focusedStyle
		m.tokenNameInput.PromptStyle = blurredStyle
		m.tokenNameInput.TextStyle = blurredStyle
		m.tokenNameInput.Cursor.Style = blurredStyle
	}
}

// Update input fields based on user input
func (m *appModel) updateInputs(msg tea.Msg) tea.Cmd {
	cmds := make([]tea.Cmd, 2)
	m.tokenNameInput, cmds[0] = m.tokenNameInput.Update(msg)
	m.tokenValueInput, cmds[1] = m.tokenValueInput.Update(msg)
	return tea.Batch(cmds...)
}

// Save the token value to the keychain
func saveToken(m *appModel) bool {
	currentUser, err := user.Current()
	if err != nil {
		log.Fatalf(err.Error())
	}

	username := currentUser.Username

	cmd := exec.Command("security", "add-generic-password", "-a", username, "-s", m.tokenNameInput.Value(), "-w", m.tokenValueInput.Value())
	if err := cmd.Run(); err != nil {
		fmt.Printf("Failed to save token to keychain: %s\n", err)
		log.Fatalf(err.Error())
	}
	fmt.Printf("%s=%s\n", m.tokenNameInput.Value(), m.tokenValueInput.Value())
	return true
}

// Render the view
func (m appModel) View() string {
	var b strings.Builder
	b.WriteString(m.tokenNameInput.View())
	b.WriteRune('\n')
	b.WriteString(m.tokenValueInput.View())
	b.WriteRune('\n')
	return b.String()
}

// Main function to run the application
func main() {
	if _, err := tea.NewProgram(newModel()).Run(); err != nil {
		fmt.Printf("could not start program: %s\n", err)
		os.Exit(1)
	}
}
