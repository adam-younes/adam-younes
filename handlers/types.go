package handlers

import "html/template"

// PageData represents the data structure for rendering pages
type PageData struct {
	Title   string
	Content template.HTML
}

// PageHandler is a generic handler for simple pages
type PageHandler struct {
	Template string
	Title    string
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
	Items       []NoteItem
	CurrentNote *NoteData
	SearchQuery string
}

// ProjectsPageData represents the data for the projects page
type ProjectsPageData struct {
	Projects []Project
}

// ExperiencePageData represents the data for the experience page
type ExperiencePageData struct {
	Experience []Experience
}

// Project represents a project item
type Project struct {
	Title       string
	Description string
	Link        string
}

// Experience represents an experience item
type Experience struct {
	Title       string
	Description string
	Link        string
}

