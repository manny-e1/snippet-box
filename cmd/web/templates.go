package main

import (
	"github.com/manny-e1/snippetbox/internal/models"
	"github.com/manny-e1/snippetbox/ui"
	"html/template"
	"io/fs"
	"path/filepath"
	"time"
)

type templateData struct {
	CurrentYear     int
	Snippet         *models.Snippet
	Snippets        []*models.Snippet
	Form            any
	Flash           string
	IsAuthenticated bool
	CSRFToken       string
}

func humanDate(time time.Time) string {
	if time.IsZero() {
		return ""
	}
	return time.UTC().Format("02 Jan 2006 at 15:04")
}

var functions = template.FuncMap{
	"humanDate": humanDate,
}

func newTemplateCache() (map[string]*template.Template, error) {
	cache := map[string]*template.Template{}
	//pages, err := filepath.Glob("./ui/html/pages/*.tmpl")
	pages, err := fs.Glob(ui.Files, "html/pages/*.tmpl")
	if err != nil {
		return nil, err
	}
	for _, page := range pages {
		name := filepath.Base(page)
		patterns := []string{"html/layout/base.tmpl", "html/partials/*.tmpl", page}
		ts, err := template.New(name).Funcs(functions).ParseFS(ui.Files, patterns...)
		//if err != nil {
		//	return nil, err
		//}
		//ts, err = ts.ParseGlob("./ui/html/partials/*.tmpl")
		//if err != nil {
		//	return nil, err
		//}
		//ts, err = ts.ParseFiles(page)
		if err != nil {
			return nil, err
		}
		cache[name] = ts
	}
	return cache, nil
}
