package models

import (
	"crypto/md5"
	"encoding/hex"
	"gorm.io/gorm"
)

type ShortLink struct {
	gorm.Model
	ShortCode   string
	OriginalUrl string
}

func GetShortCodeForUrl(url string) string {
	hash := md5.Sum([]byte(url))
	stringHash := hex.EncodeToString(hash[:])
	return stringHash[:8]
}
