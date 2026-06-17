package routes

import (
	"libraone/internal/services/profile"
	"libraone/routes/campus"
	"libraone/routes/oauth"
)

type Routes struct {
	oauth.OAuth
}

func (Routes) Campus(profileService *profile.ProfileService) *campus.Campus {
	return &campus.Campus{ProfileService: *profileService}
}
