package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	var (
		command string
		model   tea.Model
		opts    []tea.ProgramOption
	)

	if len(os.Args) < 2 {
		command = "help" // Default command
	} else {
		command = os.Args[1]
	}

	switch command {
	case "save":
		model = newSavePage()
	case "load":
		model = newLoadPage()
	default:
		model = newHelpPage()
	}

	program := tea.NewProgram(model, opts...)
	if _, err := program.Run(); err != nil {
		fmt.Printf("could not start program: %s\n", err)
		panic(err)
	}
}
