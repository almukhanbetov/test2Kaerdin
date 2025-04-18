package models

import (
	"test2/db"
	"log"
)

type LiveGame struct {
	ID         string `db:"id" json:"id"`
	Name       string `db:"name" json:"name"`
	League     string `db:"league" json:"league"`
	Score      string `db:"score" json:"score"`
	MatchTime  string `db:"match_time" json:"match_time"`
	UpdateTime string `db:"update_time" json:"update_time"`
}

func GetAllLiveGames() ([]LiveGame, error) {
	var games []LiveGame
	err := db.DB.Select(&games, "SELECT * FROM live_games ORDER BY update_time DESC")
	if err != nil {
		log.Printf("❌ Ошибка при запросе матчей: %v", err)
		return nil, err
	}
	return games, nil
}
