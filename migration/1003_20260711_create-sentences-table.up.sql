CREATE TABLE IF NOT EXISTS sentences (
    id UUID PRIMARY KEY DEFAULT uuidv7(),
    card_id UUID NOT NULL,
    position INTEGER NOT NULL,
    text VARCHAR(255) NOT NULL,
    reading VARCHAR(255) NOT NULL DEFAULT '',
    is_target BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_sentences_card_id ON sentences(card_id);

ALTER TABLE sentences
ADD CONSTRAINT fk_sentences_card_id_cards_id FOREIGN KEY(card_id) REFERENCES cards(id);
