package main

import (
	"net/http"
	"net/url"
	"encoding/json"
)

func (app *Application) handleShorten (w http.ResponseWriter, r *http.Request) {
	var request ShortenRequest
	var response LinkResponse
	// ensure correct http method is being used (POST)
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	// decode incoming JSON body 
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, "Invalid JSON payload", http.StatusBadRequest)
		return
	}

	// validate url
	parsedURL, err := url.Parse(request.URL)
	if err != nil || parsedURL.Scheme == "" || parsedURL.Host == "" {
		http.Error(w, "URL is not valid", http.StatusBadRequest)
		return
	}

	// call InsertLink() and save to db 
	response.ShortCode, response.CreatedAt, err = InsertLink(app.db, request.URL)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	response.ShortURL = "http://localhost:8080/" + response.ShortCode
	response.OriginalURL = request.URL

	// return response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, "Failed to encode response JSON", http.StatusInternalServerError)
		return
	}
}