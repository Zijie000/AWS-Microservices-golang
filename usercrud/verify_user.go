package usercrud

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"gorm.io/gorm"
)

func VerifyUser(c *gin.Context, db *gorm.DB) {

	token := c.Param("token")

	var user User
	var existingUser User
	var userCache UserCache

	expiredHtmlContent := `
		<!DOCTYPE html>
		<html>
		<head>
			<title>Verification</title>
		</head>
		<body>
			<h1>This link is no longer valid!</h1>
			<p>Verification took more than 2 minutes or the user has been created, Please go back to the application</p>
		</body>
		</html>
	`

	if err := db.Where("token = ?", token).First(&userCache).Error; err != nil {
		c.Data(http.StatusNotFound, "text/html; charset=utf-8", []byte(expiredHtmlContent))
		return
	}

	if err := db.Where("email = ?", userCache.Email).First(&existingUser).Error; err == nil {
		c.Data(http.StatusNotFound, "text/html; charset=utf-8", []byte(expiredHtmlContent))
		return
	}

	user.Email = userCache.Email
	user.FirstName = userCache.FirstName
	user.LastName = userCache.LastName
	user.Password = userCache.Password

	if result := db.Create(&user); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	successHtmlContent := `
			<!DOCTYPE html>
			<html>
			<head>
				<title>Verification</title>
			</head>
			<body>
				<h1>Verify Success! Welcome to CSYE6225</h1>
				<p>Please go back to the application</p>
			</body>
			</html>
		`
	c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(successHtmlContent))

}
