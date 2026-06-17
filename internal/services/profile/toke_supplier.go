package profile

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
)

const LOGIN_ENDPOINT = PROFILE_API_BASE + "/login"

type TokenSupplier struct {
	Login    string
	Password string
	mx       sync.Mutex
	token    string
}

func (ts *TokenSupplier) Get() string {
	return ts.token
}

func (ts *TokenSupplier) refreshToken() error {
	payload, err := json.Marshal(map[string]string{"username": ts.Login, "password": ts.Password})
	if err != nil {
		return err
	}
	resp, err := http.Post(LOGIN_ENDPOINT, "application/json", bytes.NewBuffer(payload))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("login failed with status: %s", resp.Status)
	}

	var loginResult struct {
		Token string `json:"token"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&loginResult); err != nil {
		return err
	}

	ts.mx.Lock()
	defer ts.mx.Unlock()
	ts.token = loginResult.Token

	return nil
}
