package usercrud

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"gorm.io/gorm"

	"webapp/observability"

	"log"

	"time"
)

func HealthCheck(c *gin.Context, db *gorm.DB) {

	start := time.Now()

	err := observability.Client.Incr("Health Check API", nil, 1)
    if err != nil {
        log.Printf("Error incrementing Health Check API count: %v", err)
    }

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

	duration := time.Since(start).Milliseconds()

	err = observability.Client.Timing("api.response_time.HealthCheckAPI", time.Duration(duration)*time.Millisecond, nil, 1)
    if err != nil {
        log.Printf("Error recording Health Check API timing: %v", err)
    }

	c.Status(http.StatusOK)

}
