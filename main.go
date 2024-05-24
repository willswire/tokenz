package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Model representing the application state
type model struct {
	focusIndex int
	inputs     []textinput.Model
	focusedStyle, unfocusedStyle lipgloss.Style
	focusedSubmitButton, unfocusedSubmitButton string
}

// Initialize the input model with two text inputs
func inputModel() model {
	focusedStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	unfocusedStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("240"))

	m := model{
		inputs: make([]textinput.Model, 2),
		focusedStyle: focusedStyle,
		unfocusedStyle: unfocusedStyle,
		focusedSubmitButton: focusedStyle.Copy().Render("[ Submit ]"),
		unfocusedSubmitButton: fmt.Sprintf("[ %s ]", unfocusedStyle.Render("Submit")),
	}

	for i := range m.inputs {
		t := textinput.New()
		t.Cursor.Style = focusedStyle.Copy()
		t.CharLimit = 32

		switch i {
		case 0:
			t.Placeholder = "Description"
			t.Focus()
			t.PromptStyle = focusedStyle
			t.TextStyle = focusedStyle
		case 1:
			t.Placeholder = "Value"
			t.EchoMode = textinput.EchoPassword
			t.EchoCharacter = '•'
		}

		m.inputs[i] = t
	}

	return m
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

			return m, m.updateFocus()
		}
	}

	cmd := m.updateInputs(msg)
	return m, cmd
}

// Handle navigation keys and focus changes
func handleNavigation(key string, m *model) bool {
	if key == "enter" && m.focusIndex == len(m.inputs) {
		return true
	}

	if key == "up" || key == "shift+tab" {
		m.focusIndex--
	} else {
		m.focusIndex++
	}

	if m.focusIndex > len(m.inputs) {
		m.focusIndex = 0
	} else if m.focusIndex < 0 {
		m.focusIndex = len(m.inputs)
	}

	return false
}

// Update the focus of input fields
func (m *model) updateFocus() tea.Cmd {
	cmds := make([]tea.Cmd, len(m.inputs))

	for i := 0; i < len(m.inputs); i++ {
		if i == m.focusIndex {
			cmds[i] = m.inputs[i].Focus()
			m.inputs[i].PromptStyle = m.focusedStyle
			m.inputs[i].TextStyle = m.focusedStyle
		} else {
			m.inputs[i].Blur()
			m.inputs[i].PromptStyle = lipgloss.NewStyle()
			m.inputs[i].TextStyle = lipgloss.NewStyle()
		}
	}

	return tea.Batch(cmds...)
}

// Update input fields based on user input
func (m *model) updateInputs(msg tea.Msg) tea.Cmd {
	cmds := make([]tea.Cmd, len(m.inputs))

	for i := range m.inputs {
		m.inputs[i], cmds[i] = m.inputs[i].Update(msg)
	}

	return tea.Batch(cmds...)
}

// Render the view
func (m model) View() string {
	var b strings.Builder

	for i := range m.inputs {
		b.WriteString(m.inputs[i].View())
		if i < len(m.inputs)-1 {
			b.WriteRune('\n')
		}
	}

	button := m.unfocusedSubmitButton
	if m.focusIndex == len(m.inputs) {
		button = m.focusedSubmitButton
	}
	fmt.Fprintf(&b, "\n\n%s\n\n", button)

	return b.String()
}

// Main function to run the application
func main() {
	if _, err := tea.NewProgram(inputModel()).Run(); err != nil {
		fmt.Printf("could not start program: %s\n", err)
		os.Exit(1)
	}
}
