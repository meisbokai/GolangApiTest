package v1

import (
	"context"
	"errors"
	"net/http"
	"time"

	V1Domains "github.com/meisbokai/GolangApiTest/internal/http/domains/v1"
	"github.com/meisbokai/GolangApiTest/internal/util"
	"github.com/meisbokai/GolangApiTest/pkg/jwt"
)

type userUsecase struct {
	repo       V1Domains.UserRepository
	jwtService jwt.JWTService
}

func NewUserUsecase(repo V1Domains.UserRepository, jwtService jwt.JWTService) V1Domains.UserUsecase {
	return &userUsecase{
		repo:       repo,
		jwtService: jwtService,
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
	_, err = util.IsOldEmailMatchNew(user.Email, newEmail)
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

func (userUC *userUsecase) DeleteUser(ctx context.Context, id string) (outDom V1Domains.UserDomain, statusCode int, err error) {
	user, err := userUC.repo.GetUserByID(ctx, &V1Domains.UserDomain{ID: id})

	err = userUC.repo.DeleteUser(ctx, &V1Domains.UserDomain{ID: id})
	if err != nil {
		return V1Domains.UserDomain{}, http.StatusInternalServerError, err
	}

	return user, http.StatusOK, nil
}

func (userUC *userUsecase) Login(ctx context.Context, inDom *V1Domains.UserDomain) (outDom V1Domains.UserDomain, statusCode int, err error) {
	userDomain, err := userUC.repo.GetUserByEmail(ctx, inDom)
	if err != nil {
		return V1Domains.UserDomain{}, http.StatusUnauthorized, errors.New("invalid email or password") // for security purpose better use generic error message
	}

	if !util.ValidateHash(inDom.Password, userDomain.Password) {
		return V1Domains.UserDomain{}, http.StatusUnauthorized, errors.New("invalid email or password(hash)")
	}

	// RodeID 1 = Admin
	if userDomain.RoleID == 1 {
		userDomain.Token, err = userUC.jwtService.GenerateToken(userDomain.ID, userDomain.Username, true, userDomain.Email, userDomain.Password)
	} else {
		userDomain.Token, err = userUC.jwtService.GenerateToken(userDomain.ID, userDomain.Username, false, userDomain.Email, userDomain.Password)
	}

	if err != nil {
		return V1Domains.UserDomain{}, http.StatusInternalServerError, err
	}

	return userDomain, http.StatusOK, nil
}

func (userUC *userUsecase) GetUserByID(ctx context.Context, id string) (outDom V1Domains.UserDomain, statusCode int, err error) {
	user, err := userUC.repo.GetUserByID(ctx, &V1Domains.UserDomain{ID: id})
	if err != nil {
		return V1Domains.UserDomain{}, http.StatusNotFound, errors.New("id not found")
	}

	return user, http.StatusOK, nil
}
