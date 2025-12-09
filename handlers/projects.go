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

	data := ProjectsPageData{Projects: projects}
	h.renderer.Render(w, "projectsContent", "My Projects", data)
}
