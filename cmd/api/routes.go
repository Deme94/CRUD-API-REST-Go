package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (app *application) routes() http.Handler {
	router := httprouter.New()

	router.HandlerFunc(http.MethodGet, "/status", app.statusHandler)

	router.HandlerFunc(http.MethodGet, "/v1/genres", app.getAllGenres)
	router.HandlerFunc(http.MethodGet, "/v1/modes", app.getAllModes)

	router.HandlerFunc(http.MethodGet, "/v1/games", app.getAllGames)
	router.HandlerFunc(http.MethodGet, "/v1/game/:id", app.getOneGame)
	router.HandlerFunc(http.MethodGet, "/v1/game/:id/image", app.getGameImage)
	router.HandlerFunc(http.MethodGet, "/v1/games/images", app.getAllImages)
	router.HandlerFunc(http.MethodGet, "/v1/games/genre/:genre", app.getAllGamesByGenre)

	router.HandlerFunc(http.MethodPut, "/v1/games/insert", app.insertGame)
	router.HandlerFunc(http.MethodPut, "/v1/games/update/:id", app.updateGame)
	router.HandlerFunc(http.MethodDelete, "/v1/games/delete/:id", app.deleteGame)

	return app.enableCORS(router)
}
