package candidate

import (
	db "libraone/db/generated"
	"libraone/internal/dto"
	"libraone/internal/lib/trail"
	"libraone/internal/model"
	"net/http"
	"strconv"

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

func (cmp *Candidate) Candidate(c *trail.Context, selfCandidate dto.Candidate) (trail.Success, *trail.Error) {
	idParam := c.Request.PathValue("id")
	if idParam == "" {
		return c.Success(http.StatusOK, nil, selfCandidate)
	}
	candidateId, err := strconv.Atoi(idParam)
	if err != nil {
		return model.ErrInvalidCandidateIdParam(err)
	}

	dbGiteaToken, err := cmp.queries.GetGiteaTokenByCandidateId(c, int64(candidateId))
	if err != nil {
		return model.ErrDatabase(err)
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
		return model.ErrFetchCandidateInfo(err)
	}
	return c.Success(http.StatusOK, nil, dto.CandidateFromZ01auth(candidate))
}
