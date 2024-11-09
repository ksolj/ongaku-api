package main

import (
	"fmt"
	"net/http"
)

func (app *application) createTrackHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "create a new track")
}

func (app *application) showTrackHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil || id < 1 {
		http.NotFound(w, r)
		return
	}

	fmt.Fprintf(w, "show the details of track %d\n", id)
}
