package v1

import (
	"context"
	"time"
)

type UserDomain struct {
	ID        string
	RoleID    int
	Username  string
	Email     string
	Password  string
	CreatedAt time.Time
	UpdatedAt *time.Time
	Token     string
}

type UserUsecase interface {
	GetAllUsers(ctx context.Context) (outDom []UserDomain, statusCode int, err error)
}

type UserRepository interface {
	GetAllUsers(ctx context.Context) (outDoms []UserDomain, err error)
}
