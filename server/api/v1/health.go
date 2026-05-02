package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type HealthApi struct{}

var Health = new(HealthApi)

func (a *HealthApi) GetHealth(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}
