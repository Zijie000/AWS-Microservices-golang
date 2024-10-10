package test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"

	"github.com/stretchr/testify/assert"

	"encoding/base64"

	"bytes"

	"webapp/usercrud"
)

func TestRegisterUser(t *testing.T) {
	db, err := setupTestDB()
	assert.NoError(t, err)

	// 使用 httptest 创建一个请求响应记录器

	jsonBody := `{"email":"rynzzj997@gmail.com", "password":"12345678", "first_name":"Zijie", "last_name":"Zhou"}`
	reqBody := bytes.NewBuffer([]byte(jsonBody))

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/v1/user", reqBody)

	// 创建一个 Gin 引擎
	r := gin.Default()

	// 定义路由并处理
	r.POST("/v1/user", func(c *gin.Context) {
		usercrud.RegisterUser(c, db)
	})

	// 运行 Gin 引擎处理请求
	r.ServeHTTP(w, req)

	// 断言返回状态码 201
	assert.Equal(t, http.StatusCreated, w.Code)

}

func TestGetUser(t *testing.T) {
	db, err := setupTestDB()
	assert.NoError(t, err)

	username := "rynzzj997@gmail.com"
	password := "12345678"

	// 将用户名和密码组合为 "username:password"
	auth := username + ":" + password

	// 将 auth 进行 Base64 编码
	encodedAuth := base64.StdEncoding.EncodeToString([]byte(auth))

	// 使用 httptest 创建一个请求响应记录器

	jsonBody := `{"email":"rynzzj997@gmail.com", "password":"12345678", "first_name":"Zijie", "last_name":"Zhou"}`
	reqBody := bytes.NewBuffer([]byte(jsonBody))

	w1 := httptest.NewRecorder()
	w2 := httptest.NewRecorder()
	req1, _ := http.NewRequest("POST", "/v1/user", reqBody)

	req2, _ := http.NewRequest("GET", "/v1/user/self", nil)

	// 设置 Authorization 头，格式为 "Basic base64(username:password)"
	req2.Header.Add("Authorization", "Basic "+encodedAuth)

	// 创建一个 Gin 引擎
	r := gin.Default()

	// 定义路由并处理
	r.POST("/v1/user", func(c *gin.Context) {
		usercrud.RegisterUser(c, db)
	})

	r.GET("/v1/user/self", func(c *gin.Context) {
		usercrud.GetUser(c, db)
	})

	// 运行 Gin 引擎处理请求
	r.ServeHTTP(w1, req1)
	r.ServeHTTP(w2, req2)

	// 断言返回状态码 200
	assert.Equal(t, http.StatusOK, w2.Code)

}

func TestUpdateUser(t *testing.T) {
	db, err := setupTestDB()
	assert.NoError(t, err)

	username := "rynzzj997@gmail.com"
	password := "12345678"

	passwordUpdated := "87654321"

	// 将用户名和密码组合为 "username:password"
	authUpdate := username + ":" + password

	authGet := username + ":" + passwordUpdated

	// 将 auth 进行 Base64 编码
	encodedAuthUpdate := base64.StdEncoding.EncodeToString([]byte(authUpdate))
	encodedAuthGet := base64.StdEncoding.EncodeToString([]byte(authGet))

	// 使用 httptest 创建一个请求响应记录器

	jsonBody1 := `{"email":"rynzzj997@gmail.com", "password":"12345678", "first_name":"Zijie", "last_name":"Zhou"}`
	reqBody1 := bytes.NewBuffer([]byte(jsonBody1))

	jsonBody2 := `{"password":"87654321"}`
	reqBody2 := bytes.NewBuffer([]byte(jsonBody2))

	w1 := httptest.NewRecorder()
	w2 := httptest.NewRecorder()
	req1, _ := http.NewRequest("POST", "/v1/user", reqBody1)

	req2, _ := http.NewRequest("GET", "/v1/user/self", nil)

	req3, _ := http.NewRequest("PUT", "/v1/user/self", reqBody2)

	// 设置 Authorization 头，格式为 "Basic base64(username:password)"

	//UPDATE USER
	req3.Header.Add("Authorization", "Basic "+encodedAuthUpdate)

	//GET USER
	req2.Header.Add("Authorization", "Basic "+encodedAuthGet)

	// 创建一个 Gin 引擎
	r := gin.Default()

	// 定义路由并处理
	r.POST("/v1/user", func(c *gin.Context) {
		usercrud.RegisterUser(c, db)
	})

	r.GET("/v1/user/self", func(c *gin.Context) {
		usercrud.GetUser(c, db)
	})
	r.PUT("/v1/user/self", func(c *gin.Context) {
		usercrud.UpdateUser(c, db)
	})

	// 运行 Gin 引擎处理请求
	r.ServeHTTP(w1, req1)
	r.ServeHTTP(w1, req3)
	r.ServeHTTP(w2, req2)

	// 断言返回状态码 200
	assert.Equal(t, http.StatusOK, w2.Code)

}
