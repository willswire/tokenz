package main

import (
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type loadPage struct {
	count    int
	finished bool
}

func newLoadPage() loadPage {
	return loadPage{}
}

func (m loadPage) Init() tea.Cmd {
	return tea.Batch(
		loadTokens,
	)
}

func (m loadPage) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case finalCount:
		m.finished = true
		m.count = int(msg)
		return m, tea.Quit
	default:
		return m, nil
	}
}

func (m loadPage) View() string {
	var s string

	if m.finished {
		s = "🔑 " + strconv.Itoa(m.count) + " tokenz sourced!\n"
	} else {
		s = "🔎 Finding tokenz...\n"
	}

	return lipgloss.NewStyle().Render(s)
}

// processFinishedMsg is sent when process completes.
type finalCount int

func loadTokens() tea.Msg {
	count := load()
	return finalCount(count)
}

func load() int {
	// Command to list all generic passwords
	listCmd := exec.Command("security", "dump-keychain")
	output, err := listCmd.Output()
	if err != nil {
		log.Println("Error listing keychain items:", err)
		return 0
	}

	// Split the output into lines
	lines := strings.Split(string(output), "\n")
	count := 0

	for _, line := range lines {
		// Look for lines with the TOKENZ_ prefix in the service name
		if strings.Contains(line, "\"svce\"<blob>=\"TOKENZ_") {
			// Extract the service name
			start := strings.Index(line, "\"svce\"<blob>=\"") + len("\"svce\"<blob>=\"")
			end := strings.Index(line[start:], "\"") + start
			if start < end {
				key := line[start:end]
				value := "$(security find-generic-password -a " + os.Getenv("USER") + " -s " + key + " -w)"

				// Set the environment variable without the TOKENZ_ prefix
				envVarKey := strings.ToUpper(strings.TrimPrefix(key, "TOKENZ_"))
				err := os.Setenv(envVarKey, value)
				if err != nil {
					log.Println(err)
				}

				count++
			}
		}
	}

	return count
}
