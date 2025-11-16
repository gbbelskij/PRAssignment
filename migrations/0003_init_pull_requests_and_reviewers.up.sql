CREATE EXTENSION IF NOT EXISTS pgcrypto;

CREATE TYPE pull_request_status AS ENUM ('OPEN', 'MERGED');

CREATE TABLE IF NOT EXISTS pull_requests(
    pull_request_id UUID PRIMARY KEY DEFAULT gen_random_uuid() NOT NULL,
    pull_request_name VARCHAR(255) NOT NULL,
    author_id UUID REFERENCES team_members(user_id) ON DELETE CASCADE,
    status pull_request_status NOT NULL,
    updated_at TIMESTAMP DEFAULT NOW() NOT NULL,
    merged_at TIMESTAMP DEFAULT NOW() NOT NULL,
    created_at TIMESTAMP DEFAULT NOW() NOT NULL
);

CREATE TABLE IF NOT EXISTS pull_request_reviewers (
    pull_request_id UUID REFERENCES pull_requests(pull_request_id) ON DELETE CASCADE,
    user_id UUID REFERENCES team_members(user_id) ON DELETE CASCADE,
    PRIMARY KEY (pull_request_id, user_id)
);

CREATE INDEX IF NOT EXISTS pull_request_reviewers_user_id ON pull_request_reviewers(user_id);