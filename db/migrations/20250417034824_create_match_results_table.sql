-- +goose Up
CREATE TABLE IF NOT EXISTS match_results (
    id SERIAL PRIMARY KEY,
    match_id INTEGER REFERENCES matches(id),
    score TEXT,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT unique_result_match UNIQUE (match_id)
);

-- +goose Down
DROP TABLE IF EXISTS match_results;
