-- +goose Up
-- +goose StatementBegin
CREATE TABLE orders (
    id BIGINT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    user_id REFERENCES users(id) NOT NULL, 
    number VARCHAR(255) UNIQUE NOT NULL, 
    status VARCHAR(255) DEFAULT 'NEW',
    accrual INT,
    uploaded_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP, 
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS orders;
-- +goose StatementEnd
