package utils

import (
	"encoding/json"
	"os"
	"log"
	"path/filepath"

	"test2/db"
	"test2/models"
)

func ImportMatchesFromJSON() {
	path := filepath.Join("data", "allLiveGAmes.json")
	content, err := os.ReadFile(path)
	if err != nil {
		log.Fatalf("❌ Не удалось прочитать JSON: %v", err)
	}

	var raw models.AllLiveGames
	if err := json.Unmarshal(content, &raw); err != nil {
		log.Fatalf("❌ Ошибка парсинга JSON: %v", err)
	}

	count := 0
	for _, group := range raw.Results {
		for _, item := range group {
			if item.Type == "EV" {
				match := models.Match{
					MatchID:   item.ID,
					League:    item.CT,
					MatchName: item.NA,
					Score:     item.SS,
					GameTime:  item.TM,
					Half:      item.GO,
				}
				if err := db.SaveMatch(match); err == nil {
					count++
				}
			}
		}
	}

	log.Printf("📥 Импортировано матчей: %d", count)
}
