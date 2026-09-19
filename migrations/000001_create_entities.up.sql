CREATE TABLE IF NOT EXISTS entities (
    id          VARCHAR(32)    PRIMARY KEY,
    name        VARCHAR(255)   NOT NULL,
    kind        VARCHAR(50)    NOT NULL,
    status      VARCHAR(50)    NOT NULL,
    latitude    DECIMAL(10,6)  NOT NULL,
    longitude   DECIMAL(10,6)  NOT NULL,
    description TEXT,
    created_at  TIMESTAMPTZ    NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ    NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS entity_attributes (
    id         BIGSERIAL     PRIMARY KEY,
    entity_id  VARCHAR(32)   NOT NULL REFERENCES entities(id) ON DELETE CASCADE,
    label      VARCHAR(100)  NOT NULL,
    value      VARCHAR(255)  NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_entities_kind ON entities(kind);
CREATE INDEX IF NOT EXISTS idx_entities_status ON entities(status);
CREATE INDEX IF NOT EXISTS idx_entity_attributes_entity_id ON entity_attributes(entity_id);
