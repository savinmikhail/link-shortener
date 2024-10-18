package main

import (
	"crypto/md5"
	"fmt"
	"log"
	"net/http"
	"os"
)
import "encoding/json"

func shorten(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		fmt.Fprintf(w, err.Error())
		return
	}
	url := r.PostForm.Get("url")
	shortUrl := md5.Sum([]byte(url))
	mappedUrls := make(map[string]string)
	mappedUrls[string(shortUrl[:])] = url
	jsonContent, err := json.Marshal(mappedUrls)
	if err != nil {
		fmt.Fprintf(w, err.Error())
		return
	}
	os.WriteFile("list.json", jsonContent, 0777)
	fmt.Fprintf(w, "saved %v", url)
}

func redirect(w http.ResponseWriter, r *http.Request) {
	shortCode := r.URL.Path[1:]
	jsonContent, err := os.ReadFile("list.json")
	if err != nil {
		fmt.Fprintf(w, err.Error())
	}
	mappedUrls := make(map[string]string)
	err1 := json.Unmarshal(jsonContent, &mappedUrls)
	if err1 != nil {
		fmt.Fprintf(w, err1.Error())
	}
	origUrl := mappedUrls[shortCode]
	http.Redirect(w, r, origUrl, http.StatusFound)
}

func main() {
	http.HandleFunc("/shorten", shorten)
	http.HandleFunc("/", redirect)
	//http.HandleFunc("/", handler)

	log.Fatal(http.ListenAndServe(":8079", nil))
}
