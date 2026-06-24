CREATE TABLE
    gitea_tokens (
        candidate_id INTEGER PRIMARY KEY NOT NULL REFERENCES candidates (id) ON DELETE CASCADE,
        access_token TEXT UNIQUE NOT NULL,
        token_type TEXT NOT NULL,
        refresh_token TEXT NOT NULL,
        expiry TIMESTAMP NOT NULL,
        expires_in INTEGER NOT NULL
    );