CREATE EXTENSION IF NOT EXISTS pgcrypto;

CREATE TABLE IF NOT EXISTS team_members (
    user_id UUID PRIMARY KEY DEFAULT gen_random_uuid() NOT NULL,
    team_id UUID REFERENCES teams(team_id) ON DELETE CASCADE,
    username VARCHAR(255) UNIQUE NOT NULL,
    is_active BOOLEAN DEFAULT false NOT NULL
);