package main

import (
	"crypto/sha256"
	"encoding/base64"
	"fmt"

	"github.com/gin-gonic/gin"
)

const (
	// URLKeySize defines the number of characters in a URL key.
	// As URL keys are encoded in Base 64 there are 64 ^ URLKeySize possible keys.
	URLKeySize = 8

	// URLPrefix defines the prefix for shortened URLs.
	URLPrefix = "https://bajo"
)

// ShortenRequest represents a request to the URL shortening route.
type ShortenRequest struct {
	URL string `form:"url"`
}

func shorten(c *gin.Context) {
	var shortenRequest ShortenRequest

	if c.ShouldBind(&shortenRequest) != nil {
		c.String(400, "Bad Request")
		return
	}

	hasher := sha256.New()
	hasher.Write([]byte(shortenRequest.URL))
	base64Hash := base64.RawURLEncoding.EncodeToString(hasher.Sum(nil))

	URLKey := base64Hash[:URLKeySize]
	shortenedURL := fmt.Sprintf("%s/%s", URLPrefix, URLKey)

	c.JSON(200, gin.H{
		"shortened_url": shortenedURL,
	})
}
