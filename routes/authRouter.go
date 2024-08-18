package routes

import (
	_ "net/http"

	controller "github.com/acurilem/SID-UMAG-AuthService/controllers"

	"github.com/acurilem/SID-UMAG-AuthService/middleware"
	encuestaDocenteController "github.com/acurilem/encuesta-docente-backend/controller"
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
	surveyGroup := r.Group("/api/v1/encuesta-docente")
	{
		surveyGroup.GET("/status", encuestaDocenteController.GetConfirmacionEncuestasPorContestarRespuestas)
	}
}
