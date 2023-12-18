package main

import (
	"math/rand"
	"time"
)

type User struct {
	ID         int       `json:"id"`
	Username   string    `json:"username"`
	Email      string    `json:"email"`
	CreateDate time.Time `json:"created_at"`
}

func NewUser(username string, email string) *User {
	return &User{
		ID:         rand.Intn(10000),
		Username:   username,
		Email:      email,
		CreateDate: time.Now(),
	}
}
