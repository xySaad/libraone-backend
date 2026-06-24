package routes

import (
	db "libraone/db/generated"
	gql "libraone/internal/services/graphql"
	"libraone/internal/services/profile"
	"libraone/routes/campus"
	"libraone/routes/candidate"
	"libraone/routes/graphql"
	"libraone/routes/oauth"

	"github.com/xySaad/z01auth"
)

type Routes struct {
	oauth.OAuth
}

func (Routes) Campus(profileService *profile.ProfileService) *campus.Campus {
	return &campus.Campus{ProfileService: *profileService}
}
func (Routes) Candidate(queries *db.Queries, z01authConfig z01auth.Config) *candidate.Candidate {
	return candidate.New(queries, z01authConfig)
}

func (Routes) GraphQL(token *gql.TokenSupplier) *graphql.GraphQL {
	return graphql.New(token)
}
