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

type RequestData struct {
	URL string `json:"url"`
}

func shorten(w http.ResponseWriter, r *http.Request) {
	// get orig url
	var data RequestData
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
	hash := md5.Sum([]byte(url))
	stringHash := hex.EncodeToString(hash[:])
	shortUrl := stringHash[:8]
	//write to the file
	mappedUrls, err := getMappedUrls()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	mappedUrls[shortUrl] = url
	err = saveMappedUrls(mappedUrls)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	//respond
	fmt.Fprintf(w, "saved %v", url)
	fmt.Fprintf(w, "\n short code: http://localhost:8079/%v", shortUrl)
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
		http.Error(w, err.Error(), http.StatusBadRequest)
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
