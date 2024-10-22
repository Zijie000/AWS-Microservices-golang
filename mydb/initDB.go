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
	dbPort := os.Getenv("DB_PORT")
	fmt.Print("password is:" + dbPassword)
	dsn := "mysql://csye6225:" + dbPassword + "@" + dbHost + ":" + dbPort + "/csye6225"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	db.AutoMigrate(&usercrud.User{})

	return db, nil
}
