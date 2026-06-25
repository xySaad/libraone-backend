ALTER TABLE candidates
RENAME TO DEPRECATED_candidates;

CREATE TABLE
    candidates (
        id INTEGER PRIMARY KEY NOT NULL,
        created_at TIMESTAMP NOT NULL
    );

INSERT INTO
    candidates (id, created_at)
SELECT
    id,
    created_at
FROM
    DEPRECATED_candidates;