CREATE TABLE category (
    id         BIGSERIAL    NOT NULL UNIQUE,
    guid       UUID         NOT NULL PRIMARY KEY DEFAULT gen_random_uuid(),
    name       VARCHAR(255) NOT NULL UNIQUE,
    created_at TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);