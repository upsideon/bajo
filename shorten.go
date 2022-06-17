package main

import (
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	dberror "github.com/syndtr/goleveldb/leveldb/errors"
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
	URL string `form:"url" json:"url" binding:"required"`
}

// ShortenController contains logic and data related to the /shorten route.
type ShortenController struct {
	URLDatabase URLDatabase
}

// Shorten implements the logic for the /shorten route.
func (c *ShortenController) Shorten(context *gin.Context) {
	var shortenRequest ShortenRequest

	if context.BindJSON(&shortenRequest) != nil {
		context.String(http.StatusBadRequest, "Bad Request")
		return
	}

	shortenRequestURLBytes := []byte(shortenRequest.URL)

	hasher := sha256.New()
	hasher.Write(shortenRequestURLBytes)
	base64Hash := base64.RawURLEncoding.EncodeToString(hasher.Sum(nil))
	URLKey := base64Hash[:URLKeySize]

	URLKeyBytes := []byte(URLKey)

	if _, err := c.URLDatabase.Get(URLKeyBytes, nil); err != nil {
		if err == dberror.ErrNotFound {
			// When not already present, the mapping between the URL key and URL is stored.
			if err = c.URLDatabase.Put(URLKeyBytes, shortenRequestURLBytes, nil); err != nil {
				context.String(http.StatusInternalServerError, "Internal Server Error")
				return
			}
		} else {
			// Any other error signals something unrecoverable, so we terminate the request.
			context.String(http.StatusInternalServerError, "Internal Server Error")
			return
		}
	}

	shortenedURL := fmt.Sprintf("%s/%s", URLPrefix, URLKey)
	context.JSON(http.StatusOK, gin.H{
		"shortened_url": shortenedURL,
	})
}
