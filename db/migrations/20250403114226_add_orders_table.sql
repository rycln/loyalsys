-- +goose Up
-- +goose StatementBegin
CREATE TABLE orders (
    id BIGINT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    number VARCHAR(255) UNIQUE NOT NULL, 
    user_id BIGINT NOT NULL REFERENCES users(id), 
    status VARCHAR(255) DEFAULT 'NEW',
    accrual INT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP 
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS orders;
-- +goose StatementEnd
