package v1

import (
	"github.com/gin-gonic/gin"

	V1Domains "github.com/meisbokai/GolangApiTest/internal/http/domains/v1"
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
