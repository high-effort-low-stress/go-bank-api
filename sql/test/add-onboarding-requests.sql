INSERT INTO onboarding.onboarding_requests (
    public_id,
    full_name,
    email,
    document_number,
    verification_token_hash,
    token_expires_at
) VALUES (
    '01ARZ3NDEKTSV4RRFFQ69G5FAV',
    'Jane Doe',
    'jane.doe@example.com',
    '12345678901',
    '9f86d081884c7d659a2feaa0c55ad015a3bf4f1b2b0b822cd15d6c15b0f00a08',
    NOW() + INTERVAL '1 hour'
);
