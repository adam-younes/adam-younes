package handlers

import (
	"bytes"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
	"github.com/microcosm-cc/bluemonday"
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
func (r *Renderer) Render(w http.ResponseWriter, contentTmpl, title, description, path string, data interface{}) {
	r.RenderWithNavbar(w, contentTmpl, title, description, path, data, true)
}

// RenderWithNavbar renders a page with optional navbar
func (r *Renderer) RenderWithNavbar(w http.ResponseWriter, contentTmpl, title, description, path string, data interface{}, showNavbar bool) {
	var buf bytes.Buffer
	if err := r.tmpl.ExecuteTemplate(&buf, contentTmpl, data); err != nil {
		log.Printf("template render error: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	pageData := PageData{
		Title:       title,
		Description: description,
		Path:        path,
		Content:     template.HTML(buf.String()),
		ShowNavbar:  showNavbar,
	}
	if err := r.tmpl.ExecuteTemplate(w, "base", pageData); err != nil {
		log.Printf("base template error: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
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

// ExtractTOC extracts table of contents from markdown content
func ExtractTOC(content []byte) []TOCItem {
	var items []TOCItem

	contentStr := string(content)

	// Find the "## Table of Contents" section
	tocHeaderPattern := regexp.MustCompile(`(?mi)^##\s+Table of Contents\s*$`)
	tocHeaderLoc := tocHeaderPattern.FindStringIndex(contentStr)
	if tocHeaderLoc == nil {
		return items
	}

	// Get content after the TOC header
	afterHeader := contentStr[tocHeaderLoc[1]:]

	// Find where the TOC section ends (next heading or ---)
	endPattern := regexp.MustCompile(`(?m)^(##|---)`)
	endLoc := endPattern.FindStringIndex(afterHeader)

	var tocSection string
	if endLoc != nil {
		tocSection = afterHeader[:endLoc[0]]
	} else {
		tocSection = afterHeader
	}

	// Extract links from numbered list items: "1. [Title](#anchor)"
	linkPattern := regexp.MustCompile(`\[([^\]]+)\]\(([^)]+)\)`)
	matches := linkPattern.FindAllStringSubmatch(tocSection, -1)

	for _, match := range matches {
		if len(match) >= 3 {
			items = append(items, TOCItem{
				Title:  match[1],
				Anchor: match[2],
			})
		}
	}

	return items
}

// ReadMarkdownFile reads and renders a markdown file
func ReadMarkdownFile(filePath string) (*NoteData, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	// Extract TOC before rendering
	toc := ExtractTOC(content)

	// Create parser with extensions for auto heading IDs
	extensions := parser.CommonExtensions | parser.AutoHeadingIDs
	p := parser.NewWithExtensions(extensions)

	// Create renderer with heading IDs enabled
	htmlFlags := html.CommonFlags
	opts := html.RendererOptions{Flags: htmlFlags}
	renderer := html.NewRenderer(opts)

	// Convert markdown to HTML and sanitize
	htmlContent := markdown.ToHTML(content, p, renderer)
	sanitizer := bluemonday.UGCPolicy()
	htmlContent = sanitizer.SanitizeBytes(htmlContent)

	// Extract title from first heading or use filename
	title := strings.TrimSuffix(filepath.Base(filePath), filepath.Ext(filePath))

	return &NoteData{
		Title:   title,
		Content: template.HTML(htmlContent),
		Path:    GetWebPath(filePath),
		TOC:     toc,
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

