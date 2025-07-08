package main

import (
	"bytes"
	"embed"
	"html/template"
	"io/fs"
	"log"
	"net/http"
	"github.com/adam-younes/adam-younes/models"
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

	tmpl = template.Must(template.ParseFS(
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

func render(w http.ResponseWriter, contentTmpl, title string) {
	var buf bytes.Buffer
	if err := tmpl.ExecuteTemplate(&buf, contentTmpl, nil); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data := PageData{
		Title:   title,
		Content: template.HTML(buf.String()),
	}
	if err := tmpl.ExecuteTemplate(w, "base", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

type PageHandler struct {
	Template string
	Title    string
}

func (p PageHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	render(w, p.Template, p.Title)
}

func main() {
	mux := http.NewServeMux()
	staticFS, err := fs.Sub(assets, "static")
	if err != nil {
		log.Fatal(err)
	}
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.FS(staticFS))))

	mux.Handle("/", PageHandler{"indexContent", "Home"})

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

	// execute the content template with our ProjectsPageData as “dot”
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

	// execute the content template with our ExperienceData as “dot”
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

