package test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"

	"github.com/stretchr/testify/assert"

	"gorm.io/driver/sqlite"

	"gorm.io/gorm"

	"webapp/usercrud"
)

func setupTestDB() (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	db.AutoMigrate(&usercrud.User{})

	return db, nil
}

// database healthy 200ok
func TestHealthCheckSuccess(t *testing.T) {
	db, err := setupTestDB()
	assert.NoError(t, err)

	// 使用 httptest 创建一个请求响应记录器
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/healthz", nil)

	// 创建一个 Gin 引擎
	r := gin.Default()

	// 定义路由并处理
	r.GET("/healthz", func(c *gin.Context) {
		usercrud.HealthCheck(c, db)
	})

	// 运行 Gin 引擎处理请求
	r.ServeHTTP(w, req)

	// 断言返回状态码 200
	assert.Equal(t, http.StatusOK, w.Code)

}

// database unhealthy 503 unavailable
func TestHealthCheckFailure(t *testing.T) {

	db, err := setupTestDB()
	assert.NoError(t, err)

	// 使用 httptest 创建一个请求响应记录器
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/healthz", nil)

	// 创建一个 Gin 引擎
	r := gin.Default()

	mysqlDB, _ := db.DB()

	mysqlDB.Close()

	// 定义路由并处理
	r.GET("/healthz", func(c *gin.Context) {
		usercrud.HealthCheck(c, db)
	})

	// 运行 Gin 引擎处理请求
	r.ServeHTTP(w, req)

	// 断言返回状态码 503
	assert.Equal(t, http.StatusServiceUnavailable, w.Code)

}
