package profile

import (
	"bytes"
	"encoding/json"
	"fmt"
	tokensupplier "libraone/internal/token_supplier"
	"net/http"
)

func newTokenFetcher(login, password string) tokensupplier.Fetcher {
	return func() (toke string, err error) {
		payload, err := json.Marshal(map[string]string{"username": login, "password": password})
		if err != nil {
			return
		}
		resp, err := http.Post(LOGIN_ENDPOINT, "application/json", bytes.NewBuffer(payload))
		if err != nil {
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			err = fmt.Errorf("login failed with status: %s", resp.Status)
			return
		}

		var loginResult struct {
			Token string `json:"token"`
		}
		if err = json.NewDecoder(resp.Body).Decode(&loginResult); err != nil {
			return
		}
		return loginResult.Token, nil
	}
}
