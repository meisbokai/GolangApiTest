package v1

import (
	"context"
	"errors"
	"net/http"
	"time"

	V1Domains "github.com/meisbokai/GolangApiTest/internal/http/domains/v1"
	"github.com/meisbokai/GolangApiTest/internal/util"
)

type userUsecase struct {
	repo V1Domains.UserRepository
}

func NewUserUsecase(repo V1Domains.UserRepository) V1Domains.UserUsecase {
	return &userUsecase{
		repo: repo,
	}
}

func (userUC *userUsecase) GetAllUsers(ctx context.Context) (outDom []V1Domains.UserDomain, statusCode int, err error) {
	users, err := userUC.repo.GetAllUsers(ctx)
	if err != nil {
		return []V1Domains.UserDomain{}, http.StatusNotFound, errors.New("Unable to get full list of users")
	}

	return users, http.StatusOK, nil
}
func (userUC *userUsecase) GetUserByEmail(ctx context.Context, email string) (outDom V1Domains.UserDomain, statusCode int, err error) {
	user, err := userUC.repo.GetUserByEmail(ctx, &V1Domains.UserDomain{Email: email})
	if err != nil {
		return V1Domains.UserDomain{}, http.StatusNotFound, errors.New("email not found")
	}

	return user, http.StatusOK, nil
}
