package candidate

import (
	db "libraone/db/generated"
	"libraone/internal/middlewares"
	"libraone/internal/services/profile"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Candidate struct {
	ProfileService profile.ProfileService
}

func (cmp *Candidate) Candidate(queries *db.Queries) middlewares.HandlerFunc[db.Candidate] {
	return func(c *gin.Context, selfCandidate *db.Candidate) {
		login := c.Param("login")
		if login == "" {
			c.JSON(http.StatusOK, selfCandidate)
			return
		}

		candidate, err := queries.GetCandidateByGraphqlLogin(c, selfCandidate.GraphqlLogin)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "Not Found"})
			return
		}
		c.JSON(http.StatusOK, candidate)
	}
}
