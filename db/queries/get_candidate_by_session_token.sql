-- name: GetCandidateBySessionToken :one
SELECT
    candidates.*
FROM
    sessions
    JOIN candidates ON candidates.id = sessions.candidate_id
WHERE
    sessions.token_hash = ?
    AND sessions.expires_at > ?;