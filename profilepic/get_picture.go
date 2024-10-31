package profilepic

import (
	"fmt"
	"net/http"
	"webapp/usercrud"

	"webapp/observability"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"golang.org/x/crypto/bcrypt"

	"log"

	"time"
)

func GetImage(c *gin.Context, db *gorm.DB) {

	start := time.Now()

	err := observability.Client.Incr("Get Image API", nil, 1)
	if err != nil {
		log.Printf("Error incrementing Get Image API count: %v", err)
	}

	if len(c.Request.URL.Query()) > 0 {
		c.Status(http.StatusBadRequest)
		return
	}

	if c.Request.ContentLength > 0 {
		c.Status(http.StatusBadRequest)
		return
	}

	email, password, hasAuth := c.Request.BasicAuth()
	if !hasAuth {
		c.Status(http.StatusUnauthorized)
		return
	}

	var user usercrud.User
	if err := db.Where("email = ?", email).First(&user).Error; err != nil {
		c.Status(http.StatusUnauthorized)
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		c.Status(http.StatusUnauthorized)
		return
	}

	startdb := time.Now()

	var picture Picture
	if err := db.Where("id = ?", user.ID).First(&picture).Error; err != nil {
		c.Status(http.StatusNotFound)
		return
	}

	durationdb := time.Since(startdb).Milliseconds()

	err = observability.Client.Timing("api.response_time.GetImageAPIDataBase", time.Duration(durationdb)*time.Millisecond, nil, 1)

	if err != nil {
		log.Printf("Error recording Get Image API DataBase timing: %v", err)
	}

	pictureInfo := struct {
		FileName   string `json:"file_name"`
		ID         string `json:"id"`
		USERID     string `json:"user_id"`
		URL        string `json:"url"`
		UploadDate string `json:"upload_date"`
	}{
		FileName:   picture.FileName,
		ID:         fmt.Sprintf("%d", user.ID),
		USERID:     fmt.Sprintf("%d", user.ID),
		URL:        picture.URL,
		UploadDate: picture.UploadDate,
	}

	duration := time.Since(start).Milliseconds()

	err = observability.Client.Timing("api.response_time.GetImageAPI", time.Duration(duration)*time.Millisecond, nil, 1)

	if err != nil {
		log.Printf("Error recording Get Image API timing: %v", err)
	}

	c.JSON(http.StatusOK, pictureInfo)

}
