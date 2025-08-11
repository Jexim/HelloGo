-- +goose Up
CREATE TABLE IF NOT EXISTS hello (
  id SERIAL PRIMARY KEY,
  message TEXT NOT NULL
);

-- +goose Down
DROP TABLE IF EXISTS hello;

