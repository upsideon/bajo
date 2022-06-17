package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/syndtr/goleveldb/leveldb"
)

func main() {
	urlDatabase, err := leveldb.OpenFile("url_database", nil)
	if err != nil {
		errorMessage := fmt.Sprintf("Error: Unable to access URL database: %s", err)
		panic(errorMessage)
	}
	defer urlDatabase.Close()

	router := initializeRouter(urlDatabase)
	router.Run(":8080")
}

func initializeRouter(urlDatabase URLDatabase) *gin.Engine {
	shortenController := ShortenController{
		URLDatabase: urlDatabase,
	}

	router := gin.Default()
	router.POST("/shorten", shortenController.Shorten)
	return router
}
