package main

import (
	"net/http"
	"path/filepath"
)

// This method will probably be changed in the future due to possible use of S3 storage
func (app *application) getTabsHandler(w http.ResponseWriter, r *http.Request) {
	filename, err := app.readFilenameParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	filePath := filepath.Join(uploadDir, filename)
	http.ServeFile(w, r, filePath)
}
