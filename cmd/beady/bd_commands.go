package main

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"path/filepath"
	"runtime"
)

// executeBDCommand executes a bd command with the given arguments.
// It searches for the bd binary in PATH or in the same directory as the beady executable.
// Returns the combined stdout/stderr output and any error.
func executeBDCommand(args ...string) ([]byte, error) {
	// Try to find bd in PATH first
	bdPath, err := exec.LookPath("bd")
	if err != nil {
		// If not in PATH, try same directory as beady executable
		exePath, err := getBinaryPath()
		if err == nil {
			bdPath = filepath.Join(filepath.Dir(exePath), bdBinaryName())
			// Check if it exists
			if _, err := exec.LookPath(bdPath); err != nil {
				return nil, fmt.Errorf("bd binary not found in PATH or alongside beady executable")
			}
		} else {
			return nil, fmt.Errorf("bd binary not found: %w", err)
		}
	}

	cmd := exec.Command(bdPath, args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return output, fmt.Errorf("bd command failed: %w\nOutput: %s", err, string(output))
	}
	return output, nil
}

// executeBDCommandJSON executes a bd command with --json flag and parses the JSON response.
// Returns the parsed JSON as a raw message for flexible downstream handling.
func executeBDCommandJSON(args ...string) (*json.RawMessage, error) {
	// Append --json flag if not already present
	hasJSON := false
	for _, arg := range args {
		if arg == "--json" {
			hasJSON = true
			break
		}
	}
	if !hasJSON {
		args = append(args, "--json")
	}

	output, err := executeBDCommand(args...)
	if err != nil {
		return nil, err
	}

	var result json.RawMessage
	if err := json.Unmarshal(output, &result); err != nil {
		return nil, fmt.Errorf("failed to parse JSON output: %w\nOutput: %s", err, string(output))
	}
	return &result, nil
}

// getBinaryPath returns the path to the currently running executable.
func getBinaryPath() (string, error) {
	return filepath.Abs(filepath.Dir(runtime.GOROOT()))
}

// bdBinaryName returns the appropriate bd binary name for the current platform.
func bdBinaryName() string {
	if runtime.GOOS == "windows" {
		return "bd.exe"
	}
	return "bd"
}

// BDCommandResult represents a generic result from a bd command.
type BDCommandResult struct {
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
	Error   string `json:"error,omitempty"`
}

// CreateIssueRequest represents the request body for creating a new issue.
type CreateIssueRequest struct {
	Title          string   `json:"title"`
	Description    string   `json:"description,omitempty"`
	Type           string   `json:"type,omitempty"`
	Priority       int      `json:"priority,omitempty"`
	Labels         []string `json:"labels,omitempty"`
	Assignee       string   `json:"assignee,omitempty"`
	Design         string   `json:"design,omitempty"`
	Acceptance     string   `json:"acceptance,omitempty"`
	Username       string   `json:"username,omitempty"` // For attribution
}

// UpdateStatusRequest represents the request body for updating issue status.
type UpdateStatusRequest struct {
	Status   string `json:"status"`
	Username string `json:"username,omitempty"`
}

// UpdatePriorityRequest represents the request body for updating issue priority.
type UpdatePriorityRequest struct {
	Priority int    `json:"priority"`
	Username string `json:"username,omitempty"`
}

// CloseIssueRequest represents the request body for closing an issue.
type CloseIssueRequest struct {
	Reason   string `json:"reason,omitempty"`
	Username string `json:"username,omitempty"`
}

// AddCommentRequest represents the request body for adding a comment.
type AddCommentRequest struct {
	Text     string `json:"text"`
	Username string `json:"username,omitempty"`
}

// UpdateNotesRequest represents the request body for updating issue notes.
type UpdateNotesRequest struct {
	Notes    string `json:"notes"`
	Username string `json:"username,omitempty"`
}

// AddLabelsRequest represents the request body for adding labels.
type AddLabelsRequest struct {
	Labels   []string `json:"labels"`
	Username string   `json:"username,omitempty"`
}

// AddDependencyRequest represents the request body for adding a dependency.
type AddDependencyRequest struct {
	DependencyType string `json:"dependency_type"` // e.g., "blocks", "depends-on"
	TargetID       string `json:"target_id"`
	Username       string `json:"username,omitempty"`
}
