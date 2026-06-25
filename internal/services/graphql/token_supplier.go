package graphql

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	tokensupplier "libraone/internal/token_supplier"
	"strings"
	"time"
)

type TokenSupplier struct {
	*tokensupplier.Supplier
	exp time.Time
}

func MustNewTokenSupplier(login, password string) *TokenSupplier {
	sp, err := NewTokenSupplier(login, password)
	if err != nil {
		panic(err)
	}
	return sp
}

func NewTokenSupplier(login, password string) (*TokenSupplier, error) {
	sp, err := tokensupplier.NewSupplier(NewTokenFetcher(login, password))
	if err != nil {
		return nil, err
	}

	exp, err := parseTokenExpiry(sp.Get())
	if err != nil {
		return nil, fmt.Errorf("parse token expiry: %w", err)
	}
	return &TokenSupplier{Supplier: sp, exp: exp}, nil
}

func (ts *TokenSupplier) Get() string {
	if time.Now().Before(ts.exp) {
		return ts.Supplier.Get()
	}
	err := ts.RefreshToken()
	if err != nil {

		return ts.Supplier.Get()
	}
	newToken := ts.Supplier.Get()
	newExp, err := parseTokenExpiry(newToken)
	if err != nil {

		return newToken
	}

	ts.exp = newExp
	return newToken
}

type jwtClaims struct {
	Exp int64 `json:"exp"`
}

// parseTokenExpiry extracts the "exp" claim from a
// JWT without verifying its signature.
func parseTokenExpiry(token string) (time.Time, error) {
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return time.Time{}, fmt.Errorf("invalid JWT format")
	}

	payload, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return time.Time{}, fmt.Errorf("decode payload: %w", err)
	}

	var claims jwtClaims
	if err := json.Unmarshal(payload, &claims); err != nil {
		return time.Time{}, fmt.Errorf("unmarshal claims: %w", err)
	}
	if claims.Exp == 0 {
		return time.Time{}, fmt.Errorf("token has no exp claim")
	}

	return time.Unix(claims.Exp, 0), nil
}
