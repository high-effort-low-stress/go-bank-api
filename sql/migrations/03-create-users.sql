create schema "user";

CREATE TYPE "user".user_status AS ENUM ('ACTIVE', 'INACTIVE', 'BLOCKED');

CREATE TABLE "user".users (
    id BIGINT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    public_id VARCHAR(26) NOT NULL UNIQUE,
    full_name VARCHAR(255) NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    document_number VARCHAR(11) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    status "user".user_status NOT NULL DEFAULT 'ACTIVE',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deactivated_at TIMESTAMPTZ
);

CREATE TABLE "user".accounts (
    id BIGINT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    public_id VARCHAR(26) NOT NULL UNIQUE, 
    user_id BIGINT NOT NULL,
    account_number VARCHAR(10) UNIQUE NOT NULL,
    agency_number VARCHAR(6) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    FOREIGN KEY (user_id) REFERENCES "user".users(id) ON DELETE CASCADE
);
