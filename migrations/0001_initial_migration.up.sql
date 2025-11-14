CREATE EXTENSION IF NOT EXISTS pgcrypto;

CREATE TABLE IF NOT EXISTS teams(
    team_id UUID PRIMARY KEY DEFAULT gen_random_uuid() NOT NULL,
    team_name VARCHAR(255) UNIQUE NOT NULL,
    updated_at TIMESTAMP DEFAULT NOW() NOT NULL,
    created_at TIMESTAMP DEFAULT NOW() NOT NULL
);

CREATE INDEX IF NOT EXISTS teams_team_name_idx ON teams(team_name);

CREATE TABLE IF NOT EXISTS team_members (
    user_id UUID PRIMARY KEY DEFAULT gen_random_uuid() NOT NULL,
    team_id UUID REFERENCES teams(team_id) ON DELETE CASCADE,
    username VARCHAR(255) UNIQUE NOT NULL,
    is_active BOOLEAN DEFAULT false NOT NULL
);

CREATE TYPE pull_request_status AS ENUM ('open', 'closed');

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