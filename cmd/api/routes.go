package main

import (
	"expvar"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (app *application) routes() http.Handler {
	router := httprouter.New()

	router.NotFound = http.HandlerFunc(app.notFoundResponse)
	router.MethodNotAllowed = http.HandlerFunc(app.methodNotAllowedResponse)

	router.HandlerFunc(http.MethodGet, "/v1/healthcheck", app.healthcheckHandler)

	router.HandlerFunc(http.MethodGet, "/v1/tracks", app.listTracksHandler)
	router.HandlerFunc(http.MethodPost, "/v1/tracks", app.requirePermission("tracks:write", app.createTrackHandler))
	router.HandlerFunc(http.MethodGet, "/v1/tracks/:id", app.showTrackHandler)
	router.HandlerFunc(http.MethodPatch, "/v1/tracks/:id", app.requirePermission("tracks:write", app.updateTrackHandler))
	router.HandlerFunc(http.MethodDelete, "/v1/tracks/:id", app.requirePermission("tracks:write", app.deleteTrackHandler))

	router.HandlerFunc(http.MethodGet, "/v1/tabs/:filename", app.getTabsHandler)

	router.HandlerFunc(http.MethodPost, "/v1/users", app.registerUserHandler)
	router.HandlerFunc(http.MethodPut, "/v1/users/activated", app.activateUserHandler)
	router.HandlerFunc(http.MethodPut, "/v1/users/password", app.updateUserPasswordHandler)

	router.HandlerFunc(http.MethodPost, "/v1/tokens/authentication", app.createAuthenticationTokenHandler)
	router.HandlerFunc(http.MethodPost, "/v1/tokens/password-reset", app.createPasswordResetTokenHandler)

	router.Handler(http.MethodGet, "/debug/vars", expvar.Handler())

	return app.metrics(app.recoverPanic(app.enableCORS(app.rateLimit(app.authenticate(router)))))
}
