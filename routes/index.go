package routes

import "github.com/gin-gonic/gin"

func Init(r *gin.Engine) {
	initAuthRoutes(r)
}
