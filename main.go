package main

import (
  "embed"
  "html/template"
  "log"
  "net/http"
  "time"
)

// embed everything under templates/ and static/css/
//go:embed templates/* static/css/*
var assets embed.FS

var tmpl = template.Must(template.ParseFS(assets, "templates/*.html"))

func main() {
  mux := http.NewServeMux()

  // Serve CSS (and any other static files you add there)
  mux.Handle("/static/", http.FileServer(http.FS(assets)))

  // Landing page
  mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
    if err := tmpl.ExecuteTemplate(w, "index.html", nil); err != nil {
      http.Error(w, err.Error(), http.StatusInternalServerError)
    }
  })

  // HTMX endpoint for current time
  mux.HandleFunc("/api/time", func(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "text/html; charset=utf-8")
    now := time.Now().Format("Mon, 02 Jan 2006 15:04:05 MST")
    // return just the fragment that HTMX will swap into <div id="time">
    w.Write([]byte(now))
  })

  log.Println("Listening on :8080")
  log.Fatal(http.ListenAndServe(":8080", mux))
}

