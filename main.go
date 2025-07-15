package main

import (
	"bytes"
	"embed"
	"html/template"
	"io/fs"
	"log"
	"net/http"
	"strings"
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

	mux.HandleFunc("/notes", notesHandler)
	mux.HandleFunc("/notes/", notesHandler)

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

type NoteSection struct {
   ID      string
   Title   string
   Content template.HTML
}
type NoteChapter struct {
   Title    string
   Sections []NoteSection
}
type BookData struct {
   Slug     string
   Title    string
   Chapters []NoteChapter
}
type Category struct {
   Slug  string
   Title string
}
type NotesPageData struct {
   Categories  []Category
   CurrentSlug string       // empty on /notes
   Book        *BookData    // nil on /notes
}

var allCategories = []Category{
  {Slug: "software",     Title: "Software"},
  {Slug: "math",         Title: "Math"},
  {Slug: "science",      Title: "Science"},
  {Slug: "example-book", Title: "The Go Primer"},  // ← add this
}
// in main.go, replace or append to your allBooks:
var allBooks = []BookData{
  {
    Slug:  "example-book",
    Title: "The Go Primer",
    Chapters: []NoteChapter{
      {
        Title: "Chapter 1: Getting Started",
        Sections: []NoteSection{
          {
            ID:      "example-intro",
            Title:   "Introduction",
            Content: `<p>Go (or Golang) is an open-source, statically typed language designed at Google. It's great for building fast, concurrent services and command-line tools.</p>`,
          },
          {
            ID:      "example-install",
            Title:   "Installation",
            Content: `<p>Download the binary from <a href="https://golang.org/dl/">golang.org/dl</a> and follow the installer for your OS. After installing, verify with <code>go version</code>.</p>`,
          },
        },
      },
      {
        Title: "Chapter 2: Core Concepts",
        Sections: []NoteSection{
          {
            ID:      "example-variables",
            Title:   "Variables & Types",
            Content: `<p>Declare variables with <code>var x int</code> or use the shorthand <code>x := 42</code>. Go has built-in types like <code>int</code>, <code>string</code>, <code>bool</code>, plus slices, maps, structs, and more.</p>`,
          },
          {
            ID:      "example-functions",
            Title:   "Functions",
            Content: `<p>Functions are first-class. Declare with <code>func add(a, b int) int { return a + b }</code>. Multiple return values are supported: <code>func swap(a, b string) (string, string)</code>.</p>`,
          },
        },
      },
    },
  },
}


func notesHandler(w http.ResponseWriter, r *http.Request) {
   slug := strings.TrimPrefix(r.URL.Path, "/notes/")
   var book *BookData
   for i := range allBooks {
     if allBooks[i].Slug == slug {
       book = &allBooks[i]
       break
     }
   }
   data := NotesPageData{
     Categories:  allCategories,
     CurrentSlug: slug,
     Book:        book,
   }

   // render into notesContent
   var buf bytes.Buffer
   if err := tmpl.ExecuteTemplate(&buf, "notesContent", data); err != nil {
     http.Error(w, err.Error(), 500)
     return
   }
   // wrap in base
   pageData := PageData{
     Title:   func() string { if book!=nil { return book.Title }; return "Notes" }(),
     Content: template.HTML(buf.String()),
   }
   if err := tmpl.ExecuteTemplate(w, "base", pageData); err != nil {
     http.Error(w, err.Error(), 500)
   }
}
