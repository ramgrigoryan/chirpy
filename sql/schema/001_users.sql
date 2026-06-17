-- +goose Up
CREATE TABLE users(
    id varchar(36) primary key NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    email varchar(255) NOT NULL UNIQUE
);

-- +goose Down
DROP TABLE users;