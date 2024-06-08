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
	listCmd := exec.Command("security", "dump-keychain")
	output, err := listCmd.Output()
	if err != nil {
		log.Println("Error listing keychain items:", err)
		return nil
	}

	lines := strings.Split(string(output), "\n")
	var exportStatements []string

	for _, line := range lines {
		if strings.Contains(line, "\"svce\"<blob>=\"TOKENZ_") {
			start := strings.Index(line, "\"svce\"<blob>=\"") + len("\"svce\"<blob>=\"")
			end := strings.Index(line[start:], "\"") + start
			if start < end {
				key := line[start:end]
				envKey := strings.ToUpper(strings.TrimPrefix(key, "TOKENZ_"))
				envVal := fmt.Sprintf("'$(security find-generic-password -a %s -s %s -w)'", os.Getenv("USER"), key)
				exportStatement := fmt.Sprintf("export %s='%s'", envKey, envVal)
				exportStatements = append(exportStatements, exportStatement)
			}
		}
	}

	return exportStatements
}
