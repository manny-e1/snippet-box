package main

import (
	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
	"net/http"
)

func (app *application) routes() http.Handler {
	//mux := http.NewServeMux()
	router := httprouter.New()

	router.NotFound = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		app.notFound(w)
	})

	fileServer := http.FileServer(http.Dir("./ui/static/"))
	router.Handler(http.MethodGet, "/static/*filepath", http.StripPrefix("/static", fileServer))

	dynamic := alice.New(app.sessionManager.LoadAndSave)

	router.Handler(http.MethodGet, "/", dynamic.ThenFunc(app.home))
	router.Handler(http.MethodGet, "/snippet/view/:id", dynamic.ThenFunc(app.showSnippet))
	router.Handler(http.MethodGet, "/snippets", dynamic.ThenFunc(app.showSnippets))
	router.Handler(http.MethodGet, "/transaction-example", dynamic.ThenFunc(app.transactionTrial))
	router.Handler(http.MethodGet, "/snippet/create", dynamic.ThenFunc(app.createSnippet))
	router.Handler(http.MethodPost, "/snippet/create", dynamic.ThenFunc(app.createSnippetPost))

	//mux.Handle("/static", http.NotFoundHandler())
	standard := alice.New(app.recoverPanic, app.logRequest, secureHeaders)
	return standard.Then(router)
}
