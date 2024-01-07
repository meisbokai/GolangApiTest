package middlewares

import (
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/meisbokai/GolangApiTest/internal/constants"
	V1Handler "github.com/meisbokai/GolangApiTest/internal/http/handlers/v1"
	"github.com/meisbokai/GolangApiTest/pkg/jwt"
	"github.com/meisbokai/GolangApiTest/pkg/logger"
	"github.com/sirupsen/logrus"
)

type AuthMiddleware struct {
	jwtService jwt.JWTService
	isAdmin    bool
}

func NewAuthMiddleware(jwtService jwt.JWTService, isAdmin bool) gin.HandlerFunc {
	return (&AuthMiddleware{
		jwtService: jwtService,
		isAdmin:    isAdmin,
	}).Handle
}

func (m *AuthMiddleware) Handle(ctx *gin.Context) {
	authHeader := ctx.GetHeader("Authorization")
	if authHeader == "" {
		V1Handler.NewAbortResponse(ctx, "missing authorization header")
		return
	}

	headerParts := strings.Split(authHeader, " ")
	if len(headerParts) != 2 {
		V1Handler.NewAbortResponse(ctx, "invalid header format")
		return
	}

	if headerParts[0] != "jwt" {
		V1Handler.NewAbortResponse(ctx, "token must contain 'jtw'")
		return
	}

	claim, err := m.jwtService.ParseToken(headerParts[1])
	if err != nil {
		logger.ErrorF(fmt.Sprintf("parse token error: %s [end]", err), logrus.Fields{constants.LoggerCategory: constants.LoggerCategory})

		V1Handler.NewAbortResponse(ctx, "invalid token")
		return
	}

	if claim.IsAdmin != m.isAdmin && !claim.IsAdmin {
		V1Handler.NewAbortResponse(ctx, "you don't have access for this action")
		return
	}

	ctx.Set(constants.AuthenticatedClaimKey, claim)
	ctx.Next()
}
