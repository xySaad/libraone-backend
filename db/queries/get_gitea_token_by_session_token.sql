-- name: GetGiteaTokenBySessionToken :one
SELECT
    gitea_tokens.*
FROM
    sessions
    JOIN gitea_tokens ON gitea_tokens.candidate_id = sessions.candidate_id
WHERE
    sessions.token_hash = ?
    AND sessions.expires_at > ?;