package usercrud

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"golang.org/x/crypto/bcrypt"

	"fmt"

	"gorm.io/gorm"
)

func RegisterUser(c *gin.Context, db *gorm.DB) {

	if c.Request.Method != http.MethodPost {
		c.Header("Allow", "POST")
		c.Status(http.StatusMethodNotAllowed)
		return
	}

	if len(c.Request.URL.Query()) > 0 {
		c.Status(http.StatusBadRequest)
		return
	}

	var user User

	var tmpUser UserRegisterForm

	if err := c.ShouldBindJSON(&tmpUser); err != nil {
		c.Status(http.StatusBadRequest)
		fmt.Print(1)
		return
	}

	if tmpUser.Email == "" {

		c.Status(http.StatusBadRequest)
		fmt.Print(2)
		return

	}

	var existingUser User
	if err := db.Where("email = ?", tmpUser.Email).First(&existingUser).Error; err == nil {
		c.Status(http.StatusBadRequest)
		fmt.Print(3)
		return
	}

	user.Email = tmpUser.Email
	user.FirstName = tmpUser.FirstName
	user.LastName = tmpUser.LastName

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(tmpUser.Password), bcrypt.DefaultCost)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	user.Password = string(hashedPassword)

	if result := db.Create(&user); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	userCreateInfo := struct {
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
		Email     string `json:"email"`
	}{
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
	}

	c.JSON(http.StatusCreated, userCreateInfo)
}
