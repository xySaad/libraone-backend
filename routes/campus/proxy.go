package campus

import (
	"io"
	"libraone/internal/dto"
	"libraone/internal/services/profile"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Campus struct {
	ProfileService profile.ProfileService
}

func (cmp *Campus) ProxyHandler(c *gin.Context, candidate *dto.Candidate) {
	path := c.Param("path")
	resp, err := cmp.ProfileService.ForwardRequest(c.Request, path)
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
