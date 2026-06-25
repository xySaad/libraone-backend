package model

import (
	"libraone/internal/lib/trail"
	"net/http"
)

type ErrorBuilder func(err error) (trail.Success, *trail.Error)

func e(status int, message string) ErrorBuilder {
	return func(err error) (trail.Success, *trail.Error) {
		publicErr := trail.NewPublicError(status, message)
		trailErr := trail.NewError(publicErr, err)
		return trail.Success{Status: status}, &trailErr
	}
}

// Internal Error
var (
	ErrFetchCandidateInfo = e(http.StatusInternalServerError, "failed to fetch candidate info from z01auth")
	ErrFetchCandidateId   = e(http.StatusInternalServerError, "failed to fetch candidate ID from z01auth")
	ErrOAuthCodeExchange  = e(http.StatusInternalServerError, "failed to exchange OAuth code for token")
	ErrInsertGiteaToken   = e(http.StatusInternalServerError, "failed to persist Gitea token")
	ErrCreateSession      = e(http.StatusInternalServerError, "failed to create session")
	ErrCreateCandidate    = e(http.StatusInternalServerError, "failed to create candidate")
	ErrDatabase           = e(http.StatusInternalServerError, "database error")
	ErrGraphqlProxy       = e(http.StatusInternalServerError, "failed to proxy request to GraphQL engine")
	ErrCampusProfileProxy = e(http.StatusInternalServerError, "failed to proxy request to campus profile service")
)

// BadRequest
var (
	ErrInvalidCandidateIdParam = e(http.StatusBadRequest, "invalid candidate id parameter")
)
