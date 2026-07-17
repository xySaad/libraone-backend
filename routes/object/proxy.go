package object

import (
	"fmt"
	"libraone/internal/dto"
	"libraone/internal/lib/trail"
	"net/http"
	"strings"
)

const apiBase = "https://learn.zone01oujda.ma/api/object"

type Object struct{}

func (o *Object) ProxyHandler(c *trail.Context, candidate dto.Candidate) (trail.Success, *trail.Error) {
	path := strings.TrimPrefix(c.Request.URL.Path, "/object")

	endpoint := apiBase + path
	if c.Request.URL.RawQuery != "" {
		endpoint += "?" + c.Request.URL.RawQuery
	}

	req, err := http.NewRequest(c.Request.Method, endpoint, c.Request.Body)
	if err != nil {
		return c.Error(
			trail.NewPublicError(http.StatusBadGateway, "failed to build upstream request"),
			fmt.Errorf("object proxy: build request: %w", err),
		)
	}
	for k, vv := range c.Request.Header {
		for _, v := range vv {
			req.Header.Add(k, v)
		}
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return c.Error(
			trail.NewPublicError(http.StatusBadGateway, "upstream request failed"),
			fmt.Errorf("object proxy: do request: %w", err),
		)
	}

	headers := make(http.Header)
	for k, vv := range resp.Header {
		headers[k] = vv
	}

	return c.Success(resp.StatusCode, headers, resp.Body)
}
