package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (app *application) routes() *httprouter.Router {
	router := httprouter.New()

	router.HandlerFunc(http.MethodGet, "/v1/healthcheck", app.healthcheckHandler)
	router.HandlerFunc(http.MethodPost, "/v1/tracks", app.createTrackHandler)
	router.HandlerFunc(http.MethodGet, "/v1/tracks/:id", app.showTrackHandler)

	return router
}
