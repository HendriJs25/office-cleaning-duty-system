package health

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Health struct{}

func NewHandler() *Health {
	return &Health{}
}

func (h *Health) HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "OK",
		"Message": "Welcome to cleaning app",
	})
}
