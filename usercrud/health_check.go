package usercrud

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"gorm.io/gorm"
)

func HealthCheck(c *gin.Context, db *gorm.DB) {

	if c.Request.Method != http.MethodGet {
		c.Header("Allow", "GET")
		c.Status(http.StatusMethodNotAllowed)
		return
	}

	c.Header("Cache-Control", "no-cache")

	if len(c.Request.URL.Query()) > 0 {
		c.Status(http.StatusBadRequest)
		return
	}

	if c.Request.ContentLength > 0 {
		c.Status(http.StatusBadRequest)
		return
	}

	sqlDB, err := db.DB()
	if err != nil || sqlDB.Ping() != nil {
		c.Status(http.StatusServiceUnavailable)
		return
	}

	c.Status(http.StatusOK)
}
