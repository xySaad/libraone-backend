CREATE TABLE
    sessions (
        token_hash TEXT PRIMARY KEY,
        candidate_id INTEGER NOT NULL REFERENCES candidates (id) ON DELETE CASCADE,
        expires_at DATETIME NOT NULL,
        created_at DATETIME NOT NULL
    );

CREATE INDEX idx_sessions_candidate_id ON sessions (candidate_id);