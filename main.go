package main

import (
  "embed"
  "html/template"
  "log"
  "net/http"
  "time"
)

//go:embed templates/* static/css/*
var assets embed.FS

var tmpl = template.Must(template.ParseFS(assets, "templates/*.html"))

func main() {
  mux := http.NewServeMux()

  // Serve static files
  mux.Handle("/static/", http.FileServer(http.FS(assets)))

  // Landing page
  mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
    tmpl.ExecuteTemplate(w, "index.html", nil)
  })

  // HTMX endpoint
  mux.HandleFunc("/api/time", func(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "text/html; charset=utf-8")
    w.Write([]byte(time.Now().Format("Mon, 02 Jan 2006 15:04:05 MST")))
  })

  log.Println("Listening on :8080")
  log.Fatal(http.ListenAndServe(":8080", mux))
}

