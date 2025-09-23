-- +goose Up
-- +goose StatementBegin

CREATE TABLE IF NOT EXISTS users
(
    id          BIGSERIAL   NOT NULL,
    email       VARCHAR     NOT NULL,
    password    VARCHAR     NOT NULL,
    full_name   VARCHAR     NOT NULL,
    create_time TIMESTAMPTZ NOT NULL DEFAULT now(),

    PRIMARY KEY (id)
);

CREATE UNIQUE INDEX IF NOT EXISTS users_email_idx ON users (email);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP TABLE IF EXISTS users;

-- +goose StatementEnd
