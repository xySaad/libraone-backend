package graphql

import (
	"fmt"
	"io"
	db "libraone/db/generated"
	"libraone/internal/services/graphql"
	"net/http"

	"github.com/gin-gonic/gin"
)

const API_BASE = "https://learn.zone01oujda.ma/api/graphql-engine/v1/graphql"

type GraphQL struct {
	token *graphql.TokenSupplier
}

func New(token *graphql.TokenSupplier) *GraphQL {
	return &GraphQL{token: token}
}

func (gql *GraphQL) ProxyHandler(c *gin.Context, candidate *db.Candidate) {
	path := c.Param("path")
	resp, err := gql.request(c.Request, path)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to proxy request"})
		return
	}
	defer resp.Body.Close()

	for key, values := range resp.Header {
		for _, value := range values {
			c.Writer.Header().Add(key, value)
		}
	}

	c.Status(resp.StatusCode)
	io.Copy(c.Writer, resp.Body)
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
