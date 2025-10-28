// Build is a utility that compiles the beady binary with branch-aware naming.
// It appends the current git branch name to the output binary (e.g., beady-feature)
// unless building from the main branch.
package main

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

// main builds the beady binary and names the output by appending the current Git branch (except when on "main") and an OS-specific extension.
//
// It obtains the current Git branch, sanitizes the branch by replacing characters not in [A-Za-z0-9-_] with '_', appends "-<branch>" to the base name "beady" when the branch is not "main", adds ".exe" on Windows, and runs `go build -o <output> ./cmd/beady` to produce the binary. If obtaining the branch or the build process fails, it prints an error and exits with a non-zero status.
func main() {
	// Get current branch
	cmd := exec.Command("git", "branch", "--show-current")
	out, err := cmd.Output()
	if err != nil {
		fmt.Printf("Error getting git branch: %v\nNote: build.go requires a git repository with an active branch.\n", err)
		os.Exit(1)
	}
	branch := strings.TrimSpace(string(out))

	// Determine extension based on OS
	ext := ""
	if runtime.GOOS == "windows" {
		ext = ".exe"
	}

	// Sanitize branch name to prevent path traversal and special characters
	branch = strings.Map(func(r rune) rune {
		if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') || r == '-' || r == '_' {
			return r
		}
		return '_'
	}, branch)

	// Determine output name
	output := "beady"
	if branch != "main" {
		output += "-" + branch
	}
	output += ext

	fmt.Printf("Building for branch '%s' -> %s\n", branch, output)

	// Build
	buildCmd := exec.Command("go", "build", "-o", "bin/"+output, "./cmd/beady")
	buildOutput, err := buildCmd.CombinedOutput()
	if err != nil {
		fmt.Printf("Error building: %v\n%s\n", err, string(buildOutput))
		os.Exit(1)
	}

	fmt.Printf("Built %s successfully\n", output)
}
