package importer

import (
	"encoding/json"	
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/jmoiron/sqlx"
)

var DB *sqlx.DB // передаём из main

type Event struct {
	Type string `json:"type"`
	ID   string `json:"ID"`
	NA   string `json:"NA"` // Match name
	CT   string `json:"CT"` // League name
	CL   string `json:"CL"` // League ID (optional)
	SS   string `json:"SS"` // Score
	TM   string `json:"TM"` // Game time
	GO   string `json:"GO"` // Half
}

type AllData struct {
	Results [][]Event `json:"results"`
}

func ImportAllLiveGames(db *sqlx.DB, jsonPath string) {
	DB = db // сохраняем ссылку

	fullPath := filepath.Join(jsonPath)
	content, err := os.ReadFile(fullPath)
	if err != nil {
		log.Fatalf("❌ Ошибка чтения файла %s: %v", fullPath, err)
	}

	var data AllData
	if err := json.Unmarshal(content, &data); err != nil {
		log.Fatalf("❌ Ошибка парсинга JSON: %v", err)
	}

	count := 0
	for _, group := range data.Results {
		for _, ev := range group {
			if ev.Type != "EV" {
				continue
			}

			home, away := parseTeams(ev.NA)
			if home == "" || away == "" {
				log.Printf("⚠️ Не удалось распарсить команды из: %s", ev.NA)
				continue
			}

			leagueID, err := getOrCreateLeague(ev.CT)
			if err != nil {
				continue
			}

			homeID, err := getOrCreateTeam(home)
			if err != nil {
				continue
			}
			awayID, err := getOrCreateTeam(away)
			if err != nil {
				continue
			}

			matchID, err := getOrCreateMatch(ev.ID, leagueID, homeID, awayID)
			if err != nil {
				continue
			}

			_ = updateMatchStatus(matchID, ev.TM, ev.GO)
			_ = updateMatchScore(matchID, ev.SS)

			count++
		}
	}

	log.Printf("✅ Импорт завершён. Добавлено матчей: %d", count)
}

func parseTeams(name string) (string, string) {
	delims := []string{" v ", " vs ", " VS ", " V "}
	for _, delim := range delims {
		if strings.Contains(name, delim) {
			parts := strings.SplitN(name, delim, 2)
			if len(parts) == 2 {
				return strings.TrimSpace(parts[0]), strings.TrimSpace(parts[1])
			}
		}
	}
	return "", ""
}

func getOrCreateLeague(name string) (int, error) {
	var id int
	err := DB.Get(&id, `SELECT id FROM leagues WHERE name = $1`, name)
	if err == nil {
		return id, nil
	}

	err = DB.QueryRowx(`INSERT INTO leagues (name) VALUES ($1) RETURNING id`, name).Scan(&id)
	if err != nil {
		log.Printf("❌ Лига '%s' не создана: %v", name, err)
	}
	return id, err
}

func getOrCreateTeam(name string) (int, error) {
	var id int
	err := DB.Get(&id, `SELECT id FROM teams WHERE name = $1`, name)
	if err == nil {
		return id, nil
	}

	err = DB.QueryRowx(`INSERT INTO teams (name) VALUES ($1) RETURNING id`, name).Scan(&id)
	if err != nil {
		log.Printf("❌ Команда '%s' не создана: %v", name, err)
	}
	return id, err
}

func getOrCreateMatch(matchUID string, leagueID, homeID, awayID int) (int, error) {
	var id int
	err := DB.Get(&id, `SELECT id FROM matches WHERE match_id = $1`, matchUID)
	if err == nil {
		return id, nil
	}

	query := `
		INSERT INTO matches (match_id, league_id, home_team_id, away_team_id)
		VALUES ($1, $2, $3, $4)
		RETURNING id
	`
	err = DB.QueryRowx(query, matchUID, leagueID, homeID, awayID).Scan(&id)
	if err != nil {
		log.Printf("❌ Матч '%s' не создан: %v", matchUID, err)
	}
	return id, err
}

func updateMatchStatus(matchID int, gameTime, half string) error {
	_, err := DB.Exec(`
		INSERT INTO match_status (match_id, game_time, half)
		VALUES ($1, $2, $3)
		ON CONFLICT (match_id) DO UPDATE
		SET game_time = EXCLUDED.game_time,
		    half = EXCLUDED.half,
		    updated_at = CURRENT_TIMESTAMP
	`, matchID, gameTime, half)
	if err != nil {
		log.Printf("❌ Статус матча %d не обновлён: %v", matchID, err)
	}
	return err
}

func updateMatchScore(matchID int, score string) error {
	_, err := DB.Exec(`
		INSERT INTO match_results (match_id, score)
		VALUES ($1, $2)
		ON CONFLICT (match_id) DO UPDATE
		SET score = EXCLUDED.score,
		    updated_at = CURRENT_TIMESTAMP
	`, matchID, score)
	if err != nil {
		log.Printf("❌ Счёт матча %d не обновлён: %v", matchID, err)
	}
	return err
}
