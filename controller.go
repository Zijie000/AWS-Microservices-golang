package main

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"gorm.io/gorm"

	"webapp/usercrud"

	"webapp/mydb"
)

var db *gorm.DB

func main() {
	var err error
	db, err = mydb.InitDB()
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

	usercrud.HealthCheck(c, db)
}

func registerUser(c *gin.Context) {

	usercrud.RegisterUser(c, db)
}

func getUser(c *gin.Context) {

	usercrud.GetUser(c, db)
}

func updateUser(c *gin.Context) {

	usercrud.UpdateUser(c, db)
}
