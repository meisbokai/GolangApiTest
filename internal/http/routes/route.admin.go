package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"

	V1PostgresRepository "github.com/meisbokai/GolangApiTest/internal/datasources/repositories"
	V1Handler "github.com/meisbokai/GolangApiTest/internal/http/handlers/v1"
	V1Usecase "github.com/meisbokai/GolangApiTest/internal/usecases/v1"
	"github.com/meisbokai/GolangApiTest/pkg/jwt"
)

type adminRoutes struct {
	V1Handler      V1Handler.UserHandler
	router         *gin.RouterGroup
	db             *sqlx.DB
	authMiddleware gin.HandlerFunc
}

func NewAdminRoute(router *gin.RouterGroup, db *sqlx.DB, jwtService jwt.JWTService, authMiddleware gin.HandlerFunc) *adminRoutes {
	V1UserRepository := V1PostgresRepository.NewUserRepository(db)
	V1UserUsecase := V1Usecase.NewUserUsecase(V1UserRepository, jwtService)
	V1UserHandler := V1Handler.NewUserHandler(V1UserUsecase)

	return &adminRoutes{V1Handler: V1UserHandler, router: router, db: db, authMiddleware: authMiddleware}
}

func (r *adminRoutes) Routes() {
	// Routes V1
	V1Route := r.router.Group("/v1")
	{
		// users
		V1AdminRoute := V1Route.Group("/admin/users")
		V1AdminRoute.Use(r.authMiddleware)
		{
			V1AdminRoute.GET("/all", r.V1Handler.GetAllUserData)
			V1AdminRoute.GET("/email", r.V1Handler.GetUserByEmail)
		}
	}

}
