package dto

import (
	"github.com/xySaad/z01auth"
)

type Candidate struct {
	ID           int64  `json:"id"`
	Role         string `json:"role"`
	AvatarUrl    string `json:"avatar_url"`
	Description  string `json:"description"`
	GiteaLogin   string `json:"gitea_login"`
	GraphqlLogin string `json:"graphql_login"`
	GraphqlID    int64  `json:"graphql_id"`
	Campus       string `json:"campus"`
}

func CandidateFromZ01auth(candidate *z01auth.Candidate) Candidate {
	return Candidate{
		ID:           int64(candidate.GiteaID),
		Role:         string(candidate.Role),
		AvatarUrl:    candidate.AvatarURL,
		Description:  candidate.Description,
		GiteaLogin:   candidate.GiteaLogin,
		GraphqlLogin: candidate.GraphqlLogin,
		GraphqlID:    int64(candidate.GraphqlId),
		Campus:       candidate.Campus,
	}
}
