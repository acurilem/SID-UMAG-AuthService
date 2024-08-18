package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/acurilem/SID-UMAG-AuthService/routes"
	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

//	@title			Servicio de autenticación UMAG
//	@version		1.0
//	@description	API Server para el serivicio de autenticación de usuarios con LDAP para el proyecto UMAG
//	@termsOfService	http://swagger.io/terms/

//	@contact.name	API Support
//	@contact.url	http://www.swagger.io/support
//	@contact.email	support@swagger.io

// lincense.name  Apache 2.0
//	@license.url	http://www.apache.org/licenses/LICENSE-2.0.html

//	@tag.name			auth
//	@tag.description	Autenticación, refresco y sesión del usuario

//	@host		localhost:8080
//	@BasePath	/api/v1/auth

//	@securityDefinitions.apikey	ApiKeyAuth
//	@in							header
//	@name						Authorization
//	@description				BearerJWTToken in Authorization Header

//	@accept		json
//	@produce	json

// @schemes	http https
func main() {
	router := gin.New()
	// Zap logger
	if _, err := os.Stat("logs"); os.IsNotExist(err) {
		err := os.Mkdir("logs", os.ModePerm)
		if err != nil {
			panic(err)
		}
	}
	// Log file
	logEncoder := zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig())
	fileCore := zapcore.NewCore(logEncoder, zapcore.AddSync(&lumberjack.Logger{
		Filename:   "logs/logs.log",
		MaxSize:    10,
		MaxBackups: 3,
		MaxAge:     7,
	}), zap.InfoLevel)
	// Log console
	consoleEncoder := zapcore.NewConsoleEncoder(zap.NewProductionEncoderConfig())
	consoleCore := zapcore.NewCore(consoleEncoder, zapcore.AddSync(os.Stdout), zap.InfoLevel)
	// Combine cores for multi-output logging
	teeCore := zapcore.NewTee(fileCore, consoleCore)
	zapLogger := zap.New(teeCore)

	router.Use(ginzap.GinzapWithConfig(zapLogger, &ginzap.Config{
		TimeFormat: time.RFC3339,
		UTC:        true,
		SkipPaths:  []string{"/swagger"},
	}))
	router.Use(ginzap.RecoveryWithZap(zapLogger, true))

	router.Use(gin.CustomRecovery(func(c *gin.Context, recovered interface{}) {
		if err, ok := recovered.(string); ok {
			c.String(http.StatusInternalServerError, fmt.Sprintf("Server Internal Error: %s", err))
		}
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": "Error en el servidor",
		})
	}))
	// Docs
	//docs.SwaggerInfo.BasePath = "/api/v1/auth"
	//docs.SwaggerInfo.Version = "v1"
	//docs.SwaggerInfo.Host = "localhost:8080"
	//router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	// CORS
	//router.Use(middleware.CorsMiddleware())
	// Init routes
	routes.Init(router)

	router.NoRoute(func(ctx *gin.Context) {
		ctx.JSON(404, gin.H{
			"error": "ss not found",
		})
	})
	// Init
	if err := router.Run(); err != nil {
		log.Fatalf("Error al iniciar servidor")
	}
}
