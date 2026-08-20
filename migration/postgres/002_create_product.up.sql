CREATE TABLE product (
    id            BIGSERIAL    NOT NULL UNIQUE,
    guid          UUID         NOT NULL PRIMARY KEY DEFAULT gen_random_uuid(),
    name          VARCHAR(255) NOT NULL UNIQUE,
    description   TEXT,
    price         BIGINT       NOT NULL,
    category_guid UUID         NOT NULL REFERENCES category(guid) ON DELETE RESTRICT,
    created_at    TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at    TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);