package main

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"gorm.io/gorm"

	"webapp/usercrud"

	"webapp/mydb"

	"webapp/observability"

	"webapp/profilepic"

	"log"

	"os"
)

var db *gorm.DB

func main() {

	logFile, erro := os.OpenFile("/var/log/webapp.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if erro != nil {
		log.Fatalf("Failed to open log file: %v", erro)
	}

	gin.DefaultWriter = logFile
	gin.DefaultErrorWriter = logFile

	var err error
	db, err = mydb.InitDB()

	observability.Init()
	defer observability.Close()

	if err != nil {
		panic("Failed to connect to database!")
	}

	r := gin.Default()

	r.Any("/healthz", healthCheck)

	r.Any("/v1/user", registerUser)

	r.GET("/v1/user/self", getUser)

	r.PUT("/v1/user/self", updateUser)

	r.POST("/v1/user/self/pic", uploadPicture)

	r.GET("/v1/user/self/pic", getPicture)

	r.DELETE("/v1/user/self/pic", deletePicture)

	r.POST("/v1/user/self", NotSupported)
	r.DELETE("/v1/user/self", NotSupported)
	r.PATCH("/v1/user/self", NotSupported)
	r.HEAD("/v1/user/self", NotSupported)
	r.OPTIONS("/v1/user/self", NotSupported)

	r.PATCH("/v1/user/self/pic", NotSupported)
	r.HEAD("/v1/user/self/pic", NotSupported)
	r.OPTIONS("/v1/user/self/pic", NotSupported)

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

func uploadPicture(c *gin.Context) {

	profilepic.UploadImage(c, db)
}

func getPicture(c *gin.Context) {

	profilepic.GetImage(c, db)
}

func deletePicture(c *gin.Context) {

	profilepic.DeleteImage(c, db)
}
