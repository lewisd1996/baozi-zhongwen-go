-- +goose Up
-- +goose StatementBegin
ALTER TABLE card DROP CONSTRAINT card_deck_id_fkey;
ALTER TABLE card
ADD CONSTRAINT card_deck_id_fkey FOREIGN KEY (deck_id) REFERENCES deck (id) ON DELETE CASCADE;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE card DROP CONSTRAINT card_deck_id_fkey;
ALTER TABLE card
  ADD CONSTRAINT card_deck_id_fkey FOREIGN KEY (deck_id) REFERENCES deck (id);
-- +goose StatementEnd
