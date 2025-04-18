-- +goose Up
-- +goose StatementBegin
CREATE TABLE withdrawals (
    id BIGINT PRIMARY KEY GENERATED ALWAYS AS IDENTITY, 
    "order" VARCHAR(255) UNIQUE NOT NULL,
    user_id BIGINT NOT NULL REFERENCES users(id),  
    sum DECIMAL(10, 2) NOT NULL,
    processed_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP  
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS withdrawals;
-- +goose StatementEnd
