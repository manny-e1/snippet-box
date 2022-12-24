package main

import (
	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
	"github.com/manny-e1/snippetbox/ui"
	"net/http"
)

func (app *application) routes() http.Handler {
	//mux := http.NewServeMux()
	router := httprouter.New()

	router.NotFound = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		app.notFound(w)
	})

	//fileServer := http.FileServer(http.Dir("./ui/static/"))
	//Using embedded filesystem
	fileServer := http.FileServer(http.FS(ui.Files))

	//router.Handler(http.MethodGet, "/static/*filepath", http.StripPrefix("/static", fileServer))
	//We no longer need to strip the prefix from the request url
	router.Handler(http.MethodGet, "/static/*filepath", fileServer)
	router.Handler(http.MethodGet, "/ping", http.HandlerFunc(Ping))

	dynamic := alice.New(app.sessionManager.LoadAndSave, noSurf, app.authenticate)

	router.Handler(http.MethodGet, "/", dynamic.ThenFunc(app.home))
	router.Handler(http.MethodGet, "/snippet/view/:id", dynamic.ThenFunc(app.showSnippet))
	router.Handler(http.MethodGet, "/transaction-example", dynamic.ThenFunc(app.transactionTrial))
	router.Handler(http.MethodGet, "/user/signup", dynamic.ThenFunc(app.userSignup))
	router.Handler(http.MethodPost, "/user/signup", dynamic.ThenFunc(app.userSignupPost))
	router.Handler(http.MethodGet, "/user/login", dynamic.ThenFunc(app.userLogin))
	router.Handler(http.MethodPost, "/user/login", dynamic.ThenFunc(app.userLoginPost))

	protected := dynamic.Append(app.requireAuthentication)

	router.Handler(http.MethodGet, "/snippet/create", protected.ThenFunc(app.createSnippet))
	router.Handler(http.MethodPost, "/snippet/create", protected.ThenFunc(app.createSnippetPost))
	router.Handler(http.MethodPost, "/user/logout", protected.ThenFunc(app.userLogoutPost))

	//mux.Handle("/static", http.NotFoundHandler())
	standard := alice.New(app.recoverPanic, app.logRequest, secureHeaders)
	return standard.Then(router)
}
