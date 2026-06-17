package routes

import (
	"libraone/routes/campus"
	"libraone/routes/oauth"
)

type Routes struct {
	oauth.OAuth
	campus.Campus
}
