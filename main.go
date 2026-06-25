package main

//go:generate go run github.com/xySaad/cfgo
//go:generate rm -r db/generated
//go:generate mkdir db/generated
//go:generate go run github.com/sqlc-dev/sqlc/cmd/sqlc generate -f db/sqlc.json
//go:generate go run -tags sqlite3 github.com/golang-migrate/migrate/v4/cmd/migrate -source file://db/migrations -database sqlite3://db/database.db up

import (
	"database/sql"
	"fmt"
	config "libraone/config/generated"
	db "libraone/db/generated"
	"libraone/internal/lib/trail"
	"libraone/internal/middlewares"
	"libraone/internal/services/graphql"
	"libraone/internal/services/profile"
	"libraone/routes"
	"log"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
	"github.com/xySaad/z01auth"
)

func main() {

	config := config.GetConfig()
	graphqlToken := graphql.MustNewTokenSupplier(config.GRAPHQL_LOGIN, config.GRAPHQL_PASSWORD)
	z01authConfig := z01auth.New(config.GiteaClientID, config.GiteaClientSecret, config.GiteaRedirectURL, graphqlToken)
	sqlDB, err := sql.Open("sqlite3", "db/database.db")
	if err != nil {
		log.Fatal("failed to open SQLite database", err)
	}
	queries := db.New(sqlDB)

	routes := routes.Routes{}
	router := trail.DefaultRouter()
	router.AddRoute("GET /oauth/gitea", routes.OAuth.Gitea.Entry(z01authConfig))
	router.AddRoute("GET /oauth/gitea/callback", routes.OAuth.Gitea.Callback(config, z01authConfig, queries))

	talentOnly := trail.Extend(router, middlewares.EnsureTalentRole(queries, z01authConfig))
	talentOnly.AddRoute("/graphql/", routes.GraphQL(graphqlToken).ProxyHandler)

	profileTokenSupplier := profile.MustNewService(config.PROFILE_LOGIN, config.PROFILE_PASSWORD)
	talentOnly.AddRoute("/campus/", routes.Campus(&profileTokenSupplier).ProxyHandler)

	candidateRoutes := routes.Candidate(queries, z01authConfig)
	talentOnly.AddRoute("GET /candidate/", candidateRoutes.Candidate)
	talentOnly.AddRoute("GET /candidate/{id}", candidateRoutes.Candidate)

	addr := ":5051"
	fmt.Println("HTTP server listening", addr)
	err = http.ListenAndServe(addr, router)
	if err != nil {
		log.Fatal("HTTP server exited with error", err)
	}
}
