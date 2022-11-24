package main

import (
	"errors"
	"fmt"
	"github.com/manny-e1/snippetbox/internal/models"
	"html/template"
	"net/http"
	"strconv"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	files := []string{
		"./ui/html/pages/home.tmpl",
		"./ui/html/layout/base.tmpl",
		"./ui/html/partials/footer.tmpl",
		"./ui/html/partials/nav.tmpl",
	}

	snippets, err := app.snippets.Latest()
	if err != nil {
		app.errorLogger.Println(err.Error())
		app.serverError(w, err)
		return
	}

	ts, err := template.ParseFiles(files...)
	if err != nil {
		app.errorLogger.Println(err.Error())
		app.serverError(w, err)
		return
	}

	data := &templateData{
		Snippets: snippets,
	}

	err = ts.ExecuteTemplate(w, "base", data)
	if err != nil {
		app.errorLogger.Println(err.Error())
		app.serverError(w, err)
	}
}
func (app *application) showSnippet(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || id < 1 {
		app.notFound(w)
		return
	}
	snippet, err := app.snippets.Get(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(w, err)
		}
	}
	files := []string{
		"./ui/html/layout/base.tmpl",
		"./ui/html/partials/nav.tmpl",
		"./ui/html/pages/view.tmpl",
		"./ui/html/partials/footer.tmpl",
	}
	ts, err := template.ParseFiles(files...)
	if err != nil {
		app.serverError(w, err)
		return
	}

	data := &templateData{
		Snippet: snippet,
	}

	err = ts.ExecuteTemplate(w, "base", data)
	if err != nil {
		app.serverError(w, err)
	}

}
func (app *application) createSnippet(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.Header().Set("Allow", http.MethodPost)
		app.clientError(w, http.StatusMethodNotAllowed)
		return
	}
	title := "Trial"
	content := "Example \nClimb Mount Fuji,\nBut slowly, slowly!\n\n– Kobayashi Issa"
	expires := 7
	id, err := app.snippets.Insert(title, content, expires)
	if err != nil {
		app.errorLogger.Println(err.Error())
		app.serverError(w, err)
		return
	}
	http.Redirect(w, r, fmt.Sprintf("/snippet?id=%d", id), http.StatusSeeOther)
}

func (app *application) transactionTrial(w http.ResponseWriter, r *http.Request) {
	err := app.transactionExample.InsertAndUpdate()
	if err != nil {
		app.errorLogger.Println(err.Error())
		app.serverError(w, err)
	}
	http.Redirect(w, r, fmt.Sprintf("/snippets"), http.StatusSeeOther)

}

func (app *application) showSnippets(w http.ResponseWriter, r *http.Request) {

	snippets, err := app.snippets.Latest()
	if err != nil {
		app.errorLogger.Println(err.Error())
		app.serverError(w, err)
		return
	}
	for _, snippet := range snippets {

		fmt.Fprintf(w, "%+v", snippet)
	}
}
