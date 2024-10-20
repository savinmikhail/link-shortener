package main

import (
	"github.com/savinmikhail/go_crud/initializers"
	"github.com/savinmikhail/link-shortener/models"
)

func init() {
	initializers.LoadEnvVariables()
	initializers.ConnectToDB()
}

func main() {
	initializers.DB.AutoMigrate(&models.ShortLink{})
}
