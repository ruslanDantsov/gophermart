-- +goose Up
CREATE TABLE user_data (
    id     uuid NOT NULL PRIMARY KEY,
    login   varchar(64) NOT NULL UNIQUE,
    password  varchar(64) NOT NULL,
    created_at TIMESTAMP
);

CREATE INDEX user_login_index on user_data USING btree(login);

-- +goose Down
DROP TABLE IF EXISTS user_data;
DROP INDEX IF Exists user_login_index;