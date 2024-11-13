package main

import (
	"errors"
	"fmt"
	"net/http"

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

	err = app.models.Tracks.Insert(track)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	// When sending a HTTP response, we want to include a Location header to let the
	// client know which URL they can find the newly-created resource at. We make an
	// empty http.Header map and then use the Set() method to add a new Location header,
	// interpolating the system-generated ID for our new track in the URL.
	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/v1/tracks/%d", track.ID))

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

	track, err := app.models.Tracks.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"track": track}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) updateTrackHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	track, err := app.models.Tracks.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	var input struct {
		Name     *string  `json:"name"`
		Duration *int32   `json:"duration"`
		Artists  []string `json:"artists"`
		Album    *string  `json:"album"`
		// Tabs     []data.Tab `json:"tabs"` // TODO: implement TAB PROPERLY!!!!
	}

	err = app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if input.Name != nil {
		track.Name = *input.Name
	}

	if input.Duration != nil {
		track.Duration = *input.Duration
	}

	if input.Artists != nil {
		track.Artists = input.Artists
	}

	if input.Album != nil {
		track.Album = *input.Album
	}

	// track.Tabs = input.Tabs

	v := validator.New()

	if data.ValidateTrack(v, track); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = app.models.Tracks.Update(track)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrEditConflict):
			app.editConflictResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"track": track}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) deleteTrackHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	err = app.models.Tracks.Delete(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"message": "track successfully deleted"}, nil) // May change to '204 No Content' code in the future
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
