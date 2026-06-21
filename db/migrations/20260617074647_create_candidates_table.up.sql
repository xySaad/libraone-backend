CREATE TABLE
    candidates (
        id INTEGER PRIMARY KEY NOT NULL,
        role TEXT NOT NULL,
        avatar_url TEXT NOT NULL,
        description TEXT NOT NULL,
        gitea_login TEXT NOT NULL UNIQUE,
        graphql_login TEXT NOT NULL UNIQUE,
        graphql_id INTEGER NOT NULL UNIQUE,
        campus TEXT NOT NULL,
        platform_access BOOlEAN NOT NULL,
        created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
    );