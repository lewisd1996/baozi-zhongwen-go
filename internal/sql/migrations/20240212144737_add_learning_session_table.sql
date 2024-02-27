-- +goose Up
-- +goose StatementBegin
CREATE TABLE learning_session (
  id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
  created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
  user_id uuid NOT NULL,
  deck_id uuid NOT NULL,
  ended_at TIMESTAMP WITH TIME ZONE,
  review_count INT NOT NULL DEFAULT 0,
  FOREIGN KEY (user_id) REFERENCES "user" (id) ON DELETE CASCADE,
  FOREIGN KEY (deck_id) REFERENCES deck (id) ON DELETE CASCADE
);

CREATE TRIGGER update_learning_session_modtime
BEFORE UPDATE ON "learning_session" 
FOR EACH ROW EXECUTE FUNCTION update_modified_column();
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TRIGGER update_learning_session_modtime ON "learning_session";
DROP TABLE learning_session;
-- +goose StatementEnd
