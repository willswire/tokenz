package main

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
)

// MODEL

type helpPage struct{ text string }

func newHelpPage() helpPage {
	helpText := "" +
		"usage: tokenz <command>\n" +
		"\n" +
		"   shellenv    Load environment variables for all tokens: $(tokenz shellenv)\n" +
		"   save        Record a new token into secure storage\n" +
		"   audit       Review tokens save in secure storage\n"
	return helpPage{text: helpText}
}

func (s helpPage) Init() tea.Cmd { return nil }

// VIEW

func (s helpPage) View() string {
	return fmt.Sprintf(s.text)
}

// UPDATE

func (s helpPage) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return s, tea.Quit
}
