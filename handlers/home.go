package handlers

import "net/http"

// HomeHandler handles the home page
type HomeHandler struct {
	renderer *Renderer
}

// NewHomeHandler creates a new home handler
func NewHomeHandler(renderer *Renderer) *HomeHandler {
	return &HomeHandler{renderer: renderer}
}

// ServeHTTP implements the http.Handler interface
func (h *HomeHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.renderer.RenderWithNavbar(w, "indexContent",
		"Adam Younes",
		"Adam Younes — Palestinian software engineer. Founding engineer at Mod AI (YC F25). Projects, experience, and notes.",
		"/",
		nil, false)
}

