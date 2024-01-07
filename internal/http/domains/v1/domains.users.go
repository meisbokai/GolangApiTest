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
	CreateUser(ctx context.Context, inDom *UserDomain) (outDom UserDomain, statusCode int, err error)
	GetUserByEmail(ctx context.Context, email string) (outDom UserDomain, statusCode int, err error)
	UpdateUserEmail(ctx context.Context, oldEmail string, newEmail string) (outDom UserDomain, statusCode int, err error)
	DeleteUser(ctx context.Context, id string) (outDom UserDomain, statusCode int, err error)
	Login(ctx context.Context, inDom *UserDomain) (outDom UserDomain, statusCode int, err error)
	GetUserByID(ctx context.Context, id string) (outDom UserDomain, statusCode int, err error)
}

type UserRepository interface {
	GetAllUsers(ctx context.Context) (outDoms []UserDomain, err error)
	CreateUser(ctx context.Context, inDom *UserDomain) (err error)
	GetUserByEmail(ctx context.Context, inDom *UserDomain) (outDom UserDomain, err error)
	UpdateUserEmail(ctx context.Context, inDom *UserDomain, newEmail string) (err error)
	DeleteUser(ctx context.Context, inDom *UserDomain) (err error)
	GetUserByID(ctx context.Context, inDom *UserDomain) (outDom UserDomain, err error)
}
