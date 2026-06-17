package middlewares

import (
	db "libraone/db/generated"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/xySaad/z01auth"
)

func EnsureTalentRole(queries *db.Queries) MiddlewareFunc[db.Candidate] {
	return func(c *gin.Context) *db.Candidate {
		candidate := EnsureAuthenticated(queries)(c)
		if c.IsAborted() {
			return nil
		}

		if candidate.Role != string(z01auth.CandidateRole_TALENT) {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Forbidden"})
			return nil
		}

		return candidate
	}
}
