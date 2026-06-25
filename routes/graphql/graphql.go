package graphql

import (
	"fmt"
	"libraone/internal/dto"
	"libraone/internal/lib/trail"
	"libraone/internal/model"
	"libraone/internal/services/graphql"
	"net/http"
	"strings"
)

const API_BASE = "https://learn.zone01oujda.ma/api/graphql-engine/v1/graphql"

type GraphQL struct {
	token *graphql.TokenSupplier
}

func New(token *graphql.TokenSupplier) *GraphQL {
	return &GraphQL{token: token}
}

func (gql *GraphQL) ProxyHandler(c *trail.Context, candidate dto.Candidate) (trail.Success, *trail.Error) {
	path := strings.TrimPrefix(c.Request.URL.Path, "/graphql")
	resp, err := gql.request(c.Request, path)
	if err != nil {
		return model.ErrGraphqlProxy(err)
	}

	headers := make(http.Header)
	for key, values := range resp.Header {
		for _, value := range values {
			headers.Add(key, value)
		}
	}

	return c.Success(http.StatusOK, headers, resp.Body)
}

func (gqls *GraphQL) request(originalReq *http.Request, targetPath string) (*http.Response, error) {
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
	req.Header.Set("Authorization", "Bearer "+gqls.token.Get())
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("executing request: %w", err)
	}

	return resp, nil
}
