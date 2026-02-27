package handlers

import (
	"fmt"
	"net/http"
)

// RobotsTxt serves the robots.txt file
func RobotsTxt(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	fmt.Fprint(w, "User-agent: *\nAllow: /\n\nSitemap: https://adamyounes.com/sitemap.xml\n")
}

// SitemapXML serves the sitemap.xml file
func SitemapXML(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/xml; charset=utf-8")
	fmt.Fprint(w, `<?xml version="1.0" encoding="UTF-8"?>
<urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">
  <url><loc>https://adamyounes.com/</loc><priority>1.0</priority></url>
  <url><loc>https://adamyounes.com/about</loc><priority>0.8</priority></url>
  <url><loc>https://adamyounes.com/projects</loc><priority>0.8</priority></url>
  <url><loc>https://adamyounes.com/experience</loc><priority>0.8</priority></url>
  <url><loc>https://adamyounes.com/notes</loc><priority>0.5</priority></url>
</urlset>
`)
}
