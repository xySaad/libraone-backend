package middlewares

import (
	"database/sql"
	"libraone/internal/lib/trail"
	"time"

	db "libraone/db/generated"
	"libraone/internal/dto"
	"libraone/internal/session"

	"github.com/xySaad/z01auth"
	"golang.org/x/oauth2"
)

// EnsureAuthenticated reads the session cookie, loads the candidate
// it belongs to, and only lets the request through if the session is
// valid (not expired) and the candidate has platform_access = true.
func EnsureAuthenticated(queries *db.Queries, z01authConfig z01auth.Config) trail.Middleware[dto.Candidate] {
	return func(c *trail.Context) (dto.Candidate, *trail.Error) {
		cookieToken, err := c.Request.Cookie(session.CookieName)
		if err != nil {
			return ErrMissingCookie(err)
		}

		dbGiteaToken, err := queries.GetGiteaTokenBySessionToken(c, db.GetGiteaTokenBySessionTokenParams{
			TokenHash: session.Hash(cookieToken.String()),
			ExpiresAt: time.Now(),
		})
		if err == sql.ErrNoRows {
			return ErrInvalidSession(err)
		}
		if err != nil {
			return ErrDatabase(err)
		}
		token := &oauth2.Token{
			AccessToken:  dbGiteaToken.AccessToken,
			TokenType:    dbGiteaToken.TokenType,
			RefreshToken: dbGiteaToken.RefreshToken,
			Expiry:       dbGiteaToken.Expiry,
			ExpiresIn:    dbGiteaToken.ExpiresIn,
		}
		candidate, err := z01authConfig.FetchCandidate(token)
		if err != nil {
			return ErrFetchCandidateInfo(err)
		}

		if !candidate.PlatformAccess {
			return ErrPlatformBanned(nil)
		}
		return dto.CandidateFromZ01auth(candidate), nil
	}
}
