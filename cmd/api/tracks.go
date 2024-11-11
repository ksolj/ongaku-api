package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/ksolj/ongaku-api/internal/data"
	"github.com/ksolj/ongaku-api/internal/data/validator"
)

func (app *application) createTrackHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Name     string   `json:"name"`
		Duration int32    `json:"duration"`
		Artists  []string `json:"artists"`
		Album    string   `json:"album"`
		// Tabs []Tab `json:"tabs"` // TODO: will be implemented later
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	track := &data.Track{
		Name:     input.Name,
		Duration: input.Duration,
		Artists:  input.Artists,
		Album:    input.Album,
	}

	v := validator.New()

	if !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	if data.ValidateTrack(v, track); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	fmt.Fprintf(w, "%+v\n", input)
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
