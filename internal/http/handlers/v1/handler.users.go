package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/meisbokai/GolangApiTest/internal/constants"
	V1Domains "github.com/meisbokai/GolangApiTest/internal/http/domains/v1"
	"github.com/meisbokai/GolangApiTest/internal/http/requests"
	"github.com/meisbokai/GolangApiTest/internal/http/responses"
	"github.com/meisbokai/GolangApiTest/internal/util"
	"github.com/meisbokai/GolangApiTest/pkg/jwt"
)

type UserHandler struct {
	usecase V1Domains.UserUsecase
}

func NewUserHandler(usecase V1Domains.UserUsecase) UserHandler {
	return UserHandler{
		usecase: usecase,
	}
}

func (userHandler UserHandler) GetAllUserData(ctx *gin.Context) {
	ctxx := ctx.Request.Context()
	userDom, statusCode, err := userHandler.usecase.GetAllUsers(ctxx)
	if err != nil {
		NewErrorResponse(ctx, statusCode, err.Error())
		return
	}

	userResponse := responses.ToResponseList(userDom)

	NewSuccessResponse(ctx, statusCode, "user data fetched successfully", map[string]interface{}{
		"user": userResponse,
	})
}

func (userHandler UserHandler) CreateUser(ctx *gin.Context) {
	var userCreateRequest requests.UserCreateRequest
	if err := ctx.ShouldBindJSON(&userCreateRequest); err != nil {
		NewErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	userDomain := userCreateRequest.ToV1Domain()
	userDomainn, statusCode, err := userHandler.usecase.CreateUser(ctx.Request.Context(), userDomain)

	if err != nil {
		NewErrorResponse(ctx, statusCode, err.Error())
		return
	}

	NewSuccessResponse(ctx, statusCode, "registration user success", map[string]interface{}{
		"user": responses.FromV1Domain(userDomainn),
	})
}

func (userHandler UserHandler) GetUserByEmail(ctx *gin.Context) {
	ctxx := ctx.Request.Context()

	email := ctx.Query("email")

	userDom, statusCode, err := userHandler.usecase.GetUserByEmail(ctxx, email)
	if err != nil {
		NewErrorResponse(ctx, statusCode, err.Error())
		return
	}

	userResponse := responses.FromV1Domain(userDom)

	NewSuccessResponse(ctx, statusCode, "user data fetched successfully", map[string]interface{}{
		"user": userResponse,
	})

}

func (userHandler UserHandler) UpdateUserEmail(ctx *gin.Context) {
	// Get authenticated user from context
	userClaims := ctx.MustGet(constants.AuthenticatedClaimKey).(jwt.JwtCustomClaim)

	ctxx := ctx.Request.Context()

	var UserUpdateEmailRequest requests.UserUpdateEmailRequest
	if err := ctx.ShouldBindJSON(&UserUpdateEmailRequest); err != nil {
		NewErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	claimEmail := userClaims.Email
	oldEmail := UserUpdateEmailRequest.OldEmail
	newEmail := UserUpdateEmailRequest.NewEmail

	if _, err := util.IsOldEmailMatchClaim(oldEmail, claimEmail); err != nil {
		NewErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	userDom, statusCode, err := userHandler.usecase.UpdateUserEmail(ctxx, oldEmail, newEmail)
	if err != nil {
		NewErrorResponse(ctx, statusCode, err.Error())
		return
	}

	NewSuccessResponse(ctx, statusCode, "Update success", responses.FromV1Domain(userDom))

}

func (userHandler UserHandler) DeleteUser(ctx *gin.Context) {
	// Get authenticated user from context
	userClaims := ctx.MustGet(constants.AuthenticatedClaimKey).(jwt.JwtCustomClaim)

	ctxx := ctx.Request.Context()

	_, statusCode, err := userHandler.usecase.DeleteUser(ctxx, userClaims.UserID)
	if err != nil {
		NewErrorResponse(ctx, statusCode, err.Error())
		return
	}

	NewSuccessResponse(ctx, statusCode, "user deleted", map[string]interface{}{
		"user": userClaims.Username,
	})

}

func (userHandler UserHandler) Login(ctx *gin.Context) {
	var UserLoginRequest requests.UserLoginRequest
	if err := ctx.ShouldBindJSON(&UserLoginRequest); err != nil {
		NewErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	userDomain, statusCode, err := userHandler.usecase.Login(ctx.Request.Context(), UserLoginRequest.ToV1Domain())
	if err != nil {
		NewErrorResponse(ctx, statusCode, err.Error())
		return
	}

	NewSuccessResponse(ctx, statusCode, "login success", responses.FromV1Domain(userDomain))
}
