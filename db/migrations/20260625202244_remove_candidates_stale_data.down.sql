ALTER TABLE candidates
RENAME TO REVERTED_NEW_candidates;

ALTER TABLE DEPRECATED_candidates
RENAME TO candidates;