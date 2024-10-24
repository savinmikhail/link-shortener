package main

import (
	"github.com/gin-gonic/gin"
	"github.com/savinmikhail/link-shortener/controllers"
	"log"
)

func main() {
	r := gin.Default()

	r.POST("/shorten", controllers.Shorten)
	r.GET("/:shortCode", controllers.Redirect)

	err := r.Run(":8080")
	if err != nil {
		log.Fatal(err)
		return
	}
}
