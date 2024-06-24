package main

import (
	"log"
	"os/exec"
	"strings"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var baseStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.NormalBorder()).
	BorderForeground(lipgloss.Color("240"))

type auditPage struct {
	table table.Model
}

func (m auditPage) Init() tea.Cmd { return nil }

func (m auditPage) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			if m.table.Focused() {
				m.table.Blur()
			} else {
				m.table.Focus()
			}
		case "q", "ctrl+c":
			return m, tea.Quit
		case "enter":
			return m, tea.Batch(
				tea.Printf("Let's go to %s!", m.table.SelectedRow()[1]),
			)
		}
	}
	m.table, cmd = m.table.Update(msg)
	return m, cmd
}

func (m auditPage) View() string {
	return baseStyle.Render(m.table.View()) + "\n"
}

func fetchTokenKeys() []string {
	listCmd := exec.Command("security", "dump-keychain")
	output, err := listCmd.Output()
	if err != nil {
		log.Println("Error listing keychain items:", err)
		return nil
	}

	lines := strings.Split(string(output), "\n")
	var tokenKeys []string

	for _, line := range lines {
		if strings.Contains(line, "\"svce\"<blob>=\"TOKENZ_") {
			start := strings.Index(line, "\"svce\"<blob>=\"") + len("\"svce\"<blob>=\"")
			end := strings.Index(line[start:], "\"") + start
			if start < end {
				key := line[start:end]
				envKey := strings.ToUpper(strings.TrimPrefix(key, "TOKENZ_"))
				tokenKeys = append(tokenKeys, envKey)
			}
		}
	}

	return tokenKeys
}

func newAuditPage() auditPage {
	columns := []table.Column{
		{Title: "Key", Width: 20},
		{Title: "Last Updated", Width: 10},
	}

	rows := []table.Row{
		{"1", "Tokyo", "Japan", "37,274,000"},
		{"2", "Delhi", "India", "32,065,760"},
		{"3", "Shanghai", "China", "28,516,904"},
		{"4", "Dhaka", "Bangladesh", "22,478,116"},
		{"5", "São Paulo", "Brazil", "22,429,800"},
		{"6", "Mexico City", "Mexico", "22,085,140"},
		{"7", "Cairo", "Egypt", "21,750,020"},
		{"8", "Beijing", "China", "21,333,332"},
		{"9", "Mumbai", "India", "20,961,472"},
		{"10", "Osaka", "Japan", "19,059,856"},
	}

	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(7),
	)

	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		BorderBottom(true).
		Bold(false)
	s.Selected = s.Selected.
		Foreground(lipgloss.Color("229")).
		Background(lipgloss.Color("57")).
		Bold(false)
	t.SetStyles(s)

	m := auditPage{t}
	return m
}
