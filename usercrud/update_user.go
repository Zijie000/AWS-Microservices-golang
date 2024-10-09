package usercrud

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"golang.org/x/crypto/bcrypt"

	"fmt"

	"gorm.io/gorm"
)

func UpdateUser(c *gin.Context, db *gorm.DB) {

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

	if err := db.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user"})
		return
	}

	c.Status(http.StatusNoContent)
}
