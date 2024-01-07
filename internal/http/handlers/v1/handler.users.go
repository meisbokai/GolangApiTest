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
	"github.com/meisbokai/GolangApiTest/pkg/validators"
)

type UserHandler struct {
	usecase V1Domains.UserUsecase
}

func NewUserHandler(usecase V1Domains.UserUsecase) UserHandler {
	return UserHandler{
		usecase: usecase,
	}
}

// GetAllUserData godoc
// @Summary Get all user data
// @Description Get all user data
// @Tags Admin
// @Accept json
// @Produce json
// @Success 200 {array} responses.UserResponse "User data"
// @Failure 401 {object} object{message=string,status=bool} "Unauthorized"
// @Failure 500 {object} object{message=string,status=bool} "Internal Server Error"
// @Router /v1/admin/users/all [get]
// @Security jwtToken
// @Param Authorization header string true "Insert your access token" default(jwt <Add access token here>)
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

// CreateUser godoc
// @Summary Create a user
// @Description Create a new user
// @Tags Auth
// @Accept json
// @Produce json
// @Param user body requests.UserCreateRequest true "User data"
// @Success 201 {object} responses.UserResponse "User data"
// @Failure 400 {object} object{message=string,status=bool}  "Bad Request"
// @Failure 500 {object} object{message=string,status=bool}  "Internal Server Error"
// @Router /v1/auth/signup [post]
func (userHandler UserHandler) CreateUser(ctx *gin.Context) {
	var userCreateRequest requests.UserCreateRequest
	if err := ctx.ShouldBindJSON(&userCreateRequest); err != nil {
		NewErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	if err := validators.ValidatePayloads(userCreateRequest); err != nil {
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

// GetUserByEmail godoc
// @Summary Get user by email
// @Description Get user by email
// @Tags Admin
// @Accept json
// @Produce json
// @Param email query string true "email"
// @Success 200 {object} responses.UserResponse "User data"
// @Failure 404 {object} object{message=string,status=bool}  "Not Found"
// @Failure 500 {object} object{message=string,status=bool}  "Internal Server Error"
// @Router /v1/admin/users/email [get]
// @Security jwtToken
// @Param Authorization header string true "Insert your access token" default(jwt <Add access token here>)
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

// UpdateUserEmail godoc
// @Summary Update user email
// @Description Update user email based on the authenticated user
// @Tags Users
// @Accept json
// @Produce json
// @Param body body requests.UserUpdateEmailRequest true "Update email request"
// @Success 200 {object} responses.UserResponse "User data"
// @Failure 400 {object} object{message=string,status=bool}  "Bad Request"
// @Failure 401 {object} object{message=string,status=bool}  "Unauthorized"
// @Failure 500 {object} object{message=string,status=bool}  "Internal Server Error"
// @Router /v1/users/updateEmail [put]
// @Security jwtToken
// @Param Authorization header string true "Insert your access token" default(jwt <Add access token here>)
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

	if err := validators.ValidatePayloads(UserUpdateEmailRequest); err != nil {
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

// DeleteUser godoc
// @Summary Delete a user
// @Description Delete a user based on the authenticated user
// @Tags Users
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{} "User data"
// @Failure 401 {object} object{message=string,status=bool}  "Unauthorized"
// @Failure 500 {object} object{message=string,status=bool}  "Internal Server Error"
// @Router /v1/users/delete [delete]
// @Security jwtToken
// @Param Authorization header string true "Insert your access token" default(jwt <Add access token here>)
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

// Login godoc
// @Summary Login a user
// @Description Login a user
// @Tags Auth
// @Accept json
// @Produce json
// @Param user body requests.UserLoginRequest true "User data"
// @Success 200 {object} responses.UserResponse "User data"
// @Failure 400 {object} object{message=string,status=bool}  "Bad Request"
// @Failure 500 {object} object{message=string,status=bool}  "Internal Server Error"
// @Router /v1/auth/login [post]
func (userHandler UserHandler) Login(ctx *gin.Context) {
	var UserLoginRequest requests.UserLoginRequest
	if err := ctx.ShouldBindJSON(&UserLoginRequest); err != nil {
		NewErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	if err := validators.ValidatePayloads(UserLoginRequest); err != nil {
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

// GetSelfUser godoc
// @Summary Get user data
// @Description Get user data based on the authenticated user
// @Tags Users
// @Accept json
// @Produce json
// @Success 200 {object} responses.UserResponse "User data"
// @Failure 401 {object} object{message=string,status=bool}  "Unauthorized"
// @Failure 500 {object} object{message=string,status=bool}  "Internal Server Error"
// @Router /v1/users/self [get]
// @Security jwtTokenring
// @Param Authorization header string true "Insert your access token" default(jwt <Add access token here>)
func (userHandler UserHandler) GetSelfUser(ctx *gin.Context) {
	// get authenticated user from context
	userClaims := ctx.MustGet(constants.AuthenticatedClaimKey).(jwt.JwtCustomClaim)

	ctxx := ctx.Request.Context()
	userDom, statusCode, err := userHandler.usecase.GetUserByID(ctxx, userClaims.UserID)
	if err != nil {
		NewErrorResponse(ctx, statusCode, err.Error())
		return
	}

	userResponse := responses.FromV1Domain(userDom)

	NewSuccessResponse(ctx, statusCode, "user data fetched successfully", map[string]interface{}{
		"user": userResponse,
	})
}
