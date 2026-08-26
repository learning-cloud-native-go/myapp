-- +goose Up
-- SQL in this section is executed when the migration is applied.
CREATE TABLE IF NOT EXISTS books
(
    created_at     TIMESTAMPTZ NOT NULL,
    updated_at     TIMESTAMPTZ NOT NULL,
    id             UUID        NOT NULL,
    published_date DATE        NOT NULL,
    status         SMALLINT    NOT NULL,
    title          TEXT        NOT NULL,
    description    TEXT,
    image_url      TEXT,
    PRIMARY KEY (id)
);

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
DROP TABLE IF EXISTS books;