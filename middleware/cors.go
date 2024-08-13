package middleware

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func CorsMiddleware() gin.HandlerFunc {
	config := cors.DefaultConfig()

	config.AllowMethods = append(config.AllowMethods, "DELETE", "OPTIONS", "POST", "GET", "PUT")
	config.AllowHeaders = append(config.AllowHeaders, "Authorization", "Pagination-Count", "Origin", "Content-Type")
	config.ExposeHeaders = append(config.ExposeHeaders, "Pagination-Count")
	//	config.AllowOrigins = strings.Split(os.Getenv("CORS_URLS"), ",")
	config.AllowAllOrigins = true
	config.AllowCredentials = false

	return cors.New(config)
}
