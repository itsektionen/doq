// Package pages contain code for serving pages
package pages

import (
	"embed"
	"html/template"
	"log"
	"math/rand/v2"
	"net/http"
	"strings"

	"github.com/itsektionen/doq/app/static"
)

var quotes = []string{
	"Carpe Diem",
	"Pics or it didn't happen",
	"It's DOQing time",
	"Sorry, DOQtor's orders",
	"DOQ DOQ DOQ",
}

//go:embed templates
var templateFS embed.FS

type Pages struct {
	base *template.Template
}

type PageData struct {
	Title string
	Quote string
	Data  any
}

func New() Pages {
	base, err := template.ParseFS(templateFS, "templates/layouts/*.html")
	if err != nil {
		log.Fatalf("failed to parse base layouts: %v", err)
	}
	return Pages{base: base}
}

func (p Pages) render(w http.ResponseWriter, pageFile string, data any, statusCode ...int) {
	tmpl, err := p.base.Clone()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	tmpl, err = tmpl.ParseFS(templateFS, "templates/"+pageFile)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	pageData := PageData{
		Title: "DOQ",
	}

	switch v := data.(type) {
	case nil:
		// keep default PageData with Title: "DOQ"
	case PageData:
		pageData = v
		if pageData.Title == "" {
			pageData.Title = "DOQ"
		}
	case *PageData:
		if v != nil {
			pageData = *v
			if pageData.Title == "" {
				pageData.Title = "DOQ"
			}
		}
	case string:
		if v != "" {
			pageData.Quote = v
		}
	case map[string]any:
		if title, ok := v["Title"].(string); ok && title != "" {
			pageData.Title = title
		}
		if quote, ok := v["Quote"].(string); ok && quote != "" {
			pageData.Quote = quote
		}
		pageData.Data = v
	default:
		pageData.Data = v
	}

	status := http.StatusOK
	if len(statusCode) > 0 && statusCode[0] > 0 {
		status = statusCode[0]
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(status)
	if err := tmpl.ExecuteTemplate(w, "base", pageData); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (p Pages) HandleIndex(w http.ResponseWriter, r *http.Request) {
	randomQuote := quotes[rand.N(len(quotes))]

	p.render(w, "index.html", PageData{
		Quote: randomQuote,
	})
}

func (p Pages) HandleStatic() http.Handler {
	fileServer := http.StripPrefix("/static/", static.Handler())
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path := strings.TrimPrefix(r.URL.Path, "/static/")
		if !static.Exists(path) {
			p.render(w, "404.html", PageData{Title: "404 - Not Found"}, http.StatusNotFound)
			return
		}
		fileServer.ServeHTTP(w, r)
	})
}

func (p Pages) HandleNotFound(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/")
	if static.Exists(path) {
		static.Handler().ServeHTTP(w, r)
		return
	}
	p.render(w, "404.html", PageData{Title: "404 - Not Found"}, http.StatusNotFound)
}
