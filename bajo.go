package main

import (
	"github.com/gin-gonic/gin"
)

func main() {
	databaseManager := &LevelDBDatabaseManager{}
	urlDatabase := GetURLDatabase(databaseManager)
	defer urlDatabase.Close()
	router := initializeRouter(urlDatabase)
	router.Run(":8080")
}

func initializeRouter(urlDatabase URLDatabase) *gin.Engine {
	shortenController := ShortenController{
		URLDatabase: urlDatabase,
	}

	redirectController := RedirectController{
		URLDatabase: urlDatabase,
	}

	router := gin.Default()
	router.POST("/shorten", shortenController.Shorten)
	router.GET("/:key", redirectController.Redirect)
	return router
}
