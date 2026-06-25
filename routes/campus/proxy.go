package campus

import (
	"libraone/internal/dto"
	"libraone/internal/lib/trail"
	"libraone/internal/model"
	"libraone/internal/services/profile"
	"net/http"
	"strings"
)

type Campus struct {
	ProfileService profile.ProfileService
}

func (cmp *Campus) ProxyHandler(c *trail.Context, candidate dto.Candidate) (trail.Success, *trail.Error) {
	path := strings.TrimPrefix(c.Request.URL.Path, "/campus")
	resp, err := cmp.ProfileService.ForwardRequest(c.Request, path)
	if err != nil {
		return model.ErrCampusProfileProxy(err)
	}
	defer resp.Body.Close()

	headers := make(http.Header)
	for key, values := range resp.Header {
		for _, value := range values {
			headers.Add(key, value)
		}
	}

	return c.Success(http.StatusOK, headers, resp.Body)
}
