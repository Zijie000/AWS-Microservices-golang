package usercrud

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"golang.org/x/crypto/bcrypt"

	"time"

	"gorm.io/gorm"

	"webapp/observability"

	"log"
)

func GetUser(c *gin.Context, db *gorm.DB) {

	start := time.Now()

	err := observability.Client.Incr("Get User API", nil, 1)
	if err != nil {
		log.Printf("Error incrementing Get User API count: %v", err)
	}

	if c.Request.Method != http.MethodGet {
		c.Header("Allow", "GET")
		c.Status(http.StatusMethodNotAllowed)
		return
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

	startdb := time.Now()

	var user User
	if err := db.Where("email = ?", email).First(&user).Error; err != nil {
		c.Status(http.StatusUnauthorized)
		return
	}

	durationdb := time.Since(startdb).Milliseconds()

	err = observability.Client.Timing("api.response_time.GetUserAPIDataBase", time.Duration(durationdb)*time.Millisecond, nil, 1)

	if err != nil {
		log.Printf("Error recording Get User API DataBase timing: %v", err)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		c.Status(http.StatusUnauthorized)
		return
	}

	userInfo := struct {
		ID             uint      `json:"id"`
		FirstName      string    `json:"first_name"`
		LastName       string    `json:"last_name"`
		Email          string    `json:"email"`
		AccountCreated time.Time `json:"account_created"`
		AccountUpdated time.Time `json:"account_updated"`
	}{
		ID:             user.ID,
		FirstName:      user.FirstName,
		LastName:       user.LastName,
		Email:          user.Email,
		AccountCreated: user.AccountCreated,
		AccountUpdated: user.AccountUpdated,
	}

	duration := time.Since(start).Milliseconds()

	err = observability.Client.Timing("api.response_time.GetUserAPI", time.Duration(duration)*time.Millisecond, nil, 1)

	if err != nil {
		log.Printf("Error recording Get User API timing: %v", err)
	}

	c.JSON(http.StatusOK, userInfo)

}
