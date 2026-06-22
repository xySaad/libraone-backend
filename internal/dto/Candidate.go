package dto

import (
	db "libraone/db/generated"
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

func CandidateFromDB(candidate db.Candidate) *Candidate {
	return &Candidate{
		ID:           candidate.ID,
		Role:         candidate.Role,
		AvatarUrl:    candidate.AvatarUrl,
		Description:  candidate.Description,
		GiteaLogin:   candidate.GiteaLogin,
		GraphqlLogin: candidate.GraphqlLogin,
		GraphqlID:    candidate.GraphqlID,
		Campus:       candidate.Campus,
	}
}
