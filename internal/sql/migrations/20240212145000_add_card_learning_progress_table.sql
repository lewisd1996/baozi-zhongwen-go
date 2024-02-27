-- +goose Up
-- +goose StatementBegin
CREATE TABLE card_learning_progress (
  id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
  created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
  session_id uuid NOT NULL,
  card_id uuid NOT NULL,
  user_id uuid NOT NULL,
  review_count INTEGER NOT NULL DEFAULT 0,
  success_count INTEGER NOT NULL DEFAULT 0,
  last_reviewed_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY (session_id) REFERENCES learning_session (id) ON DELETE CASCADE,
  FOREIGN KEY (card_id) REFERENCES card (id) ON DELETE CASCADE,
  FOREIGN KEY (user_id) REFERENCES "user" (id) ON DELETE CASCADE
);


CREATE TRIGGER update_card_learning_progress_modtime
BEFORE UPDATE ON "card_learning_progress" 
FOR EACH ROW EXECUTE FUNCTION update_modified_column();
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TRIGGER update_card_learning_progress_modtime ON "card_learning_progress";
DROP TABLE card_learning_progress;
-- +goose StatementEnd
