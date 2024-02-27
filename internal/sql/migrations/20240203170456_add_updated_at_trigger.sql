-- +goose Up
-- +goose StatementBegin
CREATE OR REPLACE FUNCTION update_modified_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;   
END;
$$ language 'plpgsql';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP FUNCTION update_modified_column();
-- +goose StatementEnd
