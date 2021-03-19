-- +goose Up
-- SQL in this section is executed when the migration is applied.
CREATE TABLE  IF NOT EXISTS signals (
     signal_id  serial PRIMARY KEY,
     device_id  text
);
-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
DROP TABLE IF EXISTS signals;