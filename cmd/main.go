package main

import (
	"fmt"
	"kosakata/internal/game/sambungkata"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("failed load environment")
	}

	DB_HOST := os.Getenv("DB_HOST")
	DB_PORT := os.Getenv("DB_PORT")
	DB_NAME := os.Getenv("DB_NAME")
	DB_USER := os.Getenv("DB_USER")
	DB_PASS := os.Getenv("DB_PASS")

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", DB_USER, DB_PASS, DB_HOST, DB_PORT, DB_NAME) 
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
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
