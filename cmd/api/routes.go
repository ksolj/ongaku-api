package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (app *application) routes() *httprouter.Router {
	router := httprouter.New()

	router.NotFound = http.HandlerFunc(app.notFoundResponse)

	router.MethodNotAllowed = http.HandlerFunc(app.methodNotAllowedResponse)

	router.HandlerFunc(http.MethodGet, "/v1/healthcheck", app.healthcheckHandler)

	router.HandlerFunc(http.MethodGet, "/v1/tracks", app.listTracksHandler)
	router.HandlerFunc(http.MethodPost, "/v1/tracks", app.createTrackHandler)
	router.HandlerFunc(http.MethodGet, "/v1/tracks/:id", app.showTrackHandler)
	router.HandlerFunc(http.MethodPatch, "/v1/tracks/:id", app.updateTrackHandler)
	router.HandlerFunc(http.MethodDelete, "/v1/tracks/:id", app.deleteTrackHandler)

	router.HandlerFunc(http.MethodGet, "/v1/tabs/:filename", app.getTabsHandler)
	return router
}
