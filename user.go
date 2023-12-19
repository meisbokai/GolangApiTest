package main

import (
	"time"
)

type CreateUserRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
}

type User struct {
	ID        int       `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}

func NewUser(username string, email string) *User {
	// Set Singapore Timezone
	// TODO: Check where to initialize timezone. Shouldnt be in this constructor
	loc, _ := time.LoadLocation("Asia/Singapore")

	return &User{
		// ID:        rand.Intn(10000), // TODO: Use postgres incremental ID generation
		Username:  username,
		Email:     email,
		CreatedAt: time.Now().UTC().In(loc),
	}
}
