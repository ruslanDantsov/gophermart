-- +goose Up
CREATE TABLE withdraw (
    id       uuid NOT NULL PRIMARY KEY,
    sum  numeric(12, 4) DEFAULT 0,
    created_at TIMESTAMP,
    order_id uuid NOT NULL REFERENCES "order"(id)
);

-- +goose Down
DROP TABLE IF EXISTS withdraw;
