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
}

// Initialize the input model with two text inputs
func inputModel() model {
	m := model{
		inputs: make([]textinput.Model, 2),
	}

	defaultStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("205"))

	for i := range m.inputs {
		t := textinput.New()
		t.Cursor.Style = defaultStyle.Copy()
		t.CharLimit = 32

		switch i {
		case 0:
			t.Placeholder = "Description"
			t.Focus()
			t.PromptStyle = defaultStyle
			t.TextStyle = defaultStyle
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
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			return m, tea.Quit

		case "tab", "shift+tab", "enter", "up", "down":
			if handleNavigation(msg.String(), &m) {
				return m, tea.Quit
			}

			cmd := m.updateFocus()
			cmds = append(cmds, cmd)
		}
	}

	cmd := m.updateInputs(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
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

	focusedStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	unfocusedStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("240"))

	for i := 0; i < len(m.inputs); i++ {
		if i == m.focusIndex {
			m.inputs[i].Focus()
			m.inputs[i].PromptStyle = focusedStyle
			m.inputs[i].TextStyle = focusedStyle
			cmds[i] = m.inputs[i].Focus()
		} else {
			m.inputs[i].Blur()
			m.inputs[i].PromptStyle = unfocusedStyle
			m.inputs[i].TextStyle = unfocusedStyle
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

	// Apply focused/unfocused styles to the submit button
	focusedStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	unfocusedStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("240"))

	button := unfocusedStyle.Render("[ Submit ]")
	if m.focusIndex == len(m.inputs) {
		button = focusedStyle.Render("[ Submit ]")
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
