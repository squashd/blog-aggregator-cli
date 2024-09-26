-- +goose Up
CREATE TABLE IF NOT EXISTS posts (
    id UUID PRIMARY KEY,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    title VARCHAR(256) NOT NULL,
    url TEXT NOT NULL,
    description TEXT NOT NULL,
    published_at TIMESTAMP NOT NULL,
    feed_id UUID NOT NULL
);

-- +goose Down
DROP TABLE IF EXISTS posts;

