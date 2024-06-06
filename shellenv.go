package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
)

func shellenv() {
	exportStatements := generateExportStatements()
	for _, statement := range exportStatements {
		fmt.Println(statement)
	}
}

func generateExportStatements() []string {
	// Command to list all generic passwords
	listCmd := exec.Command("security", "dump-keychain")
	output, err := listCmd.Output()
	if err != nil {
		log.Println("Error listing keychain items:", err)
		return nil
	}

	// Split the output into lines
	lines := strings.Split(string(output), "\n")
	var exportStatements []string

	for _, line := range lines {
		// Look for lines with the TOKENZ_ prefix in the service name
		if strings.Contains(line, "\"svce\"<blob>=\"TOKENZ_") {
			// Extract the service name
			start := strings.Index(line, "\"svce\"<blob>=\"") + len("\"svce\"<blob>=\"")
			end := strings.Index(line[start:], "\"") + start
			if start < end {
				key := line[start:end]
				envVal := fmt.Sprintf("'$(security find-generic-password -a %s -s %s -w)'", os.Getenv("USER"), key)

				// Set the environment variable without the TOKENZ_ prefix
				envKey := strings.ToUpper(strings.TrimPrefix(key, "TOKENZ_"))
				exportStatement := fmt.Sprintf("export %s='%s'", envKey, envVal)
				exportStatements = append(exportStatements, exportStatement)
			}
		}
	}

	return exportStatements
}
