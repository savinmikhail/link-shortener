package main

import (
	"github.com/gin-gonic/gin"
	"github.com/savinmikhail/link-shortener/controllers"
)

func main() {
	r := gin.Default()

	r.POST("/shorten", controllers.Shorten)
	r.GET("/:shortCode", controllers.Redirect)

	r.Run(":8080")
}
