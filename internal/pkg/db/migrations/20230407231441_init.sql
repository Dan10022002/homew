-- +goose Up
-- +goose StatementBegin
CREATE TABLE users (
    id serial PRIMARY KEY NOT NULL,
    name varchar(50) NOT NULL,
    surname varchar(50) NOT NULL,
    age int NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT now() NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE NULL
);

CREATE TABLE tickets (
    id serial PRIMARY KEY NOT NULL,
    user_id int NOT NULL,
    cost int NOT NULL DEFAULT 0,
    place int NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT now() NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE users;
DROP TABLE tickets;
-- +goose StatementEnd
