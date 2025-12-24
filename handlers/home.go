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
	h.renderer.Render(w, "indexContent", "Home", nil)
}

