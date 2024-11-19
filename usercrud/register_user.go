package usercrud

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"golang.org/x/crypto/bcrypt"

	"fmt"

	"os"

	"gorm.io/gorm"

	"regexp"

	"webapp/observability"

	"log"

	"time"

	"context"

	"github.com/google/uuid"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sns"
	//"github.com/aws/aws-sdk-go-v2/service/sns/types"
)

func publishMessage(client *sns.Client, topicArn, message, subject string) error {
	input := &sns.PublishInput{
		Message:  aws.String(message),
		TopicArn: aws.String(topicArn),
		Subject:  aws.String(subject),
	}

	_, err := client.Publish(context.TODO(), input)
	return err
}

func sendMessageToSNS(email string, token string) {
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("us-east-1"))
	if err != nil {
		log.Fatalf("unable to load SDK config, %v", err)
	}

	// 创建 SNS 客户端
	client := sns.NewFromConfig(cfg)

	accountId := os.Getenv("ACCOUNT_ID")

	// 发布消息到 SNS
	topicArn := "arn:aws:sns:us-east-1:" + accountId + ":email-sns-topic" // 替换为你的 Topic ARN
	subject := "Email Verify"
	messageBody := fmt.Sprintf(`{
		"email": "%s",
		"token": "%s"
	}`, email, token)

	err = publishMessage(client, topicArn, messageBody, subject)
	if err != nil {
		log.Fatalf("failed to publish message: %v", err)
	}
	fmt.Println("Message published successfully!")
}

func RegisterUser(c *gin.Context, db *gorm.DB) {

	start := time.Now()

	err := observability.Client.Incr("Register User API", nil, 1)
	if err != nil {
		log.Printf("Error incrementing Register User API count: %v", err)
	}

	//检查request格式

	if c.Request.Method != http.MethodPost {
		c.Header("Allow", "POST")
		c.Status(http.StatusMethodNotAllowed)
		return
	}

	if len(c.Request.URL.Query()) > 0 {
		c.Status(http.StatusBadRequest)
		return
	}

	//检查数据格式
	var userCache UserCache

	var registerForm UserRegisterForm

	if err := c.ShouldBindJSON(&registerForm); err != nil {
		c.Status(http.StatusBadRequest)
		//c.JSON(http.StatusBadRequest, gin.H{"error": "failed to bind json"})
		log.Printf("failed to bind json: %v", err)
		return
	}

	if registerForm.Email == "" {

		c.Status(http.StatusBadRequest)

		//c.JSON(http.StatusBadRequest, gin.H{"error": "email is empty"})
		log.Printf("email is empty")
		return

	}

	const emailRegex = `^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`

	re := regexp.MustCompile(emailRegex)

	if !re.MatchString(registerForm.Email) {
		c.Status(http.StatusBadRequest)
		//c.JSON(http.StatusBadRequest, gin.H{"error": "email format error"})
		log.Printf("email format error")
		return
	}

	//检查用户是否存在
	var existingUser User
	var existingCache UserCache

	if err := db.Where("email = ?", registerForm.Email).First(&existingUser).Error; err == nil {
		c.Status(http.StatusBadRequest)
		//c.JSON(http.StatusBadRequest, gin.H{"error": "user exist"})
		log.Printf("user exist: %v", err)
		return
	}

	if err := db.Where("email = ?", registerForm.Email).First(&existingCache).Error; err == nil {
		//c.Status(http.StatusBadRequest)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Do it again in 2 minutes!"})
		log.Printf("user cache exist: %v", err)
		return
	}

	userCache.Email = registerForm.Email
	userCache.FirstName = registerForm.FirstName
	userCache.LastName = registerForm.LastName
	userCache.Token = uuid.New().String()

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(registerForm.Password), bcrypt.DefaultCost)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	userCache.Password = string(hashedPassword)

	startdb := time.Now()

	if result := db.Create(&userCache); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	sendMessageToSNS(userCache.Email, userCache.Token) //sns

	durationdb := time.Since(startdb).Milliseconds()

	err = observability.Client.Timing("api.response_time.RegisterUserAPIDataBase", time.Duration(durationdb)*time.Millisecond, nil, 1)

	if err != nil {
		log.Printf("Error recording Register User API DataBase timing: %v", err)
	}

	userCreateInfo := struct {
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
		Email     string `json:"email"`
		Test      string `json:"test"`
	}{
		FirstName: userCache.FirstName,
		LastName:  userCache.LastName,
		Email:     userCache.Email,
		//Test:      "hello",
	}

	duration := time.Since(start).Milliseconds()

	err = observability.Client.Timing("api.response_time.RegisterUserAPI", time.Duration(duration)*time.Millisecond, nil, 1)
	if err != nil {
		log.Printf("Error recording Register User API timing: %v", err)
	}

	c.JSON(http.StatusCreated, userCreateInfo)

}
