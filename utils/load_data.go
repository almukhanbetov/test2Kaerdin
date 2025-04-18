package utils

import (
	"encoding/json"
	"io"
	"log"
	"os"

	"test2/db"
	"test2/models"
)

// rawEntry — для предварительного парсинга JSON-элементов
type rawEntry struct {
	Type string `json:"type"`
	ID   string `json:"ID"`
	NA   string `json:"NA"`
	CT   string `json:"CT"`
	SS   string `json:"SS"`
	TM   string `json:"TM"`
	TU   string `json:"TU"`
}

func LoadLiveGamesFromJSON(filepath string) {
	file, err := os.Open(filepath)
	if err != nil {
		log.Fatalf("❌ Ошибка открытия файла: %v", err)
	}
	defer file.Close()

	byteValue, err := io.ReadAll(file)
	if err != nil {
		log.Fatalf("❌ Ошибка чтения файла: %v", err)
	}

	var jsonData map[string]interface{}
	if err := json.Unmarshal(byteValue, &jsonData); err != nil {
		log.Fatalf("❌ Ошибка парсинга JSON: %v", err)
	}

	// results → []interface{} → []interface{} внутри
	rawEvents := jsonData["results"].([]interface{})[0].([]interface{})

	for _, item := range rawEvents {
		entryBytes, _ := json.Marshal(item)
		var ev rawEntry
		if err := json.Unmarshal(entryBytes, &ev); err != nil {
			continue
		}

		if ev.Type == "EV" {
			game := models.LiveGame{
				ID:         ev.ID,
				Name:       ev.NA,
				League:     ev.CT,
				Score:      ev.SS,
				MatchTime:  ev.TM,
				UpdateTime: ev.TU,
			}

			if _, err := db.DB.NamedExec(`INSERT INTO live_games (id, name, league, score, match_time, update_time)
				VALUES (:id, :name, :league, :score, :match_time, :update_time)
				ON CONFLICT (id) DO UPDATE SET
				name = EXCLUDED.name,
				league = EXCLUDED.league,
				score = EXCLUDED.score,
				match_time = EXCLUDED.match_time,
				update_time = EXCLUDED.update_time`, game); err != nil {
				log.Printf("⚠️ Не удалось вставить ID %s: %v", ev.ID, err)
			} else {
				log.Printf("✅ Добавлен матч: %s [%s]", ev.NA, ev.ID)
			}
		}
	}
}
