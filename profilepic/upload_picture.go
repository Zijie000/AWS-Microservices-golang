package profilepic

import (
	"context"
	"fmt"
	"net/http"
	"time"
	"webapp/usercrud"

	"webapp/observability"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"golang.org/x/crypto/bcrypt"

	"log"
)

func UploadImage(c *gin.Context, db *gorm.DB) {

	start := time.Now()

	err := observability.Client.Incr("Upload Image API", nil, 1)
	if err != nil {
		log.Printf("Error incrementing Upload Image API count: %v", err)
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

	var user usercrud.User
	if err := db.Where("email = ?", email).First(&user).Error; err != nil {
		c.Status(http.StatusUnauthorized)
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		c.Status(http.StatusUnauthorized)
		return
	}

	var existingPicture Picture
	if err := db.Where("id = ?", user.ID).First(&existingPicture).Error; err == nil {
		c.Status(http.StatusBadRequest)
		fmt.Print("existingPicture")
		return
	}

	var picture Picture

	file, header, err := c.Request.FormFile("image")

	fileName := header.Filename

	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	if header.Size == 0 {
		c.Status(http.StatusBadRequest)
		return
	}

	defer file.Close()

	today := time.Now()
	dateString := today.Format("2006-01-02")

	picture.FileName = fileName
	picture.ID = user.ID
	picture.USERID = user.ID
	picture.URL = fmt.Sprintf("%s/%d/%s", bucketName, user.ID, fileName)
	picture.UploadDate = dateString

	//imageID := fmt.Sprintf("%d-%s", time.Now().Unix(), header.Filename)

	starts3 := time.Now()

	_, err = s3Client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(fmt.Sprintf("%d", user.ID)),
		Body:   file,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Upload picture fail: " + err.Error()})
		return
	}

	durations3 := time.Since(starts3).Milliseconds()

	err = observability.Client.Timing("api.response_time.UploadImageAPIS3", time.Duration(durations3)*time.Millisecond, nil, 1)

	if err != nil {
		log.Printf("Error recording Upload Image API S3 bucket timing: %v", err)
	}

	if result := db.Create(&picture); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Create picture data fail"})
		return
	}

	pictureCreateInfo := struct {
		FileName   string `json:"file_name"`
		ID         string `json:"id"`
		USERID     string `json:"user_id"`
		URL        string `json:"url"`
		UploadDate string `json:"upload_date"`
	}{
		FileName:   fileName,
		ID:         fmt.Sprintf("%d", user.ID),
		USERID:     fmt.Sprintf("%d", user.ID),
		URL:        fmt.Sprintf("%s/%d/%s", bucketName, user.ID, fileName),
		UploadDate: dateString,
	}

	duration := time.Since(start).Milliseconds()

	err = observability.Client.Timing("api.response_time.UploadImageAPI", time.Duration(duration)*time.Millisecond, nil, 1)

	if err != nil {
		log.Printf("Error recording Upload Image API timing: %v", err)
	}

	c.JSON(http.StatusCreated, pictureCreateInfo)
}
