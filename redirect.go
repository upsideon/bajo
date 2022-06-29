package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	dberror "github.com/syndtr/goleveldb/leveldb/errors"
)

// RedirectController manages URL redirection.
type RedirectController struct {
	URLDatabase URLDatabase
}

// Redirect implements the logic for URL redirection.
func (c *RedirectController) Redirect(context *gin.Context) {
	URLKey := context.Param("key")

	urlBytes, err := c.URLDatabase.Get([]byte(URLKey), nil)

	if err != nil {
		if err == dberror.ErrNotFound {
			context.String(http.StatusNotFound, "Not Found")
			return
		} else {
			context.String(http.StatusInternalServerError, "Internal Server Error")
			return
		}
	}

	context.Redirect(http.StatusFound, string(urlBytes))
}
