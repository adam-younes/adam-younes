package handlers

import "net/http"

// ProjectsHandler handles the projects page
type ProjectsHandler struct {
	renderer *Renderer
}

// NewProjectsHandler creates a new projects handler
func NewProjectsHandler(renderer *Renderer) *ProjectsHandler {
	return &ProjectsHandler{renderer: renderer}
}

// ServeHTTP implements the http.Handler interface
func (h *ProjectsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	projects := []Project{
		{
			Title:       "calculator",
			Description: "I built a graphing calculator in C from scratch for a scholarship competition in high school",
			Link:        "https://github.com/adam-younes/calculator",
		},
		{
			Title:       "atom",
			Description: "A Linux based Wayland game engine",
			Link:        "https://github.com/adam-younes/atom",
		},
		{
			Title:       "last-fare",
			Description: "A thriller driving simulator about a rideshare driver completing their shift night after night, inspired by Papers, Please. Built with Godot 4.5",
			Link:        "https://github.com/adam-younes/last-fare",
		},
	}

	data := ProjectsPageData{Projects: projects}
	h.renderer.Render(w, "projectsContent", "My Projects", data)
}
