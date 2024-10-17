package mydb

import (
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"fmt"

	"webapp/usercrud"
)

func InitDB() (*gorm.DB, error) {
	pwd := os.Getenv("mysql_pwd")
	fmt.Print("password is:" + pwd)
	dsn := "root:" + pwd + "@tcp(127.0.0.1:3306)/test?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	db.AutoMigrate(&usercrud.User{})

	return db, nil
}
