package campus

import (
	db "libraone/db/generated"
	"libraone/internal/services/profile"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Campus struct {
	profileservice profile.ProfileService
}

func (cmp *Campus) Online(c *gin.Context, candidate *db.Candidate) {
	c.JSON(http.StatusOK, gin.H{
		"message": "not implemented",
		"login":   candidate.GiteaLogin,
	})
}
