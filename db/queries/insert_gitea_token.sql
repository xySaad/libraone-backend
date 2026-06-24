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
    (?, ?, ?, ?, ?, ?);