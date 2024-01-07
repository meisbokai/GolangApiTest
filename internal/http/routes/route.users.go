package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	V1PostgresRepository "github.com/meisbokai/GolangApiTest/internal/datasources/repositories"
	V1Handler "github.com/meisbokai/GolangApiTest/internal/http/handlers/v1"
	V1Usecase "github.com/meisbokai/GolangApiTest/internal/usecases/v1"
	"github.com/meisbokai/GolangApiTest/pkg/jwt"
)

type usersRoutes struct {
	V1Handler      V1Handler.UserHandler
	router         *gin.RouterGroup
	db             *sqlx.DB
	authMiddleware gin.HandlerFunc
}

func NewUsersRoute(router *gin.RouterGroup, db *sqlx.DB, jwtService jwt.JWTService, authMiddleware gin.HandlerFunc) *usersRoutes {
	V1UserRepository := V1PostgresRepository.NewUserRepository(db)
	V1UserUsecase := V1Usecase.NewUserUsecase(V1UserRepository, jwtService)
	V1UserHandler := V1Handler.NewUserHandler(V1UserUsecase)

	return &usersRoutes{V1Handler: V1UserHandler, router: router, db: db, authMiddleware: authMiddleware}
}

func (r *usersRoutes) Routes() {
	// Routes V1
	V1Route := r.router.Group("/v1")
	{
		// Authentications routes
		V1AuhtRoute := V1Route.Group("/auth")
		V1AuhtRoute.POST("/signup", r.V1Handler.CreateUser)
		V1AuhtRoute.POST("/login", r.V1Handler.Login)

		// Users
		V1UserRoute := V1Route.Group("/users")
		V1UserRoute.Use(r.authMiddleware)
		{
			V1UserRoute.GET("/self", r.V1Handler.GetUserByEmail) // TODO: Find a way to identify the sender without param/body
			V1UserRoute.PUT("/updateEmail", r.V1Handler.UpdateUserEmail)
			V1UserRoute.DELETE("/delete", r.V1Handler.DeleteUser)
		}

	}

}
