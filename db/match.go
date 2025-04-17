package db

import (
	"log"

	"test2/models"
)

func SaveMatch(match models.Match) error {
	query := `
		INSERT INTO matches (match_id, league, match_name, score, game_time, half)
		VALUES (:match_id, :league, :match_name, :score, :game_time, :half)
		ON CONFLICT (match_id) DO NOTHING
	`
	_, err := DB.NamedExec(query, match)
	if err != nil {
		log.Printf("❌ Ошибка вставки матча %s: %v", match.MatchID, err)
	}
	return err
}
