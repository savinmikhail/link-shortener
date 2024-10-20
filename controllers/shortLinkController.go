package controllers

import (
	"encoding/json"
	"github.com/savinmikhail/link-shortener/models"
	"github.com/savinmikhail/link-shortener/repository"
	"log"
	"net/http"
)

type ShortenRequestData struct {
	URL string `json:"url"`
}

type ShortenResponseData struct {
	OriginalUrl  string `json:"originalUrl"`
	ShortenedUrl string `json:"shortenedUrl"`
}

func Shorten(w http.ResponseWriter, r *http.Request) {
	// get orig url
	var data ShortenRequestData
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	url := data.URL
	if url == "" {
		http.Error(w, "URL is required", http.StatusBadRequest)
		return
	}
	//get short url
	shortCode := models.GetShortCodeForUrl(url)
	//write to the file
	mappedUrls, err := repository.GetMappedUrls()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	mappedUrls[shortCode] = url
	err = repository.SaveMappedUrls(mappedUrls)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	//respond
	resp := ShortenResponseData{url, shortCode}
	jsonResp, err := json.Marshal(resp)
	if err != nil {
		log.Fatalf("Error happened in JSON marshal. Err: %s", err)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(jsonResp)
	return
}

func Redirect(w http.ResponseWriter, r *http.Request) {
	shortCode := r.URL.Path[1:]
	origUrl, err := repository.GetOrigUrlByShortCode(shortCode)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)

		jsonResp := map[string]string{"error": err.Error()}
		json.NewEncoder(w).Encode(jsonResp)
		return
	}
	http.Redirect(w, r, origUrl, http.StatusFound)
}
