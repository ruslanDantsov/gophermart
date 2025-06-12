-- +goose Up
CREATE TABLE "order" (
    id       uuid NOT NULL PRIMARY KEY,
    number   varchar(256) NOT NULL UNIQUE,
    status   varchar(32) NOT NULL,
    accrual  numeric(12, 4) DEFAULT 0,
    created_at TIMESTAMP,
    user_id uuid NOT NULL REFERENCES "user_data"(id)
);

-- +goose Down
DROP TABLE IF EXISTS "order";
