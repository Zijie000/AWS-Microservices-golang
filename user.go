package main

import (
	"time"
)

type User struct {
	ID             uint      `gorm:"primaryKey"`
	Email          string    `json:"email" binding:"required" gorm:"unique;not null"`
	Password       string    `json:"password" binding:"required" gorm:"not null"`
	FirstName      string    `json:"first_name" binding:"required" gorm:"not null"`
	LastName       string    `json:"last_name" binding:"required" gorm:"not null"`
	AccountCreated time.Time `json:"account_created" gorm:"autoCreateTime"`
	AccountUpdated time.Time `json:"account_updated" gorm:"autoUpdateTime"`
}

type UserUpdateForm struct {
	ID             uint   `json:"id"`
	Email          string `json:"email"`
	Password       string `json:"password"`
	FirstName      string `json:"first_name"`
	LastName       string `json:"last_name"`
	AccountCreated string `json:"account_created"`
	AccountUpdated string `json:"account_updated"`
}

type UserRegisterForm struct {
	Email     string `json:"email"`
	Password  string `json:"password"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}
