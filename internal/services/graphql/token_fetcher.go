package graphql

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	tokensupplier "libraone/internal/token_supplier"
	"net/http"
)

const (
	API_BASE       = "https://learn.zone01oujda.ma/api"
	LOGIN_ENDPOINT = API_BASE + "/auth/signin"
)

func NewTokenFetcher(login, password string) tokensupplier.Fetcher {
	return func() (token string, err error) {
		req, err := http.NewRequest("POST", LOGIN_ENDPOINT, nil)
		if err != nil {
			return
		}
		authorization := fmt.Appendf(nil, "%s:%s", login, password)
		authorizationB64 := base64.StdEncoding.EncodeToString(authorization)
		req.Header.Set("Authorization", "Basic "+authorizationB64)
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			err = fmt.Errorf("login failed with status: %s", resp.Status)
			return
		}

		return token, json.NewDecoder(resp.Body).Decode(&token)
	}
}
