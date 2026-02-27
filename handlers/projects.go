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
			Link:        "/notes/games/ideas/last-fare.md",
		},
		{
			Title:       "qbtimer",
			Description: "A competition-grade Rubik's Cube timer with automatic scramble generation for 2x2-5x5, solve tracking, and user accounts. Built with React, Express, Go, and PostgreSQL",
			Link:        "https://qbtimer.com",
		},
	}

	data := ProjectsPageData{Projects: projects}
	h.renderer.Render(w, "projectsContent",
		"Projects - Adam Younes",
		"Software projects by Adam Younes — graphing calculator in C, Wayland game engine, driving simulator, and Rubik's Cube timer.",
		"/projects",
		data)
}
