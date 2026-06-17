package middlewares

import (
	"net/http"
	"time"

	db "libraone/db/generated"
	"libraone/internal/session"

	"github.com/gin-gonic/gin"
)

// EnsureAuthenticated reads the session cookie, loads the candidate
// it belongs to, and only lets the request through if the session is
// valid (not expired) and the candidate has platform_access = true.
func EnsureAuthenticated(queries *db.Queries) MiddlewareFunc[db.Candidate] {
	return func(c *gin.Context) *db.Candidate {
		token, err := c.Cookie(session.CookieName)
		if err != nil || token == "" {
			c.AbortWithStatus(http.StatusUnauthorized)
			return nil
		}

		candidate, err := queries.GetCandidateBySessionToken(c, db.GetCandidateBySessionTokenParams{
			TokenHash: session.Hash(token),
			ExpiresAt: time.Now(),
		})
		if err != nil {
			c.AbortWithStatus(http.StatusUnauthorized)
			return nil
		}

		if !candidate.PlatformAccess {
			c.AbortWithStatus(http.StatusForbidden)
			return nil
		}
		return &candidate
	}
}
