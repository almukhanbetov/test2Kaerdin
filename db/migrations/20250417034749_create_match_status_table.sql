-- +goose Up
CREATE TABLE IF NOT EXISTS match_status (
    id SERIAL PRIMARY KEY,
    match_id INTEGER REFERENCES matches(id),
    game_time TEXT,
    half TEXT,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT unique_status_match UNIQUE (match_id)
);

-- +goose Down
DROP TABLE IF EXISTS match_status;
