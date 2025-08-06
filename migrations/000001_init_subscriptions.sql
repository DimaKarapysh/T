-- +goose Up
-- +goose StatementBegin
CREATE TABLE subscriptions
(
    id           UUID PRIMARY KEY,
    service_name TEXT        NOT NULL,
    price        INT         NOT NULL,
    user_id      UUID        NOT NULL,
    start_date   DATE        NOT NULL,
    end_date     DATE,
    created_at   TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at   TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at   TIMESTAMPTZ DEFAULT NULL
);

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TABLE "subscriptions";
-- +goose StatementEnd