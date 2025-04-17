package handlers

import (
	"encoding/json"
	"net/http"
	"os"
	"path/filepath"

	"test2/models"

	"github.com/gin-gonic/gin"
)

type MatchDTO struct {
	MatchID   string `json:"match_id"`
	League    string `json:"league"`
	MatchName string `json:"match_name"`
	Score     string `json:"score"`
	GameTime  string `json:"game_time"`
	Half      string `json:"half"`
}

func GetLiveEvents(c *gin.Context) {
	path := filepath.Join("data", "allLiveGAmes.json")
	content, err := os.ReadFile(path)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "cannot read file"})
		return
	}

	var raw models.AllLiveGames
	if err := json.Unmarshal(content, &raw); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "cannot parse JSON"})
		return
	}

	var result []MatchDTO
	for _, group := range raw.Results {
		for _, item := range group {
			if item.Type == "EV" {
				result = append(result, MatchDTO{
					MatchID:   item.ID,
					League:    item.CT,
					MatchName: item.NA,
					Score:     item.SS,
					GameTime:  item.TM,
					Half:      item.GO,
				})
			}
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"matches": result,
		"count":   len(result),
	})
}
