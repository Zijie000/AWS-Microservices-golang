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

	r.Any("/healthz", healthCheck)

	r.Any("/v1/user", registerUser)

	r.GET("/v1/user/self", getUser)

	r.PUT("/v1/user/self", updateUser)

	r.POST("/v1/user/self", NotSupported)
	r.DELETE("/v1/user/self", NotSupported)
	r.PATCH("/v1/user/self", NotSupported)
	r.HEAD("/v1/user/self", NotSupported)
	r.OPTIONS("/v1/user/self", NotSupported)

	r.NoRoute(func(c *gin.Context) {
		c.Status(http.StatusBadRequest)
	})

	r.Run(":8080")
}

func NotSupported(c *gin.Context) {

	c.Header("Allow", "GET/PUT")
	c.Status(http.StatusMethodNotAllowed)

	c.Header("Cache-Control", "no-cache")

}

func healthCheck(c *gin.Context) {

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
