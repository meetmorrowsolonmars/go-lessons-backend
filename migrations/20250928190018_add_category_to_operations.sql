-- +goose Up
-- +goose StatementBegin

CREATE TABLE IF NOT EXISTS categories
(
    id   BIGSERIAL NOT NULL,
    name VARCHAR   NOT NULL,

    PRIMARY KEY (id)
);

ALTER TABLE IF EXISTS operations
    ADD COLUMN IF NOT EXISTS category_id BIGINT,
    ADD CONSTRAINT operations_category_id_fkey FOREIGN KEY (category_id) REFERENCES categories (id);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

ALTER TABLE IF EXISTS operations
    DROP COLUMN IF EXISTS category_id;

DROP TABLE IF EXISTS categories;

-- +goose StatementEnd
