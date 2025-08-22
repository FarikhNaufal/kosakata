package main

import (
	"kosakata/internal/database"
	"kosakata/internal/game/sambungkata"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {

	db, err := database.InitDB()
	
	if err != nil {
		log.Fatal("Failed load db.")
	}

	err = db.AutoMigrate(&sambungkata.Word{})

	if err != nil {
		log.Fatal("Failed to migrate")
	}

	router := gin.Default()
	router.GET("/", rootHandler)

	sambungkata.InitModule(router.Group("/word"), db)
	

	router.Run()
}

func rootHandler(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"name": "muhammad farikh",
		"age":  "23",
	})
}
