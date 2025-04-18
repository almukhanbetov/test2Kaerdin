-- +goose Up
CREATE TABLE IF NOT EXISTS live_games (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    league TEXT NOT NULL,
    score TEXT,
    match_time TEXT,
    update_time TEXT
);

-- +goose Down
DROP TABLE IF EXISTS live_games;
