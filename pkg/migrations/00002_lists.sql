-- +goose Up
-- +goose StatementBegin
CREATE TABLE lists (
    id SERIAL PRIMARY KEY,
    title TEXT NOT NULL,
    description TEXT
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE lists;
-- +goose StatementEnd
