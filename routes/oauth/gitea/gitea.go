package gitea

import (
	"context"
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
		bgCtx := context.Background()
		giteaToken, err := z01authConfig.Exchange(bgCtx, c.Query("code"))
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "OAuth failure"})
			return
		}
		candidateId, err := z01authConfig.FetchCandidateId(giteaToken)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "UserInfo failure"})
			return
		}
		cookieToken, expiresAt, err := session.New(c, queries, int64(candidateId))
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Session error"})
			return
		}
		err = queries.InsertGiteaToken(bgCtx, db.InsertGiteaTokenParams{
			CandidateID:  int64(candidateId),
			AccessToken:  giteaToken.AccessToken,
			TokenType:    giteaToken.TokenType,
			RefreshToken: giteaToken.RefreshToken,
			Expiry:       giteaToken.Expiry,
			ExpiresIn:    giteaToken.ExpiresIn,
		})
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
			return
		}

		c.SetSameSite(http.SameSiteLaxMode)
		c.SetCookie(
			session.CookieName,
			cookieToken,
			int(time.Until(expiresAt).Seconds()),
			"/",
			"",
			false, // secure: set true once serving over HTTPS
			true,
		)

		c.Redirect(http.StatusSeeOther, config.CallbackRedirectURL)
	}
}
