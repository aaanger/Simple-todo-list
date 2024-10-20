-- +goose Up
-- +goose StatementBegin
CREATE TABLE list_items (
    id SERIAL PRIMARY KEY,
    list_id INT REFERENCES lists (id) ON DELETE CASCADE NOT NULL,
    item_id INT REFERENCES items (id) ON DELETE CASCADE NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE list_items;
-- +goose StatementEnd
