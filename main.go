package main

import (
	
	"test2/config"
	"test2/db"
	"test2/handlers"
	"test2/utils"
	"test2/importer"	
	_ "github.com/jackc/pgx/v5/stdlib"
    "github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	config.LoadEnv()
	db.InitDB()
	defer db.CloseDB()

	db.RunMigrations()            // ← миграции применятся при запуске
	// utils.ImportMatchesFromJSON() // ← ВОТ ЭТА СТРОКА 📌
	importer.ImportAllLiveGames(db.DB, "data/allLiveGAmes.json")

	router := gin.Default()
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"},
		AllowMethods:     []string{"GET", "POST", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type"},
		AllowCredentials: true,
	}))
	router.GET("/live", handlers.GetLiveEvents)

	utils.StartServerGracefully(router, "8686")
}
