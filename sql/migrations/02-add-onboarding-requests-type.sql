CREATE TYPE onboarding.request_status AS ENUM ('PENDING', 'VERIFIED', 'COMPLETED', 'EXPIRED');

ALTER TABLE onboarding.onboarding_requests ALTER COLUMN status DROP DEFAULT;
ALTER TABLE onboarding.onboarding_requests ALTER COLUMN status SET DATA TYPE onboarding.request_status USING (status::onboarding.request_status);
ALTER TABLE onboarding.onboarding_requests ALTER COLUMN status SET DEFAULT 'PENDING';