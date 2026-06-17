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
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return nil
		}

		candidate, err := queries.GetCandidateBySessionToken(c, db.GetCandidateBySessionTokenParams{
			TokenHash: session.Hash(token),
			ExpiresAt: time.Now(),
		})
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return nil
		}

		if !candidate.PlatformAccess {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Forbidden"})
			return nil
		}
		return &candidate
	}
}
