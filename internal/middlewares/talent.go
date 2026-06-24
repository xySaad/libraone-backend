package middlewares

import (
	db "libraone/db/generated"
	"libraone/internal/dto"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/xySaad/z01auth"
)

func EnsureTalentRole(queries *db.Queries, z01authConfig z01auth.Config) MiddlewareFunc[dto.Candidate] {
	return func(c *gin.Context) *dto.Candidate {
		candidate := EnsureAuthenticated(queries, z01authConfig)(c)
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
