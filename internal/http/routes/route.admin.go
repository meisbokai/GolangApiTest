package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"

	V1PostgresRepository "github.com/meisbokai/GolangApiTest/internal/datasources/repositories"
	V1Handler "github.com/meisbokai/GolangApiTest/internal/http/handlers/v1"
	V1Usecase "github.com/meisbokai/GolangApiTest/internal/usecases/v1"
)

type adminRoutes struct {
	V1Handler V1Handler.UserHandler
	router    *gin.RouterGroup
	db        *sqlx.DB
}

func NewAdminRoute(router *gin.RouterGroup, db *sqlx.DB) *adminRoutes {
	V1UserRepository := V1PostgresRepository.NewUserRepository(db)
	V1UserUsecase := V1Usecase.NewUserUsecase(V1UserRepository)
	V1UserHandler := V1Handler.NewUserHandler(V1UserUsecase)

	return &adminRoutes{V1Handler: V1UserHandler, router: router, db: db}
}

func (r *adminRoutes) Routes() {
	// Routes V1
	V1Route := r.router.Group("/v1")
	{
		// users
		userRoute := V1Route.Group("/admin/users")
		userRoute.GET("/all", r.V1Handler.GetAllUserData)
		userRoute.GET("/email", r.V1Handler.GetUserByEmail)

	}

}
