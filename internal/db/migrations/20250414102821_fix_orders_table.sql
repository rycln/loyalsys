-- +goose Up
-- +goose StatementBegin
ALTER TABLE orders ALTER COLUMN accrual TYPE DECIMAL(10, 2);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE orders ALTER COLUMN accrual TYPE INT;
-- +goose StatementEnd
