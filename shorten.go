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
	// CustomKeySizeLimit determines the maximum size of a custom key.
	CustomKeySizeLimit = 32

	// URLKeySize defines the number of characters in a URL key.
	// As URL keys are encoded in Base 64 there are 64 ^ URLKeySize possible keys.
	URLKeySize = 8

	// URLPrefix defines the prefix for shortened URLs.
	URLPrefix = "https://bajo"
)

// ShortenRequest represents a request to the URL shortening route.
type ShortenRequest struct {
	// Key contains an optional custom key with which to index the provided URL.
	Key string `form:"key" json:"key,omitempty" binding:"-"`
	// URL contains the URL to be shortened.
	URL string `form:"url" json:"url" binding:"required"`
}

// ShortenController contains logic and data related to the /shorten route.
type ShortenController struct {
	URLDatabase URLDatabase
}

// Shorten implements the logic for the /shorten route.
func (c *ShortenController) Shorten(context *gin.Context) {
	var shortenRequest ShortenRequest
	var URLKey string

	if err := context.BindJSON(&shortenRequest); err != nil {
		fmt.Println("Error: ", err)
		context.String(http.StatusBadRequest, "Bad Request")
		return
	}

	shortenRequestURLBytes := []byte(shortenRequest.URL)

	// When a custom key has not been provided, we generate one from the URL.
	if shortenRequest.Key == "" {
		hasher := sha256.New()
		hasher.Write(shortenRequestURLBytes)
		base64Hash := base64.RawURLEncoding.EncodeToString(hasher.Sum(nil))
		URLKey = base64Hash[:URLKeySize]
	} else {
		if len(shortenRequest.Key) > CustomKeySizeLimit {
			fmt.Println("Error: Custom key size is too large")
			context.String(http.StatusBadRequest, "Bad Request")
			return
		}
		URLKey = shortenRequest.Key
	}

	URLKeyBytes := []byte(URLKey)

	if _, err := c.URLDatabase.Get(URLKeyBytes, nil); err != nil {
		if err == dberror.ErrNotFound {
			// When not already present, the mapping between the URL key and URL is stored.
			if err = c.URLDatabase.Put(URLKeyBytes, shortenRequestURLBytes, nil); err != nil {
				fmt.Println("Error: ", err)
				context.String(http.StatusInternalServerError, "Internal Server Error")
				return
			}
		} else {
			// Any other error signals something unrecoverable, so we terminate the request.
			fmt.Println("Error: ", err)
			context.String(http.StatusInternalServerError, "Internal Server Error")
			return
		}
	}

	shortenedURL := fmt.Sprintf("%s/%s", URLPrefix, URLKey)
	context.JSON(http.StatusOK, gin.H{
		"shortened_url": shortenedURL,
	})
}
