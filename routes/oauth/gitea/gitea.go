package gitea

import (
	"fmt"
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
		if err != nil {
			fmt.Println("error:", err)
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
		})

		if err != nil {
			fmt.Println("db error:", err)
			return
		}

		token, expiresAt, err := session.New(c, queries, int64(candidate.GiteaID))
		if err != nil {
			fmt.Println("session error:", err)
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
