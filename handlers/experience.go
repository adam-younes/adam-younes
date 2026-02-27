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
			Title:       "Founding Software Engineer",
			Company:     "Mod AI (YC F25)",
			DateRange:   "August 2025 - Present",
			Description: "Building AI agents that automate accounts payable workflows, enabling finance teams to process significantly more invoices without adding headcount.",
			Link:        "https://www.usemod.ai/",
		},
		{
			Title:       "Software Engineer",
			Company:     "Nixie Tech",
			DateRange:   "March - August 2025",
			Description: "Built digital signage management software for a contracting company.",
			Link:        "https://nixietech.com/",
		},
		{
			Title:       "Runtime Software Engineer",
			Company:     "Aechelon Technology",
			DateRange:   "May - August 2024",
			Description: "Contributed to real-time flight simulation software.",
			Link:        "https://www.aechelon.com/",
		},
	}

	data := ExperiencePageData{Experience: experience}
	h.renderer.Render(w, "experienceContent",
		"Experience - Adam Younes",
		"Professional experience of Adam Younes — Founding Software Engineer at Mod AI (YC F25), Nixie Tech, and Aechelon Technology.",
		"/experience",
		data)
}
