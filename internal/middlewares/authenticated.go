package middlewares

import (
	"database/sql"
	"net/http"
	"time"

	db "libraone/db/generated"
	"libraone/internal/dto"
	"libraone/internal/session"

	"github.com/gin-gonic/gin"
	"github.com/xySaad/z01auth"
	"golang.org/x/oauth2"
)

// EnsureAuthenticated reads the session cookie, loads the candidate
// it belongs to, and only lets the request through if the session is
// valid (not expired) and the candidate has platform_access = true.
func EnsureAuthenticated(queries *db.Queries, z01authConfig z01auth.Config) MiddlewareFunc[dto.Candidate] {
	return func(c *gin.Context) *dto.Candidate {
		cookieToken, err := c.Cookie(session.CookieName)
		if err != nil || cookieToken == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return nil
		}
		dbGiteaToken, err := queries.GetGiteaTokenBySessionToken(c, db.GetGiteaTokenBySessionTokenParams{
			TokenHash: session.Hash(cookieToken),
			ExpiresAt: time.Now(),
		})
		if err == sql.ErrNoRows {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return nil
		}
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
			return nil
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
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "UserInfo failure"})
			return nil
		}

		if !candidate.PlatformAccess {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Forbidden"})
			return nil
		}
		return dto.CandidateFromZ01auth(candidate)
	}
}
