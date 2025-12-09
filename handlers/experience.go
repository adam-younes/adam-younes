package handlers

import "net/http"

// ExperienceHandler handles the experience page
type ExperienceHandler struct {
	renderer *Renderer
}

// NewExperienceHandler creates a new experience handler
func NewExperienceHandler(renderer *Renderer) *ExperienceHandler {
	return &ExperienceHandler{renderer: renderer}
}

// ServeHTTP implements the http.Handler interface
func (h *ExperienceHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	experience := []Experience{
		{
			Title:       "Runtime Software Engineer",
			Description: "Contributed to real-time flight simulation software",
			Link:        "https://github.com/you/graphviz",
		},
	}

	data := ExperiencePageData{Experience: experience}
	h.renderer.Render(w, "experienceContent", "Experience", data)
}
