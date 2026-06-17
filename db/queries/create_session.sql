-- name: CreateSession :exec
INSERT INTO
    sessions (token_hash, candidate_id, expires_at, created_at)
VALUES
    (?, ?, ?, ?);