-- name: GetGiteaTokenByCandidateId :one
SELECT
    *
FROM
    gitea_tokens
WHERE
    candidate_id = ?