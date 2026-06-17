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
	"libraone/routes"
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
	"github.com/xySaad/z01auth"
)

func main() {
	config := config.GetConfig()
	z01authConfig := z01auth.New(config.GiteaClientID, config.GiteaClientSecret, config.GRAPHQL_TOKEN)
	sqlDB, err := sql.Open("sqlite3", "db/database.db")
	if err != nil {
		panic(err)
	}

	queries := db.New(sqlDB)

	routes := routes.Routes{}
	router := gin.Default()
	router.GET("/oauth/gitea", routes.OAuth.Gitea.Entry(z01authConfig))
	//TODO: change /api/auth to /oauth/gitea/callback
	router.GET("/api/auth/callback", routes.OAuth.Gitea.Callback(config, z01authConfig, queries))

	authenticated := middlewares.Group(router.RouterGroup, "", middlewares.EnsureAuthenticated(queries))
	authenticated.GET("/campus/online", routes.Campus.Online)

	err = http.ListenAndServe(":5051", router)
	if err != nil {
		fmt.Println(err)
		return
	}
}
