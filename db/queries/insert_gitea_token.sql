-- name: InsertGiteaToken :exec
INSERT INTO
    gitea_tokens (
        candidate_id,
        access_token,
        token_type,
        refresh_token,
        expiry,
        expires_in
    )
VALUES
    (?, ?, ?, ?, ?, ?) ON CONFLICT (candidate_id) DO
UPDATE
SET
    access_token = excluded.access_token,
    token_type = excluded.token_type,
    refresh_token = excluded.refresh_token,
    expiry = excluded.expiry,
    expires_in = excluded.expires_in;