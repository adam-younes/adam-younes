package main

import (
	"embed"
	"html/template"
	"io/fs"
	"log"
	"net/http"

	"github.com/adam-younes/adam-younes/handlers"
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

	// Create template with custom functions
	tmpl = template.Must(template.New("").Funcs(template.FuncMap{
		"getWebPath": handlers.GetWebPath,
		"getRelativePath": handlers.GetRelativePath,
	}).ParseFS(
		tplFS,
		"base.html",
		"elements/*.html",
		"pages/*.html",
	))
}



func main() {
	mux := http.NewServeMux()
	staticFS, err := fs.Sub(assets, "static")
	if err != nil {
		log.Fatal(err)
	}
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.FS(staticFS))))

	// Create renderer
	renderer := handlers.NewRenderer(tmpl)

	// Create handlers
	homeHandler := handlers.NewHomeHandler(renderer)
	aboutHandler := handlers.NewAboutHandler(renderer)
	projectsHandler := handlers.NewProjectsHandler(renderer)
	experienceHandler := handlers.NewExperienceHandler(renderer)
	notesHandler := handlers.NewNotesHandler(renderer)

	// Register routes
	mux.Handle("/", homeHandler)
	mux.Handle("/about", aboutHandler)
	mux.Handle("/projects", projectsHandler)
	mux.Handle("/experience", experienceHandler)
	mux.HandleFunc("/notes", notesHandler.ServeHTTP)
	mux.HandleFunc("/notes/", notesHandler.ServeHTTP)

	// redirects
	mux.Handle("/experience/", http.StripPrefix("/experience/", mux))
	mux.Handle("/projects/", http.StripPrefix("/projects/", mux))
	mux.Handle("/about/", http.StripPrefix("/about/", mux))

	log.Println("Listening on :9999")
	log.Fatal(http.ListenAndServe(":9999", mux))
}


