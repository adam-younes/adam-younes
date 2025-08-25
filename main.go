package main

import (
	"bytes"
	"embed"
	"html/template"
	"io/fs"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"github.com/adam-younes/adam-younes/models"
	"github.com/gomarkdown/markdown"
)

//go:embed templates/base.html
//go:embed templates/elements/*
//go:embed templates/pages/*
//go:embed static/css/*
//go:embed static/assets/*
//go:embed static/fonts/*
var assets embed.FS

var tmpl *template.Template

func init() {
	// FS rooted at ./templates/
	tplFS, err := fs.Sub(assets, "templates")
	if err != nil {
		log.Fatal(err)
	}

	// Create template with custom functions
	tmpl = template.Must(template.New("").Funcs(template.FuncMap{
		"getWebPath": getWebPath,
		"getRelativePath": getRelativePath,
	}).ParseFS(
		tplFS,
		"base.html",
		"elements/*.html",
		"pages/*.html",
	))
}

type PageData struct {
	Title   string
	Content template.HTML
}

func render(w http.ResponseWriter, contentTmpl, title string, data interface{}) {
	var buf bytes.Buffer
	if err := tmpl.ExecuteTemplate(&buf, contentTmpl, data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	pageData := PageData{
		Title:   title,
		Content: template.HTML(buf.String()),
	}
	if err := tmpl.ExecuteTemplate(w, "base", pageData); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

type PageHandler struct {
	Template string
	Title    string
}

func (p PageHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	render(w, p.Template, p.Title, nil)
}

// NoteItem represents a file or directory in the notes structure
type NoteItem struct {
	Name     string
	Path     string
	IsDir    bool
	Children []NoteItem
}

// NoteData represents the data for a specific note
type NoteData struct {
	Title   string
	Content template.HTML
	Path    string
}

// NotesPageData represents the data for the notes page
type NotesPageData struct {
	Items      []NoteItem
	CurrentNote *NoteData
	SearchQuery string
}

// scanNotesDirectory recursively scans the notes directory and builds the structure
func scanNotesDirectory(rootPath string) ([]NoteItem, error) {
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
			children, err := scanNotesDirectory(item.Path)
			if err != nil {
				log.Printf("Error scanning directory %s: %v", item.Path, err)
				continue
			}
			item.Children = children
		}
		
		items = append(items, item)
	}
	
	return items, nil
}

// getRelativePath returns the path relative to the notes directory
func getRelativePath(fullPath string) string {
	notesDir := "static/notes"
	if strings.HasPrefix(fullPath, notesDir) {
		return strings.TrimPrefix(fullPath, notesDir)
	}
	return fullPath
}

// getWebPath converts a file system path to a web URL path
func getWebPath(fsPath string) string {
	relativePath := getRelativePath(fsPath)
	if relativePath == "" {
		return "/notes"
	}
	return "/notes" + relativePath
}

// readMarkdownFile reads and renders a markdown file
func readMarkdownFile(filePath string) (*NoteData, error) {
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
		Path:    getWebPath(filePath),
	}, nil
}

func main() {
	mux := http.NewServeMux()
	staticFS, err := fs.Sub(assets, "static")
	if err != nil {
		log.Fatal(err)
	}
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.FS(staticFS))))

	mux.Handle("/", PageHandler{"indexContent", "Home"})

	mux.HandleFunc("/notes", notesHandler)
	mux.HandleFunc("/notes/", notesHandler)

	mux.Handle("/about", PageHandler{"aboutContent", "About Me"})

	mux.HandleFunc("/projects", listProjects)
	mux.HandleFunc("/experience", listExperience)

	// redirects
	mux.Handle("/experience/", 	http.StripPrefix("/experience/", mux))
	mux.Handle("/projects/", 	http.StripPrefix("/projects/", mux))
	mux.Handle("/about/", 		http.StripPrefix("/about/", mux))

	log.Println("Listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}

type ProjectsPageData struct {
	Projects []models.Project
}

func listProjects(w http.ResponseWriter, r *http.Request) {
	projects := []models.Project{
		{
			Title:       "Calculator",
			Description: "A terminal based graphing calculator",
			Link:        "https://github.com/adam-younes/calculator",
		},
		{
			Title:       "atom",
			Description: "A Linux based Wayland game engine",
			Link:        "https://github.com/adam-younes/atom",
		},
	}

	// execute the content template with our ProjectsPageData as "dot"
	var buf bytes.Buffer
	if err := tmpl.ExecuteTemplate(&buf, "projectsContent", ProjectsPageData{projects}); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	// now inject that HTML into your base layout
	data := PageData{
		Title:   "My Projects",
		Content: template.HTML(buf.String()),
	}
	if err := tmpl.ExecuteTemplate(w, "base", data); err != nil {
		http.Error(w, err.Error(), 500)
	}
}

type ExperiencePageData struct {
	Experience []models.Experience
}

func listExperience(w http.ResponseWriter, r *http.Request) {
	experience := []models.Experience {
		{
			Title:       "Runtime Software Engineer",
			Description: "Contributed to real-time flight simulation software",
			Link:        "https://github.com/you/graphviz",
		},
	}

	// execute the content template with our ExperienceData as "dot"
	var buf bytes.Buffer
	if err := tmpl.ExecuteTemplate(&buf, "experienceContent", ExperiencePageData{experience}); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	// now inject that HTML into your base layout
	data := PageData{
		Title:   "Experience",
		Content: template.HTML(buf.String()),
	}

	if err := tmpl.ExecuteTemplate(w, "base", data); err != nil {
		http.Error(w, err.Error(), 500)
	}
}

func notesHandler(w http.ResponseWriter, r *http.Request) {
	// Extract the path after /notes/
	requestPath := strings.TrimPrefix(r.URL.Path, "/notes")
	
	// Handle search query
	searchQuery := r.URL.Query().Get("search")
	
	// Scan the notes directory
	notesDir := "static/notes"
	items, err := scanNotesDirectory(notesDir)
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
			noteData, err := readMarkdownFile(filePath)
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
		items = filterNotesBySearch(items, searchQuery)
	}
	
	data := NotesPageData{
		Items:       items,
		CurrentNote: currentNote,
		SearchQuery: searchQuery,
	}
	
	// Render the notes page
	render(w, "notesContent", "Notes", data)
}

// filterNotesBySearch recursively filters notes based on search query
func filterNotesBySearch(items []NoteItem, query string) []NoteItem {
	var filtered []NoteItem
	query = strings.ToLower(query)
	
	for _, item := range items {
		// Check if current item matches
		matches := strings.Contains(strings.ToLower(item.Name), query)
		
		// If it's a directory, check children
		if item.IsDir {
			filteredChildren := filterNotesBySearch(item.Children, query)
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
