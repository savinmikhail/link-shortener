package main

import (
	"github.com/savinmikhail/link-shortener/controllers"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/shorten", controllers.Shorten)
	http.HandleFunc("/", controllers.Redirect)

	log.Fatal(http.ListenAndServe(":3000", nil))
}
