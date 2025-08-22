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

	wordRepository := sambungkata.NewRepository(db)
	wordService := sambungkata.NewService(wordRepository)
	wordHandler := sambungkata.NewHandler(wordService)

	router := gin.Default()

	router.GET("/", rootHandler)
	router.Group("/word").
		GET("/", wordHandler.ShowAllWord).
		GET("/:id", wordHandler.ShowWord).
		POST("/store", wordHandler.StoreWord).
		GET("/today", wordHandler.GetTodayWord).
		POST("/check", wordHandler.CheckingWord)

	router.Run()
}

func rootHandler(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"name": "muhammad farikh",
		"age":  "23",
	})
}
