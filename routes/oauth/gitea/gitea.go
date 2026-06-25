package gitea

import (
	config "libraone/config/generated"
	db "libraone/db/generated"
	"libraone/internal/lib/trail"
	"libraone/internal/model"
	"libraone/internal/session"
	"net/http"

	"github.com/xySaad/z01auth"
)

type Gitea struct{}

func (Gitea) Entry(z01authConfig z01auth.Config) trail.NoDepHandler {
	return func(c *trail.Context, _ trail.Null) (trail.Success, *trail.Error) {
		//TOOD: add state
		headers := http.Header{
			"Location": {z01authConfig.AuthCodeURL("")},
		}

		return c.Success(http.StatusMovedPermanently, headers, nil)
	}
}

func (Gitea) Callback(config config.Config, z01authConfig z01auth.Config, queries *db.Queries) trail.NoDepHandler {
	return func(c *trail.Context, _ trail.Null) (trail.Success, *trail.Error) {
		giteaToken, err := z01authConfig.Exchange(c, c.Query("code"))
		if err != nil {
			return model.ErrOAuthCodeExchange(err)
		}
		candidateId, err := z01authConfig.FetchCandidateId(giteaToken)
		if err != nil {
			return model.ErrFetchCandidateId(err)
		}
		cookieToken, expiresAt, err := session.New(c, queries, int64(candidateId))
		if err != nil {
			return model.ErrCreateSession(err)
		}
		err = queries.InsertGiteaToken(c, db.InsertGiteaTokenParams{
			CandidateID:  int64(candidateId),
			AccessToken:  giteaToken.AccessToken,
			TokenType:    giteaToken.TokenType,
			RefreshToken: giteaToken.RefreshToken,
			Expiry:       giteaToken.Expiry,
			ExpiresIn:    giteaToken.ExpiresIn,
		})
		if err != nil {
			return model.ErrInsertGiteaToken(err)
		}

		cookie := http.Cookie{
			Name:     session.CookieName,
			Value:    cookieToken,
			Expires:  expiresAt,
			Path:     "/",
			Domain:   "",
			Secure:   false,
			HttpOnly: true,
		}

		headers := http.Header{
			"Location":   {config.CallbackRedirectURL},
			"Set-Cookie": {cookie.String()},
		}

		return c.Success(http.StatusSeeOther, headers, nil)
	}
}
