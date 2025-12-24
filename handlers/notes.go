package handlers

import (
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

// NotesHandler handles the notes page
type NotesHandler struct {
	renderer *Renderer
}

// NewNotesHandler creates a new notes handler
func NewNotesHandler(renderer *Renderer) *NotesHandler {
	return &NotesHandler{renderer: renderer}
}

// ServeHTTP implements the http.Handler interface
func (h *NotesHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Extract the path after /notes/
	requestPath := strings.TrimPrefix(r.URL.Path, "/notes")

	// Handle search query
	searchQuery := r.URL.Query().Get("search")

	// Scan the notes directory
	notesDir := "static/notes"
	items, err := ScanNotesDirectory(notesDir)
	if err != nil {
		log.Printf("Error scanning notes directory: %v", err)
		http.Error(w, "Error reading notes", http.StatusInternalServerError)
		return
	}

	var currentNote *NoteData

	// If a specific file is requested, read and render it
	if requestPath != "" && requestPath != "/" {
		filePath := filepath.Join(notesDir, requestPath)

		// Check if file exists and is a markdown file
		if info, err := os.Stat(filePath); err == nil && !info.IsDir() && strings.HasSuffix(filePath, ".md") {
			noteData, err := ReadMarkdownFile(filePath)
			if err != nil {
				log.Printf("Error reading markdown file: %v", err)
				http.Error(w, "Error reading note", http.StatusInternalServerError)
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

	// Render the notes page
	h.renderer.Render(w, "notesContent", "Notes", data)
}

