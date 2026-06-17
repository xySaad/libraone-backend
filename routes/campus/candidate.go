package campus

import (
	db "libraone/db/generated"
	"libraone/internal/middlewares"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (cmp *Campus) Candidate(queries *db.Queries) middlewares.HandlerFunc[db.Candidate] {
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
