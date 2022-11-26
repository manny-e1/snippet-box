package main

import (
	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
	"net/http"
)

func (app *application) routes() http.Handler {
	//mux := http.NewServeMux()
	router := httprouter.New()

	router.HandlerFunc(http.MethodGet, "/", app.home)
	router.HandlerFunc(http.MethodGet, "/snippet/view/:id", app.showSnippet)
	router.HandlerFunc(http.MethodGet, "/snippets", app.showSnippets)
	router.HandlerFunc(http.MethodGet, "/transaction-example", app.transactionTrial)

	router.HandlerFunc(http.MethodPost, "/snippet/create", app.createSnippet)

	fileServer := http.FileServer(http.Dir("./ui/static/"))
	//mux.Handle("/static", http.NotFoundHandler())
	router.Handler(http.MethodGet, "/static/*filepath", http.StripPrefix("/static", fileServer))
	standard := alice.New(app.recoverPanic, app.logRequest, secureHeaders)
	return standard.Then(router)
}
