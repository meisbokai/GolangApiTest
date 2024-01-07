package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"

	V1Domains "github.com/meisbokai/GolangApiTest/internal/http/domains/v1"
	"github.com/meisbokai/GolangApiTest/internal/http/requests"
	"github.com/meisbokai/GolangApiTest/internal/http/responses"
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
