package mydb

import (
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"fmt"

	"webapp/usercrud"
)

func InitDB() (*gorm.DB, error) {
	dbPassword := os.Getenv("MYSQL_PWD")
	dbHost := os.Getenv("DB_HOST")
	fmt.Print("password is:" + dbPassword)
	fmt.Print("host is:" + dbHost)
	dsn := "csye6225:" + dbPassword + "@tcp(" + dbHost + ")/csye6225?parseTime=true"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	db.AutoMigrate(&usercrud.User{})

	return db, nil
}
