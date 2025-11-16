CREATE EXTENSION IF NOT EXISTS pgcrypto;

CREATE TABLE IF NOT EXISTS teams(
    team_id UUID PRIMARY KEY DEFAULT gen_random_uuid() NOT NULL,
    team_name VARCHAR(255) UNIQUE NOT NULL,
    updated_at TIMESTAMP DEFAULT NOW() NOT NULL,
    created_at TIMESTAMP DEFAULT NOW() NOT NULL
);

CREATE INDEX IF NOT EXISTS teams_team_name_idx ON teams(team_name);