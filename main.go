package main

//go:generate go run github.com/xySaad/cfgo
//go:generate rm -r db/generated
//go:generate go run github.com/sqlc-dev/sqlc/cmd/sqlc generate -f db/sqlc.json
//go:generate go run -tags sqlite3 github.com/golang-migrate/migrate/v4/cmd/migrate -source file://db/migrations -database sqlite3://db/database.db up

import (
	"database/sql"
	"fmt"
	config "libraone/config/generated"
	db "libraone/db/generated"
	"libraone/internal/middlewares"
	"libraone/internal/services/graphql"
	"libraone/internal/services/profile"
	"libraone/routes"
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
	"github.com/xySaad/z01auth"
)

func main() {
	config := config.GetConfig()
	graphqlToken := graphql.MustNewTokenSupplier(config.GRAPHQL_LOGIN, config.GRAPHQL_PASSWORD)
	z01authConfig := z01auth.New(config.GiteaClientID, config.GiteaClientSecret, config.GiteaRedirectURL, graphqlToken)
	sqlDB, err := sql.Open("sqlite3", "db/database.db")
	if err != nil {
		panic(err)
	}
	queries := db.New(sqlDB)

	routes := routes.Routes{}
	router := gin.Default()
	router.GET("/oauth/gitea", routes.OAuth.Gitea.Entry(z01authConfig))
	router.GET("/oauth/gitea/callback", routes.OAuth.Gitea.Callback(config, z01authConfig, queries))

	talentOnly := middlewares.Group(router.RouterGroup, "", middlewares.EnsureTalentRole(queries))
	profileTokenSupplier := profile.MustNewService(config.PROFILE_LOGIN, config.PROFILE_PASSWORD)
	campusRoutes := routes.Campus(&profileTokenSupplier)
	talentOnly.Any("/campus/*path", campusRoutes.ProxyHandler)
	talentOnly.GET("/candidate/", campusRoutes.Candidate(queries))
	talentOnly.GET("/candidate/:login", campusRoutes.Candidate(queries))

	err = http.ListenAndServe(":5051", router)
	if err != nil {
		fmt.Println(err)
		return
	}
}
