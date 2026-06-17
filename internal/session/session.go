package session

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"time"

	db "libraone/db/generated"
)

const (
	CookieName = "session_token"
	ttl        = 7 * 24 * time.Hour
)

// New generates a random session token, stores its hash in the DB,
// and returns the raw token (to be set in a cookie) plus its expiry.
func New(ctx context.Context, queries *db.Queries, candidateID int64) (token string, expiresAt time.Time, err error) {
	raw := make([]byte, 32)
	if _, err = rand.Read(raw); err != nil {
		return "", time.Time{}, err
	}
	token = base64.RawURLEncoding.EncodeToString(raw)
	now := time.Now()
	expiresAt = now.Add(ttl)

	err = queries.CreateSession(ctx, db.CreateSessionParams{
		TokenHash:   Hash(token),
		CandidateID: candidateID,
		ExpiresAt:   expiresAt,
		CreatedAt:   now,
	})
	return token, expiresAt, err
}

// Hash returns the hex-encoded SHA-256 hash of a raw session token.
func Hash(token string) string {
	sum := sha256.Sum256([]byte(token))
	return hex.EncodeToString(sum[:])
}
