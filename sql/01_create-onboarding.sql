create schema onboarding;

CREATE TABLE onboarding.onboarding_requests (
    id BIGINT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    public_id VARCHAR(26) UNIQUE NOT NULL,
    full_name VARCHAR(255) NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    document_number VARCHAR(11) UNIQUE NOT NULL,
    verification_token_hash VARCHAR(255) NOT NULL UNIQUE,
    token_expires_at TIMESTAMPTZ NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'PENDING',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
); 