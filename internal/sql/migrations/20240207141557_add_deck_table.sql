-- +goose Up
-- +goose StatementBegin
CREATE TABLE deck (
  id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
  created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
  name VARCHAR(255) NOT NULL,
  description VARCHAR(255) NOT NULL,
  owner_id uuid NOT NULL
);

ALTER TABLE deck
  ADD CONSTRAINT deck_owner_id_fkey FOREIGN KEY (owner_id) REFERENCES "user" (id);

CREATE TRIGGER update_deck_modtime
BEFORE UPDATE ON "deck"
FOR EACH ROW EXECUTE FUNCTION update_modified_column();
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TRIGGER update_deck_modtime ON "deck";
DROP TABLE deck;
-- +goose StatementEnd
