package gitea

import (
	config "libraone/config/generated"
	db "libraone/db/generated"
	"libraone/internal/session"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/xySaad/z01auth"
)

type Gitea struct{}

func (Gitea) Entry(z01authConfig z01auth.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		//TOOD: add state
		c.Redirect(http.StatusMovedPermanently, z01authConfig.AuthCodeURL(""))
	}
}

func (Gitea) Callback(config config.Config, z01authConfig z01auth.Config, queries *db.Queries) gin.HandlerFunc {
	return func(c *gin.Context) {
		candidate, err := z01authConfig.Callback(c.Query("code"))
		if err == z01auth.ErrMultipleUsers {
			c.AbortWithStatusJSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "OAuth failure"})
			return
		}

		err = queries.CreateOrUpdateCandidate(c, db.CreateOrUpdateCandidateParams{
			ID:             int64(candidate.GiteaID),
			GiteaLogin:     candidate.GiteaLogin,
			AvatarUrl:      candidate.AvatarURL,
			Description:    candidate.Description,
			Role:           string(candidate.Role),
			GraphqlLogin:   candidate.GraphqlLogin,
			Campus:         candidate.Campus,
			PlatformAccess: candidate.PlatformAccess,
			GraphqlID:      int64(candidate.GraphqlId),
		})
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
			return
		}

		token, expiresAt, err := session.New(c, queries, int64(candidate.GiteaID))
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Session error"})
			return
		}

		c.SetSameSite(http.SameSiteLaxMode)
		c.SetCookie(
			session.CookieName,
			token,
			int(time.Until(expiresAt).Seconds()),
			"/",
			"",
			false, // secure: set true once serving over HTTPS
			true,
		)

		c.Redirect(http.StatusSeeOther, config.CallbackRedirectURL)
	}
}
