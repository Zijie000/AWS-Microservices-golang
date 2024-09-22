package main

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"gorm.io/gorm"
)

var db *gorm.DB

func main() {
	var err error
	db, err = initDB()
	if err != nil {
		panic("Failed to connect to database!")
	}

	r := gin.Default()

	r.Any("/healthz", HealthCheck)

	r.Run(":8080")
}

func HealthCheck(c *gin.Context) {

	if c.Request.Method != http.MethodGet {
		c.Header("Allow", "GET")
		c.Status(http.StatusMethodNotAllowed)
		return
	}

	c.Header("Cache-Control", "no-cache")

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
