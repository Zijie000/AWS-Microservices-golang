package main

import (
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func initDB() (*gorm.DB, error) {
	pwd := os.Getenv("mysql_pwd")
	dsn := "root:" + pwd + "@tcp(127.0.0.1:3306)/test?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	db.AutoMigrate(&User{})

	return db, nil
}
