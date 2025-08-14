package route

import (
	"github.com/gin-gonic/gin"
	httpdelivery "apiprofile/internal/delivery/http"
	"apiprofile/internal/middleware"
)

func SetupRouter(authHandler *httpdelivery.AuthHandler, userHandler *httpdelivery.UserHandler, jwtSecret string) *gin.Engine {
	r := gin.Default()

	r.POST("/auth/register", authHandler.Register)
	r.POST("/auth/login", authHandler.Login)
	r.POST("/auth/refresh", authHandler.Refresh)

	ug := r.Group("/users")
	ug.Use(middleware.JWTMiddleware(jwtSecret))
	{
		ug.POST("/", userHandler.Create)
		ug.GET("/", userHandler.GetAll)
		ug.GET("/:id", userHandler.GetByID)
		ug.PUT("/:id", userHandler.Update)
		ug.DELETE("/:id", userHandler.Delete)
	}
	return r
}
