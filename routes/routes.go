package routes

import (
	gql "libraone/internal/services/graphql"
	"libraone/internal/services/profile"
	"libraone/routes/campus"
	"libraone/routes/graphql"
	"libraone/routes/oauth"
)

type Routes struct {
	oauth.OAuth
}

func (Routes) Campus(profileService *profile.ProfileService) *campus.Campus {
	return &campus.Campus{ProfileService: *profileService}
}

func (Routes) GraphQL(token *gql.TokenSupplier) *graphql.GraphQL {
	return graphql.New(token)
}
