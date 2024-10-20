package controllers

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/savinmikhail/link-shortener/models"
	"github.com/savinmikhail/link-shortener/repository"
	"net/http"
)

type ShortenRequestData struct {
	URL string `json:"url"`
}

type ShortenResponseData struct {
	OriginalUrl  string `json:"originalUrl"`
	ShortenedUrl string `json:"shortenedUrl"`
}

func Shorten(c *gin.Context) {
	// get orig url
	var data ShortenRequestData
	err := c.BindJSON(&data)
	url := data.URL
	if url == "" {
		c.Error(errors.New("url is empty"))
		return
	}
	//get short url
	shortCode := models.GetShortCodeForUrl(url)
	//write to the file
	mappedUrls, err := repository.GetMappedUrls()
	if err != nil {
		c.Error(err)
		return
	}
	mappedUrls[shortCode] = url
	err = repository.SaveMappedUrls(mappedUrls)
	if err != nil {
		c.Error(err)
		return
	}
	//respond
	resp := ShortenResponseData{url, shortCode}

	c.JSON(http.StatusOK, resp)
}

func Redirect(c *gin.Context) {
	shortCode := c.Param("shortCode")
	origUrl, err := repository.GetOrigUrlByShortCode(shortCode)
	if err != nil {
		c.Error(err)
		return
	}
	c.Redirect(http.StatusMovedPermanently, origUrl)
}
