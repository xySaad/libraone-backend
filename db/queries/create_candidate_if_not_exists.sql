-- name: CreateCandidateIfNotExists :exec
INSERT INTO
    candidates (id, created_at)
VALUES
    (?, ?) ON CONFLICT (id) DO NOTHING;