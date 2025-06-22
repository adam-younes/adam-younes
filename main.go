package main

import (
	"bytes"
	"embed"
	"html/template"
	"io/fs"
	"log"
	"net/http"
)

//go:embed templates/base.html
//go:embed templates/elements/*
//go:embed templates/pages/*
//go:embed static/css/*
//go:embed static/assets/*
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
	mux.Handle("/about/", PageHandler{"aboutContent", "About Me"})

	log.Println("Listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}

