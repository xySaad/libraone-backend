package profile

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	tokensupplier "libraone/internal/token_supplier"
	"net/http"
)

const (
	API_BASE       = "https://mapl.zone01oujda.ma"
	LOGIN_ENDPOINT = API_BASE + "/login"
)

type ProfileService struct {
	token *tokensupplier.Supplier
}

func MustNewService(login, password string) ProfileService {
	return ProfileService{token: tokensupplier.MustNewSupplier(newTokenFetcher(login, password))}
}

func NewService(login, password string) (ProfileService, error) {
	token, err := tokensupplier.NewSupplier(newTokenFetcher(login, password))
	return ProfileService{token: token}, err
}

func (ps *ProfileService) request(originalReq *http.Request, targetPath string) (*http.Response, error) {
	endpoint := API_BASE + targetPath + "?" + originalReq.URL.RawQuery
	req, err := http.NewRequest(originalReq.Method, endpoint, originalReq.Body)
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}
	for k, vv := range originalReq.Header {
		for _, v := range vv {
			req.Header.Add(k, v)
		}
	}
	req.Header.Set("X-TOKEN", ps.token.Get())
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("executing request: %w", err)
	}

	return resp, nil
}

func (ps *ProfileService) ForwardRequest(originalReq *http.Request, targetPath string) (*http.Response, error) {
	resp, err := ps.request(originalReq, targetPath)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusForbidden {
		return resp, nil
	}

	bodyCopy, err := io.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		return nil, fmt.Errorf("reading 403 response body: %w", err)
	}
	var jsonResp map[string]string
	err = json.Unmarshal(bodyCopy, &jsonResp)
	if err != nil {
		return nil, fmt.Errorf("parsing 403 body: %w", err)
	}
	if jsonResp["detail"] != "Not authenticated" {
		resp.Body = io.NopCloser(bytes.NewBuffer(bodyCopy))
		return resp, nil
	}

	if err := ps.token.RefreshToken(); err != nil {
		return nil, fmt.Errorf("refreshing token: %w", err)
	}
	return ps.request(originalReq, targetPath)
}
