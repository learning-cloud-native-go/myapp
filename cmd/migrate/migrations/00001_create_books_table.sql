-- +goose Up
-- SQL in this section is executed when the migration is applied.
CREATE TABLE IF NOT EXISTS books
(
    id             UUID PRIMARY KEY,
    title          TEXT NOT NULL,
    author         TEXT NOT NULL,
    published_date DATE NOT NULL,
    image_url      TEXT NULL,
    description    TEXT NULL,
    created_at     TIMESTAMPTZ NOT NULL,
    updated_at     TIMESTAMPTZ NOT NULL
);

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
DROP TABLE IF EXISTS books;