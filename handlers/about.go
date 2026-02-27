package handlers

import "net/http"

// AboutHandler handles the about page
type AboutHandler struct {
	renderer *Renderer
}

// NewAboutHandler creates a new about handler
func NewAboutHandler(renderer *Renderer) *AboutHandler {
	return &AboutHandler{renderer: renderer}
}

// ServeHTTP implements the http.Handler interface
func (h *AboutHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.renderer.Render(w, "aboutContent",
		"About - Adam Younes",
		"About Adam Younes, a Palestinian software engineer. Founding engineer at Mod AI (YC F25), previously at Nixie Tech and Aechelon Technology.",
		"/about",
		nil)
}

