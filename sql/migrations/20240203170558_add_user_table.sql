-- +goose Up
-- +goose StatementBegin
CREATE TABLE "user" (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    email VARCHAR(255) NOT NULL
);
CREATE TRIGGER update_user_modtime
BEFORE UPDATE ON "user"
FOR EACH ROW EXECUTE FUNCTION update_modified_column();
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TRIGGER update_user_modtime ON "user";
DROP TABLE "user";
-- +goose StatementEnd
