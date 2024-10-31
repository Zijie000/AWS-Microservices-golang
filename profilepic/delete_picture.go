package profilepic

import (
	"context"
	"fmt"
	"net/http"
	"webapp/usercrud"

	"gorm.io/gorm"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/gin-gonic/gin"

	"webapp/observability"

	"log"

	"golang.org/x/crypto/bcrypt"

	"time"
)

func DeleteImage(c *gin.Context, db *gorm.DB) {

	start := time.Now()

	erro := observability.Client.Incr("Delete Image API", nil, 1)
	if erro != nil {
		log.Printf("Error incrementing Delete Image API count: %v", erro)
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

	var picture Picture
	if err := db.Where("id = ?", user.ID).First(&picture).Error; err != nil {
		c.Status(http.StatusNotFound)
		return
	}

	starts3 := time.Now()

	_, err := s3Client.DeleteObject(context.TODO(), &s3.DeleteObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(fmt.Sprintf("%d", picture.ID)),
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Deleted picture fail: " + err.Error()})
		return
	}

	durations3 := time.Since(starts3).Milliseconds()

	err = observability.Client.Timing("api.response_time.DeleteImageAPIS3", time.Duration(durations3)*time.Millisecond, nil, 1)

	if err != nil {
		log.Printf("Error recording Delete Image API S3 bucket timing: %v", err)
	}

	if err := db.Delete(&picture).Error; err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	duration := time.Since(start).Milliseconds()

	err = observability.Client.Timing("api.response_time.DeleteImageAPI", time.Duration(duration)*time.Millisecond, nil, 1)

	if err != nil {
		log.Printf("Error recording Delete Image API timing: %v", err)
	}

	c.Status(http.StatusNoContent)
}
