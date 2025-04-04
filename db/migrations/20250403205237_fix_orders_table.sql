-- +goose Up
-- +goose StatementBegin
ALTER TABLE orders
ALTER COLUMN accrual SET DEFAULT 0;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE orders
ALTER COLUMN accrual DROP DEFAULT;
-- +goose StatementEnd

