package routes

import (
	controller "github.com/citiaps/SID-UMAG-AuthService/controllers"
	"github.com/citiaps/SID-UMAG-AuthService/middleware"
	"github.com/gin-gonic/gin"
)

func initAuthRoutes(r *gin.Engine) {
	authGroup := r.Group("/api/v1/auth")
	{
		authGroup.POST("/login", controller.LoginFunc)
		authGroup.POST("/refresh", controller.RefreshToken)
		authenticatedRoutes := r.Group("/api/v1/auth")
		authenticatedRoutes.Use(middleware.AuthMiddleware())
		{
			authenticatedRoutes.GET("/user", controller.User)
		}
	}
}
