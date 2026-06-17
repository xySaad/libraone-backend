-- name: GetCandidateByGraphqlLogin :one
SELECT
    candidates.*
FROM
    candidates
WHERE
    graphql_login = ?;