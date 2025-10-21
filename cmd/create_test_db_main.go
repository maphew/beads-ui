package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/steveyegge/beads"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run test_db.go <db_path>")
		os.Exit(1)
	}

	dbPath := os.Args[1]
	
	// Create .beads directory if it doesn't exist
	beadsDir := filepath.Dir(dbPath)
	if err := os.MkdirAll(beadsDir, 0750); err != nil {
		fmt.Fprintf(os.Stderr, "Error: failed to create %s directory: %v\n", beadsDir, err)
		os.Exit(1)
	}

	// Create database
	store, err := beads.NewSQLiteStorage(dbPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: failed to create database: %v\n", err)
		os.Exit(1)
	}

	// Create a few test issues
	ctx := context.Background()
	
	now := time.Now()
	
	// Create some test issues
	issues := []*beads.Issue{
		{
			ID:          "test-1",
			Title:       "Test issue 1",
			Description: "This is a test issue",
			Status:      beads.StatusOpen,
			Priority:    1,
			IssueType:   beads.TypeTask,
			CreatedAt:   now,
			UpdatedAt:   now,
		},
		{
			ID:          "test-2",
			Title:       "Test issue 2",
			Description: "This is another test issue",
			Status:      beads.StatusInProgress,
			Priority:    2,
			IssueType:   beads.TypeFeature,
			CreatedAt:   now,
			UpdatedAt:   now,
		},
		{
			ID:          "test-3",
			Title:       "Test issue 3",
			Description: "This is a third test issue",
			Status:      beads.StatusClosed,
			Priority:    0,
			IssueType:   beads.TypeBug,
			CreatedAt:   now,
			UpdatedAt:   now,
			ClosedAt:    &now,
		},
	}

	for _, issue := range issues {
		// Create the issue in the database
		if err := store.CreateIssue(ctx, issue, "webui-test"); err != nil {
			fmt.Fprintf(os.Stderr, "Error: failed to create issue %s: %v\n", issue.ID, err)
			os.Exit(1)
		}
	}

	fmt.Println("Test database created successfully with sample issues")
}