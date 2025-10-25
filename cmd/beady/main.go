package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"html/template"
	"io/fs"
	"log"
	"net"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/gorilla/websocket"
	"github.com/maphew/beads-ui/assets/beady"
	"github.com/steveyegge/beads"
)

var embedFS = beady.FS

var tmplFS fs.FS

// Pre-parse templates at package init for performance
var (
	tmplIndex       *template.Template
	tmplDetail      *template.Template
	tmplGraph       *template.Template
	tmplReady       *template.Template
	tmplBlocked     *template.Template
	tmplIssuesTbody *template.Template
)

func init() {
	tmplFS = embedFS
	flag.BoolVar(&devMode, "dev", false, "")
	flag.BoolVar(&devMode, "d", false, "Enable development mode with live reload")
	// Templates will be parsed after flag parsing
}

func parseTemplates() {
	tmplIndex = template.Must(template.ParseFS(tmplFS, "templates/index.html"))
	tmplDetail = template.Must(template.ParseFS(tmplFS, "templates/detail.html"))
	tmplGraph = template.Must(template.ParseFS(tmplFS, "templates/graph.html"))
	tmplReady = template.Must(template.ParseFS(tmplFS, "templates/ready.html"))
	tmplBlocked = template.Must(template.ParseFS(tmplFS, "templates/blocked.html"))
	tmplIssuesTbody = template.Must(template.ParseFS(tmplFS, "templates/issues_tbody.html"))
}

var store beads.Storage

var devMode bool

var help = flag.Bool("help", false, "Show help")

func printUsage() {
	fmt.Fprintf(os.Stderr, "Usage: %s [database-path] [port] [-d] [--help]\n", os.Args[0])
	fmt.Fprintf(os.Stderr, "Options:\n")
	fmt.Fprintf(os.Stderr, "  -d, --dev       Enable development mode with live reload\n")
	fmt.Fprintf(os.Stderr, "  -h, --help      Show help\n")
	fmt.Fprintf(os.Stderr, "Examples:\n")
	fmt.Fprintf(os.Stderr, "  %s                    # autodiscover database\n", os.Args[0])
	fmt.Fprintf(os.Stderr, "  %s .beads/name.db   # specify database path\n", os.Args[0])
	fmt.Fprintf(os.Stderr, "  %s .beads/name.db 8080  # specify path and port\n", os.Args[0])
	fmt.Fprintf(os.Stderr, "  %s -d .beads/name.db 8080  # enable live reload\n", os.Args[0])
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true }, // Allow all origins for dev
}

var clients = make(map[*websocket.Conn]bool)
var clientsMu sync.Mutex
var shutdownTimer *time.Timer

// broadcast sends the given text message to all registered WebSocket clients.
// If writing to a client fails, the function closes that connection and removes it from the client set.
func broadcast(message string) {
	clientsMu.Lock()
	defer clientsMu.Unlock()
	for conn := range clients {
		if err := conn.WriteMessage(websocket.TextMessage, []byte(message)); err != nil {
			conn.Close()
			delete(clients, conn)
		}
	}
}

// handleWS upgrades an HTTP connection to a WebSocket and manages the connection lifecycle for live reloads.
// 
// It registers the new client and cancels any pending shutdown timer while connected. When the client
// disconnects it is removed; if no clients remain a 5-second timer is started to terminate the process.
// The handler keeps the connection alive by continuously reading messages until an error occurs.
func handleWS(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("WebSocket upgrade failed:", err)
		return
	}
	defer conn.Close()

	clientsMu.Lock()
	clients[conn] = true
	// Cancel shutdown timer if running
	if shutdownTimer != nil {
		if shutdownTimer.Stop() {
			shutdownTimer = nil
		}
	}
	clientsMu.Unlock()

	defer func() {
		clientsMu.Lock()
		delete(clients, conn)
		// Start shutdown timer if no clients left
		if len(clients) == 0 {
			shutdownTimer = time.AfterFunc(5*time.Second, func() {
				clientsMu.Lock()
				defer clientsMu.Unlock()
				if len(clients) != 0 {
					return
				}
				log.Println("No clients connected, shutting down...")
				os.Exit(0)
			})
		}
		clientsMu.Unlock()
	}()

	for {
		_, _, err := conn.ReadMessage()
		if err != nil {
			break
		}
	}
}

// startFileWatcher watches the embedded assets/beady/templates and assets/beady/static directories for file changes and triggers live-reload actions.
// 
// When a change is detected it logs the change, re-parses HTML templates if a template file was modified, and broadcasts a "reload" message to connected WebSocket clients.
// It also logs watcher errors.
func startFileWatcher() {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal("Failed to create file watcher:", err)
	}
	defer watcher.Close()

	// Add all files in templates and static directories
	addFiles := func(dir string) {
		filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if !info.IsDir() {
				watcher.Add(path)
			}
			return nil
		})
	}

	addFiles("assets/beady/templates")
	addFiles("assets/beady/static")

	// Verify paths exist
	if _, err := os.Stat("assets/beady/templates"); os.IsNotExist(err) {
		log.Fatal("Development mode requires running from repository root (assets/beady/templates not found)")
	}

	for {
		select {
		case event, ok := <-watcher.Events:
			if !ok {
				return
			}
			if event.Has(fsnotify.Write) || event.Has(fsnotify.Create) {
				log.Printf("File changed: %s", event.Name)
				// Re-parse templates if a template file changed
				if strings.HasPrefix(event.Name, "assets/beady/templates/") && strings.HasSuffix(event.Name, ".html") {
					log.Printf("Re-parsing templates")
					parseTemplates()
				}
				log.Printf("Broadcasting reload to clients")
				broadcast("reload")
			}
		case err, ok := <-watcher.Errors:
			if !ok {
				return
			}
			log.Println("Watcher error:", err)
		}
	}
}

// main is the program entrypoint. It parses command-line flags, loads templates and the beads database (using the provided path or autodiscovery), configures HTTP routes and server timeouts, and starts the web UI server. In development mode it enables live-reload (file watcher and websocket), opens the default browser to the UI, and logs relevant startup info. The function blocks indefinitely.
func main() {
	flag.Usage = printUsage
	flag.Parse()
	if *help {
		flag.Usage()
		os.Exit(0)
	}

	// Set filesystem for templates and static files
	if devMode {
		tmplFS = os.DirFS("assets/beady")
	}
	parseTemplates()

	args := flag.Args()

	if len(args) > 2 {
		printUsage()
		os.Exit(1)
	}

	var dbPath string
	port := "8080"
	if len(args) > 0 {
		dbPath = args[0]
	}
	if len(args) > 1 {
		port = args[1]
	}

	// Open database
	var err error
	if dbPath == "" {
		// No path provided, try autodiscovery first
		if foundDB := beads.FindDatabasePath(); foundDB != "" {
			store, err = beads.NewSQLiteStorage(foundDB)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error opening database: %v\n", err)
				os.Exit(1)
			}
		} else {
			fmt.Fprintf(os.Stderr, "No database path provided and no database found via autodiscovery\n")
			os.Exit(1)
		}
	} else {
		// Path provided, try it first
		store, err = beads.NewSQLiteStorage(dbPath)
		if err != nil {
			// Try autodiscovery
			if foundDB := beads.FindDatabasePath(); foundDB != "" {
				store, err = beads.NewSQLiteStorage(foundDB)
				if err != nil {
					fmt.Fprintf(os.Stderr, "Error opening database: %v\n", err)
					os.Exit(1)
				}
			} else {
				fmt.Fprintf(os.Stderr, "Error opening database: %v\n", err)
				os.Exit(1)
			}
		}
	}

	addr := net.JoinHostPort("127.0.0.1", port)

	mux := http.NewServeMux()
	mux.HandleFunc("/", handleIndex)
	mux.HandleFunc("/ready", handleReady)
	mux.HandleFunc("/blocked", handleBlocked)
	mux.HandleFunc("/issue/", handleIssueDetail)
	mux.HandleFunc("/graph/", handleGraph)
	mux.HandleFunc("/api/issues", handleAPIIssues)
	mux.HandleFunc("/api/issue/", handleAPIIssue)
	mux.HandleFunc("/api/stats", handleAPIStats)
	if devMode {
		mux.HandleFunc("/ws", handleWS)
	}
	mux.HandleFunc("/static/", handleStatic)

	srv := &http.Server{
		Addr:         addr,
		Handler:      mux,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	fmt.Printf("Starting beads web UI at http://%s\n", addr)
	if devMode {
		fmt.Printf("Development mode enabled with live reload\n")
	}
	fmt.Printf("Press Ctrl+C to stop\n")

	if devMode {
		log.Printf("Starting file watcher for live reload")
		go startFileWatcher()
	}

	// Start server in goroutine
	errCh := make(chan error, 1)
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			errCh <- err
		}
	}()

	// Give server a moment to start
	time.Sleep(100 * time.Millisecond)
	select {
	case err := <-errCh:
		fmt.Fprintf(os.Stderr, "Error starting server: %v\n", err)
		os.Exit(1)
	default:
		// Server started successfully
	}

	if devMode {
		// Open browser (best-effort)
		url := "http://" + addr
		fmt.Printf("Opening browser to %s\n", url)
		if err := openBrowser(url); err != nil {
			log.Printf("Open browser failed: %v", err)
		}
	}

	// Wait for interrupt
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	<-sigCh

	// Graceful shutdown
	log.Println("Shutting down server...")
	if err := srv.Shutdown(context.Background()); err != nil {
		log.Printf("Server shutdown error: %v", err)
	}
	log.Println("Server stopped")
}

// handleIndex serves the main index page showing issues and statistics.
// It validates that the request path is "/" and the method is GET, then
// fetches up to 100 issues, enriches them with labels and dependency counts,
// obtains overall statistics, and renders the index template.
// Responds with 404 for non-root paths, 405 for non-GET methods, and 500 for
// storage or template rendering errors.
func handleIndex(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	ctx := r.Context()
	issues, err := store.SearchIssues(ctx, "", beads.IssueFilter{})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Limit to first 100 issues
	if len(issues) > 100 {
		issues = issues[:100]
	}

	issuesWithLabels := enrichIssuesWithLabels(ctx, issues)

	stats, err := store.GetStatistics(ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data := map[string]interface{}{
		"Issues": issuesWithLabels,
		"Stats":  stats,
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := tmplIndex.Execute(w, data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func handleIssueDetail(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	issueID := strings.TrimPrefix(r.URL.Path, "/issue/")
	if issueID == "" {
		http.Error(w, "Issue ID required", http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	issue, err := store.GetIssue(ctx, issueID)
	if err != nil {
		http.Error(w, "Issue not found", http.StatusNotFound)
		return
	}

	deps, _ := store.GetDependencies(ctx, issueID)
	dependents, _ := store.GetDependents(ctx, issueID)
	labels, _ := store.GetLabels(ctx, issueID)
	events, _ := store.GetEvents(ctx, issueID, 50)

	data := map[string]interface{}{
		"Issue":      issue,
		"Deps":       deps,
		"Dependents": dependents,
		"Labels":     labels,
		"Events":     events,
		"HasDeps":    len(deps) > 0 || len(dependents) > 0,
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := tmplDetail.Execute(w, data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func handleGraph(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	issueID := strings.TrimPrefix(r.URL.Path, "/graph/")
	if issueID == "" {
		http.Error(w, "Issue ID required", http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	issue, err := store.GetIssue(ctx, issueID)
	if err != nil {
		http.Error(w, "Issue not found", http.StatusNotFound)
		return
	}

	dotGraph := generateDotGraph(ctx, issue)

	data := map[string]interface{}{
		"Issue":    issue,
		"DotGraph": dotGraph,
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := tmplGraph.Execute(w, data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func handleReady(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	ctx := r.Context()

	ready, err := store.GetReadyWork(ctx, beads.WorkFilter{})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Filter out issues with excluded labels
	excludeLabel := r.URL.Query().Get("exclude")
	var filtered []*beads.Issue
	if excludeLabel != "" {
		for _, issue := range ready {
			labels, _ := store.GetLabels(ctx, issue.ID)
			hasExcluded := false
			for _, label := range labels {
				if label == excludeLabel {
					hasExcluded = true
					break
				}
			}
			if !hasExcluded {
				filtered = append(filtered, issue)
			}
		}
	} else {
		filtered = ready
	}

	issuesWithLabels := enrichIssuesWithLabels(ctx, filtered)
	stats, _ := store.GetStatistics(ctx)

	data := map[string]interface{}{
		"Issues":       issuesWithLabels,
		"Stats":        stats,
		"ExcludeLabel": excludeLabel,
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := tmplReady.Execute(w, data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func handleBlocked(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	ctx := r.Context()

	blocked, err := store.GetBlockedIssues(ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	stats, _ := store.GetStatistics(ctx)

	data := map[string]interface{}{
		"Blocked": blocked,
		"Stats":   stats,
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := tmplBlocked.Execute(w, data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func handleAPIIssues(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	ctx := r.Context()

	searchQuery := ""

	// Create a filter with default limit of 1000
	filter := beads.IssueFilter{}

	// We'll handle filtering manually since we can't set the limit directly
	if status := r.URL.Query().Get("status"); status != "" {
		s := beads.Status(status)
		filter.Status = &s
	}

	if priority := r.URL.Query().Get("priority"); priority != "" {
		p, err := strconv.Atoi(priority)
		if err != nil {
			http.Error(w, "Invalid priority", http.StatusBadRequest)
			return
		}
		filter.Priority = &p
	}

	issues, err := store.SearchIssues(ctx, searchQuery, filter)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Apply limit manually
	if len(issues) > 1000 {
		issues = issues[:1000]
	}

	// Check if htmx request (return partial HTML)
	if r.Header.Get("HX-Request") == "true" {
		issuesWithLabels := enrichIssuesWithLabels(ctx, issues)
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		if err := tmplIssuesTbody.Execute(w, issuesWithLabels); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	// Regular JSON response
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(issues); err != nil {
		http.Error(w, "Failed to encode JSON", http.StatusInternalServerError)
	}
}

func handleAPIIssue(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	issueID := strings.TrimPrefix(r.URL.Path, "/api/issue/")

	ctx := r.Context()
	issue, err := store.GetIssue(ctx, issueID)
	if err != nil {
		http.Error(w, "Issue not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(issue); err != nil {
		http.Error(w, "Failed to encode JSON", http.StatusInternalServerError)
	}
}

func handleAPIStats(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	ctx := r.Context()
	stats, err := store.GetStatistics(ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(stats); err != nil {
		http.Error(w, "Failed to encode JSON", http.StatusInternalServerError)
	}
}

type IssueWithLabels struct {
	*beads.Issue
	Labels        []string
	DepsCount     int
	BlockersCount int
}

func enrichIssuesWithLabels(ctx context.Context, issues []*beads.Issue) []*IssueWithLabels {
	result := make([]*IssueWithLabels, len(issues))
	for i, issue := range issues {
		labels, _ := store.GetLabels(ctx, issue.ID)
		deps, _ := store.GetDependencies(ctx, issue.ID)
		dependents, _ := store.GetDependents(ctx, issue.ID)
		result[i] = &IssueWithLabels{
			Issue:         issue,
			Labels:        labels,
			DepsCount:     len(deps),
			BlockersCount: len(dependents),
		}
	}
	return result
}

// generateDotGraph builds a DOT-format directed graph for the given root issue,
// including the root's dependencies and dependents as nodes and edges.
// The returned string is a complete DOT graph where each node is styled and
// colored according to the issue's status and contains the issue ID, title, and priority.
func generateDotGraph(ctx context.Context, root *beads.Issue) string {
	var sb strings.Builder
	sb.WriteString("digraph G {\n")
	sb.WriteString("  rankdir=TB;\n")
	sb.WriteString("  node [shape=box, style=filled];\n\n")

	// Build node and edge maps to avoid duplicates
	nodes := make(map[string]*beads.Issue)
	edges := make(map[string]bool)

	// Add root
	nodes[root.ID] = root

	// Get dependencies and dependents to build relationships
	deps, _ := store.GetDependencies(ctx, root.ID)
	dependents, _ := store.GetDependents(ctx, root.ID)

	// Add all dependencies as nodes and edges
	for _, dep := range deps {
		nodes[dep.ID] = dep
		edgeKey := fmt.Sprintf("%s->%s", root.ID, dep.ID)
		edges[edgeKey] = true
	}

	// Add all dependents as nodes and edges
	for _, dependent := range dependents {
		nodes[dependent.ID] = dependent
		edgeKey := fmt.Sprintf("%s->%s", dependent.ID, root.ID)
		edges[edgeKey] = true
	}

	// Render all nodes
	for _, issue := range nodes {
		color := "#7b9e87" // open
		if issue.Status == beads.StatusClosed {
			color = "#8a8175"
		} else if issue.Status == beads.StatusInProgress {
			color = "#c17a3c"
		}

		// Escape title for DOT format
		title := strings.ReplaceAll(issue.Title, "\\", "\\\\")
		title = strings.ReplaceAll(title, "\"", "'")

		label := fmt.Sprintf("%s\\n%s\\nP%d", issue.ID, title, issue.Priority)

		sb.WriteString(fmt.Sprintf("  \"%s\" [label=\"%s\", fillcolor=\"%s\", fontcolor=\"white\"];\n",
			issue.ID, label, color))
	}

	sb.WriteString("\n")

	// Render all edges
	for edge := range edges {
		parts := strings.Split(edge, "->")
		sb.WriteString(fmt.Sprintf("  \"%s\" -> \"%s\";\n", parts[0], parts[1]))
	}

	sb.WriteString("}\n")
	return sb.String()
}

// openBrowser opens the specified URL in the user's default web browser.
// It returns an error if the platform command used to launch the browser cannot be started.
func openBrowser(url string) error {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "windows":
		// Empty title arg avoids treating URL as window title; quote-safe
		cmd = exec.Command("cmd", "/c", "start", "", url)
	case "darwin":
		cmd = exec.Command("open", url)
	default: // linux, etc.
		cmd = exec.Command("xdg-open", url)
	}
	return cmd.Start()
}

// handleStatic serves files requested under the /static/ path from the configured template filesystem.
// It looks up the resource under "static/{path}" with a fallback to "templates/{path}", sets the
// Content-Type for ".css" and ".js" files, responds with 404 if the file cannot be found, and returns
// 405 Method Not Allowed for any non-GET request.
func handleStatic(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	path := strings.TrimPrefix(r.URL.Path, "/static/")

	var contentType string
	if strings.HasSuffix(path, ".css") {
		contentType = "text/css; charset=utf-8"
	} else if strings.HasSuffix(path, ".js") {
		contentType = "application/javascript; charset=utf-8"
	}

	content, err := fs.ReadFile(tmplFS, "static/"+path)
	if err != nil {
		// Try templates directory as fallback (for backward compatibility)
		content, err = fs.ReadFile(tmplFS, "templates/"+path)
		if err != nil {
			http.NotFound(w, r)
			return
		}
	}

	w.Header().Set("Content-Type", contentType)
	w.Write(content)
}