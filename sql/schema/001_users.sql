-- +goose Up
CREATE TABLE users (
    id UUID PRIMARY KEY,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    name TEXT UNIQUE NOT NULL
);
CREATE TABLE feeds (
    id UUID PRIMARY KEY,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    name TEXT NOT NULL,
    url TEXT UNIQUE NOT NULL,
    user_id UUID REFERENCES users(id) ON DELETE CASCADE NOT NULL
);

-- +goose Down
DROP TABLE users;
DROP TABLE feeds;