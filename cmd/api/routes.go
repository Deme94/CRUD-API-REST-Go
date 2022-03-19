package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (app *application) routes() http.Handler {
	router := httprouter.New()

	router.HandlerFunc(http.MethodGet, "/status", app.statusHandler)

	router.HandlerFunc(http.MethodGet, "/v1/games", app.getAllGames)
	router.HandlerFunc(http.MethodGet, "/v1/game/:id", app.getOneGame)

	return app.enableCORS(router)
}
