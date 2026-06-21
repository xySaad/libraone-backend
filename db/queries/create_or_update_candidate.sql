-- name: CreateOrUpdateCandidate :exec
INSERT INTO
    candidates (
        id,
        role,
        avatar_url,
        description,
        gitea_login,
        graphql_login,
        graphql_id,
        campus,
        platform_access
    )
VALUES
    (?, ?, ?, ?, ?, ?, ?, ?, ?) ON CONFLICT (id) DO
UPDATE
SET
    role = excluded.role,
    avatar_url = excluded.avatar_url,
    description = excluded.description,
    gitea_login = excluded.gitea_login,
    graphql_login = excluded.graphql_login,
    graphql_id = excluded.graphql_id,
    campus = excluded.campus,
    platform_access = excluded.platform_access;