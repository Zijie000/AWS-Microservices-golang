package usercrud

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"golang.org/x/crypto/bcrypt"

	"fmt"

	"gorm.io/gorm"

	"webapp/observability"

	"log"

	"time"
)

func UpdateUser(c *gin.Context, db *gorm.DB) {

	start := time.Now()

	err := observability.Client.Incr("Update User API", nil, 1)
	if err != nil {
		log.Printf("Error incrementing Update User API count: %v", err)
	}

	if len(c.Request.URL.Query()) > 0 {
		c.Status(http.StatusBadRequest)
		return
	}

	email, password, hasAuth := c.Request.BasicAuth()
	if !hasAuth {
		c.Status(http.StatusUnauthorized)
		return
	}

	var user User

	var tmpUser UserUpdateForm

	tmpUser.AccountCreated = "N/A"
	tmpUser.AccountUpdated = "N/A"
	tmpUser.Email = "N/A"
	tmpUser.ID = 0

	if err := db.Where("email = ?", email).First(&user).Error; err != nil {
		c.Status(http.StatusUnauthorized)
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		c.Status(http.StatusUnauthorized)
		return
	}

	if err := c.ShouldBindJSON(&tmpUser); err != nil {
		c.Status(http.StatusBadRequest)

		fmt.Print(1)
		fmt.Print(err.Error())

		return
	}

	if tmpUser.ID != 0 || tmpUser.Email != "N/A" || tmpUser.AccountCreated != "N/A" || tmpUser.AccountUpdated != "N/A" {
		c.Status(http.StatusBadRequest)
		fmt.Print(2)
		return
	}

	if tmpUser.FirstName == "" && tmpUser.LastName == "" && tmpUser.Password == "" {
		c.Status(http.StatusBadRequest)

		fmt.Print(3)

		return
	}

	if tmpUser.FirstName != "" {

		user.FirstName = tmpUser.FirstName

	}

	if tmpUser.LastName != "" {

		user.LastName = tmpUser.LastName

	}

	if tmpUser.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(tmpUser.Password), bcrypt.DefaultCost)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
			return
		}
		user.Password = string(hashedPassword)
	}

	startdb := time.Now()

	if err := db.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user"})
		return
	}

	durationdb := time.Since(startdb).Milliseconds()

	err = observability.Client.Timing("api.response_time.UpdateUserAPIDataBase", time.Duration(durationdb)*time.Millisecond, nil, 1)

	if err != nil {
		log.Printf("Error recording Update User API DataBase timing: %v", err)
	}

	duration := time.Since(start).Milliseconds()

	err = observability.Client.Timing("api.response_time.UpdateUserAPI", time.Duration(duration)*time.Millisecond, nil, 1)
	if err != nil {
		log.Printf("Error recording Update User API timing: %v", err)
	}

	c.Status(http.StatusNoContent)

}
