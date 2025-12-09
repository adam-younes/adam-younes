package handlers

import (
	"bytes"
	"html/template"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gomarkdown/markdown"
)

// Renderer handles template rendering
type Renderer struct {
	tmpl *template.Template
}

// NewRenderer creates a new renderer with the given template
func NewRenderer(tmpl *template.Template) *Renderer {
	return &Renderer{tmpl: tmpl}
}

// Render renders a page using the template
func (r *Renderer) Render(w http.ResponseWriter, contentTmpl, title string, data interface{}) {
	var buf bytes.Buffer
	if err := r.tmpl.ExecuteTemplate(&buf, contentTmpl, data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	pageData := PageData{
		Title:   title,
		Content: template.HTML(buf.String()),
	}
	if err := r.tmpl.ExecuteTemplate(w, "base", pageData); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// ScanNotesDirectory recursively scans the notes directory and builds the structure
func ScanNotesDirectory(rootPath string) ([]NoteItem, error) {
	var items []NoteItem

	entries, err := os.ReadDir(rootPath)
	if err != nil {
		return nil, err
	}

	for _, entry := range entries {
		item := NoteItem{
			Name:  entry.Name(),
			Path:  filepath.Join(rootPath, entry.Name()),
			IsDir: entry.IsDir(),
		}

		if entry.IsDir() {
			children, err := ScanNotesDirectory(item.Path)
			if err != nil {
				continue
			}
			item.Children = children
		}

		items = append(items, item)
	}

	return items, nil
}

// GetRelativePath returns the path relative to the notes directory
func GetRelativePath(fullPath string) string {
	notesDir := "static/notes"
	if strings.HasPrefix(fullPath, notesDir) {
		return strings.TrimPrefix(fullPath, notesDir)
	}
	return fullPath
}

// GetWebPath converts a file system path to a web URL path
func GetWebPath(fsPath string) string {
	relativePath := GetRelativePath(fsPath)
	if relativePath == "" {
		return "/notes"
	}
	return "/notes" + relativePath
}

// ReadMarkdownFile reads and renders a markdown file
func ReadMarkdownFile(filePath string) (*NoteData, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	// Convert markdown to HTML
	html := markdown.ToHTML(content, nil, nil)

	// Extract title from first heading or use filename
	title := strings.TrimSuffix(filepath.Base(filePath), filepath.Ext(filePath))

	return &NoteData{
		Title:   title,
		Content: template.HTML(html),
		Path:    GetWebPath(filePath),
	}, nil
}

// FilterNotesBySearch recursively filters notes based on search query
func FilterNotesBySearch(items []NoteItem, query string) []NoteItem {
	var filtered []NoteItem
	query = strings.ToLower(query)

	for _, item := range items {
		// Check if current item matches
		matches := strings.Contains(strings.ToLower(item.Name), query)

		// If it's a directory, check children
		if item.IsDir {
			filteredChildren := FilterNotesBySearch(item.Children, query)
			if len(filteredChildren) > 0 || matches {
				item.Children = filteredChildren
				filtered = append(filtered, item)
			}
		} else if matches {
			filtered = append(filtered, item)
		}
	}

	return filtered
}

