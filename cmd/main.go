package main

import (
	"kosakata/internal/database"
	"kosakata/internal/game/sambungkata"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	_ "kosakata/docs" 
)

// @title Kosakata API
// @version 1.0
// @description This is the Kosakata API server.
// @host localhost:8080
// @BasePath /
func main() {

	db, err := database.InitDB()

	if err != nil {
		log.Fatal(err.Error())
	}

	err = db.AutoMigrate(&sambungkata.Word{})

	if err != nil {
		log.Fatal("Failed to migrate")
	}

	router := gin.Default()
	router.GET("/", rootHandler)
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	sambungkata.InitModule(router.Group("/word"), db)

	router.Run()
}

func rootHandler(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"name": "muhammad farikh",
		"age":  "23",
	})
}
