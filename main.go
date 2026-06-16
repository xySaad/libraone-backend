package main

//go:generate go run github.com/xySaad/cfgo@v0.0.5
import (
	"encoding/json"
	"fmt"
	config "libraone/config/generated"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/xySaad/z01auth"
)

func main() {
	router := gin.Default()
	config := config.GetConfig()
	z01authConfig := z01auth.New(config.GiteaClientID, config.GiteaClientSecret, config.GRAPHQL_TOKEN)
	router.GET("/oauth/gitea", func(c *gin.Context) {
		//TOOD: add state
		c.Redirect(http.StatusMovedPermanently, z01authConfig.AuthCodeURL(""))
	})

	//TODO: change /api/auth to /oauth/gitea/callback
	router.GET("/api/auth/callback", func(c *gin.Context) {
		z01User, err := z01authConfig.Callback(c.Query("code"))
		if err != nil {
			fmt.Println("error:", err)
			return
		}
		//TODO: save basic information in DB
		//TODO: generate JWT cookie, save it in DB, and set it in the response headers
		json.NewEncoder(os.Stdout).Encode(z01User)
		c.Redirect(http.StatusSeeOther, config.CallbackRedirectURL)
	})

	s := &http.Server{Addr: ":5051", Handler: router}
	s.ListenAndServe()
}
