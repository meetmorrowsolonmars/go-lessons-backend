-- +goose Up
-- +goose StatementBegin

CREATE TABLE IF NOT EXISTS operations
(
    id          UUID           NOT NULL,
    user_id     BIGINT         NOT NULL,
    account_id  BIGINT         NOT NULL,
    type        VARCHAR        NOT NULL,
    amount      DECIMAL(15, 2) NOT NULL,
    description VARCHAR        NOT NULL,
    create_time TIMESTAMPTZ    NOT NULL DEFAULT now(),

    PRIMARY KEY (id),
    FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE,
    FOREIGN KEY (account_id) REFERENCES accounts (id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS operations_account_id_create_time_idx ON operations (account_id, create_time);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP TABLE IF EXISTS operations;

-- +goose StatementEnd
