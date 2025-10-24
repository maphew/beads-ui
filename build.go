package main

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

func main() {
	// Get current branch
	cmd := exec.Command("git", "branch", "--show-current")
	out, err := cmd.Output()
	if err != nil {
		fmt.Printf("Error getting git branch: %v\n", err)
		os.Exit(1)
	}
	branch := strings.TrimSpace(string(out))

	// Determine extension based on OS
	ext := ""
	if runtime.GOOS == "windows" {
		ext = ".exe"
	}

	// Determine output name
	output := "beady"
	if branch != "main" {
		output += "-" + branch
	}
	output += ext

	fmt.Printf("Building for branch '%s' -> %s\n", branch, output)

	// Build
	buildCmd := exec.Command("go", "build", "-o", output, "./cmd/beady")
	err = buildCmd.Run()
	if err != nil {
		fmt.Printf("Error building: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Built %s successfully\n", output)
}
