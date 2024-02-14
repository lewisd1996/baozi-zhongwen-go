-- +goose Up
-- +goose StatementBegin
CREATE OR REPLACE VIEW v_card_with_last_reviewed AS
SELECT
    c.id,
    c.created_at,
    c.updated_at,
    c.deck_id,
    c.content,
    c.translation,
    clp.last_reviewed_at
FROM
    card c
LEFT JOIN LATERAL (
    SELECT
        clp.last_reviewed_at
    FROM
        card_learning_progress clp
    WHERE
        clp.card_id = c.id
    ORDER BY
        clp.last_reviewed_at DESC
    LIMIT 1
) clp ON true
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP VIEW v_card_with_last_reviewed;
-- +goose StatementEnd
