package middlewares

import (
	db "libraone/db/generated"
	"libraone/internal/dto"
	"libraone/internal/lib/trail"

	"github.com/xySaad/z01auth"
)

func EnsureTalentRole(queries *db.Queries, z01authConfig z01auth.Config) trail.Middleware[dto.Candidate] {
	return func(c *trail.Context) (dto.Candidate, *trail.Error) {
		candidate, err := EnsureAuthenticated(queries, z01authConfig)(c)
		if err != nil {
			return candidate, err
		}

		if candidate.Role != string(z01auth.CandidateRole_TALENT) {
			return ErrCandidateIsNotTalent(nil)
		}

		return candidate, nil
	}
}
