-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
-- +goose StatementEnd
CREATE TABLE IF NOT EXISTS administrators (
    user_id INT PRIMARY KEY,
    username VARCHAR NOT NULL DEFAULT ('')
);

CREATE TABLE IF NOT EXISTS chats (
    id SERIAL PRIMARY KEY,
    chat_id BIGINT NOT NULL,
    chat_name VARCHAR NOT NULL
);

CREATE TABLE IF NOT EXISTS users (
    user_id INT PRIMARY KEY,
    username VARCHAR NOT NULL DEFAULT (''),
    chat_ids BIGINT[] NOT NULL
);

CREATE TABLE IF NOT EXISTS subscriptions (
    id SERIAL PRIMARY KEY,
    chat_id BIGINT NOT NULL,
    chat_name VARCHAR NOT NULL,
    user_id INT NOT NULL,
    feature INT NOT NULL,
    store jsonb NOT NULL,
    created_at TIMESTAMP NOT NULL,
    last_run_at TIMESTAMP NOT NULL
);

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
DROP TABLE IF EXISTS administrators;
DROP TABLE IF EXISTS chats;
DROP TABLE IF EXISTS users CASCADE;
DROP TABLE IF EXISTS subscriptions;