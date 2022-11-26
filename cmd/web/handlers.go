package main

import (
	"errors"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"github.com/manny-e1/snippetbox/internal/models"
	"net/http"
	"strconv"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	snippets, err := app.snippets.Latest()
	if err != nil {
		app.errorLogger.Println(err.Error())
		app.serverError(w, err)
		return
	}

	data := app.newTemplateData(r)
	data.Snippets = snippets

	app.render(w, http.StatusOK, "home.tmpl", data)
}
func (app *application) showSnippet(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())

	id, err := strconv.Atoi(params.ByName("id"))
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

	data := app.newTemplateData(r)
	data.Snippet = snippet
	app.render(w, http.StatusOK, "view.tmpl", data)

}
func (app *application) createSnippet(w http.ResponseWriter, r *http.Request) {

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
