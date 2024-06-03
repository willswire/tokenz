package main

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

// MODEL

type helpPage struct{ text string }

func newHelpPage() helpPage {
	helpText := "" +
		"Usage: tokenz [command]\n" +
		"\n" +
		" load	Load environment variables for all tokens\n" +
		" save	Record a new token into secure storage\n"
	return helpPage{text: helpText}
}

func (s helpPage) Init() tea.Cmd { return nil }

// VIEW

func (s helpPage) View() string {
	textLen := len(s.text)
	topAndBottomBar := strings.Repeat("*", textLen+4)
	return fmt.Sprintf(
		"%s\n%s\n%s\n\n",
		topAndBottomBar, s.text, topAndBottomBar,
	)
}

// UPDATE

func (s helpPage) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return s, tea.Quit
}
