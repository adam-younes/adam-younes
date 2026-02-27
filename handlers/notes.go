package handlers

import (
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

// NotesHandler handles the notes page
type NotesHandler struct {
	renderer  *Renderer
	mu        sync.RWMutex
	cache     []NoteItem
	cacheTime time.Time
	cacheTTL  time.Duration
}

// NewNotesHandler creates a new notes handler
func NewNotesHandler(renderer *Renderer) *NotesHandler {
	return &NotesHandler{
		renderer: renderer,
		cacheTTL: 5 * time.Minute,
	}
}

// getItems returns the cached notes tree, refreshing if expired.
func (h *NotesHandler) getItems(notesDir string) ([]NoteItem, error) {
	h.mu.RLock()
	if h.cache != nil && time.Since(h.cacheTime) < h.cacheTTL {
		items := h.cache
		h.mu.RUnlock()
		return items, nil
	}
	h.mu.RUnlock()

	h.mu.Lock()
	defer h.mu.Unlock()

	// Double-check after acquiring write lock
	if h.cache != nil && time.Since(h.cacheTime) < h.cacheTTL {
		return h.cache, nil
	}

	items, err := ScanNotesDirectory(notesDir)
	if err != nil {
		return nil, err
	}
	h.cache = items
	h.cacheTime = time.Now()
	return items, nil
}

// ServeHTTP implements the http.Handler interface
func (h *NotesHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Extract the path after /notes/
	requestPath := strings.TrimPrefix(r.URL.Path, "/notes")

	// Handle search query
	searchQuery := r.URL.Query().Get("search")

	// Scan the notes directory (cached)
	notesDir := "static/notes"
	items, err := h.getItems(notesDir)
	if err != nil {
		log.Printf("Error scanning notes directory: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	var currentNote *NoteData

	// If a specific file is requested, read and render it
	if requestPath != "" && requestPath != "/" {
		filePath := filepath.Join(notesDir, requestPath)

		// Prevent path traversal: ensure resolved path stays within notesDir
		absNotes, err := filepath.Abs(notesDir)
		if err != nil {
			log.Printf("Error resolving notes directory: %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		absFile, err := filepath.Abs(filePath)
		if err != nil {
			http.NotFound(w, r)
			return
		}
		if !strings.HasPrefix(absFile, absNotes+string(os.PathSeparator)) {
			http.NotFound(w, r)
			return
		}

		// Check if file exists and is a markdown file
		if info, err := os.Stat(filePath); err == nil && !info.IsDir() && strings.HasSuffix(filePath, ".md") {
			noteData, err := ReadMarkdownFile(filePath)
			if err != nil {
				log.Printf("Error reading markdown file: %v", err)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}
			currentNote = noteData
		}
	}

	// Filter items based on search query if provided
	if searchQuery != "" {
		items = FilterNotesBySearch(items, searchQuery)
	}

	data := NotesPageData{
		Items:       items,
		CurrentNote: currentNote,
		SearchQuery: searchQuery,
	}

	// Build SEO fields based on whether a specific note is being viewed
	title := "Notes - Adam Younes"
	description := "Technical notes and documentation by Adam Younes."
	path := "/notes"
	if currentNote != nil {
		title = currentNote.Title + " - Adam Younes"
		description = currentNote.Title + " — technical notes by Adam Younes."
		path = currentNote.Path
	}

	// Render the notes page
	h.renderer.Render(w, "notesContent", title, description, path, data)
}

