package models

type GameItem struct {
	Type string `json:"type"`
	ID   string `json:"ID"`
	NA   string `json:"NA,omitempty"`
	CT   string `json:"CT,omitempty"`
	SS   string `json:"SS,omitempty"`
	TM   string `json:"TM,omitempty"`
	GO   string `json:"GO,omitempty"`
	// Добавь другие поля при необходимости
}

type AllLiveGames struct {
	Success int           `json:"success"`
	Results [][]GameItem  `json:"results"`
}
