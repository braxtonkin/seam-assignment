package main

import (
	"encoding/json"
	"net/http"
)

type wrapper map[string]any

// WriteJSON Utility function for sending JSON in an http response
// w 	   - the http response to be written to
// status  - the status code of
// data    - the data to encode to JSON
//
//	By default, the only header included is { "Content-Type" : "application/json" }
func (app *application) writeJSON(w http.ResponseWriter, status int, data any) error {
	// Using marshalindent for greater readibility when making curl requests
	js, err := json.MarshalIndent(data, "", "\t")
	if err != nil {
		return err
	}
	// appending new line for readability
	js = append(js, '\n')

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if _, err := w.Write(js); err != nil {
		return err
	}

	return nil
}

// Sends a generic error response with the provdied status code and message
func (app *application) writeError(w http.ResponseWriter, r *http.Request, status int, message any) {
	errorResponse := wrapper{"error": message}

	err := app.writeJSON(w, status, errorResponse)
	if err != nil {
		app.logger.Error(err.Error(), "method", r.Method, "uri", r.URL.RequestURI())
		w.WriteHeader(http.StatusInternalServerError)
	}
}

// To log unexpected problems at runtime
// Largely the same as writeError, but logs the error and
func (app *application) writeServerError(w http.ResponseWriter, r *http.Request, err error) {
	app.logger.Error(err.Error(), "method", r.Method, "uri", r.URL.RequestURI())

	app.writeError(w, r, http.StatusInternalServerError, "server encountered a problem")
}
