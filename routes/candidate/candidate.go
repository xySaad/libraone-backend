package candidate

import (
	db "libraone/db/generated"
	"libraone/internal/dto"
	"libraone/internal/middlewares"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/xySaad/z01auth"
	"golang.org/x/oauth2"
)

type Candidate struct {
	queries       *db.Queries
	z01authConfig z01auth.Config
}

func New(queries *db.Queries, z01authConfig z01auth.Config) *Candidate {
	return &Candidate{queries: queries, z01authConfig: z01authConfig}

}

func (cmp *Candidate) Candidate() middlewares.HandlerFunc[dto.Candidate] {
	return func(c *gin.Context, selfCandidate *dto.Candidate) {
		idParam := c.Param("id")
		if idParam == "" {
			c.JSON(http.StatusOK, selfCandidate)
			return
		}
		candidateId, err := strconv.Atoi(idParam)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid id"})
			return
		}

		dbGiteaToken, err := cmp.queries.GetGiteaTokenByCandidateId(c, int64(candidateId))
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
			return
		}

		token := &oauth2.Token{
			AccessToken:  dbGiteaToken.AccessToken,
			TokenType:    dbGiteaToken.TokenType,
			RefreshToken: dbGiteaToken.RefreshToken,
			Expiry:       dbGiteaToken.Expiry,
			ExpiresIn:    dbGiteaToken.ExpiresIn,
		}
		candidate, err := cmp.z01authConfig.FetchCandidate(token)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "UserInfo failure"})
			return
		}
		c.JSON(http.StatusOK, dto.CandidateFromZ01auth(candidate))
	}
}
