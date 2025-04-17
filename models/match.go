package models

type Match struct {
	MatchID   string `db:"match_id"`
	League    string `db:"league"`
	MatchName string `db:"match_name"`
	Score     string `db:"score"`
	GameTime  string `db:"game_time"`
	Half      string `db:"half"`
}
