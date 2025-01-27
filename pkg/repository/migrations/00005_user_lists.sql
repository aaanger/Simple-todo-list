-- +goose Up
-- +goose StatementBegin
CREATE TABLE user_lists (
    id SERIAL PRIMARY KEY,
    user_id INT REFERENCES users (id) ON DELETE CASCADE,
    list_id INT REFERENCES lists (id) ON DELETE CASCADE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE user_lists;
-- +goose StatementEnd
