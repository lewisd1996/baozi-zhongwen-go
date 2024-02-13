-- +goose Up
-- +goose StatementBegin
CREATE TABLE card (
  id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
  created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
  deck_id uuid NOT NULL,
  content VARCHAR(500) NOT NULL,
  translation VARCHAR(500) NOT NULL
);

ALTER TABLE card
  ADD CONSTRAINT card_deck_id_fkey FOREIGN KEY (deck_id) REFERENCES deck (id);

CREATE TRIGGER update_card_modtime
BEFORE UPDATE ON "card" 
FOR EACH ROW EXECUTE FUNCTION update_modified_column();
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TRIGGER update_card_modtime ON "card";
DROP TABLE card;
-- +goose StatementEnd
