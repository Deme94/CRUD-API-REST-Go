package main

import (
	"CRUDWeb/models"
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/julienschmidt/httprouter"
)

// statusHandler handles /status
func (app *application) statusHandler(w http.ResponseWriter, r *http.Request) {
	currentStatus := AppStatus{
		Status:      "Available",
		Environment: app.config.env,
		Version:     version,
	}

	js, err := json.MarshalIndent(currentStatus, "", "\t")
	if err != nil {
		app.logger.Println(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(js)
}

// getAllGenres handles /v1/genres
func (app *application) getAllGenres(w http.ResponseWriter, r *http.Request) {
	genres, err := app.models.DB.GetAllGenres()
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	app.writeJSON(w, http.StatusOK, genres, "genres")
}

// getAllModes handles /v1/modes
func (app *application) getAllModes(w http.ResponseWriter, r *http.Request) {
	modes, err := app.models.DB.GetAllModes()
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	app.writeJSON(w, http.StatusOK, modes, "modes")
}

// getAllGames handles /v1/games
func (app *application) getAllGames(w http.ResponseWriter, r *http.Request) {
	games, err := app.models.DB.GetAllGames()
	if err != nil {
		app.errorJSON(w, err)
		return
	}
	for _, game := range games {
		game.ImageUrl = ""
	}

	app.writeJSON(w, http.StatusOK, games, "games")
}

// getOneGame handles /v1/game/:id
func (app *application) getOneGame(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())

	id, err := strconv.Atoi(params.ByName("id"))
	if err != nil {
		app.logger.Print(errors.New("invalid id parameter"))
		app.errorJSON(w, err)
		return
	}

	game, err := app.models.DB.GetOneGame(id)
	if err != nil {
		app.errorJSON(w, err)
		return
	}
	game.ImageUrl = ""

	app.writeJSON(w, http.StatusOK, game, "game")
}

// getGameImage handles /v1/game/:id/image
func (app *application) getGameImage(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())

	id, err := strconv.Atoi(params.ByName("id"))
	if err != nil {
		app.logger.Print(errors.New("invalid id parameter"))
		app.errorJSON(w, err)
		return
	}

	image, err := app.models.DB.GetGameImage(id)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	// abrimos imagen y la devolvemos al cliente
	file, err := ioutil.ReadFile("./images/" + image)
	if err != nil {
		app.logger.Print(err)
		app.errorJSON(w, err)
		return
	}

	w.Header().Set("Content-Type", "image/png")
	w.WriteHeader(http.StatusOK)
	w.Write(file)
}

// getAllImages handles /v1/game/images
func (app *application) getAllImages(w http.ResponseWriter, r *http.Request) {

	images, err := app.models.DB.GetAllImages()
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	arrayImagesFile := make(map[int][]byte)
	for id, imageUrl := range images {
		file, err := ioutil.ReadFile("./images/" + imageUrl)
		if err != nil {
			app.logger.Print(err)
			app.errorJSON(w, err)
			return
		}

		arrayImagesFile[id] = file
	}
	app.writeJSON(w, http.StatusOK, arrayImagesFile, "images")
}

// getAllGames handles /v1/games/genre/:genre
func (app *application) getAllGamesByGenre(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())

	genreID, err := strconv.Atoi(params.ByName("genre"))
	if err != nil {
		app.logger.Print(errors.New("invalid genre_id parameter"))
		app.errorJSON(w, err)
		return
	}

	games, err := app.models.DB.GetAllGamesByGenre(genreID)
	if err != nil {
		app.errorJSON(w, err)
		return
	}
	for _, game := range games {
		game.ImageUrl = ""
	}

	app.writeJSON(w, http.StatusOK, games, "games")
}

type gamePayload struct {
	Title       string    `json:"title"`
	Genres      []int     `json:"genres"`
	Modes       []int     `json:"modes"`
	Developers  []string  `json:"developers"`
	Publishers  []string  `json:"publishers"`
	ReleaseDate time.Time `json:"release_date"`
	Storage     int       `json:"storage"`
	Likes       int       `json:"likes"`
}

// insertGame handles /v1/games/insert
func (app *application) insertGame(w http.ResponseWriter, r *http.Request) {
	var payload gamePayload

	gameString := r.PostFormValue("game")
	err := json.Unmarshal([]byte(gameString), &payload)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	// leer imagen y crear archivo imagen en carpeta del proyecto (servidor)
	r.ParseMultipartForm(32 << 20)
	file, _, err := r.FormFile("image")
	if err != nil {
		app.errorJSON(w, err)
		return
	}
	defer file.Close()

	imageName := payload.Title + ".png"
	f, err := os.OpenFile("./images/"+imageName, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		app.errorJSON(w, err)
		return
	}
	defer f.Close()
	io.Copy(f, file)

	var game models.Game
	game.Title = payload.Title
	game.ImageUrl = imageName // cambiar por ruta del archivo creado
	game.Developers = payload.Developers
	game.Publishers = payload.Publishers
	game.ReleaseDate = payload.ReleaseDate
	game.Storage = payload.Storage
	game.Likes = 0

	err = app.models.DB.InsertGame(&game, payload.Genres, payload.Modes)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	type jsonResp struct {
		OK bool `json:"ok"`
	}

	ok := jsonResp{
		OK: true,
	}
	app.writeJSON(w, http.StatusOK, ok, "OK")
}

// updateGame handles /v1/games/update/:id
func (app *application) updateGame(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())

	id, err := strconv.Atoi(params.ByName("id"))
	if err != nil {
		app.logger.Print(errors.New("invalid id parameter"))
		app.errorJSON(w, err)
		return
	}

	var payload gamePayload

	gameString := r.PostFormValue("game")
	err = json.Unmarshal([]byte(gameString), &payload)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	// leer imagen y crear archivo imagen en carpeta del proyecto (servidor)
	r.ParseMultipartForm(32 << 20)
	file, _, err := r.FormFile("image")
	if err != nil {
		app.errorJSON(w, err)
		return
	}
	defer file.Close()

	imageName := payload.Title + ".png"
	f, err := os.OpenFile("./images/"+imageName, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		app.errorJSON(w, err)
		return
	}
	defer f.Close()
	io.Copy(f, file)

	var game models.Game
	game.Title = payload.Title
	game.ImageUrl = imageName // cambiar por ruta del archivo creado
	game.Developers = payload.Developers
	game.Publishers = payload.Publishers
	game.ReleaseDate = payload.ReleaseDate
	game.Storage = payload.Storage
	game.Likes = 0

	err = app.models.DB.UpdateGame(id, &game, payload.Genres, payload.Modes)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	type jsonResp struct {
		OK bool `json:"ok"`
	}

	ok := jsonResp{
		OK: true,
	}
	app.writeJSON(w, http.StatusOK, ok, "OK")
}

// deleteGame handles /v1/games/delete
func (app *application) deleteGame(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())

	id, err := strconv.Atoi(params.ByName("id"))
	if err != nil {
		app.logger.Print(errors.New("invalid id parameter"))
		app.errorJSON(w, err)
		return
	}

	err = app.models.DB.DeleteGame(id)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	type jsonResp struct {
		OK bool `json:"ok"`
	}

	ok := jsonResp{
		OK: true,
	}
	app.writeJSON(w, http.StatusOK, ok, "OK")
}
