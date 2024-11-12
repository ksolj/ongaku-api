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
		Tabs:     []data.Tab{"null"}, // TODO: implement!
	}

	v := validator.New()

	if data.ValidateTrack(v, track); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	// Call the Insert() method on our movies model, passing in a pointer to the
	// validated movie struct. This will create a record in the database and update the
	// movie struct with the system-generated information.
	err = app.models.Tracks.Insert(track)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	// When sending a HTTP response, we want to include a Location header to let the
	// client know which URL they can find the newly-created resource at. We make an
	// empty http.Header map and then use the Set() method to add a new Location header,
	// interpolating the system-generated ID for our new movie in the URL.
	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/v1/tracks/%d", track.ID))

	// Write a JSON response with a 201 Created status code, the movie data in the
	// response body, and the Location header.
	err = app.writeJSON(w, http.StatusCreated, envelope{"track": track}, headers)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
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
