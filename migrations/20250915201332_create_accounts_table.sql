-- +goose Up
-- +goose StatementBegin

CREATE TABLE IF NOT EXISTS accounts
(
    id          BIGSERIAL   NOT NULL,
    user_id     BIGINT      NOT NULL,
    title       VARCHAR     NOT NULL,
    is_default  BOOLEAN     NOT NULL DEFAULT false,
    create_time TIMESTAMPTZ NOT NULL DEFAULT now(),

    PRIMARY KEY (id),
    FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP TABLE IF EXISTS accounts;

-- +goose StatementEnd
