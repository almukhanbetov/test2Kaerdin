package main

import (
	"github.com/gin-contrib/cors"
	"net/http"
	"test2/config"
	"test2/db"
	"test2/models"
	"test2/utils"
	
	"github.com/gin-gonic/gin"
	
)

func main() {

	config.LoadEnv()
	db.InitDB()
	defer db.CloseDB()
	db.RunMigrations()
	utils.LoadLiveGamesFromJSON("allLiveGAmes.json")
	utils.EnsureUploadsFolder()
	
	router := gin.Default()
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5174"},
		AllowMethods:     []string{"GET", "POST", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))
	// 👉 Новый эндпоинт
	router.GET("/games", func(c *gin.Context) {
		games, err := models.GetAllLiveGames()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при получении матчей"})
			return
		}
		c.JSON(http.StatusOK, games)
	})
	utils.StartServerGracefully(router, "8082")
}
