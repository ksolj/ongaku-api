package main

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/ksolj/ongaku-api/internal/data"
	"github.com/ksolj/ongaku-api/internal/data/validator"
)

const uploadDir = "./uploads/tabs"
const baseFileURL = "http://localhost:4000/tabs"

func (app *application) createTrackHandler(w http.ResponseWriter, r *http.Request) {
	contentType := r.Header.Get("Content-Type")
	if !strings.HasPrefix(contentType, "multipart/form-data") {
		app.badRequestResponse(w, r, fmt.Errorf("Content-Type must be multipart/form-data"))
		return
	}

	var input struct {
		Name     string   `json:"name"`
		Duration int32    `json:"duration"`
		Artists  []string `json:"artists"`
		Album    string   `json:"album"`
	}

	// Get json data from the request and convert it to the ReaderCloser to pass it to our readJSON helper.
	jsonPart := r.FormValue("json_data")
	jsonPartReader := strings.NewReader(jsonPart)
	r.Body = io.NopCloser(jsonPartReader)

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		app.badRequestResponse(w, r, fmt.Errorf("File must be provided"))
	}
	defer file.Close()

	// This part is temporary 'cause in the future cloud object storage will be used (S3)
	filePath := filepath.Join(uploadDir, header.Filename) // TODO: get rid of this later. Use S3
	dst, err := os.Create(filePath)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	defer dst.Close()

	if _, err := io.Copy(dst, file); err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	fileURL := fmt.Sprintf("%s/%s", baseFileURL, header.Filename)

	track := &data.Track{
		Name:     input.Name,
		Duration: input.Duration,
		Artists:  input.Artists,
		Album:    input.Album,
		Tabs:     fileURL,
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
