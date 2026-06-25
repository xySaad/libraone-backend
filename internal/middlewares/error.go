package middlewares

import (
	"libraone/internal/dto"
	"libraone/internal/lib/trail"
	"libraone/internal/model"
	"net/http"
)

func e(status int, message string) func(err error) (dto.Candidate, *trail.Error) {
	return func(err error) (dto.Candidate, *trail.Error) {
		publicErr := trail.NewPublicError(status, message)
		trailErr := trail.NewError(publicErr, err)
		return dto.Candidate{}, &trailErr
	}
}

func fromBuilder(builder model.ErrorBuilder) func(err error) (dto.Candidate, *trail.Error) {
	return func(err error) (dto.Candidate, *trail.Error) {
		_, trailErr := builder(err)
		return dto.Candidate{}, trailErr
	}
}

// Forbidden
var (
	ErrCandidateIsNotTalent = e(http.StatusForbidden, "candidate is not a talent")
	ErrPlatformBanned       = e(http.StatusForbidden, "candidate does not have platform access")
)

// Unauthorized
var (
	ErrMissingCookie  = e(http.StatusUnauthorized, "missing or empty session cookie")
	ErrInvalidSession = e(http.StatusUnauthorized, "session token not found or expired")
)

// Internal
var (
	ErrFetchCandidateInfo = fromBuilder(model.ErrFetchCandidateInfo)
	ErrDatabase           = fromBuilder(model.ErrDatabase)
)
