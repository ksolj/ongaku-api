package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/ksolj/ongaku-api/internal/data"
)

func (app *application) createTrackHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "create a new track")
}

func (app *application) showTrackHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil || id < 1 {
		app.notFoundResponse(w, r)
		return
	}

	track := data.Track{
		ID:        id,
		CreatedAt: time.Now(),
		Name:      "moment A rhythm",
		Duration:  1009,
		Artists:   []string{"Ling Tosite Sigure"},
		Album:     "Just a Moment",
		Tabs:      []data.Tab{"Test"},
		Version:   1,
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"track": track}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
