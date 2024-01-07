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

func (userUC *userUsecase) CreateUser(ctx context.Context, inDom *V1Domains.UserDomain) (outDom V1Domains.UserDomain, statusCode int, err error) {
	inDom.Password, err = util.GenerateHash(inDom.Password)
	if err != nil {
		return V1Domains.UserDomain{}, http.StatusInternalServerError, err
	}

	inDom.CreatedAt = time.Now().In(time.FixedZone("GMT+8", 8*60*60))

	err = userUC.repo.CreateUser(ctx, inDom)
	if err != nil {
		return V1Domains.UserDomain{}, http.StatusInternalServerError, err
	}

	outDom, err = userUC.repo.GetUserByEmail(ctx, inDom)
	if err != nil {
		return V1Domains.UserDomain{}, http.StatusInternalServerError, err
	}

	return outDom, http.StatusCreated, nil
}

func (userUC *userUsecase) GetUserByEmail(ctx context.Context, email string) (outDom V1Domains.UserDomain, statusCode int, err error) {
	user, err := userUC.repo.GetUserByEmail(ctx, &V1Domains.UserDomain{Email: email})
	if err != nil {
		return V1Domains.UserDomain{}, http.StatusNotFound, errors.New("email not found")
	}

	return user, http.StatusOK, nil
}

func (userUC *userUsecase) UpdateUserEmail(ctx context.Context, oldEmail string, newEmail string) (outDom V1Domains.UserDomain, statusCode int, err error) {
	user, err := userUC.repo.GetUserByEmail(ctx, &V1Domains.UserDomain{Email: oldEmail})
	if err != nil {
		return V1Domains.UserDomain{}, http.StatusNotFound, errors.New("email not found")
	}

	// Check if new email is valid
	_, err = util.ValidateEmail(newEmail)
	if err != nil {
		return V1Domains.UserDomain{}, http.StatusBadRequest, err
	}

	// Check if new email is same as old
	_, err = util.IsOldEmailMatch(user.Email, newEmail)
	if err != nil {
		return V1Domains.UserDomain{}, http.StatusBadRequest, err
	}

	err = userUC.repo.UpdateUserEmail(ctx, &V1Domains.UserDomain{Email: oldEmail}, newEmail)
	if err != nil {
		return V1Domains.UserDomain{}, http.StatusInternalServerError, err
	}

	user, err = userUC.repo.GetUserByEmail(ctx, &V1Domains.UserDomain{Email: newEmail})
	if err != nil {
		return V1Domains.UserDomain{}, http.StatusNotFound, errors.New("New email not found")
	}

	return user, http.StatusOK, nil
}

func (userUC *userUsecase) DeleteUser(ctx context.Context, email string) (outDom V1Domains.UserDomain, statusCode int, err error) {
	user, err := userUC.repo.GetUserByEmail(ctx, &V1Domains.UserDomain{Email: email})

	err = userUC.repo.DeleteUser(ctx, &V1Domains.UserDomain{Email: email})
	if err != nil {
		return V1Domains.UserDomain{}, http.StatusInternalServerError, err
	}

	return user, http.StatusOK, nil
}
