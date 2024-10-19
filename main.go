package main

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"log"
	"net/http"
	"os"
)
import "encoding/json"

type ShortenRequestData struct {
	URL string `json:"url"`
}

type ShortenResponseData struct {
	OriginalUrl  string `json:"originalUrl"`
	ShortenedUrl string `json:"shortenedUrl"`
}

func shorten(w http.ResponseWriter, r *http.Request) {
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
	shortCode := getShortCodeForUrl(url)
	//write to the file
	mappedUrls, err := getMappedUrls()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	mappedUrls[shortCode] = url
	err = saveMappedUrls(mappedUrls)
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

func getShortCodeForUrl(url string) string {
	hash := md5.Sum([]byte(url))
	stringHash := hex.EncodeToString(hash[:])
	return stringHash[:8]
}

func saveMappedUrls(mappedUrls map[string]string) error {
	jsonContent, err := json.Marshal(mappedUrls)
	if err != nil {
		return err
	}
	err = os.WriteFile("list.json", jsonContent, 0777)
	if err != nil {
		return err
	}

	return nil
}

func redirect(w http.ResponseWriter, r *http.Request) {
	shortCode := r.URL.Path[1:]
	origUrl, err := getOrigUrlByShortCode(shortCode)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)

		jsonResp := map[string]string{"error": err.Error()}
		json.NewEncoder(w).Encode(jsonResp)
		return
	}
	http.Redirect(w, r, origUrl, http.StatusFound)
}

func getOrigUrlByShortCode(shortCode string) (string, error) {
	mappedUrls, err := getMappedUrls()
	if err != nil {
		return "", err
	}
	origUrl, exists := mappedUrls[shortCode]
	if !exists {
		return "", fmt.Errorf("short code %s not found", shortCode)
	}
	return origUrl, nil
}

func getMappedUrls() (map[string]string, error) {
	fileContent, err := os.ReadFile("list.json")
	if err != nil {
		return nil, err
	}
	mappedUrls := make(map[string]string)
	err = json.Unmarshal(fileContent, &mappedUrls)
	if err != nil {
		return nil, err
	}
	return mappedUrls, nil
}

func main() {
	http.HandleFunc("/shorten", shorten)
	http.HandleFunc("/", redirect)

	log.Fatal(http.ListenAndServe(":8079", nil))
}
